package swarm

import (
	"context"
	"fmt"

	"github.com/da4nik/swanager/config"
	"github.com/da4nik/swanager/core/entities"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// ServiceCreate creates swarm service form Service entity
func ServiceCreate(service *entities.Service) string {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	CreateNetwork(getNetworkName(service))

	mounts, _ := getServiceMounts(service)

	containerSpec := swarm.ContainerSpec{
		Image:  service.Image,
		Mounts: mounts,
	}

	updateConfig := swarm.UpdateConfig{
		Parallelism:     service.Parallelism,
		FailureAction:   "pause",
		MaxFailureRatio: 0.5,
	}

	serviceSpec := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: dockerServiceName(service),
			Labels: map[string]string{
				"swanager_id":    service.ID,
				"application_id": service.Application.Name,
			},
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: containerSpec,
			Networks: []swarm.NetworkAttachmentConfig{
				swarm.NetworkAttachmentConfig{Target: getNetworkName(service)},
			},
		},
		Mode: swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{Replicas: service.Replicas},
		},
		UpdateConfig: &updateConfig,
	}

	serviceCreateOptions := types.ServiceCreateOptions{}

	responce, err := cli.ServiceCreate(context.Background(), serviceSpec, serviceCreateOptions)
	if err != nil {
		panic(err)
	}

	if len(responce.Warnings) > 0 {
		fmt.Println("Wranings:")
		fmt.Println(responce.Warnings)
	}

	return responce.ID
}

// ServiceRemove removes service
func ServiceRemove(service *entities.Service) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	err = cli.ServiceRemove(context.Background(), dockerServiceName(service))
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = cli.NetworksPrune(context.Background(), filters.Args{})
	return err
}

func dockerServiceName(service *entities.Service) string {
	return fmt.Sprintf("%s-%s-%s", service.Application.Name, service.Name, service.ID)
}

func getServiceMounts(service *entities.Service) ([]mount.Mount, error) {
	result := make([]mount.Mount, 0)
	vols, err := ImageVolumes(service.Image)
	if err != nil {
		return result, err
	}

	volumeOptions := mount.VolumeOptions{
		Labels: map[string]string{
			"application_id": service.Application.ID,
			"service_id":     service.ID,
		},
	}

	for _, vol := range *vols {
		result = append(result, mount.Mount{
			Type:          "bind",
			Source:        getMountPathPrefix() + vol,
			Target:        vol,
			VolumeOptions: &volumeOptions,
		})
	}

	return result, nil
}

func getMountPathPrefix() string {
	return config.MountPathPrefix
}
