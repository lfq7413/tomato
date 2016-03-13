package rest

// Destroy ...
type Destroy struct {
	auth         *Auth
	className    string
	query        map[string]interface{}
	originalData map[string]interface{}
}

// NewDestroy ...
func NewDestroy(
	auth *Auth,
	className string,
	query map[string]interface{},
	originalData map[string]interface{},
) *Destroy {
	destroy := &Destroy{
		auth:         auth,
		className:    className,
		query:        query,
		originalData: originalData,
	}
	return destroy
}

// Execute ...
func (d *Destroy) Execute() map[string]interface{} {
	return nil
}
