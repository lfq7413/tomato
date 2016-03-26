package utils

import "regexp"

// AppendString 连接两个 []string
func AppendString(slice1 []string, slice2 []string) []string {
	s := slice1
	if slice2 == nil {
		return s
	}
	for _, v := range slice2 {
		s = append(s, v)
	}
	return s
}

// AppendInterface 连接两个 []interface{}
func AppendInterface(slice1 []interface{}, slice2 []interface{}) []interface{} {
	s := slice1
	if slice2 == nil {
		return s
	}
	for _, v := range slice2 {
		s = append(s, v)
	}
	return s
}

// HasResults Find() 返回数据中是否有结果
func HasResults(response map[string]interface{}) bool {
	if response == nil ||
		response["results"] == nil ||
		SliceInterface(response["results"]) == nil ||
		len(SliceInterface(response["results"])) == 0 {
		return false
	}
	return true
}

// IsEmail ...
func IsEmail(email string) bool {
	b, _ := regexp.MatchString("^.+@.+$", email)
	return b
}

// DeepCopy 简易版的内存复制
func DeepCopy(i interface{}) interface{} {
	if i == nil {
		return nil
	}
	if s, ok := i.([]interface{}); ok {
		return CopySlice(s)
	}
	if m, ok := i.(map[string]interface{}); ok {
		return CopyMap(m)
	}
	return i
}

// CopyMap 复制 map
func CopyMap(m map[string]interface{}) map[string]interface{} {
	d := map[string]interface{}{}
	for k, v := range m {
		d[k] = DeepCopy(v)
	}
	return d
}

// CopySlice 复制 slice
func CopySlice(s []interface{}) []interface{} {
	d := []interface{}{}
	for _, v := range s {
		d = append(d, DeepCopy(v))
	}
	return d
}
