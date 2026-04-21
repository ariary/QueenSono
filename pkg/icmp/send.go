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
)

// IcmpSendRaw sends a single ICMP echo request without waiting for a reply.
func IcmpSendRaw(listeningReplyAddr string, remoteAddr string, data string) (*net.IPAddr, error) {
	c, err := icmp.ListenPacket("ip4:icmp", listeningReplyAddr)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	dst, err := net.ResolveIPAddr("ip4", remoteAddr)
	if err != nil {
		return nil, err
	}

	m := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte(data),
		},
	}
	b, err := m.Marshal(nil)
	if err != nil {
		return dst, err
	}
	n, err := c.WriteTo(b, dst)
	if err != nil {
		return dst, err
	} else if n != len(b) {
		return dst, fmt.Errorf("wrote %v bytes; want %v", n, len(b))
	}
	return dst, nil
}

// IcmpSendAndWaitForReply sends an ICMP echo request and waits for the echo reply.
// Both send and receive share a single PacketConn so the reply is guaranteed to
// arrive on the same socket that issued the request.
func IcmpSendAndWaitForReply(listeningReplyAddr string, remoteAddr string, data string) (*net.IPAddr, time.Duration, error) {
	c, err := icmp.ListenPacket("ip4:icmp", listeningReplyAddr)
	if err != nil {
		return nil, 0, err
	}
	defer c.Close()

	dst, err := net.ResolveIPAddr("ip4", remoteAddr)
	if err != nil {
		return nil, 0, err
	}

	m := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte(data),
		},
	}
	b, err := m.Marshal(nil)
	if err != nil {
		return dst, 0, err
	}

	start := time.Now()
	n, err := c.WriteTo(b, dst)
	if err != nil {
		return dst, 0, err
	} else if n != len(b) {
		return dst, 0, fmt.Errorf("wrote %v bytes; want %v", n, len(b))
	}

	reply := make([]byte, 65535)
	if err = c.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return dst, 0, err
	}
	n, peer, err := c.ReadFrom(reply)
	if err != nil {
		return dst, 0, err
	}
	duration := time.Since(start)

	rm, err := icmp.ParseMessage(ProtocolICMP, reply[:n])
	if err != nil {
		return dst, 0, err
	}
	if rm.Type != ipv4.ICMPTypeEchoReply {
		return dst, 0, fmt.Errorf("got %+v from %v; want echo reply", rm, peer)
	}
	return dst, duration, nil
}

// Chunks splits s into substrings of at most chunkSize bytes.
// Splitting is byte-based; multi-byte UTF-8 sequences may be split mid-character
// (correct for binary data).
func Chunks(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	chunks := make([]string, 0, (len(s)-1)/chunkSize+1)
	for i := 0; i < len(s); i += chunkSize {
		chunks = append(chunks, s[i:min(i+chunkSize, len(s))])
	}
	return chunks
}

// SendWhileNoEchoReply retransmits data until an echo reply is received.
func SendWhileNoEchoReply(listeningReplyAddr string, remoteAddr string, data string) {
	for {
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

// SendReply sends data chunked over ICMP, waiting for echo reply after each chunk.
func SendReply(listeningReplyAddr string, remoteAddr string, chunkSize int, delay int, data string) {
	dataSlice := Chunks(data, chunkSize)
	dataSlice = message.QueenSonoMarshall(dataSlice)
	time.Sleep(time.Duration(delay) * time.Second)

	SendWhileNoEchoReply(listeningReplyAddr, remoteAddr, strconv.Itoa(len(dataSlice)))
	for i := range len(dataSlice) {
		time.Sleep(time.Duration(delay) * time.Second)
		SendWhileNoEchoReply(listeningReplyAddr, remoteAddr, dataSlice[i])
	}
}

// SendNoReply sends data chunked over ICMP without waiting for echo replies.
func SendNoReply(listeningReplyAddr string, remoteAddr string, chunkSize int, delay int, data string) {
	dataSlice := Chunks(data, chunkSize)
	dataSlice = message.QueenSonoMarshall(dataSlice)
	time.Sleep(time.Duration(delay) * time.Second)

	if _, err := IcmpSendRaw(listeningReplyAddr, remoteAddr, strconv.Itoa(len(dataSlice))); err != nil {
		fmt.Fprintf(os.Stderr, "send size announcement: %v\n", err)
	}
	for i := range len(dataSlice) {
		time.Sleep(time.Duration(delay) * time.Second)
		if _, err := IcmpSendRaw(listeningReplyAddr, remoteAddr, dataSlice[i]); err != nil {
			fmt.Fprintf(os.Stderr, "send chunk %d: %v\n", i, err)
		}
	}
}

// SendHashedmessage sends a hash to remoteAddr for integrity verification.
func SendHashedmessage(msg string, remoteAddr string, listenAddr string) {
	fmt.Println("addr", remoteAddr)
	dst, dur, err := IcmpSendAndWaitForReply(listenAddr, remoteAddr, msg)
	if err != nil {
		log.Fatalf("SendHashedmessage: %v", err)
	}
	log.Printf("PING %s (%s): %s\n", remoteAddr, dst, dur)
}
