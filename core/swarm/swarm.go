package swarm

import (
	"github.com/da4nik/swanager/core/entities"
	"github.com/da4nik/swanager/core/swarm/network"
	swarm_service "github.com/da4nik/swanager/core/swarm/service"
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
