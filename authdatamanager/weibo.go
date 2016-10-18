package authdatamanager

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

type weibo struct{}

func (a weibo) ValidateAuthData(authData types.M, options types.M) error {
	// 具体接口参考： http://open.weibo.com/wiki/Oauth2/get_token_info
	host := "https://api.weibo.com/oauth2/"
	path := "get_token_info"
	requestData := map[string]string{
		"access_token": utils.S(authData["access_token"]),
	}
	data, err := post(host+path, nil, requestData)
	if err != nil {
		return errs.E(errs.ObjectNotFound, "Failed to validate this access token with Weibo.")
	}
	if data["uid"] != nil && utils.S(data["uid"]) == utils.S(authData["id"]) {
		return nil
	}
	return errs.E(errs.ObjectNotFound, "Weibo auth is invalid for this user.")
}
