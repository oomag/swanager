package db

import (
	"github.com/da4nik/swanager/config"
	mgo "gopkg.in/mgo.v2"
)

// GetSession return opened mongo session
func GetSession() *mgo.Session {
	session, err := mgo.Dial(config.MongoURL)
	if err != nil {
		panic(err)
	}
	return session
}
