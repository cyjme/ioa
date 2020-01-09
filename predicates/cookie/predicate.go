package cookie

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
	cookie, err := ctx.Request.Cookie(p.config.key)

	if err != nil {
		return false
	}
	if p.config.hasRegexp && p.config.valueRegexp.MatchString(cookie.Value) {
		return true
	}

	return false
}
