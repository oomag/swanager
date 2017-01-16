package entities

import "time"

// Token represent auth token for User
type Token struct {
	Token    string
	LastUsed time.Time `json:"last_used,omitempty"`
}
