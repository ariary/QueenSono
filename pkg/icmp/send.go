package icmp

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/ariary/QueenSono/pkg/message"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	ProtocolICMP = 1
	//ProtocolIPv6ICMP = 58
)

// Mostly based on https://github.com/golang/net/blob/master/icmp/ping_test.go
// Send ICMP echo request packet (code 8) to remote and waiting for the echo reply (code 0)
func IcmpSendRaw(listeningReplyAddr string, remoteAddr string, data string) (*net.IPAddr, error) {
	// Listen is used to have a PacketConn but we won't wait for reply
	c, err := icmp.ListenPacket("ip4:icmp", listeningReplyAddr)
	if err != nil {
		return nil, err
	}

	defer c.Close()
	// Resolve any DNS (if used) and get the real IP of the target
	dst, err := net.ResolveIPAddr("ip4", remoteAddr)
	if err != nil {
		return nil, err
	}

	// Make a new ICMP message
	m := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1, //<< uint(seq), // TODO
			Data: []byte(data),
		},
	}
	b, err := m.Marshal(nil)
	if err != nil {
		return dst, err
	}

	// Send it
	n, err := c.WriteTo(b, dst)
	if err != nil {
		return dst, err
	} else if n != len(b) {
		return dst, fmt.Errorf("got %v; want %v", n, len(b))
	}

	return dst, nil
}

// Mostly based on https://github.com/golang/net/blob/master/icmp/ping_test.go
// Send ICMP echo request packet (code 8) to remote and waiting for the echo reply (code 0)
func IcmpSendAndWaitForReply(listeningReplyAddr string, remoteAddr string, data string) (*net.IPAddr, time.Duration, error) {
	// Start listening for icmp replies
	c, err := icmp.ListenPacket("ip4:icmp", listeningReplyAddr)
	if err != nil {
		return nil, 0, err
	}
	defer c.Close()

	start := time.Now()
	dst, err := IcmpSendRaw(listeningReplyAddr, remoteAddr, data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Wait for a reply
	reply := make([]byte, 65507)
	err = c.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return dst, 0, err
	}
	n, peer, err := c.ReadFrom(reply) //n, peer, err := c.ReadFrom(reply)
	if err != nil {
		return dst, 0, err
	}
	duration := time.Since(start)

	// Pack it up boys, we're done here
	rm, err := icmp.ParseMessage(ProtocolICMP, reply[:n])
	if err != nil {
		return dst, 0, err
	}

	//We do not received echo reply
	if rm.Type != ipv4.ICMPTypeEchoReply {
		return dst, 0, fmt.Errorf("got %+v from %v; want echo reply", rm, peer)
	}
	return dst, duration, nil
}

// Return a slice of a the string chunked with specific length (string length of each chunk)
//Thanks to https://stackoverflow.comProtocolICMP/questions/25686109/split-string-by-length-in-golang
func Chunks(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string = make([]string, 0, (len(s)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
		if currentLen == chunkSize {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}

//Send echo the same echo packet while we do not received an echo reply
func SendWhileNoEchoReply(listeningReplyAddr string, remoteAddr string, data string) {
	for { //while we do not received echo reply ~ACK, resend it
		dst, dur, err := IcmpSendAndWaitForReply(listeningReplyAddr, remoteAddr, data)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Retrying...")
		} else {
			fmt.Printf("PING %s (%s): %s\n", remoteAddr, dst, dur)
			break
		}
	}
}

//Send string to remote using ICMP and waiting for echo reply. You must specify the delay between each packet and the size
func SendReply(listeningReplyAddr string, remoteAddr string, chunkSize int, delay int, data string) {
	dataSlice := Chunks(data, chunkSize) //1 character = 1byte , max size of icmp data 65507
	dataSlice = message.QueenSonoMarshall(dataSlice)
	time.Sleep(time.Duration(delay) * time.Second)

	// Announce the data size
	SendWhileNoEchoReply(listeningReplyAddr, remoteAddr, strconv.Itoa(len(dataSlice)))

	//Send the data
	for i := 0; i < len(dataSlice); i++ {
		time.Sleep(time.Duration(delay) * time.Second)
		SendWhileNoEchoReply(listeningReplyAddr, remoteAddr, dataSlice[i])
	}
}

//Send string to remote using ICMP and waiting for echo reply. You must specify the delay between each packet and the size
func SendNoReply(listeningReplyAddr string, remoteAddr string, chunkSize int, delay int, data string) {
	dataSlice := Chunks(data, chunkSize) //1 character = 1byte , max size of icmp data 65507
	dataSlice = message.QueenSonoMarshall(dataSlice)
	time.Sleep(time.Duration(delay) * time.Second)

	// Announce the data size
	IcmpSendRaw(listeningReplyAddr, remoteAddr, strconv.Itoa(len(dataSlice)))

	//Send the data
	for i := 0; i < len(dataSlice); i++ {
		time.Sleep(time.Duration(delay) * time.Second)
		IcmpSendRaw(listeningReplyAddr, remoteAddr, dataSlice[i])
	}
}

func SendHashedmessage(msg string, remoteAddr string, listenAddr string) {
	fmt.Println("addr", remoteAddr)
	// hash := utils.Sha1(msg)
	// fmt.Println("hash", hash)
	dst, dur, err := IcmpSendAndWaitForReply(listenAddr, msg, remoteAddr)
	if err != nil {
		panic(err)
	}
	log.Printf("PING %s (%s): %s\n", remoteAddr, dst, dur)
}
