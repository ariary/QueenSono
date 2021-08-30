package main

import (
	"fmt"

	"github.com/ariary/QueenSono/pkg/icmp"
)

func main() {
	listenAddr := "10.0.2.15"
	size, sender := icmp.GetMessageSizeAndSender(listenAddr)
	fmt.Println("Sender:", sender, ", Number of packet wanted:", size)
	message := icmp.Serve(listenAddr, size)
	fmt.Println("Message received:", message)
	//icmp.SendHashedmessage(message, sender) //Integrity check
}
