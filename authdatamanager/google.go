package authdatamanager

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

type google struct{}

func (a google) ValidateAuthData(authData types.M, options types.M) error {
	var err error
	if utils.S(authData["id_token"]) != "" {
		err = a.validateIDToken(utils.S(authData["id"]), utils.S(authData["id_token"]))
	} else {
		err = a.validateAuthToken(utils.S(authData["id"]), utils.S(authData["access_token"]))
		if err != nil {
			err = a.validateIDToken(utils.S(authData["id"]), utils.S(authData["access_token"]))
		}
	}
	return err
}

func (a google) validateIDToken(id, token string) error {
	host := "https://www.googleapis.com/oauth2/v3/"
	path := "tokeninfo?id_token=" + token
	data, err := request(host+path, nil)
	if err != nil {
		return errs.E(errs.ObjectNotFound, "Failed to validate this access token with Google.")
	}
	if data != nil && (utils.S(data["sub"]) == id || utils.S(data["user_id"]) == id) {
		return nil
	}
	return errs.E(errs.ObjectNotFound, "Google auth is invalid for this user.")
}

func (a google) validateAuthToken(id, token string) error {
	host := "https://www.googleapis.com/oauth2/v3/"
	path := "tokeninfo?access_token=" + token
	data, err := request(host+path, nil)
	if err != nil {
		return errs.E(errs.ObjectNotFound, "Failed to validate this access token with Google.")
	}
	if data != nil && (utils.S(data["sub"]) == id || utils.S(data["user_id"]) == id) {
		return nil
	}
	return errs.E(errs.ObjectNotFound, "Google auth is invalid for this user.")
}
