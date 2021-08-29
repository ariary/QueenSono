package main

import (
	"fmt"

	"github.com/ariary/QueenSono/pkg/icmp"
)

func main() {
	size, sender := icmp.GetMessageSizeAndSender()
	fmt.Println("Sender:", sender, ", size:", size)
	message := icmp.Serve(size)
	fmt.Println("Message received:", message)
	icmp.SendHashedmessage(message, sender)
}
