package main

import (
	"fmt"
	"net"
)

type ipVersion uint8

const (
	ipV4 ipVersion = iota
	ipV6
)

func HostnameToIp(hostname string, ipVer ipVersion) ([]string, error) {
	if ipVer != ipV4 && ipVer != ipV6 {
		return nil, fmt.Errorf("error: used verion %d is not konwn: %w", ipVer, ErrUnknownIPVersion)
	}

	found, err := net.LookupIP(hostname)
	if err != nil {
		return nil, fmt.Errorf("error: looking up ip address: %w", err)
	}

	ips := make([]string, 0)

	switch ipVer {
	case ipV4:
		for _, ip := range found {
			ipv4 := ip.To4()

			if ipv4 != nil {
				ips = append(ips, ipv4.String())
			}
		}
	case ipV6:
		for _, ip := range found {
			ipv6 := ip.To16()
			if ipv6 != nil && ip.To4() == nil {
				ips = append(ips, ipv6.String())
			}
		}
	}

	return ips, nil
}
