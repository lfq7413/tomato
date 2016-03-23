package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LoginController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LoginController"],
		beego.ControllerComments{
			"HandleLogIn",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LoginController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LoginController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LoginController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LoginController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LoginController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LoginController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogoutController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogoutController"],
		beego.ControllerComments{
			"HandleLogOut",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogoutController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogoutController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogoutController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogoutController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogoutController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogoutController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"],
		beego.ControllerComments{
			"HandleCreate",
			`/:className`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"],
		beego.ControllerComments{
			"HandleGet",
			`/:className/:objectId`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"],
		beego.ControllerComments{
			"HandleUpdate",
			`/:className/:objectId`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"],
		beego.ControllerComments{
			"HandleFind",
			`/:className`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"],
		beego.ControllerComments{
			"HandleDelete",
			`/:className/:objectId`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ResetController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ResetController"],
		beego.ControllerComments{
			"HandleResetRequest",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ResetController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ResetController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ResetController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ResetController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ResetController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ResetController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"],
		beego.ControllerComments{
			"HandleFind",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"],
		beego.ControllerComments{
			"HandleGet",
			`/:objectId`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"],
		beego.ControllerComments{
			"HandleCreate",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"],
		beego.ControllerComments{
			"HandleUpdate",
			`/:objectId`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"],
		beego.ControllerComments{
			"HandleDelete",
			`/:objectId`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"],
		beego.ControllerComments{
			"HandleMe",
			`/me`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

}
