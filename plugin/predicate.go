package plugin

import (
	"errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"github.com/cyjme/ioa/context"
	"plugin"
)

// 加载.so出错时,请参考：
// https://stackoverflow.com/questions/42388090/go-1-8-plugin-use-custom-interface (建议使用方案2,方案1会出现版本不一致的错误)

type Predicate interface {
	Name() string
	Apply(ctx *context.Context) bool
}

var allPredicate = make(map[string]func(arg string) (interface{}, error))

func LoadPredicate(dirPath string) {
	if dirPath == "" {
		return
	}
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}
	for _, file := range dir {
		fileName := file.Name()

		if fileName[len(fileName)-3:] == ".so" {
			predicatePlugin, err := plugin.Open(dirPath + "/" + fileName)
			if err != nil {
				logrus.Error("open *.so plugin file error", err)
				return
			}
			symbol, err := predicatePlugin.Lookup("New")

			if err != nil {
				logrus.Error("lookup plugin err: ", err)
				return
			}
			getPredicate, ok := symbol.(func(arg string) (interface{}, error))
			if !ok {
				logrus.Error("predicate plugin type conv err ", ok)
				return
			}
			allPredicate[fileName[:len(fileName)-3]] = getPredicate
		}
	}
}

func CreatePredicateByName(name, arg string) (Predicate, error) {
	getPredicate, ok := allPredicate[name]
	if !ok {
		logrus.Error("predicate  " + name + " not exist")
		return nil, errors.New("predicate " + name + " not exist")
	}
	predicate, err := getPredicate(arg)
	return predicate.(Predicate), err
}
