package entities

import (
	"fmt"
	"time"

	"hash/crc32"

	"github.com/Sirupsen/logrus"
	"github.com/da4nik/swanager/config"
	"github.com/da4nik/swanager/core/db"
	"github.com/da4nik/swanager/lib"
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
}

// Service describes service entity
type Service struct {
	ID            string                `bson:"_id,omitempty" json:"id"`
	Name          string                `json:"name"`
	Image         string                `json:"image"`
	NSName        string                `json:"ns_name" bson:"ns_name"`
	Replicas      *uint64               `json:"replicas"`
	Parallelism   uint64                `json:"parallelism"`
	Status        []ServiceStatusStruct `bson:"-" json:"status,omitempty"`
	ApplicationID string                `bson:"application_id,omitempty" json:"application_id,omitempty"`
	Application   Application           `bson:"-" json:"-"`
	UserID        string                `bson:"user_id" json:"-"`
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
func (s *Service) UpdateParams(newService *Service) error {
	s.Name = newService.Name
	s.Image = newService.Image
	s.Replicas = newService.Replicas
	s.Parallelism = newService.Parallelism
	return nil
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

// AddServiceStatus -sdfkjsdf
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

func logService() *logrus.Entry {
	return logrus.WithField("module", "entities.Service")
}
