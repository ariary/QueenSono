package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ariary/QueenSono/pkg/icmp"
	"github.com/ariary/QueenSono/pkg/utils"
	"github.com/spf13/cobra"
)

func main() {
	//TODO integrityCheck
	//CMD SEND
	//send var
	var remoteAddr string
	var listenAddr string
	var chunkSize int
	var delay int
	var noreply bool
	var encryption string

	var cmdSend = &cobra.Command{ //basic send (send string from stdin)
		Use:   "send [string to send]",
		Short: "Send data to a remote with ICMP protocol",
		Long: `send is for sending  data to a remote queensono receiver waiting.
it uses the icmp protocol.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var data string
			if encryption != "" {
				pubKey := utils.Base64ToPublicKey(encryption)
				data = string(utils.EncryptWithPublicKey([]byte(args[0]), pubKey)) //send byte instead? (compare byte at both endpoint)
			} else {
				data = args[0]
			}
			if noreply {
				icmp.SendNoReply(listenAddr, remoteAddr, chunkSize, delay, data)
			} else {
				icmp.SendReply(listenAddr, remoteAddr, chunkSize, delay, data)
			}

		},
	}

	//cmdSend flag handling
	cmdSend.PersistentFlags().StringVarP(&remoteAddr, "remote", "r", "", "address of remote QueenSono receiver  (required)")
	cmdSend.MarkFlagRequired("remote")

	cmdSend.PersistentFlags().StringVarP(&listenAddr, "listen", "l", "0.0.0.0", "address used for listening echo reply")

	cmdSend.PersistentFlags().IntVarP(&chunkSize, "size", "s", 65507, "Size of each ICMP data section send by packet")

	cmdSend.PersistentFlags().IntVarP(&delay, "delay", "d", 4, "delay between each packet sent")

	cmdSend.PersistentFlags().BoolVarP(&noreply, "noreply", "N", false, "do not wait for echo reply")

	cmdSend.PersistentFlags().StringVarP(&encryption, "encrypt", "e", "", "use encryption for data exchange (provide public key)")

	//CMD SEND FILE
	var cmdSendFile = &cobra.Command{
		Use:   "file [path of the file you want to send]",
		Short: "Send the content of a file to a remote with ICMP protocol",
		Long: `file is for sending  the content of a file to a remote queensono receiver waiting.
it uses the icmp protocol.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			//Retreive file content
			b, err := ioutil.ReadFile(args[0])
			if err != nil {
				fmt.Print(err)
				os.Exit(1)
			}
			data := string(b)
			//send it
			if noreply {
				icmp.SendNoReply(listenAddr, remoteAddr, chunkSize, delay, data)
			} else {
				icmp.SendReply(listenAddr, remoteAddr, chunkSize, delay, data)
			}

		},
	}

	var rootCmd = &cobra.Command{Use: "qssender"}
	rootCmd.AddCommand(cmdSend)
	cmdSend.AddCommand(cmdSendFile)
	rootCmd.Execute()
}
