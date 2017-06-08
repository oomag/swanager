package command

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"

	"github.com/dokkur/swanager/core/entities"
	swarm_service "github.com/dokkur/swanager/core/swarm/service"
)

// ServiceLogs returns service logs
type ServiceLogs struct {
	CommonCommand

	User      *entities.User
	ServiceID string

	responseChan chan<- []string
}

// NewServiceLogsCommand created service start command
func NewServiceLogsCommand(command ServiceLogs) (ServiceLogs, chan []string, chan error) {
	response := make(chan []string, 1)
	err := make(chan error, 1)

	command.responseChan = response
	command.errorChan = err

	return command, response, err
}

// Process service logs
func (sl ServiceLogs) Process() {
	service, err := entities.GetService(bson.M{
		"user_id": sl.User.ID,
		"_id":     sl.ServiceID,
	})
	if err != nil {
		sl.errorChan <- err
		return
	}

	_, err = swarm_service.Inspect(service)
	if err != nil {
		sl.errorChan <- fmt.Errorf("Service stoppedtar")
		return
	}

	logs, err := swarm_service.Logs(service)
	if err != nil {
		sl.errorChan <- err
		return
	}

	sl.responseChan <- logs
}
