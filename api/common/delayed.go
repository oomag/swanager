package common

import (
	"fmt"
	"net/http"

	"github.com/dokkur/swanager/core/entities"
	"github.com/gin-gonic/gin"
)

// AsyncJobContext - context for delayed job
type AsyncJobContext struct {
	User       *entities.User
	Job        *entities.Job
	Process    func() (interface{}, error)
	GinContext *gin.Context
}

// RunAsync starts delayed job
func RunAsync(context AsyncJobContext) {
	job, err := entities.CreateJob(context.User)
	if err != nil {
		RenderError(context.GinContext, http.StatusInternalServerError, err.Error())
		return
	}

	context.Job = job

	go process(&context)

	context.GinContext.JSON(http.StatusOK, gin.H{
		"job_id": job.ID,
		"url":    fmt.Sprintf("http://%s/api/v1/jobs/%s", context.GinContext.Request.Host, job.ID),
	})
}

func process(context *AsyncJobContext) {
	result, err := context.Process()
	if err != nil {
		context.Job.SetState(entities.JobStateError, err)
		return
	}
	context.Job.SetState(entities.JobStateSuccess, result)
}
