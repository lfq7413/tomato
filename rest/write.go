package rest

import (
	"github.com/lfq7413/tomato/auth"
)

// Write ...
type Write struct {
	auth         *auth.Auth
	className    string
	query        map[string]interface{}
	data         map[string]interface{}
	originalData map[string]interface{}
}

// NewWrite ...
func NewWrite(
	auth *auth.Auth,
	className string,
	query map[string]interface{},
	data map[string]interface{},
	originalData map[string]interface{},
) *Write {
	write := &Write{
		auth:         auth,
		className:    className,
		query:        query,
		data:         data,
		originalData: originalData,
	}
	return write
}

// Execute ...
func (w *Write) Execute() map[string]interface{} {
	return map[string]interface{}{}
}
