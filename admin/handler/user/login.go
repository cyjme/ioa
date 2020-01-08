package user

import (
	"net/http"

	"ioa/admin/config"
	"ioa/admin/handler"
	"ioa/admin/pkg/errno"

	"github.com/sirupsen/logrus"
)

type AdminLoginRequest struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	adminLoginRequest := AdminLoginRequest{}
	if err := handler.BindWith(r, &adminLoginRequest); err != nil {
		logrus.Println(errno.ErrBind.Msg)
		handler.ResponseJson(w, "", errno.ErrBind)
		return
	}

	if !(adminLoginRequest.Name == config.Config.Admin.Name && adminLoginRequest.Pass == config.Config.Admin.Pass) {
		logrus.Println(errno.AdminLoginError.Msg)
		handler.ResponseJson(w, "", errno.AdminLoginError)
		return
	}

	r.Header.Set("user_id", adminLoginRequest.Name)
	handler.ResponseJson(w, "", errno.OK)
}