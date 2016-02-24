package tomato

import (
	_ "github.com/lfq7413/tomato/routers"

	"github.com/astaxie/beego"
)

//Run ...
func Run() {
    if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}