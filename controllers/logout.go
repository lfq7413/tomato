package controllers

import (
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// LogoutController ...
type LogoutController struct {
	ObjectsController
}

// HandleLogOut ...
// @router / [post]
func (l *LogoutController) HandleLogOut() {
	if l.Info != nil && l.Info.SessionToken != "" {
		where := types.M{
			"sessionToken": l.Info.SessionToken,
		}
		// TODO 处理错误
		records, _ := rest.Find(rest.Master(), "_Session", where, types.M{})
		if utils.HasResults(records) {
			results := utils.SliceInterface(records["results"])
			obj := utils.MapInterface(results[0])
			rest.Delete(rest.Master(), "_Session", utils.String(obj["objectId"]))
		}
	}
	l.Data["json"] = types.M{}
	l.ServeJSON()
}

// Get ...
// @router / [get]
func (l *LogoutController) Get() {
	l.ObjectsController.Get()
}

// Delete ...
// @router / [delete]
func (l *LogoutController) Delete() {
	l.ObjectsController.Delete()
}

// Put ...
// @router / [put]
func (l *LogoutController) Put() {
	l.ObjectsController.Put()
}
