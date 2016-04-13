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
func (c *ErrorController) Error404() {
	e := types.M{
		"code":  404,
		"error": "Method Not Allowed",
	}
	c.Data["json"] = e
	c.ServeJSON()
}

// Error405 ...
func (c *ErrorController) Error405() {
	e := types.M{
		"code":  405,
		"error": "Method Not Allowed",
	}
	c.Data["json"] = e
	c.ServeJSON()
}

// Error501 ...
func (c *ErrorController) Error501() {
	e := types.M{
		"code":  501,
		"error": "server error",
	}
	c.Data["json"] = e
	c.ServeJSON()
}
