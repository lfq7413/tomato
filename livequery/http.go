package livequery

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var tomatoInfo map[string]string

func getUser(sessionToken string) (M, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", tomatoInfo["serverURL"]+"/users/me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Parse-Application-Id", tomatoInfo["appId"])
	req.Header.Add("X-Parse-Client-Key", tomatoInfo["clientKey"])
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
	var user M
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
