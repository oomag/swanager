package stat

import (
	"net/http"
	"time"

	"github.com/da4nik/swanager/api/common"
	"github.com/da4nik/swanager/core/entities"
	"github.com/da4nik/swanager/core/swarm"
	"github.com/gin-gonic/gin"
)

type statService struct {
	UserEmail       string    `json:"user_email"`
	ApplicationName string    `json:"application_name"`
	ApplicationID   string    `json:"application_id"`
	ServiceID       string    `json:"service_id"`
	ServiceName     string    `json:"service_name"`
	ReplicaID       string    `json:"replica_id"`
	Status          string    `json:"status"`
	Timestamp       time.Time `json:"timestamp"`
}

// GetRoutesForRouter adds resource routes to api router
func GetRoutesForRouter(router *gin.RouterGroup) {

	stats := router.Group("/stat")
	{
		stats.GET("", common.AuthLocal(), stat)
	}
}

func stat(c *gin.Context) {
	apps, err := entities.GetApplications(gin.H{})
	if err != nil {
		common.RenderError(c, http.StatusInternalServerError, err)
		return
	}

	users, err := entities.GetUsers(nil)
	if err != nil {
		common.RenderError(c, http.StatusInternalServerError, err)
		return
	}
	usersMap := getUsersMap(users)

	result := make([]statService, 0)
	for _, app := range apps {
		app.LoadServices()
		for _, service := range app.Services {
			swarm.GetServiceStatuses(&service)

			for _, status := range service.Status {
				result = append(result, statService{
					UserEmail:       usersMap[app.UserID].Email,
					ApplicationID:   app.ID,
					ApplicationName: app.Name,
					ServiceID:       service.ID,
					ServiceName:     service.Name,
					ReplicaID:       status.ReplicaID,
					Status:          status.Status,
					Timestamp:       status.Timestamp,
				})
			}
		}

	}

	c.JSON(http.StatusOK, gin.H{"stats": result})
}

func getUsersMap(users []entities.User) (result map[string]entities.User) {
	result = map[string]entities.User{}
	for _, user := range users {
		result[user.ID] = user
	}
	return
}
