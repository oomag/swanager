package entities

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"hash/crc32"

	"github.com/Sirupsen/logrus"
	"github.com/dokkur/swanager/config"
	"github.com/dokkur/swanager/core/db"
	"github.com/dokkur/swanager/lib"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const servicesCollectionName = "services"

// ServiceStatusStruct represents current service's task state
type ServiceStatusStruct struct {
	ReplicaID string    `json:"replica_id"`
	Node      string    `json:"node"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"error"`
}

// ServiceEnvVariable - represents name value struct
type ServiceEnvVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ServicePublishedPort - represents serivce publishing port
type ServicePublishedPort struct {
	Internal uint32 `json:"internal"`
	External uint32 `json:"external"`
	Protocol string `json:"protocol"`
	Disabled bool   `json:"disabled"`
}

// ServiceVolume - represents mounted volume
type ServiceVolume struct {
	Service string `json:"service" bson:"service"`
	Backend string `json:"backend,omitempty" bson:"backend"`
	AppWide bool   `json:"app_wide,omitempty" bson:"app_wide"`
	Size    int64  `json:"size,omitempty" bson:"-"`
}

// FrontendEndpoint - represents frontend endpoint
type FrontendEndpoint struct {
	Domain       string `json:"domain" bson:"domain"`
	InternalPort uint32 `json:"internal_port" bson:"internal_port"`
	ExternalPort uint32 `json:"external_port" bsin:"external_port"`
	Disabled     bool   `json:"disabled" bson:"disabled"`
}

// Service describes service entity
type Service struct {
	ID                string                 `bson:"_id,omitempty" json:"id"`
	Name              string                 `json:"name"`
	Image             string                 `json:"image"`
	Command           string                 `json:"command"`
	NSName            string                 `json:"ns_name" bson:"ns_name"`
	Replicas          *uint64                `json:"replicas"`
	Parallelism       uint64                 `json:"parallelism"`
	EnvVariables      []ServiceEnvVariable   `json:"env" bson:"env_vars"`
	PublishedPorts    []ServicePublishedPort `json:"published_ports" bson:"published_ports"`
	FrontendEndpoints []FrontendEndpoint     `json:"frontend_endpoints" bson:"frontend_endpoints"`
	ApplicationID     string                 `bson:"application_id,omitempty" json:"application_id,omitempty"`
	UserID            string                 `bson:"user_id" json:"-"`
	Volumes           []ServiceVolume        `bson:"volumes" json:"volumes"`
	Application       Application            `bson:"-" json:"-"`
	Status            []ServiceStatusStruct  `bson:"-" json:"status,omitempty"`
}

// GetService return service if it exists
func GetService(params map[string]interface{}) (*Service, error) {
	session := db.GetSession()
	defer session.Close()
	c := getServicesCollection(session)

	service := Service{}

	if err := c.Find(params).One(&service); err != nil {
		return nil, fmt.Errorf("GetService error: %s", err)
	}

	return &service, nil
}

// GetServices returns services by filters
func GetServices(params map[string]interface{}) ([]Service, error) {
	session := db.GetSession()
	defer session.Close()
	c := getServicesCollection(session)

	services := make([]Service, 0)
	if err := c.Find(params).All(&services); err != nil {
		return nil, err
	}
	return services, nil
}

// MustGetPublicPort returns random free public port in 10000 - 60000
func MustGetPublicPort() uint32 {
	existingPorts := getPublishedPorts(bson.M{})
	var port uint32
	for {
		port = uint32(10000 + 50000*(rand.Float64()))
		if _, ok := existingPorts[port]; ok {
			continue
		}
		break
	}
	return port
}

// PublicPortExists checks if public port exists
func PublicPortExists(port uint32, params map[string]interface{}) bool {
	existingPorts := getPublishedPorts(params)
	_, exists := existingPorts[port]
	return exists
}

type publishedPorts struct {
	Ports []ServicePublishedPort `bson:"published_ports"`
}

// Delete - removed entity from database
func (s *Service) Delete() error {
	session := db.GetSession()
	defer session.Close()

	c := getServicesCollection(session)
	if err := c.RemoveId(s.ID); err != nil {
		return err
	}
	return nil
}

