package mongo

import (
	"github.com/da4nik/swanager/config"
	mgo "gopkg.in/mgo.v2"
)

// Mongo driver struct
type Mongo struct {
}

func (m *Mongo) getSession() *mgo.Session {
	session, err := mgo.Dial(config.MongoURL)
	if err != nil {
		panic(err)
	}
	return session
}
