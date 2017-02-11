package events

import (
	"io"

	"github.com/Sirupsen/logrus"
	"github.com/da4nik/swanager/api/ws"
	"github.com/da4nik/swanager/core/swarm"
	"github.com/davecgh/go-spew/spew"
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

			ws.NotifyServiceState(ws.NotifyServiceStateMessage{
				ServiceName: message.Actor.Attributes["com.docker.swarm.service.name"],
				Action:      message.Action,
			})

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
