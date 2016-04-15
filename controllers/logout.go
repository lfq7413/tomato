package controllers

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// LogoutController 处理 /logout 接口的请求
type LogoutController struct {
	ObjectsController
}

// HandleLogOut 处理用户退出请求
// @router / [post]
func (l *LogoutController) HandleLogOut() {
	if l.Info != nil && l.Info.SessionToken != "" {
		where := types.M{
			"sessionToken": l.Info.SessionToken,
		}
		records, err := rest.Find(rest.Master(), "_Session", where, types.M{})

		if err != nil {
			l.Data["json"] = errs.ErrorToMap(err)
			l.ServeJSON()
			return
		}
		if utils.HasResults(records) {
			results := utils.SliceInterface(records["results"])
			obj := utils.MapInterface(results[0])
			err := rest.Delete(rest.Master(), "_Session", utils.String(obj["objectId"]))
			if err != nil {
				l.Data["json"] = errs.ErrorToMap(err)
				l.ServeJSON()
				return
			}
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
