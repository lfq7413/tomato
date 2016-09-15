package tomato

import (
	"github.com/lfq7413/tomato/controllers"
	_ "github.com/lfq7413/tomato/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/lfq7413/tomato/livequery"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/storage"
)

// Run ...
func Run() {

	defer storage.CloseDB()

	// 创建必要的索引
	orm.TomatoDBController.PerformInitialization()

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.ErrorController(&controllers.ErrorController{})

	allowCrossDomain()

	beego.Run()
}

// RunLiveQueryServer 运行 LiveQuery 服务
func RunLiveQueryServer() {
	args := map[string]string{}
	args["pattern"] = "/livequery"
	args["addr"] = "/8089"
	livequery.Run(args)
}

func allowCrossDomain() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Authorization", "Access-Control-Allow-Origin"},
		ExposeHeaders: []string{"Access-Control-Allow-Headers", "X-Parse-Master-Key", "X-Parse-REST-API-Key",
			"X-Parse-Javascript-Key", "X-Parse-Application-Id", "X-Parse-Client-Version", "X-Parse-Session-Token",
			"X-Requested-With", "X-Parse-Revocable-Session", "Content-Type"},
		AllowCredentials: true,
	}))
	beego.InsertFilter("*", beego.BeforeRouter, func(ctx *context.Context) {
		if ctx.Input.Method() == "OPTIONS" {
			ctx.Output.SetStatus(200)
			ctx.ResponseWriter.Started = true
		}
	})
}
