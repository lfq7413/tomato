package tomato

import (
	"github.com/lfq7413/tomato/controllers"
	_ "github.com/lfq7413/tomato/routers"

	"github.com/astaxie/beego"
	"github.com/lfq7413/tomato/livequery"
	"github.com/lfq7413/tomato/storage"
)

// Run ...
func Run() {

	storage.OpenDB()
	defer storage.CloseDB()

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.ErrorController(&controllers.ErrorController{})
	beego.Run()
}

// RunLiveQueryServer 运行 LiveQuery 服务
func RunLiveQueryServer() {
	args := map[string]string{}
	args["pattern"] = "/livequery"
	args["addr"] = "/8089"
	livequery.Run(args)
}
