package utils

import (
	"reflect"
	"regexp"

	"github.com/lfq7413/tomato/types"
)

// HasResults Find() 返回数据中是否有结果
func HasResults(response types.M) bool {
	if response == nil ||
		response["results"] == nil ||
		A(response["results"]) == nil ||
		len(A(response["results"])) == 0 {
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
	return Copy(i)
}

// CopyMap 复制 map
func CopyMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		return nil
	}
	d := map[string]interface{}{}
	for k, v := range m {
		d[k] = DeepCopy(v)
	}
	return d
}

// CopySlice 复制 slice
func CopySlice(s []interface{}) []interface{} {
	if s == nil {
		return nil
	}
	d := []interface{}{}
	for _, v := range s {
		d = append(d, DeepCopy(v))
	}
	return d
}

// CopyMapM 复制 map
func CopyMapM(m types.M) types.M {
	if m == nil {
		return nil
	}
	d := types.M{}
	for k, v := range m {
		d[k] = DeepCopy(v)
	}
	return d
}

// CopySliceS 复制 slice
func CopySliceS(s types.S) types.S {
	if s == nil {
		return nil
	}
	d := types.S{}
	for _, v := range s {
		d = append(d, DeepCopy(v))
	}
	return d
}

// CompareArray 比较两个数组是否相等，忽略数组顺序
func CompareArray(i1, i2 interface{}) bool {
	if i1 == nil && i2 == nil {
		return true
	}
	if v1 := A(i1); v1 != nil {
		if v2 := A(i2); v2 != nil {
			// TODO 去重
			if len(v1) != len(v2) {
				return false
			}

			for _, a := range v1 {
				match := false
				for _, b := range v2 {
					if reflect.DeepEqual(a, b) {
						match = true
						break
					}
				}
				if match == false {
					return false
				}
			}
			return true
		}
		return false
	}
	return false
}
