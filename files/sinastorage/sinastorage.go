// Package sinastorage ...
/*
https://github.com/s3storage/sinastoragegosdk
Golang SDK for 新浪云存储
 S3官方API接口文档地址:
 	http://open.sinastorage.com/doc/scs/api
 Contact:
 	s3storage@sina.com
*/
package sinastorage

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

/*
SCS 快捷ACL
 	private 		Bucket和Object 	Owner权限 = FULL_CONTROL，其他人没有任何权限
 	public-read 		Bucket和Object 	Owner权限 = FULL_CONTROL，GRPS000000ANONYMOUSE权限 = READ
 	public-read-write 	Bucket和Object 	Owner权限 = FULL_CONTROL，GRPS000000ANONYMOUSE权限 = READ + WRITE
 	authenticated-read 	Bucket和Object 	Owner权限 = FULL_CONTROL，GRPS0000000CANONICAL权限 = READ
 	GRPS0000000CANONICAL：此组表示所有的新浪云存储注册帐户。所有的请求必须签名（认证），如果签名认证通过，即可按照已设置的权限规则进行访问。
 	GRPS000000ANONYMOUSE：匿名用户组，对应的请求可以不带签名。
 	SINA000000000000IMGX：图片处理服务，将您的bucket的ACL设置为对SINA000000000000IMGX的读写权限，在您使用图片处理服务的时候可以免签名。
SCS 所有方法返回值均是Status Code 和Response Body, Response Body 均是json格式。
*/
type SCS struct {
	Accessk string
	Secretk string
	URI     string
}

type sig struct {
	httpVerb                string
	contentMD5              string
	contentType             string
	date                    string
	canonicalizedAmzHeaders string
	canonicalizedResource   string
}

func (sg sig) getSsig(sk string) string {
	sig := fmt.Sprintf("%s\n%s\n%s\n%s\n%s%s", sg.httpVerb, sg.contentMD5, sg.contentType, sg.date, sg.canonicalizedAmzHeaders, sg.canonicalizedResource)
	return getHash(sig, sk)
}

/*
ListBucket 列出用户账户下所有的bucket。
返回status code 和response body， response body 是json格式。
*/
func (scs SCS) ListBucket() (statusCode int, respBody []uint8) {
	var head = make(map[string]string)

	sig := &sig{"GET", "", "", getDateTime(), "", "/"}
	ssig := sig.getSsig(scs.Secretk)

	uri := fmt.Sprintf("http://%s?formatter=json", scs.URI)
	head["Date"] = getDateTime()
	head["Authorization"] = fmt.Sprintf("SINA %s:%s", scs.Accessk, ssig)

	resp := getResponse(uri, "GET", nil, head)
	body, errIO := ioutil.ReadAll(resp.Body)
	errExcept(errIO)
	defer resp.Body.Close()

	return resp.StatusCode, body
}

/*
ListObject 列出Bucket 下的所有object。
	delimiter 	折叠显示字符,通常使用：'/'
			"" 时，以"join/mailaddresss.txt" 这种”目录+object“的形式展示,
			"/" 时，以"join" 这种"目录"的形式展示，不会展开目录
	prefix 		列出以指定字符为开头的Key,可为""空字符串
	marker 		Key的初始位置，系统将列出比Key大的值，通常用作‘分页’的场景,可为""空字符串
  	max-keys 	返回值的最大Key的数量。
返回status code 和response body，response body 是json格式。
*/
func (scs SCS) ListObject(bucket string, delimiter, prefix, marker string, maxKeys int) (statusCode int, respBody []uint8) {
	var head = make(map[string]string)

	if delimiter != "" {
		delimiter = fmt.Sprintf("delimiter=%s&", delimiter)
	}

	if prefix != "" {
		prefix = fmt.Sprintf("prefix=%s&", prefix)
	}

	if marker != "" {
		marker = fmt.Sprintf("marker=%s&", marker)
	}

	sig := &sig{"GET", "", "", getDateTime(), "", fmt.Sprintf("/%s/", bucket)}
	ssig := sig.getSsig(scs.Secretk)

	uri := fmt.Sprintf("http://%s/%s/?%s%s%s%s&formatter=json", scs.URI, bucket, delimiter, prefix, marker, fmt.Sprintf("max-keys=%d&", maxKeys))
	head["Date"] = getDateTime()
	head["Authorization"] = fmt.Sprintf("SINA %s:%s", scs.Accessk, ssig)

	resp := getResponse(uri, "GET", nil, head)
	body, errIO := ioutil.ReadAll(resp.Body)
	errExcept(errIO)
	defer resp.Body.Close()

	return resp.StatusCode, body
}

