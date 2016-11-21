package files

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/files/tencentcos"
)

// tencentAdapter 腾讯云存储
// TODO 测试
type tencentAdapter struct {
	cos *tencentcos.COS
}

func newTencentAdapter() *tencentAdapter {
	cos := &tencentcos.COS{
		AppID:     config.TConfig.TencentAppID,
		SecretID:  config.TConfig.TencentSecretID,
		SecretKey: config.TConfig.TencentSecretKey,
		Bucket:    config.TConfig.TencentBucket,
	}
	t := &tencentAdapter{
		cos: cos,
	}
	return t
}

func (t *tencentAdapter) createFile(filename string, data []byte, contentType string) error {
	code, err := t.cos.PutObject(filename, data)
	if code != 200 || err != nil {
		return errs.E(errs.FileSaveError, "createFile failed.")
	}
	return nil
}

func (t *tencentAdapter) deleteFile(filename string) error {
	code, err := t.cos.DeleteObject(filename)
	if code != 200 || err != nil {
		return errs.E(errs.FileDeleteError, "deleteFile failed.")
	}
	return nil
}

func (t *tencentAdapter) getFileData(filename string) ([]byte, error) {
	return t.download(filename)
}

func (t *tencentAdapter) getFileLocation(filename string) string {
	if config.TConfig.FileDirectAccess {
		return fmt.Sprintf("http://%s-%s.file.myqcloud.com/%s", t.cos.Bucket, t.cos.AppID, url.QueryEscape(filename))
	}
	return config.TConfig.ServerURL + "/files/" + config.TConfig.AppID + "/" + url.QueryEscape(filename)
}

func (t *tencentAdapter) getFileStream(filename string) (FileStream, error) {
	return nil, errs.E(errs.FileReadError, "no such file or directory")
}

func (t *tencentAdapter) getAdapterName() string {
	return "tencentAdapter"
}

func (t *tencentAdapter) download(filename string) ([]byte, error) {
	path := fmt.Sprintf("http://%s-%s.file.myqcloud.com/%s", t.cos.Bucket, t.cos.AppID, url.QueryEscape(filename))
	request, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errs.E(errs.FileReadError, "no such file or directory")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
