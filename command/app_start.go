package command

import (
	"github.com/dokkur/swanager/core/entities"
)

// AppStart - starts app
type AppStart struct {
	CommonCommand

	Application *entities.Application
	User        *entities.User

	responseChan chan<- entities.Job
}

// NewAppStartCommand creates app start command
func NewAppStartCommand(command AppStart) (AppStart, chan entities.Job, chan error) {
	response := make(chan entities.Job, 1)
	err := make(chan error, 1)

	command.responseChan = response
	command.errorChan = err

	return command, response, err
}

// Process starts application
func (appStart AppStart) Process() {
	job, err := entities.CreateJob(appStart.User)
	if err != nil {
		appStart.errorChan <- err
		return
	}
	appStart.responseChan <- *job

	app := appStart.Application
	app.LoadServices()

	for _, service := range app.Services {
		cmd, _, _ := NewServiceStartCommand(ServiceStart{
			User:    appStart.User,
			Service: service,
		})
		RunAsync(cmd)
	}

	// job.SetState(entities.JobStateError, "Error stoping service: "+err.Error())
	job.SetState(entities.JobStateSuccess, app)
}
