package auth

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

type janrainengage struct{}

func (a janrainengage) ValidateAuthData(authData types.M, options types.M) error {
	// 具体接口参考： http://developers.janrain.com/overview/social-login/identity-providers/user-profile-data/#normalized-user-profile-data
	host := "https://rpxnow.com"
	path := "/api/v2/auth_info"
	requestData := map[string]string{
		"token":  utils.S(authData["auth_token"]),
		"apiKey": utils.S(options["api_key"]),
		"format": "json",
	}
	data, err := post(host+path, nil, requestData)
	if err != nil {
		return errs.E(errs.ObjectNotFound, "Failed to validate this access token with Janrain.")
	}
	if utils.S(data["stat"]) == "ok" {
		if profile := utils.M(data["profile"]); profile != nil {
			if utils.S(profile["identifier"]) == utils.S(authData["id"]) {
				return nil
			}
		}
	}
	return errs.E(errs.ObjectNotFound, "Janrain auth is invalid for this user.")
}
