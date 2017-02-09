package job

import (
	"net/http"

	"github.com/da4nik/swanager/api/common"
	"github.com/da4nik/swanager/core/entities"
	"github.com/gin-gonic/gin"
)

// GetRoutesForRouter adds resource routes to api router
func GetRoutesForRouter(router *gin.RouterGroup) {
	apps := router.Group("/jobs/:job_id", common.Auth(true))
	{
		apps.GET("", show)
	}
}

func show(c *gin.Context) {
	job, err := entities.GetJob(c.Param("job_id"))
	if err != nil {
		common.RenderError(c, http.StatusNotFound, "not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"job": job})
}
