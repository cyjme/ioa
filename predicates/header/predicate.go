package header

import (
	"github.com/cyjme/ioa/context"
	"regexp"
	"strings"
)

type predicate struct {
	config config
}

type config struct {
	key         string
	valueRegexp *regexp.Regexp
	hasRegexp   bool
}

func New(arg string) (predicate, error) {
	p := predicate{
		config: config{},
	}

	p.config.hasRegexp = false

	query := strings.Split(arg, ",")
	p.config.key = query[0]
	if len(query) == 2 {
		p.config.hasRegexp = true
		value := query[1]

		var err error

		p.config.valueRegexp, err = regexp.Compile(value)

		if err != nil {
			return p, err
		}
	}

	return p, nil
}

func (p *predicate) Name() string {
	return "query"
}

func (p *predicate) Apply(ctx *context.Context) bool {
	reqValue := ctx.Request.Header.Get(p.config.key)
	if p.config.hasRegexp && p.config.valueRegexp.MatchString(reqValue) {
		return true
	}

	return false
}
