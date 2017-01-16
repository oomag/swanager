package session

import (
	"net/http"

	"github.com/da4nik/swanager/api/common"
	"github.com/gin-gonic/gin"
)

// GetRoutesForRouter adds resource routes to api router
func GetRoutesForRouter(router *gin.RouterGroup) *gin.RouterGroup {

	auth := router.Group("/session")
	{
		auth.POST("", common.Auth(false), login)
		auth.DELETE("", common.Auth(true), logout)
	}

	return auth
}

func login(c *gin.Context) {
	c.AbortWithStatus(http.StatusOK)
}

func logout(c *gin.Context) {
	c.AbortWithStatus(http.StatusOK)
}
