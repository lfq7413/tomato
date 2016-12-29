package auth

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

type janraincapture struct{}

func (a janraincapture) ValidateAuthData(authData types.M, options types.M) error {
	// 具体接口参考： https://docs.janrain.com/api/registration/entity/#entity
	host := utils.S(options["janrain_capture_host"])
	if host == "" {
		return errs.E(errs.ObjectNotFound, "Janrain auth is invalid for this user.")
	}
	path := "/entity?attribute_name=uuid&access_token=" + utils.S(authData["access_token"])
	data, err := request(host+path, nil)
	if err != nil {
		return errs.E(errs.ObjectNotFound, "Failed to validate this access token with Janrain.")
	}
	if utils.S(data["stat"]) == "ok" && utils.S(data["result"]) == utils.S(authData["id"]) {
		return nil
	}
	return errs.E(errs.ObjectNotFound, "Janrain auth is invalid for this user.")
}
