package trace

import (
	"errors"
	"net"
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

	mapTest := map[string][]struct {
		name string
		in   input
		out  []*IP
		err  error
	}{
		"Ok": {
			{
				"should only return ipv4 addresses",
				input{hostname: "localhost", ipVer: ipV4},
				[]*IP{
					{
						bytes:   net.ParseIP("127.0.0.1"),
						version: ipV4,
					},
				},
				nil,
			},
			{
				"should only return ipv6 addresses",
				input{hostname: "localhost", ipVer: ipV6},
				[]*IP{
					{
						bytes:   net.ParseIP("::1"),
						version: ipV6,
					},
				},
				nil,
			},
		},
		"Error": {
			{
				"unknown host",
				input{hostname: "unknown hostname", ipVer: ipV4},
				nil,
				errors.New("error: looking up ip address: lookup unknown host: no such host"),
			},
			{
				"wrong ip version",
				input{hostname: "localhost", ipVer: "10"},
				nil,
				errors.New("error: used verion 10 is not konwn: unknown ip version"),
			},
		},
	}

	for _, row := range mapTest["Ok"] {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()

			ips, err := hostnameToIps(row.in.hostname, row.in.ipVer)

			require.NoError(t, err)

			assert.Equal(t, len(row.out), len(ips))

			for i, val := range row.out {
				assert.Equal(t, string(val.bytes), string(ips[i].bytes))
				assert.Equal(t, val.version, ips[i].version)
			}
		})
	}

	for _, row := range mapTest["Error"] {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()

			_, err := hostnameToIps(row.in.hostname, row.in.ipVer)

			require.Error(t, err)

			assert.Equal(t, row.err.Error(), err.Error())
		})
	}
}
