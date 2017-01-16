package auth

import (
	"fmt"

	"github.com/da4nik/swanager/core/entities"
)

// WithToken authenticates with token
func WithToken(token string) (*entities.User, error) {
	if token == "" {
		return nil, fmt.Errorf("Empty token")
	}

	user, err := entities.GetUserByToken(token)
	if err != nil {
		return nil, fmt.Errorf("AuthWithToken error: %s", err)
	}

	return user, nil
}
