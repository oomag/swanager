package command

import (
	"sync"

	"github.com/dokkur/swanager/core/entities"
	"github.com/dokkur/swanager/core/swarm"
	swarm_service "github.com/dokkur/swanager/core/swarm/service"
	"github.com/gin-gonic/gin"
)

// ServiceList returns services list
type ServiceList struct {
	CommonCommand

	User            *entities.User
	ApplicationID   string
	WithStatuses    bool
	WithVolumeSizes bool

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
	searchOptions := make(gin.H)

	if c.User != nil {
		searchOptions["user_id"] = c.User.ID
	}

	if c.ApplicationID != "" {
		searchOptions["application_id"] = c.ApplicationID
	}

	services, err := entities.GetServices(searchOptions)
	if err != nil {
		c.errorChan <- err
		return
	}

	if c.WithStatuses || c.WithVolumeSizes {
		var wg sync.WaitGroup
		for serviceIndex := range services {
			wg.Add(1)
			go func(service *entities.Service) {
				defer wg.Done()
				if c.WithStatuses {
					swarm.GetServiceStatuses(service)
				}

				if c.WithVolumeSizes {
					swarm_service.LoadVolumeSizes(service)
				}
			}(&services[serviceIndex])
		}
		wg.Wait()
	}

	c.responseChan <- services
}
