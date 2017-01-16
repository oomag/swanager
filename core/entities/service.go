package entities

import (
	"fmt"

	"github.com/da4nik/swanager/config"
	"github.com/da4nik/swanager/core/db"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const servicesCollectionName = "services"

// Service describes service entity
type Service struct {
	ID            string `bson:"_id,omitempty"`
	Name          string
	Image         string
	Replicas      int
	ApplicationID string `bson:"application_id,omitempty"`
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

	s.ID = generateUUID()

	if err := c.Insert(s); err != nil {
		return fmt.Errorf("Unable to create service: %s", err)
	}
	return nil
}

func getServicesCollection(session *mgo.Session) *mgo.Collection {
	return session.DB(config.DatabaseName).C(servicesCollectionName)
}
