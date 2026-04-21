package main

import (
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ariary/QueenSono/pkg/icmp"
	"github.com/ariary/QueenSono/pkg/utils"
	"github.com/spf13/cobra"
)

func main() {
	var listenAddr string
	var filename string
	var progressBar bool
	var encryption bool

	var cmdReceive = &cobra.Command{
		Use:   "receive",
		Short: "receive data from icmp packet",
		Long: `receive is for receiving data from a remote queensono sender.
it uses the icmp protocol.`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			var privKey *rsa.PrivateKey
			if encryption {
				var pubKey *rsa.PublicKey
				privKey, pubKey = utils.GenerateKeyPair(4096)
				fmt.Println("Public Key (copy and paste it in qssender):")
				fmt.Println(utils.PublicKeyToBase64(pubKey))
			}
			size, sender, err := icmp.GetMessageSizeAndSender(listenAddr)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Sender:", sender, ", Number of packet wanted:", size)
			message, missingPacketsIndexes, err := icmp.Serve(listenAddr, size, progressBar)
			if err != nil {
				log.Fatal(err)
			}
			if encryption {
				message = string(utils.Base64DecryptWithPrivateKey(message, privKey))
			}
			if len(missingPacketsIndexes) > 0 {
				fmt.Println("Missing packet:")
				for i := range len(missingPacketsIndexes) {
					fmt.Print(missingPacketsIndexes[i], " ")
				}
				fmt.Println()
			}
			if filename != "" {
				if err := os.WriteFile(filename, []byte(message), 0644); err != nil {
					log.Fatal(err)
				}
				fmt.Println("data saved in", filename)
			} else {
				fmt.Println("qsreceiver received:", message)
			}
		},
	}

	cmdReceive.PersistentFlags().StringVarP(&listenAddr, "listen", "l", "0.0.0.0", "address used for listening icmp packet")
	cmdReceive.PersistentFlags().StringVarP(&filename, "filename", "f", "", "filename where stored the data received")
	cmdReceive.PersistentFlags().BoolVarP(&progressBar, "progress-bar", "p", false, "print progression of the data reception")
	cmdReceive.PersistentFlags().BoolVarP(&encryption, "encrypt", "e", false, "use encryption for data exchange")

	var cmdTruncated = &cobra.Command{
		Use:   "truncated [delay]",
		Short: "receive data without waiting indefinitely for all packets",
		Long:  `Receive data from icmp packets and stop after (n+2)*delay seconds. Reports missing packet indexes.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			size, sender, err := icmp.GetMessageSizeAndSender(listenAddr)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Sender:", sender, ", Number of packet wanted:", size)
			delay, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatalf("invalid delay %q: %v", args[0], err)
			}
			message, missingPacketsIndexes, err := icmp.ServeTemporary(listenAddr, size, progressBar, delay)
			if err != nil {
				log.Fatal(err)
			}
			if len(missingPacketsIndexes) > 0 {
				fmt.Println("Missing packet:")
				for i := range len(missingPacketsIndexes) {
					fmt.Print(missingPacketsIndexes[i], " ")
				}
				fmt.Println()
			}
			if filename != "" {
				if err := os.WriteFile(filename, []byte(message), 0644); err != nil {
					log.Fatal(err)
				}
				fmt.Println("data saved in", filename)
			} else {
				fmt.Println("qsreceiver received:", message)
			}
		},
	}

	var replySendListenAddr string
	var replySendChunkSize int
	var replySendDelay int

	var cmdReplySend = &cobra.Command{
		Use:   "reply-send [data to send]",
		Short: "Send data to a remote qssender via ICMP echo replies",
		Long: `reply-send waits for a QS_READY trigger from qssender, then sends data
back as chunked ICMP echo replies.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := icmp.ServeWithEchoReply(replySendListenAddr, args[0], replySendChunkSize, replySendDelay); err != nil {
				log.Fatal(err)
			}
		},
	}
	cmdReplySend.PersistentFlags().StringVarP(&replySendListenAddr, "listen", "l", "0.0.0.0", "address to listen on for trigger")
	cmdReplySend.PersistentFlags().IntVarP(&replySendChunkSize, "size", "s", 65488, "size of each ICMP data chunk")
	cmdReplySend.PersistentFlags().IntVarP(&replySendDelay, "delay", "d", 4, "delay between reply packets in seconds")

	var rootCmd = &cobra.Command{Use: "qsreceiver"}
	rootCmd.AddCommand(cmdReceive)
	cmdReceive.AddCommand(cmdTruncated)
	rootCmd.AddCommand(cmdReplySend)
	rootCmd.Execute()
}
