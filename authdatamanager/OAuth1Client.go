package authdatamanager

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/lfq7413/tomato/types"
)

// OAuth ...
type OAuth struct {
	ConsumerKey     string
	ConsumerSecret  string
	AuthToken       string
	AuthTokenSecret string
	Host            string
	OAuthParams     map[string]string
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
	// TODO
	if len(params) > 0 {
		path = path + "?" + buildParameterString(params)
	}

	request, err := http.NewRequest(method, o.Host+path, nil)
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

	request = signRequest(request, oauthParams, o.ConsumerSecret, o.AuthTokenSecret)

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

func encode(str string) string {
	// TODO
	return str
}

func signRequest(req *http.Request, oauthParameters map[string]string, consumerSecret, authTokenSecret string) *http.Request {
	// TODO
	return req
}
