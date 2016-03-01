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
	)
	beego.AddNamespace(ns)
}
