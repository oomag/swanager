package mongo

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"

	"github.com/da4nik/swanager/core/entities"
)

// GetUser retrieves user from db
func (m *Mongo) GetUser(emailOrID string) (*entities.User, error) {
	session := m.getSession()
	defer session.Close()

	user := entities.User{}

	c := session.DB("swanager").C("users")
	err := c.Find(bson.M{"email": emailOrID, "_id": emailOrID}).One(&user)
	if err != nil {
		return nil, fmt.Errorf("User not found")
	}
	return &user, nil
}
