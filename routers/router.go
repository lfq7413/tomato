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
		beego.NSNamespace("/logout",
			beego.NSInclude(
				&controllers.LogoutController{},
			),
		),
		beego.NSNamespace("/requestPasswordReset",
			beego.NSInclude(
				&controllers.ResetController{},
			),
		),
		beego.NSNamespace("/sessions",
			beego.NSInclude(
				&controllers.SessionsController{},
			),
		),
		beego.NSNamespace("/roles",
			beego.NSInclude(
				&controllers.RolesController{},
			),
		),
		beego.NSNamespace("/files",
			beego.NSInclude(
				&controllers.FilesController{},
			),
		),
		beego.NSNamespace("/events",
			beego.NSInclude(
				&controllers.AnalyticsController{},
			),
		),
		beego.NSNamespace("/push",
			beego.NSInclude(
				&controllers.PushController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
