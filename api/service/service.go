package service

import (
	"net/http"

	"github.com/da4nik/swanager/core/swarm"
	"github.com/gin-gonic/gin"
)

// GetRoutesForRouter adds resource routes to api router
func GetRoutesForRouter(router *gin.RouterGroup) *gin.RouterGroup {

	services := router.Group("/services")
	{
		services.GET("", getList)
		services.POST("", createService)
	}

	service := services.Group("/:id")
	{
		service.GET("", showService)
	}

	return services
}

func getList(c *gin.Context) {
	c.AbortWithStatus(http.StatusOK)
}

func createService(c *gin.Context) {
	networkID := swarm.CreateNetwork("overlay", "test_network")
	c.JSON(http.StatusCreated, struct{ networkID string }{networkID: networkID})
	// c.AbortWithStatus(http.StatusOK)
}

func showService(c *gin.Context) {
	c.AbortWithStatus(http.StatusOK)
}
