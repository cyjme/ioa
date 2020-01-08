package between

import (
	"errors"
	"ioa/context"
	"strings"
	"time"
)

type predicate struct {
	config config
}
type config struct {
	startDate time.Time
	endDate   time.Time
}

var timeLayout = "2006-01-02 15:04:05"

func New(arg string) (predicate, error) {
	predicate := predicate{}
	configDate := strings.Split(arg, ",")
	if len(configDate) != 2 {
		return predicate, errors.New("you must provide two date, and split by ,")
	}

	startDate, err := time.Parse(timeLayout, configDate[0])
	if err != nil {
		return predicate, err
	}
	endDate, err := time.Parse(timeLayout, configDate[1])
	if err != nil {
		return predicate, err
	}

	predicate.config.startDate = startDate
	predicate.config.endDate = endDate

	return predicate, nil
}

func (p *predicate) Name() string {
	return "predicate"
}

func (p *predicate) Apply(ctx *context.Context) bool {
	now := time.Now()

	if now.After(p.config.startDate) && now.Before(p.config.endDate) {
		return true
	}

	return false
}
