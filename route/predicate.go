package route

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/cyjme/ioa/context"
	"github.com/cyjme/ioa/plugin"
	"github.com/cyjme/ioa/predicates/after"
	"github.com/cyjme/ioa/predicates/before"
	"github.com/cyjme/ioa/predicates/between"
	"github.com/cyjme/ioa/predicates/cookie"
	"github.com/cyjme/ioa/predicates/header"
	"github.com/cyjme/ioa/predicates/host"
	"github.com/cyjme/ioa/predicates/method"
	"github.com/cyjme/ioa/predicates/path"
	"github.com/cyjme/ioa/predicates/query"
	"github.com/cyjme/ioa/predicates/remoteAddr"
)

type Predicate interface {
	Name() string
	Apply(ctx *context.Context) bool
}

func createPredicateByName(name string, arg string) (Predicate, error) {
	var predicate Predicate
	switch name {
	case "Host":
		pd, err := host.New(arg)
		if err != nil {
			return &pd, err
		}
		return &pd, nil
	case "Method":
		pd, err := method.New(arg)
		if err != nil {
			return &pd, err
		}
		return &pd, nil
	case "After":
		pd, err := after.New(arg)
		if err != nil {
			return &pd, err
		}
		return &pd, nil
	case "Before":
		pd, err := before.New(arg)
		if err != nil {
			return &pd, err
		}
		return &pd, nil
	case "Between":
		pd, err := between.New(arg)
		if err != nil {
			return &pd, err
		}
		return &pd, nil

	case "Path":
		pd, err := path.New(arg)
		if err != nil {
			return &pd, err
		}
		return &pd, nil

	case "Query":
		pd, err := query.New(arg)
		if err != nil {
			return &pd, err
		}
		return &pd, nil

	case "Cookie":
		pd, err := cookie.New(arg)
		if err != nil {
			return &pd, err
		}
		return &pd, nil

	case "Header":
		pd, err := header.New(arg)
		if err != nil {
			return &pd, err
		}
		return &pd, nil
	case "RemoteAddr":
		pd, err := remoteAddr.New(arg)
		if err != nil {
			return &pd, err
		}
		return &pd, nil

	default:
		var err error
		predicate, err = plugin.CreatePredicateByName(name, arg)
		if err != nil {
			return predicate, err
		}
		if predicate == nil {
			logrus.Error("not found Predicate by name:" + name)
			return predicate, errors.New("not found Predicate by name:" + name)
		}
		return predicate, nil
	}
}
