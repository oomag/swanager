package network

import (
	"context"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/da4nik/swanager/core/entities"
	"github.com/da4nik/swanager/lib"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// Create creates swarm network, not working due to different version api and client
func Create(name string) string {
	log().Debugf("Creating network %s", name)
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	createOptions := types.NetworkCreate{Driver: "overlay"}

	// TODO: Check error if unable to create network, but not with duplication error
	response, err := cli.NetworkCreate(context.Background(), name, createOptions)
	if err != nil {
		log().Debugf("Unable to create network: %e", err)
		return ""
	}

	return response.ID
}

// Remove removes swarm network
func Remove(name string) error {
	log().Debugf("Removing network %s", name)
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	return cli.NetworkRemove(context.Background(), name)
}

// Prune removes unused networks
func Prune() error {
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
	name := fmt.Sprintf("%s-%s", lib.IdentifierName(app.Name), app.ID)[:32]
	return strings.Trim(name, " -")
}

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{"module": "swanager.Network"})
}
