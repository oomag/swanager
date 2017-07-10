package events

import (
	"io"

	"github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"github.com/dokkur/swanager/api/ws"
	"github.com/dokkur/swanager/core/swarm"
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
	log().Debug("Starting listening swarm events.")
	for {
		select {
		case message := <-messageChan:
			log().Debugf("Got docker event [%s] %s", message.Action, message.Actor.Attributes["com.docker.swarm.service.name"])
			// log().WithField("message", fmt.Sprintf("%+v", message)).Debug("Got docker event")

			ws.SendService(message.Actor.Attributes["com.docker.swarm.service.name"])

			break
		case err := <-errChan:
			if err == io.EOF {
				return
			}
			log().Debug("Docker error:")
			spew.Dump(err)
			break
		case <-quit:
			cancelListening()
			return
		}
	}
}

func log() *logrus.Entry {
	return logrus.WithField("module", "events")
}
