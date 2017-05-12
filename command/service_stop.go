package command

import (
	"github.com/da4nik/swanager/core/entities"
	"github.com/da4nik/swanager/core/swarm"
)

// ServiceStop stops service
type ServiceStop struct {
	User    *entities.User
	Service *entities.Service

	errorChan    chan<- error
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

	if err = swarm.StopService(ss.Service); err != nil {
		job.SetState(entities.JobStateError, "Error stoping service: "+err.Error())
		return
	}
	job.SetState(entities.JobStateSuccess, ss.Service)
}
