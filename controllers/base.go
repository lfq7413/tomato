package controllers

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/astaxie/beego"
	"github.com/lfq7413/tomato/client"
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
)

// BaseController ...
// Info 当前请求的权限信息
// Auth 当前请求的用户权限
// JSONBody 由 JSON 格式转换来的请求数据
// RawBody 原始请求数据
type BaseController struct {
	beego.Controller
	Info     *RequestInfo
	Auth     *rest.Auth
	JSONBody types.M
	RawBody  []byte
}

// RequestInfo http 请求的权限信息
type RequestInfo struct {
	AppID          string
	MasterKey      string
	ClientKey      string
	JavascriptKey  string
	DotNetKey      string
	RestAPIKey     string
	SessionToken   string
	InstallationID string
	ClientVersion  string
	ClientSDK      map[string]string
}

// Prepare 对请求权限进行处理
// 1. 从请求头中获取各种 key
// 2. 尝试按 json 格式转换 body
// 3. 尝试从 body 中获取各种 key
// 4. 校验请求权限
// 5. 生成用户信息
func (b *BaseController) Prepare() {
	info := &RequestInfo{}
	info.AppID = b.Ctx.Input.Header("X-Parse-Application-Id")
	info.MasterKey = b.Ctx.Input.Header("X-Parse-Master-Key")
	info.ClientKey = b.Ctx.Input.Header("X-Parse-Client-Key")
	info.JavascriptKey = b.Ctx.Input.Header("X-Parse-Javascript-Key")
	info.DotNetKey = b.Ctx.Input.Header("X-Parse-Windows-Key")
	info.RestAPIKey = b.Ctx.Input.Header("X-Parse-REST-API-Key")
	info.SessionToken = b.Ctx.Input.Header("X-Parse-Session-Token")
	info.InstallationID = b.Ctx.Input.Header("X-Parse-Installation-Id")
	info.ClientVersion = b.Ctx.Input.Header("X-Parse-Client-Version")

	basicAuth := httpAuth(b.Ctx.Input.Header("authorization"))
	if basicAuth != nil {
		info.AppID = basicAuth["appId"]
		if basicAuth["masterKey"] != "" {
			info.MasterKey = basicAuth["masterKey"]
		}
		if basicAuth["javascriptKey"] != "" {
			info.ClientKey = basicAuth["javascriptKey"]
		}
	}

	if b.Ctx.Input.RequestBody != nil {
		contentType := b.Ctx.Input.Header("Content-type")
		if strings.HasPrefix(contentType, "application/json") {
			// 请求数据为 json 格式，进行转换，转换出错则返回错误
			var object types.M
			err := json.Unmarshal(b.Ctx.Input.RequestBody, &object)
			if err != nil {
				b.Data["json"] = errs.ErrorMessageToMap(errs.InvalidJSON, "invalid JSON")
				b.ServeJSON()
				return
			}
			b.JSONBody = object
		} else {
			// TODO 转换 json 之前，可能需要判断一下数据大小，以确保不会去转换超大数据
			// 其他格式的请求数据，仅尝试转换，转换失败不返回错误
			var object types.M
			err := json.Unmarshal(b.Ctx.Input.RequestBody, &object)
			if err != nil {
				b.RawBody = b.Ctx.Input.RequestBody
			} else {
				b.JSONBody = object
			}
		}
	}

	if b.JSONBody != nil {
		// Unity SDK sends a _noBody key which needs to be removed.
		// Unclear at this point if action needs to be taken.
		delete(b.JSONBody, "_noBody")

		delete(b.JSONBody, "_RevocableSession")
	}

	if info.AppID == "" {
		if b.JSONBody != nil {
			delete(b.JSONBody, "_RevocableSession")
		}
		// 从请求数据中获取各种 key
		if b.JSONBody != nil && b.JSONBody["_ApplicationId"] != nil {
			info.AppID = b.JSONBody["_ApplicationId"].(string)
			delete(b.JSONBody, "_ApplicationId")
			if b.JSONBody["_ClientKey"] != nil {
				info.ClientKey = b.JSONBody["_ClientKey"].(string)
				delete(b.JSONBody, "_ClientKey")
			}
			if b.JSONBody["_InstallationId"] != nil {
				info.InstallationID = b.JSONBody["_InstallationId"].(string)
				delete(b.JSONBody, "_InstallationId")
			}
			if b.JSONBody["_SessionToken"] != nil {
				info.SessionToken = b.JSONBody["_SessionToken"].(string)
				delete(b.JSONBody, "_SessionToken")
			}
			if b.JSONBody["_MasterKey"] != nil {
				info.MasterKey = b.JSONBody["_MasterKey"].(string)
				delete(b.JSONBody, "_MasterKey")
			}
			if b.JSONBody["_ContentType"] != nil {
				b.Ctx.Input.Context.Request.Header.Set("Content-type", b.JSONBody["_ContentType"].(string))
				delete(b.JSONBody, "_ContentType")
			}
		} else {
			// 请求数据中也不存在 APPID 时，返回错误
			b.Data["json"] = errs.ErrorMessageToMap(403, "unauthorized")
			b.Ctx.Output.SetStatus(403)
			b.ServeJSON()
			return
		}
	}

	if info.ClientVersion != "" {
		info.ClientSDK = client.FromString(info.ClientVersion)
	}

	if b.JSONBody != nil && b.JSONBody["base64"] != nil {
		// 请求数据中存在 base64 字段，表明为文件上传，解码并设置到 RawBody 上
		data, err := base64.StdEncoding.DecodeString(b.JSONBody["base64"].(string))
		if err == nil {
			b.RawBody = data
		}
	}

	b.Info = info

	// 校验请求权限
	if info.AppID != config.TConfig.AppID {
		b.Data["json"] = errs.ErrorMessageToMap(403, "unauthorized")
		b.Ctx.Output.SetStatus(403)
		b.ServeJSON()
		return
	}
	if info.MasterKey == config.TConfig.MasterKey {
		b.Auth = &rest.Auth{InstallationID: info.InstallationID, IsMaster: true}
		return
	}
	var allow = false
	if (len(info.ClientKey) > 0 && info.ClientKey == config.TConfig.ClientKey) ||
		(len(info.JavascriptKey) > 0 && info.JavascriptKey == config.TConfig.JavascriptKey) ||
		(len(info.RestAPIKey) > 0 && info.RestAPIKey == config.TConfig.RestAPIKey) ||
		(len(info.DotNetKey) > 0 && info.DotNetKey == config.TConfig.DotNetKey) {
		allow = true
	}
	if allow == false {
		b.Data["json"] = errs.ErrorMessageToMap(403, "unauthorized")
		b.Ctx.Output.SetStatus(403)
		b.ServeJSON()
		return
	}
	// 登录时删除 Token
	url := b.Ctx.Input.URL()
	if strings.HasSuffix(url, "/login/") {
		info.SessionToken = ""
	}
	// 生成当前会话用户权限信息
	if info.SessionToken == "" {
		b.Auth = &rest.Auth{InstallationID: info.InstallationID, IsMaster: false}
	} else {
		var err error
		b.Auth, err = rest.GetAuthForSessionToken(info.SessionToken, info.InstallationID)
		if err != nil {
			b.Data["json"] = errs.ErrorToMap(err)
			b.ServeJSON()
			return
		}
	}
}

func httpAuth(authorization string) map[string]string {
	if authorization == "" {
		return nil
	}

	header := authorization
	var appID, masterKey, javascriptKey string
	authPrefix := "basic "

	match := strings.HasPrefix(strings.ToLower(header), authPrefix)
	if match {
		encodedAuth := header[len(authPrefix):len(header)]
		credentials := strings.Split(decodeBase64(encodedAuth), ":")

		if len(credentials) == 2 {
			appID = credentials[0]
			key := credentials[1]
			jsKeyPrefix := "javascript-key="

			matchKey := strings.HasPrefix(key, jsKeyPrefix)
			if matchKey {
				javascriptKey = key[len(jsKeyPrefix):len(key)]
			} else {
				masterKey = key
			}
			return map[string]string{
				"appId":         appID,
				"masterKey":     masterKey,
				"javascriptKey": javascriptKey,
			}
		}
		return nil
	}

	return nil
}

func decodeBase64(str string) string {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	return string(data)
}
