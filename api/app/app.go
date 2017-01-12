package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetRoutesForRouter adds resource routes to api router
func GetRoutesForRouter(router *gin.RouterGroup) *gin.RouterGroup {

	apps := router.Group("/apps")
	{
		apps.GET("", list)
		apps.POST("", create)
	}

	app := apps.Group("/:app_id")
	{
		app.GET("", show)
	}

	return apps
}

func list(c *gin.Context) {
	c.AbortWithStatus(http.StatusOK)
}

func show(c *gin.Context) {
	c.AbortWithStatus(http.StatusOK)
}

func create(c *gin.Context) {
	c.AbortWithStatus(http.StatusOK)
}
