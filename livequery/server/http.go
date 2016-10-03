package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/lfq7413/tomato/livequery/t"
)

// TomatoInfo ...
var TomatoInfo map[string]string

func getUser(sessionToken string) (t.M, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", TomatoInfo["serverURL"]+"/users/me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Parse-Application-Id", TomatoInfo["appId"])
	req.Header.Add("X-Parse-Client-Key", TomatoInfo["clientKey"])
	req.Header.Add("X-Parse-Session-Token", sessionToken)

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
