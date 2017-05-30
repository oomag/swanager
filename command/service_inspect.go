package command

import (
	"github.com/dokkur/swanager/core/entities"
	"github.com/dokkur/swanager/core/swarm"
	swarm_service "github.com/dokkur/swanager/core/swarm/service"
	"github.com/gin-gonic/gin"
)

// ServiceInspect returns services list
type ServiceInspect struct {
	CommonCommand

	User      *entities.User
	ServiceID string

	responseChan chan<- entities.Service
}

// NewServiceInspectCommand inspects service
func NewServiceInspectCommand(command ServiceInspect) (ServiceInspect, chan entities.Service, chan error) {
	response := make(chan entities.Service, 1)
	err := make(chan error, 1)

	command.errorChan = err
	command.responseChan = response
	return command, response, err
}

// Process executes command
func (c ServiceInspect) Process() {
	searchOptions := gin.H{
		"user_id": c.User.ID,
		"_id":     c.ServiceID,
	}

	service, err := entities.GetService(searchOptions)
	if err != nil {
		c.errorChan <- err
		return
	}

	swarm.GetServiceStatuses(service)
	swarm_service.LoadVolumeSizes(service)

	c.responseChan <- *service
}
