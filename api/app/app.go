package app

import (
	"net/http"

	"github.com/da4nik/swanager/api/common"
	"github.com/da4nik/swanager/core/entities"
	swarm_service "github.com/da4nik/swanager/core/swarm/service"
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
		app.POST("/start", start)
		app.POST("/stop", stop)
	}
}

func list(c *gin.Context) {
	currentUser := common.MustGetCurrentUser(c)

	applications, err := currentUser.GetApplications()
	if err != nil {
		common.RenderError(c, http.StatusUnprocessableEntity, err)
		return
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

}

func stop(c *gin.Context) {

}

func getApplication(c *gin.Context, appID string) (app *entities.Application, err error) {
	currentUser := common.MustGetCurrentUser(c)
	app, err = entities.GetApplication(gin.H{"_id": appID, "user_id": currentUser.ID})
	return
}
