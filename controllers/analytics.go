package controllers

import (
	"github.com/lfq7413/tomato/analytics"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// AnalyticsController ...
type AnalyticsController struct {
	ClassesController
}

// AppOpened ...
// @router /AppOpened [post]
func (a *AnalyticsController) AppOpened() {
	if a.JSONBody == nil {
		a.Data["json"] = types.M{}
		a.ServeJSON()
		return
	}
	a.addTags(a.JSONBody)
	response := analytics.AppOpened(a.JSONBody)
	a.Data["json"] = response
	a.ServeJSON()
}

// HandleEvent ...
// @router /:eventName [post]
func (a *AnalyticsController) HandleEvent() {
	if a.JSONBody == nil {
		a.Data["json"] = types.M{}
		a.ServeJSON()
		return
	}
	a.addTags(a.JSONBody)
	response := analytics.TrackEvent(a.Ctx.Input.Param(":eventName"), a.JSONBody)
	a.Data["json"] = response
	a.ServeJSON()
}

func (a *AnalyticsController) addTags(event types.M) {
	tags := utils.M(event["tags"])
	if tags == nil {
		tags = types.M{}
	}

	if a.Info.ClientKey != "" {
		tags["from"] = "Client"
	} else if a.Info.JavaScriptKey != "" {
		tags["from"] = "JavaScript"
	} else if a.Info.DotNetKey != "" {
		tags["from"] = "DotNet"
	} else if a.Info.RestAPIKey != "" {
		tags["from"] = "RestAPI"
	}

	if a.Info.ClientVersion != "" {
		tags["version"] = a.Info.ClientVersion
	}

	event["tags"] = tags
}

// Get ...
// @router / [get]
func (a *AnalyticsController) Get() {
	a.ClassesController.Get()
}

// Post ...
// @router / [post]
func (a *AnalyticsController) Post() {
	a.ClassesController.Post()
}

// Delete ...
// @router / [delete]
func (a *AnalyticsController) Delete() {
	a.ClassesController.Delete()
}

// Put ...
// @router / [put]
func (a *AnalyticsController) Put() {
	a.ClassesController.Put()
}
