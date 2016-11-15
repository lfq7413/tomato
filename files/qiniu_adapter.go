package files

import (
	"bytes"

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
	return nil, errs.E(errs.FileReadError, "no such file or directory")
}

func (q *qiniuAdapter) getFileLocation(filename string) string {
	return q.url + "/" + filename
}

func (q *qiniuAdapter) getFileStream(filename string) (FileStream, error) {
	return nil, errs.E(errs.FileReadError, "no such file or directory")
}

func (q *qiniuAdapter) getAdapterName() string {
	return "qiniuAdapter"
}
