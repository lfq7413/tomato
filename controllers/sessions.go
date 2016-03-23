package controllers

// SessionsController ...
type SessionsController struct {
	ObjectsController
}

// HandleFind ...
// @router / [get]
func (s *SessionsController) HandleFind() {

}

// HandleGet ...
// @router /:objectId [get]
func (s *SessionsController) HandleGet() {

}

// HandleCreate ...
// @router / [post]
func (s *SessionsController) HandleCreate() {

}

// HandleUpdate ...
// @router /:objectId [put]
func (s *SessionsController) HandleUpdate() {

}

// HandleDelete ...
// @router /:objectId [delete]
func (s *SessionsController) HandleDelete() {
	objectID := s.Ctx.Input.Param(":objectId")
	if objectID == "me" {
		s.ObjectsController.Delete()
		return
	}
}

// HandleGetMe ...
// @router /me [get]
func (s *SessionsController) HandleGetMe() {

}

// HandleUpdateMe ...
// @router /me [put]
func (s *SessionsController) HandleUpdateMe() {

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
