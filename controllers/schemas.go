package controllers

// SchemasController ...
type SchemasController struct {
	ObjectsController
}

// HandleFind ...
// @router / [get]
func (s *SchemasController) HandleFind() {
	s.ObjectsController.Get()
}

// HandleGet ...
// @router /:className [get]
func (s *SchemasController) HandleGet() {
	s.ObjectsController.Get()
}

// HandleCreate ...
// @router / [post]
func (s *SchemasController) HandleCreate() {
	s.ObjectsController.Post()
}

// HandleUpdate ...
// @router /:className [put]
func (s *SchemasController) HandleUpdate() {
	s.ObjectsController.Put()
}

// HandleDelete ...
// @router /:className [delete]
func (s *SchemasController) HandleDelete() {
	s.ObjectsController.Delete()
}

// Delete ...
// @router / [delete]
func (s *SchemasController) Delete() {
	s.ObjectsController.Delete()
}

// Put ...
// @router / [put]
func (s *SchemasController) Put() {
	s.ObjectsController.Put()
}
