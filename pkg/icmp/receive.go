package icmp

import (
	"fmt"
	"os"
	"strconv"

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
func Serve(listenAddr string, n int, progressBar bool) (data string) {
	var bar *progressbar.ProgressBar
	if progressBar {
		bar = progressbar.Default(int64(n))
	}
	c, err := icmp.ListenPacket("ip4:icmp", listenAddr)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer c.Close()
	for i := 0; i < n; i++ {
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
			m := string(echo[2:]) //clean
			if progressBar {
				bar.Add(1)
			} else {
				fmt.Println(m)
			}
			data += m
		default:
			fmt.Errorf("got %+v from %v; want echo request", message, peer)
		}

	}
	return data
}

//ICMP server waiting for specific number of packet (waiting n packet)
// func ServeTemporary(listenAddr string, n int, data chan string) {
// 	c, err := icmp.ListenPacket("ip4:icmp", listenAddr)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 	}
// 	defer c.Close()
// 	for i := 0; i < n; i++ {
// 		packet := make([]byte, 65507)
// 		n, peer, err := c.ReadFrom(packet)
// 		if err != nil {
// 			fmt.Println("Error while reading icmp packet:", err)
// 		}

// 		message, err := icmp.ParseMessage(ProtocolICMP, packet[:n])
// 		if err != nil {
// 			fmt.Println("Error while parsing icmp message:", err)
// 		}

// 		switch message.Type {
// 		case ipv4.ICMPTypeEcho:
// 			echo, _ := message.Body.Marshal(1)
// 			m := string(echo[2:]) //clean
// 			fmt.Println(m)
// 			data <- m
// 		default:
// 			fmt.Errorf("got %+v from %v; want echo request", message, peer)
// 		}
// 	}
// 	return data
// }

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
			fmt.Println("DEFAULT!!!!!!!!!!")
			fmt.Errorf("got %+v from %v; want echo request", message, peer)
		}
	}
}
