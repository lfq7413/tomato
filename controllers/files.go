package controllers

import "github.com/lfq7413/tomato/files"

// FilesController ...
type FilesController struct {
	ObjectsController
}

// HandleGet ...
// @router /:appId/:filename [get]
func (f *FilesController) HandleGet() {
	// TODO
	filename := f.Ctx.Input.Param(":filename")
	data := files.GetFileData(filename)
	f.Ctx.Output.SetStatus(200)
	f.Ctx.Output.ContentType("")
	f.Ctx.Output.Body(data)
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
