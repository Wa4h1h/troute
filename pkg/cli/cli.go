package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/Wa4h1h/troute/pkg/trace"
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
	flag.UintVar(&chops, "ch", 1, "number of concurrent ttls")
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
	if len(hosts) < 1 {
		fmt.Println("hostname missing")
		os.Exit(1)
	}

	if len(hosts) > 1 {
		fmt.Println("provide only one hostname")
		os.Exit(1)
	}

	if udpp {
		port = 33434
	}

	t := trace.NewTrace(&trace.Config{
		Host:         hosts[0],
		Ipv:          resolveIpVersion(),
		Proto:        resolveProto(),
		StartTTL:     startTTL,
		MaxTTL:       maxTTL,
		Port:         uint16(port),
		Nprobes:      nprobes,
		Cprobes:      cprobes,
		Chops:        chops,
		Probetimeout: probetimeout,
		Debug:        debug,
	})

	if err := t.Trace(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func resolveIpVersion() string {
	if ipv6 {
		return "ip6"
	}

	return "ip4"
}

func resolveProto() string {
	if tcpp {
		return "tcp"
	}

	if icmpp {
		return "icmp"
	}

	return "udp"
}
