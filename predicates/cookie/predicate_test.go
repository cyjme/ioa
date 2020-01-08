package cookie

import (
	"ioa/context"
	"net/http"
	"strings"
	"testing"
)

func TestPredicate(t *testing.T) {
	request, _ := http.NewRequest("GET", "https://ioa.letsgo.tech", strings.NewReader(""))
	request.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "1234",
	})

	data := []struct {
		arg         string
		request     *http.Request
		expectation bool
	}{
		{
			arg:         "token,\\d+",
			request:     request,
			expectation: true,
		},
		{
			arg:         "token,1234",
			request:     request,
			expectation: true,
		},
		{
			arg:         "token,abc",
			request:     request,
			expectation: false,
		},
		{
			arg:         "Auth",
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
