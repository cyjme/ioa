package path

import (
	"ioa/context"
	"net/http"
	"strings"
	"testing"
)

//http://www.csyangchen.com/go-http-router.html

func TestPredicate(t *testing.T) {
	data := []struct {
		arg         string
		request     *http.Request
		expartation bool
	}{
		{
			arg:         "/users",
			request:     getRequest(http.NewRequest("GET", "/users", strings.NewReader(""))),
			expartation: true,
		},
		{
			arg:         "/users/:userId",
			request:     getRequest(http.NewRequest("GET", "/users/1234", strings.NewReader(""))),
			expartation: true,
		},
		{
			arg:         "/users/:userId([0-9]+)",
			request:     getRequest(http.NewRequest("GET", "/users/1234", strings.NewReader(""))),
			expartation: true,
		},
		{
			arg:         "/users/:userId([0-9]+)",
			request:     getRequest(http.NewRequest("GET", "/users/abc", strings.NewReader(""))),
			expartation: false,
		},
	}

	for _, item := range data {
		p, err := New(item.arg)
		if err != nil {
			t.Error("new predicate error", err)
		}
		ctx := context.Context{
			Request: item.request,
		}

		if item.expartation != p.Apply(&ctx) {
			t.Error("unexpacted item", item.arg)
		}
	}

}

func getRequest(req *http.Request, err error) *http.Request {

	return req
}
