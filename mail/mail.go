package mail

import "github.com/lfq7413/tomato/types"

// Adapter ...
type Adapter interface {
	// SendMail 包含三个参数：
	// to 接收方地址
	// text 邮件内容
	// subject 邮件主题
	SendMail(types.M) error
}
