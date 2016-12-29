package auth

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

type yixin struct{}

func (a yixin) ValidateAuthData(authData types.M, options types.M) error {
	// 具体接口参考： https://open.yixin.im/document/oauth/api
	host := "https://open.yixin.im/api/"
	path := "userinfo?access_token=" + utils.S(authData["access_token"])
	data, err := request(host+path, nil)
	if err != nil {
		return errs.E(errs.ObjectNotFound, "Failed to validate this access token with Yixin.")
	}
	if code, ok := data["code"].(float64); ok && code == 1 {
		if userinfo := utils.M(data["userinfo"]); userinfo != nil {
			if utils.S(userinfo["accountId"]) == utils.S(authData["id"]) {
				return nil
			}
		}
	}
	return errs.E(errs.ObjectNotFound, "Yixin auth is invalid for this user.")
}
