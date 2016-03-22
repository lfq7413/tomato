package routers

import (
	"github.com/astaxie/beego"
)

func init() {

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

}
