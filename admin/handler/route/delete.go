package route

import (
	"net/http"

	"github.com/cyjme/ioa/route"
	"github.com/cyjme/ioa/admin/config"
	"github.com/cyjme/ioa/admin/handler"
	"github.com/cyjme/ioa/admin/pkg/errno"
	
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func DeleteRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	routeId := vars["id"]

	if routeId == "" {
		logrus.Println(errno.BadRequestError.Msg)
		handler.ResponseJson(w, "", errno.BadRequestError)
		return
	}

	routeWriter, _ := route.NewRouteWriter(config.Config.RoutesReader.RoutesReaderType, config.Config.RoutesReader.RoutesReaderArg)
	if err := routeWriter.DeleteRouteBy(routeId); err != nil {
		logrus.Println(errno.DeleteRouteError.Msg)
		handler.ResponseJson(w, "", errno.DeleteRouteError)
		return
	}

	handler.ResponseJson(w, "", errno.OK)
}
