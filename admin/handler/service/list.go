package service

import (
	"net/http"
	"strconv"

	"github.com/cyjme/ioa/admin/config"
	"github.com/cyjme/ioa/admin/handler"
	"github.com/cyjme/ioa/admin/pkg/errno"
	"github.com/cyjme/ioa/service"

	"github.com/sirupsen/logrus"
)

func ListService(w http.ResponseWriter, r *http.Request) {
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

	listService, err := service.GetAllServiceDefinitionBy(config.Config.ServicesReader.ServicesReaderType, config.Config.ServicesReader.ServicesReaderArg)
	if err != nil {
		logrus.Println(errno.ListServiceError.Msg)
		handler.ResponseJson(w, "", errno.ListServiceError)
		return
	}

	total := len(listService)
	result := make([]service.ServiceDefinition, 0)
	if limit >= 0 {
		var end int
		if offset+limit < total {
			end = offset + limit
		} else {
			end = total
		}

		for k := offset; k < end; k++ {
			result = append(result, listService[k])
		}
	} else {
		result = listService
	}

	response := make(map[string]interface{})
	response["data"] = result
	response["total"] = len(listService)

	handler.ResponseJson(w, response, errno.OK)
}
