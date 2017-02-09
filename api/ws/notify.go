package ws

import (
	"github.com/da4nik/swanager/core/entities"
	"github.com/da4nik/swanager/core/swarm"
	"gopkg.in/mgo.v2/bson"
)

// NotifyServiceStateMessage params for NotifyServiceState
type NotifyServiceStateMessage struct {
	ServiceName string
	Action      string
}

// NotifyServiceState notify user about state change
func NotifyServiceState(message NotifyServiceStateMessage) {
	log().Debugf("Got service state notification: [%s] %s", message.Action, message.ServiceName)
	service, err := entities.GetService(bson.M{"ns_name": message.ServiceName})
	if err != nil {
		log().WithField("function", "NotifyServiceState").Debugf("Got error: %s", err.Error())
		return
	}

	service.LoadApplication()

	if clientContext, ok := clients[service.Application.UserID]; ok {
		swarm.GetServiceStatuses(service)
		clientContext.Incoming <- *service
	}
}
