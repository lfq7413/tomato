package controllers

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// SessionsController 处理 /sessions 接口的请求
type SessionsController struct {
	ObjectsController
}

// HandleFind 处理查找 session 请求
// @router / [get]
func (s *SessionsController) HandleFind() {
	s.ClassName = "_Session"
	s.ObjectsController.HandleFind()
}

// HandleGet 处理获取指定 session 请求
// @router /:objectId [get]
func (s *SessionsController) HandleGet() {
	s.ClassName = "_Session"
	s.ObjectID = s.Ctx.Input.Param(":objectId")
	s.ObjectsController.HandleGet()
}

// HandleCreate 处理 session 创建请求
// @router / [post]
func (s *SessionsController) HandleCreate() {
	s.ClassName = "_Session"
	s.ObjectsController.HandleCreate()
}

// HandleUpdate 处理更新指定 session 请求
// @router /:objectId [put]
func (s *SessionsController) HandleUpdate() {
	s.ClassName = "_Session"
	s.ObjectID = s.Ctx.Input.Param(":objectId")
	s.ObjectsController.HandleUpdate()
}

// HandleDelete 处理删除指定 session 请求
// @router /:objectId [delete]
func (s *SessionsController) HandleDelete() {
	objectID := s.Ctx.Input.Param(":objectId")
	if objectID == "me" {
		s.ObjectsController.Delete()
		return
	}
	s.ClassName = "_Session"
	s.ObjectID = objectID
	s.ObjectsController.HandleDelete()
}

// HandleGetMe 处理当前请求 session
// @router /me [get]
func (s *SessionsController) HandleGetMe() {
	if s.Info == nil || s.Info.SessionToken == "" {
		s.Data["json"] = errs.ErrorMessageToMap(errs.InvalidSessionToken, "Session token required.")
		s.ServeJSON()
		return
	}
	where := types.M{
		"_session_token": s.Info.SessionToken,
	}
	response, err := rest.Find(rest.Master(), "_Session", where, types.M{})
	if err != nil {
		s.Data["json"] = errs.ErrorToMap(err)
		s.ServeJSON()
		return
	}
	if utils.HasResults(response) == false {
		s.Data["json"] = errs.ErrorMessageToMap(errs.InvalidSessionToken, "Session token not found.")
		s.ServeJSON()
		return
	}
	results := utils.SliceInterface(response["results"])
	s.Data["json"] = results[0]
	s.ServeJSON()
}

// HandleUpdateMe ...
// @router /me [put]
func (s *SessionsController) HandleUpdateMe() {
	// TODO
}

// Put ...
// @router / [put]
func (s *SessionsController) Put() {
	s.ObjectsController.Put()
}

// Delete ...
// @router / [delete]
func (s *SessionsController) Delete() {
	s.ObjectsController.Delete()
}
