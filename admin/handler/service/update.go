package service

import (
	"net/http"
	
	"ioa/service"
	"ioa/admin/config"
	"ioa/admin/handler"
	"ioa/admin/pkg/errno"

	"github.com/sirupsen/logrus"
)

type UpdateServiceRequest struct {
	Qps  int     `json:"qps"`
	Id   string  `json:"id"`
	Uri  string  `json:"uri"`
	
	Backends []Backend `json:"backends"`
}

func UpdateService(w http.ResponseWriter, r *http.Request) {
	updateServiceRequest := UpdateServiceRequest{}
	if err := handler.BindWith(r, &updateServiceRequest); err != nil {
		logrus.Println(errno.ErrBind)
		handler.ResponseJson(w, "", errno.UpdateServiceError)
		return
	}

	updateServiceDefintion := service.ServiceDefinition{}
	updateServiceDefintion.Id = updateServiceRequest.Id
	updateServiceDefintion.Uri = updateServiceRequest.Uri
	updateServiceDefintion.Qps = updateServiceRequest.Qps

	serviceDefinitionBackend := service.Backend{}
	serviceDefinitionBackends := make([]service.Backend, 0)
	for _, backend := range updateServiceRequest.Backends {
		serviceDefinitionBackend.Id = backend.Id
		serviceDefinitionBackend.Uri = backend.Uri
		serviceDefinitionBackend.Qps = backend.Qps
		serviceDefinitionBackend.Weight = backend.Weight
		
		serviceDefinitionBackends = append(serviceDefinitionBackends, serviceDefinitionBackend)
	}
	updateServiceDefintion.Backends = serviceDefinitionBackends
	
	serviceWriter, _ := service.NewServiceWriter(config.Config.ServicesReader.ServicesReaderType, config.Config.ServicesReader.ServicesReaderArg)
	if err := serviceWriter.UpdateServiceBy(updateServiceDefintion); err != nil {
		logrus.Println(errno.UpdateServiceError.Msg)
		handler.ResponseJson(w, "", errno.UpdateServiceError)
		return
	}

	handler.ResponseJson(w, "", errno.OK)
}
