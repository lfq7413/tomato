package files

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

// qiniuAdapter 七牛云存储
// 参考文档： http://developer.qiniu.com/code/v7/sdk/go.html
// 参考文档： https://developer.qiniu.com/kodo/sdk/1238/go
type qiniuAdapter struct {
	bucket    string
	url       string
	accessKey string
	secretKey string
}

func newQiniuAdapter() *qiniuAdapter {
	q := &qiniuAdapter{
		bucket:    config.TConfig.QiniuBucket,
		url:       config.TConfig.QiniuDomain,
		accessKey: config.TConfig.QiniuAccessKey,
		secretKey: config.TConfig.QiniuSecretKey,
	}
	return q
}

func (q *qiniuAdapter) createFile(filename string, data []byte, contentType string) error {
	putPolicy := storage.PutPolicy{
		Scope: q.bucket,
	}
	mac := qbox.NewMac(q.accessKey, q.secretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	// 空间对应的机房
	switch config.TConfig.QiniuZone {
	case "Huadong":
		cfg.Zone = &storage.ZoneHuadong
	case "Huabei":
		cfg.Zone = &storage.ZoneHuabei
	case "Huanan":
		cfg.Zone = &storage.ZoneHuanan
	case "Beimei":
		cfg.Zone = &storage.ZoneBeimei
	}
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}

	err := formUploader.Put(context.Background(), &ret, upToken, filename, bytes.NewReader(data), int64(len(data)), &putExtra)
	return err
}

func (q *qiniuAdapter) deleteFile(filename string) error {
	mac := qbox.NewMac(q.accessKey, q.secretKey)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	//cfg.Zone=&storage.ZoneHuabei
	bucketManager := storage.NewBucketManager(mac, &cfg)
	err := bucketManager.Delete(q.bucket, filename)
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

	if response.StatusCode != 200 {
		return nil, errs.E(errs.FileReadError, "no such file or directory")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
