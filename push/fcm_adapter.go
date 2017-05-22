package push

import "github.com/lfq7413/tomato/types"

type fcmPushAdapter struct {
	validPushTypes []string
	serverKey      string
}

func newFCMPush() *fcmPushAdapter {
	f := &fcmPushAdapter{
		validPushTypes: []string{"ios", "android"},
	}
	return f
}

func (f *fcmPushAdapter) send(body types.M, installations types.S, pushStatus string) []types.M {
	return []types.M{}
}

func (f *fcmPushAdapter) getValidPushTypes() []string {
	return f.validPushTypes
}
