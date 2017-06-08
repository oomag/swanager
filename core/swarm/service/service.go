package service

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/dokkur/swanager/config"
	"github.com/dokkur/swanager/core/entities"
	"github.com/dokkur/swanager/core/swarm/task"
)

// StatusStruct represents service state
type StatusStruct struct {
	TaskID    string
	Node      string
	Status    string
	Timestamp time.Time
	Error     string
}

// SpecOptions service create params
type SpecOptions struct {
	Service     *entities.Service
	NetworkName string
	Index       uint64
}

// Create creates swarm service form Service entity
func Create(opts SpecOptions) (string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	serviceSpec := getServiceSpec(opts)
	serviceCreateOptions := types.ServiceCreateOptions{}

	log().WithField("spec", fmt.Sprintf("%+v", serviceSpec)).Debug("Creating swarm service")

	responce, err := cli.ServiceCreate(context.Background(), serviceSpec, serviceCreateOptions)
	if err != nil {
		panic(err)
	}

	if len(responce.Warnings) > 0 {
		log().Debug("Warnings:")
		log().Debugf("%+v", responce.Warnings)
	}

	return responce.ID, nil
}

// Update - updates existing service
func Update(opts SpecOptions) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	serviceSpec := getServiceSpec(opts)
	serviceUpdateOptions := types.ServiceUpdateOptions{}

	responce, err := cli.ServiceUpdate(context.Background(), opts.Service.NSName, swarm.Version{Index: opts.Index}, serviceSpec, serviceUpdateOptions)
	if err != nil {
		panic(err)
	}

	if len(responce.Warnings) > 0 {
		log().Debug("Warnings:")
		log().Debugf("%+v", responce.Warnings)
	}

	return nil
}

// Remove removes service
func Remove(service *entities.Service) error {
	log().Debugf("Removing service. %s", service.NSName)
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	defer cli.Close()

	err = cli.ServiceRemove(context.Background(), service.NSName)
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

	serviceInspection, _, err := cli.ServiceInspectWithRaw(context.Background(), service.NSName, types.ServiceInspectOptions{InsertDefaults: true})
	if err != nil {
		return nil, err
	}

	return &serviceInspection, nil
}

// Status returns service status
func Status(service *entities.Service) ([]StatusStruct, error) {
	tasks, err := task.ListFor(service.NSName)
	if err != nil {
		return nil, err
	}

	result := make([]StatusStruct, 0)
	for _, task := range *tasks {
		result = append(result, StatusStruct{
			TaskID:    task.ID,
			Node:      task.NodeID,
			Status:    string(task.Status.State),
			Timestamp: task.Status.Timestamp,
			Error:     task.Status.Err,
		})
	}
	return result, nil
}

// Logs return service logs
func Logs(service *entities.Service) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	reader, err := cli.ServiceLogs(ctx, service.NSName, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Follow:     false,
	})
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(reader)
	result := make([]string, 0)

	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return result, nil
}

func getServiceSpec(opts SpecOptions) swarm.ServiceSpec {
	opts.Service.LoadApplication()

	mounts := getServiceVolumes(opts.Service)

	containerSpec := swarm.ContainerSpec{
		Image:   opts.Service.Image,
		Mounts:  mounts,
		Env:     prepareEnvVars(opts.Service),
		Command: prepareCommand(opts.Service),
	}

	updateConfig := swarm.UpdateConfig{
		Parallelism:     opts.Service.Parallelism,
		FailureAction:   "pause",
		MaxFailureRatio: 0.5,
	}

	return swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: opts.Service.NSName,
			Labels: map[string]string{
				"swanager_id":    opts.Service.ID,
				"application_id": opts.Service.Application.ID,
			},
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: containerSpec,
			Resources: &swarm.ResourceRequirements{
				Limits: &swarm.Resources{
					NanoCPUs:    0, // CPU ratio * 10^9 :)
					MemoryBytes: 0, // in bytes
				},
			},
			Networks: []swarm.NetworkAttachmentConfig{
				swarm.NetworkAttachmentConfig{Target: opts.NetworkName},
			},
		},
		Mode: swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{Replicas: opts.Service.Replicas},
		},
		UpdateConfig: &updateConfig,
		EndpointSpec: &swarm.EndpointSpec{
			Mode:  swarm.ResolutionModeVIP,
			Ports: preparePorts(opts.Service),
		},
	}
}

// Update service in db with currently running service params
// e.g. Autoassigned published ports
func updateWithRunningSpec(service *entities.Service) error {
	running, err := Inspect(service)
	if err != nil {
		return err
	}

	// runningPorts := running.
	spew.Dump(running)

	return nil
}

func getServiceVolumes(service *entities.Service) []mount.Mount {
	service.LoadApplication()

	result := make([]mount.Mount, 0)
	for _, vol := range service.Volumes {
		sourcePath := getMountPath(service, vol)
		os.MkdirAll(sourcePath, 0777)

		result = append(result, mount.Mount{
			Type:     mount.TypeBind,
			Source:   sourcePath,
			Target:   vol.Service,
			ReadOnly: false,
		})
	}
	return result
}

func prepareCommand(service *entities.Service) []string {
	if service.Command == "" {
		return make([]string, 0)
	}
	return regexp.MustCompile("\\s+").Split(service.Command, -1)
}

func prepareEnvVars(service *entities.Service) (vars []string) {
	for _, envVar := range service.EnvVariables {
		vars = append(vars, fmt.Sprintf("%s=%s", envVar.Name, envVar.Value))
	}
	return
}

func preparePorts(service *entities.Service) (ports []swarm.PortConfig) {

	for _, port := range service.PublishedPorts {
		// Don't publish port if disabled
		if port.Disabled {
			continue
		}

		ports = append(ports, swarm.PortConfig{
			Name:          "swanager_port",
			Protocol:      stringToProtocol(port.Protocol),
			TargetPort:    port.Internal,
			PublishedPort: port.External,
			PublishMode:   swarm.PortConfigPublishModeIngress,
		})
	}
	return
}

func stringToProtocol(protocol string) swarm.PortConfigProtocol {
	switch strings.ToLower(protocol) {
	case "udp":
		return swarm.PortConfigProtocolUDP
	}
	return swarm.PortConfigProtocolTCP
}

func getMountPathPrefix(service *entities.Service, appWide bool) string {
	path := service.NSName
	if appWide {
		path = "app_wide"
	}

	return filepath.Join(config.MountPathPrefix, service.ApplicationID, path)
}

func getMountPath(service *entities.Service, vol entities.ServiceVolume) string {
	path := vol.Service
	if vol.Backend != "" {
		path = vol.Backend
	}

	return filepath.Join(getMountPathPrefix(service, vol.AppWide), path)
}

func log() *logrus.Entry {
	return logrus.WithField("module", "swarm.service")
}

// LoadVolumeSizes loads volume sized into struct
func LoadVolumeSizes(service *entities.Service) {
	var wg sync.WaitGroup
	for index := range service.Volumes {
		wg.Add(1)

		go dirSize(service, &service.Volumes[index], &wg)
	}
	wg.Wait()
}

func dirSize(service *entities.Service, vol *entities.ServiceVolume, wg *sync.WaitGroup) {
	defer wg.Done()
	var size int64

	root := getMountPath(service, *vol)

	filepath.Walk(root, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	vol.Size = size
}
