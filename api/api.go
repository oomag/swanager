package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/da4nik/swanager/api/app"
	"github.com/da4nik/swanager/api/job"
	"github.com/da4nik/swanager/api/service"
	"github.com/da4nik/swanager/api/session"
	"github.com/da4nik/swanager/api/user"
	"github.com/da4nik/swanager/api/ws"
	"github.com/da4nik/swanager/config"
	"github.com/gin-gonic/gin"
)

// Start starts listening incoming connections
func Start() {
	router := gin.New()

	if gin.Mode() == gin.ReleaseMode {
		router.Use(ginlogrus(logrus.StandardLogger(), time.RFC3339, true))
	} else {
		router.Use(gin.Logger())
	}

	router.Use(corsMiddleware())
	router.Use(gin.Recovery())

	ws.InitWS(router)

	apiGroup := router.Group("/api/v1")
	app.GetRoutesForRouter(apiGroup)
	service.GetRoutesForRouter(apiGroup)
	user.GetRoutesForRouter(apiGroup)
	session.GetRoutesForRouter(apiGroup)
	job.GetRoutesForRouter(apiGroup)

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

func ginlogrus(logger *logrus.Logger, timeFormat string, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}

		ip := c.ClientIP()

		entry := logger.WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"ip":         ip,
			"latency":    latency,
			"user-agent": c.Request.UserAgent(),
			"time":       end.Format(timeFormat),
		})

		msg := fmt.Sprintf("%s -> %s %s (%s)", ip, c.Request.Method, path, latency)

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.String())
		} else {
			entry.Info(msg)
		}
	}
}