/*
GetBucketInfo 获取bucket 的meta 或acl 信息，"info"值为 "meta" or "acl"。
bucket是一个链接项目时，无法获取具体信息。
返回status code 和response body， response body 是json格式。
*/
func (scs SCS) GetBucketInfo(bucket string, info string) (statusCode int, respBody []uint8) {
	var head = make(map[string]string)

	if info == "" {
		info = "meta"
	}

	info = strings.ToLower(info)
	if info != "meta" && info != "acl" {
		log.Fatal("Info Incorrect! InvalidArgument")
	}

	sig := &sig{"GET", "", "", getDateTime(), "", fmt.Sprintf("/%s/?%s", bucket, info)}
	ssig := sig.getSsig(scs.Secretk)

	uri := fmt.Sprintf("http://%s/%s/?%s&formatter=json", scs.URI, bucket, info)
	head["Date"] = getDateTime()
	head["Authorization"] = fmt.Sprintf("SINA %s:%s", scs.Accessk, ssig)

	resp := getResponse(uri, "GET", nil, head)
	body, errIO := ioutil.ReadAll(resp.Body)
	errExcept(errIO)
	defer resp.Body.Close()

	return resp.StatusCode, body
}

/*
PutBucket 创建bucket，acl 是快捷ACL。
acl 值为""，对应的快捷ACL 为private。
成功返回状态码200， 返回response body 为空。
*/
func (scs SCS) PutBucket(bucket string, acl string) (statusCode int, respBody []uint8) {
	var head = make(map[string]string)

	if acl == "" {
		acl = "private"
	}

	if acl != "private" && acl != "public-read" && acl != "public-read-write" && acl != "authenticated-read" {
		log.Fatal("ACL Incorrect! InvalidArgument")
	}

	sig := &sig{"PUT", "", "", getDateTime(), fmt.Sprintf("x-amz-acl:%s\n", acl), fmt.Sprintf("/%s/", bucket)}
	ssig := sig.getSsig(scs.Secretk)

	uri := fmt.Sprintf("http://%s/%s/?formatter=json", scs.URI, bucket)
	head["Date"] = getDateTime()
	head["Authorization"] = fmt.Sprintf("SINA %s:%s", scs.Accessk, ssig)
	head["x-amz-acl"] = acl

	resp := getResponse(uri, "PUT", nil, head)
	body, errIO := ioutil.ReadAll(resp.Body)
	errExcept(errIO)
	defer resp.Body.Close()

	return resp.StatusCode, body
}

/*
DeleteBucket 删除bucket。
删除成功返回状态码204， 返回response body 为空。
*/
func (scs SCS) DeleteBucket(bucket string) (statusCode int, respBody []uint8) {
	var head = make(map[string]string)

	sig := &sig{"DELETE", "", "", getDateTime(), "", fmt.Sprintf("/%s/", bucket)}
	ssig := sig.getSsig(scs.Secretk)

	uri := fmt.Sprintf("http://%s/%s/?formatter=json", scs.URI, bucket)
	head["Date"] = getDateTime()
	head["Authorization"] = fmt.Sprintf("SINA %s:%s", scs.Accessk, ssig)

	resp := getResponse(uri, "DELETE", nil, head)
	body, errIO := ioutil.ReadAll(resp.Body)
	errExcept(errIO)
	defer resp.Body.Close()

	return resp.StatusCode, body
}

