package main

import (
	"log"

	"github.com/ariary/QueenSono/pkg/icmp"
)

func main() {
	addr := "google.com"
	data := "toto"
	dst, dur, err := icmp.IcmpSendRaw(data, addr)
	if err != nil {
		panic(err)
	}
	log.Printf("Ping %s (%s): %s\n", addr, dst, dur)
}
