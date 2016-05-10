package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
)

// Hash ...
func Hash(password string) string {
	h := sha256.New()
	io.WriteString(h, password)
	s := fmt.Sprintf("%x", h.Sum(nil))
	return s
}

// Compare ...
func Compare(password string, hashedPassword string) bool {
	s := Hash(password)
	if s == hashedPassword {
		return true
	}
	return false
}

// MD5Hash ...
func MD5Hash(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	r := fmt.Sprintf("%x", h.Sum(nil))
	return r
}
