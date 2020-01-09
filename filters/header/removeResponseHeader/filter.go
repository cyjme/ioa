package removeResponseHeader

import (
	"github.com/cyjme/ioa/context"
)

type filter struct {
	headerKey string
}

func New(arg string) (*filter, error) {
	filter := filter{}
	filter.headerKey = arg

	return &filter, nil
}

func (f *filter) Name() string {
	return "removeResponseHeader"
}

func (f *filter) Request(ctx *context.Context) error {
	return nil
}

func (f *filter) Response(ctx *context.Context) error {
	ctx.Response.Header.Del(f.headerKey)

	return nil
}
