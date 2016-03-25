package utils

import "github.com/astaxie/beego/utils"

// CreateObjectID ...
func CreateObjectID() string {
	return string(utils.RandomCreateBytes(32))
}

// CreateToken ...
func CreateToken() string {
	alphabets := []byte("0123456789ABCDEF")
	return string(utils.RandomCreateBytes(32, alphabets...))
}
