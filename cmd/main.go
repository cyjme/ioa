package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"ioa/proxy"

	adminConfig "ioa/admin/config"
	adminRouter "ioa/admin/router"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var configFile string

func main() {

	logrus.SetLevel(logrus.DebugLevel)
	//logrus.SetLevel(logrus.ErrorLevel)

	cfg, err := readConfigFromFile()
	if err != nil {
		logrus.Panic(err)
	}

	proxy := proxy.New(cfg)

	// admin
	adminConfig.Init(cfg)
	go adminRouter.Init()

	go watchSignal(proxy)
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	proxy.Run()

}

func readConfigFromFile() (proxy.Config, error) {
	cfg := proxy.Config{}

	flag.StringVar(&configFile, "configFile", "./config.yml", "config file")
	flag.Parse()

	rawConfig, err := ioutil.ReadFile(configFile)
	if err != nil {
		logrus.Panic("read config file error", err)
		return cfg, err
	}

	viper.SetConfigType("yaml")
	if err = viper.ReadConfig(bytes.NewBuffer(rawConfig)); err != nil {
		logrus.Panic("read config error", err)
		return cfg, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		logrus.Panic(err)
		return cfg, err
	}

	return cfg, nil
}

func watchSignal(p *proxy.Proxy) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP)
	for {
		s := <-signalChan
		if s == syscall.SIGHUP {
			cfg, err := readConfigFromFile()
			if err != nil {
				logrus.Error("read config file error", err)
			}
			p.ReloadConfig(cfg)
		}

		if s == syscall.SIGINT {
			os.Exit(9)
		}
	}
}
