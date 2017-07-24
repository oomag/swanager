package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dokkur/swanager/api/common"
	"github.com/dokkur/swanager/command"
	"github.com/dokkur/swanager/config"
	"github.com/dokkur/swanager/core/entities"
	"github.com/dokkur/swanager/core/swarm"
	swarm_service "github.com/dokkur/swanager/core/swarm/service"
	"github.com/dokkur/swanager/frontend"
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
		service.GET("/logs", logs)
		service.DELETE("", delete)
		service.PUT("", update)
		service.PUT("/start", start)
		service.PUT("/stop", stop)

	}
}

func list(c *gin.Context) {
	currentUser := common.MustGetCurrentUser(c)

	withStatuses := c.Query("with_statuses") == "1"
	withVolumeSizes := c.Query("with_volume_sizes") == "1"

	cmd, respChan, errChan := command.NewServiceListCommand(command.ServiceList{
		User:            currentUser,
		ApplicationID:   c.Query("app_id"),
		WithStatuses:    withStatuses,
		WithVolumeSizes: withVolumeSizes,
	})
	command.RunAsync(cmd)

	select {
	case services := <-respChan:
		c.JSON(http.StatusOK, gin.H{"services": services})
	case err := <-errChan:
		common.RenderError(c, http.StatusInternalServerError, err.Error())
	case <-time.After(time.Second * time.Duration(config.RequestTimeout)):
		common.RenderError(c, http.StatusRequestTimeout, "Timeout")
	}
}

func delete(c *gin.Context) {
	cmd, respChan, errChan := command.NewServiceDeleteCommand(command.ServiceDelete{
		User:      common.MustGetCurrentUser(c),
		ServiceID: c.Param("service_id"),
	})
	command.RunAsync(cmd)

	select {
	case service := <-respChan:
		c.JSON(http.StatusOK, gin.H{"service": service})
	case err := <-errChan:
		common.RenderError(c, http.StatusUnprocessableEntity, err.Error())
	case <-time.After(time.Second * time.Duration(config.RequestTimeout)):
		common.RenderError(c, http.StatusRequestTimeout, "Timeout")
	}
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

	notes := service.UpdateParams(&newService)
	service.Save()

	swarm.UpdateService(service)
	go frontend.Update()

	serviceStatus, err := swarm_service.Status(service)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"service": service, "status_error": err, "notes": notes})
		return
	}

	c.JSON(http.StatusOK, gin.H{"service": service, "status": serviceStatus, "notes": notes})
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

	service.UpdateParams(&service)

	if err := service.Save(); err != nil {
		common.RenderError(c, http.StatusUnprocessableEntity, err)
		return
	}

	go frontend.Update()
	c.JSON(http.StatusCreated, gin.H{"service": service})
}

func show(c *gin.Context) {
	cmd, respChan, errChan := command.NewServiceInspectCommand(command.ServiceInspect{
		User:      common.MustGetCurrentUser(c),
		ServiceID: c.Param("service_id"),
	})
	command.RunAsync(cmd)

	select {
	case service := <-respChan:
		c.JSON(http.StatusOK, gin.H{"service": service})
	case err := <-errChan:
		common.RenderError(c, http.StatusUnprocessableEntity, err.Error())
	case <-time.After(time.Second * time.Duration(config.RequestTimeout)):
		common.RenderError(c, http.StatusRequestTimeout, "Timeout")
	}
}

func start(c *gin.Context) {
	currentUser := common.MustGetCurrentUser(c)
	service, err := getService(c, c.Param("service_id"))
	if err != nil {
		common.RenderError(c, http.StatusBadRequest, "Service not found: "+err.Error())
		return
	}

	cmd, respChan, errChan := command.NewServiceStartCommand(command.ServiceStart{
		User:    currentUser,
		Service: *service,
	})
	command.RunAsync(cmd)

	select {
	case job := <-respChan:
		c.JSON(http.StatusAccepted, gin.H{
			"job_id": job.ID,
			"url":    fmt.Sprintf("http://%s/api/v1/jobs/%s", c.Request.Host, job.ID),
		})
	case err = <-errChan:
		common.RenderError(c, http.StatusInternalServerError, err.Error())
	case <-time.After(time.Second * time.Duration(config.RequestTimeout)):
		common.RenderError(c, http.StatusRequestTimeout, "Timeout")
	}
}

func stop(c *gin.Context) {
	currentUser := common.MustGetCurrentUser(c)
	service, err := getService(c, c.Param("service_id"))
	if err != nil {
		common.RenderError(c, http.StatusNotFound, "Service not found: "+err.Error())
		return
	}

	cmd, respChan, errChan := command.NewServiceStopCommand(command.ServiceStop{
		User:    currentUser,
		Service: *service,
	})
	command.RunAsync(cmd)

	select {
	case job := <-respChan:
		c.JSON(http.StatusAccepted, gin.H{
			"job_id": job.ID,
			"url":    fmt.Sprintf("http://%s/api/v1/jobs/%s", c.Request.Host, job.ID),
		})
	case err = <-errChan:
		common.RenderError(c, http.StatusInternalServerError, err.Error())
	case <-time.After(time.Second * time.Duration(config.RequestTimeout)):
		common.RenderError(c, http.StatusRequestTimeout, "Timeout")
	}
}

func logs(c *gin.Context) {
	currentUser := common.MustGetCurrentUser(c)

	cmd, respChan, errChan := command.NewServiceLogsCommand(command.ServiceLogs{
		User:      currentUser,
		ServiceID: c.Param("service_id"),
	})
	command.RunAsync(cmd)

	select {
	case logs := <-respChan:
		c.JSON(http.StatusOK, gin.H{"logs": logs})
	case err := <-errChan:
		common.RenderError(c, http.StatusUnprocessableEntity, err.Error())
	case <-time.After(time.Second * time.Duration(config.RequestTimeout)):
		common.RenderError(c, http.StatusRequestTimeout, "Timeout")
	}
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
