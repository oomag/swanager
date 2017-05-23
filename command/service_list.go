package command

import (
	"github.com/da4nik/swanager/core/entities"
	"github.com/da4nik/swanager/core/swarm"
	"github.com/gin-gonic/gin"
)

// ServiceList returns services list
type ServiceList struct {
	CommonCommand

	User          *entities.User
	ApplicationID string
	WithStatuses  bool

	responseChan chan<- []entities.Service
}

// NewServiceListCommand create command
func NewServiceListCommand(command ServiceList) (ServiceList, chan []entities.Service, chan error) {
	response := make(chan []entities.Service, 1)
	err := make(chan error, 1)

	command.errorChan = err
	command.responseChan = response
	return command, response, err
}

// Process executes command
func (c ServiceList) Process() {
	searchOptions := gin.H{"user_id": c.User.ID}
	if c.ApplicationID != "" {
		searchOptions["application_id"] = c.ApplicationID
	}

	services, err := entities.GetServices(searchOptions)
	if err != nil {
		c.errorChan <- err
		return
	}

	if c.WithStatuses {
		for serviceIndex := range services {
			swarm.GetServiceStatuses(&services[serviceIndex])
		}
	}

	c.responseChan <- services
}
