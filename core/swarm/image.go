package swarm

import (
	"context"

	"github.com/docker/docker/client"
)

func ImageVolumes(name string) (*[]string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	imageInfo, _, err := cli.ImageInspectWithRaw(context.Background(), name)
	if err != nil {
		return nil, err
	}

	volumes := make([]string, 0)
	for k := range imageInfo.Config.Volumes {
		volumes = append(volumes, k)
	}

	return &volumes, nil
}
