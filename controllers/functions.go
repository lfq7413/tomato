package controllers

import (
	"encoding/json"

	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
)

// FunctionsController ...
type FunctionsController struct {
	ObjectsController
}

// HandleCloudFunction ...
// @router /:functionName [post]
func (f *FunctionsController) HandleCloudFunction() {
	functionName := f.Ctx.Input.Param(":functionName")
	theFunction := rest.GetFunction(functionName)
	if theFunction == nil {
		// TODO 无效函数
		return
	}

	params := types.M{}
	if f.Ctx.Input.RequestBody != nil {
		err := json.Unmarshal(f.Ctx.Input.RequestBody, &params)
		if err != nil {
			// TODO 参数错误
			return
		}
	}

	f.Auth.IsMaster = true
	request := rest.RequestInfo{
		Auth:      f.Auth,
		NewObject: params,
	}
	resp := theFunction(request)

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
