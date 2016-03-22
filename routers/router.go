package routers

import (
	"github.com/astaxie/beego"
	"github.com/lfq7413/tomato/controllers"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/classes",
			beego.NSInclude(
				&controllers.ObjectsController{},
			),
		),
		beego.NSNamespace("/users",
			beego.NSInclude(
				&controllers.UsersController{},
			),
		),
		beego.NSNamespace("/login",
			beego.NSInclude(
				&controllers.LoginController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
