package utils

import "github.com/lfq7413/tomato/types"

// A 将 interface{} 转换为 []interface{}
func A(i interface{}) []interface{} {
	if v, ok := i.([]interface{}); ok {
		return v
	}
	if v, ok := i.(types.S); ok {
		return v
	}
	return nil
}

// M 将 interface{} 转换为 map[string]interface{}
func M(i interface{}) map[string]interface{} {
	if v, ok := i.(map[string]interface{}); ok {
		return v
	}
	if v, ok := i.(types.M); ok {
		return v
	}
	return nil
}

// String 将 interface{} 转换为 string
func S(i interface{}) string {
	if v, ok := i.(string); ok {
		return v
	}
	return ""
}
