package cors

import (
	"github.com/cyjme/ioa/context"
	"net/http"
	"strings"
)

type filter struct {
	allowOrigin      string
	allowMethods     string
	allowHeaders     string
	allowCredentials string
	exposeHeaders    string
	maxAge           string
}

func New(arg string) (*filter, error) {
	filter := filter{}
	kvs := make(map[string]string)

	kvsString := strings.Split(arg, ",")
	for _, kvString := range kvsString {
		i := strings.Index(kvString, ":")
		kvs[kvString[:i]] = kvString[i+1:]
	}

	if val, ok := kvs["allowOrigin"]; ok {
		filter.allowOrigin = val
	}

	if val, ok := kvs["allowMethods"]; ok {
		filter.allowMethods = val
	}

	if val, ok := kvs["allowHeaders"]; ok {
		filter.allowHeaders = val
	}

	if val, ok := kvs["allowCredentials"]; ok {
		filter.allowCredentials = val
	}

	if val, ok := kvs["exposeHeaders"]; ok {
		filter.exposeHeaders = val
	}

	if val, ok := kvs["maxAge"]; ok {
		filter.maxAge = val
	}

	return &filter, nil
}

func (f *filter) Request(ctx *context.Context) error {

	if ctx.Request.Method == http.MethodOptions {
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", f.allowOrigin)
		ctx.ResponseWriter.Header().Set("Access-Control-Max-Age", f.maxAge)
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", f.allowMethods)
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", f.allowHeaders)
		ctx.ResponseWriter.Header().Set("Access-Control-Expose-Headers", f.exposeHeaders)
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", f.allowCredentials)
	}

	ctx.ResponseWriter.WriteHeader(200)
	ctx.Cancel()

	return nil
}

func (f *filter) Name() string {
	return "cors"
}

func (f *filter) Response(ctx *context.Context) error {
	return nil
}
