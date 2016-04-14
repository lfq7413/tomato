package controllers

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// ObjectsController 对象操作 API 的基础结构
// Info 当前请求的权限信息
// Auth 当前请求的用户权限
// JSONBody 由 JSON 格式转换来的请求数据
// RawBody 原始请求数据
// ClassName 要操作的类名
// ObjectID 要操作的对象 id
type ObjectsController struct {
	beego.Controller
	Info      *RequestInfo
	Auth      *rest.Auth
	JSONBody  types.M
	RawBody   []byte
	ClassName string
	ObjectID  string
}

// RequestInfo http 请求的权限信息
type RequestInfo struct {
	AppID          string
	MasterKey      string
	ClientKey      string
	SessionToken   string
	InstallationID string
}

// Prepare 对请求权限进行处理
// 1. 从请求头中获取各种 key
// 2. 尝试按 json 格式转换 body
// 3. 尝试从 body 中获取各种 key
// 4. 校验请求权限
// 5. 生成用户信息
func (o *ObjectsController) Prepare() {
	info := &RequestInfo{}
	info.AppID = o.Ctx.Input.Header("X-Parse-Application-Id")
	info.MasterKey = o.Ctx.Input.Header("X-Parse-Master-Key")
	info.ClientKey = o.Ctx.Input.Header("X-Parse-Client-Key")
	info.SessionToken = o.Ctx.Input.Header("X-Parse-Session-Token")
	info.InstallationID = o.Ctx.Input.Header("X-Parse-Installation-Id")

	if o.Ctx.Input.RequestBody != nil {
		contentType := o.Ctx.Input.Header("Content-type")
		if strings.HasPrefix(contentType, "application/json") {
			// 请求数据为 json 格式，进行转换，转换出错则返回错误
			var object types.M
			err := json.Unmarshal(o.Ctx.Input.RequestBody, &object)
			if err != nil {
				o.Data["json"] = errs.ErrorMessageToMap(errs.InvalidJSON, "invalid JSON")
				o.ServeJSON()
				return
			}
			o.JSONBody = object
		} else {
			// TODO 转换 json 之前，可能需要判断一下数据大小，以确保不会去转换超大数据
			// 其他格式的请求数据，仅尝试转换，转换失败不返回错误
			var object types.M
			err := json.Unmarshal(o.Ctx.Input.RequestBody, &object)
			if err != nil {
				o.RawBody = o.Ctx.Input.RequestBody
			} else {
				o.JSONBody = object
			}
		}
	}

	if o.JSONBody != nil {
		// Unity SDK sends a _noBody key which needs to be removed.
		// Unclear at this point if action needs to be taken.
		delete(o.JSONBody, "_noBody")
	}

	if info.AppID == "" {
		// 从请求数据中获取各种 key
		if o.JSONBody != nil && o.JSONBody["_ApplicationId"] != nil {
			info.AppID = o.JSONBody["_ApplicationId"].(string)
			delete(o.JSONBody, "_ApplicationId")
			if o.JSONBody["_ClientKey"] != nil {
				info.ClientKey = o.JSONBody["_ClientKey"].(string)
				delete(o.JSONBody, "_ClientKey")
			}
			if o.JSONBody["_InstallationId"] != nil {
				info.InstallationID = o.JSONBody["_InstallationId"].(string)
				delete(o.JSONBody, "_InstallationId")
			}
			if o.JSONBody["_SessionToken"] != nil {
				info.SessionToken = o.JSONBody["_SessionToken"].(string)
				delete(o.JSONBody, "_SessionToken")
			}
			if o.JSONBody["_MasterKey"] != nil {
				info.MasterKey = o.JSONBody["_MasterKey"].(string)
				delete(o.JSONBody, "_MasterKey")
			}
		} else {
			// 请求数据中也不存在 APPID 时，返回错误
			o.Data["json"] = errs.ErrorMessageToMap(403, "unauthorized")
			o.Ctx.Output.SetStatus(403)
			o.ServeJSON()
			return
		}
	}

	if o.JSONBody != nil && o.JSONBody["base64"] != nil {
		// 请求数据中存在 base64 字段，表明为文件上传，解码并设置到 RawBody 上
		data, err := base64.StdEncoding.DecodeString(o.JSONBody["base64"].(string))
		if err == nil {
			o.RawBody = data
		}
	}

	o.Info = info

	// 校验请求权限
	if info.AppID != config.TConfig.AppID {
		o.Data["json"] = errs.ErrorMessageToMap(403, "unauthorized")
		o.Ctx.Output.SetStatus(403)
		o.ServeJSON()
		return
	}
	if info.MasterKey == config.TConfig.MasterKey {
		o.Auth = &rest.Auth{InstallationID: info.InstallationID, IsMaster: true}
		return
	}
	if info.ClientKey != config.TConfig.ClientKey {
		o.Data["json"] = errs.ErrorMessageToMap(403, "unauthorized")
		o.Ctx.Output.SetStatus(403)
		o.ServeJSON()
		return
	}
	// 生成当前会话用户权限信息
	if info.SessionToken == "" {
		o.Auth = &rest.Auth{InstallationID: info.InstallationID, IsMaster: false}
	} else {
		o.Auth = rest.GetAuthForSessionToken(info.SessionToken, info.InstallationID)
	}
}

