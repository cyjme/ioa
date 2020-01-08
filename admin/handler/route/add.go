package route

import (
	"net/http"

	"ioa/route"
	"ioa/admin/config"
	"ioa/admin/handler"
	"ioa/admin/pkg/errno"

	"github.com/sirupsen/logrus"
)

type AddRouteRequest struct {
	Id         string   `json:"id"`
	Uri        string   `json:"uri"`
	Filters    []string `json:"filters"`
	Predicates []string `json:"predicates"`
}

func AddRoute(w http.ResponseWriter, r *http.Request) {
	addRouteRequest := AddRouteRequest{}
	if err := handler.BindWith(r, &addRouteRequest); err != nil {
		logrus.Println(errno.ErrBind.Msg)
		handler.ResponseJson(w, "", errno.ErrBind)
		return
	}

	addRouteDefinition := route.RouteDefinition{}
	addRouteDefinition.Id = addRouteRequest.Id
	addRouteDefinition.Uri = addRouteRequest.Uri
	addRouteDefinition.Filters = addRouteRequest.Filters
	addRouteDefinition.Predicates = addRouteRequest.Predicates

	routeWriter, _ := route.NewRouteWriter(config.Config.ServicesReader.ServicesReaderType, config.Config.ServicesReader.ServicesReaderArg)
	if err := routeWriter.AddRouteBy(addRouteDefinition); err != nil {
		logrus.Println(errno.AddRouteError.Msg)
		handler.ResponseJson(w, "", errno.AddRouteError)
		return
	}

	handler.ResponseJson(w, "", errno.OK)
}
