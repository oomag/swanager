package swarm

import (
	"context"
	"fmt"

	"github.com/da4nik/swanager/core/entities"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// CreateNetwork creates swarm network, not working due to different version api and client
func CreateNetwork(name string) string {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	createOptions := types.NetworkCreate{Driver: "overlay"}

	// TODO: Check error if unable to create network, but not with duplication error
	response, _ := cli.NetworkCreate(context.Background(), name, createOptions)

	return response.ID
}

// RemoveNetwork removes swarm network
func RemoveNetwork(name string) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	return cli.NetworkRemove(context.Background(), name)
}

func getNetworkName(service *entities.Service) string {
	return fmt.Sprintf("%s_%s", service.Application.Name, service.Application.ID)
}
