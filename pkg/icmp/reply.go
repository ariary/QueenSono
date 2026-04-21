package icmp

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ariary/QueenSono/pkg/message"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const triggerPayload = "QS_READY"

// ServeWithEchoReply waits for a QS_READY echo request from the client, then sends
// data back as chunked ICMP echo replies. listenAddr is the local address to bind.
func ServeWithEchoReply(listenAddr string, data string, chunkSize, delay int) error {
	c, err := icmp.ListenPacket("ip4:icmp", listenAddr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	defer c.Close()

	buf := make([]byte, 65535)
	n, peer, err := c.ReadFrom(buf)
	if err != nil {
		return fmt.Errorf("read trigger: %w", err)
	}
	parsed, err := icmp.ParseMessage(ProtocolICMP, buf[:n])
	if err != nil {
		return fmt.Errorf("parse trigger: %w", err)
	}
	if parsed.Type != ipv4.ICMPTypeEcho {
		return fmt.Errorf("expected echo request, got %v", parsed.Type)
	}
	echo, ok := parsed.Body.(*icmp.Echo)
	if !ok {
		return fmt.Errorf("unexpected ICMP body type")
	}
	if string(echo.Data) != triggerPayload {
		return fmt.Errorf("unexpected trigger payload: %q, want %q", string(echo.Data), triggerPayload)
	}

	chunks := Chunks(data, chunkSize)
	chunks = message.QueenSonoMarshall(chunks)

	sendReply := func(payload string, seq int) error {
		m := icmp.Message{
			Type: ipv4.ICMPTypeEchoReply,
			Code: 0,
			Body: &icmp.Echo{
				ID:   echo.ID,
				Seq:  seq,
				Data: []byte(payload),
			},
		}
		b, err := m.Marshal(nil)
		if err != nil {
			return fmt.Errorf("marshal reply: %w", err)
		}
		if _, err := c.WriteTo(b, peer); err != nil {
			return fmt.Errorf("write reply: %w", err)
		}
		return nil
	}

	if err := sendReply(strconv.Itoa(len(chunks)), 0); err != nil {
		return fmt.Errorf("send count announcement: %w", err)
	}
	for i, chunk := range chunks {
		time.Sleep(time.Duration(delay) * time.Second)
		if err := sendReply(chunk, i+1); err != nil {
			return fmt.Errorf("send chunk %d: %w", i, err)
		}
	}
	return nil
}

// TriggerAndReceiveReplies sends a QS_READY echo request to remoteAddr and reassembles
// data sent back as chunked ICMP echo replies. listenAddr is the local address to bind.
func TriggerAndReceiveReplies(listenAddr, remoteAddr string, delay int) (string, error) {
	c, err := icmp.ListenPacket("ip4:icmp", listenAddr)
	if err != nil {
		return "", fmt.Errorf("listen: %w", err)
	}
	defer c.Close()

	dst, err := net.ResolveIPAddr("ip4", remoteAddr)
	if err != nil {
		return "", fmt.Errorf("resolve %s: %w", remoteAddr, err)
	}

	trigger := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte(triggerPayload),
		},
	}
	b, err := trigger.Marshal(nil)
	if err != nil {
		return "", fmt.Errorf("marshal trigger: %w", err)
	}
	if _, err := c.WriteTo(b, dst); err != nil {
		return "", fmt.Errorf("send trigger: %w", err)
	}

	buf := make([]byte, 65535)

	// Read until we get a valid integer count reply.
	// Kernel auto-replies mirror "QS_READY" which fails integer parsing and is skipped.
	if err := c.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
		return "", fmt.Errorf("set deadline: %w", err)
	}
	n := 0
	for {
		pktLen, _, err := c.ReadFrom(buf)
		if err != nil {
			return "", fmt.Errorf("read count reply: %w", err)
		}
		parsed, err := icmp.ParseMessage(ProtocolICMP, buf[:pktLen])
		if err != nil || parsed.Type != ipv4.ICMPTypeEchoReply {
			continue
		}
		echo, ok := parsed.Body.(*icmp.Echo)
		if !ok {
			continue
		}
		count, err := strconv.Atoi(string(echo.Data))
		if err != nil {
			continue
		}
		n = count
		break
	}
	if n == 0 {
		return "", fmt.Errorf("received chunk count of 0")
	}

	timeout := time.Duration(n+2)*time.Duration(delay)*time.Second + 30*time.Second
	if err := c.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return "", fmt.Errorf("set deadline: %w", err)
	}
	dataChunked := make([]string, n)
	received := 0
	for received < n {
		pktLen, _, err := c.ReadFrom(buf)
		if err != nil {
			return "", fmt.Errorf("read chunk: %w", err)
		}
		parsed, err := icmp.ParseMessage(ProtocolICMP, buf[:pktLen])
		if err != nil || parsed.Type != ipv4.ICMPTypeEchoReply {
			continue
		}
		echo, ok := parsed.Body.(*icmp.Echo)
		if !ok {
			continue
		}
		content, idx, err := message.QueenSonoUnmarshall(string(echo.Data))
		if err != nil {
			continue
		}
		if idx >= 0 && idx < n {
			dataChunked[idx] = content
			received++
		}
	}

	return strings.Join(dataChunked, ""), nil
}
