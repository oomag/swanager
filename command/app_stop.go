package command

import (
	"github.com/dokkur/swanager/core/entities"
)

// AppStop - stops app
type AppStop struct {
	CommonCommand

	Application *entities.Application
	User        *entities.User

	responseChan chan<- entities.Job
}

// NewAppStopCommand creates app start command
func NewAppStopCommand(command AppStop) (AppStop, chan entities.Job, chan error) {
	response := make(chan entities.Job, 1)
	err := make(chan error, 1)

	command.responseChan = response
	command.errorChan = err

	return command, response, err
}

// Process stops application
func (appStop AppStop) Process() {
	job, err := entities.CreateJob(appStop.User)
	if err != nil {
		appStop.errorChan <- err
		return
	}
	appStop.responseChan <- *job

	app := appStop.Application
	app.LoadServices()

	for _, service := range app.Services {
		cmd, _, _ := NewServiceStopCommand(ServiceStop{
			User:    appStop.User,
			Service: service,
		})
		RunAsync(cmd)
	}

	// job.SetState(entities.JobStateError, "Error stoping service: "+err.Error())
	job.SetState(entities.JobStateSuccess, app)
}
