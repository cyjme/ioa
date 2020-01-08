package service

import (
	"net/http"
	
	"ioa/service"
	"ioa/admin/config"
	"ioa/admin/handler"
	"ioa/admin/pkg/errno"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func GetService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceId := vars["id"]

	if serviceId == "" {
		logrus.Println(errno.BadRequestError.Msg)
		handler.ResponseJson(w, "", errno.BadRequestError)
		return
	}

	listService, err := service.GetAllServiceDefinitionBy(config.Config.ServicesReader.ServicesReaderType, config.Config.ServicesReader.ServicesReaderArg)
	if err != nil {
		logrus.Println(errno.ListServiceError.Msg)
		handler.ResponseJson(w, "", errno.ListServiceError)
		return
	}

	exisits := false
	serviceRes := service.ServiceDefinition{}
	for _, value := range listService {
		if value.Id == serviceId {
			serviceRes = value
			exisits = true
			break
		}
	}

	if !exisits {
		logrus.Println(errno.NotExistsServiceError.Msg)
		handler.ResponseJson(w, "", errno.NotExistsServiceError)
		return
	}

	handler.ResponseJson(w, serviceRes, errno.OK)
}
