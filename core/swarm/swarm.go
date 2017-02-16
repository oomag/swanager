package swarm

import (
	"context"

	"github.com/da4nik/swanager/core/entities"
	"github.com/da4nik/swanager/core/swarm/network"
	swarm_service "github.com/da4nik/swanager/core/swarm/service"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// StartApplication starts a whole application
func StartApplication(app *entities.Application) error {
	app.LoadServices()

	networkName := network.NameForDocker(app)

	network.Create(networkName)

	for _, service := range app.Services {
		swarm_service.Create(swarm_service.CreateOptions{
			Service:     &service,
			NetworkName: networkName,
		})
	}
	return nil
}

// StopApplication starts a whole application
func StopApplication(app *entities.Application) error {
	app.LoadServices()

	for _, service := range app.Services {
		swarm_service.Remove(&service)
	}

	network.Remove(network.NameForDocker(app))

	return nil
}

// StartService - starts individual service
func StartService(service *entities.Service) error {
	service.LoadApplication()

	networkName := network.NameForDocker(&service.Application)
	network.Create(networkName)

	_, err := swarm_service.Create(swarm_service.CreateOptions{
		Service:     service,
		NetworkName: networkName,
	})

	return err
}

// StopService - stops/removes service
func StopService(service *entities.Service) (err error) {
	err = swarm_service.Remove(service)
	network.Prune()

	return
}

// GetServiceStatuses - loads service statused to service.Status field
func GetServiceStatuses(service *entities.Service) {
	states, err := swarm_service.Status(service)
	if err != nil {
		service.AddServiceStatus(entities.ServiceStatusStruct{Status: "not_exists"})
		return
	}

	for _, state := range states {
		service.AddServiceStatus(entities.ServiceStatusStruct{
			ReplicaID: state.TaskID,
			Node:      state.Node,
			Status:    state.Status,
			Timestamp: state.Timestamp,
		})
	}
}

// Events - returns events and error channel for docker events
func Events() (eventChan <-chan events.Message, errChan <-chan error, cancelFunc func()) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	ctx, cancelFunc := context.WithCancel(context.Background())

	// TODO: Show only container events
	filters := filters.NewArgs()
	filters.Add("type", "container")

	eventChan, errChan = cli.Events(ctx, types.EventsOptions{Filters: filters})
	return
}
