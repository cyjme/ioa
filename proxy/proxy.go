package proxy

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"ioa/config"
	"ioa/context"
	"ioa/filters/default/routeToRequestUrl"
	"ioa/plugin"
	"ioa/route"
	"ioa/service"

	"github.com/sirupsen/logrus"
)

type Proxy struct {
	routes         []route.Route
	services       *map[string]service.Service
	client         *http.Client
	defaultFilters []route.Filter
	certs          map[string]tls.Certificate
	pool           sync.Pool
	cfg            *Config
}

var accessLogger *AccessLogger

type Config struct {
	Addr               string `mapstructure:"addr"`
	HttpsAddr          string `mapstructure:"https_addr"`
	FilterPluginsDir   string `mapstructure:"filter_plugin_dir"`
	PredicatePluginDir string `mapstructure:"predicate_plugin_dir"`
	ConfigReaderType   string `mapstructure:"config_reader_type"`
	ConfigReaderArg    string `mapstructure:"config_reader_arg"`
	RoutesReaderType   string `mapstructure:"routes_reader_type"`
	RoutesReaderArg    string `mapstructure:"routes_reader_arg"`
	ServicesReaderType string `mapstructure:"services_reader_type"`
	ServicesReaderArg  string `mapstructure:"services_reader_arg"`

	MaxIdleConns        int    `mapstructure:"max_idle_conns"`
	MaxIdleConnsPerHost int    `mapstructure:"max_idle_conns_per_host"`
	AccessLogFile       string `mapstructure:"access_log_file"`
	EnableHttps         bool   `mapstructure:"enable_https"`
	RedirectHttps       bool   `mapstructure:"redirect_https"`

	AdminDomain string `mapstructure:"admin_domain"`
	AdminPort   string `mapstructure:"admin_port"`
	AdminName   string `mapstructure:"admin_name"`
	AdminPass   string `mapstructure:"admin_pass"`
}

func New(cfg Config) *Proxy {
	//fmt.Println("cfg reader type", cfg.ConfigReaderType)
	//fmt.Println("cfg reader arg", cfg.ConfigReaderArg)
	//err := config.InitConfig(cfg.ConfigReaderType, cfg.ConfigReaderArg)
	//if err != nil {
	//	logrus.Panic(err)
	//}

	plugin.Init(cfg.FilterPluginsDir)
	plugin.LoadPredicate(cfg.PredicatePluginDir)

	routes, err := route.GetRoutesBy(cfg.RoutesReaderType, cfg.RoutesReaderArg)
	if err != nil {
		logrus.Panic(err)
	}
	services, err := service.GetAllServicesBy(cfg.ServicesReaderType, cfg.ServicesReaderArg)
	if err != nil {
		logrus.Panic(err)
	}

	proxy := &Proxy{
		routes:         routes,
		services:       &services,
		defaultFilters: getDefaultFilters(),
		cfg:            &cfg,
	}

	proxy.client = proxy.initHttpClient()
	proxy.pool.New = func() interface{} {
		return &context.Context{}
	}

	if cfg.EnableHttps {
		proxy.certs = getCerts()
	}

	return proxy
}

func getCerts() map[string]tls.Certificate {
	certs := make(map[string]tls.Certificate)
	cfg := config.Config.List()
	for k, v := range cfg {
		if strings.HasPrefix(k, "ioa_tls_private") {
			host := strings.ReplaceAll(k, "ioa_tls_private_", "")
			public := config.Config.Get("ioa_tls_public_" + host)
			private := v
			if public == "" {
				logrus.Panic("can not found ioa_tls_public_" + host + "in config")
			}

			var err error
			certs[host], err = tls.LoadX509KeyPair(public, private)
			if err != nil {
				logrus.Panic("load tls cert err ", err)
			}
		}
	}

	return certs
}

func (p *Proxy) ReloadConfig(cfg Config) {
	err := config.InitConfig(cfg.ConfigReaderType, cfg.ConfigReaderArg)
	if err != nil {
		logrus.Panic(err)
	}

	routes, err := route.GetRoutesBy(cfg.RoutesReaderType, cfg.RoutesReaderArg)
	if err != nil {
		logrus.Panic(err)
	}
	services, err := service.GetAllServicesBy(cfg.ServicesReaderType, cfg.ServicesReaderArg)
	if err != nil {
		logrus.Panic(err)
	}

	p.routes = routes
	p.services = &services
}

