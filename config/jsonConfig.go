package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

type jsonConfig struct {
	filePath string
	kvs      map[string]string
}

func NewJsonConfig(arg string) (*jsonConfig, error) {
	jsonConfig := jsonConfig{
		filePath: arg,
	}

	rawConfig, err := ioutil.ReadFile(arg)
	if err != nil {
		return &jsonConfig, err
	}

	if err := json.Unmarshal(rawConfig, &jsonConfig.kvs); err != nil {
		logrus.Panic(err)
	}

	return &jsonConfig, nil
}

func (c *jsonConfig) Add(key string, value string) bool {
	return true
}

func (c *jsonConfig) Set(key string, value string) bool {
	return true
}

func (c *jsonConfig) Delete(key string) bool {
	return true
}

func (c *jsonConfig) Get(key string) string {
	return c.kvs[key]
}

func (c *jsonConfig) List() map[string]string {
	return c.kvs
}
