package utils

// SliceInterface 将 interface{} 转换为 []interface{}
func SliceInterface(i interface{}) []interface{} {
	if v, ok := i.([]interface{}); ok {
		return v
	}
	return nil
}

// MapInterface 将 interface{} 转换为 map[string]interface{}
func MapInterface(i interface{}) map[string]interface{} {
	if v, ok := i.(map[string]interface{}); ok {
		return v
	}
	return nil
}

// String 将 interface{} 转换为 string
func String(i interface{}) string {
	if v, ok := i.(string); ok {
		return v
	}
	return ""
}