// HandleCreate 处理对象创建请求，返回对象 id 与对象位置
// @router /:className [post]
func (o *ObjectsController) HandleCreate() {

	if o.ClassName == "" {
		o.ClassName = o.Ctx.Input.Param(":className")
	}

	if o.JSONBody == nil {
		o.Data["json"] = errs.ErrorMessageToMap(errs.InvalidJSON, "request body is empty")
		o.ServeJSON()
		return
	}

	result, err := rest.Create(o.Auth, o.ClassName, o.JSONBody)
	if err != nil {
		o.Data["json"] = errs.ErrorToMap(err)
		o.ServeJSON()
		return
	}

	o.Data["json"] = result["response"]
	o.Ctx.Output.SetStatus(201)
	o.Ctx.Output.Header("Location", result["location"].(string))
	o.ServeJSON()

}

// HandleGet 处理查询指定对象请求，返回查询到的对象
// @router /:className/:objectId [get]
func (o *ObjectsController) HandleGet() {
	if o.ClassName == "" {
		o.ClassName = o.Ctx.Input.Param(":className")
	}
	if o.ObjectID == "" {
		o.ObjectID = o.Ctx.Input.Param(":objectId")
	}

	options := types.M{}
	where := types.M{"objectId": o.ObjectID}
	response, err := rest.Find(o.Auth, o.ClassName, where, options)

	if err != nil {
		o.Data["json"] = errs.ErrorToMap(err)
		o.ServeJSON()
		return
	}

	results := utils.SliceInterface(response["results"])
	if results == nil && len(results) == 0 {
		o.Data["json"] = errs.ErrorMessageToMap(errs.ObjectNotFound, "Object not found.")
		o.ServeJSON()
		return
	}

	result := utils.MapInterface(results[0])
	if o.ClassName == "_User" {
		delete(result, "sessionToken")
		if o.Auth.User != nil && result["objectId"].(string) == o.Auth.User["objectId"].(string) {
			// 重新设置 session token
			result["sessionToken"] = o.Info.SessionToken
		}
	}

	o.Data["json"] = result
	o.ServeJSON()

}

// HandleUpdate 处理更新指定对象请求
// @router /:className/:objectId [put]
func (o *ObjectsController) HandleUpdate() {

	if o.ClassName == "" {
		o.ClassName = o.Ctx.Input.Param(":className")
	}
	if o.ObjectID == "" {
		o.ObjectID = o.Ctx.Input.Param(":objectId")
	}

	if o.JSONBody == nil {
		o.Data["json"] = errs.ErrorMessageToMap(errs.InvalidJSON, "request body is empty")
		o.ServeJSON()
		return
	}

	result, err := rest.Update(o.Auth, o.ClassName, o.ObjectID, o.JSONBody)
	if err != nil {
		o.Data["json"] = errs.ErrorToMap(err)
		o.ServeJSON()
		return
	}

	o.Data["json"] = result["response"]
	o.ServeJSON()

}

