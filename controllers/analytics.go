package controllers

import (
	"github.com/lfq7413/tomato/analytics"
	"github.com/lfq7413/tomato/types"
)

// AnalyticsController ...
type AnalyticsController struct {
	ClassesController
}

// HandleEvent ...
// @router /:eventName [post]
func (a *AnalyticsController) HandleEvent() {
	if a.JSONBody == nil {
		a.Data["json"] = types.M{}
		a.ServeJSON()
		return
	}
	response := analytics.TrackEvent(a.Ctx.Input.Param(":eventName"), a.JSONBody)
	a.Data["json"] = response
	a.ServeJSON()
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
