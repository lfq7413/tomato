package tencentcos

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// COS ...
// 参考链接：https://www.qcloud.com/doc/api/264/3627
type COS struct {
	AppID     string
	SecretID  string
	SecretKey string
	Bucket    string
}

// URI 文件云服务器域名
var URI = "web.file.myqcloud.com"

type sig struct {
	appID       string
	bucket      string
	secretID    string
	expiredTime string
	currentTime string
	rand        string
	fileid      string
}

func (s *sig) getMultiEffectSignature(sk string) string {
	sig := fmt.Sprintf("a=%s&b=%s&k=%s&e=%s&t=%s&r=%s&f=", s.appID, s.bucket, s.secretID, s.expiredTime, s.currentTime, s.rand)
	return getHash(sig, sk)
}

func (s *sig) getOnceSignature(sk string) string {
	sig := fmt.Sprintf("a=%s&b=%s&k=%s&e=%s&t=%s&r=%s&f=%s", s.appID, s.bucket, s.secretID, "0", s.currentTime, s.rand, s.fileid)
	return getHash(sig, sk)
}

// PutObject 创建文件
func (cos *COS) PutObject(object string, data []byte) (statusCode int, err error) {
	s := &sig{
		appID:       cos.AppID,
		bucket:      cos.Bucket,
		secretID:    cos.SecretID,
		expiredTime: getExpiredTime(),
		currentTime: getCurrentTime(),
		rand:        getRand(),
	}
	header := map[string]string{
		"Host":          URI,
		"Authorization": s.getMultiEffectSignature(cos.SecretKey),
	}
	uri := fmt.Sprintf("http://%s/files/v1/%s/%s/%s", URI, cos.AppID, cos.Bucket, object)

	resp, err := uploadFile(uri, data, header)
	if err != nil {
		return 500, err
	}
	return resp.StatusCode, nil
}

// DeleteObject 删除文件
func (cos *COS) DeleteObject(object string) (statusCode int, err error) {
	s := &sig{
		appID:       cos.AppID,
		bucket:      cos.Bucket,
		secretID:    cos.SecretID,
		expiredTime: getExpiredTime(),
		currentTime: getCurrentTime(),
		rand:        getRand(),
		fileid:      fmt.Sprintf("/%s/%s/%s", cos.AppID, cos.Bucket, url.QueryEscape(object)),
	}
	header := map[string]string{
		"Host":          URI,
		"Authorization": s.getOnceSignature(cos.SecretKey),
	}
	uri := fmt.Sprintf("http://%s/files/v1/%s/%s/%s", URI, cos.AppID, cos.Bucket, object)

	resp, err := deleteFile(uri, header)
	if err != nil {
		return 500, err
	}
	return resp.StatusCode, nil
}

// getResponse 用于创建与删除文件
func uploadFile(uri string, body []byte, header map[string]string) (*http.Response, error) {
	fileName := getFileName(uri)
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	fw, err := writer.CreateFormFile("filecontent", fileName)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(fw, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	fw, err = writer.CreateFormField("op")
	if err != nil {
		return nil, err
	}
	_, err = fw.Write([]byte("upload"))
	if err != nil {
		return nil, err
	}

	writer.Close()

	req, err := http.NewRequest("POST", uri, &buf)
	if err != nil {
		return nil, err
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	htCli := &http.Client{}
	return htCli.Do(req)
}

func deleteFile(uri string, header map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("POST", uri, bytes.NewReader([]byte(`{"op":"delete"}`)))
	if err != nil {
		return nil, err
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")

	htCli := &http.Client{}
	return htCli.Do(req)
}

func getHash(sig, sk string) string {
	mac := hmac.New(sha1.New, []byte(sk))
	mac.Write([]byte(sig))
	signTmp := mac.Sum(nil)
	sign := base64.StdEncoding.EncodeToString(append(signTmp, []byte(sig)...))
	return sign
}

func getCurrentTime() string {
	timeNow := time.Now().Unix()
	return fmt.Sprintf("%d", timeNow)
}

func getExpiredTime() string {
	timeNow := time.Now().Unix() + 60
	return fmt.Sprintf("%d", timeNow)
}

func getRand() string {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(9999999999)
	return fmt.Sprintf("%d", r)
}

func getFileName(uri string) string {
	i := strings.LastIndex(uri, "/")
	if i == -1 {
		return uri
	}
	return uri[(i + 1):]
}
