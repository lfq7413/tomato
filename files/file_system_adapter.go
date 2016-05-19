package files

import (
	"net/url"
	"os"

	"github.com/astaxie/beego/utils"
	"github.com/lfq7413/tomato/config"
)

// fileSystemAdapter 本地文件存储模块
type fileSystemAdapter struct {
	filesDir string
}

func newFileSystemAdapter(filesSubDirectory string) *fileSystemAdapter {
	f := &fileSystemAdapter{
		filesDir: filesSubDirectory,
	}
	f.filesDir = filesSubDirectory
	if f.applicationDirExist() == false {
		err := f.mkdir(f.getApplicationDir())
		if err != nil {
			panic(err)
		}
	}

	return f
}

// createFile 在磁盘上创建文件
func (f *fileSystemAdapter) createFile(filename string, data []byte, contentType string) error {
	filepath := f.getLocalFilePath(filename)
	os.Remove(filepath)

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// deleteFile 从磁盘删除文件
func (f *fileSystemAdapter) deleteFile(filename string) error {
	filepath := f.getLocalFilePath(filename)
	err := os.Remove(filepath)
	if err != nil {
		return err
	}

	return nil
}

// getFileData 从磁盘获取文件数据，出错时返回空数据
func (f *fileSystemAdapter) getFileData(filename string) []byte {
	filepath := f.getLocalFilePath(filename)

	file, err := os.Open(filepath)
	if err != nil {
		return []byte{}
	}
	defer file.Close()

	data := []byte{}
	buf := make([]byte, 1024)
	for {
		n, _ := file.Read(buf)
		if n == 0 {
			break
		}
		data = append(data, buf[:n]...)
	}

	return data
}

// getFileLocation 获取文件路径
func (f *fileSystemAdapter) getFileLocation(filename string) string {
	return config.TConfig.ServerURL + "/files/" + config.TConfig.AppID + "/" + url.QueryEscape(filename)
}

func (f *fileSystemAdapter) getApplicationDir() string {
	if f.filesDir != "" {
		return utils.SelfDir() + string(os.PathSeparator) + "files" + string(os.PathSeparator) + f.filesDir
	}
	return utils.SelfDir() + string(os.PathSeparator) + "files"
}

func (f *fileSystemAdapter) applicationDirExist() bool {
	return utils.FileExists(f.getApplicationDir())
}

func (f *fileSystemAdapter) getLocalFilePath(filename string) string {
	applicationDir := f.getApplicationDir()
	if utils.FileExists(applicationDir) == false {
		f.mkdir(applicationDir)
	}
	return applicationDir + string(os.PathSeparator) + filename
}

func (f *fileSystemAdapter) mkdir(dirPath string) error {
	return os.MkdirAll(dirPath, 0777)
}
