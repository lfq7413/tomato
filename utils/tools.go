package utils

import "github.com/lfq7413/tomato/utils"

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
		utils.SliceInterface(response["results"]) == nil ||
		len(utils.SliceInterface(response["results"])) == 0 {
		return false
	}
	return true
}
