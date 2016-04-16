package rest

import (
	"net/url"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/mail"
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
	user = getUserIfNeeded(user)
	user["className"] = "_User"
	token := url.QueryEscape(user["_email_verify_token"].(string))
	username := url.QueryEscape(user["username"].(string))
	link := config.TConfig.VerifyEmailURL + "?token=" + token + "&username=" + username
	options := types.M{
		"appName": config.TConfig.AppName,
		"link":    link,
		"user":    user,
	}
	adapter.SendMail(defaultVerificationEmail(options))
}

func getUserIfNeeded(user types.M) types.M {
	// TODO
	return nil
}

func defaultVerificationEmail(options types.M) types.M {
	// TODO
	return nil
}

// SendPasswordResetEmail ...
func SendPasswordResetEmail(email string) error {
	// TODO
	return nil
}
