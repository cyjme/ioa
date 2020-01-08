package remoteAddr

import (
	"ioa/context"
	"net"
	"strings"
)

type predicate struct {
	config config
}

type config struct {
	ipNets []*net.IPNet
}

func New(arg string) (predicate, error) {
	p := predicate{
		config: config{
			ipNets: make([]*net.IPNet, 0),
		},
	}

	cidrStrings := strings.Split(arg, ",")
	for _, cidr := range cidrStrings {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			return p, err
		}
		p.config.ipNets = append(p.config.ipNets, ipNet)
	}

	return p, nil
}

func (p *predicate) Name() string {
	return "remoteAddr"
}

func (p *predicate) Apply(ctx *context.Context) bool {
	for _, ipNet := range p.config.ipNets {
		if ipNet.Contains(net.ParseIP(ctx.Request.RemoteAddr)) {
			return true
		}
	}

	return false
}
