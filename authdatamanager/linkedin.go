package authdatamanager

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

type linkedin struct{}

func (a linkedin) ValidateAuthData(authData types.M, options types.M) error {
	host := "https://api.linkedin.com/"
	path := "v1/people/~:(id)"
	headers := map[string]string{
		"Authorization": "Bearer " + utils.S(authData["access_token"]),
		"x-li-format":   "json",
	}
	if v, ok := authData["is_mobile_sdk"].(bool); ok && v {
		headers["x-li-src"] = "msdk"
	}
	data, err := request(host+path, headers)
	if err != nil {
		return errs.E(errs.ObjectNotFound, "Failed to validate this access token with Linkedin.")
	}
	if data["id"] != nil && utils.S(data["id"]) == utils.S(authData["id"]) {
		return nil
	}
	return errs.E(errs.ObjectNotFound, "Linkedin auth is invalid for this user.")
}
