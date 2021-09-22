package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ariary/QueenSono/pkg/icmp"
	"github.com/ariary/QueenSono/pkg/utils"
	"github.com/spf13/cobra"
)

func main() {
	//CMD RECEIVE
	//receive var
	var listenAddr string
	var filename string
	var progressBar bool
	var encryption bool

	var cmdReceive = &cobra.Command{
		Use:   "receive",
		Short: "receive data from icmp packet",
		Long: `receive is for receiving  data from a remote queensono sender.
it uses the icmp protocol.`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			if encryption {
				//privKey, pubKey := utils.GenerateKeyPair(1024)
				_, pubKey := utils.GenerateKeyPair(1024)
				pubKeyEnc := utils.PublicKeyToBase64(pubKey)
				fmt.Println("Public Key (copy and paste it in qsreceiver):")
				fmt.Println(pubKeyEnc)
			}
			size, sender := icmp.GetMessageSizeAndSender(listenAddr)
			fmt.Println("Sender:", sender, ", Number of packet wanted:", size)
			message, missingPacketsIndexes := icmp.Serve(listenAddr, size, progressBar)
			//icmp.SendHashedmessage(message, sender)
			//Integrity check
			//Print missing packet
			if len(missingPacketsIndexes) > 0 {
				fmt.Println("Missing packet:")
				for i := 0; i < len(missingPacketsIndexes); i++ {
					fmt.Print(missingPacketsIndexes[i], " ")
				}
				fmt.Println()
			}
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
	cmdReceive.PersistentFlags().BoolVarP(&encryption, "encrypt", "e", false, "use encryption for data exchange")

	//CMD TRUNCATED
	var cmdTruncated = &cobra.Command{
		Use:   "truncated [delay]",
		Short: "receive data from icmp packet and do not wait indefinitively for all the packet.(indicate the packet missing at the end)",
		Long:  `receive data from icmp packet and do not wait indefinitively for all the packet. Hence you could receive truncated data. If the data isn't fully retrieved  it return the index of the missing packet`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			size, sender := icmp.GetMessageSizeAndSender(listenAddr)
			fmt.Println("Sender:", sender, ", Number of packet wanted:", size)
			delay, err := strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}
			message, missingPacketsIndexes := icmp.ServeTemporary(listenAddr, size, progressBar, delay)
			//icmp.SendHashedmessage(message, sender) //Integrity check
			//Print missing packet
			if len(missingPacketsIndexes) > 0 {
				fmt.Println("Missing packet:")
				for i := 0; i < len(missingPacketsIndexes); i++ {
					fmt.Print(missingPacketsIndexes[i], " ")
				}
				fmt.Println()
			}

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

	var rootCmd = &cobra.Command{Use: "qsreceiver"}
	rootCmd.AddCommand(cmdReceive)
	cmdReceive.AddCommand(cmdTruncated)
	rootCmd.Execute()
}
