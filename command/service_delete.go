package command

import (
	"github.com/dokkur/swanager/core/entities"
	swarm_service "github.com/dokkur/swanager/core/swarm/service"
	"github.com/gin-gonic/gin"
)

// ServiceDelete deleted service
type ServiceDelete struct {
	CommonCommand

	User      *entities.User
	ServiceID string

	responseChan chan<- entities.Service
}

// NewServiceDeleteCommand create command
func NewServiceDeleteCommand(command ServiceDelete) (ServiceDelete, chan entities.Service, chan error) {
	response := make(chan entities.Service, 1)
	err := make(chan error, 1)

	command.errorChan = err
	command.responseChan = response
	return command, response, err
}

// Process deletes command
func (sd ServiceDelete) Process() {
	service, err := entities.GetService(gin.H{
		"_id":     sd.ServiceID,
		"user_id": sd.User.ID,
	})

	if err != nil {
		sd.errorChan <- err
		return
	}

	swarm_service.Remove(service)
	if err := service.Delete(); err != nil {
		sd.errorChan <- err
		return
	}

	sd.responseChan <- *service
	RunAsync(FrontendUpdate{})
}
