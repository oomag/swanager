package swarm

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/dokkur/swanager/core/entities"
	"github.com/dokkur/swanager/core/swarm/network"
	swarm_service "github.com/dokkur/swanager/core/swarm/service"
)

// ServiceStatusNotExists - string represents service absense
const ServiceStatusNotExists = "not_exists"

// StartApplication starts a whole application
func StartApplication(app *entities.Application) error {
	app.LoadServices()

	networkName := network.NameForDocker(app)

	network.Create(networkName)

	for _, service := range app.Services {
		swarm_service.Create(swarm_service.SpecOptions{
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

	_, err := swarm_service.Create(swarm_service.SpecOptions{
		Service:     service,
		NetworkName: networkName,
	})

	return err
}

// UpdateService - updates running service
func UpdateService(service *entities.Service) error {
	if !ServiceExists(service) {
		return nil
	}

	serviceInspection, _ := swarm_service.Inspect(service)

	service.LoadApplication()
	networkName := network.NameForDocker(&service.Application)

	return swarm_service.Update(swarm_service.SpecOptions{
		Service:     service,
		NetworkName: networkName,
		Index:       serviceInspection.Meta.Version.Index,
	})
}

// StopService - stops/removes service
func StopService(service *entities.Service) (err error) {
	err = swarm_service.Remove(service)
	network.Prune()

	return
}

// ServiceExists - detects wether or not service exists in swarm
func ServiceExists(service *entities.Service) bool {
	GetServiceStatuses(service)

	// if service is not running, just return
	if len(service.Status) == 1 &&
		service.Status[0].Status == ServiceStatusNotExists {
		return false
	}
	return true
}

// GetServiceStatuses - loads service statused to service.Status field
func GetServiceStatuses(service *entities.Service) {
	states, err := swarm_service.Status(service)
	if err != nil {
		service.AddServiceStatus(entities.ServiceStatusStruct{
			Status: ServiceStatusNotExists,
		})
		return
	}

	for _, state := range states {
		service.AddServiceStatus(entities.ServiceStatusStruct{
			ReplicaID: state.TaskID,
			Node:      state.Node,
			Status:    state.Status,
			Timestamp: state.Timestamp,
			Error:     state.Error,
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
