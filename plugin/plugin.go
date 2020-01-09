package plugin

import (
	"errors"
	"io/ioutil"
	"github.com/cyjme/ioa/context"
	"plugin"

	"github.com/sirupsen/logrus"
)

type Filter interface {
	Name() string
	Request(context *context.Context) error
	Response(context *context.Context) error
}

var allFilters = make(map[string]func(arg string) (Filter, error))

func Init(dirPath string) {
	if dirPath == "" {
		return
	}

	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}
	allFilters = make(map[string]func(arg string) (Filter, error))

	for _, file := range dir {
		fileName := file.Name()

		if fileName[len(fileName)-3:] == ".so" {
			filterPlugin, err := plugin.Open(dirPath + "/" + fileName)
			if err != nil {
				logrus.Error("open *.so plugin file error", err)
				return
			}
			symbol, err := filterPlugin.Lookup("New")

			if err != nil {
				logrus.Error("lookup plugin err: ", err)
				return
			}
			getFilter, ok := symbol.(func(arg string) (Filter, error))

			if !ok {
				logrus.Error("filter plugin type conv err", ok)
				return
			}

			allFilters[fileName[:len(fileName)-3]] = getFilter
		}
	}
}

func CreateFilterByName(name, arg string) (Filter, error) {
	getFilter, ok := allFilters[name]
	if !ok {
		logrus.Error("filter " + name + " not exist")
		return nil, errors.New("filter " + name + " not exist")
	}
	filter, err := getFilter(arg)

	return filter, err
}
