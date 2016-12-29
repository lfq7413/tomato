package auth

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

type instagram struct{}

func (a instagram) ValidateAuthData(authData types.M, options types.M) error {
	host := "https://api.instagram.com/v1/"
	path := "users/self/?access_token=" + utils.S(authData["access_token"])
	data, err := request(host+path, nil)
	if err != nil {
		return errs.E(errs.ObjectNotFound, "Failed to validate this access token with Instagram.")
	}
	if d := utils.M(data["data"]); d != nil {
		if d["id"] != nil && utils.S(d["id"]) == utils.S(authData["id"]) {
			return nil
		}
	}
	return errs.E(errs.ObjectNotFound, "Instagram auth is invalid for this user.")
}
