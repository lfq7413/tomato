package controllers

import (
	"regexp"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// UsersController 处理 /users 接口的请求
type UsersController struct {
	ClassesController
}

// HandleFind 处理查找用户请求
// @router / [get]
func (u *UsersController) HandleFind() {
	u.ClassName = "_User"
	u.ClassesController.HandleFind()
}

// HandleGet 处理查找指定用户请求
// @router /:objectId [get]
func (u *UsersController) HandleGet() {
	u.ClassName = "_User"
	u.ObjectID = u.Ctx.Input.Param(":objectId")
	u.ClassesController.HandleGet()
}

// HandleCreate 处理创建用户请求
// @router / [post]
func (u *UsersController) HandleCreate() {
	u.ClassName = "_User"
	u.ClassesController.HandleCreate()
}

// HandleUpdate 处理更新用户信息请求
// 过滤对 /me 接口的 put 请求
// @router /:objectId [put]
func (u *UsersController) HandleUpdate() {
	objectID := u.Ctx.Input.Param(":objectId")
	if objectID == "me" {
		u.ClassesController.Put()
		return
	}
	u.ClassName = "_User"
	u.ObjectID = objectID
	u.ClassesController.HandleUpdate()
}

// HandleDelete 处理删除用户请求
// 过滤对 /me 接口的 delete 请求
// @router /:objectId [delete]
func (u *UsersController) HandleDelete() {
	objectID := u.Ctx.Input.Param(":objectId")
	if objectID == "me" {
		u.ClassesController.Delete()
		return
	}
	u.ClassName = "_User"
	u.ObjectID = objectID
	u.ClassesController.HandleDelete()
}

// HandleMe 处理获取当前用户信息的请求
// @router /me [get]
func (u *UsersController) HandleMe() {
	if u.Info == nil || u.Info.SessionToken == "" {
		u.HandleError(errs.E(errs.InvalidSessionToken, "invalid session token"), 0)
		return
	}
	sessionToken := u.Info.SessionToken
	where := types.M{
		"sessionToken": sessionToken,
	}
	option := types.M{
		"include": "user",
	}
	response, err := rest.Find(rest.Master(), "_Session", where, option, u.Info.ClientSDK)

	if err != nil {
		u.HandleError(err, 0)
		return
	}

	if utils.HasResults(response) == false {
		u.HandleError(errs.E(errs.InvalidSessionToken, "invalid session token"), 0)
		return
	}
	results := utils.A(response["results"])
	result := utils.M(results[0])
	if result["user"] == nil {
		u.HandleError(errs.E(errs.InvalidSessionToken, "invalid session token"), 0)
		return
	}
	user := utils.M(result["user"])
	user["sessionToken"] = sessionToken

	// 删除隐藏字段
	for key := range user {
		if key != "__type" {
			b, _ := regexp.MatchString("^[A-Za-z][0-9A-Za-z_]*$", key)
			if b == false {
				delete(user, key)
			}
		}
	}

	u.Data["json"] = user
	u.ServeJSON()
}

// Put ...
// @router / [put]
func (u *UsersController) Put() {
	u.ClassesController.Put()
}

// Delete ...
// @router / [delete]
func (u *UsersController) Delete() {
	u.ClassesController.Delete()
}
