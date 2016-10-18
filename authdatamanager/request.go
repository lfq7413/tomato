package authdatamanager

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/lfq7413/tomato/types"
)

func request(path string, headers map[string]string) (types.M, error) {
	request, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		request.Header.Set(k, v)
	}

	client := http.DefaultClient
	response, err := client.Do(request)
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

func post(path string, headers map[string]string, data map[string]string) (types.M, error) {
	values := url.Values{}
	for k, v := range data {
		values.Set(k, v)
	}
	request, err := http.NewRequest("POST", path, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	for k, v := range headers {
		request.Header.Set(k, v)
	}

	client := http.DefaultClient
	response, err := client.Do(request)
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

func requestQQ(path string, headers map[string]string) (types.M, error) {
	request, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		request.Header.Set(k, v)
	}

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	// 去除 callback( );
	b := string(body)
	b = strings.Replace(b, "\n", "", -1)
	b = strings.Replace(b, "\r", "", -1)
	b = strings.Replace(b, " ", "", -1)
	b = b[9 : len(b)-2]
	body = []byte(b)
	var result types.M
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
