package files

import "github.com/lfq7413/tomato/config"

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
func CreateFile(filename string, data []byte, contentType string) map[string]interface{} {
	adapter.createFile(filename, data, contentType)
	return nil
}

// DeleteFile ...
func DeleteFile(filename string) {
	adapter.deleteFile(filename)
}

// ExpandFilesInObject ...
func ExpandFilesInObject(object interface{}) {

}

type filesAdapter interface {
	createFile(filename string, data []byte, contentType string)
	deleteFile(filename string)
	getFileData(filename string) []byte
	getFileLocation(filename string) string
}
