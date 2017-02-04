package main

import (
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
	"github.com/da4nik/swanager/api"
	"github.com/da4nik/swanager/config"
	"github.com/gin-gonic/gin"
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
	initLogger()

	// TODO: current there no way to gracefully stop net/http server,
	// this feature was added in go 1.8, waiting for release
	go api.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	logFile.Close()
}
