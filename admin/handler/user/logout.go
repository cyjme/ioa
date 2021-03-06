package user

import (
	"net/http"

	"github.com/cyjme/ioa/admin/handler"
	"github.com/cyjme/ioa/admin/pkg/errno"

	"github.com/sirupsen/logrus"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("user_id", "")

	logrus.Println("admin user logout.")
	
	handler.ResponseJson(w, "", errno.OK)
}