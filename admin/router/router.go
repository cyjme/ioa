package router

import (
	"net/http"
	"time"

	"github.com/cyjme/ioa/admin/config"
	"github.com/cyjme/ioa/admin/handler/route"
	"github.com/cyjme/ioa/admin/handler/service"
	"github.com/cyjme/ioa/admin/handler/user"
	"github.com/cyjme/ioa/admin/router/middleware"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Init 路由初始化
func Init() {
	r := mux.NewRouter()
	// 日志中间件
	r.Use(middleware.LoggingMiddleware)

	// 登录检查中间件
	r.Use(middleware.LoginCheckMiddleware)

	// admin login
	r.HandleFunc("/admin/login", user.Login).Methods("post")
	r.HandleFunc("/admin/logout", user.Logout).Methods("post")

	// admin operate route
	r.HandleFunc("/admin/route", route.AddRoute).Methods("post")
	r.HandleFunc("/admin/route", route.ListRoute).Methods("get")
	r.HandleFunc("/admin/route", route.UpdateRoute).Methods("put")
	r.HandleFunc("/admin/route/{id}", route.GetRoute).Methods("get")
	r.HandleFunc("/admin/route/{id}", route.DeleteRoute).Methods("delete")

	// admin operate service
	r.HandleFunc("/admin/service", service.AddService).Methods("post")
	r.HandleFunc("/admin/service", service.ListService).Methods("get")
	r.HandleFunc("/admin/service", service.UpdateService).Methods("put")
	r.HandleFunc("/admin/service/{id}", service.GetService).Methods("get")
	r.HandleFunc("/admin/service/{id}", service.DeleteService).Methods("delete")

	addr := config.Config.Http.Domain + ":" + config.Config.Http.Port
	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logrus.Info("admin api runing on: ", addr)
	logrus.Fatal(srv.ListenAndServe())
}
