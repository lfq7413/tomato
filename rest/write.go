package rest

// Write ...
type Write struct {
	auth         *Auth
	className    string
	query        map[string]interface{}
	data         map[string]interface{}
	originalData map[string]interface{}
}

// NewWrite ...
func NewWrite(
	auth *Auth,
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
