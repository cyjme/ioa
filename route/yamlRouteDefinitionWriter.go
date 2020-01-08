package route

import (
	"bytes"
	"errors"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

type yamlRouteDefinitionWriter struct {
}

func (writer *yamlRouteDefinitionWriter) AddRouteDefinition(config string, routeDefinition RouteDefinition) error {
	rawConfig, err := ioutil.ReadFile(config)
	logrus.Debug("read config file ", string(rawConfig))

	if err != nil {
		return err
	}

	viper.SetConfigType("yml")
	err = viper.ReadConfig(bytes.NewBuffer(rawConfig))
	if err != nil {
		return err
	}

	var yamlConfig struct {
		GlobalFilters []string          `mapstructure:"globalFilters"`
		Routes        []RouteDefinition `mapstructure:"routes"`
	}

	if err := viper.Unmarshal(&yamlConfig); err != nil {
		return err
	}

	logrus.Debug("RouteDefinition in reader", yamlConfig.Routes)

	//check rid is repeat
	for _, route := range yamlConfig.Routes {
		if route.Id == routeDefinition.Id {
			return errors.New("route id repeat")
		}
	}

	//write new route into yaml config
	yamlConfig.Routes = append(yamlConfig.Routes, routeDefinition)

	//yaml marshal routes
	newRoutesBytes, err := yaml.Marshal(yamlConfig)
	if err != nil {
		return err
	}

	logrus.Println("new routes bytes:", string(newRoutesBytes))

	//write config into file
	if err := ioutil.WriteFile(config, newRoutesBytes, 0644); err != nil {
		return err
	}

	return nil
}

func (writer *yamlRouteDefinitionWriter) UpdateRouteDefinition(config string, routeDefinition RouteDefinition) error {
	if routeDefinition.Id == "" {
		return errors.New("RouteDefinition Id can not empty")
	}

	rawConfig, err := ioutil.ReadFile(config)
	logrus.Debug("read config file ", string(rawConfig))

	if err != nil {
		return err
	}

	viper.SetConfigType("yml")
	err = viper.ReadConfig(bytes.NewBuffer(rawConfig))
	if err != nil {
		return err
	}

	var yamlConfig struct {
		GlobalFilters []string          `mapstructure:"globalFilters"`
		Routes        []RouteDefinition `mapstructure:"routes"`
	}

	if err := viper.Unmarshal(&yamlConfig); err != nil {
		return err
	}

	logrus.Debug("RouteDefinition in reader", yamlConfig.Routes)

	//check rid is exists
	routeExists := false
	for _, route := range yamlConfig.Routes {
		if route.Id == routeDefinition.Id {
			routeExists = true
			break
		}
	}

	if !routeExists {
		return errors.New("route not exists")
	}

	//update route config by route id
	for i := range yamlConfig.Routes {
		if yamlConfig.Routes[i].Id == routeDefinition.Id {
			yamlConfig.Routes[i] = routeDefinition
		}
	}

	//yaml marshal
	newRoutesBytes, err := yaml.Marshal(yamlConfig)
	if err != nil {
		return err
	}

	//write config into file
	if err := ioutil.WriteFile(config, newRoutesBytes, 0644); err != nil {
		return err
	}

	return nil
}

func (writer *yamlRouteDefinitionWriter) DeleteRouteDefinition(config string, routDefinitionId string) error {
	rawConfig, err := ioutil.ReadFile(config)
	logrus.Debug("read config file ", string(rawConfig))

	if err != nil {
		return err
	}

	viper.SetConfigType("yml")
	err = viper.ReadConfig(bytes.NewBuffer(rawConfig))
	if err != nil {
		return err
	}

	var yamlConfig struct {
		GlobalFilters []string          `mapstructure:"globalFilters"`
		Routes        []RouteDefinition `mapstructure:"routes"`
	}

	if err := viper.Unmarshal(&yamlConfig); err != nil {
		return err
	}

	logrus.Debug("RouteDefinition in reader", yamlConfig.Routes)

	// delete route by route id
	for i, route := range yamlConfig.Routes {
		if route.Id == routDefinitionId {
			yamlConfig.Routes = append(yamlConfig.Routes[:i], yamlConfig.Routes[i+1:]...)
			break
		}
	}

	//yaml marshal routes
	newRoutesBytes, err := yaml.Marshal(yamlConfig)
	if err != nil {
		return err
	}

	//write config into file
	if err := ioutil.WriteFile(config, newRoutesBytes, 0644); err != nil {
		return err
	}

	return nil
}
