package controllers

import (
	"encoding/json"
	"log"
	"time"

	"github.com/astaxie/beego"
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/utils"
)

// ObjectsController ...
type ObjectsController struct {
	beego.Controller
}

// Post ...
// @router /:className [post]
func (o *ObjectsController) Post() {
	className := o.Ctx.Input.Param(":className")

	var cls map[string]interface{}
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

	data := make(map[string]string)
	data["objectId"] = objectId
	data["createdAt"] = utils.TimetoString(now)

	o.Data["json"] = data
	o.Ctx.Output.SetStatus(201)
	o.Ctx.Output.Header("Location", config.TConfig.URL+"classes/"+className+"/"+objectId)
	o.ServeJSON()
}

// Get ...
// @router /:className/:objectId [get]
func (o *ObjectsController) Get() {
	className := o.Ctx.Input.Param(":className")
	objectId := o.Ctx.Input.Param(":objectId")

	cls := make(map[string]interface{})
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
	data := make(map[string]string)
	data["method"] = "Put"
	data["className"] = className
	data["objectId"] = objectId
	o.Data["json"] = data
	o.ServeJSON()
}

// GetAll ...
// @router /:className [get]
func (o *ObjectsController) GetAll() {
	className := o.Ctx.Input.Param(":className")
	data := make(map[string]string)
	data["method"] = "GetAll"
	data["className"] = className
	o.Data["json"] = data
	o.ServeJSON()
}

// Delete ...
// @router /:className/:objectId [delete]
func (o *ObjectsController) Delete() {
	className := o.Ctx.Input.Param(":className")
	objectId := o.Ctx.Input.Param(":objectId")
	data := make(map[string]string)
	data["method"] = "Delete"
	data["className"] = className
	data["objectId"] = objectId
	o.Data["json"] = data
	o.ServeJSON()
}
