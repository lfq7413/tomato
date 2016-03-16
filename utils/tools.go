package utils

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
