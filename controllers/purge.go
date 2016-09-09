package controllers

import (
	"github.com/lfq7413/tomato/cache"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/types"
)

// PurgeController 处理 /purge 接口的请求
type PurgeController struct {
	ClassesController
}

// HandleDelete 处理删除指定类数据请求
// @router /:className [delete]
func (p *PurgeController) HandleDelete() {
	if p.EnforceMasterKeyAccess() == false {
		return
	}
	className := p.Ctx.Input.Param(":className")
	err := orm.TomatoDBController.PurgeCollection(className)
	if err != nil {
		p.HandleError(err, 0)
		return
	}

	if className == "_Session" {
		cache.User.Clear()
	} else if className == "_Role" {
		cache.Role.Clear()
	}

	p.Data["json"] = types.M{}
	p.ServeJSON()
	return
}

// Get ...
// @router / [get]
func (p *PurgeController) Get() {
	p.ClassesController.Get()
}

// Put ...
// @router / [put]
func (p *PurgeController) Put() {
	p.ClassesController.Put()
}

// Post ...
// @router / [post]
func (p *PurgeController) Post() {
	p.ClassesController.Post()
}
