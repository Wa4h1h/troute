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
	port         int
	nprobes      int
	cprobes      int
	chops        int
	probetimeout int
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
	flag.IntVar(&port, "p", 33434, "UDP destination port starts at 33434 or ICMP "+
		"initial sequence number or TCP dst port (Defaults 80)")
	flag.IntVar(&nprobes, "n", 3, "number of probes pro ttl")
	flag.IntVar(&cprobes, "cp", 3, "number of concurrent probes pro ttl")
	flag.IntVar(&chops, "ch", 1, "number of concurrent ttls (only for UDP and ICMP)")
	flag.IntVar(&probetimeout, "t", 3, "probe timeout in seconds")
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

	if len(hosts) != 1 {
		fmt.Println("only one hostname can be traced")
		os.Exit(1)
	}

	if udpp {
		port = 33434
	}

	t := trace.DefaultTracer

	if len(os.Args) > 2 {
		t = trace.NewTracer(&trace.TracerConfig{
			IpVer:        resolveIpVersion(),
			Proto:        resolveProto(),
			StartTTL:     startTTL,
			MaxTTL:       maxTTL,
			Port:         port,
			NProbes:      nprobes,
			CProbes:      cprobes,
			CHopes:       chops,
			ProbeTimeout: probetimeout,
		})
	}

	err := t.Trace(hosts[0])

	fmt.Println(err)
}

func resolveIpVersion() trace.IpVer {
	if ipv6 {
		return trace.IPv6
	}

	if ipv4 {
		return trace.IPv4
	}

	return trace.IPv4
}

func resolveProto() string {
	if icmpp {
		return trace.ICMP
	}

	if tcpp {
		return trace.TCP
	}

	if udpp {
		if ipv6 {
			return trace.UDP6
		}
	}

	return trace.UDP4
}
