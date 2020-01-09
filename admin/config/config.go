package config

import (
	"github.com/cyjme/ioa/proxy"
)


var Config config

type config struct {
	Http struct {
		Domain string
		Port   string
	}
	Admin struct {
		Name  string
		Pass  string 
	}
	RoutesReader struct {
		RoutesReaderType string
		RoutesReaderArg  string
	}
	ServicesReader struct {
		ServicesReaderType string
		ServicesReaderArg  string
	}
}


func Init(cfg proxy.Config) {
	Config.RoutesReader.RoutesReaderArg = cfg.RoutesReaderArg
	Config.RoutesReader.RoutesReaderType = cfg.RoutesReaderType
	Config.ServicesReader.ServicesReaderArg = cfg.ServicesReaderArg
	Config.ServicesReader.ServicesReaderType = cfg.ServicesReaderType

	Config.Http.Domain = cfg.AdminDomain
	Config.Http.Port = cfg.AdminPort
	Config.Admin.Name = cfg.AdminName
	Config.Admin.Pass = cfg.AdminPass
}
