package service

import (
	"net/http"

	"github.com/da4nik/swanager/api/common"
	"github.com/da4nik/swanager/core/entities"
	"github.com/da4nik/swanager/core/swarm"
	swarm_service "github.com/da4nik/swanager/core/swarm/service"
	"github.com/gin-gonic/gin"
)

// GetRoutesForRouter adds resource routes to api router
func GetRoutesForRouter(router *gin.RouterGroup) {

	services := router.Group("/services", common.Auth(true))
	{
		services.GET("", list)
		services.POST("", create)
	}

	service := services.Group("/:service_id")
	{
		service.GET("", show)
		service.PUT("", update)
		service.DELETE("", delete)
		service.PUT("/start", start)
		service.PUT("/stop", stop)
	}
}

func list(c *gin.Context) {
	currentUser := common.MustGetCurrentUser(c)

	searchOptions := gin.H{"user_id": currentUser.ID}
	if len(c.Param("app_id")) > 0 {
		searchOptions["application_id"] = c.Param("app_id")
	}

	services, err := entities.GetServices(searchOptions)
	if err != nil {
		common.RenderError(c, http.StatusInternalServerError, err.Error())
		return
	}

	for serviceIndex := range services {
		swarm.GetServiceStatuses(&services[serviceIndex])
	}

	c.JSON(http.StatusOK, gin.H{"services": services})
}

func delete(c *gin.Context) {
	service, err := getService(c, c.Param("service_id"))
	if err != nil {
		common.RenderError(c, http.StatusBadRequest, "Service not found")
		return
	}

	swarm_service.Remove(service)
	if err := service.Delete(); err != nil {
		common.RenderError(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"service": service})
}

func update(c *gin.Context) {
	var newService entities.Service
	if err := c.BindJSON(&newService); err != nil {
		common.RenderError(c, http.StatusBadRequest, err)
		return
	}

	service, err := getService(c, c.Param("service_id"))
	if err != nil {
		common.RenderError(c, http.StatusNotFound, "Service not found")
		return
	}

	service.UpdateParams(&newService)
	service.Save()

	swarm.UpdateService(service)

	serviceStatus, err := swarm_service.Status(service)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"service": service, "status_error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"service": service, "status": serviceStatus})
}

func create(c *gin.Context) {
	currentUser := common.MustGetCurrentUser(c)
	var service entities.Service
	if err := c.BindJSON(&service); err != nil {
		common.RenderError(c, http.StatusBadRequest, err)
		return
	}

	if len(service.Name) == 0 {
		common.RenderError(c, http.StatusBadRequest, gin.H{"name": "Name is empty"})
		return
	}

	service.UserID = currentUser.ID

	if len(service.ApplicationID) == 0 {
		if len(c.Param("app_id")) > 0 {
			service.ApplicationID = c.Param("app_id")
		} else {
			common.RenderError(c, http.StatusBadRequest, gin.H{"app_id": "Application ID (app_id) is empty"})
			return
		}
	}

	if err := service.Save(); err != nil {
		common.RenderError(c, http.StatusUnprocessableEntity, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"service": service})
}

func show(c *gin.Context) {
	currentUser := common.MustGetCurrentUser(c)

	service, err := entities.GetService(gin.H{
		"user_id": currentUser.ID,
		"_id":     c.Param("service_id"),
	})

	if err != nil {
		common.RenderError(c, http.StatusNotFound, "Service not found")
		return
	}

	serviceStatus, err := swarm_service.Status(service)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"service": service, "status_error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"service": service, "status": serviceStatus})
}

func start(c *gin.Context) {
	currentUser := common.MustGetCurrentUser(c)
	service, err := getService(c, c.Param("service_id"))
	if err != nil {
		common.RenderError(c, http.StatusBadRequest, "Service not found: "+err.Error())
		return
	}

	common.RunAsync(common.AsyncJobContext{
		User:       currentUser,
		GinContext: c,
		Process: func() (interface{}, error) {
			if err := swarm.StartService(service); err != nil {
				return "", err
			}
			return service, nil
		},
	})
}

func stop(c *gin.Context) {
	currentUser := common.MustGetCurrentUser(c)
	service, err := getService(c, c.Param("service_id"))
	if err != nil {
		common.RenderError(c, http.StatusBadRequest, "Service not found: "+err.Error())
		return
	}

	common.RunAsync(common.AsyncJobContext{
		User:       currentUser,
		GinContext: c,
		Process: func() (interface{}, error) {
			if err := swarm.StopService(service); err != nil {
				common.RenderError(c, http.StatusBadRequest, "Error stoping service: "+err.Error())
				return "", err
			}
			return service, nil
		},
	})
}

// getService returns service by it's id and current user id
func getService(c *gin.Context, serviceID string) (app *entities.Service, err error) {
	currentUser := common.MustGetCurrentUser(c)
	app, err = entities.GetService(gin.H{"_id": serviceID, "user_id": currentUser.ID})
	return
}

func loadServiceStatus(service *entities.Service) {
	states, err := swarm_service.Status(service)

	if err != nil {
		service.AddServiceStatus(entities.ServiceStatusStruct{Status: "Not exists"})
		return
	}

	for _, state := range states {
		service.AddServiceStatus(entities.ServiceStatusStruct{
			Node:      state.Node,
			Status:    state.Status,
			Timestamp: state.Timestamp,
		})
	}
}
