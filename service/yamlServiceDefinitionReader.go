package service

import (
	"bytes"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type yamlServiceDefinitionReader struct {
}

func (s *yamlServiceDefinitionReader) getServiceDefinitions(config string) ([]ServiceDefinition, error) {
	serviceDefinitions := make([]ServiceDefinition, 0)

	rawConfig, err := ioutil.ReadFile(config)
	logrus.Debug("read Service config file ", string(rawConfig))

	if err != nil {
		logrus.Panic(err)
		return serviceDefinitions, err
	}

	viper.SetConfigType("yml")
	err = viper.ReadConfig(bytes.NewBuffer(rawConfig))
	if err != nil {
		return serviceDefinitions, err
	}

	var yamlServiceDefinitions struct {
		Services []ServiceDefinition `mapstructure:"services"`
	}

	if err := viper.Unmarshal(&yamlServiceDefinitions); err != nil {
		logrus.Panic(err)
		return serviceDefinitions, err
	}

	// check is repeat
	existsMap := make(map[string]bool)
	for _, service := range yamlServiceDefinitions.Services {
		if _, ok := existsMap[service.Id]; ok {
			logrus.Panic("service id repeat")
		}
		existsMap[service.Id] = true
	}

	return yamlServiceDefinitions.Services, err
}
