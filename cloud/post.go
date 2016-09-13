package cloud

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// post 请求网络接口
// 接口返回格式如下：
// {
// 	"success":{},
// 	"error":{},
// }
func post(params types.M, URL string) (r types.M, e types.M) {
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return types.M{}, types.M{"code": -1, "message": "Malformed response"}
	}
	request, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonParams))
	if err != nil {
		return types.M{}, types.M{"code": -1, "message": "Malformed response"}
	}

	request.Header.Set("Content-Type", "application/json")
	if config.TConfig.WebhookKey != "" {
		request.Header.Add("X-Parse-Webhook-Key", config.TConfig.WebhookKey)
	}

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return types.M{}, types.M{"code": -1, "message": "Malformed response"}
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return types.M{}, types.M{"code": -1, "message": "Malformed response"}
	}

	var result types.M
	err = json.Unmarshal(body, &result)
	if err != nil {
		return types.M{}, types.M{"code": -1, "message": "Malformed response"}
	}

	if result["error"] != nil {
		return types.M{}, types.M{"code": 0, "message": utils.S(result["error"])}
	}

	return utils.M(result["success"]), nil
}
