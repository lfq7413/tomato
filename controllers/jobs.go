package controllers

import "encoding/json"
import (
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
)

// JobsController ...
type JobsController struct {
	ObjectsController
}

// HandleCloudJob ...
// @router /:jobName [post]
func (j *JobsController) HandleCloudJob() {
	if j.Auth.IsMaster == false {
		// TODO 需要 Master 权限
		return
	}
	jobName := j.Ctx.Input.Param(":jobName")
	theJob := rest.GetJob(jobName)
	if theJob == nil {
		// TODO 无效函数
		return
	}

	params := types.M{}
	if j.Ctx.Input.RequestBody != nil {
		err := json.Unmarshal(j.Ctx.Input.RequestBody, &params)
		if err != nil {
			// TODO 参数错误
			return
		}
	}

	request := rest.RequestInfo{
		Auth:      j.Auth,
		NewObject: params,
	}
	go theJob(request)

	j.Data["json"] = "{}"
	j.ServeJSON()

}

// Get ...
// @router / [get]
func (j *JobsController) Get() {
	j.ObjectsController.Get()
}

// Post ...
// @router / [post]
func (j *JobsController) Post() {
	j.ObjectsController.Post()
}

// Delete ...
// @router / [delete]
func (j *JobsController) Delete() {
	j.ObjectsController.Delete()
}

// Put ...
// @router / [put]
func (j *JobsController) Put() {
	j.ObjectsController.Put()
}
