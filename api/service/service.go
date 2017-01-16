package service

import (
	"net/http"

	"github.com/da4nik/swanager/api/common"
	"github.com/gin-gonic/gin"
)

// GetRoutesForRouter adds resource routes to api router
func GetRoutesForRouter(router *gin.RouterGroup) *gin.RouterGroup {

	services := router.Group("/apps/:app_id/services", common.Auth(true))
	{
		services.GET("", list)
		services.POST("", create)
	}

	service := services.Group("/:service_id")
	{
		service.GET("", show)
	}

	return services
}

func list(c *gin.Context) {
	c.AbortWithStatus(http.StatusOK)
}

func create(c *gin.Context) {
	// networkID := swarm.CreateNetwork("overlay", "test_network")
	c.JSON(http.StatusCreated, struct{ ID string }{ID: "test_app_id"})
	c.AbortWithStatus(http.StatusOK)
}

func show(c *gin.Context) {
	c.AbortWithStatus(http.StatusOK)
}
