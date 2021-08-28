package main

import (
	"fmt"
	"strings"

	"github.com/ariary/QueenSono/pkg/icmp"
	"github.com/spf13/cobra"
)

func main() {
	// 	var echoTimes int

	// 	var cmdPrint = &cobra.Command{
	// 		Use:   "print [string to print]",
	// 		Short: "Print anything to the screen",
	// 		Long: `print is for printing anything back to the screen.
	// For many years people have printed back to the screen.`,
	// 		Args: cobra.MinimumNArgs(1),
	// 		Run: func(cmd *cobra.Command, args []string) {
	// 			fmt.Println("Print: " + strings.Join(args, " "))
	// 		},
	// 	}

	// 	var cmdEcho = &cobra.Command{
	// 		Use:   "echo [string to echo]",
	// 		Short: "Echo anything to the screen",
	// 		Long: `echo is for echoing anything back.
	// Echo works a lot like print, except it has a child command.`,
	// 		Args: cobra.MinimumNArgs(1),
	// 		Run: func(cmd *cobra.Command, args []string) {
	// 			fmt.Println("Echo: " + strings.Join(args, " "))
	// 		},
	// 	}

	// 	var cmdTimes = &cobra.Command{
	// 		Use:   "times [string to echo]",
	// 		Short: "Echo anything to the screen more times",
	// 		Long: `echo things multiple times back to the user by providing
	// a count and a string.`,
	// 		Args: cobra.MinimumNArgs(1),
	// 		Run: func(cmd *cobra.Command, args []string) {
	// 			for i := 0; i < echoTimes; i++ {
	// 				fmt.Println("Echo: " + strings.Join(args, " "))
	// 			}
	// 		},
	// 	}

	// 	cmdTimes.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")

	var cmdICMP = &cobra.Command{
		Use:   "icmp [mode] [flags] data",
		Short: "Sending/Receiving data trough ICMP",
		Long: `Use this command to sending data trough ICMP protocol.
Prefer short data.`,
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("ICMP: " + strings.Join(args, " "))
		},
	}

	var cmdICMPSend = &cobra.Command{
		Use:   "send",
		Short: "specify which mode you want",
		Long:  `'send' for exfiltrate data, 'receive' for receiving the data`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			icmp.IcmpSendRaw(args[0])
		},
	}

	var rootCmd = &cobra.Command{Use: "queensono"}
	rootCmd.AddCommand(cmdICMP)
	cmdICMP.AddCommand(cmdICMPSend)
	rootCmd.Execute()
}
