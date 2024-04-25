package trace

import (
	"cmp"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net"
	"os"
	"slices"
	"strings"
	"time"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"

	"golang.org/x/net/icmp"
)

type Tracer struct {
	cfg     *TracerConfig
	lconn   *icmp.PacketConn
	sconn   net.PacketConn
	nextTTL int
	dst     *IP
	laddr   string
	lsubnet string
}

var DefaultTracer = &Tracer{cfg: defaultConfig}

func NewTracer(cfg *TracerConfig) *Tracer {
	return &Tracer{cfg: cfg, nextTTL: cfg.StartTTL}
}

func (t *Tracer) Trace(host string) error {
	ips, err := HostToIp(host, t.cfg.IpVer)
	if err != nil {
		return err
	}

	if len(ips) < 1 {
		fmt.Fprintf(os.Stdout, "troute: %s has no %s ip\n",
			host, t.cfg.IpVer)

		return nil
	}

	ipIndex := rand.IntN(len(ips))
	t.dst = ips[ipIndex]

	if len(ips) > 1 {
		fmt.Fprintf(os.Stdout, "troute: %s has more than one ip. "+
			"%s will be used\n", host, t.dst.Ip.String())
	}

	fmt.Fprintf(os.Stdout, "troute %s (%s) with max hops %d\n", host, t.dst.Ip.String(), t.cfg.MaxTTL)

	return t.trace()
}

func (t *Tracer) trace() error {
	t.laddr, t.lsubnet = GetOutboundIPAndSubnet(t.cfg.IpVer)
	conn, err := icmp.ListenPacket(t.cfg.Proto, t.laddr)
	if err != nil {
		return fmt.Errorf("error: starting icmp listener: %w", err)
	}

	t.lconn = conn

	switch t.cfg.Proto {
	case UDP4, UDP6:
		conn, err := net.ListenPacket(t.cfg.Proto, ":")
		if err != nil {
			return fmt.Errorf("error: dialing %s %s: %w", t.cfg.Proto, t.dst.Ip.String(), err)
		}

		t.sconn = conn

	case TCP:
	case ICMP:
	default:
		return fmt.Errorf("error: unkown protocol: %s", t.cfg.Proto)
	}

	if err := t.hops(); err != nil {
		return err
	}

	return nil
}

func (t *Tracer) setTTL() error {
	switch t.cfg.IpVer {
	case IPv4:
		if err := ipv4.NewPacketConn(t.sconn).SetTTL(t.nextTTL); err != nil {
			return fmt.Errorf("error: setting ipv4 ttl: %w", err)
		}
	case IPv6:
		if err := ipv6.NewPacketConn(t.sconn).SetHopLimit(t.nextTTL); err != nil {
			return fmt.Errorf("error: setting ipv6 ttl: %w", err)
		}
	}

	t.nextTTL++

	return nil
}

func (t *Tracer) hops() error {
	limit := make(chan struct{}, t.cfg.CHopes)
	results := make([]*Hop, 0)
	errChan := make(chan error)
	hopItems := make(chan *Hop)

	go func() {
		for i := t.cfg.StartTTL; i <= t.cfg.MaxTTL; i++ {
			limit <- struct{}{}

			go func() {
				var (
					result *Hop
					err    error
				)

				result, err = t.execHop()
				if err != nil {
					errChan <- fmt.Errorf("error: executing hop with ttl %d: %w", i, err)

					return
				}

				hopItems <- result
				<-limit
			}()
		}
	}()

	for {
		select {
		case res := <-hopItems:
			results = append(results, res)

			slices.SortStableFunc(results, func(a, b *Hop) int {
				return cmp.Compare(a.index, b.index)
			})

			t.printHop(results[len(results)-1])

			if (t.istLastHop(res) && len(results) >= res.index+1) || len(results) >= t.cfg.MaxTTL {
				return nil
			}

			t.cfg.Port++
		case err := <-errChan:
			return err
		}
	}
}

