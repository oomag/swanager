package db

import (
	"github.com/da4nik/swanager/config"
	"github.com/da4nik/swanager/core/db/drivers"
	"github.com/da4nik/swanager/core/entities"
)

var db drivers.DBDriver

// DB returns db driver
func DB() *drivers.DBDriver {
	return &db
}

func init() {
	var err error
	db, err = drivers.GetDBDriver(config.DatabaseDriver)
	if err != nil {
		panic(err)
	}
}

// GetUser get user from db
func GetUser(emailOrID string) (*entities.User, error) {
	user, err := db.GetUser(emailOrID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
