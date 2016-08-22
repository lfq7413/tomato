package cloud

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/lfq7413/tomato/types"
)

func post(params types.M, URL string) types.M {
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return types.M{}
	}
	encodeParams := url.QueryEscape(string(jsonParams))
	request, err := http.NewRequest("POST", URL, strings.NewReader(encodeParams))
	if err != nil {
		return types.M{}
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml,application/json;q=0.9,*/*;q=0.8")
	request.Header.Add("Accept-Language", "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return types.M{}
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return types.M{}
	}

	var result types.M
	err = json.Unmarshal(body, &result)
	if err != nil {
		return types.M{}
	}

	return result
}
