package mail

import "github.com/lfq7413/tomato/types"

// Adapter ...
type Adapter interface {
	SendMail(types.M) error
}
