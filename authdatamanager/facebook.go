package authdatamanager

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

type facebook struct{}

func (a facebook) ValidateAuthData(authData types.M, options types.M) error {
	accessToken := utils.S(authData["access_token"])
	host := "https://graph.facebook.com/v2.5/"
	path := "me?fields=id&access_token=" + accessToken
	data, err := request(host+path, nil)
	if err != nil {
		return errs.E(errs.ObjectNotFound, "Failed to validate this access token with Facebook.")
	}
	if data["id"] == nil || utils.S(data["id"]) != utils.S(authData["id"]) {
		return errs.E(errs.ObjectNotFound, "Facebook auth is invalid for this user.")
	}

	// 校验 appIDs
	if options == nil {
		return errs.E(errs.ObjectNotFound, "Facebook auth is not configured.")
	}
	var appIDs []string
	if v, ok := options["appIds"].([]string); ok == true && len(appIDs) > 0 {
		appIDs = v
	} else {
		return errs.E(errs.ObjectNotFound, "Facebook auth is not configured.")
	}
	path = "app?access_token=" + accessToken
	data, err = request(host+path, nil)
	if err != nil {
		return errs.E(errs.ObjectNotFound, "Failed to validate this access token with Facebook.")
	}
	if data["id"] != nil {
		id := utils.S(data["id"])
		for _, v := range appIDs {
			if id == v {
				return nil
			}
		}
	}

	return errs.E(errs.ObjectNotFound, "Facebook auth is invalid for this user.")
}
