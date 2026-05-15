package icmp

import (
	"net"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// Proto holds the protocol-specific values that differ between ICMPv4 and ICMPv6.
type Proto struct {
	Network   string    // "ip4:icmp" or "ip6:ipv6-icmp"
	IPNetwork string    // "ip4" or "ip6"
	Protocol  int       // 1 (ICMPv4) or 58 (ICMPv6)
	EchoType  icmp.Type // echo request type
	ReplyType icmp.Type // echo reply type
}

var (
	proto4 = Proto{
		Network:   "ip4:icmp",
		IPNetwork: "ip4",
		Protocol:  1,
		EchoType:  ipv4.ICMPTypeEcho,
		ReplyType: ipv4.ICMPTypeEchoReply,
	}
	proto6 = Proto{
		Network:   "ip6:ipv6-icmp",
		IPNetwork: "ip6",
		Protocol:  58,
		EchoType:  ipv6.ICMPTypeEchoRequest,
		ReplyType: ipv6.ICMPTypeEchoReply,
	}
)

// DetectProto returns the ICMPv4 or ICMPv6 protocol config based on the address.
// Empty string and unparseable addresses default to IPv4 for backward compatibility.
func DetectProto(addr string) Proto {
	if addr == "" {
		return proto4
	}
	ip := net.ParseIP(addr)
	if ip == nil {
		return proto4
	}
	if ip.To4() != nil {
		return proto4
	}
	return proto6
}

// ResolveListenAddr adjusts listenAddr to match the protocol version.
// If listenAddr is a wildcard or empty and doesn't match proto, swap to the
// correct wildcard. Specific addresses are returned as-is.
func ResolveListenAddr(listenAddr string, proto Proto) string {
	switch listenAddr {
	case "", "0.0.0.0":
		if proto.IPNetwork == "ip6" {
			return "::"
		}
		return "0.0.0.0"
	case "::":
		if proto.IPNetwork == "ip4" {
			return "0.0.0.0"
		}
		return "::"
	default:
		return listenAddr
	}
}
