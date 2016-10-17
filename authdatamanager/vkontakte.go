package authdatamanager

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

type vkontakte struct{}

func (a vkontakte) ValidateAuthData(authData types.M, options types.M) error {
	host := "https://api.vk.com/method/"
	path := "users.get?v=V&access_token=" + utils.S(authData["access_token"])
	data, err := request(host+path, nil)
	if err != nil {
		return errs.E(errs.ObjectNotFound, "Failed to validate this access token with Vk.")
	}
	if response := utils.A(data["response"]); len(response) > 0 {
		if res := utils.M(response[0]); res != nil {
			if utils.S(res["uid"]) == utils.S(authData["id"]) {
				return nil
			}
		}
	}
	return errs.E(errs.ObjectNotFound, "Vk auth is invalid for this user.")
}
