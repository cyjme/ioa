package header

import (
	"github.com/cyjme/ioa/context"
	"net/http"
	"strings"
	"testing"
)

func TestPredicate(t *testing.T) {
	request, _ := http.NewRequest("GET", "https://ioa.letsgo.tech", strings.NewReader(""))
	request.Header.Add("X-Request-Id", "1234")
	request.Header.Add("Token", "abcdef")

	data := []struct {
		arg         string
		request     *http.Request
		expectation bool
	}{
		{
			arg:         "X-Request-Id,\\d+",
			request:     request,
			expectation: true,
		},
		{
			arg:         "Token,abcdef",
			request:     request,
			expectation: true,
		},
		{
			arg:         "Token,\\d+",
			request:     request,
			expectation: false,
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
