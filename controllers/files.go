package controllers

import (
	"strconv"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/files"
	"github.com/lfq7413/tomato/utils"
)

// FilesController 处理 /files 接口的请求
// TODO 可能不需要权限验证
type FilesController struct {
	ClassesController
}

// HandleGet 处理下载文件请求
// @router /:appId/:filename [get]
func (f *FilesController) HandleGet() {
	filename := f.Ctx.Input.Param(":filename")
	data := files.GetFileData(filename)
	if data != nil {
		contentType := utils.LookupContentType(filename)
		f.Ctx.Output.SetStatus(200)
		f.Ctx.Output.Header("Content-Type", contentType)
		f.Ctx.Output.Header("Content-Length", strconv.Itoa(len(data)))
		f.Ctx.Output.Body(data)
	} else {
		f.Ctx.Output.SetStatus(404)
		f.Ctx.Output.Header("Content-Type", "text/plain")
		f.Ctx.Output.Body([]byte("File not found."))
	}
}

// HandleCreate 处理上传文件请求
// @router /:filename [post]
func (f *FilesController) HandleCreate() {
	filename := f.Ctx.Input.Param(":filename")
	data := f.Ctx.Input.RequestBody
	if data == nil && len(data) == 0 {
		f.Data["json"] = errs.ErrorMessageToMap(errs.FileSaveError, "Invalid file upload.")
		f.ServeJSON()
		return
	}
	if len(filename) > 128 {
		f.Data["json"] = errs.ErrorMessageToMap(errs.InvalidFileName, "Filename too long.")
		f.ServeJSON()
		return
	}
	if utils.IsFileName(filename) == false {
		f.Data["json"] = errs.ErrorMessageToMap(errs.InvalidFileName, "Filename contains invalid characters.")
		f.ServeJSON()
		return
	}
	contentType := f.Ctx.Input.Header("Content-type")
	result := files.CreateFile(filename, data, contentType)
	if result != nil && result["url"] != "" {
		f.Ctx.Output.SetStatus(201)
		f.Ctx.Output.Header("location", result["url"])
		f.Data["json"] = result
		f.ServeJSON()
	} else {
		f.Data["json"] = errs.ErrorMessageToMap(errs.FileSaveError, "Could not store file.")
		f.ServeJSON()
	}
}

// HandleDelete 处理删除文件请求
// @router /:filename [delete]
func (f *FilesController) HandleDelete() {
	if f.Auth.IsMaster == false {
		f.Data["json"] = errs.ErrorMessageToMap(errs.FileDeleteError, "This user is not allowed to delete file.")
		f.ServeJSON()
		return
	}
	filename := f.Ctx.Input.Param(":filename")
	err := files.DeleteFile(filename)
	if err != nil {
		f.Data["json"] = errs.ErrorMessageToMap(errs.FileDeleteError, "Could not delete file.")
		f.ServeJSON()
		return
	}
	f.Data["json"] = "{}"
	f.ServeJSON()
}

// Get ...
// @router / [get]
func (f *FilesController) Get() {
	f.ClassesController.Get()
}

// Post ...
// @router / [post]
func (f *FilesController) Post() {
	f.ClassesController.Post()
}

// Put ...
// @router / [put]
func (f *FilesController) Put() {
	f.ClassesController.Put()
}

// Delete ...
// @router / [delete]
func (f *FilesController) Delete() {
	f.ClassesController.Delete()
}