// UpdateParams - updates service entity with other service entity
func (s *Service) UpdateParams(newService *Service) (errors []string) {
	s.Name = newService.Name
	s.Image = newService.Image
	s.Command = newService.Command
	s.Replicas = newService.Replicas
	s.Parallelism = newService.Parallelism

	var volumes = make([]ServiceVolume, 0)
	for _, vol := range newService.Volumes {
		if vol.Service == "" {
			errors = append(errors, "Empty service volume provided, ignoring")
			continue
		}

		if vol.AppWide && vol.Backend == "" {
			errors = append(errors, "Empty backend volume provided, ignoring")
			continue
		}

		volumes = append(volumes, ServiceVolume{
			Service: vol.Service,
			Backend: vol.Backend,
			AppWide: vol.AppWide,
		})
	}
	s.Volumes = volumes

	var variables = make([]ServiceEnvVariable, 0)
	for _, variable := range newService.EnvVariables {
		if len(variable.Name) > 0 && len(variable.Value) > 0 {
			variables = append(variables, ServiceEnvVariable{
				Name:  prepareEnvVariableName(variable.Name),
				Value: variable.Value,
			})
		}
	}
	s.EnvVariables = variables

	var ports = make([]ServicePublishedPort, 0)
	for _, port := range newService.PublishedPorts {
		if port.Internal == 0 || port.Internal > 65535 {
			errors = append(errors, fmt.Sprintf("Internal port (%d) is 0 or greather than 65535", port.Internal))
			continue
		}

		// Don't save empty or privileged internal port.
		externalPort := port.External
		if externalPort <= 1024 ||
			PublicPortExists(externalPort, bson.M{"_id": bson.M{"$ne": s.ID}}) {
			externalPort = MustGetPublicPort()
			errors = append(errors, fmt.Sprintf("External port (%d) less than 1024 or exists, changed to %d", port.External, externalPort))
		}

		ports = append(ports, ServicePublishedPort{
			Internal: port.Internal,
			External: externalPort,
			Protocol: port.Protocol,
			Disabled: port.Disabled,
		})
	}
	s.PublishedPorts = ports

	var frontends = make([]FrontendEndpoint, 0)
	for _, frontend := range newService.FrontendEndpoints {
		if frontend.InternalPort == 0 || frontend.InternalPort > 65535 ||
			frontend.ExternalPort == 0 || frontend.ExternalPort > 65535 {
			errors = append(errors, fmt.Sprintf("Internal port (%d) or External port (%d) is 0 or greather than 65535", frontend.InternalPort, frontend.ExternalPort))
			continue
		}

		frontends = append(frontends, FrontendEndpoint{
			Domain:       frontend.Domain,
			InternalPort: frontend.InternalPort,
			ExternalPort: frontend.ExternalPort,
			Disabled:     frontend.Disabled,
		})
	}
	s.FrontendEndpoints = frontends

	return
}

// Save saves user entity in db
func (s *Service) Save() error {
	if s.ID == "" {
		return s.Create()
	}

	session := db.GetSession()
	defer session.Close()
	c := getServicesCollection(session)

	if err := c.Update(bson.M{"_id": s.ID}, bson.M{"$set": s}); err != nil {
		return fmt.Errorf("Unable to save service: %s", err)
	}
	return nil
}

// Create creates user in db
func (s *Service) Create() error {
	logService().WithField("service", s).Debugf("Creating service")
	session := db.GetSession()
	defer session.Close()
	c := getServicesCollection(session)

	s.LoadApplication()

	s.ID = lib.GenerateUUID()
	s.NSName = nsName(s)

	logService().WithField("service", s).Debugf("Creating service")

	if err := c.Insert(s); err != nil {
		return fmt.Errorf("Unable to create service: %s", err)
	}
	return nil
}

// LoadApplication loads application to Application field
func (s *Service) LoadApplication() error {
	// If it's already loaded, return
	if &s.Application != nil && len(s.Application.Name) > 0 {
		return nil
	}

	session := db.GetSession()
	defer session.Close()
	c := getApplicationsCollection(session)

	application := Application{}

	if err := c.Find(bson.M{"_id": s.ApplicationID}).One(&application); err != nil {
		panic(fmt.Errorf("LoadApplication error: %s", err))
	}

	s.Application = application

	return nil
}

// AddServiceStatus -add service task status to list
func (s *Service) AddServiceStatus(status ServiceStatusStruct) {
	s.Status = append(s.Status, status)
}

func getServicesCollection(session *mgo.Session) *mgo.Collection {
	return session.DB(config.DatabaseName).C(servicesCollectionName)
}

func nsName(service *Service) string {
	service.LoadApplication()

	app := lib.IdentifierName(service.Application.Name)
	serv := lib.IdentifierName(service.Name)

	name := fmt.Sprintf("%s-%s-%s", app, serv, service.ID)
	crc32q := crc32.MakeTable(0xD5828281)

	return fmt.Sprintf("%s-%s-%08x", app, serv, crc32.Checksum([]byte(name), crc32q))
}

func prepareEnvVariableName(name string) (result string) {
	result = strings.ToUpper(name)
	result = strings.Replace(result, " ", "_", -1)
	return
}

func getPublishedPorts(params map[string]interface{}) map[uint32]bool {
	session := db.GetSession()
	defer session.Close()

	c := getServicesCollection(session)
	res := make([]publishedPorts, 0)
	if err := c.Find(params).
		Select(bson.M{"published_ports.external": 1, "_id": 0}).
		All(&res); err != nil {
		panic(err)
	}

	ports := make(map[uint32]bool, 0)
	for _, pubPort := range res {
		for _, port := range pubPort.Ports {
			ports[port.External] = true
		}
	}
	return ports
}

func logService() *logrus.Entry {
	return logrus.WithField("module", "entities.Service")
}
