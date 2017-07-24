package main

import (
	"fmt"
	"os"
	"os/signal"

	_ "github.com/joho/godotenv/autoload"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/dokkur/swanager/api"
	"github.com/dokkur/swanager/config"
	"github.com/dokkur/swanager/core"
	"github.com/dokkur/swanager/events"
	"github.com/dokkur/swanager/frontend"
)

var logFile *os.File

func initLogger() {
	// Setting up logger
	if gin.Mode() == gin.DebugMode {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.InfoLevel)
	}

	log.SetOutput(os.Stdout)
	if len(config.LogFileName) > 0 {
		var err error
		logFile, err = os.OpenFile(config.LogFileName, os.O_WRONLY|os.O_CREATE, 0664)
		if err != nil {
			log.Warningf("File %s, can't be opened, using STDOUT for logging.", config.LogFileName)
		} else {
			log.SetOutput(logFile)
		}
	}
}

func main() {
	if core.Version != "" || core.BuildTime != "" {
		fmt.Printf("Swanager build version %s, build time %s\n\n", core.Version, core.BuildTime)
	}

	config.Init()

	initLogger()
	defer logFile.Close()

	// TODO: current there no way to gracefully stop net/http server,
	// this feature was added in go 1.8, waiting for release
	go api.Start()

	frontend.Init()

	events.Start()
	defer events.Stop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
