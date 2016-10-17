package authdatamanager

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

type google struct{}

func (a google) ValidateAuthData(authData types.M) error {
	host := "https://www.googleapis.com/oauth2/v3/"
	path := "tokeninfo?id_token=" + utils.S(authData["access_token"])
	data, err := request(host+path, nil)
	if err != nil {
		return errs.E(errs.ObjectNotFound, "Failed to validate this access token with Google.")
	}
	if data["sub"] != nil && utils.S(data["sub"]) == utils.S(authData["id"]) {
		return nil
	}
	return errs.E(errs.ObjectNotFound, "Google auth is invalid for this user.")
}
