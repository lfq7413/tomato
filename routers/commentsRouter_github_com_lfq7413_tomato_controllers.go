package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ObjectsController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

}
