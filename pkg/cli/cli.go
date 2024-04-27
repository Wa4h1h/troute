package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/Wa4h1h/troute/pkg/trace"
)

var (
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
		fmt.Fprintln(os.Stdout, `Usage: troute [options] host
Use troute -h or --help for more information.`)
		fmt.Fprintln(os.Stdout, "Options:")
		flag.PrintDefaults()
	}

	flag.Parse()
}

func Run() {
	ParseFlags()

	hosts := flag.Args()

	if len(hosts) != 1 {
		fmt.Fprintln(os.Stderr, "only one hostname can be traced")
		os.Exit(1)
	}

	switch {
	case tcpp:
		port = 80
	case icmpp:
		port = 0
	default:
		port = 33434
	}

	t := trace.DefaultTracer

	if len(os.Args) > 2 {
		t = trace.NewTracer(&trace.TracerConfig{
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

	if err := t.Trace(hosts[0]); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
}

func resolveProto() string {
	if icmpp {
		return trace.ICMP
	}

	if tcpp {
		return trace.TCP
	}

	return trace.UDP
}
