package swarm

import (
	"github.com/da4nik/swanager/core/entities"
	"github.com/da4nik/swanager/core/swarm/network"
	swarm_service "github.com/da4nik/swanager/core/swarm/service"
)

// CreateService creates service and all stuff around it
func CreateService(service *entities.Service) {
	dockerServiceName := swarm_service.NameForDocker(service)

	network.Create(dockerServiceName)
	swarm_service.Create(swarm_service.CreateOptions{
		Service:     service,
		NetworkName: network.NameForDocker(service),
	})
}