/*
SetBucketAcl 设置bucket 的ACL。
	ACL 格式举例： acl := map[string][]string{"GRPS000000ANONYMOUSE": []string{"read", "read_acp", "write", "write_acp"}}
成功返回状态码200， 返回response body 为空。
*/
func (scs SCS) SetBucketAcl(bucket string, acl map[string][]string) (statusCode int, respBody []uint8) {
	var head = make(map[string]string)

	if len(acl) == 0 {
		log.Fatal("ACL Incorrect! InvalidArgument")
	}

	var aclJ, errJ = json.Marshal(acl)
	errExcept(errJ)

	sig := &sig{"PUT", "", "", getDateTime(), "", fmt.Sprintf("/%s/?acl", bucket)}
	ssig := sig.getSsig(scs.Secretk)

	uri := fmt.Sprintf("http://%s/%s/?acl&formatter=json", scs.URI, bucket)
	head["Date"] = getDateTime()
	head["Authorization"] = fmt.Sprintf("SINA %s:%s", scs.Accessk, ssig)

	resp := getResponse(uri, "PUT", strings.NewReader(fmt.Sprintf("%s", aclJ)), head)
	body, errIO := ioutil.ReadAll(resp.Body)
	errExcept(errIO)
	defer resp.Body.Close()

	return resp.StatusCode, body
}

/*
GetObjectInfo 获取object 的meta 或acl 信息，"info"值为 "meta" or "acl"。
返回status code 和response body，response body 是json格式。
*/
func (scs SCS) GetObjectInfo(bucket, object, info string) (statusCode int, respBody []uint8) {
	var head = make(map[string]string)

	if info == "" {
		info = "meta"
	}

	info = strings.ToLower(info)
	if info != "meta" && info != "acl" {
		log.Fatal("Info Incorrect! InvalidArgument")
	}

	sig := &sig{"GET", "", "", getDateTime(), "", fmt.Sprintf("/%s/%s?%s", bucket, object, info)}
	ssig := sig.getSsig(scs.Secretk)

	uri := fmt.Sprintf("http://%s/%s/%s?%s&formatter=json", scs.URI, bucket, object, info)
	head["Date"] = getDateTime()
	head["Authorization"] = fmt.Sprintf("SINA %s:%s", scs.Accessk, ssig)

	resp := getResponse(uri, "GET", nil, head)
	body, errIO := ioutil.ReadAll(resp.Body)
	errExcept(errIO)
	defer resp.Body.Close()

	return resp.StatusCode, body
}

// GetObject 获取object内容。
func (scs SCS) GetObject(bucket, object string) (statusCode int, respBody []uint8) {
	var head = make(map[string]string)

	sig := &sig{"GET", "", "", getDateTime(), "", fmt.Sprintf("/%s/%s", bucket, object)}
	ssig := sig.getSsig(scs.Secretk)

	uri := fmt.Sprintf("http://%s/%s/%s?formatter=json", scs.URI, bucket, object)
	head["Date"] = getDateTime()
	head["Authorization"] = fmt.Sprintf("SINA %s:%s", scs.Accessk, ssig)

	resp := getResponse(uri, "GET", nil, head)
	body, errIO := ioutil.ReadAll(resp.Body)
	errExcept(errIO)
	defer resp.Body.Close()

	return resp.StatusCode, body
}

/*
ObjectCopy 通过拷贝方式创建Object（不上传具体的文件内容。而是通过COPY方式对系统内另一文件进行复制）。
Copy 成功返回状态码200， 返回response body 为空。
*/
func (scs SCS) ObjectCopy(dstbucket, dstobject, srcbucket, srcobject string) (statusCode int, respBody []uint8) {
	var head = make(map[string]string)

	sig := &sig{"PUT", "", "", getDateTime(), fmt.Sprintf("x-amz-copy-source:/%s/%s\n", srcbucket, srcobject), fmt.Sprintf("/%s/%s", dstbucket, dstobject)}
	ssig := sig.getSsig(scs.Secretk)

	uri := fmt.Sprintf("http://%s/%s/%s?formatter=json", scs.URI, dstbucket, dstobject)
	head["Date"] = getDateTime()
	head["Authorization"] = fmt.Sprintf("SINA %s:%s", scs.Accessk, ssig)
	head["x-amz-copy-source"] = fmt.Sprintf("/%s/%s", srcbucket, srcobject)

	resp := getResponse(uri, "PUT", nil, head)
	body, errIO := ioutil.ReadAll(resp.Body)
	errExcept(errIO)
	defer resp.Body.Close()

	return resp.StatusCode, body
}

