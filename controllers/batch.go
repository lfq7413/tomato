package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// BatchController ...
type BatchController struct {
	ClassesController
}

// HandleBatch ...
// @router / [post]
func (b *BatchController) HandleBatch() {
	if b.JSONBody == nil {
		b.HandleError(errs.E(errs.InvalidJSON, "requests must be an array"), 0)
		return
	}
	requests := utils.A(b.JSONBody["requests"])
	if requests == nil {
		b.HandleError(errs.E(errs.InvalidJSON, "requests must be an array"), 0)
		return
	}

	headers := map[string]string{
		"X-Parse-Application-Id": b.Info.AppID,
	}
	if b.Info.MasterKey != "" {
		headers["X-Parse-Master-Key"] = b.Info.MasterKey
	}
	if b.Info.ClientKey != "" {
		headers["X-Parse-Client-Key"] = b.Info.ClientKey
	}
	if b.Info.JavaScriptKey != "" {
		headers["X-Parse-Javascript-Key"] = b.Info.JavaScriptKey
	}
	if b.Info.DotNetKey != "" {
		headers["X-Parse-Windows-Key"] = b.Info.DotNetKey
	}
	if b.Info.RestAPIKey != "" {
		headers["X-Parse-REST-API-Key"] = b.Info.RestAPIKey
	}
	if b.Info.SessionToken != "" {
		headers["X-Parse-Session-Token"] = b.Info.SessionToken
	}
	if b.Info.InstallationID != "" {
		headers["X-Parse-Installation-Id"] = b.Info.InstallationID
	}
	if b.Info.ClientVersion != "" {
		headers["X-Parse-Client-Version"] = b.Info.ClientVersion
	}
	if b.Ctx.Input.Header("Authorization") != "" {
		headers["Authorization"] = b.Ctx.Input.Header("Authorization")
	}

	b.HandleRequest(requests, headers, b.Ctx.Input.Scheme())
}

// HandleRequest ...
func (b *BatchController) HandleRequest(requests types.S, headers map[string]string, scheme string) {
	methods := []string{}
	paths := []string{}
	bodys := []interface{}{}
	results := types.S{}

	for _, v := range requests {
		request := utils.M(v)
		if request == nil {
			b.HandleError(errs.E(errs.InvalidJSON, "Invalid request"), 0)
			return
		}

		method := utils.S(request["method"])
		if method == "" {
			b.HandleError(errs.E(errs.InvalidJSON, "Invalid method"), 0)
			return
		}
		methods = append(methods, method)

		path := utils.S(request["path"])
		if path == "" {
			b.HandleError(errs.E(errs.InvalidJSON, "Invalid path"), 0)
			return
		}
		if strings.HasPrefix(path, "/") {
			path = scheme + "://127.0.0.1:" + beego.AppConfig.String("httpport") + path
		} else if strings.HasPrefix(path, scheme+"://") {
			path = path[len(scheme+"://"):]
			p := strings.Index(path, "/")
			if p == -1 {
				b.HandleError(errs.E(errs.InvalidJSON, "Invalid path"), 0)
				return
			}
			path = path[p:]
			path = scheme + "://127.0.0.1:" + beego.AppConfig.String("httpport") + path
		} else {
			b.HandleError(errs.E(errs.InvalidJSON, "Invalid path"), 0)
			return
		}
		paths = append(paths, path)

		bodys = append(bodys, request["body"])
	}
	for i := 0; i < len(requests); i++ {
		r := request(methods[i], paths[i], headers, bodys[i])
		results = append(results, r)
	}
	b.Data["json"] = results
	b.ServeJSON()
}

func request(method, path string, headers map[string]string, body interface{}) types.M {
	var requestBody io.Reader
	if body == nil {
		requestBody = nil
	} else {
		jsonParams, err := json.Marshal(body)
		if err != nil {
			return types.M{"error": errs.ErrorMessageToMap(errs.InvalidJSON, "Invalid body")}
		}
		requestBody = bytes.NewBuffer(jsonParams)
	}

	request, err := http.NewRequest(method, path, requestBody)
	if err != nil {
		return types.M{"error": errs.ErrorMessageToMap(errs.InvalidJSON, "Invalid request")}
	}

	if method == "POST" || method == "PUT" || method == "PATCH" {
		request.Header.Set("Content-Type", "application/json")
	}
	for header, value := range headers {
		request.Header.Set(header, value)
	}

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return types.M{"error": errs.ErrorMessageToMap(errs.InvalidJSON, "Invalid request")}
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return types.M{"error": errs.ErrorMessageToMap(errs.InvalidJSON, "Invalid request")}
	}

	var result types.M
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return types.M{"error": types.M{"error": string(responseBody)}}
	}

	if result["error"] != nil {
		return types.M{"error": result}
	}

	return types.M{"success": result}
}

// Get ...
// @router / [get]
func (b *BatchController) Get() {
	b.HandleError(errors.New("Method Not Allowed"), 405)
}

// Delete ...
// @router / [delete]
func (b *BatchController) Delete() {
	b.HandleError(errors.New("Method Not Allowed"), 405)
}

// Put ...
// @router / [put]
func (b *BatchController) Put() {
	b.HandleError(errors.New("Method Not Allowed"), 405)
}