func getDefaultFilters() []route.Filter {
	filters := make([]route.Filter, 0)

	routeToRequestFilter, err := routeToRequestUrl.New("")
	if err != nil {
		logrus.Panic(err)
	}

	filters = append(filters, routeToRequestFilter)

	return filters
}

type RedirectHttps struct {
	HttpsAddr string
}

func (redirectHttps *RedirectHttps) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	httpsPort := ""
	host := strings.Split(r.Host, ":")
	if strings.Contains(redirectHttps.HttpsAddr, ":") {
		domainPort := strings.Split(redirectHttps.HttpsAddr, ":")
		httpsPort = domainPort[1]
	}

	if httpsPort != "" {
		host[1] = httpsPort
	} else {
		host[1] = "443"
	}

	target := "https://" + strings.Join(host, ":") + r.URL.Path
	if len(r.URL.RawQuery) > 0 {
		target += "?" + r.URL.RawQuery
	}

	http.Redirect(w, r, target, http.StatusTemporaryRedirect)
}

func (p *Proxy) Run() {
	if p.cfg.AccessLogFile != "" {
		file := findOrCreateAccessLogFile(p.cfg.AccessLogFile)
		defer file.Close()
		accessLogger = NewAccessLogger()
		accessLogger.SetOutput(file)
	}
	logrus.Info("proxy Running")

	server := http.Server{
		Addr:    p.cfg.Addr,
		Handler: p,
	}

	if p.cfg.EnableHttps {
		if p.cfg.RedirectHttps {

			redirect := &RedirectHttps{
				HttpsAddr: p.cfg.HttpsAddr,
			}

			httpRedirectServer := http.Server{
				Addr:    p.cfg.Addr,
				Handler: redirect,
			}

			go func() {
				if err := httpRedirectServer.ListenAndServe(); err != nil {
					log.Panic(err)
				}
			}()
		}

		cfg := &tls.Config{
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				if cert, ok := p.certs[info.ServerName]; ok {
					return &cert, nil
				}
				return nil, errors.New("not exist ca")
			},
		}

		httpsServer := http.Server{
			Addr:    p.cfg.HttpsAddr,
			Handler: p,
		}
		httpsServer.TLSConfig = cfg
		err := httpsServer.ListenAndServeTLS("", "")
		if err != nil {
			logrus.Panic(err)
		}
		return
	}

	err := server.ListenAndServe()
	if err != nil {
		logrus.Panic(err)
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := p.pool.Get().(*context.Context)

	//ctx := &context.Context{}
	ctx.ResponseWriter = w
	ctx.Request = r
	ctx.Stop = false
	ctx.ServiceMap = p.services
	ctx.Client = p.client

	for _, route := range p.routes {
		matchSuccess := true
		for _, predicate := range route.Predicates {
			if !predicate.Apply(ctx) {
				matchSuccess = false
			}
		}

		if matchSuccess {
			ctx.MatchedRoute = true
			ctx.RouteUri = route.Uri
			ctx.RouteId = route.Id

			//todo jason, 解决类型转换的问题，因为 proxy route context 循环依赖, 考虑让proxy 不再依赖 route 包,但是 route 的匹配写在哪里？
			//把发送请求的动作如果也放到 context 中，那么又会产生 proxy 和 context 的互相依赖
			for _, filter := range p.defaultFilters {
				ctx.Filters = append(ctx.Filters, filter)
			}

			for _, filter := range route.Filters {
				ctx.Filters = append(ctx.Filters, filter)
			}

			ctx.Do()
		}
	}

	if !ctx.MatchedRoute {
		p.handlerNoMatchRoute(ctx)
	}

	ctx.Reset()
	p.pool.Put(ctx)
}

func (p *Proxy) handlerNoMatchRoute(context *context.Context) {
	context.ResponseWriter.WriteHeader(http.StatusNotFound)
}

func (p *Proxy) initHttpClient() *http.Client {
	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}

	// dereference it to get a copy of the struct that the pointer points to
	defaultTransport := defaultTransportPointer
	defaultTransport.MaxIdleConns = p.cfg.MaxIdleConns
	defaultTransport.MaxIdleConnsPerHost = p.cfg.MaxIdleConnsPerHost

	return &http.Client{Transport: defaultTransport}
}

func findOrCreateAccessLogFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			file, err = os.Create(path)

			if err != nil {
				logrus.Panic(err)
			}
		}
	}

	return file
}
