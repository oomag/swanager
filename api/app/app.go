package app

import (
	"net/http"

	"github.com/da4nik/swanager/api/common"
	"github.com/gin-gonic/gin"
)

// GetRoutesForRouter adds resource routes to api router
func GetRoutesForRouter(router *gin.RouterGroup) *gin.RouterGroup {

	apps := router.Group("/apps", common.Auth(true))
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
	currentUser := common.MustGetCurrentUser(c)

	applications, err := currentUser.GetApplications()
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"applications": applications})
}

func show(c *gin.Context) {
	c.AbortWithStatus(http.StatusOK)
}

func create(c *gin.Context) {
	c.AbortWithStatus(http.StatusOK)
}
