package authdatamanager

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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
	return nil, nil
}
