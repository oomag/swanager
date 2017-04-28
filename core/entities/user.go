package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/da4nik/swanager/config"
	"github.com/da4nik/swanager/core/db"
	"github.com/da4nik/swanager/lib"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const usersCollectionName = "users"

// User describes service entity
type User struct {
	ID       string  `json:"id" bson:"_id,omitempty"`
	Email    string  `json:"email"`
	Password string  `json:"-"`
	Tokens   []Token `json:"-"`
}

// GetUser returns User object from db
func GetUser(email string) (*User, error) {
	session := db.GetSession()
	defer session.Close()
	c := getUsersCollection(session)

	user := User{}

	if err := c.Find(bson.M{"email": email}).One(&user); err != nil {
		return nil, fmt.Errorf("GetUser error: %s", err)
	}

	return &user, nil
}

// GetUsers returns all users
func GetUsers(params bson.M) ([]User, error) {
	session := db.GetSession()
	defer session.Close()
	c := getUsersCollection(session)

	users := make([]User, 0)
	if err := c.Find(params).All(&users); err != nil {
		return nil, fmt.Errorf("GetUsers error: %s", err)
	}
	return users, nil
}

// GetUserByToken returns user by associated token
func GetUserByToken(token string) (*User, error) {
	session := db.GetSession()
	defer session.Close()
	c := getUsersCollection(session)

	user := User{}

	// Token expires in 1 day
	expireDate := time.Now().Add(time.Duration(-24) * time.Hour)

	findParams := bson.M{"tokens": bson.M{"$elemMatch": bson.M{"token": token, "lastused": bson.M{"$gt": expireDate}}}}

	err := c.Find(findParams).One(&user)
	if err != nil {
		return nil, fmt.Errorf("GetUser error: %s", err)
	}

	return &user, nil
}

// Save saves user entity in db
func (u *User) Save() error {
	if u.ID == "" {
		return u.Create()
	}

	session := db.GetSession()
	defer session.Close()
	c := getUsersCollection(session)

	if err := c.Update(bson.M{"email": u.Email}, bson.M{"$set": u}); err != nil {
		return fmt.Errorf("Unable to save user: %s", err)
	}
	return nil
}

// Create creates user in db
func (u *User) Create() error {
	session := db.GetSession()
	defer session.Close()
	c := getUsersCollection(session)

	u.ID = lib.GenerateUUID()
	u.Password = lib.CalculateMD5(strings.TrimSpace(u.Password))

	if err := c.Insert(u); err != nil {
		return fmt.Errorf("Unable to create user: %s", err)
	}
	return nil
}

// GetApplications return user's applications
func (u *User) GetApplications() ([]Application, error) {
	session := db.GetSession()
	defer session.Close()
	c := getApplicationsCollection(session)

	applications := make([]Application, 0)
	err := c.Find(bson.M{"user_id": u.ID}).All(&applications)
	if err != nil {
		return nil, err
	}

	return applications, nil
}

func getUsersCollection(session *mgo.Session) *mgo.Collection {
	return session.DB(config.DatabaseName).C(usersCollectionName)
}
