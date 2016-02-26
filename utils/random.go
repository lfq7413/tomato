package utils

import "github.com/astaxie/beego/utils"

// CreateObjectId ...
func CreateObjectId() string {
    return string(utils.RandomCreateBytes(32))
}