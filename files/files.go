package files

import "github.com/lfq7413/tomato/config"
import "github.com/lfq7413/tomato/utils"

var adapter filesAdapter

func init() {
	a := config.TConfig.FileAdapter
	if a == "disk" {
		adapter = &diskAdapter{}
	} else {
		adapter = &diskAdapter{}
	}
}

// GetFileData ...
func GetFileData(filename string) []byte {
	return adapter.getFileData(filename)
}

// CreateFile ...
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

// DeleteFile ...
func DeleteFile(filename string) error {
	return adapter.deleteFile(filename)
}

// ExpandFilesInObject ...
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

type filesAdapter interface {
	createFile(filename string, data []byte, contentType string) error
	deleteFile(filename string) error
	getFileData(filename string) []byte
	getFileLocation(filename string) string
}
