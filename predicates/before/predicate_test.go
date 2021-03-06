package before

import (
	"github.com/cyjme/ioa/context"
	"testing"
	"time"
)

func TestPredicate(t *testing.T) {

	data := []struct {
		arg         string
		expectation bool
	}{
		{
			arg:         time.Now().Add(10 * time.Second).Format(timeLayout),
			expectation: true,
		},
		{
			arg:         time.Now().Truncate(10 * time.Second).Format(timeLayout),
			expectation: false,
		},
	}

	for _, item := range data {
		p, err := New(item.arg)
		if err != nil {
			t.Error(err)
		}
		if item.expectation != p.Apply(&context.Context{}) {
			t.Error("unexpected")
		}
	}
}
