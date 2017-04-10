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
	"github.com/lfq7413/tomato/utils"
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
	Query    map[string]string
	JSONBody types.M
	RawBody  []byte
}

// RequestInfo http 请求的权限信息
type RequestInfo struct {
	AppID          string
	MasterKey      string
	ClientKey      string
	JavaScriptKey  string
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
	info.JavaScriptKey = b.Ctx.Input.Header("X-Parse-Javascript-Key")
	info.DotNetKey = b.Ctx.Input.Header("X-Parse-Windows-Key")
	info.RestAPIKey = b.Ctx.Input.Header("X-Parse-REST-API-Key")
	info.SessionToken = b.Ctx.Input.Header("X-Parse-Session-Token")
	info.InstallationID = b.Ctx.Input.Header("X-Parse-Installation-Id")
	info.ClientVersion = b.Ctx.Input.Header("X-Parse-Client-Version")

	basicAuth := httpAuth(b.Ctx.Input.Header("Authorization"))
	if basicAuth != nil {
		info.AppID = basicAuth["appId"]
		if basicAuth["masterKey"] != "" {
			info.MasterKey = basicAuth["masterKey"]
		}
		if basicAuth["javascriptKey"] != "" {
			info.ClientKey = basicAuth["javascriptKey"]
		}
	}

	b.Query = map[string]string{}
	input := b.Input()
	for key := range input {
		b.Query[key] = input.Get(key)
	}

	if b.Ctx.Input.RequestBody != nil {
		contentType := b.Ctx.Input.Header("Content-type")
		if strings.HasPrefix(contentType, "application/json") {
			// 请求数据为 json 格式，进行转换，转换出错则返回错误
			var object types.M
			err := json.Unmarshal(b.Ctx.Input.RequestBody, &object)
			if err != nil {
				b.HandleError(errs.E(errs.InvalidJSON, "invalid JSON"), 0)
				return
			}
			b.JSONBody = object
		} else {
			// 当 AppID 不存在时，尝试转换，转换失败不返回错误
			if info.AppID == "" {
				var object types.M
				err := json.Unmarshal(b.Ctx.Input.RequestBody, &object)
				if err == nil {
					b.JSONBody = object
				}
			}
		}
	}

	if b.JSONBody != nil {
		// Unity SDK sends a _noBody key which needs to be removed.
		// Unclear at this point if action needs to be taken.
		delete(b.JSONBody, "_noBody")
		delete(b.JSONBody, "_method")
	}

	if info.AppID == "" {
		if b.JSONBody != nil {
			delete(b.JSONBody, "_RevocableSession")
		}
		// 从请求数据中获取各种 key
		if b.JSONBody != nil && b.JSONBody["_ApplicationId"] != nil {
			info.AppID = utils.S(b.JSONBody["_ApplicationId"])
			info.JavaScriptKey = utils.S(b.JSONBody["_JavaScriptKey"])
			delete(b.JSONBody, "_ApplicationId")
			delete(b.JSONBody, "_JavaScriptKey")

			if b.JSONBody["_ClientVersion"] != nil {
				info.ClientVersion = utils.S(b.JSONBody["_ClientVersion"])
				delete(b.JSONBody, "_ClientVersion")
			}
			if b.JSONBody["_InstallationId"] != nil {
				info.InstallationID = utils.S(b.JSONBody["_InstallationId"])
				delete(b.JSONBody, "_InstallationId")
			}
			if b.JSONBody["_SessionToken"] != nil {
				info.SessionToken = utils.S(b.JSONBody["_SessionToken"])
				delete(b.JSONBody, "_SessionToken")
			}
			if b.JSONBody["_MasterKey"] != nil {
				info.MasterKey = utils.S(b.JSONBody["_MasterKey"])
				delete(b.JSONBody, "_MasterKey")
			}
			if b.JSONBody["_ContentType"] != nil {
				b.Ctx.Input.Context.Request.Header.Set("Content-type", utils.S(b.JSONBody["_ContentType"]))
				delete(b.JSONBody, "_ContentType")
			}
		} else {
			// 请求数据中也不存在 APPID 时，返回错误
			b.InvalidRequest()
			return
		}
	}

	// 兼容 Android SDK ，查询参数保存在 body 中，并且值的类型均为 string
	if b.Ctx.Input.Method() == "GET" {
		for key, value := range b.JSONBody {
			if str, ok := value.(string); ok {
				b.Query[key] = str
				delete(b.JSONBody, key)
			}
		}
	}

	if info.ClientVersion != "" {
		info.ClientSDK = client.FromString(info.ClientVersion)
	}

	b.Info = info

	// 校验请求权限
	if info.AppID != config.TConfig.AppID {
		b.InvalidRequest()
		return
	}
	if info.MasterKey == config.TConfig.MasterKey {
		b.Auth = &rest.Auth{InstallationID: info.InstallationID, IsMaster: true}
		return
	}
	var allow = false
	if (len(info.ClientKey) > 0 && info.ClientKey == config.TConfig.ClientKey) ||
		(len(info.JavaScriptKey) > 0 && info.JavaScriptKey == config.TConfig.JavaScriptKey) ||
		(len(info.RestAPIKey) > 0 && info.RestAPIKey == config.TConfig.RestAPIKey) ||
		(len(info.DotNetKey) > 0 && info.DotNetKey == config.TConfig.DotNetKey) {
		allow = true
	}
	if allow == false {
		b.InvalidRequest()
		return
	}
	// TODO 登录时删除 Token ，如何处理接口地址？
	url := b.Ctx.Input.URL()
	if url == "/v1/login" || url == "/v1/login/" {
		info.SessionToken = ""
	}
	// 生成当前会话用户权限信息
	if info.SessionToken == "" {
		b.Auth = &rest.Auth{InstallationID: info.InstallationID, IsMaster: false}
		return
	}
	var auth *rest.Auth
	var err error
	if (url == "/v1/upgradeToRevocableSession" || url == "/v1/upgradeToRevocableSession/") &&
		strings.Index(info.SessionToken, "r:") != 0 {
		auth, err = rest.GetAuthForLegacySessionToken(info.SessionToken, info.InstallationID)
	} else {
		auth, err = rest.GetAuthForSessionToken(info.SessionToken, info.InstallationID)
	}
	if err != nil {
		b.HandleError(err, 0)
		return
	}
	b.Auth = auth
}