/*
DeleteObject 删除object。
删除成功返回状态码204， 返回response body 为空。
*/
func (scs SCS) DeleteObject(bucket, object string) (statusCode int, respBody []uint8) {
	var head = make(map[string]string)

	sig := &sig{"DELETE", "", "", getDateTime(), "", fmt.Sprintf("/%s/%s", bucket, object)}
	ssig := sig.getSsig(scs.Secretk)

	uri := fmt.Sprintf("http://%s/%s/%s?formatter=json", scs.URI, bucket, object)
	head["Date"] = getDateTime()
	head["Authorization"] = fmt.Sprintf("SINA %s:%s", scs.Accessk, ssig)

	resp := getResponse(uri, "DELETE", nil, head)
	body, errIO := ioutil.ReadAll(resp.Body)
	errExcept(errIO)
	defer resp.Body.Close()

	return resp.StatusCode, body
}

/*
DeleteObjectData 删除object。
删除成功返回状态码204， 返回response body 为空。
*/
func (scs SCS) DeleteObjectData(bucket, object string) (statusCode int, respBody []uint8, err error) {
	var head = make(map[string]string)

	sig := &sig{"DELETE", "", "", getDateTime(), "", fmt.Sprintf("/%s/%s", bucket, object)}
	ssig := sig.getSsig(scs.Secretk)

	uri := fmt.Sprintf("http://%s/%s/%s?formatter=json", scs.URI, bucket, object)
	head["Date"] = getDateTime()
	head["Authorization"] = fmt.Sprintf("SINA %s:%s", scs.Accessk, ssig)

	resp, err := getResponseAndError(uri, "DELETE", nil, head)
	if err != nil {
		return 500, nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 500, nil, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, body, nil
}

/*
SetObjectMeta 更新一个已经存在的文件的附加meta信息。
 meta 格式举例： meta := map[string]string{"x-amz-meta-name": "sandbox", "x-amz-meta-age": "13"}
注意：这个接口无法更新文件的基本信息，如文件的大小和类型等。
成功返回状态码200， 返回response body 为空。
*/
func (scs SCS) SetObjectMeta(bucket, object string, meta map[string]string) (statusCode int, respBody []uint8) {
	var head = make(map[string]string)
	var headers string
	var keys []string

	for k := range meta {
		keys = append(keys, k)
	}
	sort.Strings(keys) //sort keys
	for _, k := range keys {
		head[k] = meta[k]
		headers = fmt.Sprintf("%s%s:%s\n", headers, k, meta[k])
	}

	sig := &sig{"PUT", "", "", getDateTime(), headers, fmt.Sprintf("/%s/%s?meta", bucket, object)}
	ssig := sig.getSsig(scs.Secretk)

	uri := fmt.Sprintf("http://%s/%s/%s?meta&formatter=json", scs.URI, bucket, object)
	head["Date"] = getDateTime()
	head["Authorization"] = fmt.Sprintf("SINA %s:%s", scs.Accessk, ssig)

	resp := getResponse(uri, "PUT", nil, head)
	body, errIO := ioutil.ReadAll(resp.Body)
	errExcept(errIO)
	defer resp.Body.Close()

	return resp.StatusCode, body
}

/*
SetObjectAcl 设置指定object 的ACL。
 ACL 格式举例： acl := map[string][]string{"GRPS000000ANONYMOUSE": []string{"read", "read_acp", "write", "write_acp"}}
成功返回状态码200， 返回response body 为空。
*/
func (scs SCS) SetObjectAcl(bucket, object string, acl map[string][]string) (statusCode int, respBody []uint8) {
	var head = make(map[string]string)
	var aclJ, errJ = json.Marshal(acl)
	errExcept(errJ)

	sig := &sig{"PUT", "", "", getDateTime(), "", fmt.Sprintf("/%s/%s?acl", bucket, object)}
	ssig := sig.getSsig(scs.Secretk)

	uri := fmt.Sprintf("http://%s/%s/%s?acl&formatter=json", scs.URI, bucket, object)
	head["Date"] = getDateTime()
	head["Authorization"] = fmt.Sprintf("SINA %s:%s", scs.Accessk, ssig)

	resp := getResponse(uri, "PUT", strings.NewReader(fmt.Sprintf("%s", aclJ)), head)
	body, errIO := ioutil.ReadAll(resp.Body)
	errExcept(errIO)
	defer resp.Body.Close()

	return resp.StatusCode, body
}

/*
PutObject 上传object, acl 值为快捷ACL。
成功返回状态码200， 返回response body 为空。
*/
func (scs SCS) PutObject(bucket, object string, uploadfile string, acl string) (statusCode int, respBody []uint8) {
	var head = make(map[string]string)

	if acl == "" {
		acl = "private"
	}

	if acl != "private" && acl != "public-read" && acl != "public-read-write" && acl != "authenticated-read" {
		log.Fatal("ACL Incorrect! InvalidArgument")
	}

	data, errR := ioutil.ReadFile(uploadfile)
	errExcept(errR)
	contentType := http.DetectContentType(data)
	contentMD5 := contentMd5(uploadfile)
	sig := &sig{"PUT", contentMD5, contentType, getDateTime(), fmt.Sprintf("x-amz-acl:%s\nx-amz-meta-uploadlocation:/%s\n", acl, bucket), fmt.Sprintf("/%s/%s", bucket, object)}
	ssig := sig.getSsig(scs.Secretk)

	uri := fmt.Sprintf("http://%s/%s/%s?formatter=json", scs.URI, bucket, object)
	head["Date"] = getDateTime()
	head["Authorization"] = fmt.Sprintf("SINA %s:%s", scs.Accessk, ssig)
	head["Content-Type"] = contentType
	head["Content-MD5"] = contentMD5
	head["x-amz-acl"] = acl
	head["x-amz-meta-uploadlocation"] = fmt.Sprintf("/%s", bucket)

	resp := getResponse(uri, "PUT", strings.NewReader(string(data)), head)
	body, errIO := ioutil.ReadAll(resp.Body)
	errExcept(errIO)
	defer resp.Body.Close()

	return resp.StatusCode, body
}

/*
PutObjectData 上传object, acl 值为快捷ACL。
成功返回状态码200， 返回response body 为空。
*/
func (scs SCS) PutObjectData(bucket, object string, data []byte, acl string) (statusCode int, respBody []uint8, err error) {
	var head = make(map[string]string)

	if acl == "" {
		acl = "private"
	}

	if acl != "private" && acl != "public-read" && acl != "public-read-write" && acl != "authenticated-read" {
		log.Fatal("ACL Incorrect! InvalidArgument")
	}

	contentType := http.DetectContentType(data)
	contentMD5 := conteneMd5Byte(data)
	date := getDateTime()
	sig := &sig{"PUT", contentMD5, contentType, date, fmt.Sprintf("x-amz-acl:%s\nx-amz-meta-uploadlocation:/%s\n", acl, bucket), fmt.Sprintf("/%s/%s", bucket, object)}
	ssig := sig.getSsig(scs.Secretk)

	uri := fmt.Sprintf("http://%s/%s/%s?formatter=json", scs.URI, bucket, object)
	head["Date"] = date
	head["Authorization"] = fmt.Sprintf("SINA %s:%s", scs.Accessk, ssig)
	head["Content-Type"] = contentType
	head["Content-MD5"] = contentMD5
	head["x-amz-acl"] = acl
	head["x-amz-meta-uploadlocation"] = fmt.Sprintf("/%s", bucket)

	resp, err := getResponseAndError(uri, "PUT", bytes.NewReader(data), head)
	if err != nil {
		return 500, nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 500, nil, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, body, nil
}

/*
PutObjectRelax 通过“秒传”方式创建Object（不上传具体的文件内容。而是通过SHA-1值对系统内文件进行复制）。
成功返回状态码200， 返回response body 为空。
*/
func (scs SCS) PutObjectRelax(bucket, object string, uploadfile string) (statusCode int, respBody []uint8) {
	var head = make(map[string]string)

	contentSHA1 := contentSha1(uploadfile)
	contentTYPE := contentType(uploadfile)
	sig := *&sig{"PUT", contentSHA1, contentTYPE, getDateTime(), "", fmt.Sprintf("/%s/%s?relax", bucket, object)}
	ssig := sig.getSsig(scs.Secretk)

	uri := fmt.Sprintf("http://%s/%s/%s?relax&formatter=json", scs.URI, bucket, object)
	head["Date"] = getDateTime()
	head["Authorization"] = fmt.Sprintf("SINA %s:%s", scs.Accessk, ssig)
	head["s-sina-sha1"] = contentSHA1
	head["Content-Type"] = contentTYPE
	head["s-sina-length"] = strconv.Itoa(contentLen(uploadfile))

	resp := getResponse(uri, "PUT", nil, head)
	body, errIO := ioutil.ReadAll(resp.Body)
	errExcept(errIO)
	defer resp.Body.Close()

	return resp.StatusCode, body
}

/*
MultUpload 大文件分片上传
*/
type MultUpload struct {
	SCS
	Bucket     string // 需要上传的bucket
	Object     string // 需要上传的object
	UploadFile string // 上传文件
	SliceCount int    // 分片大小，单位字节
}

// Part ...
type Part struct {
	PartNumber int
	ETag       string
}

type multUploadResponse struct {
	statusCode int
	eTag       string
}

/*
InitiateMultipartUpload 大文件分片上传初始化，返回uploadId。
注意：在初始化上传接口中要求必须进行用户认证，匿名用户无法使用该接口。
在初始化上传时需要给定文件上传所需要的meta绑定信息，在后续的上传中该信息将被保留，并在最终完成时写入云存储系统。
 响应（示例）：
 HTTP/1.1 200 OK
 Date: Tue, 08 Apr 2014 02:59:47 GMT
 Connection: keep-alive
 X-RequestId: 00078d50-1404-0810-5947-782bcb10b128
 X-Requester: Your UserId
 {
    	"Bucket": "<Your-Bucket-Name>",
    	"Key": "<ObjectName>",
    	"UploadId": "7517c1c49a3b4b86a5f08858290c5cf6"
 }
*/
func (mlud MultUpload) InitiateMultipartUpload() (statusCode int, respBody []uint8) {
	var head = make(map[string]string)

	sig := &sig{"POST", "", "", getDateTime(), "", fmt.Sprintf("/%s/%s?multipart", mlud.Bucket, mlud.Object)}
	ssig := sig.getSsig(mlud.Secretk)

	uri := fmt.Sprintf("http://%s/%s/%s?multipart&formatter=json", mlud.URI, mlud.Bucket, mlud.Object)
	head["Date"] = getDateTime()
	head["Authorization"] = fmt.Sprintf("SINA %s:%s", mlud.Accessk, ssig)

	resp := getResponse(uri, "POST", nil, head)
	body, errIO := ioutil.ReadAll(resp.Body)
	errExcept(errIO)
	defer resp.Body.Close()

	return resp.StatusCode, body
}

func (scs SCS) putObjectByte(bucket, object string, data []byte, contentType, acl string, uploadID string, partNumber int, ch chan multUploadResponse) {
	var head = make(map[string]string)

	contentMD5 := conteneMd5Byte(data)
	sig := &sig{"PUT", contentMD5, contentType, getDateTime(), fmt.Sprintf("x-amz-acl:%s\nx-amz-meta-uploadlocation:/%s\n", acl, bucket), fmt.Sprintf("/%s/%s?partNumber=%d&uploadId=%s", bucket, object, partNumber, uploadID)}
	ssig := sig.getSsig(scs.Secretk)

	uri := fmt.Sprintf("http://%s/%s/%s?partNumber=%d&uploadId=%s&formatter=json", scs.URI, bucket, object, partNumber, uploadID)
	head["Date"] = getDateTime()
	head["Authorization"] = fmt.Sprintf("SINA %s:%s", scs.Accessk, ssig)
	head["Content-Type"] = contentType
	head["Content-MD5"] = contentMD5
	head["x-amz-acl"] = acl
	head["x-amz-meta-uploadlocation"] = fmt.Sprintf("/%s", bucket)
	head["Content-Length"] = strconv.Itoa(len(data))

	eTag, errDe := base64.StdEncoding.DecodeString(contentMD5)
	errExcept(errDe)
	resp := getResponse(uri, "PUT", strings.NewReader(string(data)), head)
	//return resp.StatusCode, fmt.Sprintf("%x", eTag)
	ch <- multUploadResponse{resp.StatusCode, fmt.Sprintf("%x", eTag)}
}

/*
UploadPart 上传分片, 注意：分片数不能超过2048。
acl 为快捷ACL。
返回 partInfo []part
 type part struct {
	 PartNumber int		分片id, 从1开始累加
	 ETag       string	分片的md5值
 }
*/
func (mlud MultUpload) UploadPart(uploadID string, acl string) []Part {
	buftmp := make([]byte, mlud.SliceCount)
	var chs []chan multUploadResponse
	var partInfo []Part
	var i = 1

	contentTYPE := contentType(mlud.UploadFile)

	fd, errFl := os.Open(mlud.UploadFile)
	errExcept(errFl)
	defer fd.Close()

	for {
		n, errR := fd.Read(buftmp)
		if errR != nil {
			break
		}
		buftmpVia := make([]byte, n)
		copy(buftmpVia, buftmp)
		buf := buftmpVia

		if n != 0 {
			ch := make(chan multUploadResponse)
			defer close(ch)
			go mlud.putObjectByte(mlud.Bucket, mlud.Object, buf, contentTYPE, acl, uploadID, i, ch)
			chs = append(chs, ch)
			i++
		} else {
			break
		}
	}

	for k, v := range chs {
		timeout := make(chan bool)
		go func() {
			time.Sleep(time.Second * 10)
			timeout <- true
		}()
		select {
		case respInfo := <-v:
			if respInfo.statusCode == 200 {
				partInfo = append(partInfo, Part{k + 1, respInfo.eTag})
			} else {
				log.Fatal("MultUpload Parts Status Code Not Equal 200 !")
			}
		case <-timeout:
			log.Fatal("MultUpload Parts, Read Channel Data Timeout !")
		}

	}

	return partInfo
}

/*
ListParts 列出已经上传的所有分片信息
成功放回状态码200，分片的Parts信息
*/
func (mlud MultUpload) ListParts(uploadID string) (statusCode int, listParts []map[string]interface{}) {
	var head = make(map[string]string)

	sig := &sig{"GET", "", "", getDateTime(), "", fmt.Sprintf("/%s/%s?uploadId=%s", mlud.Bucket, mlud.Object, uploadID)}
	ssig := sig.getSsig(mlud.Secretk)

	uri := fmt.Sprintf("http://%s/%s/%s?uploadId=%s&formatter=json", mlud.URI, mlud.Bucket, mlud.Object, uploadID)
	head["Date"] = getDateTime()
	head["Authorization"] = fmt.Sprintf("SINA %s:%s", mlud.Accessk, ssig)

	resp := getResponse(uri, "GET", nil, head)
	body, errIO := ioutil.ReadAll(resp.Body)
	errExcept(errIO)
	defer resp.Body.Close()

	var listTmp interface{}

	json.Unmarshal(body, &listTmp)
	partsTmp := listTmp.(map[string]interface{})["Parts"].([]interface{})

	var listedParts = make([]map[string]interface{}, len(partsTmp))
	for _, v := range partsTmp {
		listedParts[int(v.(map[string]interface{})["PartNumber"].(float64))-1] = v.(map[string]interface{})
		//fmt.Println(int(v.(map[string]interface{})["PartNumber"].(float64)))
	}

	//fmt.Println(len(ListedParts))
	return resp.StatusCode, listedParts
}

// CompleteMultUpload 大文件分片上传拼接（合并）。
func (mlud MultUpload) CompleteMultUpload(uploadID string, partInfo []Part) (statusCode int, respBody []uint8) {
	var head = make(map[string]string)
	paInfo, errJ := json.Marshal(partInfo)
	errExcept(errJ)

	sig := &sig{"POST", "", "", getDateTime(), "", fmt.Sprintf("/%s/%s?uploadId=%s", mlud.Bucket, mlud.Object, uploadID)}
	ssig := sig.getSsig(mlud.Secretk)

	uri := fmt.Sprintf("http://%s/%s/%s?uploadId=%s&formatter=json", mlud.URI, mlud.Bucket, mlud.Object, uploadID)
	head["Date"] = getDateTime()
	head["Authorization"] = fmt.Sprintf("SINA %s:%s", mlud.Accessk, ssig)

	resp := getResponse(uri, "POST", strings.NewReader(fmt.Sprintf("%s", paInfo)), head)
	body, errIO := ioutil.ReadAll(resp.Body)
	errExcept(errIO)
	defer resp.Body.Close()

	return resp.StatusCode, body

}

// base function
func getDateTime() string {
	timeNow := time.Now()
	dateTime := timeNow.UTC().Format(time.RFC1123Z)
	//dateTimeUnix := timeNow.Unix()
	return dateTime
}

// 过期时间
func getEpDateTime(ep time.Duration) (string, int64) {
	timeNow := time.Now()
	dateEpDT := timeNow.UTC().Add(time.Second * ep).Format(time.RFC1123Z)
	dateEpDTUnix := timeNow.Add(time.Second * ep).Unix()
	return dateEpDT, dateEpDTUnix
}

func getHash(sig, sk string) string {
	mac := hmac.New(sha1.New, []byte(sk))
	mac.Write([]byte(sig))
	ssig := base64.StdEncoding.EncodeToString(mac.Sum(nil))[5:15]
	return ssig
}

func contentMd5(absFile string) string {
	fd, errFile := os.Open(absFile)
	errExcept(errFile)
	defer fd.Close()
	md5Fl := md5.New()
	io.Copy(md5Fl, fd)
	return base64.StdEncoding.EncodeToString(md5Fl.Sum(nil))
}

func conteneMd5Byte(data []byte) string {
	md := md5.New()
	md.Write(data)
	return base64.StdEncoding.EncodeToString(md.Sum(nil))
}

func contentSha1(absFile string) string {
	fd, errFile := os.Open(absFile)
	errExcept(errFile)
	defer fd.Close()
	sha1Fl := sha1.New()
	io.Copy(sha1Fl, fd)
	return fmt.Sprintf("%x", sha1Fl.Sum(nil))

}

func contentType(absFile string) string {
	data, errR := ioutil.ReadFile(absFile)
	errExcept(errR)
	return http.DetectContentType(data)

}

func contentLen(absFile string) int {
	fd, errFile := os.Open(absFile)
	errExcept(errFile)
	defer fd.Close()
	length, _ := ioutil.ReadAll(fd)
	return len(length)
}

func getResponse(uri, method string, body io.Reader, header map[string]string) *http.Response {
	htCli := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*10) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(10 * time.Second)) //设置发送接收数据超时
				return c, nil
			},
		},
	}
	req, erreq := http.NewRequest(method, uri, body)
	errExcept(erreq)

	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, errep := htCli.Do(req)
	errExcept(errep)
	return resp
}

