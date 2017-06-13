package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"net/url"

	"github.com/lfq7413/tomato/livequery/t"
)

// TomatoInfo ...
var TomatoInfo = map[string]string{}

// userForSessionToken 访问接口 获取用户信息
func userForSessionToken(sessionToken string) (t.M, error) {
	// TODO 后续使用 go SDK 实现
	where := url.QueryEscape(`{"sessionToken":"` + sessionToken + `"}`)
	req, err := http.NewRequest("GET", TomatoInfo["serverURL"]+"/classes/_Session"+"?where="+where, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Parse-Application-Id", TomatoInfo["appId"])
	req.Header.Add("X-Parse-Master-Key", TomatoInfo["masterKey"])

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
	var session t.M
	err = json.Unmarshal(body, &session)
	if err != nil {
		return nil, errors.New("No session found for session token")
	}

	if user, ok := session["user"].(map[string]interface{}); ok && user != nil {
		return user, nil
	}

	return t.M{}, nil
}

// GetUserRoles 获取用户对应的角色列表
func GetUserRoles(userID string) []string {
	p := url.QueryEscape(`{"users":{"__type":"Pointer","className":"User","objectId":"` + userID + `"}}`)
	req, err := http.NewRequest("GET", TomatoInfo["serverURL"]+"/roles?where="+p, nil)
	if err != nil {
		return []string{}
	}

	req.Header.Add("X-Parse-Application-Id", TomatoInfo["appId"])
	req.Header.Add("X-Parse-Client-Key", TomatoInfo["clientKey"])

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}
	}
	var response t.M
	err = json.Unmarshal(body, &response)
	if err != nil {
		return []string{}
	}
	r := []string{}
	if results, ok := response["results"].([]interface{}); ok {
		for _, result := range results {
			if role, ok := result.(map[string]interface{}); ok {
				if name, ok := role["name"].(string); ok && name != "" {
					r = append(r, "role:"+name)
				}
			}
		}
	}

	return r
}
