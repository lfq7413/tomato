package controllers

import (
	"github.com/astaxie/beego"
)

// ErrorController ...
type ErrorController struct {
	beego.Controller
}

// Error404 ...
func (c *ErrorController) Error404() {
	e := map[string]interface{}{
		"code":  404,
		"error": "Method Not Allowed",
	}
	c.Data["json"] = e
	c.ServeJSON()
}

// Error405 ...
func (c *ErrorController) Error405() {
	e := map[string]interface{}{
		"code":  405,
		"error": "Method Not Allowed",
	}
	c.Data["json"] = e
	c.ServeJSON()
}

// Error501 ...
func (c *ErrorController) Error501() {
	e := map[string]interface{}{
		"code":  501,
		"error": "server error",
	}
	c.Data["json"] = e
	c.ServeJSON()
}
