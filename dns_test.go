package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type input struct {
	hostname string
	ipVer    ipVersion
}

func TestHostnameToIp(t *testing.T) {
	t.Parallel()

	table := []struct {
		name string
		in   input
		out  []string
	}{
		{
			"unknown host",
			input{hostname: "unknown hostname", ipVer: ipV4},
			nil,
		},
		{
			"should only return ipv4 addresses",
			input{hostname: "localhost", ipVer: ipV4},
			[]string{"127.0.0.1"},
		},
		{
			"should only return ipv6 addresses",
			input{hostname: "localhost", ipVer: ipV6},
			[]string{"::1"},
		},
	}

	for _, row := range table {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()

			ips, err := HostnameToIp(row.in.hostname, row.in.ipVer)

			require.NoError(t, err)

			assert.Equal(t, row.out, ips)
		})
	}
}
