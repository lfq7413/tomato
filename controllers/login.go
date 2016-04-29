package controllers

import (
	"time"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/files"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// LoginController 处理 /login 接口的请求
type LoginController struct {
	ObjectsController
}

// HandleLogIn 处理登录请求
// @router / [get]
func (l *LoginController) HandleLogIn() {
	username := l.GetString("username")
	password := l.GetString("password")
	if username == "" {
		l.Data["json"] = errs.ErrorMessageToMap(errs.UsernameMissing, "username is required.")
		l.ServeJSON()
		return
	}
	if password == "" {
		l.Data["json"] = errs.ErrorMessageToMap(errs.PasswordMissing, "password is required.")
		l.ServeJSON()
		return
	}

	where := types.M{
		"username": username,
	}
	results, err := orm.Find("_User", where, types.M{})
	if err != nil {
		l.Data["json"] = errs.ErrorToMap(err)
		l.ServeJSON()
		return
	}
	if results == nil || len(results) == 0 {
		l.Data["json"] = errs.ErrorMessageToMap(errs.ObjectNotFound, "Invalid username/password.")
		l.ServeJSON()
		return
	}
	user := utils.MapInterface(results[0])

	// TODO 换用高强度的加密方式
	correct := utils.Compare(password, utils.String(user["password"]))
	if correct == false {
		l.Data["json"] = errs.ErrorMessageToMap(errs.ObjectNotFound, "Invalid username/password.")
		l.ServeJSON()
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

	// 展开文件信息
	files.ExpandFilesInObject(user)

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
	// 为新登录用户创建 sessionToken
	write, err := rest.NewWrite(rest.Master(), "_Session", nil, sessionData, nil)
	if err != nil {
		l.Data["json"] = errs.ErrorToMap(err)
		l.ServeJSON()
		return
	}
	_, err = write.Execute()
	if err != nil {
		l.Data["json"] = errs.ErrorToMap(err)
		l.ServeJSON()
		return
	}

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
