package rateLimit

import (
	"fmt"
	"github.com/cyjme/ioa/context"
	"strconv"
	"strings"

	"golang.org/x/time/rate"
)

type filter struct {
	limiter *rate.Limiter
	config  config
}

type config struct {
	limit rate.Limit
	burst int
}

//arg: limiter:100,burst:1000
func New(arg string) (*filter, error) {
	filter := filter{}
	rawConfig := make(map[string]string)

	args := strings.Split(arg, ",")
	for _, kv := range args {
		tmp := strings.Split(kv, ":")
		rawConfig[tmp[0]] = tmp[1]
	}

	limit, err := strconv.ParseFloat(rawConfig["limit"], 64)
	if err != nil {
		return &filter, err
	}
	burst, err := strconv.ParseInt(rawConfig["burst"], 10, 64)
	if err != nil {
		return &filter, err
	}

	filter.config.limit = rate.Limit(limit)
	filter.config.burst = int(burst)

	filter.limiter = rate.NewLimiter(filter.config.limit, filter.config.burst)

	return &filter, nil
}

func (f *filter) Name() string {
	return "rateLimit"
}

func (f *filter) Request(ctx *context.Context) error {
	if !f.limiter.Allow() {
		fmt.Println("not allow")
		_, err := ctx.ResponseWriter.Write([]byte("rateLimit not allow"))
		if err != nil {
			return err
		}

		ctx.Cancel()
	}
	return nil
}

func (f *filter) Response(ctx *context.Context) error {
	return nil
}
