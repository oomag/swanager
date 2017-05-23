package command

import (
	"github.com/da4nik/swanager/core/entities"
	"github.com/da4nik/swanager/core/swarm"
)

// ServiceStart starts service
type ServiceStart struct {
	CommonCommand

	User    *entities.User
	Service *entities.Service

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

	if err = swarm.StartService(ss.Service); err != nil {
		job.SetState(entities.JobStateError, "Error stoping service: "+err.Error())
		return
	}
	job.SetState(entities.JobStateSuccess, ss.Service)
}
