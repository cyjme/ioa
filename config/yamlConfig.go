package config

import (
	"bytes"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type yamlConfig struct {
	filePath string
	kvs      map[string]string
}

func NewYamlConfig(arg string) (*yamlConfig, error) {
	yamlConfig := yamlConfig{
		filePath: arg,
	}

	rawConfig, err := ioutil.ReadFile(arg)
	if err != nil {
		return &yamlConfig, err
	}

	viper.SetConfigType("yaml")
	if err = viper.ReadConfig(bytes.NewBuffer(rawConfig)); err != nil {
		return &yamlConfig, err
	}

	if err := viper.Unmarshal(&yamlConfig.kvs); err != nil {
		logrus.Panic(err)
	}

	return &yamlConfig, nil
}

func (c *yamlConfig) Add(key string, value string) bool {
	return true
}

func (c *yamlConfig) Set(key string, value string) bool {
	return true
}

func (c *yamlConfig) Delete(key string) bool {
	return true
}

func (c *yamlConfig) Get(key string) string {
	return c.kvs[key]
}

func (c *yamlConfig) List() map[string]string {
	return c.kvs
}
