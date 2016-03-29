package controllers

// SchemasController ...
type SchemasController struct {
	ObjectsController
}

// Get ...
// @router / [get]
func (s *SchemasController) Get() {
	s.ObjectsController.Get()
}

// Post ...
// @router / [post]
func (s *SchemasController) Post() {
	s.ObjectsController.Post()
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
