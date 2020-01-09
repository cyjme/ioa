package setPath

import (
	"github.com/cyjme/ioa/context"
	"strings"
)

type filter struct {
	setPath string
}

func New(arg string) (*filter, error) {
	filter := filter{}
	filter.setPath = arg

	return &filter, nil
}

func (f *filter) Name() string {

	return "setPath"
}

func (f *filter) Request(ctx *context.Context) error {
	newPath := f.setPath
	for param, paramValue := range ctx.PathParam {
		newPath = strings.Replace(f.setPath, param, paramValue, 1)
	}

	ctx.Target.Path = newPath

	return nil
}

func (f *filter) Response(ctx *context.Context) error {

	return nil
}
