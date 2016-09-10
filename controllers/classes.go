package controllers

import (
	"encoding/json"
	"errors"
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
func (c *ClassesController) HandleCreate() {
	if c.ClassName == "" {
		c.ClassName = c.Ctx.Input.Param(":className")
	}
	if c.JSONBody == nil {
		c.HandleError(errs.E(errs.InvalidJSON, "request body is empty"), 0)
		return
	}

	result, err := rest.Create(c.Auth, c.ClassName, c.JSONBody, c.Info.ClientSDK)
	if err != nil {
		c.HandleError(err, 0)
		return
	}

	c.Data["json"] = result["response"]
	c.Ctx.Output.SetStatus(201)
	c.Ctx.Output.Header("Location", utils.S(result["location"]))
	c.ServeJSON()
}

// HandleGet 处理查询指定对象请求，返回查询到的对象
// @router /:className/:objectId [get]
func (c *ClassesController) HandleGet() {
	if c.ClassName == "" {
		c.ClassName = c.Ctx.Input.Param(":className")
	}
	if c.ObjectID == "" {
		c.ObjectID = c.Ctx.Input.Param(":objectId")
	}

	options := types.M{}
	if c.GetString("keys") != "" {
		options["keys"] = c.GetString("keys")
	}
	if c.GetString("include") != "" {
		options["include"] = c.GetString("include")
	}

	response, err := rest.Get(c.Auth, c.ClassName, c.ObjectID, options, c.Info.ClientSDK)
	if err != nil {
		c.HandleError(err, 0)
		return
	}

	results := utils.A(response["results"])
	if results == nil || len(results) == 0 {
		c.HandleError(errs.E(errs.ObjectNotFound, "Object not found."), 0)
		return
	}

	result := utils.M(results[0])

	if c.ClassName == "_User" {
		delete(result, "sessionToken")
		if c.Auth.User != nil && utils.S(result["objectId"]) == utils.S(c.Auth.User["objectId"]) {
			// 重新设置 session token
			result["sessionToken"] = c.Info.SessionToken
		}
	}

	c.Data["json"] = result
	c.ServeJSON()
}

// HandleUpdate 处理更新指定对象请求
// @router /:className/:objectId [put]
func (c *ClassesController) HandleUpdate() {
	if c.ClassName == "" {
		c.ClassName = c.Ctx.Input.Param(":className")
	}
	if c.ObjectID == "" {
		c.ObjectID = c.Ctx.Input.Param(":objectId")
	}
	if c.JSONBody == nil {
		c.HandleError(errs.E(errs.InvalidJSON, "request body is empty"), 0)
		return
	}

	result, err := rest.Update(c.Auth, c.ClassName, c.ObjectID, c.JSONBody, c.Info.ClientSDK)
	if err != nil {
		c.HandleError(err, 0)
		return
	}

	c.Data["json"] = result["response"]
	c.ServeJSON()
}

// HandleFind 处理查找对象请求
// @router /:className [get]
func (c *ClassesController) HandleFind() {
	if c.ClassName == "" {
		c.ClassName = c.Ctx.Input.Param(":className")
	}

	// 获取查询参数，并组装
	options := types.M{}
	if c.GetString("skip") != "" {
		if i, err := strconv.Atoi(c.GetString("skip")); err == nil {
			options["skip"] = i
		} else {
			c.HandleError(errs.E(errs.InvalidQuery, "skip should be int"), 0)
			return
		}
	}
	if c.GetString("limit") != "" {
		if i, err := strconv.Atoi(c.GetString("limit")); err == nil {
			options["limit"] = i
		} else {
			c.HandleError(errs.E(errs.InvalidQuery, "limit should be int"), 0)
			return
		}
	} else {
		options["limit"] = 100
	}
	if c.GetString("order") != "" {
		options["order"] = c.GetString("order")
	}
	if c.GetString("count") != "" {
		options["count"] = true
	}
	if c.GetString("keys") != "" {
		options["keys"] = c.GetString("keys")
	}
	if c.GetString("include") != "" {
		options["include"] = c.GetString("include")
	}
	if c.GetString("redirectClassNameForKey") != "" {
		options["redirectClassNameForKey"] = c.GetString("redirectClassNameForKey")
	}

	where := types.M{}
	if c.GetString("where") != "" {
		err := json.Unmarshal([]byte(c.GetString("where")), &where)
		if err != nil {
			c.HandleError(errs.E(errs.InvalidJSON, "where should be valid json"), 0)
			return
		}
	}

	response, err := rest.Find(c.Auth, c.ClassName, where, options, c.Info.ClientSDK)
	if err != nil {
		c.HandleError(err, 0)
		return
	}
	if utils.HasResults(response) {
		results := utils.A(response["results"])
		for _, v := range results {
			result := utils.M(v)
			if result["sessionToken"] != nil && c.Info.SessionToken != "" {
				result["sessionToken"] = c.Info.SessionToken
			}
		}
	}

	c.Data["json"] = response
	c.ServeJSON()
}

// HandleDelete 处理删除指定对象请求
// @router /:className/:objectId [delete]
func (c *ClassesController) HandleDelete() {
	if c.ClassName == "" {
		c.ClassName = c.Ctx.Input.Param(":className")
	}
	if c.ObjectID == "" {
		c.ObjectID = c.Ctx.Input.Param(":objectId")
	}

	err := rest.Delete(c.Auth, c.ClassName, c.ObjectID, c.Info.ClientSDK)
	if err != nil {
		c.HandleError(err, 0)
		return
	}

	c.Data["json"] = types.M{}
	c.ServeJSON()
}

// Get ...
// @router / [get]
func (c *ClassesController) Get() {
	c.HandleError(errors.New("Method Not Allowed"), 405)
}

// Post ...
// @router / [post]
func (c *ClassesController) Post() {
	c.HandleError(errors.New("Method Not Allowed"), 405)
}

// Delete ...
// @router / [delete]
func (c *ClassesController) Delete() {
	c.HandleError(errors.New("Method Not Allowed"), 405)
}

// Put ...
// @router / [put]
func (c *ClassesController) Put() {
	c.HandleError(errors.New("Method Not Allowed"), 405)
}
