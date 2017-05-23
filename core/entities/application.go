package entities

import (
	"fmt"

	"github.com/dokkur/swanager/config"
	"github.com/dokkur/swanager/core/db"
	"github.com/dokkur/swanager/lib"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const applicationsCollectionName = "applications"

// Application describes application entity
type Application struct {
	ID         string    `json:"id" bson:"_id,omitempty"`
	Name       string    `json:"name"`
	UserID     string    `json:"-" bson:"user_id"`
	Services   []Service `json:"services,omitempty" bson:"-"`
	ServiceIDS []string  `json:"service_ids,omitempty" bson:"-"`
}

// GetApplications returns applications by request
func GetApplications(params map[string]interface{}) ([]Application, error) {
	session := db.GetSession()
	defer session.Close()
	c := getApplicationsCollection(session)

	applications := make([]Application, 0)
	if err := c.Find(params).All(&applications); err != nil {
		return nil, err
	}
	return applications, nil
}

// GetApplication return application if it exists
func GetApplication(params map[string]interface{}) (*Application, error) {
	session := db.GetSession()
	defer session.Close()
	c := getApplicationsCollection(session)

	application := Application{}

	if err := c.Find(params).One(&application); err != nil {
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

	a.ID = lib.GenerateUUID()

	if err := c.Insert(a); err != nil {
		return fmt.Errorf("Unable to create application: %s", err)
	}
	return nil
}

// Delete - deletes applications and it's services
func (a *Application) Delete() error {
	session := db.GetSession()
	defer session.Close()
	c := getApplicationsCollection(session)

	a.LoadServices()

	for _, service := range a.Services {
		if err := service.Delete(); err != nil {
			fmt.Printf("Problems deleting service: %s", err.Error())
		}
	}

	if err := c.RemoveId(a.ID); err != nil {
		return fmt.Errorf("Error deleting application: %s", err.Error())
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

// LoadServices retrieves services associated with application
func (a *Application) LoadServices() {
	if len(a.Services) > 0 {
		return
	}

	services, err := a.GetServices()
	if err != nil {
		return
	}

	a.Services = services
}

func getApplicationsCollection(session *mgo.Session) *mgo.Collection {
	return session.DB(config.DatabaseName).C(applicationsCollectionName)
}
