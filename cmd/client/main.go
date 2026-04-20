package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ariary/QueenSono/pkg/icmp"
	"github.com/ariary/QueenSono/pkg/utils"
	"github.com/spf13/cobra"
)

func main() {
	var remoteAddr string
	var listenAddr string
	var chunkSize int
	var delay int
	var noreply bool
	var key string

	var cmdSend = &cobra.Command{
		Use:   "send [string to send]",
		Short: "Send data to a remote with ICMP protocol",
		Long: `send is for sending data to a remote queensono receiver waiting.
it uses the icmp protocol.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var data string
			if key != "" {
				pubKey := utils.Base64ToPublicKey(key)
				data = utils.Base64EncryptWithPublicKey([]byte(args[0]), pubKey)
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

	cmdSend.PersistentFlags().StringVarP(&remoteAddr, "remote", "r", "", "address of remote QueenSono receiver (required)")
	cmdSend.MarkFlagRequired("remote")
	cmdSend.PersistentFlags().StringVarP(&listenAddr, "listen", "l", "0.0.0.0", "address used for listening echo reply")
	cmdSend.PersistentFlags().IntVarP(&chunkSize, "size", "s", 65488, "Size of each ICMP data section sent per packet")
	cmdSend.PersistentFlags().IntVarP(&delay, "delay", "d", 4, "delay between each packet sent in seconds")
	cmdSend.PersistentFlags().BoolVarP(&noreply, "noreply", "N", false, "do not wait for echo reply")
	cmdSend.PersistentFlags().StringVarP(&key, "key", "k", "", "key used for data encryption (provide public key)")

	var cmdSendFile = &cobra.Command{
		Use:   "file [path of the file you want to send]",
		Short: "Send the content of a file to a remote with ICMP protocol",
		Long: `file is for sending the content of a file to a remote queensono receiver waiting.
it uses the icmp protocol.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			b, err := os.ReadFile(args[0])
			if err != nil {
				log.Fatal(err)
			}
			var data string
			if key != "" {
				pubKey := utils.Base64ToPublicKey(key)
				data = utils.Base64EncryptWithPublicKey(b, pubKey)
			} else {
				data = string(b)
			}
			if noreply {
				icmp.SendNoReply(listenAddr, remoteAddr, chunkSize, delay, data)
			} else {
				icmp.SendReply(listenAddr, remoteAddr, chunkSize, delay, data)
			}
		},
	}

	var receiveRemoteAddr string
	var receiveListenAddr string
	var receiveDelay int
	var receiveFilename string

	var cmdReceive = &cobra.Command{
		Use:   "receive",
		Short: "Receive data from a remote qsreceiver via ICMP echo replies",
		Long: `receive sends a QS_READY trigger to the remote qsreceiver and collects
data sent back as ICMP echo replies.`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			result, err := icmp.TriggerAndReceiveReplies(receiveListenAddr, receiveRemoteAddr, receiveDelay)
			if err != nil {
				log.Fatal(err)
			}
			if receiveFilename != "" {
				if err := os.WriteFile(receiveFilename, []byte(result), 0644); err != nil {
					log.Fatal(err)
				}
				fmt.Println("data saved in", receiveFilename)
			} else {
				fmt.Println("qssender received:", result)
			}
		},
	}
	cmdReceive.PersistentFlags().StringVarP(&receiveRemoteAddr, "remote", "r", "", "address of remote qsreceiver (required)")
	cmdReceive.MarkFlagRequired("remote")
	cmdReceive.PersistentFlags().StringVarP(&receiveListenAddr, "listen", "l", "0.0.0.0", "address used for listening echo replies")
	cmdReceive.PersistentFlags().IntVarP(&receiveDelay, "delay", "d", 4, "delay between packets in seconds")
	cmdReceive.PersistentFlags().StringVarP(&receiveFilename, "filename", "f", "", "filename to save received data")

	var rootCmd = &cobra.Command{Use: "qssender"}
	rootCmd.AddCommand(cmdSend)
	cmdSend.AddCommand(cmdSendFile)
	rootCmd.AddCommand(cmdReceive)
	rootCmd.Execute()
}
