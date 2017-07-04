package ws

import (
	"github.com/dokkur/swanager/core/entities"
	"github.com/dokkur/swanager/core/swarm"
	"github.com/gin-gonic/gin"
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
	if !isUserConnected(service.Application.UserID) {
		return
	}
	swarm.GetServiceStatuses(service)

	sendNotification(service.Application.UserID, gin.H{"service": service})
}
