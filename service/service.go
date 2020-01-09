package service

import (
	"errors"
	"github.com/cyjme/ioa/pkg/roundRobin"

	"github.com/sirupsen/logrus"
)

type Service struct {
	Id  string
	Qps int
	// start with lb://  http:// https://
	// lb:// means loadBlanceEnable is true
	Uri string

	//if loadBlanceEnable is false, means uri is a url
	LoadBlanceEnable bool
	Backends         []Backend
	BackendMap       map[string]Backend

	roundrobin *roundRobin.RoundRobin
}

type ServiceWriter struct {
	Method string
	Config string
	Writer ServiceDefinitionWriter
}

func NewServiceWriter(method, config string) (*ServiceWriter, error) {
	var writer ServiceDefinitionWriter
	serviceWriter := ServiceWriter{}
	if method == "" || config == "" {
		return &serviceWriter, errors.New("params can not empty")
	}

	switch method {
	case "yaml":
		writer = &yamlServiceDefinitionWriter{}
	}

	serviceWriter.Method = method
	serviceWriter.Config = config
	serviceWriter.Writer = writer

	return &serviceWriter, nil
}

func (s *Service) GetUrl() string {
	if s.LoadBlanceEnable {
		backendId := s.roundrobin.Next()

		return s.BackendMap[backendId].Uri
	}

	return s.Uri
}

func GetAllServicesBy(method, config string) (map[string]Service, error) {
	servicesMap := make(map[string]Service)
	var serviceDefinitionReader serviceDefinitionReader
	switch method {
	case "yaml":
		serviceDefinitionReader = &yamlServiceDefinitionReader{}
	default:
		panic("method ServiceDefinitionReader not exist")
	}

	serviceDefinitions, err := serviceDefinitionReader.getServiceDefinitions(config)
	if err != nil {
		return servicesMap, err
	}

	for _, serviceDefinition := range serviceDefinitions {
		service := serviceDefinitionToService(serviceDefinition)
		servicesMap[service.Id] = service
	}

	return servicesMap, nil
}

func GetAllServiceDefinitionBy(method, config string) ([]ServiceDefinition, error) {
	var serviceDefinitionReader serviceDefinitionReader
	switch method {
	case "yaml":
		serviceDefinitionReader = &yamlServiceDefinitionReader{}
	default:
		panic("method ServiceDefinitionReader not exist")
	}

	serviceDefinitions, err := serviceDefinitionReader.getServiceDefinitions(config)
	if err != nil {
		return nil, err
	}

	return serviceDefinitions, nil
}

func serviceDefinitionToService(serviceDefinition ServiceDefinition) Service {
	service := Service{
		Id:         serviceDefinition.Id,
		Uri:        serviceDefinition.Uri,
		Qps:        serviceDefinition.Qps,
		Backends:   serviceDefinition.Backends,
		BackendMap: make(map[string]Backend),

		LoadBlanceEnable: len(serviceDefinition.Backends) != 0,
	}

	for _, backend := range serviceDefinition.Backends {
		service.BackendMap[backend.Id] = backend
	}

	rb := roundRobin.New()
	for _, backend := range service.Backends {
		rb.Add(backend.Id, backend.Weight)
	}
	service.roundrobin = rb

	return service
}

func (serviceWriter *ServiceWriter) AddServiceBy(serviceDefinition ServiceDefinition) error {
	err := serviceWriter.Writer.AddServiceDefinition(serviceWriter.Config, serviceDefinition)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (serviceWriter *ServiceWriter) UpdateServiceBy(serviceDefinition ServiceDefinition) error {
	err := serviceWriter.Writer.UpdateServiceDefinition(serviceWriter.Config, serviceDefinition)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (serviceWriter *ServiceWriter) DeleteServiceBy(serviceDefinitionId string) error {
	err := serviceWriter.Writer.DeleteServiceDefinition(serviceWriter.Config, serviceDefinitionId)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}
