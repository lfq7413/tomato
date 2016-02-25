package routers

import (
	"github.com/lfq7413/tomato/controllers"
	"github.com/astaxie/beego"
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