func getResponseAndError(uri, method string, body io.Reader, header map[string]string) (*http.Response, error) {
	htCli := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*10) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(10 * time.Second)) //设置发送接收数据超时
				return c, nil
			},
		},
	}
	req, erreq := http.NewRequest(method, uri, body)
	if erreq != nil {
		return nil, erreq
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}
	return htCli.Do(req)
}

func errExcept(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func urlEncode(uri string) string {
	ur, errPar := url.Parse(uri)
	errExcept(errPar)
	ur.RawQuery = ur.Query().Encode()
	uri = ur.String()
	return uri
}

/*
Display 结构化显示返回的json数据,主要用于response body的格式化显示。 eg:
	{
		"ACL": {
			"GRPS000000ANONYMOUSE": [
				"read"
			],
			"SINA000000RUIKUNTEST": [
				"read",
				"write",
				"read_acp",
				"write_acp"
			]
		},
		"Owner": "SINA000000OWNER"
	}
*/
func Display(dat []byte) {
	if len(dat) > 0 {
		var data interface{}
		errU := json.Unmarshal(dat, &data)
		errExcept(errU)
		data, errJ := json.MarshalIndent(data, "", "\t")
		errExcept(errJ)
		fmt.Printf("%s\n", data)
	}
}

// GetUploadID 用于大文件分片上传，从InitiateMultipartUpload() 返回的response body中拿出uploadId, 并返回。
func GetUploadID(body []byte) string {
	var buf interface{}
	json.Unmarshal(body, &buf)
	uploadID := buf.(map[string]interface{})["UploadId"].(string)
	return uploadID
}
