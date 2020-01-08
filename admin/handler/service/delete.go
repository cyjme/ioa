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

func DeleteService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceId := vars["id"]

	if serviceId == "" {
		logrus.Println(errno.BadRequestError.Msg)
		handler.ResponseJson(w, "", errno.BadRequestError)
		return
	}

	serviceWriter, _ := service.NewServiceWriter(config.Config.ServicesReader.ServicesReaderType, config.Config.ServicesReader.ServicesReaderArg)
	if err := serviceWriter.DeleteServiceBy(serviceId); err != nil {
		logrus.Println(errno.DeleteServiceError.Msg)
		handler.ResponseJson(w, "", errno.DeleteServiceError)
		return
	}

	handler.ResponseJson(w, "", errno.OK)
}
