package remoteAddr

import (
	"github.com/cyjme/ioa/context"
	"net/http"
	"strings"
	"testing"
)

func TestPredicate(t *testing.T) {
	request, _ := http.NewRequest("GET", "https://ioa.letsgo.tech", strings.NewReader(""))
	request.RemoteAddr = "192.168.1.10"

	data := []struct {
		arg         string
		request     *http.Request
		expectation bool
	}{
		{
			arg:         "192.168.1.1/24",
			request:     request,
			expectation: true,
		},
		{
			arg:         "192.168.18.1/24",
			request:     request,
			expectation: false,
		},
		{
			arg:         "192.168.18.1/24,192.168.1.1/24",
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
