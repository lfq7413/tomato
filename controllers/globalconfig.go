package controllers

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// GlobalConfigController 处理 /config 接口的请求
type GlobalConfigController struct {
	ClassesController
}

// HandleGet 获取配置信息
// @router / [get]
func (g *GlobalConfigController) HandleGet() {
	if g.Auth.IsMaster == false {
		g.Data["json"] = errs.ErrorMessageToMap(errs.OperationForbidden, "Need master key!")
		g.ServeJSON()
		return
	}

	results, _ := orm.TomatoDBController.Find("_GlobalConfig", types.M{"objectId": "1"}, types.M{"limit": 1})
	if len(results) != 1 {
		g.Data["json"] = types.M{"params": types.M{}}
		g.ServeJSON()
		return
	}
	globalConfig := utils.M(results[0])
	if globalConfig == nil {
		g.Data["json"] = types.M{"params": types.M{}}
		g.ServeJSON()
		return
	}
	g.Data["json"] = types.M{"params": globalConfig["params"]}
	g.ServeJSON()
}

// HandlePut 修改配置信息
// @router / [put]
func (g *GlobalConfigController) HandlePut() {
	if g.Auth.IsMaster == false {
		g.Data["json"] = errs.ErrorMessageToMap(errs.OperationForbidden, "Need master key!")
		g.ServeJSON()
		return
	}

	if g.JSONBody == nil || utils.M(g.JSONBody["params"]) == nil {
		g.Data["json"] = types.M{"result": true}
		g.ServeJSON()
		return
	}
	params := utils.M(g.JSONBody["params"])
	update := types.M{}
	for k, v := range params {
		update["params."+k] = v
	}
	_, err := orm.TomatoDBController.Update("_GlobalConfig", types.M{"objectId": "1"}, update, types.M{"upsert": true}, false)
	if err != nil {
		g.HandleError(err, 0)
		return
	}
	g.Data["json"] = types.M{"result": true}
	g.ServeJSON()
}

// Post ...
// @router / [post]
func (g *GlobalConfigController) Post() {
	g.ClassesController.Post()
}

// Delete ...
// @router / [delete]
func (g *GlobalConfigController) Delete() {
	g.ClassesController.Delete()
}
