package controllers

import (
	"github.com/astaxie/beego"
    "github.com/lfq7413/tomato/orm"
    "gopkg.in/mgo.v2/bson"
    "log"
)

// ObjectsController ...
type ObjectsController struct {
	beego.Controller
}

// Post ...
// @router /:className [post]
func (o *ObjectsController) Post() {
    className := o.Ctx.Input.Param(":className")
    data := make(map[string]string)
    data["method"] = "Post"
    data["className"] = className
    err := orm.TomatoDB.Insert("std", bson.M{"name": "joe"})
    if err != nil {
        log.Fatal(err)
    }
	o.Data["json"] = data
	o.ServeJSON()
}

// Get ...
// @router /:className/:objectId [get]
func (o *ObjectsController) Get() {
    className := o.Ctx.Input.Param(":className")
    objectId := o.Ctx.Input.Param(":objectId")
    data := make(map[string]string)
    data["method"] = "Get"
    data["className"] = className
    data["objectId"] = objectId
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
