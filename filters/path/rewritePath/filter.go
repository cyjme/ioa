package rewritePath

import (
	"errors"
	"github.com/cyjme/ioa/context"
	"regexp"
	"strings"
)

type filter struct {
	rewritePathRegexpString  string
	rewritePathReplaceString string
}

func New(arg string) (*filter, error) {
	filter := filter{}
	args := strings.Split(arg, ",")
	if len(args) != 2 {
		return &filter, errors.New("arg error")
	}

	filter.rewritePathRegexpString = args[0]
	filter.rewritePathReplaceString = args[1]

	var err error

	return &filter, err
}

func (f *filter) Name() string {

	return "rewritePath"
}

func (f *filter) Request(ctx *context.Context) error {
	rewritePathRegexpStringWithParam := ""
	rewritePathReplaceStringWithParam := ""
	for param, paramValue := range ctx.PathParam {
		rewritePathRegexpStringWithParam = strings.Replace(f.rewritePathRegexpString, param, paramValue, 1)
		rewritePathReplaceStringWithParam = strings.Replace(f.rewritePathReplaceString, param, paramValue, 1)
	}

	regexpC, err := regexp.Compile(rewritePathRegexpStringWithParam)
	if err != nil {
		return err
	}

	newPath := regexpC.ReplaceAllString(ctx.Request.URL.Path, rewritePathReplaceStringWithParam)
	ctx.Target.Path = newPath

	return nil
}

func (f *filter) Response(ctx *context.Context) error {

	return nil
}
