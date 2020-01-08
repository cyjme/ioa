package copyRequest

import (
	"errors"
	"ioa/context"
	"strings"
)

type filter struct {
	copyToUris []string
}

func New(cfg string) (*filter, error) {
	filter := filter{}
	filter.copyToUris = strings.Split(cfg, ",")

	return &filter, nil
}

func (f *filter) Name() string {
	return "copyRequest"
}

func (f *filter) Request(ctx *context.Context) error {

	for _, copyToUri := range f.copyToUris {
		target := context.Target{}

		if strings.HasPrefix(copyToUri, "lb://") {
			copyToUri := strings.TrimPrefix(copyToUri, "lb://")
			service := (*ctx.ServiceMap)[copyToUri]

			//todo 处理 service 的 限流
			target.Uri = service.GetUrl()
		} else {
			target.Uri = copyToUri
		}

		if target.Uri == "" {
			_, _ = ctx.ResponseWriter.Write([]byte("target uri is null"))
			ctx.Cancel()

			return errors.New("target uri is null")
		}

		target.Method = ctx.Request.Method
		target.Path = ctx.Request.URL.Path
		
		ctx.CopyTargets = append(ctx.CopyTargets, target)
	}

	return nil
}

func (f *filter) Response(ctx *context.Context) error {
	return nil
}
