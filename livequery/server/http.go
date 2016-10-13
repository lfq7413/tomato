package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/lfq7413/tomato/livequery/t"
)

// TomatoInfo ...
var TomatoInfo = map[string]string{}

// getUser 访问接口 获取用户信息
func getUser(sessionToken string) (t.M, error) {
	req, err := http.NewRequest("GET", TomatoInfo["serverURL"]+"/users/me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Parse-Application-Id", TomatoInfo["appId"])
	req.Header.Add("X-Parse-Client-Key", TomatoInfo["clientKey"])
	req.Header.Add("X-Parse-Session-Token", "r:"+sessionToken)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var user t.M
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
