package rest

import (
	"github.com/lfq7413/tomato/auth"
)

// Destroy ...
type Destroy struct {
	auth      *auth.Auth
	className string
	query     map[string]interface{}
}

// NewDestroy ...
func NewDestroy(
	auth *auth.Auth,
	className string,
	query map[string]interface{},
	originalData map[string]interface{},
) *Destroy {
	destroy := &Destroy{
		auth:      auth,
		className: className,
		query:     query,
	}
	return destroy
}

// Execute ...
func (d *Destroy) Execute() map[string]interface{} {
	return nil
}
