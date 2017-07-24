package command

import (
	"github.com/dokkur/swanager/core/entities"
	"github.com/dokkur/swanager/core/swarm"
	"github.com/dokkur/swanager/core/swarm/service"
)

// ServiceStop stops service
type ServiceStop struct {
	CommonCommand

	User    *entities.User
	Service entities.Service

	responseChan chan<- entities.Job
}

// NewServiceStopCommand created service start command
func NewServiceStopCommand(command ServiceStop) (ServiceStop, chan entities.Job, chan error) {
	response := make(chan entities.Job, 1)
	err := make(chan error, 1)

	command.responseChan = response
	command.errorChan = err

	return command, response, err
}

// Process start service
func (ss ServiceStop) Process() {
	job, err := entities.CreateJob(ss.User)
	if err != nil {
		ss.errorChan <- err
		return
	}
	ss.responseChan <- *job

	_, err = service.Inspect(&ss.Service)
	if err != nil {
		job.SetState(entities.JobStateSuccess, ss.Service)
		return
	}

	// If service exists, then this is correct error
	if err = swarm.StopService(&ss.Service); err != nil {
		job.SetState(entities.JobStateError, "Error stoping service: "+err.Error())
		return
	}
	job.SetState(entities.JobStateSuccess, ss.Service)
	RunAsync(FrontendUpdate{})
}
