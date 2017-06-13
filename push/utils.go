package push

import (
	"strings"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

func isPushIncrementing(body types.M) bool {
	if body == nil {
		return false
	}
	data := utils.M(body["data"])
	if data == nil {
		return false
	}
	badge := utils.S(data["badge"])
	if strings.ToLower(badge) != "increment" {
		return false
	}
	return true
}

// validatePushType 校验查询条件中的推送类型
// where 查询条件
// validPushTypes 当前推送模块支持的类型
func validatePushType(where types.M, validPushTypes []string) error {
	deviceTypeField := where["deviceType"]
	if deviceTypeField == nil {
		deviceTypeField = types.M{}
	}
	deviceTypes := []string{}
	if utils.S(deviceTypeField) != "" {
		deviceTypes = append(deviceTypes, utils.S(deviceTypeField))
	} else if utils.M(deviceTypeField) != nil {
		m := utils.M(deviceTypeField)
		if utils.A(m["$in"]) != nil {
			s := utils.A(m["$in"])
			for _, v := range s {
				deviceTypes = append(deviceTypes, utils.S(v))
			}
		}
	}
	for _, deviceType := range deviceTypes {
		b := false
		for _, t := range validPushTypes {
			if deviceType == t {
				b = true
				break
			}
		}
		if b == false {
			return errs.E(errs.PushMisconfigured, deviceType+" is not supported push type.")
		}
	}

	return nil
}
