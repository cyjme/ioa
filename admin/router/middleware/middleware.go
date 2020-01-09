package middleware

import (
	"net/http"

	"github.com/cyjme/ioa/admin/config"
	"github.com/cyjme/ioa/admin/handler"
	"github.com/cyjme/ioa/admin/pkg/errno"

	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		logrus.Info(r)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func LoginCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check admin user is login
		userId := r.Header.Get("user_id")
		if !(userId == config.Config.Admin.Name) {
			logrus.Println(errno.AdminNoLoginError.Msg)
			handler.ResponseJson(w, "", errno.AdminNoLoginError)
			return
		}
		next.ServeHTTP(w, r)
	})
}