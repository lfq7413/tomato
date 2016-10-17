package authdatamanager

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

type spotify struct{}

func (a spotify) ValidateAuthData(authData types.M, options types.M) error {
	host := "https://api.spotify.com/"
	path := "v1/me"
	headers := map[string]string{
		"Authorization": "Bearer " + utils.S(authData["access_token"]),
	}
	data, err := request(host+path, headers)
	if err != nil {
		return errs.E(errs.ObjectNotFound, "Failed to validate this access token with Spotify.")
	}
	if data["id"] == nil || utils.S(data["id"]) != utils.S(authData["id"]) {
		return errs.E(errs.ObjectNotFound, "Spotify auth is invalid for this user.")
	}
	return nil

	// TODO 校验 appIDs ，调用接口与获取用户信息相同，存在问题
	// if options == nil {
	// 	return errs.E(errs.ObjectNotFound, "Spotify auth is not configured.")
	// }
	// var appIDs []string
	// if v, ok := options["appIds"].([]string); ok == true && len(appIDs) > 0 {
	// 	appIDs = v
	// } else {
	// 	return errs.E(errs.ObjectNotFound, "Spotify auth is not configured.")
	// }
	// if data["id"] != nil {
	// 	id := utils.S(data["id"])
	// 	for _, v := range appIDs {
	// 		if id == v {
	// 			return nil
	// 		}
	// 	}
	// }

	// return errs.E(errs.ObjectNotFound, "Spotify auth is invalid for this user.")
}
