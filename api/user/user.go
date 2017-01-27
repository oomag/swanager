package user

import (
	"fmt"
	"net/http"

	"github.com/da4nik/swanager/core/entities"
	"github.com/gin-gonic/gin"
)

type userCreate struct {
	Email                string
	Password             string
	PasswordConfirmation string `json:"password_confirmation"`
}

// GetRoutesForRouter adds resource routes to api router
func GetRoutesForRouter(router *gin.RouterGroup) {

	apps := router.Group("/users")
	{
		apps.POST("", create)
	}

	app := apps.Group("/:user_id")
	{
		app.GET("", show)
	}
}

func show(c *gin.Context) {
	user, err := entities.GetUser(c.Param("user_id"))
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func create(c *gin.Context) {
	var userRequest userCreate
	if err := c.BindJSON(&userRequest); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if userRequest.Password != userRequest.PasswordConfirmation {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Password and confirmation are not match."))
		return
	}

	user := entities.User{
		Email:    userRequest.Email,
		Password: userRequest.Password,
	}

	if err := user.Save(); err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}
