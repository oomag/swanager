package command

import (
	"github.com/dokkur/swanager/core/entities"
	"github.com/dokkur/swanager/core/swarm"
	"github.com/dokkur/swanager/core/swarm/service"
)

// ServiceStart starts service
type ServiceStart struct {
	CommonCommand

	User    *entities.User
	Service entities.Service

	responseChan chan<- entities.Job
}

// NewServiceStartCommand created service start command
func NewServiceStartCommand(command ServiceStart) (ServiceStart, chan entities.Job, chan error) {
	response := make(chan entities.Job, 1)
	err := make(chan error, 1)

	command.responseChan = response
	command.errorChan = err

	return command, response, err
}

// Process start service
func (ss ServiceStart) Process() {
	job, err := entities.CreateJob(ss.User)
	if err != nil {
		ss.errorChan <- err
		return
	}
	ss.responseChan <- *job

	startFunction := swarm.UpdateService

	// If service exists we need to update it instead of create
	_, err = service.Inspect(&ss.Service)
	if err != nil {
		startFunction = swarm.StartService
	}

	if err = startFunction(&ss.Service); err != nil {
		job.SetState(entities.JobStateError, "Error starting service: "+err.Error())
		return
	}

	job.SetState(entities.JobStateSuccess, ss.Service)

	RunAsync(FrontendUpdate{})
}
