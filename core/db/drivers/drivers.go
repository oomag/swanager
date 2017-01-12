package drivers

import (
	"fmt"
	"strings"

	"github.com/da4nik/swanager/core/db/drivers/mongo"
	"github.com/da4nik/swanager/core/entities"
)

// DBDriver interface for backend db drivers
type DBDriver interface {
	GetUser(emailOrID string) (*entities.User, error)
}

// GetDBDriver return selected db driver
func GetDBDriver(driver string) (DBDriver, error) {
	switch strings.ToLower(driver) {
	case "mongo":
		return new(mongo.Mongo), nil
	default:
		return nil, fmt.Errorf("Unknown driver %s.", driver)
	}
}
