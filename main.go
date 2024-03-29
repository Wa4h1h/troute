package main

import "fmt"

func main() {
	i, err := HostnameToIp("localhost", ipV6)
	if err != nil {
		panic(err)
	}

	fmt.Println(i)
}
