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
	ID              string `bson:"_id,omitempty"`
	Name            string
	Image           string
	Replicas        *uint64
	Parallelism     uint64
	Status          []ServiceStatusStruct `bson:"-" json:"status,omitempty"`
	ApplicationID   string                `bson:"application_id,omitempty" json:"application_id"`
	DockerServiceID string                `bson:"docker_service_id,omitempty" json:"-"`
	Application     Application           `bson:"-" json:"-"`
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
