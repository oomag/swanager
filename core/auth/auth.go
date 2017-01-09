package auth

// WithToken authenticates with token
func WithToken(token string) bool {
	if token == "" {
		return false
	}
	return true
}
