package trace

import (
	"errors"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHostToIp(t *testing.T) {
	t.Parallel()

	table := []struct {
		name   string
		input  string
		output []*IP
		err    error
	}{
		{
			name:   "host returns list of ipv4 address",
			input:  "localhost",
			output: []*IP{{Ip: net.ParseIP("127.0.0.1")}},
		},
		{
			name: "host unknown", input: "unknown",
			output: nil, err: errors.New("error: looking up hostname unknown"),
		},
	}

	for _, row := range table {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()

			res, err := HostToIp(row.input)

			if row.err == nil {
				require.Nil(t, err)
				assert.Equal(t, len(row.output), len(res))

				for i, ip := range res {
					assert.Equal(t, row.output[i].Ip.String(), ip.Ip.String())
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
			input: "not-valid-ip", output: "not-valid-ip",
		},
	}

	for _, row := range table {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()

			host := IpToHost(row.input)

			assert.Equal(t, row.output, host)
		})
	}
}
