package method

import (
	"ioa/context"
	"strings"
)

type predicate struct {
	config config
}
type config struct {
	methods []string
}

func New(arg string) (predicate, error) {
	toUpperMethods := make([]string, 0)

	methods := strings.Split(arg, ",")
	for _, method := range methods {
		toUpperMethods = append(toUpperMethods, strings.ToUpper(method))
	}

	return predicate{
		config: config{
			methods: toUpperMethods,
		},
	}, nil
}

func (p *predicate) Name() string {
	return "predicate"
}

func (p *predicate) Apply(ctx *context.Context) bool {
	for _, method := range p.config.methods {
		if method == ctx.Request.Method {
			return true
		}
	}

	return false
}
