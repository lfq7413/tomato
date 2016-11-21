package tencentcos

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
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

func (s sig) getMultiEffectSignature(sk string) string {
	sig := fmt.Sprintf("a=%s&b=%s&k=%s&e=%s&t=%s&r=%s&f=", s.appID, s.bucket, s.secretID, s.expiredTime, s.currentTime, s.rand)
	return getHash(sig, sk)
}

func (s sig) getOnceSignature(sk string) string {
	sig := fmt.Sprintf("a=%s&b=%s&k=%s&e=%s&t=%s&r=%s&f=%s", s.appID, s.bucket, s.secretID, "0", s.currentTime, s.rand, s.fileid)
	return getHash(sig, sk)
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
