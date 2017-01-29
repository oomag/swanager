package session

import (
	"net/http"

	"github.com/da4nik/swanager/api/common"
	"github.com/da4nik/swanager/core/auth"
	"github.com/gin-gonic/gin"
)

type loginMessage struct {
	Email      string
	Password   string
	RememberMe bool `json:"remember_me,omitempty"`
}

// GetRoutesForRouter adds resource routes to api router
func GetRoutesForRouter(router *gin.RouterGroup) {

	auth := router.Group("/session")
	{
		auth.POST("", common.Auth(false), login)
		auth.DELETE("", common.Auth(true), logout)
	}
}

func login(c *gin.Context) {
	var json loginMessage
	if err := c.BindJSON(&json); err != nil {
		common.RenderError(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := auth.WithEmailAndPassword(json.Email, json.Password)
	if err != nil {
		common.RenderError(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func logout(c *gin.Context) {
	user := common.MustGetCurrentUser(c)
	if err := auth.Deauthorize(user); err != nil {
		common.RenderError(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.AbortWithStatus(http.StatusOK)
}
