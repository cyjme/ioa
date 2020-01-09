package routeToRequestUrl

import (
	"errors"
	"github.com/cyjme/ioa/context"
	"strings"
)

type filter struct {
}

func New(cfg string) (*filter, error) {
	filter := filter{}
	return &filter, nil
}

func (f *filter) Name() string {
	return "routeToRequestUrl"
}

func (f *filter) Request(ctx *context.Context) error {
	if strings.HasPrefix(ctx.RouteUri, "lb://") {
		serviceId := strings.TrimPrefix(ctx.RouteUri, "lb://")
		service := (*ctx.ServiceMap)[serviceId]

		//todo 处理 service 的 限流
		ctx.Target.Uri = service.GetUrl()
	} else {
		ctx.Target.Uri = ctx.RouteUri
	}
	if ctx.Target.Uri == "" {
		_, _ = ctx.ResponseWriter.Write([]byte("target uri is null"))
		ctx.Cancel()

		return errors.New("target uri is null")
	}

	ctx.Target.Method = ctx.Request.Method
	ctx.Target.Path = ctx.Request.URL.Path

	return nil
}

func (f *filter) Response(ctx *context.Context) error {
	return nil
}
