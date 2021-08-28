package icmp

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	// Stolen from https://godoc.org/golang.org/x/net/internal/iana,
	// can't import "internal" packages
	ProtocolICMP = 1
	//ProtocolIPv6ICMP = 58
)

// Default to listen on all IPv4 interfaces
var ListenAddr = "0.0.0.0"

// Mostly based on https://github.com/golang/net/blob/master/icmp/ping_test.go
// All ye beware, there be dragons below...

func IcmpSendRaw(data string, addr string) (*net.IPAddr, time.Duration, error) {
	// Start listening for icmp replies
	c, err := icmp.ListenPacket("ip4:icmp", ListenAddr)
	if err != nil {
		return nil, 0, err
	}
	defer c.Close()

	// Resolve any DNS (if used) and get the real IP of the target
	dst, err := net.ResolveIPAddr("ip4", addr)
	if err != nil {
		panic(err)
		return nil, 0, err
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
		return dst, 0, err
	}

	// Send it
	start := time.Now()
	n, err := c.WriteTo(b, dst)
	if err != nil {
		return dst, 0, err
	} else if n != len(b) {
		return dst, 0, fmt.Errorf("got %v; want %v", n, len(b))
	}

	// Wait for a reply
	reply := make([]byte, 1500)
	err = c.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return dst, 0, err
	}
	n, peer, err := c.ReadFrom(reply)
	if err != nil {
		return dst, 0, err
	}
	duration := time.Since(start)

	// Pack it up boys, we're done here
	rm, err := icmp.ParseMessage(ProtocolICMP, reply[:n])
	if err != nil {
		return dst, 0, err
	}
	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		echo, _ := rm.Body.Marshal(1)
		fmt.Println(string(echo))
		return dst, duration, nil
	default:
		return dst, 0, fmt.Errorf("got %+v from %v; want echo reply", rm, peer)
	}
}
