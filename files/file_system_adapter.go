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
func (f *fileSystemAdapter) getFileData(filename string) ([]byte, error) {
	filepath := f.getLocalFilePath(filename)

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
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

	return data, nil
}

// getFileLocation 获取文件路径
func (f *fileSystemAdapter) getFileLocation(filename string) string {
	return config.TConfig.ServerURL + "/files/" + config.TConfig.AppID + "/" + url.QueryEscape(filename)
}

func (f *fileSystemAdapter) getFileStream(filename string) (FileStream, error) {
	filepath := f.getLocalFilePath(filename)
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	return &diskFileStream{file: file}, nil
}

func (f *fileSystemAdapter) getAdapterName() string {
	return "fileSystemAdapter"
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

type diskFileStream struct {
	file *os.File
}

func (d *diskFileStream) Seek(offset int64, whence int) (ret int64, err error) {
	return d.file.Seek(offset, whence)
}

func (d *diskFileStream) Read(b []byte) (n int, err error) {
	return d.file.Read(b)
}

func (d *diskFileStream) Size() (bytes int64) {
	i, err := d.file.Stat()
	if err != nil {
		return 0
	}
	return i.Size()
}

func (d *diskFileStream) Close() (err error) {
	return d.file.Close()
}
