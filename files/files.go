package files

import (
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/utils"
)

var adapter filesAdapter

// init 初始化文件处理模块
// 当前只有本地文件存储模块
// 后续可增加数据库文件存储、第三方网络文件存储模块
func init() {
	a := config.TConfig.FileAdapter
	if a == "disk" {
		adapter = &diskAdapter{}
	} else {
		adapter = &diskAdapter{}
	}
}

// GetFileData 获取文件数据
func GetFileData(filename string) []byte {
	return adapter.getFileData(filename)
}

// CreateFile 创建文件，返回文件地址与文件名
func CreateFile(filename string, data []byte, contentType string) map[string]string {
	extname := utils.ExtName(filename)
	if extname == "" && contentType != "" && utils.LookupExtension(contentType) != "" {
		filename = filename + "." + utils.LookupExtension(contentType)
	} else if extname != "" && contentType == "" {
		contentType = utils.LookupContentType(filename)
	}

	filename = utils.CreateToken() + "-" + filename
	location := adapter.getFileLocation(filename)

	err := adapter.createFile(filename, data, contentType)

	if err != nil {
		return nil
	}
	return map[string]string{
		"url":  location,
		"name": filename,
	}
}

// DeleteFile 删除文件
func DeleteFile(filename string) error {
	return adapter.deleteFile(filename)
}

// ExpandFilesInObject 展开文件对象
// 展开之后的文件对象如下
// {
// 	"__type": "File",
// 	"url": "http://example.com/pic.jpg",
// 	"name": "pic.jpg",
// }
func ExpandFilesInObject(object interface{}) {
	if utils.SliceInterface(object) != nil {
		objs := utils.SliceInterface(object)
		for _, obj := range objs {
			ExpandFilesInObject(obj)
		}
	}

	if utils.MapInterface(object) == nil {
		return
	}

	obj := utils.MapInterface(object)

	for _, v := range obj {
		fileObject := utils.MapInterface(v)
		if fileObject != nil && fileObject["__type"] == "File" {
			if fileObject["url"] != nil {
				continue
			}
			filename := utils.String(fileObject["name"])
			fileObject["url"] = adapter.getFileLocation(filename)
		}
	}
}

// filesAdapter 规定了文件存储模块需要实现的接口
type filesAdapter interface {
	createFile(filename string, data []byte, contentType string) error
	deleteFile(filename string) error
	getFileData(filename string) []byte
	getFileLocation(filename string) string
}
