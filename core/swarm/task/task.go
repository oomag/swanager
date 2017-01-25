package task

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// ListFor returns tasks associated with service
func ListFor(serviceName string) (*[]swarm.Task, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	filter := filters.NewArgs()
	filter.Add("service", serviceName)

	tasks := make([]swarm.Task, 0)

	tasks, err = cli.TaskList(context.Background(), types.TaskListOptions{Filters: filter})
	if err != nil {
		return nil, err
	}

	return &tasks, err
}
