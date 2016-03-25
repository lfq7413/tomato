package push

type tomatoPushAdapter struct {
	validPushTypes []string
}

func newTomatoPush() *tomatoPushAdapter {
	t := &tomatoPushAdapter{
		validPushTypes: []string{"ios", "android"},
	}
	return t
}

func (t *tomatoPushAdapter) send() {

}

func (t *tomatoPushAdapter) getValidPushTypes() []string {
	return t.validPushTypes
}
