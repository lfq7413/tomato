package push

import "github.com/lfq7413/tomato/types"

type tomatoPushAdapter struct {
	validPushTypes []string
}

func newTomatoPush() *tomatoPushAdapter {
	t := &tomatoPushAdapter{
		validPushTypes: []string{"ios", "android"},
	}
	return t
}

func (t *tomatoPushAdapter) send(data types.M, installations types.S) {

}

func (t *tomatoPushAdapter) getValidPushTypes() []string {
	return t.validPushTypes
}
