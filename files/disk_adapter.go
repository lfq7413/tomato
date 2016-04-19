package files

import (
	"os"

	"github.com/astaxie/beego/utils"
	"github.com/lfq7413/tomato/config"
)

// diskAdapter 本地文件存储模块
type diskAdapter struct {
}

// createFile 在磁盘上创建文件
func (d *diskAdapter) createFile(filename string, data []byte, contentType string) error {
	dir := config.TConfig.FileDir + string(os.PathSeparator) + config.TConfig.AppID
	if utils.FileExists(dir) == false {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}
	}

	filepath := dir + string(os.PathSeparator) + filename
	os.Remove(filepath)

	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// deleteFile 从磁盘删除文件
func (d *diskAdapter) deleteFile(filename string) error {
	dir := config.TConfig.FileDir + string(os.PathSeparator) + config.TConfig.AppID
	filepath := dir + string(os.PathSeparator) + filename
	err := os.Remove(filepath)
	if err != nil {
		return err
	}

	return nil
}

// getFileData 从磁盘获取文件数据，出错时返回空数据
func (d *diskAdapter) getFileData(filename string) []byte {
	dir := config.TConfig.FileDir + string(os.PathSeparator) + config.TConfig.AppID
	filepath := dir + string(os.PathSeparator) + filename

	f, err := os.Open(filepath)
	if err != nil {
		return []byte{}
	}
	defer f.Close()

	data := []byte{}
	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err != nil {
			return []byte{}
		}
		if n == 0 {
			break
		}
		data = append(data, buf[:n]...)
	}

	return data
}

// getFileLocation 获取文件路径
func (d *diskAdapter) getFileLocation(filename string) string {
	return config.TConfig.ServerURL + "/files/" + config.TConfig.AppID + "/" + filename
}
