package app

import (
	"fmt"
	"net/http"

	"github.com/da4nik/swanager/api/common"
	"github.com/da4nik/swanager/core/entities"
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
		app.PUT("", update)
		app.POST("/start", start)
		app.POST("/stop", stop)
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
	app, err := getApplication(c, c.Param("app_id"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	app.LoadServices()

	c.JSON(http.StatusOK, gin.H{"application": app})
}

func create(c *gin.Context) {
	currentUser := common.MustGetCurrentUser(c)
	var app entities.Application
	if err := c.BindJSON(&app); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	app.UserID = currentUser.ID
	if err := app.Save(); err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"application": app})
}

func update(c *gin.Context) {
	var newApp entities.Application
	if err := c.BindJSON(&newApp); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	app, err := getApplication(c, c.Param("app_id"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
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

func getApplication(c *gin.Context, appID string) (*entities.Application, error) {
	currentUser := common.MustGetCurrentUser(c)
	fmt.Println(currentUser.ID)
	fmt.Println(appID)
	app, err := entities.GetApplication(gin.H{"_id": appID, "user_id": currentUser.ID})
	// app, err := entities.GetApplication(gin.H{"_id": appID})
	if err != nil {
		return nil, fmt.Errorf("Application not found")
	}
	return app, nil
}
