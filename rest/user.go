package rest

import (
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

func shouldVerifyEmails() bool {
	return config.TConfig.VerifyUserEmails
}

// SetEmailVerifyToken ...
func SetEmailVerifyToken(user types.M) {
	if shouldVerifyEmails() {
		user["_email_verify_token"] = utils.CreateToken()
		user["emailVerified"] = false
	}
}

// SendVerificationEmail ...
func SendVerificationEmail(user types.M) {
	if shouldVerifyEmails() == false {
		return
	}
	// TODO 发送验证邮件
}

// SendPasswordResetEmail ...
func SendPasswordResetEmail(email string) error {
	// TODO
	return nil
}
