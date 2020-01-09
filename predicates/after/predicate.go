package after

import (
	"github.com/cyjme/ioa/context"
	"time"
)

type predicate struct {
	config config
}
type config struct {
	date time.Time
}

var timeLayout = "2006-01-02 15:04:05"

func New(arg string) (predicate, error) {
	date, err := time.Parse(timeLayout, arg)
	if err != nil {
		return predicate{}, err
	}

	return predicate{
		config: config{
			date: date,
		},
	}, nil

}

func (p *predicate) Name() string {
	return "predicate"
}

func (p *predicate) Apply(ctx *context.Context) bool {
	now := time.Now()

	return now.After(p.config.date)
}
