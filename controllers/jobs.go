package controllers

import (
	"github.com/lfq7413/tomato/cloud"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/job"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// JobsController 处理 /jobs 接口的请求
type JobsController struct {
	ClassesController
}

// HandleCloudJob 执行后台任务
// @router /:jobName [post]
func (j *JobsController) HandleCloudJob() {
	if j.Auth.IsMaster == false {
		j.Data["json"] = errs.ErrorMessageToMap(errs.OperationForbidden, "need master key")
		j.ServeJSON()
		return
	}
	jobName := j.Ctx.Input.Param(":jobName")
	j.runJob(jobName)
}

// HandlePost ...
// @router / [post]
func (j *JobsController) HandlePost() {
	if j.Auth.IsMaster == false {
		j.Data["json"] = errs.ErrorMessageToMap(errs.OperationForbidden, "need master key")
		j.ServeJSON()
		return
	}
	jobName := utils.S(j.JSONBody["jobName"])
	j.runJob(jobName)
}

func (j *JobsController) runJob(jobName string) {
	jobFunction := cloud.GetJob(jobName)
	if jobFunction == nil {
		j.Data["json"] = errs.ErrorMessageToMap(errs.ScriptFailed, "Invalid job.")
		j.ServeJSON()
		return
	}
	jobHandler := job.NewjobStatus()

	if j.JSONBody == nil {
		j.JSONBody = types.M{}
	}

	request := cloud.JobRequest{
		Params:  j.JSONBody,
		JobName: jobName,
		Headers: j.Ctx.Input.Context.Request.Header,
	}
	response := cloud.JobResponse{
		JobStatus: jobHandler,
	}
	jobStatus := jobHandler.SetRunning(jobName, j.JSONBody)
	request.JobID = utils.S(jobStatus["objectId"])

	go jobFunction(request, response)

	j.Ctx.Output.Header("X-Parse-Job-Status-Id", utils.S(jobStatus["objectId"]))
	j.Data["json"] = types.M{}
	j.ServeJSON()
}

// Get ...
// @router / [get]
func (j *JobsController) Get() {
	j.ClassesController.Get()
}

// Delete ...
// @router / [delete]
func (j *JobsController) Delete() {
	j.ClassesController.Delete()
}

// Put ...
// @router / [put]
func (j *JobsController) Put() {
	j.ClassesController.Put()
}
