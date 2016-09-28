package controllers

import (
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// UpgradeSessionController 处理 /upgradeToRevocableSession 接口的请求
type UpgradeSessionController struct {
	ClassesController
}

// HandleUpdateToRevocableSession ...
// @router / [post]
func (u *UpgradeSessionController) HandleUpdateToRevocableSession() {
	if u.Info == nil || u.Info.SessionToken == "" {
		u.HandleError(errs.E(errs.InvalidSessionToken, "Session token required."), 0)
		return
	}

	token := "r:" + utils.CreateToken()
	userID := utils.S(u.Auth.User["objectId"])
	expiresAt := config.GenerateSessionExpiresAt()
	sessionData := types.M{
		"sessionToken": token,
		"user": types.M{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  userID,
		},
		"createdWith": types.M{
			"action": "upgrade",
		},
		"restricted":     false,
		"installationId": u.Auth.InstallationID,
		"expiresAt": types.M{
			"__type": "Date",
			"iso":    utils.TimetoString(expiresAt),
		},
	}

	create, err := rest.NewWrite(rest.Master(), "_Session", nil, sessionData, nil, nil)
	if err != nil {
		u.HandleError(err, 0)
		return
	}
	_, err = create.Execute()
	if err != nil {
		u.HandleError(err, 0)
		return
	}

	// 删除 _User 中的 session token 字段
	query := types.M{"objectId": userID}
	update := types.M{
		"sessionToken": types.M{
			"__op": "Delete",
		},
	}
	_, err = orm.TomatoDBController.Update("_User", query, update, types.M{}, false)
	if err != nil {
		u.HandleError(err, 0)
		return
	}

	u.Data["json"] = sessionData
	u.ServeJSON()
}

// Get ...
// @router / [get]
func (u *UpgradeSessionController) Get() {
	u.ClassesController.Get()
}

// Delete ...
// @router / [delete]
func (u *UpgradeSessionController) Delete() {
	u.ClassesController.Delete()
}

// Put ...
// @router / [put]
func (u *UpgradeSessionController) Put() {
	u.ClassesController.Put()
}
