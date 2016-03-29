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

func (t *tomatoPushAdapter) send(data map[string]interface{}, installations []interface{}) {

}

func (t *tomatoPushAdapter) getValidPushTypes() []string {
	return t.validPushTypes
}
