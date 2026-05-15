package icmp

import (
	"testing"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

func TestDetectProto(t *testing.T) {
	tests := []struct {
		name      string
		addr      string
		network   string
		ipNetwork string
		protocol  int
		echoType  icmp.Type
		replyType icmp.Type
	}{
		{"IPv4", "192.168.1.1", "ip4:icmp", "ip4", 1, ipv4.ICMPTypeEcho, ipv4.ICMPTypeEchoReply},
		{"IPv6 link-local", "fe80::1", "ip6:ipv6-icmp", "ip6", 58, ipv6.ICMPTypeEchoRequest, ipv6.ICMPTypeEchoReply},
		{"IPv6 full", "2001:db8::1", "ip6:ipv6-icmp", "ip6", 58, ipv6.ICMPTypeEchoRequest, ipv6.ICMPTypeEchoReply},
		{"wildcard v4", "0.0.0.0", "ip4:icmp", "ip4", 1, ipv4.ICMPTypeEcho, ipv4.ICMPTypeEchoReply},
		{"wildcard v6", "::", "ip6:ipv6-icmp", "ip6", 58, ipv6.ICMPTypeEchoRequest, ipv6.ICMPTypeEchoReply},
		{"empty defaults to v4", "", "ip4:icmp", "ip4", 1, ipv4.ICMPTypeEcho, ipv4.ICMPTypeEchoReply},
		{"loopback v4", "127.0.0.1", "ip4:icmp", "ip4", 1, ipv4.ICMPTypeEcho, ipv4.ICMPTypeEchoReply},
		{"loopback v6", "::1", "ip6:ipv6-icmp", "ip6", 58, ipv6.ICMPTypeEchoRequest, ipv6.ICMPTypeEchoReply},
		{"unparseable defaults to v4", "not-an-ip", "ip4:icmp", "ip4", 1, ipv4.ICMPTypeEcho, ipv4.ICMPTypeEchoReply},
		{"IPv4-mapped IPv6 returns v4", "::ffff:192.0.2.1", "ip4:icmp", "ip4", 1, ipv4.ICMPTypeEcho, ipv4.ICMPTypeEchoReply},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := DetectProto(tt.addr)
			if p.Network != tt.network {
				t.Errorf("Network: got %s, want %s", p.Network, tt.network)
			}
			if p.IPNetwork != tt.ipNetwork {
				t.Errorf("IPNetwork: got %s, want %s", p.IPNetwork, tt.ipNetwork)
			}
			if p.Protocol != tt.protocol {
				t.Errorf("Protocol: got %d, want %d", p.Protocol, tt.protocol)
			}
			if p.EchoType != tt.echoType {
				t.Errorf("EchoType: got %v, want %v", p.EchoType, tt.echoType)
			}
			if p.ReplyType != tt.replyType {
				t.Errorf("ReplyType: got %v, want %v", p.ReplyType, tt.replyType)
			}
		})
	}
}

func TestResolveListenAddr(t *testing.T) {
	tests := []struct {
		name   string
		listen string
		proto  Proto
		want   string
	}{
		{"v6 proto + v4 wildcard", "0.0.0.0", proto6, "::"},
		{"v6 proto + empty", "", proto6, "::"},
		{"v4 proto + v6 wildcard", "::", proto4, "0.0.0.0"},
		{"v4 proto + empty", "", proto4, "0.0.0.0"},
		{"v4 proto + v4 wildcard", "0.0.0.0", proto4, "0.0.0.0"},
		{"v6 proto + v6 wildcard", "::", proto6, "::"},
		{"specific v4 addr", "10.0.0.1", proto4, "10.0.0.1"},
		{"specific v6 addr", "fe80::1", proto6, "fe80::1"},
		{"specific v4 addr with v6 proto", "10.0.0.1", proto6, "10.0.0.1"},
		{"specific v6 addr with v4 proto", "fe80::1", proto4, "fe80::1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveListenAddr(tt.listen, tt.proto)
			if got != tt.want {
				t.Errorf("ResolveListenAddr(%q, %s): got %q, want %q", tt.listen, tt.proto.IPNetwork, got, tt.want)
			}
		})
	}
}
