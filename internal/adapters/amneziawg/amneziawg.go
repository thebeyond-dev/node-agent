package amneziawg

import (
	"fmt"
	"net/netip"
	"sync"
)

const (
	defaultInterface = "awg0"
	defaultConfPath  = "/etc/amnezia/amneziawg/awg0.conf"
)

var ErrIPPoolExhausted = fmt.Errorf("ip pool exhausted")

type Adapter struct {
	mu           sync.Mutex
	intf         string
	configPath   string
	endpoint     string
	clientDNS    string
	network      netip.Prefix
	magicParams  magicParams
	serverPubKey string
	listenPort   string
}

type magicParams struct {
	Jc, Jmin, Jmax     string
	S1, S2, S3, S4     string
	H1, H2, H3, H4     string
	I1, I2, I3, I4, I5 string
}

type Option func(*Adapter)

func New(endpoint, networkCIDR, clientDNS string, opts ...Option) (*Adapter, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint is required")
	}

	prefix, err := netip.ParsePrefix(networkCIDR)
	if err != nil {
		return nil, fmt.Errorf("invalid network cidr %q: %w", networkCIDR, err)
	}

	adapter := &Adapter{
		intf:       defaultInterface,
		configPath: defaultConfPath,
		endpoint:   endpoint,
		network:    prefix,
		clientDNS:  clientDNS,
	}

	for _, opt := range opts {
		opt(adapter)
	}
	return adapter, nil
}

func WithInterface(name string) Option {
	return func(a *Adapter) { a.intf = name }
}

func WithConfigPath(path string) Option {
	return func(a *Adapter) { a.configPath = path }
}
