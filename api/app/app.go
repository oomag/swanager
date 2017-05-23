package app

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
	"github.com/gin-gonic/gin"
)

// GetRoutesForRouter adds resource routes to api router
func GetRoutesForRouter(router *gin.RouterGroup) {

	apps := router.Group("/apps", common.Auth(true))
	{
		apps.GET("", list)
		apps.POST("", create)
	}

	app := apps.Group("/:app_id")
	{
		app.GET("", show)
		app.PUT("", update)
		app.DELETE("", destroy)
		app.PUT("/start", start)
		app.PUT("/stop", stop)
	}
}

func list(c *gin.Context) {
	currentUser := common.MustGetCurrentUser(c)

	services, err := entities.GetServices(gin.H{"user_id": currentUser.ID})
	if err != nil {
		common.RenderError(c, http.StatusInternalServerError, err.Error())
		return
	}

	appServices := make(map[string][]string)
	for _, service := range services {
		if _, exists := appServices[service.ApplicationID]; !exists {
			appServices[service.ApplicationID] = make([]string, 0)
		}
		appServices[service.ApplicationID] = append(appServices[service.ApplicationID], service.ID)
	}

	applications, err := currentUser.GetApplications()
	if err != nil {
		common.RenderError(c, http.StatusUnprocessableEntity, err)
		return
	}

	for index := range applications {
		applications[index].ServiceIDS = appServices[applications[index].ID]
	}

	c.JSON(http.StatusOK, gin.H{"applications": applications})
}

func show(c *gin.Context) {
	app, err := getApplication(c, c.Param("app_id"))
	if err != nil {
		common.RenderError(c, http.StatusNotFound, err)
		return
	}

	app.LoadServices()
	for _, service := range app.Services {
		serviceStatus, err := swarm_service.Status(&service)
		if err == nil {
			for _, status := range serviceStatus {
				service.Status = append(service.Status, entities.ServiceStatusStruct{
					Node:      status.Node,
					Status:    status.Status,
					Timestamp: status.Timestamp,
				})
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"application": app})
}

func create(c *gin.Context) {
	currentUser := common.MustGetCurrentUser(c)
	var app entities.Application
	if err := c.BindJSON(&app); err != nil {
		common.RenderError(c, http.StatusBadRequest, err)
		return
	}

	if len(app.Name) == 0 {
		common.RenderError(c, http.StatusBadRequest, gin.H{"name": "Name is empty"})
		return
	}

	app.UserID = currentUser.ID
	if err := app.Save(); err != nil {
		common.RenderError(c, http.StatusUnprocessableEntity, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"application": app})
}

func update(c *gin.Context) {
	var newApp entities.Application
	if err := c.BindJSON(&newApp); err != nil {
		common.RenderError(c, http.StatusBadRequest, err)
		return
	}

	app, err := getApplication(c, c.Param("app_id"))
	if err != nil {
		common.RenderError(c, http.StatusBadRequest, "Application not found")
		return
	}

	app.Name = newApp.Name
	app.Save()
	app.LoadServices()

	c.JSON(http.StatusOK, gin.H{"application": app})
}

func start(c *gin.Context) {
	currentUser := common.MustGetCurrentUser(c)
	app, err := getApplication(c, c.Param("app_id"))
	if err != nil {
		common.RenderError(c, http.StatusBadRequest, "Application not found: "+err.Error())
		return
	}

	cmd, respChan, errChan := command.NewAppStartCommand(command.AppStart{
		User:        currentUser,
		Application: app,
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
	app, err := getApplication(c, c.Param("app_id"))
	if err != nil {
		common.RenderError(c, http.StatusBadRequest, "Application not found: "+err.Error())
		return
	}

	cmd, respChan, errChan := command.NewAppStopCommand(command.AppStop{
		User:        currentUser,
		Application: app,
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

func destroy(c *gin.Context) {
	app, err := getApplication(c, c.Param("app_id"))
	if err != nil {
		common.RenderError(c, http.StatusBadRequest, "Application not found: "+err.Error())
		return
	}

	if err := swarm.StopApplication(app); err != nil {
		common.RenderError(c, http.StatusBadRequest, "Error stoping application: "+err.Error())
		return
	}

	if err := app.Delete(); err != nil {
		common.RenderError(c, http.StatusBadRequest, "Error deleting application: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"application": app})
}

func getApplication(c *gin.Context, appID string) (app *entities.Application, err error) {
	currentUser := common.MustGetCurrentUser(c)
	app, err = entities.GetApplication(gin.H{"_id": appID, "user_id": currentUser.ID})
	return
}
