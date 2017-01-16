package entities

import (
	"fmt"

	"github.com/da4nik/swanager/config"
	"github.com/da4nik/swanager/core/db"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const applicationsCollectionName = "applications"

// Application describes application entity
type Application struct {
	ID     string `bson:"_id,omitempty"`
	Name   string
	UserID string `json:"user_id,omitempty"`
}

// GetApplication return application if it exists
func GetApplication(id string) (*Application, error) {
	session := db.GetSession()
	defer session.Close()
	c := getApplicationsCollection(session)

	application := Application{}

	if err := c.Find(bson.M{"_id": id}).One(&application); err != nil {
		return nil, fmt.Errorf("GetApplication error: %s", err)
	}

	return &application, nil
}

// Save saves user entity in db
func (a *Application) Save() error {
	if a.ID == "" {
		return a.Create()
	}

	session := db.GetSession()
	defer session.Close()
	c := getApplicationsCollection(session)

	if err := c.Update(bson.M{"_id": a.ID}, bson.M{"$set": a}); err != nil {
		return fmt.Errorf("Unable to save application: %s", err)
	}
	return nil
}

// Create creates user in db
func (a *Application) Create() error {
	session := db.GetSession()
	defer session.Close()
	c := getApplicationsCollection(session)

	a.ID = generateUUID()

	if err := c.Insert(a); err != nil {
		return fmt.Errorf("Unable to create application: %s", err)
	}
	return nil
}

// GetServices returns application services
func (a *Application) GetServices() ([]Service, error) {
	session := db.GetSession()
	defer session.Close()
	c := getServicesCollection(session)

	services := make([]Service, 0)
	err := c.Find(bson.M{"application_id": a.ID}).All(&services)
	if err != nil {
		return nil, err
	}

	return services, nil
}

func getApplicationsCollection(session *mgo.Session) *mgo.Collection {
	return session.DB(config.DatabaseName).C(applicationsCollectionName)
}
