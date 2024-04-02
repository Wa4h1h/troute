package trace

type Config struct {
	Host         string
	Ipv4         bool
	Ipv6         bool
	Icmpp        bool
	Tcpp         bool
	Udpp         bool
	StartTTL     uint8
	MaxTTL       uint8
	Port         uint16
	Nprobes      uint
	Cprobes      uint
	Chops        uint
	Probetimeout uint
}

type Tracer struct {
	tc *Config
}

func NewTrace(tc *Config) *Tracer {
	return &Tracer{tc: tc}
}

func (t *Tracer) Trace() error {
	return nil
}
