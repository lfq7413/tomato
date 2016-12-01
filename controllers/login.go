package controllers

import (
	"time"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/files"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// LoginController 处理 /login 接口的请求
type LoginController struct {
	ClassesController
}

// HandleLogIn 处理登录请求
// @router / [get]
func (l *LoginController) HandleLogIn() {
	var username, password string
	if l.JSONBody != nil && l.JSONBody["username"] != nil {
		username = utils.S(l.JSONBody["username"])
	} else {
		username = l.Query["username"]
	}
	if l.JSONBody != nil && l.JSONBody["password"] != nil {
		password = utils.S(l.JSONBody["password"])
	} else {
		password = l.Query["password"]
	}

	if username == "" {
		l.HandleError(errs.E(errs.UsernameMissing, "username is required."), 0)
		return
	}
	if password == "" {
		l.HandleError(errs.E(errs.PasswordMissing, "password is required."), 0)
		return
	}

	where := types.M{
		"username": username,
	}
	results, err := orm.TomatoDBController.Find("_User", where, types.M{})
	if err != nil {
		l.HandleError(err, 0)
		return
	}
	if results == nil || len(results) == 0 {
		l.HandleError(errs.E(errs.ObjectNotFound, "Invalid username/password."), 0)
		return
	}
	user := utils.M(results[0])

	var emailVerified bool
	if _, ok := user["emailVerified"]; ok {
		if v, ok := user["emailVerified"].(bool); ok {
			emailVerified = v
		}
	}
	if config.TConfig.VerifyUserEmails && config.TConfig.PreventLoginWithUnverifiedEmail && emailVerified == false {
		// 拒绝未验证邮箱的用户登录
		l.HandleError(errs.E(errs.EmailNotFound, "User email is not verified."), 0)
		return
	}

	// TODO 换用高强度的加密方式
	correct := utils.Compare(password, utils.S(user["password"]))
	accountLockoutPolicy := rest.NewAccountLockout(utils.S(user["username"]))
	err = accountLockoutPolicy.HandleLoginAttempt(correct)
	if err != nil {
		l.HandleError(err, 0)
		return
	}
	if correct == false {
		l.HandleError(errs.E(errs.ObjectNotFound, "Invalid username/password."), 0)
		return
	}

	// 检测密码是否过期
	if config.TConfig.PasswordPolicy && config.TConfig.MaxPasswordAge > 0 {
		if changedAt, ok := user["_password_changed_at"].(time.Time); ok {
			// 密码过期时间戳存在，判断是否过期
			expiresAt := changedAt.Add(time.Duration(config.TConfig.MaxPasswordAge) * 24 * time.Hour)
			if expiresAt.UnixNano() < time.Now().UnixNano() {
				l.HandleError(errs.E(errs.ObjectNotFound, "Your password has expired. Please reset your password."), 0)
				return
			}
		} else {
			// 在启用密码过期之前的数据，需要增加该字段
			query := types.M{"username": user["username"]}
			update := types.M{"_password_changed_at": utils.TimetoString(time.Now().UTC())}
			orm.TomatoDBController.Update("_User", query, update, types.M{}, false)
		}
	}

	token := "r:" + utils.CreateToken()
	user["sessionToken"] = token
	delete(user, "password")

	if user["authData"] != nil {
		authData := utils.M(user["authData"])
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

	expiresAt := config.GenerateSessionExpiresAt()
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
	write, err := rest.NewWrite(rest.Master(), "_Session", nil, sessionData, nil, l.Info.ClientSDK)
	if err != nil {
		l.HandleError(err, 0)
		return
	}
	_, err = write.Execute()
	if err != nil {
		l.HandleError(err, 0)
		return
	}

	l.Data["json"] = user
	l.ServeJSON()

}

// Post ...
// @router / [post]
func (l *LoginController) Post() {
	l.ClassesController.Post()
}

// Delete ...
// @router / [delete]
func (l *LoginController) Delete() {
	l.ClassesController.Delete()
}

// Put ...
// @router / [put]
func (l *LoginController) Put() {
	l.ClassesController.Put()
}
