package trace

import (
	"errors"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type input struct {
	hostname string
	ipVer    IpVer
}

func TestHostToIp(t *testing.T) {
	t.Parallel()

	table := []struct {
		name   string
		input  input
		output []*IP
		err    error
	}{
		{
			name:   "host returns list of ipv4 address",
			input:  input{hostname: "localhost", ipVer: IPv4},
			output: []*IP{{Ip: net.ParseIP("127.0.0.1"), Verstion: IPv4}},
		},
		{
			name:   "host returns list of ipv6 address",
			input:  input{hostname: "localhost", ipVer: IPv6},
			output: []*IP{{Ip: net.ParseIP("::1"), Verstion: IPv6}},
		},
		{
			name: "host unknown", input: input{hostname: "unknown", ipVer: IPv4},
			output: nil, err: errors.New("error: looking up hostname unknown"),
		},
		{
			name: "wrong ip version", input: input{hostname: "localhost", ipVer: "10"},
			output: nil, err: errors.New("error: unknown IP version"),
		},
	}

	for _, row := range table {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()

			res, err := HostToIp(row.input.hostname, row.input.ipVer)

			if row.err == nil {
				require.Nil(t, err)
				assert.Equal(t, len(row.output), len(res))

				for i, ip := range res {
					assert.Equal(t, row.output[i].Ip.String(), ip.Ip.String())
					assert.Equal(t, row.output[i].Verstion, ip.Verstion)
				}

			} else {
				require.NotNil(t, err)
				assert.True(t, strings.HasPrefix(err.Error(), row.err.Error()))
			}
		})
	}
}

func TestIpTpHost(t *testing.T) {
	t.Parallel()

	table := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:  "ip returns a valid hostname",
			input: "127.0.0.1", output: "localhost",
		},
		{
			name:  "ip can not be mapped to a valid hostname",
			input: "not-valid-ip", err: errors.New("error: getting host from IP not-valid-ip"),
		},
	}

	for _, row := range table {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()

			host := IpToHost(row.input)

			if row.err == nil {
				require.Nil(t, err)
				assert.Equal(t, row.output, host)
			} else {
				require.NotNil(t, err)
				assert.True(t, strings.HasPrefix(err.Error(), row.err.Error()))
			}
		})
	}
}
