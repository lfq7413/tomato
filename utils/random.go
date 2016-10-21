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

// CreateFileName ...
func CreateFileName() string {
	name := CreateToken()
	name = name[0:8] + "-" + name[8:12] + "-" + name[12:16] + "-" + name[16:20] + "-" + name[20:32]
	return name
}

func CreateString(n int) string {
	return string(utils.RandomCreateBytes(n))
}
