package trace

import (
	"fmt"
	"math/rand/v2"
	"net"
)

type IP struct {
	Ip net.IP
}

func HostToIp(host string) ([]*IP, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, fmt.Errorf("error: looking up hostname %s: %w", host, err)
	}

	resIP := make([]*IP, 0)

	for _, ip := range ips {
		if ip.To4() != nil {
			resIP = append(resIP, &IP{Ip: ip})
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
