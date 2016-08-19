package controllers

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
)

// FunctionsController 处理 /functions 接口的请求
type FunctionsController struct {
	ObjectsController
}

// HandleCloudFunction 执行指定的云函数
// 返回数据格式如下：
// {
// 	"result":"func res"
// }
// @router /:functionName [post]
func (f *FunctionsController) HandleCloudFunction() {
	functionName := f.Ctx.Input.Param(":functionName")
	theFunction := rest.GetFunction(functionName)
	theValidator := rest.GetValidator(functionName)
	if theFunction == nil {
		f.Data["json"] = errs.ErrorMessageToMap(errs.ScriptFailed, "Invalid function.")
		f.ServeJSON()
		return
	}

	if f.JSONBody == nil {
		f.JSONBody = types.M{}
	}

	// TODO 补全 Headers
	request := rest.FunctionRequest{
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

	response := &functionResponse{}
	theFunction(request, response)
	if response.err != nil {
		f.Data["json"] = errs.ErrorToMap(response.err)
		f.ServeJSON()
		return
	}

	f.Data["json"] = response.response
	f.ServeJSON()
}

// Get ...
// @router / [get]
func (f *FunctionsController) Get() {
	f.ObjectsController.Get()
}

// Post ...
// @router / [post]
func (f *FunctionsController) Post() {
	f.ObjectsController.Post()
}

// Delete ...
// @router / [delete]
func (f *FunctionsController) Delete() {
	f.ObjectsController.Delete()
}

// Put ...
// @router / [put]
func (f *FunctionsController) Put() {
	f.ObjectsController.Put()
}

type functionResponse struct {
	response types.M
	err      error
}

func (f *functionResponse) Success(response interface{}) {
	f.response = types.M{
		"result": response,
	}
}

func (f *functionResponse) Error(code int, message string) {
	if code == 0 {
		code = errs.ScriptFailed
	}
	f.err = errs.E(code, message)
}
