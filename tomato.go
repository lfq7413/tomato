package tomato

import (
	"github.com/lfq7413/tomato/controllers"
	_ "github.com/lfq7413/tomato/routers"
	"github.com/lfq7413/tomato/types"

	"github.com/astaxie/beego"
	"github.com/lfq7413/tomato/livequery"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/storage"
)

// Run ...
func Run() {

	storage.OpenDB()
	defer storage.CloseDB()

	createIndexes()

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

// createIndexes 创建必要的索引
func createIndexes() {
	requiredUserFields := types.M{}
	defaultUserColumns := types.M{}
	for k, v := range orm.DefaultColumns["_Default"] {
		defaultUserColumns[k] = v
	}
	for k, v := range orm.DefaultColumns["_User"] {
		defaultUserColumns[k] = v
	}
	requiredUserFields["fields"] = defaultUserColumns
	orm.TomatoDBController.LoadSchema(nil).EnforceClassExists("_User")
	orm.Adapter.EnsureUniqueness("_User", requiredUserFields, []string{"username"})
	orm.Adapter.EnsureUniqueness("_User", requiredUserFields, []string{"email"})
}
