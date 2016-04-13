package controllers

import (
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// UsersController ...
type UsersController struct {
	ObjectsController
}

// HandleFind ...
// @router / [get]
func (u *UsersController) HandleFind() {
	u.ClassName = "_User"
	u.ObjectsController.HandleFind()
}

// HandleGet ...
// @router /:objectId [get]
func (u *UsersController) HandleGet() {
	u.ClassName = "_User"
	u.ObjectID = u.Ctx.Input.Param(":objectId")
	u.ObjectsController.HandleGet()
}

// HandleCreate ...
// @router / [post]
func (u *UsersController) HandleCreate() {
	u.ClassName = "_User"
	u.ObjectsController.HandleCreate()
}

// HandleUpdate ...
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

// HandleDelete ...
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

// HandleMe ...
// @router /me [get]
func (u *UsersController) HandleMe() {
	if u.Info == nil || u.Info.SessionToken == "" {
		// TODO SessionToken 无效
		return
	}
	sessionToken := u.Info.SessionToken
	where := types.M{
		"sessionToken": sessionToken,
	}
	option := types.M{
		"include": "user",
	}
	response := rest.Find(rest.Master(), "_Session", where, option)
	if utils.HasResults(response) == false {
		// TODO SessionToken 无效
		return
	}
	results := utils.SliceInterface(response["results"])
	result := utils.MapInterface(results[0])
	if result["user"] == nil {
		// TODO SessionToken 无效
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
