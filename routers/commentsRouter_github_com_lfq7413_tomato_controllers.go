package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"],
		beego.ControllerComments{
			"Post",
			`/:className`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"],
		beego.ControllerComments{
			"Get",
			`/:className/:objectId`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"],
		beego.ControllerComments{
			"Put",
			`/:className/:objectId`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"],
		beego.ControllerComments{
			"GetAll",
			`/:className`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"],
		beego.ControllerComments{
			"Delete",
			`/:className/:objectId`,
			[]string{"delete"},
			nil})

}
