package trace

import (
	"time"

	"golang.org/x/net/icmp"
)

const (
	UDP    = "udp4"
	TCP    = "tcp"
	ICMP   = "ip4:icmp"
	ICMPv4 = 1
)

type TracerConfig struct {
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
	Proto:        UDP,
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
	rtt      time.Duration
	valid    bool
	icmpType icmp.Type
}

type Hop struct {
	index  int
	probes []*Probe
}
