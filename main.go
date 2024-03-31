package main

import (
	"fmt"

	"github.com/Wa4h1h/troute/pkg/trace"
)

func main() {
	ips, err := trace.HostnameToIp("facebook.com", trace.IpV6)
	if err != nil {
		panic(err)
	}

	fmt.Println(ips)
}
