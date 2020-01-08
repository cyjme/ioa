package retry

import (
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"ioa/context"
)

type filter struct {
	maxAttempts int
}

func New(cfg string) (*filter, error) {
	filter := filter{}
	maxAttempts, err := strconv.Atoi(cfg)
	if err != nil {
		maxAttempts = 3
	}

	if maxAttempts <= 0 {
		maxAttempts = 3
	}

	filter.maxAttempts = maxAttempts

	return &filter, nil
}

func (f *filter) Name() string {
	return "retry"
}

func (f *filter) Request(ctx *context.Context) error {
	return nil
}

func (f *filter) Response(ctx *context.Context) error {
	if ctx.Response.StatusCode != http.StatusOK {
		for i := 0; i < f.maxAttempts; i++ {
			reqToTarget, _ := http.NewRequest(
				ctx.Target.Method,
				ctx.Target.Uri+ctx.Target.Path,
				ctx.Request.Body)
			c := &http.Client{}
			resp, err := c.Do(reqToTarget)
			resp.Body.Close()

			if err == nil && resp.StatusCode == http.StatusOK {
				body, _ := ioutil.ReadAll(resp.Body)
				ctx.ResponseBody = body
				ctx.Response = resp
				return nil
			}

			waitTime := math.Pow(2, float64(i)) * 100000000
			logrus.Info(time.Duration(waitTime))
			time.Sleep(time.Duration(waitTime))
		}

	}
	return nil
}
