package route

import (
	"net/http"

	"ioa/route"
	"ioa/admin/config"
	"ioa/admin/handler"
	"ioa/admin/pkg/errno"
	
	"github.com/sirupsen/logrus"
)

type UpdateRouteRequest struct {
	Id         string   `json:"id"`
	Uri        string   `json:"uri"`
	Filters    []string `json:"filters"`
	Predicates []string `json:"predicates"`
}

func UpdateRoute(w http.ResponseWriter, r *http.Request) {
	updateRouteRequest := UpdateRouteRequest{}
	if err := handler.BindWith(r, &updateRouteRequest); err != nil {
		logrus.Println(errno.ErrBind.Msg)
		handler.ResponseJson(w, "", errno.ErrBind)
		return
	}

	updateRouteDefinition := route.RouteDefinition{}
	updateRouteDefinition.Id = updateRouteRequest.Id
	updateRouteDefinition.Uri = updateRouteRequest.Uri
	updateRouteDefinition.Filters = updateRouteRequest.Filters
	updateRouteDefinition.Predicates = updateRouteRequest.Predicates

	routeWriter, _ := route.NewRouteWriter(config.Config.RoutesReader.RoutesReaderType, config.Config.RoutesReader.RoutesReaderArg)
	if err := routeWriter.UpdateRouteBy(updateRouteDefinition); err != nil {
		logrus.Println(errno.UpdateRouteError)
		handler.ResponseJson(w, "", errno.UpdateRouteError)
		return
	}

	handler.ResponseJson(w, "", errno.OK)
}
