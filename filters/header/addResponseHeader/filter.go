package addResponseHeader

import (
	"github.com/cyjme/ioa/context"
	"strings"
)

type filter struct {
	headerKey   string
	headerValue string
}

func New(arg string) (*filter, error) {
	filter := filter{}
	kv := strings.Split(arg, ",")
	filter.headerKey = kv[0]
	filter.headerValue = kv[1]

	return &filter, nil
}

func (f *filter) Name() string {
	return "addResponseHeader"
}

func (f *filter) Request(ctx *context.Context) error {

	return nil
}

func (f *filter) Response(ctx *context.Context) error {
	ctx.Response.Header.Add(f.headerKey, f.headerValue)

	return nil
}
