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

func GetRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	routeId := vars["id"]

	if routeId == "" {
		logrus.Println(errno.BadRequestError.Msg)
		handler.ResponseJson(w, "", errno.BadRequestError)
		return
	}

	listRoute, err := route.GetRouteDefinitionBy(config.Config.RoutesReader.RoutesReaderType, config.Config.RoutesReader.RoutesReaderArg)
	if err != nil {
		logrus.Println(errno.ListRouteError.Msg)
		handler.ResponseJson(w, "", errno.ListRouteError)
		return
	}

	exisits := false
	routeRes := route.RouteDefinition{}
	for _, value := range listRoute {
		if value.Id == routeId {
			routeRes = value
			exisits = true
			break
		}
	}

	if !exisits {
		logrus.Println(errno.NotExistsRouteError.Msg)
		handler.ResponseJson(w, "", errno.NotExistsRouteError)
		return
	}

	handler.ResponseJson(w, routeRes, errno.OK)
}
