package events

import (
	"io"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/da4nik/swanager/api/ws"
	"github.com/da4nik/swanager/core/entities"
	"github.com/da4nik/swanager/core/swarm"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

var quit = make(chan bool, 1)

// Start starts listening for docker events
func Start() {
	go listen()
}

// Stop stops listening
func Stop() {
	quit <- true
}

func listen() {
	messageChan, errChan, cancelListening := swarm.Events()
	for {
		select {
		case message := <-messageChan:
			parts := strings.Split(message.Actor.Attributes["com.docker.swarm.service.name"], "-")

			application, err := entities.GetApplication(gin.H{"name": parts[0]})
			if err != nil {
				continue
			}

			application.LoadServices()
			for serviceIndex := range application.Services {
				swarm.GetServiceStatuses(&application.Services[serviceIndex])
				ws.Notify(&application.Services[serviceIndex])
			}

			break
		case err := <-errChan:
			if err == io.EOF {
				return
			}
			logrus.Debug("Docker error:")
			spew.Dump(err)
			break
		case <-quit:
			cancelListening()
			return
		}
	}
}
