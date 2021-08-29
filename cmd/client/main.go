package main

import (
	"log"
	"strconv"

	"github.com/ariary/QueenSono/pkg/icmp"
	"github.com/ariary/QueenSono/pkg/utils"
)

func main() {
	addr := "192.168.1.39"
	data := "totoc,titi,tata,tutu,tyty,deuxfois,totoc,titi,tata,tutu,tyty,deuxfois"
	hash := utils.Sha1(data)
	go icmp.IntegrityCheck(hash)
	dataSlice := icmp.Chunks(data, 6)
	// Announce the data size
	dst, dur, err := icmp.IcmpSendRaw(strconv.Itoa(len(dataSlice)), addr)
	if err != nil {
		panic(err)
	}
	log.Printf("Ping %s (%s): %s\n", addr, dst, dur)

	//Send the data
	for i := 0; i < len(dataSlice); i++ {
		dst, dur, err := icmp.IcmpSendRaw(dataSlice[i], addr)
		if err != nil {
			panic(err)
		}
		log.Printf("Ping %s (%s): %s\n", addr, dst, dur)
	}
}
