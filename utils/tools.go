package utils

// Append 连接两个 slice
func Append(slice1 []interface{}, slice2 []interface{}) []interface{} {
	s := slice1
	for _, v := range slice2 {
		s = append(s, v)
	}
	return s
}
