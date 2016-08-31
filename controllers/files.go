package controllers

import (
	"strconv"
	"strings"

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
	contentType := utils.LookupContentType(filename)
	if f.isFileStreamable() {
		s, err := files.GetFileStream(filename)
		if err != nil {
			f.Ctx.Output.SetStatus(404)
			f.Ctx.Output.Header("Content-Type", "text/plain")
			f.Ctx.Output.Body([]byte("File not found."))
			return
		}
		f.handleFileStream(s, contentType)
		return
	}
	data, err := files.GetFileData(filename)
	if err != nil {
		f.Ctx.Output.SetStatus(404)
		f.Ctx.Output.Header("Content-Type", "text/plain")
		f.Ctx.Output.Body([]byte("File not found."))
	} else {
		f.Ctx.Output.SetStatus(200)
		f.Ctx.Output.Header("Content-Type", contentType)
		f.Ctx.Output.Header("Content-Length", strconv.Itoa(len(data)))
		f.Ctx.Output.Body(data)
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

func (f *FilesController) isFileStreamable() bool {
	if f.Ctx.Input.Header("Range") == "" {
		return false
	}
	n := files.GetAdapterName()
	if n == "fileSystemAdapter" || n == "gridStoreAdapter" {
		return true
	}
	return false
}

func (f *FilesController) handleFileStream(stream files.FileStream, contentType string) {
	defer stream.Close()
	var bufferSize = 1024 * 1024
	var r = f.Ctx.Input.Header("Range")
	r = strings.Replace(r, "bytes=", "", -1)
	r = strings.Replace(r, " ", "", -1)
	var parts = strings.Split(r, "-")
	if len(parts) != 2 {
		f.fileNotFound()
		return
	}
	// 修正读取范围
	var start, end, chunksize int
	if parts[0] == "" {
		start = 0
	} else {
		start, _ = strconv.Atoi(parts[0])
	}
	if parts[1] == "" {
		end = int(stream.Size()) - 1
	} else {
		end, _ = strconv.Atoi(parts[1])
	}
	if start > end {
		f.fileNotFound()
		return
	}
	if start > int(stream.Size())-1 {
		f.fileNotFound()
		return
	}
	chunksize = end - start + 1

	if chunksize > bufferSize {
		end = start + bufferSize - 1
		chunksize = bufferSize
	}
	// 设置 http 头部
	f.Ctx.Output.SetStatus(206)
	f.Ctx.Output.Header("Content-Range", "bytes "+strconv.Itoa(start)+"-"+strconv.Itoa(end)+"/"+strconv.Itoa(int(stream.Size())))
	f.Ctx.Output.Header("Accept-Ranges", "bytes")
	f.Ctx.Output.Header("Content-Length", strconv.Itoa(chunksize))
	f.Ctx.Output.Header("Content-Type", contentType)
	// 读取数据
	_, err := stream.Seek(int64(start), 0)
	if err != nil {
		f.fileNotFound()
		return
	}

	var bufferAvail = 0
	var size = (end - start) + 1
	var totalbyteswanted = (end - start) + 1
	var totalbyteswritten = 0

	data := []byte{}
	buf := make([]byte, 1024)
	for {
		n, _ := stream.Read(buf)
		if n == 0 {
			break
		}
		bufferAvail += n
		if bufferAvail < size {
			data = append(data, buf[:n]...)
			totalbyteswritten += n
			size -= n
			bufferAvail -= n
		} else {
			data = append(data, buf[:size]...)
			totalbyteswritten += n
			bufferAvail -= size
		}
		if totalbyteswritten >= totalbyteswanted {
			break
		}
	}
	f.Ctx.Output.Body(data)
}

func (f *FilesController) fileNotFound() {
	f.Ctx.Output.SetStatus(404)
	f.Ctx.Output.Header("Content-Type", "text/plain")
	f.Ctx.Output.Body([]byte("File not found."))
}
