package trace

import (
	"fmt"
	"net"
)

func GetOutboundIPAndSubnet(ipver IpVer) (string, string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			switch ipver {
			case IPv4:
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String(), fmt.Sprintf("%s", address)
				}
			case IPv6:
				if ipnet.IP.To16() != nil && ipnet.IP.To4() == nil {
					return ipnet.IP.String(), fmt.Sprintf("%s", address)
				}
			}
		}
	}

	panic("no available ip address")
}

func IpInCIDR(network string, ip string) bool {
	_, subnet, _ := net.ParseCIDR(network)
	pip := net.ParseIP(ip)

	return subnet.Contains(pip)
}
