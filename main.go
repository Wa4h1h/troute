package main

import (
	"fmt"
)

func main() {
	ips, err := HostnameToIp("localhost", ipV4)
	if err != nil {
		panic(err)
	}

	for _, ip := range ips {
		fmt.Println(fmt.Sprintf("%s %d", ip.bytes.String(), ip.version))
	}
}
