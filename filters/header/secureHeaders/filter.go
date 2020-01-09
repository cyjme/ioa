package secureHeaders

import (
	"github.com/cyjme/ioa/context"
)

type filter struct {
}

var secureHeaders = map[string]string{
	"X-Xss-Protection":                  "1; mode=block",
	"Strict-Transport-Security":         "max-age=631138519",
	"X-Frame-Options":                   "DENY",
	"X-Content-Type-Options":            "nosniff",
	"Referrer-Policy":                   "no-referrer",
	"Content-Security-Policy":           "default-src 'self' https:; font-src 'self' https: data:; img-src 'self' https: data:; object-src 'none'; script-src https:; style-src 'self' https: 'unsafe-inline'",
	"X-Download-Options":                "noopen",
	"X-Permitted-Cross-Domain-Policies": "none",
}

func New(arg string) (*filter, error) {
	filter := filter{}

	return &filter, nil
}

func (f *filter) Name() string {
	return "secureHeaders"
}

func (f *filter) Request(ctx *context.Context) error {
	return nil
}

func (f *filter) Response(ctx *context.Context) error {
	for headerKey, headerValue := range secureHeaders {
		ctx.Response.Header.Add(headerKey, headerValue)
	}

	return nil
}
