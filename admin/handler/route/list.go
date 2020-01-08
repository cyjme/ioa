package route

import (
	"net/http"
	"strconv"

	"ioa/admin/config"
	"ioa/admin/handler"
	"ioa/admin/pkg/errno"
	"ioa/route"

	"github.com/sirupsen/logrus"
)

func ListRoute(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	pageParam := query.Get("page")
	pageSizeParam := query.Get("pageSize")

	pageInt, _ := strconv.Atoi(pageParam)
	pageSizeInt, _ := strconv.Atoi(pageSizeParam)

	offset := pageInt*pageSizeInt - pageSizeInt
	limit := pageSizeInt

	if pageInt == 0 || pageSizeInt == 0 {
		limit = -1
	}

	listRoute, err := route.GetRouteDefinitionBy(config.Config.RoutesReader.RoutesReaderType, config.Config.RoutesReader.RoutesReaderArg)
	if err != nil {
		logrus.Println(errno.ListRouteError.Msg)
		handler.ResponseJson(w, "", errno.ListRouteError)
		return
	}
	logrus.Println("listRoute RouteDefinition:", listRoute)

	total := len(listRoute)
	result := make([]route.RouteDefinition, 0)
	if limit >= 0 {
		var end int
		if offset+limit < total {
			end = offset + limit
		} else {
			end = total
		}

		for k := offset; k < end; k++ {
			result = append(result, listRoute[k])
		}
	} else {
		result = listRoute
	}

	response := make(map[string]interface{})
	response["data"] = result
	response["total"] = len(listRoute)

	handler.ResponseJson(w, response, errno.OK)
}
