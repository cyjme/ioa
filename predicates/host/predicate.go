package host

import (
	"fmt"
	"github.com/cyjme/ioa/context"
	"regexp"
	"strings"
)

type config struct {
	hosts       []string
	hostRegexps []*regexp.Regexp
}

type predicate struct {
	config config
}

func New(arg string) (predicate, error) {
	hosts := strings.Split(arg, ",")
	hostRegexps := make([]*regexp.Regexp, 0)

	for _, host := range hosts {
		hostRegexStr := ""
		if strings.Contains(host, "*") {
			hostRegexStr = strings.Replace(host, ".", "\\.", -1)
			hostRegexStr = strings.Replace(hostRegexStr, "*", ".+", -1)
			hostRegexp, _ := regexp.Compile(fmt.Sprintf("^%s$", hostRegexStr))

			hostRegexps = append(hostRegexps, hostRegexp)
			continue
		}
		hosts = append(hosts, host)
	}

	pd := predicate{
		config: config{
			hosts:       hosts,
			hostRegexps: hostRegexps,
		},
	}

	return pd, nil
}

func (p *predicate) Name() string {
	return "host"
}

func (p *predicate) Apply(ctx *context.Context) bool {
	for _, host := range p.config.hosts {
		if host == ctx.Request.Host {
			return true
		}
	}
	for _, hostRegexp := range p.config.hostRegexps {
		if hostRegexp.MatchString(ctx.Request.Host) {
			return true
		}
	}

	return false
}
