package rest

import (
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/utils"
)

func shouldVerifyEmails() bool {
	return config.TConfig.VerifyUserEmails
}

// SetEmailVerifyToken ...
func SetEmailVerifyToken(user map[string]interface{}) {
	if shouldVerifyEmails() {
		user["_email_verify_token"] = utils.CreateToken()
		user["emailVerified"] = false
	}
}

// SendVerificationEmail ...
func SendVerificationEmail(user map[string]interface{}) {
	if shouldVerifyEmails() == false {
		return
	}
	// TODO 发送验证邮件
}
