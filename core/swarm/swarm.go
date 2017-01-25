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
