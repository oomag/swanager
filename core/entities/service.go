package entities

import (
	"fmt"
	"time"

	"github.com/da4nik/swanager/config"
	"github.com/da4nik/swanager/core/db"
	"github.com/da4nik/swanager/lib"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const servicesCollectionName = "services"

// ServiceStatusStruct represents current service's task state
type ServiceStatusStruct struct {
	Node      string
	Status    string
	Timestamp time.Time
}

// Service describes service entity
type Service struct {
	ID            string                `bson:"_id,omitempty" json:"id"`
	Name          string                `json:"name"`
	Image         string                `json:"image"`
	Replicas      *uint64               `json:"replicas"`
	Parallelism   uint64                `json:"parallelism"`
	Status        []ServiceStatusStruct `bson:"-" json:"status,omitempty"`
	ApplicationID string                `bson:"application_id,omitempty" json:"application_id,omitempty"`
	Application   Application           `bson:"-" json:"-"`
	UserID        string                `bson:"user_id" json:"user_id"`
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
	session := db.GetSession()
	defer session.Close()
	c := getServicesCollection(session)

	s.ID = lib.GenerateUUID()

	if err := c.Insert(s); err != nil {
		return fmt.Errorf("Unable to create service: %s", err)
	}
	return nil
}

func getServicesCollection(session *mgo.Session) *mgo.Collection {
	return session.DB(config.DatabaseName).C(servicesCollectionName)
}
