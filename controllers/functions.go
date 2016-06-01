package controllers

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
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
	if theFunction == nil {
		f.Data["json"] = errs.ErrorMessageToMap(errs.ScriptFailed, "Invalid function.")
		f.ServeJSON()
		return
	}

	if f.JSONBody == nil {
		f.JSONBody = types.M{}
	}

	for k, v := range f.JSONBody {
		if value, ok := v.(map[string]interface{}); ok {
			if value["__type"].(string) == "Date" {
				f.JSONBody[k], _ = utils.StringtoTime(value["iso"].(string))
			}
		}
	}

	f.Auth.IsMaster = true
	request := rest.RequestInfo{
		Auth:      f.Auth,
		NewObject: f.JSONBody,
	}
	resp := theFunction(request)
	if resp == nil {
		f.Data["json"] = errs.ErrorMessageToMap(errs.ScriptFailed, "Call function fail.")
		f.ServeJSON()
		return
	}

	f.Data["json"] = resp
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
