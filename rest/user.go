package rest

import (
	"errors"
	"net/url"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/mail"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

var adapter mail.Adapter

func init() {
	a := config.TConfig.MailAdapter
	if a == "smtp" {
		adapter = mail.NewSMTPAdapter()
	} else {
		adapter = mail.NewSMTPAdapter()
	}
}

// shouldVerifyEmails 根据配置参数确定是否需要验证邮箱
func shouldVerifyEmails() bool {
	return config.TConfig.VerifyUserEmails
}

// SetEmailVerifyToken 设置需要验证的 token
func SetEmailVerifyToken(user types.M) {
	if shouldVerifyEmails() {
		user["_email_verify_token"] = utils.CreateToken()
		user["emailVerified"] = false
	}
}

// SendVerificationEmail 发送验证邮件
func SendVerificationEmail(user types.M) {
	if shouldVerifyEmails() == false {
		return
	}
	token := url.QueryEscape(user["_email_verify_token"].(string))
	user = getUserIfNeeded(user)
	if user == nil {
		return
	}
	user["className"] = "_User"
	username := url.QueryEscape(user["username"].(string))
	link := config.TConfig.ServerURL + "app/verify_email" + "?token=" + token + "&username=" + username
	options := types.M{
		"appName": config.TConfig.AppName,
		"link":    link,
		"user":    user,
	}
	adapter.SendMail(defaultVerificationEmail(options))
}

// getUserIfNeeded 把 user 填充完整，如果无法完成则返回 nil
func getUserIfNeeded(user types.M) types.M {
	if user["username"] != nil && user["email"] != nil {
		return user
	}
	where := types.M{}
	if user["username"] != nil {
		where["username"] = user["username"]
	}
	if user["email"] != nil {
		where["email"] = user["email"]
	}

	query, err := NewQuery(Master(), "_User", where, types.M{})
	if err != nil {
		return nil
	}
	response, err := query.Execute()
	if err != nil {
		return nil
	}
	if utils.HasResults(response) == false {
		return nil
	}

	return response["results"].([]interface{})[0].(map[string]interface{})
}

func defaultVerificationEmail(options types.M) types.M {
	user := utils.MapInterface(options["user"])
	text := "Hi,\n\n"
	text += "You are being asked to confirm the e-mail address " + user["email"].(string)
	text += " with " + options["appName"].(string) + "\n\n"
	text += "Click here to confirm it:\n" + options["link"].(string)
	to := user["email"].(string)
	subject := "Please verify your e-mail for " + options["appName"].(string)
	return types.M{
		"text":    text,
		"to":      to,
		"subject": subject,
	}
}

// SendPasswordResetEmail 发送密码重置邮件
func SendPasswordResetEmail(email string) error {
	user := setPasswordResetToken(email)
	if user == nil || len(user) == 0 {
		return errs.E(errs.EmailMissing, "you must provide an email")
	}
	user["className"] = "_User"
	token := url.QueryEscape(user["_perishable_token"].(string))
	username := url.QueryEscape(user["username"].(string))
	link := config.TConfig.ServerURL + "app/request_password_reset" + "?token=" + token + "&username=" + username
	options := types.M{
		"appName": config.TConfig.AppName,
		"link":    link,
		"user":    user,
	}
	adapter.SendMail(defaultResetPasswordEmail(options))
	return nil
}

// setPasswordResetToken 设置修改密码 token
func setPasswordResetToken(email string) types.M {
	token := utils.CreateToken()
	collection := orm.AdaptiveCollection("_User")
	where := types.M{"email": email}
	update := types.M{
		"$set": types.M{"_perishable_token": token},
	}
	return collection.FindOneAndUpdate(where, update)
}

func defaultResetPasswordEmail(options types.M) types.M {
	user := utils.MapInterface(options["user"])
	text := "Hi,\n\n"
	text += "You requested to reset your password for " + options["appName"].(string) + "\n\n"
	text += "Click here to reset it:\n" + options["link"].(string)
	to := user["email"].(string)
	subject := "Password Reset for " + options["appName"].(string)
	return types.M{
		"text":    text,
		"to":      to,
		"subject": subject,
	}
}

// VerifyEmail 更新邮箱验证标志
func VerifyEmail(username, token string) bool {
	if shouldVerifyEmails() == false {
		return false
	}

	collection := orm.AdaptiveCollection("_User")
	where := types.M{
		"username":            username,
		"_email_verify_token": token,
	}
	update := types.M{
		"$set": types.M{"emailVerified": true},
	}
	document := collection.FindOneAndUpdate(where, update)
	if document == nil || len(document) == 0 {
		return false
	}

	return true
}

// CheckResetTokenValidity 检查要重置密码的用户与 token 是否存在
func CheckResetTokenValidity(username, token string) types.M {
	collection := orm.AdaptiveCollection("_User")
	where := types.M{
		"username":          username,
		"_perishable_token": token,
	}
	option := types.M{"limit": 1}
	results := collection.Find(where, option)
	if len(results) != 1 {
		return nil
	}

	return results[0]
}

// UpdatePassword 更新指定用户的密码
func UpdatePassword(username, token, newPassword string) error {
	user := CheckResetTokenValidity(username, token)
	if user == nil {
		return errors.New("Invalid token")
	}

	err := updateUserPassword(user["_id"].(string), newPassword)
	if err != nil {
		return err
	}

	// 清空重置密码 token
	collection := orm.AdaptiveCollection("_User")
	selector := types.M{"username": username}
	update := types.M{
		"$unset": types.M{"_perishable_token": nil},
	}
	collection.FindOneAndUpdate(selector, update)

	return nil
}

func updateUserPassword(userID, password string) error {
	_, err := Update(Master(), "_User", userID, types.M{"password": password})
	if err != nil {
		return err
	}
	return nil
}
