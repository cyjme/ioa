package method

import (
	"github.com/cyjme/ioa/context"
	"net/http"
	"testing"
)

var data = []struct {
	arg         string
	req         *http.Request
	expectation bool
}{
	{
		arg: "get,post",
		req: &http.Request{
			Method: http.MethodGet,
		},
		expectation: true,
	},
	{
		arg: "post",
		req: &http.Request{
			Method: http.MethodGet,
		},
		expectation: false,
	},
	{
		arg: "Post",
		req: &http.Request{
			Method: http.MethodGet,
		},
		expectation: false,
	},
	{
		arg: "Post",
		req: &http.Request{
			Method: http.MethodPost,
		},
		expectation: true,
	},
}

func TestApply(t *testing.T) {
	for _, item := range data {
		pd, err := New(item.arg)
		if err != nil {
			t.Error(err)
		}
		ctx := context.Context{
			Request: item.req,
		}
		result := pd.Apply(&ctx)
		if result != item.expectation {
			t.Error("item unexpected", item)
		}
	}
}
