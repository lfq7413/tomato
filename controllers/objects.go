package controllers

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/lfq7413/tomato/auth"
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/utils"
	"gopkg.in/mgo.v2/bson"
)

// ObjectsController ...
type ObjectsController struct {
	beego.Controller
	Info      *RequestInfo
	Auth      *auth.Auth
	ClassName string
	ObjectID  string
}

// RequestInfo ...
type RequestInfo struct {
	AppID          string
	MasterKey      string
	ClientKey      string
	SessionToken   string
	InstallationID string
}

// Prepare ...
func (o *ObjectsController) Prepare() {
	//TODO 1、获取请求头
	info := &RequestInfo{}
	info.AppID = o.Ctx.Input.Header("X-Parse-Application-Id")
	info.MasterKey = o.Ctx.Input.Header("X-Parse-Master-Key")
	info.ClientKey = o.Ctx.Input.Header("X-Parse-Client-Key")
	info.SessionToken = o.Ctx.Input.Header("X-Parse-Session-Token")
	info.InstallationID = o.Ctx.Input.Header("X-Parse-Installation-Id")
	o.Info = info
	//TODO 2、校验头部数据
	if info.AppID != config.TConfig.AppID {
		//TODO AppID 不正确
	}
	if info.MasterKey == config.TConfig.MasterKey {
		o.Auth = &auth.Auth{InstallationID: info.InstallationID, IsMaster: true}
		return
	}
	if info.ClientKey != config.TConfig.ClientKey {
		//TODO ClientKey 不正确
	}
	//TODO 3、生成当前会话用户权限信息
	if info.SessionToken == "" {
		o.Auth = &auth.Auth{InstallationID: info.InstallationID, IsMaster: false}
	} else {
		o.Auth = auth.GetAuthForSessionToken(info.SessionToken, info.InstallationID)
	}

}

// Post ...
// @router /:className [post]
func (o *ObjectsController) Post() {
	className := o.Ctx.Input.Param(":className")

	var cls bson.M
	json.Unmarshal(o.Ctx.Input.RequestBody, &cls)

	objectId := utils.CreateObjectId()
	now := time.Now().UTC()
	cls["_id"] = objectId
	cls["createdAt"] = now
	cls["updatedAt"] = now

	err := orm.TomatoDB.Insert(className, cls)
	if err != nil {
		log.Fatal(err)
	}

	data := bson.M{}
	data["objectId"] = objectId
	data["createdAt"] = utils.TimetoString(now)

	o.Data["json"] = data
	o.Ctx.Output.SetStatus(201)
	o.Ctx.Output.Header("Location", config.TConfig.ServerURL+"/classes/"+className+"/"+objectId)
	o.ServeJSON()
}

// Get ...
// @router /:className/:objectId [get]
func (o *ObjectsController) Get() {

	options := map[string]interface{}{}
	where := map[string]interface{}{"objectId": o.ObjectID}

	rest.Find(o.Auth, o.ClassName, where, options)

	className := o.Ctx.Input.Param(":className")
	objectId := o.Ctx.Input.Param(":objectId")

	cls := bson.M{}
	cls["_id"] = objectId

	data, err := orm.TomatoDB.FindOne(className, cls)
	if err != nil {
		log.Fatal(err)
	}

	data["objectId"] = data["_id"]
	delete(data, "_id")
	if createdAt, ok := data["createdAt"].(time.Time); ok {
		data["createdAt"] = utils.TimetoString(createdAt.UTC())
	}
	if updatedAt, ok := data["updatedAt"].(time.Time); ok {
		data["updatedAt"] = utils.TimetoString(updatedAt.UTC())
	}

	o.Data["json"] = data
	o.ServeJSON()
}

// Put ...
// @router /:className/:objectId [put]
func (o *ObjectsController) Put() {
	className := o.Ctx.Input.Param(":className")
	objectId := o.Ctx.Input.Param(":objectId")

	var cls bson.M
	json.Unmarshal(o.Ctx.Input.RequestBody, &cls)

	now := time.Now().UTC()
	cls["updatedAt"] = now
	update := bson.M{"$set": cls}

	err := orm.TomatoDB.Update(className, bson.M{"_id": objectId}, update)
	if err != nil {
		log.Fatal(err)
	}

	data := bson.M{}
	data["updatedAt"] = utils.TimetoString(now)
	o.Data["json"] = data
	o.ServeJSON()
}

// GetAll ...
// @router /:className [get]
func (o *ObjectsController) GetAll() {

	// TODO 获取查询参数，并组装
	options := map[string]interface{}{}
	if o.GetString("skip") != "" {
		if i, err := strconv.Atoi(o.GetString("skip")); err == nil {
			options["skip"] = i
		} else {
			// TODO return error
		}
	}
	if o.GetString("limit") != "" {
		if i, err := strconv.Atoi(o.GetString("limit")); err == nil {
			options["limit"] = i
		} else {
			// TODO return error
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

	where := map[string]interface{}{}
	if o.GetString("where") != "" {
		err := json.Unmarshal([]byte(o.GetString("where")), &where)
		if err != nil {
			// TODO return err
		}
	}

	rest.Find(o.Auth, o.ClassName, where, options)

	className := o.Ctx.Input.Param(":className")

	cls := bson.M{}

	data, err := orm.TomatoDB.Find(className, cls)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range data {
		v["objectId"] = v["_id"]
		delete(v, "_id")
		if createdAt, ok := v["createdAt"].(time.Time); ok {
			v["createdAt"] = utils.TimetoString(createdAt.UTC())
		}
		if updatedAt, ok := v["updatedAt"].(time.Time); ok {
			v["updatedAt"] = utils.TimetoString(updatedAt.UTC())
		}
	}
	o.Data["json"] = bson.M{"results": data}
	o.ServeJSON()
}

// Delete ...
// @router /:className/:objectId [delete]
func (o *ObjectsController) Delete() {
	className := o.Ctx.Input.Param(":className")
	objectId := o.Ctx.Input.Param(":objectId")

	cls := bson.M{}
	cls["_id"] = objectId

	err := orm.TomatoDB.Remove(className, cls)
	if err != nil {
		log.Fatal(err)
	}

	data := bson.M{}
	o.Data["json"] = data
	o.ServeJSON()
}
