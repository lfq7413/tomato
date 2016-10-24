package authdatamanager

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

type youdao struct{}

func (a youdao) ValidateAuthData(authData types.M, options types.M) error {
	client := NewOAuth(options)
	client.Host = "http://note.youdao.com"
	client.AuthToken = utils.S(authData["access_token"])
	client.AuthTokenSecret = utils.S(authData["auth_token_secret"])
	data, err := client.Get("/yws/open/user/get.json", nil)
	if err != nil {
		return errs.E(errs.ObjectNotFound, "Failed to validate this access token with Youdao.")
	}
	if data["user"] != nil && utils.S(data["user"]) == utils.S(authData["id"]) {
		return nil
	}
	return errs.E(errs.ObjectNotFound, "Youdao auth is invalid for this user.")
}
