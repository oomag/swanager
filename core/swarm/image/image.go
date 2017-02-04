package image

import (
	"context"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/client"
)

// Volumes return image volumes
func Volumes(name string) (*[]string, error) {
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

	log.WithFields(log.Fields{
		"name":    name,
		"volumes": volumes,
	}).Debugf("Getting '%s' image volumes.", name)

	return &volumes, nil
}
