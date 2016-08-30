package controllers

import (
	"github.com/lfq7413/tomato/cloud"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
)

// FunctionsController 处理 /functions 接口的请求
type FunctionsController struct {
	ClassesController
}

// HandleCloudFunction 执行指定的云函数
// 返回数据格式如下：
// {
// 	"result":"func res"
// }
// @router /:functionName [post]
func (f *FunctionsController) HandleCloudFunction() {
	functionName := f.Ctx.Input.Param(":functionName")
	theFunction := cloud.GetFunction(functionName)
	theValidator := cloud.GetValidator(functionName)
	if theFunction == nil {
		f.Data["json"] = errs.ErrorMessageToMap(errs.ScriptFailed, "Invalid function.")
		f.ServeJSON()
		return
	}

	if f.JSONBody == nil {
		f.JSONBody = types.M{}
	}

	// TODO 补全 Headers
	request := cloud.FunctionRequest{
		Params:         f.JSONBody,
		Master:         false,
		InstallationID: f.Info.InstallationID,
	}
	if f.Auth != nil {
		request.Master = f.Auth.IsMaster
		request.User = f.Auth.User
	}

	if theValidator != nil {
		result := theValidator(request)
		if result == false {
			f.Data["json"] = errs.ErrorMessageToMap(errs.ValidationError, "Validation failed.")
			f.ServeJSON()
			return
		}
	}

	response := &cloud.FunctionResponse{}
	theFunction(request, response)
	if response.Err != nil {
		f.Data["json"] = errs.ErrorToMap(response.Err)
		f.ServeJSON()
		return
	}

	f.Data["json"] = response.Response
	f.ServeJSON()
}

// Get ...
// @router / [get]
func (f *FunctionsController) Get() {
	f.ClassesController.Get()
}

// Post ...
// @router / [post]
func (f *FunctionsController) Post() {
	f.ClassesController.Post()
}

// Delete ...
// @router / [delete]
func (f *FunctionsController) Delete() {
	f.ClassesController.Delete()
}

// Put ...
// @router / [put]
func (f *FunctionsController) Put() {
	f.ClassesController.Put()
}
