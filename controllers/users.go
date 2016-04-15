package controllers

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// UsersController 处理 /users 接口的请求
type UsersController struct {
	ObjectsController
}

// HandleFind 处理查找用户请求
// @router / [get]
func (u *UsersController) HandleFind() {
	u.ClassName = "_User"
	u.ObjectsController.HandleFind()
}

// HandleGet 处理查找指定用户请求
// @router /:objectId [get]
func (u *UsersController) HandleGet() {
	u.ClassName = "_User"
	u.ObjectID = u.Ctx.Input.Param(":objectId")
	u.ObjectsController.HandleGet()
}

// HandleCreate 处理创建用户请求
// @router / [post]
func (u *UsersController) HandleCreate() {
	u.ClassName = "_User"
	u.ObjectsController.HandleCreate()
}

// HandleUpdate 处理更新用户信息请求
// 过滤对 /me 接口的 put 请求
// @router /:objectId [put]
func (u *UsersController) HandleUpdate() {
	objectID := u.Ctx.Input.Param(":objectId")
	if objectID == "me" {
		u.ObjectsController.Put()
		return
	}
	u.ClassName = "_User"
	u.ObjectID = objectID
	u.ObjectsController.HandleUpdate()
}

// HandleDelete 处理删除用户请求
// 过滤对 /me 接口的 delete 请求
// @router /:objectId [delete]
func (u *UsersController) HandleDelete() {
	objectID := u.Ctx.Input.Param(":objectId")
	if objectID == "me" {
		u.ObjectsController.Delete()
		return
	}
	u.ClassName = "_User"
	u.ObjectID = objectID
	u.ObjectsController.HandleDelete()
}

// HandleMe 处理获取当前用户信息的请求
// @router /me [get]
func (u *UsersController) HandleMe() {
	if u.Info == nil || u.Info.SessionToken == "" {
		u.Data["json"] = errs.ErrorMessageToMap(errs.InvalidSessionToken, "invalid session token")
		u.ServeJSON()
		return
	}
	sessionToken := u.Info.SessionToken
	where := types.M{
		"sessionToken": sessionToken,
	}
	option := types.M{
		"include": "user",
	}
	response, err := rest.Find(rest.Master(), "_Session", where, option)

	if err != nil {
		u.Data["json"] = errs.ErrorToMap(err)
		u.ServeJSON()
		return
	}

	if utils.HasResults(response) == false {
		u.Data["json"] = errs.ErrorMessageToMap(errs.InvalidSessionToken, "invalid session token")
		u.ServeJSON()
		return
	}
	results := utils.SliceInterface(response["results"])
	result := utils.MapInterface(results[0])
	if result["user"] == nil {
		u.Data["json"] = errs.ErrorMessageToMap(errs.InvalidSessionToken, "invalid session token")
		u.ServeJSON()
		return
	}
	user := utils.MapInterface(result["user"])
	user["sessionToken"] = sessionToken
	u.Data["json"] = user
	u.ServeJSON()
}

// Put ...
// @router / [put]
func (u *UsersController) Put() {
	u.ObjectsController.Put()
}

// Delete ...
// @router / [delete]
func (u *UsersController) Delete() {
	u.ObjectsController.Delete()
}
