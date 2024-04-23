package trace

import (
	"fmt"
	"math/rand/v2"
	"net"
)

type IP struct {
	Ip       net.IP
	Verstion IpVer
}

func HostToIp(host string, ipver IpVer) ([]*IP, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, fmt.Errorf("error: looking up hostname %s: %w", host, err)
	}

	if ipver != IPv4 && ipver != IPv6 {
		return nil, fmt.Errorf("error: unknown IP version %s", ipver)
	}

	resIP := make([]*IP, 0)

	for _, ip := range ips {
		if ip.To4() != nil && ipver == IPv4 {
			resIP = append(resIP, &IP{Ip: ip, Verstion: IPv4})
		} else if ip.To16() != nil && ip.To4() == nil && ipver == IPv6 {
			resIP = append(resIP, &IP{Ip: ip, Verstion: IPv6})
		}
	}

	return resIP, nil
}

func IpToHost(ip string) string {
	hosts, err := net.LookupAddr(ip)
	if err != nil {
		return ip
	}

	randIndex := rand.IntN(len(hosts))
	host := hosts[randIndex]

	return host
}