func httpAuth(authorization string) map[string]string {
	if authorization == "" {
		return nil
	}

	var appID, masterKey, javascriptKey string
	authPrefix1 := "basic "
	authPrefix2 := "Basic "

	match := strings.HasPrefix(authorization, authPrefix1) || strings.HasPrefix(authorization, authPrefix2)
	if match {
		encodedAuth := authorization[len(authPrefix1):]
		credentials := strings.Split(decodeBase64(encodedAuth), ":")

		if len(credentials) == 2 {
			appID = credentials[0]
			key := credentials[1]
			jsKeyPrefix := "javascript-key="

			matchKey := strings.HasPrefix(key, jsKeyPrefix)
			if matchKey {
				javascriptKey = key[len(jsKeyPrefix):]
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

// HandleError 返回错误信息，不指定 status 参数时，默认为 0
func (b *BaseController) HandleError(err error, status int) {
	code := errs.GetErrorCode(err)
	if code != 0 {
		var httpStatus int
		switch code {
		case errs.InternalServerError:
			httpStatus = 500
		case errs.ObjectNotFound:
			httpStatus = 404
		default:
			httpStatus = 400
		}

		b.Ctx.Output.SetStatus(httpStatus)
		b.Data["json"] = errs.ErrorToMap(err)
		b.ServeJSON()
		return
	}

	if status != 0 {
		b.Ctx.Output.SetStatus(status)
		b.Data["json"] = types.M{"error": err.Error()}
		b.ServeJSON()
		return
	}

	b.Ctx.Output.SetStatus(500)
	b.Data["json"] = errs.ErrorMessageToMap(errs.InternalServerError, "Internal server error: "+err.Error())
	b.ServeJSON()
}

// InvalidRequest 无效请求
func (b *BaseController) InvalidRequest() {
	b.Ctx.Output.SetStatus(403)
	b.Data["json"] = types.M{"error": "unauthorized"}
	b.ServeJSON()
}

// EnforceMasterKeyAccess 接口需要 Master 权限
// 返回 true 表示当前请求是 Master 权限
func (b *BaseController) EnforceMasterKeyAccess() bool {
	if b.Auth.IsMaster == false {
		b.Ctx.Output.SetStatus(403)
		b.Data["json"] = types.M{"error": "unauthorized: master key is required"}
		b.ServeJSON()
		return false
	}
	return true
}
