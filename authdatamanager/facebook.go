package authdatamanager

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

type facebook struct{}

func (a facebook) ValidateAuthData(authData types.M) error {
	host := "https://graph.facebook.com/v2.5/"
	path := "me?fields=id&access_token=" + utils.S(authData["access_token"])
	data, err := request(host+path, nil)
	if err != nil {
		return errs.E(errs.ObjectNotFound, "Failed to validate this access token with Facebook.")
	}
	if data["id"] != nil && utils.S(data["id"]) == utils.S(authData["id"]) {
		return nil
	}
	return errs.E(errs.ObjectNotFound, "Facebook auth is invalid for this user.")
}
