package service

import (
	"bytes"
	"errors"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

type yamlServiceDefinitionWriter struct {
}

func (writer *yamlServiceDefinitionWriter) AddServiceDefinition(config string, serviceDefinition ServiceDefinition) error {
	rawConfig, err := ioutil.ReadFile(config)
	logrus.Debug("read Service config file ", string(rawConfig))

	if err != nil {
		return err
	}

	viper.SetConfigType("yml")
	err = viper.ReadConfig(bytes.NewBuffer(rawConfig))
	if err != nil {
		return err
	}

	var yamlServiceDefinitions struct {
		Services []ServiceDefinition `mapstructure:"services"`
	}

	if err := viper.Unmarshal(&yamlServiceDefinitions); err != nil {
		return err
	}

	//check rid is repeat
	for _, service := range yamlServiceDefinitions.Services {
		if service.Id == serviceDefinition.Id {
			return errors.New("service id repeat")
		}
	}

	//write new Service into yaml config
	yamlServiceDefinitions.Services = append(yamlServiceDefinitions.Services, serviceDefinition)

	//yaml marshal services
	newServicesBytes, err := yaml.Marshal(yamlServiceDefinitions)
	if err != nil {
		return err
	}

	//write config into file
	if err := ioutil.WriteFile(config, newServicesBytes, 0644); err != nil {
		return err
	}

	return nil
}

func (writer *yamlServiceDefinitionWriter) UpdateServiceDefinition(config string, serviceDefinition ServiceDefinition) error {
	if serviceDefinition.Id == "" {
		return errors.New("ServiceDefinition Id can not empty")
	}

	rawConfig, err := ioutil.ReadFile(config)
	logrus.Debug("read Service config file ", string(rawConfig))

	if err != nil {
		return err
	}

	viper.SetConfigType("yml")
	err = viper.ReadConfig(bytes.NewBuffer(rawConfig))
	if err != nil {
		return err
	}

	var yamlServiceDefinitions struct {
		Services []ServiceDefinition `mapstructure:"services"`
	}

	if err := viper.Unmarshal(&yamlServiceDefinitions); err != nil {
		return err
	}

	//check sid is exists
	serviceExists := false
	for _, service := range yamlServiceDefinitions.Services {
		if service.Id == serviceDefinition.Id {
			serviceExists = true
			break
		}
	}

	if !serviceExists {
		return errors.New("service not exists")
	}

	//update Service config by Service id
	for i := range yamlServiceDefinitions.Services {
		if yamlServiceDefinitions.Services[i].Id == serviceDefinition.Id {
			yamlServiceDefinitions.Services[i] = serviceDefinition
		}
	}

	//yaml marshal services
	newServicesBytes, err := yaml.Marshal(yamlServiceDefinitions)
	if err != nil {
		return err
	}

	//write config into file
	if err := ioutil.WriteFile(config, newServicesBytes, 0644); err != nil {
		return err
	}

	return nil
}

func (writer *yamlServiceDefinitionWriter) DeleteServiceDefinition(config string, serviceDefinitionId string) error {
	rawConfig, err := ioutil.ReadFile(config)
	logrus.Debug("read Service config file ", string(rawConfig))

	if err != nil {
		return err
	}

	viper.SetConfigType("yml")
	err = viper.ReadConfig(bytes.NewBuffer(rawConfig))
	if err != nil {
		return err
	}

	var yamlServiceDefinitions struct {
		Services []ServiceDefinition `mapstructure:"services"`
	}

	if err := viper.Unmarshal(&yamlServiceDefinitions); err != nil {
		return err
	}

	//delete Service by Service id
	for i, service := range yamlServiceDefinitions.Services {
		if service.Id == serviceDefinitionId {
			yamlServiceDefinitions.Services = append(yamlServiceDefinitions.Services[:i], yamlServiceDefinitions.Services[i+1:]...)
			break
		}
	}

	//yaml marshal services
	newServicesBytes, err := yaml.Marshal(yamlServiceDefinitions)
	if err != nil {
		return err
	}

	//write config into file
	if err := ioutil.WriteFile(config, newServicesBytes, 0644); err != nil {
		return err
	}

	return nil
}
