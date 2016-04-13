package controllers

import (
	"time"

	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// LoginController ...
type LoginController struct {
	ObjectsController
}

// HandleLogIn ...
// @router / [get]
func (l *LoginController) HandleLogIn() {
	username := l.GetString("username")
	password := l.GetString("password")
	if username == "" {
		// TODO 用户名不能为空
		return
	}
	if password == "" {
		// TODO 密码不能为空
		return
	}

	where := types.M{
		"username": username,
	}
	results := orm.Find("_User", where, types.M{})
	if results == nil || len(results) == 0 {
		// TODO 用户名密码错误
		return
	}
	user := utils.MapInterface(results[0])

	correct := utils.Compare(password, utils.String(user["password"]))
	if correct == false {
		// TODO 用户名密码错误
		return
	}

	token := "r:" + utils.CreateToken()
	user["sessionToken"] = token
	delete(user, "password")

	if user["authData"] != nil {
		authData := utils.MapInterface(user["authData"])
		for k, v := range authData {
			if v == nil {
				delete(authData, k)
			}
		}
		if len(authData) == 0 {
			delete(user, "authData")
		}
	}

	// TODO 展开文件信息

	expiresAt := time.Now().UTC()
	expiresAt = expiresAt.AddDate(1, 0, 0)
	usr := types.M{
		"__type":    "Pointer",
		"className": "_User",
		"objectId":  user["objectId"],
	}
	createdWith := types.M{
		"action":       "login",
		"authProvider": "password",
	}
	sessionData := types.M{
		"sessionToken": token,
		"user":         usr,
		"createdWith":  createdWith,
		"restricted":   false,
		"expiresAt":    utils.TimetoString(expiresAt),
	}
	if l.Info.InstallationID != "" {
		sessionData["installationId"] = l.Info.InstallationID
	}

	rest.NewWrite(rest.Master(), "_Session", nil, sessionData, nil).Execute()

	l.Data["json"] = user
	l.ServeJSON()

}

// Post ...
// @router / [post]
func (l *LoginController) Post() {
	l.ObjectsController.Post()
}

// Delete ...
// @router / [delete]
func (l *LoginController) Delete() {
	l.ObjectsController.Delete()
}

// Put ...
// @router / [put]
func (l *LoginController) Put() {
	l.ObjectsController.Put()
}