func (t *Tracer) execHop() (*Hop, error) {
	limit := make(chan struct{}, t.cfg.CProbes)
	results := make([]*Probe, 0)
	errChan := make(chan error)
	probeItems := make(chan *Probe)

	if err := t.setTTL(); err != nil {
		return nil, err
	}

	go func() {
		for i := 0; i < t.cfg.NProbes; i++ {
			limit <- struct{}{}
			go func() {
				p, err := t.execProbe()
				if err != nil {
					errChan <- err

					return
				}

				probeItems <- p
				<-limit
			}()
		}
	}()

	for {
		select {
		case err := <-errChan:
			return nil, err
		case item := <-probeItems:
			results = append(results, item)

			if len(results) >= t.cfg.NProbes {
				h := new(Hop)
				h.probes = results
				h.index = t.getHopIndex()

				return h, nil
			}
		}
	}
}

func (t *Tracer) execProbe() (*Probe, error) {
	var dst net.Addr

	switch t.cfg.Proto {
	case UDP4, UDP6:
		dst = &net.UDPAddr{IP: t.dst.Ip, Port: t.cfg.Port}
	case ICMP:
		dst = &net.IPAddr{IP: t.dst.Ip}
	}

	if err := t.sconn.SetWriteDeadline(time.Now().Add(time.Duration(t.cfg.ProbeTimeout) * time.Second)); err != nil {
		return nil, fmt.Errorf("error: setting write deadline: %w", err)
	}

	start := time.Now()

	if _, err := t.sconn.WriteTo([]byte{0x0}, dst); err != nil {
		return nil, fmt.Errorf("error: sending probe: %w", err)
	}

	if err := t.lconn.SetReadDeadline(time.Now().Add(time.Duration(t.cfg.ProbeTimeout) * time.Second)); err != nil {
		return nil, fmt.Errorf("error: setting read deadline: %w", err)
	}

	reply := make([]byte, 0)

	for {
		tmp := make([]byte, 512)

		n, addr, err := t.lconn.ReadFrom(tmp)
		if err != nil {
			if os.IsTimeout(err) {
				return &Probe{src: "*", host: "*", valid: false}, nil
			}
			if !errors.Is(err, io.EOF) {
				return nil, fmt.Errorf("error: reading probe: %w", err)
			}
		}

		reply = append(reply, tmp[:n]...)
		if (n == 0 || n < 512) && len(reply) > 4 {
			p, err := t.parseICMP(reply, addr)
			if err != nil {
				return nil, err
			}

			end := time.Since(start)

			p.rtts = make([]time.Duration, 0)
			p.rtts = append(p.rtts, end)

			return p, err
		}
	}
}

func (t *Tracer) parseICMP(bytes []byte, src net.Addr) (*Probe, error) {
	icmpVer := ICMPv4
	if t.cfg.IpVer == IPv6 {
		icmpVer = ICMPv6
	}

	msg, err := icmp.ParseMessage(icmpVer, bytes)
	if err != nil {
		return nil, fmt.Errorf("error: parsing icmp message: %w", err)
	}

	probe := new(Probe)
	srcIP := src.String()[:strings.Index(src.String(), ":")]
	host := IpToHost(srcIP)

	probe.host = host
	probe.src = srcIP
	probe.valid = true
	probe.icmpType = msg.Type

	return probe, nil
}

func (t *Tracer) getHopIndex() int {
	var index int

	switch t.cfg.Proto {
	case UDP4, UDP6:
		index = t.cfg.Port - 33434
	case ICMP:
	}

	return index
}

func (t *Tracer) istLastHop(h *Hop) bool {
	last := false

Loop:

	for _, p := range h.probes {
		switch t.cfg.IpVer {
		case IPv4:
			if p.icmpType == ipv4.ICMPTypeDestinationUnreachable {
				last = true

				break Loop
			}
		case IPv6:
			if p.icmpType == ipv6.ICMPTypeDestinationUnreachable {
				last = true

				break Loop
			}
		}
	}

	return last
}

func (t *Tracer) printHop(h *Hop) {
	fmt.Fprintf(os.Stdout, "%d ", h.index+1)

	probes := make([]string, 0)
	for i, p := range h.probes {
		if !p.valid {
			fmt.Fprint(os.Stdout, " * ")
		} else {
			if !slices.Contains(probes, p.host) {
				if i > 0 {
					fmt.Fprint(os.Stdout, "\n  ")
				}
				fmt.Fprintf(os.Stdout, "%s (%s)", p.host, p.src)
			}

			for _, rtt := range p.rtts {
				fmt.Fprintf(os.Stdout, " %v ", rtt)
			}

		}

		probes = append(probes, p.host)
	}

	fmt.Fprint(os.Stdout, "\n")
}
