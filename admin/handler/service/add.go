package service

import (
	"net/http"
	
	"github.com/cyjme/ioa/service"
	"github.com/cyjme/ioa/admin/config"
	"github.com/cyjme/ioa/admin/handler"
	"github.com/cyjme/ioa/admin/pkg/errno"

	"github.com/sirupsen/logrus"
)

type Backend struct {
	Id  string `json:"id"`
	Uri string `json:"uri"`
	Qps int    `json:"qps"`

	Weight int `json:"weight"`
}

type AddServiceRequest struct {
	Id  string `json:"id"`
	Uri string `json:"uri"`
	Qps int    `json:"qps"`

	Backends []Backend `json:"backends"`
}

func AddService(w http.ResponseWriter, r *http.Request) {
	addServiceRequest := AddServiceRequest{}
	if err := handler.BindWith(r, &addServiceRequest); err != nil {
		logrus.Println(errno.ErrBind.Msg)
		handler.ResponseJson(w, "", errno.ErrBind)
		return
	}

	addServiceDefinition := service.ServiceDefinition{}
	addServiceDefinition.Id = addServiceRequest.Id
	addServiceDefinition.Uri = addServiceRequest.Uri
	addServiceDefinition.Qps = addServiceRequest.Qps

	serviceDefinitionBackend := service.Backend{}
	serviceDefinitionBackends := make([]service.Backend, 0)
	for _, backend := range addServiceRequest.Backends {
		serviceDefinitionBackend.Id = backend.Id
		serviceDefinitionBackend.Uri = backend.Uri
		serviceDefinitionBackend.Qps = backend.Qps
		serviceDefinitionBackend.Weight = backend.Weight
		
		serviceDefinitionBackends = append(serviceDefinitionBackends, serviceDefinitionBackend)
	}
	addServiceDefinition.Backends = serviceDefinitionBackends

	serviceWriter, _ := service.NewServiceWriter(config.Config.ServicesReader.ServicesReaderType, config.Config.ServicesReader.ServicesReaderArg)
	if err := serviceWriter.AddServiceBy(addServiceDefinition); err != nil {
		logrus.Println(errno.AddServiceError.Msg)
		handler.ResponseJson(w, "", errno.AddServiceError)
		return
	}

	handler.ResponseJson(w, "", errno.OK)
}
