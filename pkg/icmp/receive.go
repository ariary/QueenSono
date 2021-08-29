package icmp

import (
	"fmt"
	"strconv"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

//Wait for the fistr ICMP packet setting sized of data
func GetMessageSizeAndSender() (size int, sender string) {
	c, err := icmp.ListenPacket("ip4:icmp", "192.168.1.39")
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
func Serve(n int) (data string) {
	c, err := icmp.ListenPacket("ip4:icmp", "192.168.1.39")
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
			fmt.Println(m)
			data += m
		default:
			fmt.Errorf("got %+v from %v; want echo request", message, peer)
		}
	}
	return data
}
