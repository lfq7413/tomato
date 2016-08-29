package mail

import (
	"net/smtp"

	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)
import "github.com/lfq7413/tomato/config"

// SMTPMailAdapter ...
type SMTPMailAdapter struct {
	server   string
	username string
	password string
}

// NewSMTPAdapter ...
func NewSMTPAdapter() *SMTPMailAdapter {
	s := &SMTPMailAdapter{
		server:   config.TConfig.SMTPServer,
		username: config.TConfig.MailUsername,
		password: config.TConfig.MailPassword,
	}
	return s
}

// SendMail ...
func (s *SMTPMailAdapter) SendMail(object types.M) error {
	// TODO 打印错误日志
	if s.server == "" || s.username == "" || s.password == "" {
		return nil
	}

	receiver := utils.S(object["to"])
	subject := utils.S(object["subject"])
	text := utils.S(object["text"])

	auth := smtp.PlainAuth("", s.username, s.password, s.server)
	to := []string{receiver}
	msg := []byte("To: " + receiver + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		text + "\r\n")
	err := smtp.SendMail(s.server+":25", auth, s.username, to, msg)
	if err != nil {
		// 打印错误
	}
	return nil
}
