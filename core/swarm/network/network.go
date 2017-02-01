package network

import (
	"context"
	"fmt"
	"strings"

	"github.com/da4nik/swanager/core/entities"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// Create creates swarm network, not working due to different version api and client
func Create(name string) string {
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

// Remove removes swarm network
func Remove(name string) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	return cli.NetworkRemove(context.Background(), name)
}

// Prune removes unused networks
func Prune(name string) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	_, err = cli.NetworksPrune(context.Background(), filters.Args{})

	return err
}

// NameForDocker returns network name for docker
func NameForDocker(app *entities.Application) string {
	return fmt.Sprintf("%s_%s",
		strings.ToLower(app.Name),
		app.ID,
	)
}
