package swarm

import (
	"context"

	"github.com/da4nik/swanager/core/entities"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// GetTasks returns tasks associated with service
func GetTasks(service *entities.Service) (*[]swarm.Task, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	filter := filters.NewArgs()
	filter.Add("service", dockerServiceName(service))

	tasks := make([]swarm.Task, 0)

	tasks, err = cli.TaskList(context.Background(), types.TaskListOptions{Filters: filter})
	if err != nil {
		return nil, err
	}

	return &tasks, err
}
