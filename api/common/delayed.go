package common

import (
	"fmt"
	"net/http"

	"github.com/da4nik/swanager/core/entities"
	"github.com/gin-gonic/gin"
)

// DelayedJobContext - context for delayed job
type DelayedJobContext struct {
	User       *entities.User
	Job        *entities.Job
	Process    func() (string, error)
	GinContext *gin.Context
}

// RunDelayed starts delayed job
func RunDelayed(context DelayedJobContext) {
	job, err := entities.CreateJob(context.User)
	if err != nil {
		RenderError(context.GinContext, http.StatusInternalServerError, err.Error())
		return
	}

	context.Job = job

	go process(&context)

	// c.JSON(http.StatusOK, gin.H{"application": app})
	context.GinContext.JSON(http.StatusOK, gin.H{
		"job_id": job.ID,
		"url":    fmt.Sprintf("http://%s/api/v1/jobs/%s", context.GinContext.Request.Host, job.ID),
	})
}

func process(context *DelayedJobContext) {
	result, err := context.Process()
	if err != nil {
		context.Job.SetState(entities.JobStateError, err.Error())
		return
	}
	context.Job.SetState(entities.JobStateSuccess, result)
}
