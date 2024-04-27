package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/Wa4h1h/troute/pkg/trace"
)

var (
	icmpp        bool
	udpp         bool
	startTTL     int
	maxTTL       int
	port         int
	nprobes      int
	cprobes      int
	chops        int
	probetimeout int
)

func ParseFlags() {
	flag.BoolVar(&icmpp, "I", false, "use icmp echo probes")
	flag.BoolVar(&udpp, "U", true, "use udp packet probes")
	flag.IntVar(&startTTL, "start-ttl", 1, "specifies with what TTL to start")
	flag.IntVar(&maxTTL, "max-ttl", 30, "specifies the maximum number of hops (max ttl value)")
	flag.IntVar(&nprobes, "n", 3, "number of probes pro ttl")
	flag.IntVar(&cprobes, "cp", 3, "number of concurrent probes pro ttl")
	flag.IntVar(&chops, "ch", 1, "number of concurrent ttls")
	flag.IntVar(&probetimeout, "t", 3, "probe timeout in seconds")

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
	case icmpp:
		port = 0
	case udpp:
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

	return trace.UDP
}
