package trace

import (
	"fmt"
	"math/rand/v2"
	"net"
)

type ipVer uint8

const (
	ipV4 ipVer = iota
	ipV6
)

type IP struct {
	Ip       net.IP
	Verstion ipVer
}

func HostToIp(host string, ipver ipVer) ([]*IP, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, fmt.Errorf("error: looking up hostname %s: %w", host, err)
	}

	if ipver != ipV4 && ipver != ipV6 {
		return nil, fmt.Errorf("error: unknown IP version %d", ipver)
	}

	resIP := make([]*IP, 0)

	for _, ip := range ips {
		if ip.To4() != nil && ipver == ipV4 {
			resIP = append(resIP, &IP{Ip: ip, Verstion: ipV4})
		} else if ip.To16() != nil && ip.To4() == nil && ipver == ipV6 {
			resIP = append(resIP, &IP{Ip: ip, Verstion: ipV6})
		}
	}

	return resIP, nil
}

func IpTpHost(ip string) (string, error) {
	hosts, err := net.LookupAddr(ip)
	if err != nil {
		return "", fmt.Errorf("error: getting host from IP %s: %w", ip, err)
	}

	randIndex := rand.IntN(len(hosts))
	host := hosts[randIndex]

	return host, nil
}
