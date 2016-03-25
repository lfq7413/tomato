package controllers

import (
	"github.com/lfq7413/tomato/files"
	"github.com/lfq7413/tomato/utils"
)

// FilesController ...
type FilesController struct {
	ObjectsController
}

// HandleGet ...
// @router /:appId/:filename [get]
func (f *FilesController) HandleGet() {
	filename := f.Ctx.Input.Param(":filename")
	data := files.GetFileData(filename)
	if data != nil {
		contentType := utils.LookupContentType(filename)
		f.Ctx.Output.SetStatus(200)
		f.Ctx.Output.ContentType(contentType)
		f.Ctx.Output.Body(data)
	} else {
		f.Ctx.Output.SetStatus(404)
		f.Ctx.Output.ContentType("text/plain")
		f.Ctx.Output.Body([]byte("File not found."))
	}
}

// HandleCreate ...
// @router /:filename [post]
func (f *FilesController) HandleCreate() {
	filename := f.Ctx.Input.Param(":filename")
	data := f.Ctx.Input.RequestBody
	if data == nil && len(data) == 0 {
		// TODO 无效上传
		return
	}
	if len(filename) > 128 {
		// TODO 文件名太长
		return
	}
	if utils.IsFileName(filename) == false {
		// TODO 无效文件名
		return
	}
	contentType := f.Ctx.Input.Header("Content-type")
	result := files.CreateFile(filename, data, contentType)
	if result != nil && result["url"] != "" {
		f.Ctx.Output.SetStatus(201)
		f.Ctx.Output.Header("Location", result["url"])
		f.Data["json"] = result
		f.ServeJSON()
	} else {
		// TODO 保存文件失败
	}
}

// HandleDelete ...
// @router /:filename [delete]
func (f *FilesController) HandleDelete() {
	if f.Auth.IsMaster == false {
		// TODO 权限不足
		return
	}
	filename := f.Ctx.Input.Param(":filename")
	err := files.DeleteFile(filename)
	if err != nil {
		// TODO 删除失败
		return
	}
	f.Data["json"] = "{}"
	f.ServeJSON()
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
