package swarm

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// CreateNetwork creates swarm network, not working due to different version api and client
func CreateNetwork(networkType string, name string) string {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	response, err := cli.NetworkCreate(context.Background(), name, types.NetworkCreate{})
	if err != nil {
		panic(err)
	}

	return response.ID
}
