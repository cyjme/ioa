package context

import (
	"bytes"
	"ioa/service"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Filter interface {
	Name() string
	Request(ctx *Context) error
	Response(ctx *Context) error
}

//Context
type Context struct {
	RequestError   error
	ResponseWriter http.ResponseWriter
	Response       *http.Response
	Request        *http.Request
	MatchedRoute   bool
	Stop           bool

	Client *http.Client

	PathParam map[string]string
	RouteId   string
	RouteUri  string
	Target    Target

	Filters []Filter

	CopyTargets []Target

	ServiceMap *map[string]service.Service

	Elapsed      time.Duration
	ResponseBody []byte
}

// Target means a real host, include method, uri, path
type Target struct {
	Method string
	Uri    string
	Path   string
}

// Cancel set context.Stop is true, when filters was range, it will check context.Stop
func (ctx *Context) Cancel() {
	ctx.Stop = true
}

func (ctx *Context) Reset() {
	ctx.RequestError = nil
	ctx.ResponseWriter = nil
	ctx.Response = nil
	ctx.Request = nil
	ctx.MatchedRoute = false
	ctx.Stop = false
	ctx.PathParam = nil
	ctx.RouteId = ""
	ctx.RouteUri = ""
	ctx.Target = Target{}
	ctx.Filters = nil
	ctx.CopyTargets = nil
	ctx.ServiceMap = nil
	ctx.ResponseBody = nil
}

func (ctx *Context) Do() {
	start := time.Now()
	ctx.execRequestFilters()
	if ctx.Stop {
		return
	}

	var resp *http.Response
	var err error

	reqToTarget, err := http.NewRequest(ctx.Target.Method, ctx.Target.Uri+ctx.Target.Path, ctx.Request.Body)
	if err != nil {
		logrus.Error("do request error:", err)
	} else {
		reqToTarget.Header = ctx.Request.Header
	}

	go ctx.doReqeustToCopyTarget()

	resp, err = ctx.Client.Do(reqToTarget)
	if err != nil {
		logrus.Error("do request error:", err)
		return
	}

	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)

	ctx.ResponseBody = buf.Bytes()
	ctx.RequestError = err
	ctx.Response = resp

	ctx.execResponseFilters()
	if ctx.Stop {
		return
	}

	ctx.sendResp()
	ctx.Elapsed = time.Since(start)
	//	accessLogger.Append(ctx)
}

func (ctx *Context) execRequestFilters() {
	for _, filter := range ctx.Filters {
		err := filter.Request(ctx)
		if err != nil {
			logrus.Error(err)
		}
	}
}

func (ctx *Context) execResponseFilters() {
	for _, filter := range ctx.Filters {
		err := filter.Response(ctx)
		if err != nil {
			logrus.Error(err)
		}
	}
}

func (ctx *Context) doReqeustToCopyTarget() {
	for _, target := range ctx.CopyTargets {
		reqToCopyTarget, err := http.NewRequest(target.Method, target.Uri+target.Path, ctx.Request.Body)
		if err != nil {
			logrus.Error("do request error:", err)
		}
		reqToCopyTarget.Header = ctx.Request.Header

		resp, err := ctx.Client.Do(reqToCopyTarget)
		if err != nil {
			logrus.Error("do request error:", err)
		}
		defer resp.Body.Close()
	}
}

func (ctx *Context) sendResp() {
	for name, values := range ctx.Response.Header {
		ctx.ResponseWriter.Header()[name] = values
	}

	ctx.ResponseWriter.WriteHeader(ctx.Response.StatusCode)
	_, err := ctx.ResponseWriter.Write(ctx.ResponseBody)
	if err != nil {
		logrus.Error("response write error", err)
	}
}
