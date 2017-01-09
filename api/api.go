package api

import (
	"net/http"

	"github.com/da4nik/books/config"
	"github.com/da4nik/swanager/api/service"
	"github.com/da4nik/swanager/core/auth"
	"github.com/gin-gonic/gin"
)

func init() {
	router := gin.Default()
	router.Use(corsMiddleware())
	router.Use(tokenAuthMiddleware())

	apiGroup := router.Group("/api")
	service.GetRoutesForRouter(apiGroup)

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

func tokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if auth.WithToken(token) {
			c.Next()
			return
		}
		c.AbortWithStatus(http.StatusUnauthorized)
		c.Next()
	}
}
