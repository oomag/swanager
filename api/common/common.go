package common

import (
	"fmt"
	"net/http"

	"github.com/da4nik/swanager/core/auth"
	"github.com/da4nik/swanager/core/entities"
	"github.com/gin-gonic/gin"
)

// GetCurrentUser return current user from context
func GetCurrentUser(c *gin.Context) (*entities.User, error) {
	interfaceUser, exists := c.Get("CurrentUser")
	if !exists {
		return nil, fmt.Errorf("Current user not found")
	}

	currentUser := interfaceUser.(*entities.User)

	return currentUser, nil
}

// MustGetCurrentUser return current user from context or panic
func MustGetCurrentUser(c *gin.Context) *entities.User {
	user, err := GetCurrentUser(c)
	if err != nil {
		panic(err)
	}
	return user
}

// Auth authentication handler function
func Auth(authenticate bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !authenticate {
			c.Next()
			return
		}

		token := c.Request.Header.Get("Authorization")
		user, err := auth.WithToken(token)
		if err != nil {
			RenderError(c, http.StatusUnauthorized, "Unauthorized")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("CurrentUser", user)
		c.Next()
	}
}

// RenderError formats JSON error
func RenderError(c *gin.Context, status int, errors interface{}) {
	c.JSON(status, gin.H{"errors": errors})
}
