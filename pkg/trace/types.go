package trace

import (
	"time"

	"golang.org/x/net/icmp"
)

type IpVer string

const (
	UDP4         = "udp4"
	UDP6         = "udp6"
	TCP          = "tcp"
	ICMP         = "icmp"
	IPv4   IpVer = "ip4"
	IPv6   IpVer = "ip6"
	ICMPv4       = 1
	ICMPv6       = 58
)

type TracerConfig struct {
	IpVer        IpVer
	Proto        string
	StartTTL     int
	MaxTTL       int
	Port         int
	NProbes      int
	CProbes      int
	CHopes       int
	ProbeTimeout int
}

var defaultConfig = &TracerConfig{
	IpVer:        IPv4,
	Proto:        UDP4,
	StartTTL:     1,
	MaxTTL:       30,
	Port:         33434,
	NProbes:      3,
	CProbes:      3,
	CHopes:       1,
	ProbeTimeout: 3,
}

type Probe struct {
	src      string
	host     string
	rtts     []time.Duration
	valid    bool
	icmpType icmp.Type
}

type Hop struct {
	index  int
	probes []*Probe
}
