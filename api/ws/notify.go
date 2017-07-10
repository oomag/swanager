package ws

import (
	"github.com/dokkur/swanager/core/entities"
	"github.com/dokkur/swanager/core/swarm"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// SendService notify user about state change
func SendService(serviceName string) {
	log().Debugf("Got service notification: %s", serviceName)
	service, err := entities.GetService(bson.M{"ns_name": serviceName})
	if err != nil {
		log().WithField("function", "NotifyService").Debugf("Got error: %s", err.Error())
		return
	}

	service.LoadApplication()
	if !IsUserConnected(service.Application.UserID) {
		return
	}
	swarm.GetServiceStatuses(service)

	Send(service.Application.UserID, gin.H{"service": service})
}
