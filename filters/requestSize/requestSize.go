package requestSize

import (
	"github.com/cyjme/ioa/context"
	"net/http"
	"strconv"
)

type filter struct {
	config config
}

type config struct {
	maxSize int64
}

func New(arg string) (*filter, error) {
	filter := filter{}
	size, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		return &filter, err
	}
	filter.config.maxSize = size

	return &filter, nil
}

func (f *filter) Name() string {
	return "requestSize"
}

func (f *filter) Request(ctx *context.Context) error {

	contentLength := ctx.Request.ContentLength
	if contentLength > f.config.maxSize {
		ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		_, err := ctx.ResponseWriter.Write([]byte("RESP_CONTENT_TOO_LARGE"))
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
