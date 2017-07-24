package node

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/dokkur/swanager/core/entities"
)

// List returns list of nodes
func List() ([]entities.Node, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	filters := filters.NewArgs()
	filters.Add("membership", "accepted")

	nodes, err := cli.NodeList(
		context.Background(),
		types.NodeListOptions{
			Filters: filters,
		})
	if err != nil {
		return nil, err
	}

	result := make([]entities.Node, 0)
	for _, node := range nodes {
		result = append(result, entities.Node{
			ID:            node.ID,
			Addr:          node.Status.Addr,
			Hostname:      node.Description.Hostname,
			Availability:  string(node.Spec.Availability),
			State:         string(node.Status.State),
			Role:          string(node.Spec.Role),
			NanoCPUs:      node.Description.Resources.NanoCPUs,
			MemoryBytes:   node.Description.Resources.MemoryBytes,
			EngineVersion: node.Description.Engine.EngineVersion,
		})
	}

	return result, nil
}
