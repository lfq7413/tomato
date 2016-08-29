package mail

import (
	"testing"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/types"
)

func Test_smtp(t *testing.T) {
	config.TConfig = &config.Config{
		SMTPServer:   "smtp.163.com",
		MailUsername: "user@163.com",
		MailPassword: "password",
	}

	s := NewSMTPAdapter()
	object := types.M{
		"text":    "text from tomato",
		"to":      "user@163.com",
		"subject": "tomato send",
	}
	s.SendMail(object)
}
