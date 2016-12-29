package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

const signatureMethod = "HMAC-SHA1"
const version = "1.0"

// OAuth ...
type OAuth struct {
	ConsumerKey     string
	ConsumerSecret  string
	AuthToken       string
	AuthTokenSecret string
	Host            string
	OAuthParams     map[string]string
}

// NewOAuth ...
func NewOAuth(options types.M) *OAuth {
	o := &OAuth{
		ConsumerKey:     utils.S(options["consumer_key"]),
		ConsumerSecret:  utils.S(options["consumer_secret"]),
		AuthToken:       utils.S(options["auth_token"]),
		AuthTokenSecret: utils.S(options["auth_token_secret"]),
		Host:            utils.S(options["host"]),
	}
	return o
}

// Get ...
func (o *OAuth) Get(path string, params map[string]string) (types.M, error) {
	req, err := o.buildRequest("GET", path, params, nil)
	if err != nil {
		return nil, err
	}
	return o.Send(req)
}

// Post ...
func (o *OAuth) Post(path string, params map[string]string, body map[string]string) (types.M, error) {
	req, err := o.buildRequest("POST", path, params, body)
	if err != nil {
		return nil, err
	}
	return o.Send(req)
}

// Send ...
func (o *OAuth) Send(req *http.Request) (types.M, error) {
	client := http.DefaultClient
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var result types.M
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (o *OAuth) buildRequest(method, path string, params, body map[string]string) (*http.Request, error) {
	if len(params) > 0 {
		path = path + "?" + buildParameterString(params)
	}

	var request *http.Request
	var err error
	if len(body) > 0 {
		request, err = http.NewRequest(method, o.Host+path, strings.NewReader(buildParameterString(body)))
	} else {
		request, err = http.NewRequest(method, o.Host+path, nil)
	}
	if err != nil {
		return nil, err
	}

	oauthParams := map[string]string{}
	if o.OAuthParams != nil {
		oauthParams = o.OAuthParams
	}
	oauthParams["oauth_consumer_key"] = o.ConsumerKey
	if o.AuthToken != "" {
		oauthParams["oauth_token"] = o.AuthToken
	}

	request = signRequest(request, oauthParams, o.ConsumerSecret, o.AuthTokenSecret, request.URL.RequestURI(), params, body)

	return request, nil
}

func buildParameterString(obj map[string]string) string {
	if len(obj) == 0 {
		return ""
	}
	keys := []string{}
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Sort(sort.StringSlice(keys))
	result := []string{}
	for _, k := range keys {
		result = append(result, k+"="+encode(obj[k]))
	}
	return strings.Join(result, "&")
}

func signRequest(req *http.Request, oauthParameters map[string]string, consumerSecret, authTokenSecret, url string, params, body map[string]string) *http.Request {
	if oauthParameters == nil {
		oauthParameters = map[string]string{}
	}
	if oauthParameters["oauth_nonce"] == "" {
		oauthParameters["oauth_nonce"] = nonce()
	}
	if oauthParameters["oauth_timestamp"] == "" {
		oauthParameters["oauth_timestamp"] = strconv.Itoa(int(time.Now().Unix()))
	}
	if oauthParameters["oauth_signature_method"] == "" {
		oauthParameters["oauth_signature_method"] = signatureMethod
	}
	if oauthParameters["oauth_version"] == "" {
		oauthParameters["oauth_version"] = version
	}

	signatureParams := map[string]string{}
	for _, parameters := range []map[string]string{oauthParameters, params, body} {
		for k, v := range parameters {
			signatureParams[k] = v
		}
	}

	parameterString := buildParameterString(signatureParams)
	signatureString := buildSignatureString(req.Method, url, parameterString)
	signatureKey := strings.Join([]string{encode(consumerSecret), encode(authTokenSecret)}, "&")
	signature := signature(signatureString, signatureKey)

	oauthParameters["oauth_signature"] = signature

	keys := []string{}
	for k := range oauthParameters {
		keys = append(keys, k)
	}
	sort.Sort(sort.StringSlice(keys))
	for i, k := range keys {
		keys[i] = k + `="` + oauthParameters[k] + `"`
	}
	authHeader := strings.Join(keys, ", ")

	req.Header.Set("Authorization", "OAuth "+authHeader)
	if req.Method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return req
}

func buildSignatureString(method, url, parameters string) string {
	l := []string{method, encode(url), parameters}
	return strings.Join(l, "&")
}

func signature(text, key string) string {
	k := []byte(key)
	mac := hmac.New(sha1.New, k)
	mac.Write([]byte(text))
	src := mac.Sum(nil)
	result := base64.StdEncoding.EncodeToString(src)
	return encode(result)
}

func encode(str string) string {
	str = url.QueryEscape(str)
	str = strings.Replace(str, `!`, "%21", -1)
	str = strings.Replace(str, `'`, "%27", -1)
	str = strings.Replace(str, `(`, "%28", -1)
	str = strings.Replace(str, `)`, "%29", -1)
	str = strings.Replace(str, `*`, "%2A", -1)
	return str
}

func nonce() string {
	return utils.CreateString(30)
}
