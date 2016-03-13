package utils

// AppendString 连接两个 slice
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
