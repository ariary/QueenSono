// package main

// import (
// 	"fmt"

// 	"github.com/ariary/QueenSono/pkg/icmp"
// )

// func main() {
// 	listenAddr := "10.0.2.15"
// 	size, sender := icmp.GetMessageSizeAndSender(listenAddr)
// 	fmt.Println("Sender:", sender, ", Number of packet wanted:", size)
// 	message := icmp.Serve(listenAddr, size)
// 	fmt.Println("Message received:", message)
// 	//icmp.SendHashedmessage(message, sender) //Integrity check
// }
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ariary/QueenSono/pkg/icmp"
	"github.com/spf13/cobra"
)

func main() {
	//TODO truncated,integrity,crossbar,
	//CMD RECEIVE
	//receive var
	var listenAddr string
	var filename string
	var progressBar bool

	var cmdReceive = &cobra.Command{
		Use:   "receive",
		Short: "receive data from icmp packet",
		Long: `receive is for receiving  data from a remote queensono sender.
it uses the icmp protocol.`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			size, sender := icmp.GetMessageSizeAndSender(listenAddr)
			fmt.Println("Sender:", sender, ", Number of packet wanted:", size)
			message := icmp.Serve(listenAddr, size, progressBar)
			//icmp.SendHashedmessage(message, sender) //Integrity check
			if filename != "" {
				f, err := os.Create(filename)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()

				_, err2 := f.WriteString(message)
				if err2 != nil {
					log.Fatal(err2)
				}
				fmt.Println("data saved in", filename)
			} else {
				fmt.Println("qsreceiver received:", message)
			}
		},
	}

	var cmdTruncated = &cobra.Command{
		Use:   "truncated [delay]",
		Short: "receive data from icmp packet and do not wait indefinitively for all the packet.",
		Long:  `rreceive data from icmp packet and do not wait indefinitively for all the packet. Hence you could receive truncated data`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			size, sender := icmp.GetMessageSizeAndSender(listenAddr)
			fmt.Println("Sender:", sender, ", Number of packet wanted:", size)
			data := make(chan string)
			//go countdown(strconv.Atoi(args[0]),size)
			delay, err := strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}
			icmp.ServeTemporary(listenAddr, size, delay, data)
			message := <-data
			fmt.Println("Message received:", message)
			//icmp.SendHashedmessage(message, sender) //Integrity check

			if filename != "" {
				f, err := os.Create(filename)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()

				_, err2 := f.WriteString(message)
				if err2 != nil {
					log.Fatal(err2)
				}
				fmt.Println("data saved in", filename)
			} else {
				fmt.Println("qsreceiver received:", message)
			}
		},
	}

	//cmdReceive flag handling
	cmdReceive.PersistentFlags().StringVarP(&listenAddr, "listen", "l", "0.0.0.0", "address used for listening icmp packet")
	cmdReceive.PersistentFlags().StringVarP(&filename, "filename", "f", "", "filename where stored the data received")
	cmdReceive.PersistentFlags().BoolVarP(&progressBar, "progress-bar", "p", false, "print progression of the data reception")

	var rootCmd = &cobra.Command{Use: "qsreceiver"}
	rootCmd.AddCommand(cmdReceive)
	cmdReceive.AddCommand(cmdTruncated)
	rootCmd.Execute()
}
