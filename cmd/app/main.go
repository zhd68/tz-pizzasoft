package main

import (
	"flag"
	"log"

	"github.com/zhd68/tz-pizzasoft/internal/app"
	"github.com/zhd68/tz-pizzasoft/internal/config"
	"github.com/zhd68/tz-pizzasoft/pkg/logging"
)

var cfgFlag = flag.String("config", "", "path to config file (if not indicated, used default test config)")

func main() {
	flag.Parse()
	cfgPathString := *cfgFlag

	cfg := config.New()
	if cfgPathString != "" {
		err := cfg.ParseConfig(cfgPathString)
		if err != nil {
			log.Fatal(err)
		}
	}

	logger := logging.New()
	if err := logger.SetConfig(cfg.LogLavel); err != nil {
		log.Fatal(err)
	}
	logger.Infoln("logger initialized")

	if cfgPathString == "" {
		logger.Infoln("used default config")
	} else {
		logger.Infoln("used config file:", cfgPathString)
	}
	logger.Infoln("lod level:", cfg.LogLavel)

	app, err := app.NewApp(cfg, logger)
	if err != nil {
		log.Fatal(err)
	}
	app.Run()
}
