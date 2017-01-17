package entities

import (
	"time"

	"github.com/da4nik/swanager/lib"
)

// Token represent auth token for User
type Token struct {
	Token    string    `json:"token"`
	LastUsed time.Time `json:"last_used,omitempty"`
}

// GenerateToken generated new token
func GenerateToken() *Token {
	token := Token{Token: lib.GenerateUUID(), LastUsed: time.Now()}

	return &token
}
