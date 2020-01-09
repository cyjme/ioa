package removeRequestHeader

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
	return "removeRequestHeader"
}

func (f *filter) Request(ctx *context.Context) error {
	ctx.Request.Header.Del(f.headerKey)

	return nil
}

func (f *filter) Response(ctx *context.Context) error {
	return nil
}
