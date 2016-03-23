package controllers

// FilesController ...
type FilesController struct {
	ObjectsController
}

// HandleGet ...
// @router /:appId/:filename [get]
func (f *FilesController) HandleGet() {
	// TODO
}

// HandleCreate ...
// @router /:filename [post]
func (f *FilesController) HandleCreate() {
	// TODO
}

// HandleDelete ...
// @router /:filename [delete]
func (f *FilesController) HandleDelete() {
	// TODO
}

// Get ...
// @router / [get]
func (f *FilesController) Get() {
	f.ObjectsController.Get()
}

// Post ...
// @router / [post]
func (f *FilesController) Post() {
	f.ObjectsController.Post()
}

// Put ...
// @router / [put]
func (f *FilesController) Put() {
	f.ObjectsController.Put()
}

// Delete ...
// @router / [delete]
func (f *FilesController) Delete() {
	f.ObjectsController.Delete()
}
