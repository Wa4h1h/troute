package trace

import (
	"crypto/rand"
	"fmt"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"math/big"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/net/icmp"
)

type Config struct {
	Host         string
	Ipv          string
	Proto        string
	StartTTL     int
	MaxTTL       int
	Port         uint16
	Nprobes      uint
	Cprobes      uint
	Chops        uint
	Probetimeout uint
	Debug        bool
}

type Tracer struct {
	tc         *Config
	lconn      *icmp.PacketConn
	sconn      net.PacketConn
	dstIp      *IP
	currentTTL int
}

func NewTrace(tc *Config) *Tracer {
	return &Tracer{tc: tc}
}

func (t *Tracer) Trace() error {
	ludp := "udp4"
	if ipVersion(t.tc.Ipv) == ipV6 {
		ludp = "udp6"
	}

	lconn, err := icmp.ListenPacket(ludp, "")
	if err != nil {
		return fmt.Errorf("error: listening for icmp packets: %w", err)
	}

	t.lconn = lconn

	ips, err := hostnameToIps(t.tc.Host, ipVersion(t.tc.Ipv))
	if err != nil {
		return err
	}

	switch t.tc.Proto {
	case "udp":
		conn, err := net.ListenUDP(ludp, nil)
		if err != nil {
			return fmt.Errorf("error: creating udp listener: %w", err)
		}

		t.sconn = conn
	case "tcp":
	case "icmp":
	}

	t.currentTTL = t.tc.StartTTL

	if err := t.setTTL(); err != nil {
		return err
	}

	randIpIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(ips))))

	t.dstIp = ips[randIpIndex.Int64()]

	fmt.Println(t.dstIp.bytes)

	return t.trace()
}

func (t *Tracer) setTTL() error {
	switch t.tc.Ipv {
	case string(ipV4):
		if err := ipv4.NewPacketConn(t.sconn).SetTTL(t.currentTTL); err != nil {
			return fmt.Errorf("error: setting start ttl %d: %w", t.tc.StartTTL, err)
		}
	case string(ipV6):
		if err := ipv6.NewPacketConn(t.sconn).SetHopLimit(t.currentTTL); err != nil {
			return fmt.Errorf("error: setting start hop limit %d: %w", t.tc.StartTTL, err)
		}
	}

	t.currentTTL++

	return nil
}

func (t *Tracer) trace() error {
	limiter := make(chan struct{}, t.tc.Chops)
	results := make(chan *Hop)
	done := make(chan struct{})

	defer func() {
		close(limiter)
		close(results)
	}()

	go func(d chan<- struct{}, r chan<- *Hop, l <-chan struct{}) {
		for i := t.tc.StartTTL; i <= t.tc.MaxTTL; i++ {
			limiter <- struct{}{}
			go t.executeHop(r, d, l)
			if strings.Contains(t.tc.Proto, "udp") {
				t.tc.Port++
			}
			if err := t.setTTL(); err != nil {
				panic(err)
			}
		}

		d <- struct{}{}
	}(done, results, limiter)

	go func() {
		for range results {

		}
	}()

	<-done

	return nil
}

func (t *Tracer) executeHop(hopResult chan<- *Hop, done chan<- struct{}, l <-chan struct{}) {
	limiter := make(chan struct{}, t.tc.Cprobes)
	results := make(chan *Probe)
	probes := make([]*Probe, t.tc.Nprobes)

	defer func() {
		close(limiter)
		close(results)
	}()

	go func() {
		for range t.tc.Nprobes {
			limiter <- struct{}{}
			go t.executeProbe(results)
		}
	}()

	for range t.tc.Nprobes {
		probe := <-results
		probes = append(probes, probe)
	}

	hopResult <- &Hop{probes: probes}
	<-l
}

func (t *Tracer) executeProbe(results chan<- *Probe) {
	start := time.Now()

	switch t.tc.Proto {
	case "udp":
		dstIp := net.UDPAddr{IP: t.dstIp.bytes, Port: int(t.tc.Port)}
		if _, err := t.sconn.WriteTo([]byte{0x0}, &dstIp); err != nil {
			if t.tc.Debug {
				fmt.Fprintf(os.Stdout, "error: writing to %s:%d: %s", dstIp.IP.String(), dstIp.Port, err.Error())
			}

			results <- &Probe{host: "*", rtt: -1}

			return
		}
	case "tcp":
	case "icmp":
	}

	reply := make([]byte, 1000)
	if err := t.lconn.SetReadDeadline(time.Now().Add(time.Duration(t.tc.Probetimeout) * time.Second)); err != nil {
		if t.tc.Debug {
			fmt.Fprintf(os.Stdout, "error: setting read timeout : %s", err.Error())
		}

		results <- &Probe{host: "*", rtt: -1}

		return
	}
	_, _, err := t.lconn.ReadFrom(reply)
	if err == nil {
		/*host := srcAddr.String()
		end := time.Since(start).Seconds() / 1000
		hosts, errLook := ipToHostnames(srcAddr.String())
		if errLook != nil {
			fmt.Println(errLook)
			hostRandIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(hosts))))
			host = hosts[hostRandIndex.Int64()-1]
		}*/

		var proto int

		switch t.tc.Ipv {
		case string(ipV4):
			proto = ProtocolICMP
		case string(ipV6):
			proto = ProtocolICMPv6
		}

		msg, errp := icmp.ParseMessage(proto, reply)
		if errp != nil {
			results <- &Probe{host: "*", rtt: -1}

			return
		}

		fmt.Println(*msg)

		end := time.Since(start).Seconds() / 1000

		results <- &Probe{srcIp: "o", host: "3", rtt: end}
	} else {
		if t.tc.Debug {
			fmt.Fprintf(os.Stdout, "error: reading probe to: %s", err.Error())
		}

		results <- &Probe{host: "*", rtt: -1}
	}

	return
}
