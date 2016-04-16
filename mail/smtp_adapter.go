package mail

import "github.com/lfq7413/tomato/types"

// SMTPMailAdapter ...
type SMTPMailAdapter struct {
}

// NewSMTPAdapter ...
func NewSMTPAdapter() *SMTPMailAdapter {
	s := &SMTPMailAdapter{}
	return s
}

// SendMail ...
func (s *SMTPMailAdapter) SendMail(types.M) error {
	// TODO
	return nil
}
