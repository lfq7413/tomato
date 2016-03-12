package conv

// SliceInterface ...
func SliceInterface(i interface{}) []interface{} {
	if v, ok := i.([]interface{}); ok {
		return v
	}
	return nil
}
