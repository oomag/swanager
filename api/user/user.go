package user

import (
	"net/http"

	"github.com/da4nik/swanager/api/common"
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
		common.RenderError(c, http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func create(c *gin.Context) {
	var userRequest userCreate
	if err := c.BindJSON(&userRequest); err != nil {
		common.RenderError(c, http.StatusBadRequest, err)
		return
	}

	if userRequest.Password != userRequest.PasswordConfirmation {
		common.RenderError(c, http.StatusBadRequest, "Password and confirmation are not match.")
		return
	}

	user := entities.User{
		Email:    userRequest.Email,
		Password: userRequest.Password,
	}

	if err := user.Save(); err != nil {
		common.RenderError(c, http.StatusUnprocessableEntity, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}
