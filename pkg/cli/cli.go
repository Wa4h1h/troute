package cli

import (
	"flag"
	"fmt"
	"os"
)

var (
	ipv4         bool
	ipv6         bool
	icmpp        bool
	tcpp         bool
	udpp         bool
	startTTL     int
	maxTTL       int
	port         uint16Flag
	nprobes      uint
	cprobes      uint
	chops        uint
	probetimeout uint
	debug        bool
)

func ParseFlags() {
	flag.BoolVar(&ipv4, "4", true, "use ip version 4")
	flag.BoolVar(&ipv6, "6", false, "use ip version 6")
	flag.BoolVar(&icmpp, "I", false, "use icmp echo for probes")
	flag.BoolVar(&tcpp, "T", false, "use tcp SYN for probes")
	flag.BoolVar(&udpp, "U", true, "use udp packet for probes")
	flag.IntVar(&startTTL, "start-ttl", 1, "specifies with what TTL to start")
	flag.IntVar(&maxTTL, "max-ttl", 30, "specifies the maximum number of hops (max ttl value)")
	flag.Var(&port, "p", "UDP destination port starts at 33434 or ICMP "+
		"initial sequence number or TCP dst port (Defaults 80)")
	flag.UintVar(&nprobes, "n", 3, "number of probes pro ttl")
	flag.UintVar(&cprobes, "cp", 3, "number of concurrent probes pro ttl")
	flag.UintVar(&chops, "ch", 1, "number of concurrent ttls (only for UDP and ICMP)")
	flag.UintVar(&probetimeout, "t", 3, "probe timeout in seconds")
	flag.BoolVar(&debug, "d", false, "enable probe debug")

	flag.Usage = func() {
		fmt.Println(`Usage: troute [options] host
Use troute -h or --help for more information.`)
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	flag.Parse()
}

func Run() {
	ParseFlags()

	hosts := flag.Args()

	if len(hosts) != 0 {
		fmt.Println("inly one hostname can be traced")
		os.Exit(1)
	}

	if udpp {
		port = 33434
	}
}
