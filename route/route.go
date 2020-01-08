package route

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
)

type Route struct {
	Id         string
	Uri        string
	Filters    []Filter
	Predicates []Predicate
}

type RouteWriter struct {
	Method string
	Config string
	Writer RouteDefinitionWriter
}

func NewRouteWriter(method, config string) (*RouteWriter, error) {
	var writer RouteDefinitionWriter
	routeWriter := RouteWriter{}
	if method == "" || config == "" {
		return &routeWriter, errors.New("params can not empty")
	}

	switch method {
	case "yaml":
		writer = &yamlRouteDefinitionWriter{}
	}

	routeWriter.Method = method
	routeWriter.Config = config
	routeWriter.Writer = writer

	return &routeWriter, nil
}

func routeDefinitionToRoute(rd RouteDefinition) Route {
	route := Route{}

	route.Id = rd.Id
	route.Uri = rd.Uri

	for _, predicateString := range rd.Predicates {
		index := strings.Index(predicateString, "=")
		name := predicateString[:index]
		arg := predicateString[index+1:]

		predicate, err := createPredicateByName(name, arg)
		if err != nil {
			logrus.Errorf("load %s error:%s", name, err.Error())
			continue
		}

		route.Predicates = append(route.Predicates, predicate)
	}

	for _, filterString := range rd.Filters {
		index := strings.Index(filterString, "=")
		name := filterString[:index]
		arg := filterString[index+1:]

		filter, err := createFilterByName(name, arg)
		if err != nil {
			logrus.Error("createFilter" + name + " error:  " + err.Error())

			continue
		}

		route.Filters = append(route.Filters, filter)
	}

	return route
}

func GetRoutesBy(method, config string) ([]Route, error) {
	var reader routeDefinitionReader
	var routes = make([]Route, 0)
	switch method {
	case "yaml":
		reader = &yamlRouteDefinitionReader{}
	case "etcd":
		reader = &etcdRouteDefinitionReader{}
	}

	routeDefinitions, err := reader.GetRouteDefinitions(config)
	if err != nil {
		logrus.Error(err)
		return routes, err
	}

	for _, rd := range routeDefinitions {
		route := routeDefinitionToRoute(rd)
		routes = append(routes, route)
	}

	return routes, nil
}

func GetRouteDefinitionBy(method, config string) ([]RouteDefinition, error) {
	var reader routeDefinitionReader
	switch method {
	case "yaml":
		reader = &yamlRouteDefinitionReader{}
	case "etcd":
		reader = &etcdRouteDefinitionReader{}
	}

	routeDefinitions, err := reader.GetRouteDefinitions(config)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return routeDefinitions, nil
}


func (routeWriter *RouteWriter) AddRouteBy(routDefinition RouteDefinition) error {
	err := routeWriter.Writer.AddRouteDefinition(routeWriter.Config, routDefinition)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (routeWriter *RouteWriter) UpdateRouteBy(routDefinition RouteDefinition) error {
	err := routeWriter.Writer.UpdateRouteDefinition(routeWriter.Config, routDefinition)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (routeWriter *RouteWriter) DeleteRouteBy(routDefinitionId string) error {
	err := routeWriter.Writer.DeleteRouteDefinition(routeWriter.Config, routDefinitionId)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}