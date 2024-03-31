package cli

import (
	"flag"
	"fmt"
)

var (
	ipv4         bool
	ipv6         bool
	icmpp        bool
	tcpp         bool
	udpp         bool
	startTTL     uint8Flag = 1
	maxTTL       uint8Flag = 30
	port         uint16Flag
	nprobes      uint
	cprobes      uint
	chops        uint
	probetimeout uint
)

func ParseFlags() {
	flag.BoolVar(&ipv4, "4", true, "use ip version 4")
	flag.BoolVar(&ipv6, "6", false, "use ip version 6")
	flag.BoolVar(&icmpp, "I", false, "use icmp echo for probes")
	flag.BoolVar(&tcpp, "T", false, "use tcp SYN for probes")
	flag.BoolVar(&udpp, "U", true, "use udp packet for probes")
	flag.Var(&startTTL, "start-ttl", "specifies with what TTL to start. Defaults to 1")
	flag.Var(&maxTTL, "max-ttl", "specifies the maximum number of hops (max ttl value). Defaults 30")
	flag.Var(&port, "p", "UDP destination port starts at 33434 or ICMP initial sequence number or TCP dst port (Defaults 80)")
	flag.UintVar(&nprobes, "n", 3, "number of probes pro ttl")
	flag.UintVar(&cprobes, "cp", 3, "number of concurrent probes pro ttl")
	flag.UintVar(&chops, "ch", 3, "number of concurrent ttls")
	flag.UintVar(&probetimeout, "t", 3, "probe timeout in seconds")

	flag.Usage = func() {
		fmt.Println(`Usage: troute [options]
Use troute -h or --help for more information.`)
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	flag.Parse()
}