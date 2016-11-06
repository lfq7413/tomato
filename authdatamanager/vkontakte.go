package authdatamanager

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

type vkontakte struct{}

func (a vkontakte) ValidateAuthData(authData types.M, params types.M) error {
	response, err := a.vkOAuth2Request(params)
	if err != nil {
		return errs.E(errs.ObjectNotFound, "Failed to validate this access token with Vk.")
	}
	if response != nil && utils.S(response["access_token"]) != "" {
		host := "https://api.vk.com/"
		path := "method/secure.checkToken?token=" + utils.S(authData["access_token"]) +
			"&client_secret=" + utils.S(params["appSecret"]) +
			"&access_token=" + utils.S(response["access_token"])
		data, err := request(host+path, nil)
		if err != nil {
			return errs.E(errs.ObjectNotFound, "Failed to validate this access token with Vk.")
		}
		if res := utils.M(data["response"]); res != nil {
			if utils.S(res["user_id"]) == utils.S(authData["id"]) {
				return nil
			}
		}
		return errs.E(errs.ObjectNotFound, "Vk auth is invalid for this user.")
	}
	return errs.E(errs.ObjectNotFound, "Vk appIds or appSecret is incorrect.")
}

func (a vkontakte) vkOAuth2Request(params types.M) (types.M, error) {
	if params == nil || utils.S(params["appIds"]) == "" || utils.S(params["appSecret"]) == "" {
		return nil, errs.E(errs.ObjectNotFound, "Vk auth is not configured. Missing appIds or appSecret.")
	}
	host := "https://oauth.vk.com/"
	path := "access_token?client_id=" + utils.S(params["appIds"]) + "&client_secret=" + utils.S(params["appSecret"])
	return request(host+path, nil)
}
