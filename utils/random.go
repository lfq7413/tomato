package utils

import (
	"github.com/astaxie/beego/utils"
	"gopkg.in/mgo.v2/bson"
)

// CreateObjectID ...
func CreateObjectID() string {
	return bson.NewObjectId().Hex()
}

// CreateToken ...
func CreateToken() string {
	alphabets := []byte("0123456789ABCDEF")
	return string(utils.RandomCreateBytes(32, alphabets...))
}
