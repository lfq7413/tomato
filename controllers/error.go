package controllers

import (
	"github.com/astaxie/beego"
	"github.com/lfq7413/tomato/types"
)

// ErrorController ...
type ErrorController struct {
	beego.Controller
}

// Error404 ...
func (e *ErrorController) Error404() {
	e.Data["json"] = types.M{"error": "Method Not Allowed"}
	e.ServeJSON()
}

// Error405 ...
func (e *ErrorController) Error405() {
	e.Data["json"] = types.M{"error": "Method Not Allowed"}
	e.ServeJSON()
}

// Error501 ...
func (e *ErrorController) Error501() {
	e.Data["json"] = types.M{"error": "Method Not Allowed"}
	e.ServeJSON()
}
