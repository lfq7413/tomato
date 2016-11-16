package files

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
	"qiniupkg.com/api.v7/conf"
	"qiniupkg.com/api.v7/kodo"
	"qiniupkg.com/api.v7/kodocli"
)

// qiniuAdapter 七牛云存储
// 参考文档： http://developer.qiniu.com/code/v7/sdk/go.html
type qiniuAdapter struct {
	bucket string
	url    string
}

func newQiniuAdapter() *qiniuAdapter {
	q := &qiniuAdapter{
		bucket: config.TConfig.QiniuBucket,
		url:    config.TConfig.QiniuDomain,
	}
	conf.ACCESS_KEY = config.TConfig.QiniuAccessKey
	conf.SECRET_KEY = config.TConfig.QiniuSecretKey
	return q
}

func (q *qiniuAdapter) createFile(filename string, data []byte, contentType string) error {
	c := kodo.New(0, nil)

	policy := &kodo.PutPolicy{
		Scope:   q.bucket,
		Expires: 3600,
	}
	token := c.MakeUptoken(policy)

	zone := 0
	uploader := kodocli.NewUploader(zone, nil)

	var ret kodocli.PutRet
	err := uploader.Put(nil, &ret, token, filename, bytes.NewReader(data), int64(len(data)), nil)
	return err
}

func (q *qiniuAdapter) deleteFile(filename string) error {
	c := kodo.New(0, nil)
	p := c.Bucket(q.bucket)

	err := p.Delete(nil, filename)
	return err
}

func (q *qiniuAdapter) getFileData(filename string) ([]byte, error) {
	return q.download(filename)
}

func (q *qiniuAdapter) getFileLocation(filename string) string {
	if config.TConfig.FileDirectAccess {
		return q.url + "/" + url.QueryEscape(filename)
	}
	return config.TConfig.ServerURL + "/files/" + config.TConfig.AppID + "/" + url.QueryEscape(filename)
}

func (q *qiniuAdapter) getFileStream(filename string) (FileStream, error) {
	return nil, errs.E(errs.FileReadError, "no such file or directory")
}

func (q *qiniuAdapter) getAdapterName() string {
	return "qiniuAdapter"
}

func (q *qiniuAdapter) download(filename string) ([]byte, error) {
	path := q.url + "/" + url.QueryEscape(filename)
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

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
