package authdatamanager

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

type douban struct{}

func (a douban) ValidateAuthData(authData types.M, options types.M) error {
	// 具体接口参考： https://developers.douban.com/wiki/?title=connect
	host := "https://api.douban.com/v2/"
	path := "user/~me"
	headers := map[string]string{
		"Authorization": "Bearer " + utils.S(authData["access_token"]),
	}
	data, err := request(host+path, headers)
	if err != nil {
		return errs.E(errs.ObjectNotFound, "Failed to validate this access token with Douban.")
	}
	if data["id"] != nil && utils.S(data["id"]) == utils.S(authData["id"]) {
		return nil
	}
	return errs.E(errs.ObjectNotFound, "Douban auth is invalid for this user.")
}
