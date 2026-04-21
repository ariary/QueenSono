package icmp

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	qsmessage "github.com/ariary/QueenSono/pkg/message"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// GetMessageSizeAndSender waits for the first ICMP packet announcing the chunk count.
func GetMessageSizeAndSender(listenAddr string) (size int, sender string, err error) {
	c, err := icmp.ListenPacket("ip4:icmp", listenAddr)
	if err != nil {
		return 0, "", fmt.Errorf("listen: %w", err)
	}
	defer c.Close()

	packet := make([]byte, 65535)
	n, peer, err := c.ReadFrom(packet)
	if err != nil {
		return 0, "", fmt.Errorf("read packet: %w", err)
	}

	msg, err := icmp.ParseMessage(ProtocolICMP, packet[:n])
	if err != nil {
		return 0, "", fmt.Errorf("parse packet: %w", err)
	}

	switch msg.Type {
	case ipv4.ICMPTypeEcho:
		echo, ok := msg.Body.(*icmp.Echo)
		if !ok {
			return 0, "", fmt.Errorf("unexpected ICMP body type")
		}
		size, err = strconv.Atoi(string(echo.Data))
		if err != nil {
			return 0, "", fmt.Errorf("parse size %q: %w", string(echo.Data), err)
		}
	default:
		return 0, "", fmt.Errorf("got %+v from %v; want echo request", msg, peer)
	}
	return size, peer.String(), nil
}

// Serve receives exactly n ICMP data packets and reassembles them in order.
func Serve(listenAddr string, n int, progressBar bool) (data string, missingPacketIndexes []int, err error) {
	var bar *progressbar.ProgressBar
	if progressBar {
		bar = progressbar.Default(int64(n))
	}
	c, err := icmp.ListenPacket("ip4:icmp", listenAddr)
	if err != nil {
		return "", nil, fmt.Errorf("listen: %w", err)
	}
	defer c.Close()

	dataChunked := make([]string, n)
	indexes := make(map[int]int)
	for i := range n {
		indexes[i] = 0
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			if progressBar {
				getPacketAndBarUpdate(bar, c, dataChunked, indexes, &mu)
			} else {
				getPacket(c, dataChunked, indexes, &mu)
			}
		}()
	}
	wg.Wait()

	data = strings.Join(dataChunked, "")
	missingPacketIndexes = make([]int, 0, len(indexes))
	for index := range indexes {
		missingPacketIndexes = append(missingPacketIndexes, index)
	}
	return data, missingPacketIndexes, nil
}

// ServeTemporary receives up to n ICMP data packets, stopping after (n+2)*delay seconds.
func ServeTemporary(listenAddr string, n int, progressBar bool, delay int) (data string, missingPacketIndexes []int, err error) {
	var bar *progressbar.ProgressBar
	if progressBar {
		bar = progressbar.Default(int64(n))
	}

	c, err := icmp.ListenPacket("ip4:icmp", listenAddr)
	if err != nil {
		return "", nil, fmt.Errorf("listen: %w", err)
	}
	defer c.Close()

	dataChunked := make([]string, n)
	indexes := make(map[int]int)
	for i := range n {
		indexes[i] = 0
	}

	var mu sync.Mutex
	for range n {
		go func() {
			if progressBar {
				getPacketAndBarUpdate(bar, c, dataChunked, indexes, &mu)
			} else {
				getPacket(c, dataChunked, indexes, &mu)
			}
		}()
	}

	time.Sleep(time.Duration((n+2)*delay) * time.Second)
	data = strings.Join(dataChunked, "")
	missingPacketIndexes = make([]int, 0, len(indexes))
	for index := range indexes {
		missingPacketIndexes = append(missingPacketIndexes, index)
	}
	return data, missingPacketIndexes, nil
}

func getPacket(c *icmp.PacketConn, data []string, indexes map[int]int, mu *sync.Mutex) {
	packet := make([]byte, 65535)
	n, _, err := c.ReadFrom(packet)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read packet: %v\n", err)
		return
	}
	msg, err := icmp.ParseMessage(ProtocolICMP, packet[:n])
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse packet: %v\n", err)
		return
	}
	switch msg.Type {
	case ipv4.ICMPTypeEcho:
		echo, ok := msg.Body.(*icmp.Echo)
		if !ok {
			fmt.Fprintf(os.Stderr, "unexpected ICMP body type\n")
			return
		}
		content, index, err := qsmessage.QueenSonoUnmarshall(string(echo.Data))
		if err != nil {
			fmt.Fprintf(os.Stderr, "unmarshal: %v\n", err)
			return
		}
		mu.Lock()
		data[index] = content
		delete(indexes, index)
		mu.Unlock()
	default:
		fmt.Fprintf(os.Stderr, "got %+v; want echo request\n", msg)
	}
}

func getPacketAndBarUpdate(bar *progressbar.ProgressBar, c *icmp.PacketConn, data []string, indexes map[int]int, mu *sync.Mutex) {
	packet := make([]byte, 65535)
	n, _, err := c.ReadFrom(packet)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read packet: %v\n", err)
		return
	}
	msg, err := icmp.ParseMessage(ProtocolICMP, packet[:n])
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse packet: %v\n", err)
		return
	}
	switch msg.Type {
	case ipv4.ICMPTypeEcho:
		echo, ok := msg.Body.(*icmp.Echo)
		if !ok {
			fmt.Fprintf(os.Stderr, "unexpected ICMP body type\n")
			return
		}
		content, index, err := qsmessage.QueenSonoUnmarshall(string(echo.Data))
		if err != nil {
			fmt.Fprintf(os.Stderr, "unmarshal: %v\n", err)
			return
		}
		mu.Lock()
		data[index] = content
		delete(indexes, index)
		mu.Unlock()
		bar.Add(1)
	default:
		fmt.Fprintf(os.Stderr, "got %+v; want echo request\n", msg)
	}
}

// IntegrityCheck listens on listenAddr for a hash and exits when it matches.
func IntegrityCheck(listenAddr, hash string) {
	c, err := icmp.ListenPacket("ip4:icmp", listenAddr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
	defer c.Close()

	for {
		packet := make([]byte, 65535)
		n, _, err := c.ReadFrom(packet)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read: %v\n", err)
			continue
		}
		msg, err := icmp.ParseMessage(ProtocolICMP, packet[:n])
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse: %v\n", err)
			continue
		}
		switch msg.Type {
		case ipv4.ICMPTypeEcho:
			echo, ok := msg.Body.(*icmp.Echo)
			if !ok {
				continue
			}
			if string(echo.Data) == hash {
				fmt.Println("Communication end")
				os.Exit(0)
			}
		default:
			fmt.Fprintf(os.Stderr, "got %+v; want echo request\n", msg)
		}
	}
}
