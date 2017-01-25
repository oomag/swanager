package service

import (
	"context"
	"fmt"
	"time"

	"github.com/da4nik/swanager/config"
	"github.com/da4nik/swanager/core/entities"
	"github.com/da4nik/swanager/core/swarm/image"
	"github.com/da4nik/swanager/core/swarm/task"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// StatusStruct represents service state
type StatusStruct struct {
	Node      string
	Status    string
	Timestamp time.Time
}

// CreateOptions service create params
type CreateOptions struct {
	Service     *entities.Service
	NetworkName string
}

// Create creates swarm service form Service entity
func Create(opts CreateOptions) string {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	mounts, _ := getServiceMounts(opts.Service)

	containerSpec := swarm.ContainerSpec{
		Image:  opts.Service.Image,
		Mounts: mounts,
	}

	updateConfig := swarm.UpdateConfig{
		Parallelism:     opts.Service.Parallelism,
		FailureAction:   "pause",
		MaxFailureRatio: 0.5,
	}

	serviceSpec := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: NameForDocker(opts.Service),
			Labels: map[string]string{
				"swanager_id":    opts.Service.ID,
				"application_id": opts.Service.Application.Name,
			},
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: containerSpec,
			Networks: []swarm.NetworkAttachmentConfig{
				swarm.NetworkAttachmentConfig{Target: opts.NetworkName},
			},
		},
		Mode: swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{Replicas: opts.Service.Replicas},
		},
		UpdateConfig: &updateConfig,
	}

	serviceCreateOptions := types.ServiceCreateOptions{}

	responce, err := cli.ServiceCreate(context.Background(), serviceSpec, serviceCreateOptions)
	if err != nil {
		panic(err)
	}

	if len(responce.Warnings) > 0 {
		fmt.Println("Warnings:")
		fmt.Println(responce.Warnings)
	}

	return responce.ID
}

// Remove removes service
func Remove(service *entities.Service) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	err = cli.ServiceRemove(context.Background(), NameForDocker(service))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// Inspect return service status
func Inspect(service *entities.Service) (*swarm.Service, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	serviceInspection, _, err := cli.ServiceInspectWithRaw(context.Background(), NameForDocker(service))
	if err != nil {
		return nil, err
	}

	return &serviceInspection, nil
}

// Status returns service status
func Status(service *entities.Service) ([]StatusStruct, error) {
	tasks, err := task.ListFor(NameForDocker(service))
	if err != nil {
		return nil, err
	}

	result := make([]StatusStruct, 0)
	for _, task := range *tasks {
		result = append(result, StatusStruct{
			Node:      task.NodeID,
			Status:    string(task.Status.State),
			Timestamp: task.Status.Timestamp,
		})
	}

	return result, nil
}

// NameForDocker return service name for docker
func NameForDocker(service *entities.Service) string {
	return fmt.Sprintf("%s-%s-%s", service.Application.Name, service.Name, service.ID)
}

// getServiceMounts returns mount struct for creating new service
func getServiceMounts(service *entities.Service) ([]mount.Mount, error) {
	result := make([]mount.Mount, 0)
	vols, err := image.Volumes(service.Image)
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
