package controllers

import (
	"github.com/astaxie/beego"
)

// ObjectsController ...
type ObjectsController struct {
	beego.Controller
}

// @router / [get]
func (o *ObjectsController) GetAll() {
    data := make(map[string]string)
    data["msg"] = "hello"
	o.Data["json"] = data
	o.ServeJSON()
}