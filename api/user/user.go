package user

import (
	"net/http"

	"github.com/da4nik/swanager/core/entities"
	"github.com/gin-gonic/gin"
)

// GetRoutesForRouter adds resource routes to api router
func GetRoutesForRouter(router *gin.RouterGroup) *gin.RouterGroup {

	apps := router.Group("/users")
	{
		apps.POST("", create)
	}

	app := apps.Group("/:user_id")
	{
		app.GET("", show)
	}

	return apps
}

func show(c *gin.Context) {
	user, err := entities.GetUser(c.Param("user_id"))
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func create(c *gin.Context) {
	c.AbortWithStatus(http.StatusOK)
}
