package utils

import "github.com/astaxie/beego/utils"

// CreateObjectID ...
func CreateObjectID() string {
	return string(utils.RandomCreateBytes(32))
}

// CreateToken ...
func CreateToken() string {
	alphabets := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 'A', 'B', 'C', 'D', 'E', 'F'}
	return string(utils.RandomCreateBytes(32, alphabets...))
}