// HandleFind 处理查找对象请求
// @router /:className [get]
func (o *ObjectsController) HandleFind() {
	if o.ClassName == "" {
		o.ClassName = o.Ctx.Input.Param(":className")
	}

	// 获取查询参数，并组装
	options := types.M{}
	if o.GetString("skip") != "" {
		if i, err := strconv.Atoi(o.GetString("skip")); err == nil {
			options["skip"] = i
		} else {
			o.Data["json"] = errs.ErrorMessageToMap(errs.InvalidQuery, "skip should be int")
			o.ServeJSON()
			return
		}
	}
	if o.GetString("limit") != "" {
		if i, err := strconv.Atoi(o.GetString("limit")); err == nil {
			options["limit"] = i
		} else {
			o.Data["json"] = errs.ErrorMessageToMap(errs.InvalidQuery, "limit should be int")
			o.ServeJSON()
			return
		}
	} else {
		options["limit"] = 100
	}
	if o.GetString("order") != "" {
		options["order"] = o.GetString("order")
	}
	if o.GetString("count") != "" {
		options["count"] = true
	}
	if o.GetString("keys") != "" {
		options["keys"] = o.GetString("keys")
	}
	if o.GetString("include") != "" {
		options["include"] = o.GetString("include")
	}
	if o.GetString("redirectClassNameForKey") != "" {
		options["redirectClassNameForKey"] = o.GetString("redirectClassNameForKey")
	}

	where := types.M{}
	if o.GetString("where") != "" {
		err := json.Unmarshal([]byte(o.GetString("where")), &where)
		if err != nil {
			o.Data["json"] = errs.ErrorMessageToMap(errs.InvalidJSON, "where should be valid json")
			o.ServeJSON()
			return
		}
	}

	response, err := rest.Find(o.Auth, o.ClassName, where, options)
	if err != nil {
		o.Data["json"] = errs.ErrorToMap(err)
		o.ServeJSON()
		return
	}
	if utils.HasResults(response) {
		results := utils.SliceInterface(response["results"])
		for _, v := range results {
			result := utils.MapInterface(v)
			if result["sessionToken"] != nil && o.Info.SessionToken != "" {
				result["sessionToken"] = o.Info.SessionToken
			}
		}
	}

	o.Data["json"] = response
	o.ServeJSON()
}

// HandleDelete ...
// @router /:className/:objectId [delete]
func (o *ObjectsController) HandleDelete() {

	if o.ClassName == "" {
		o.ClassName = o.Ctx.Input.Param(":className")
	}
	if o.ObjectID == "" {
		o.ObjectID = o.Ctx.Input.Param(":objectId")
	}

	rest.Delete(o.Auth, o.ClassName, o.ObjectID)

	className := o.Ctx.Input.Param(":className")
	objectId := o.Ctx.Input.Param(":objectId")

	cls := types.M{}
	cls["_id"] = objectId

	err := orm.TomatoDB.Remove(className, cls)
	if err != nil {
		log.Fatal(err)
	}

	data := types.M{}
	o.Data["json"] = data
	o.ServeJSON()
}

// Get ...
// @router / [get]
func (o *ObjectsController) Get() {
	e := types.M{
		"code":  405,
		"error": "Method Not Allowed",
	}
	o.Ctx.Output.SetStatus(405)
	o.Data["json"] = e
	o.ServeJSON()
}

// Post ...
// @router / [post]
func (o *ObjectsController) Post() {
	e := types.M{
		"code":  405,
		"error": "Method Not Allowed",
	}
	o.Ctx.Output.SetStatus(405)
	o.Data["json"] = e
	o.ServeJSON()
}

// Delete ...
// @router / [delete]
func (o *ObjectsController) Delete() {
	e := types.M{
		"code":  405,
		"error": "Method Not Allowed",
	}
	o.Ctx.Output.SetStatus(405)
	o.Data["json"] = e
	o.ServeJSON()
}

// Put ...
// @router / [put]
func (o *ObjectsController) Put() {
	e := types.M{
		"code":  405,
		"error": "Method Not Allowed",
	}
	o.Ctx.Output.SetStatus(405)
	o.Data["json"] = e
	o.ServeJSON()
}
