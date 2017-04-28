package auth

import (
	"fmt"

	"github.com/da4nik/swanager/core/entities"
	"github.com/da4nik/swanager/lib"
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

// WithEmailAndPassword auths user with email and password, return newly created token
func WithEmailAndPassword(email, password string) (*entities.Token, error) {
	user, err := entities.GetUser(email)
	if err != nil {
		return nil, authError()
	}

	if user.Password != lib.CalculateMD5(password) {
		return nil, authError()
	}

	token := entities.GenerateToken()
	token.User = user

	user.Tokens = append(user.Tokens, *token)
	user.Save()

	return token, nil
}

// Deauthorize logs user out
func Deauthorize(user *entities.User) error {
	user.Tokens = make([]entities.Token, 0)
	if err := user.Save(); err != nil {
		return err
	}

	return nil
}

func authError() error {
	return fmt.Errorf("Email or Password are wrong")
}
