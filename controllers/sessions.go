package controllers

import (
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// SessionsController ...
type SessionsController struct {
	ObjectsController
}

// HandleFind ...
// @router / [get]
func (s *SessionsController) HandleFind() {
	s.ClassName = "_Session"
	s.ObjectsController.HandleFind()
}

// HandleGet ...
// @router /:objectId [get]
func (s *SessionsController) HandleGet() {
	s.ClassName = "_Session"
	s.ObjectID = s.Ctx.Input.Param(":objectId")
	s.ObjectsController.HandleGet()
}

// HandleCreate ...
// @router / [post]
func (s *SessionsController) HandleCreate() {
	s.ClassName = "_Session"
	s.ObjectsController.HandleCreate()
}

// HandleUpdate ...
// @router /:objectId [put]
func (s *SessionsController) HandleUpdate() {
	s.ClassName = "_Session"
	s.ObjectID = s.Ctx.Input.Param(":objectId")
	s.ObjectsController.HandleUpdate()
}

// HandleDelete ...
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

// HandleGetMe ...
// @router /me [get]
func (s *SessionsController) HandleGetMe() {
	if s.Info == nil || s.Info.SessionToken == "" {
		// TODO 需要 SessionToken
		return
	}
	where := types.M{
		"sessionToken": s.Info.SessionToken,
	}
	// TODO 处理错误
	response, _ := rest.Find(rest.Master(), "_Session", where, types.M{})
	if utils.HasResults(response) == false {
		// TODO 未找到 Session
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
