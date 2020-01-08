package route

import (
	"bytes"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type yamlRouteDefinitionReader struct {
}

func (reader *yamlRouteDefinitionReader) GetRouteDefinitions(config string) ([]RouteDefinition, error) {
	routeDefinitions := make([]RouteDefinition, 0)
	rawConfig, err := ioutil.ReadFile(config)
	logrus.Debug("read config file ", string(rawConfig))

	if err != nil {
		logrus.Panic(err)
		return routeDefinitions, err
	}

	viper.SetConfigType("yml")
	err = viper.ReadConfig(bytes.NewBuffer(rawConfig))
	if err != nil {
		return routeDefinitions, err
	}

	var yamlConfig struct {
		GlobalFilters []string          `mapstructure:"globalFilters"`
		Routes        []RouteDefinition `mapstructure:"routes"`
	}

	if err := viper.Unmarshal(&yamlConfig); err != nil {
		logrus.Panic(err)
		return routeDefinitions, err
	}

	logrus.Debug("RouteDefinition in reader:", yamlConfig.Routes)

	//if has globalFilters, add globalFilter to every RouteDefinition
	if len(yamlConfig.GlobalFilters) > 0 {
		newRoutes := make([]RouteDefinition, 0)
		for _, route := range yamlConfig.Routes {
			route.Filters = append(yamlConfig.GlobalFilters, route.Filters...)
			logrus.Debug("route filters", route.Filters)
			newRoutes = append(newRoutes, route)
		}
		yamlConfig.Routes = newRoutes
	}

	// check is repeat
	existsMap := make(map[string]bool)
	for _, route := range yamlConfig.Routes {
		if _, ok := existsMap[route.Id]; ok {
			logrus.Panic("route id repeat")
		}
		existsMap[route.Id] = true
	}

	return yamlConfig.Routes, nil
}
