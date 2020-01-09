package breaker

import (
	"errors"
	"fmt"
	"github.com/cyjme/ioa/context"
	"strconv"
	"strings"
	"time"

	"github.com/sony/gobreaker"
)

type filter struct {
	tscb         *gobreaker.TwoStepCircuitBreaker
	done         func(success bool)
	failureRatio float64
	maxRequests  int64
	interval     int64
	timeout      int64
}

/**
  args:
	FailureRatio

  MaxRequests
  Interval
  Timeout
*/
func New(arg string) (*filter, error) {
	filter := filter{
		failureRatio: 0.6,
	}
	st := gobreaker.Settings{}

	rawConfig := make(map[string]string)

	if arg != "" {
		args := strings.Split(arg, ",")

		for _, kvString := range args {
			kvs := strings.Split(kvString, ",")
			if len(kvs) != 2 {
				return &filter, errors.New("arg error")
			}
			rawConfig[kvs[0]] = kvs[1]
		}
	}

	var err error
	if failureRationString, ok := rawConfig["failureRation"]; ok {
		filter.failureRatio, err = strconv.ParseFloat(failureRationString, 64)
		if err != nil {
			return &filter, err
		}
	}

	if intervalString, ok := rawConfig["interval"]; ok {
		filter.interval, err = strconv.ParseInt(intervalString, 10, 64)

		if err != nil {
			return &filter, err
		}
		st.Interval = time.Duration(filter.interval) * time.Second
	}

	if timeoutString, ok := rawConfig["timeoutString"]; ok {
		filter.timeout, err = strconv.ParseInt(timeoutString, 10, 64)

		if err != nil {
			return &filter, err
		}
		st.Timeout = time.Duration(filter.timeout) * time.Second
	}

	if maxRequestString, ok := rawConfig["maxRequest"]; ok {
		filter.maxRequests, err = strconv.ParseInt(maxRequestString, 10, 64)

		if err != nil {
			return &filter, err
		}

		st.MaxRequests = uint32(filter.maxRequests)
	}

	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= filter.failureRatio
	}

	filter.tscb = gobreaker.NewTwoStepCircuitBreaker(st)

	return &filter, nil
}

func (f *filter) Name() string {
	return "breaker"
}

func (f *filter) Request(ctx *context.Context) error {
	var err error
	f.done, err = f.tscb.Allow()
	if err != nil {
		fmt.Println("request not allow")
		//todo 处理已经熔断
		_, err := ctx.ResponseWriter.Write([]byte(err.Error()))
		ctx.Cancel()

		if err != nil {
			return err
		}
	}

	return nil
}

func (f *filter) Response(ctx *context.Context) error {
	fmt.Println("catch response")
	if ctx.RequestError != nil {
		f.done(false)
		return nil
	}

	if ctx.Response.StatusCode > 500 {
		f.done(false)
		return nil
	}

	f.done(true)
	return nil
}
