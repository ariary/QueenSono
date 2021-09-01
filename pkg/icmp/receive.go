package icmp

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	qsmessage "github.com/ariary/QueenSono/pkg/message"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

//Wait for the fistr ICMP packet setting sized of data
func GetMessageSizeAndSender(listenAddr string) (size int, sender string) {
	c, err := icmp.ListenPacket("ip4:icmp", listenAddr)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer c.Close()

	packet := make([]byte, 65507)
	n, peer, err := c.ReadFrom(packet)
	if err != nil {
		fmt.Println("Error while reading icmp packet:", err)
	}

	message, err := icmp.ParseMessage(ProtocolICMP, packet[:n])
	if err != nil {
		fmt.Println("Error while parsing icmp message:", err)
	}

	switch message.Type {
	case ipv4.ICMPTypeEcho:
		echo, _ := message.Body.Marshal(1)
		m := string(echo[4:]) //clean
		size, err = strconv.Atoi(m)
		if err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Errorf("got %+v from %v; want echo request", message, peer)
	}
	sender = peer.String()
	return size, sender
}

//ICMP server waiting for packet (waiting n packet)
func Serve(listenAddr string, n int, progressBar bool) (data string, missingPacketIndexes []int) {
	var bar *progressbar.ProgressBar
	if progressBar {
		bar = progressbar.Default(int64(n))
	}
	c, err := icmp.ListenPacket("ip4:icmp", listenAddr)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer c.Close()

	dataChunked := make([]string, n)
	//prepare map of indexes (map cause it  complexity to delete an element is 0(1), slice 0(n))
	indexes := make(map[int]int)
	for i := 0; i < n; i++ {
		indexes[i] = 0 //the value does not have any interest
	}
	//Retrieve the n packet (waiting till we have all packets)
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(c *icmp.PacketConn, dataChunked []string, i int, indexes map[int]int, bar *progressbar.ProgressBar) {
			defer wg.Done()
			if progressBar {
				getPacketAndBarUpdate(bar, c, dataChunked, indexes)
			} else {
				getPacket(c, dataChunked, indexes)
			}
		}(c, dataChunked, i, indexes, bar)
	}

	wg.Wait()

	data = strings.Join(dataChunked, "")

	missingPacketIndexes = make([]int, 0, len(indexes))
	for index := range indexes {
		missingPacketIndexes = append(missingPacketIndexes, index)
	}
	return data, missingPacketIndexes
}

//ICMP server waiting for specific number of packet (waiting n packet) but that stop working after (n+2)*delay seconds
// return a string with the data section of each packet concatened and the index of the packet missing
func ServeTemporary(listenAddr string, n int, progressBar bool, delay int) (data string, missingPacketIndexes []int) {
	//crossBar
	var bar *progressbar.ProgressBar
	if progressBar {
		bar = progressbar.Default(int64(n))
	}

	//ICMP Serverlistening
	c, err := icmp.ListenPacket("ip4:icmp", listenAddr)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer c.Close()

	dataChunked := make([]string, n)
	indexes := make(map[int]int) //prepare map of indexes (map cause it  complexity to delete an element is 0(1), slice 0(n))
	for i := 0; i < n; i++ {
		indexes[i] = 0 //the value does not have any interest
	}

	//Retrieve the packet
	for i := 0; i < n; i++ {
		if progressBar {
			go getPacketAndBarUpdate(bar, c, dataChunked, indexes)
		} else {
			go getPacket(c, dataChunked, indexes)
		}
	}

	//Counter for not waiting indefinitively
	counter := (n + 2) * delay //let a little offset
	time.Sleep(time.Duration(counter) * time.Second)
	data = strings.Join(dataChunked, "")

	missingPacketIndexes = make([]int, 0, len(indexes))
	for index := range indexes {
		missingPacketIndexes = append(missingPacketIndexes, index)
	}
	return data, missingPacketIndexes
}

//Get a single packet then add the data to the slice (= chunked data) and remove the index from the indexes
func getPacket(c *icmp.PacketConn, data []string, indexes map[int]int) {
	packet := make([]byte, 65507)
	n, peer, err := c.ReadFrom(packet)
	if err != nil {
		fmt.Println("Error while reading icmp packet:", err)
	}

	message, err := icmp.ParseMessage(ProtocolICMP, packet[:n])
	if err != nil {
		fmt.Println("Error while parsing icmp message:", err)
	}

	switch message.Type {
	case ipv4.ICMPTypeEcho:
		echo, _ := message.Body.Marshal(1)
		//m := string(echo[2:]) //clean
		msg, index := qsmessage.QueenSonoUnmarshall(string(echo[4:])) //echo has 4 bytes which are added, don't know why
		data[index] = msg
		delete(indexes, index)
	default:
		fmt.Errorf("got %+v from %v; want echo request", message, peer)
	}
}

//Get a single packet then add the data to the slice (= chunked data), remove the index from the indexes & update the crossbar
func getPacketAndBarUpdate(bar *progressbar.ProgressBar, c *icmp.PacketConn, data []string, indexes map[int]int) {
	packet := make([]byte, 65507)
	n, peer, err := c.ReadFrom(packet)
	if err != nil {
		fmt.Println("Error while reading icmp packet:", err)
	}

	message, err := icmp.ParseMessage(ProtocolICMP, packet[:n])
	if err != nil {
		fmt.Println("Error while parsing icmp message:", err)
	}

	switch message.Type {
	case ipv4.ICMPTypeEcho:
		echo, _ := message.Body.Marshal(1)
		//m := string(echo[2:]) //clean
		msg, index := qsmessage.QueenSonoUnmarshall(string(echo[4:])) //echo has 4 bytes which are added, don't know why
		data[index] = msg
		delete(indexes, index)
		bar.Add(1)
	default:
		fmt.Errorf("got %+v from %v; want echo request", message, peer)
	}
}

//Wait ICMP message from remote to assert if the message is well received
func IntegrityCheck(hash string) {
	fmt.Println("launch integrity server")
	c, err := icmp.ListenPacket("ip4:icmp", "localhost")
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer c.Close()

	for {
		packet := make([]byte, 65507)
		n, peer, err := c.ReadFrom(packet)
		if err != nil {
			fmt.Println("Error while reading icmp packet:", err)
		}

		message, err := icmp.ParseMessage(ProtocolICMP, packet[:n])
		if err != nil {
			fmt.Println("Error while parsing icmp message:", err)
		}

		switch message.Type {
		case ipv4.ICMPTypeEcho:
			fmt.Println("Get integrity")
			echo, _ := message.Body.Marshal(1)
			fmt.Println("hash received:", string(echo[4:])) //clean
			if string(echo[4:]) == hash {
				fmt.Println("Communication end")
				os.Exit(0)
			}
		default:
			//fmt.Println("DEFAULT!!!!!!!!!!")
			fmt.Errorf("got %+v from %v; want echo request", message, peer)
		}
	}
}
