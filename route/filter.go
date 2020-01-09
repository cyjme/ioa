package route

import (
	"github.com/cyjme/ioa/context"
	"github.com/cyjme/ioa/filters/breaker"
	"github.com/cyjme/ioa/filters/cache"
	"github.com/cyjme/ioa/filters/cors"
	"github.com/cyjme/ioa/filters/header/addRequestHeader"
	"github.com/cyjme/ioa/filters/header/addResponseHeader"
	"github.com/cyjme/ioa/filters/header/removeRequestHeader"
	"github.com/cyjme/ioa/filters/header/removeResponseHeader"
	"github.com/cyjme/ioa/filters/header/secureHeaders"
	"github.com/cyjme/ioa/filters/header/setRequestHeader"
	"github.com/cyjme/ioa/filters/header/setResponseHeader"
	"github.com/cyjme/ioa/filters/kuipAuth"
	"github.com/cyjme/ioa/filters/path/prefixPath"
	"github.com/cyjme/ioa/filters/path/rewritePath"
	"github.com/cyjme/ioa/filters/path/setPath"
	"github.com/cyjme/ioa/filters/rateLimit"
	"github.com/cyjme/ioa/filters/requestSize"
	"github.com/cyjme/ioa/filters/retry"
	"github.com/cyjme/ioa/plugin"
)

type Filter interface {
	Name() string
	Request(context *context.Context) error
	Response(context *context.Context) error
}
type Filters []Filter

func createFilterByName(name, arg string) (Filter, error) {
	var filter Filter

	switch name {
	case "AddRequestHeader":
		return addRequestHeader.New(arg)
	case "AddResponseHeader":
		return addResponseHeader.New(arg)
	case "RemoveRequestHeader":
		return removeRequestHeader.New(arg)
	case "RemoveResponseHeader":
		return removeResponseHeader.New(arg)
	case "SecureHeaders":
		return secureHeaders.New(arg)
	case "SetRequestHeader":
		return setRequestHeader.New(arg)
	case "SetResponseHeader":
		return setResponseHeader.New(arg)

	case "PrefixPath":
		return prefixPath.New(arg)
	case "RewritePath":
		return rewritePath.New(arg)
	case "SetPath":
		return setPath.New(arg)

	case "RequestSize":
		return requestSize.New(arg)

	case "Breaker":
		return breaker.New(arg)

	case "RateLimit":
		return rateLimit.New(arg)

	case "KuipAuth":
		return kuipAuth.New(arg)

	case "Cors":
		return cors.New(arg)

	case "Cache":
		return cache.New(arg)

	case "Retry":
		return retry.New(arg)

	default:
		var err error
		filter, err = plugin.CreateFilterByName(name, arg)
		if err != nil {
			return filter, err
		}

		return filter, err
	}
}
