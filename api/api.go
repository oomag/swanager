package api

import (
	"net/http"

	"github.com/da4nik/swanager/api/app"
	"github.com/da4nik/swanager/api/service"
	"github.com/da4nik/swanager/api/session"
	"github.com/da4nik/swanager/api/user"
	"github.com/da4nik/swanager/config"
	"github.com/gin-gonic/gin"
)

func init() {
	router := gin.Default()

	router.Use(corsMiddleware())

	apiGroup := router.Group("/api")
	app.GetRoutesForRouter(apiGroup)
	service.GetRoutesForRouter(apiGroup)
	user.GetRoutesForRouter(apiGroup)
	session.GetRoutesForRouter(apiGroup)

	router.Run(":" + config.Port)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}
