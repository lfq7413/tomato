package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// ClassesController 对象操作 API 的基础结构
// 处理 /classes 接口的所有请求，处理内部类的部分请求
// ClassName 要操作的类名
// ObjectID 要操作的对象 id
type ClassesController struct {
	BaseController
	ClassName string
	ObjectID  string
}

// HandleCreate 处理对象创建请求，返回对象 id 与对象位置
// @router /:className [post]
func (o *ClassesController) HandleCreate() {

	if o.ClassName == "" {
		o.ClassName = o.Ctx.Input.Param(":className")
	}

	if o.JSONBody == nil {
		o.Data["json"] = errs.ErrorMessageToMap(errs.InvalidJSON, "request body is empty")
		o.ServeJSON()
		return
	}

	result, err := rest.Create(o.Auth, o.ClassName, o.JSONBody, o.Info.ClientSDK)
	if err != nil {
		o.Data["json"] = errs.ErrorToMap(err)
		o.ServeJSON()
		return
	}

	o.Data["json"] = result["response"]
	o.Ctx.Output.SetStatus(201)
	o.Ctx.Output.Header("location", result["location"].(string))
	o.ServeJSON()

}

// HandleGet 处理查询指定对象请求，返回查询到的对象
// @router /:className/:objectId [get]
func (o *ClassesController) HandleGet() {
	if o.ClassName == "" {
		o.ClassName = o.Ctx.Input.Param(":className")
	}
	if o.ObjectID == "" {
		o.ObjectID = o.Ctx.Input.Param(":objectId")
	}
	options := types.M{}
	if o.GetString("keys") != "" {
		options["keys"] = o.GetString("keys")
	}
	if o.GetString("include") != "" {
		options["include"] = o.GetString("include")
	}
	response, err := rest.Get(o.Auth, o.ClassName, o.ObjectID, options, o.Info.ClientSDK)

	if err != nil {
		o.Data["json"] = errs.ErrorToMap(err)
		o.ServeJSON()
		return
	}

	results := utils.A(response["results"])
	if results == nil && len(results) == 0 {
		o.Data["json"] = errs.ErrorMessageToMap(errs.ObjectNotFound, "Object not found.")
		o.ServeJSON()
		return
	}

	result := utils.M(results[0])

	if o.ClassName == "_User" {
		delete(result, "sessionToken")
		if o.Auth.User != nil && result["objectId"].(string) == o.Auth.User["objectId"].(string) {
			// 重新设置 session token
			result["sessionToken"] = o.Info.SessionToken
		}
	}

	o.Data["json"] = result
	o.ServeJSON()

}

// HandleUpdate 处理更新指定对象请求
// @router /:className/:objectId [put]
func (o *ClassesController) HandleUpdate() {

	if o.ClassName == "" {
		o.ClassName = o.Ctx.Input.Param(":className")
	}
	if o.ObjectID == "" {
		o.ObjectID = o.Ctx.Input.Param(":objectId")
	}

	if o.JSONBody == nil {
		o.Data["json"] = errs.ErrorMessageToMap(errs.InvalidJSON, "request body is empty")
		o.ServeJSON()
		return
	}

	result, err := rest.Update(o.Auth, o.ClassName, o.ObjectID, o.JSONBody, o.Info.ClientSDK)
	if err != nil {
		o.Data["json"] = errs.ErrorToMap(err)
		o.ServeJSON()
		return
	}

	o.Data["json"] = result["response"]
	o.ServeJSON()

}

// HandleFind 处理查找对象请求
// @router /:className [get]
func (o *ClassesController) HandleFind() {
	if o.ClassName == "" {
		o.ClassName = o.Ctx.Input.Param(":className")
	}

	// 获取查询参数，并组装
	options := types.M{}
	if o.GetString("skip") != "" {
		if i, err := strconv.Atoi(o.GetString("skip")); err == nil {
			options["skip"] = i
		} else {
			o.Data["json"] = errs.ErrorMessageToMap(errs.InvalidQuery, "skip should be int")
			o.ServeJSON()
			return
		}
	}
	if o.GetString("limit") != "" {
		if i, err := strconv.Atoi(o.GetString("limit")); err == nil {
			options["limit"] = i
		} else {
			o.Data["json"] = errs.ErrorMessageToMap(errs.InvalidQuery, "limit should be int")
			o.ServeJSON()
			return
		}
	} else {
		options["limit"] = 100
	}
	if o.GetString("order") != "" {
		options["order"] = o.GetString("order")
	}
	if o.GetString("count") != "" {
		options["count"] = true
	}
	if o.GetString("keys") != "" {
		options["keys"] = o.GetString("keys")
	}
	if o.GetString("include") != "" {
		options["include"] = o.GetString("include")
	}
	if o.GetString("redirectClassNameForKey") != "" {
		options["redirectClassNameForKey"] = o.GetString("redirectClassNameForKey")
	}

	where := types.M{}
	if o.GetString("where") != "" {
		err := json.Unmarshal([]byte(o.GetString("where")), &where)
		if err != nil {
			o.Data["json"] = errs.ErrorMessageToMap(errs.InvalidJSON, "where should be valid json")
			o.ServeJSON()
			return
		}
	}

	response, err := rest.Find(o.Auth, o.ClassName, where, options, o.Info.ClientSDK)
	if err != nil {
		o.Data["json"] = errs.ErrorToMap(err)
		o.ServeJSON()
		return
	}
	if utils.HasResults(response) {
		results := utils.A(response["results"])
		for _, v := range results {
			result := utils.M(v)
			if result["sessionToken"] != nil && o.Info.SessionToken != "" {
				result["sessionToken"] = o.Info.SessionToken
			}
		}
	}

	o.Data["json"] = response
	o.ServeJSON()
}

// HandleDelete 处理删除指定对象请求
// @router /:className/:objectId [delete]
func (o *ClassesController) HandleDelete() {

	if o.ClassName == "" {
		o.ClassName = o.Ctx.Input.Param(":className")
	}
	if o.ObjectID == "" {
		o.ObjectID = o.Ctx.Input.Param(":objectId")
	}

	err := rest.Delete(o.Auth, o.ClassName, o.ObjectID, o.Info.ClientSDK)
	if err != nil {
		o.Data["json"] = errs.ErrorToMap(err)
		o.ServeJSON()
		return
	}

	o.Data["json"] = types.M{}
	o.ServeJSON()
}

// Get ...
// @router / [get]
func (o *ClassesController) Get() {
	o.Ctx.Output.SetStatus(405)
	o.Data["json"] = errs.ErrorMessageToMap(405, "Method Not Allowed")
	o.ServeJSON()
}

// Post ...
// @router / [post]
func (o *ClassesController) Post() {
	o.Ctx.Output.SetStatus(405)
	o.Data["json"] = errs.ErrorMessageToMap(405, "Method Not Allowed")
	o.ServeJSON()
}

// Delete ...
// @router / [delete]
func (o *ClassesController) Delete() {
	o.Ctx.Output.SetStatus(405)
	o.Data["json"] = errs.ErrorMessageToMap(405, "Method Not Allowed")
	o.ServeJSON()
}

// Put ...
// @router / [put]
func (o *ClassesController) Put() {
	o.Ctx.Output.SetStatus(405)
	o.Data["json"] = errs.ErrorMessageToMap(405, "Method Not Allowed")
	o.ServeJSON()
}
