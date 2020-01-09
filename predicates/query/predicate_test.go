package query

import (
	"github.com/cyjme/ioa/context"
	"net/http"
	"strings"
	"testing"
)

func TestPredicate(t *testing.T) {
	request, _ := http.NewRequest("GET", "https://ioa.letsgo.tech?name=jerry&phone=123456", strings.NewReader(""))

	data := []struct {
		arg         string
		request     *http.Request
		expectation bool
	}{
		{
			arg:         "name,jerry",
			request:     request,
			expectation: true,
		},
		{
			arg:         "name,tom",
			request:     request,
			expectation: false,
		},
		{
			arg:         "phone,^[0-9]*$",
			request:     request,
			expectation: true,
		},
	}

	for _, item := range data {
		p, err := New(item.arg)
		if err != nil {
			t.Error(err)
		}
		ctx := context.Context{
			Request: item.request,
		}

		if item.expectation != p.Apply(&ctx) {
			t.Error("unexpacted with arg", item.arg)
		}
	}

}
