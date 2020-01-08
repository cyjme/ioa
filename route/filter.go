package route

import (
	"ioa/context"
	"ioa/filters/breaker"
	"ioa/filters/cache"
	"ioa/filters/cors"
	"ioa/filters/header/addRequestHeader"
	"ioa/filters/header/addResponseHeader"
	"ioa/filters/header/removeRequestHeader"
	"ioa/filters/header/removeResponseHeader"
	"ioa/filters/header/secureHeaders"
	"ioa/filters/header/setRequestHeader"
	"ioa/filters/header/setResponseHeader"
	"ioa/filters/kuipAuth"
	"ioa/filters/path/prefixPath"
	"ioa/filters/path/rewritePath"
	"ioa/filters/path/setPath"
	"ioa/filters/rateLimit"
	"ioa/filters/requestSize"
	"ioa/filters/retry"
	"ioa/plugin"
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
