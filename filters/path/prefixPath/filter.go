package prefixPath

import (
	"ioa/context"
	"strings"
)

type filter struct {
	prefixPath string
}

func New(arg string) (*filter, error) {
	filter := filter{}
	if !strings.HasPrefix(arg, "/") {
		arg = "/" + arg
	}

	filter.prefixPath = arg
	return &filter, nil
}

func (f *filter) Name() string {

	return "prefixPath"
}

func (f *filter) Request(ctx *context.Context) error {
	ctx.Target.Path = f.prefixPath + ctx.Target.Path

	return nil
}

func (f *filter) Response(ctx *context.Context) error {

	return nil
}
