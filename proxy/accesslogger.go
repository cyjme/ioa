package proxy

import (
	"fmt"
	"ioa/context"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type AccessLogger struct {
	file *os.File
}

func NewAccessLogger() *AccessLogger {
	return &AccessLogger{}
}
func (log *AccessLogger) SetOutput(f *os.File) {
	log.file = f
}

func (log *AccessLogger) Append(ctx *context.Context) {
	req := ctx.Request
	text := fmt.Sprintf("%v %v %v %v %v %v %v %v %v %v", time.Now().Format(time.RFC3339), ctx.RouteId, req.RemoteAddr, req.Method, req.URL.Host, req.URL.Path, ctx.Elapsed, ctx.Target.Uri, ctx.Target.Method, ctx.Target.Path)
	fmt.Println(text)

	_, err := log.file.WriteString(text + "\n")
	if err != nil {
		logrus.Error("write access log err")
	}
}
