package trace

import (
	"fmt"
	"net"
)

type ipVersion string

const (
	ipV4 ipVersion = "ip4"
	ipV6 ipVersion = "ip6"
)

type IP struct {
	bytes   net.IP
	version ipVersion
}

func hostnameToIps(hostname string, ipVer ipVersion) ([]*IP, error) {
	if ipVer != ipV4 && ipVer != ipV6 {
		return nil, fmt.Errorf("error: used verion %s is not konwn: %w", ipVer, ErrUnknownIPVersion)
	}

	found, err := net.LookupIP(hostname)
	if err != nil {
		return nil, fmt.Errorf("error: looking up ip address: %w", err)
	}

	ips := make([]*IP, 0)

	switch ipVer {
	case ipV4:
		for _, ip := range found {
			ipv4 := ip.To4()

			if ipv4 != nil {
				ips = append(ips, &IP{bytes: ipv4, version: ipVer})
			}
		}
	case ipV6:
		for _, ip := range found {
			ipv6 := ip.To16()

			if ipv6 != nil && ip.To4() == nil {
				ips = append(ips, &IP{bytes: ipv6, version: ipV6})
			}
		}
	}

	return ips, nil
}

func ipToHostnames(ip string) ([]string, error) {
	hosts, err := net.LookupAddr(ip)
	if err != nil {
		return nil, fmt.Errorf("error: looking hostname for %s: %w", ip, err)
	}

	return hosts, nil
}
