package authdatamanager

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

type qq struct{}

func (a qq) ValidateAuthData(authData types.M, options types.M) error {
	// 具体接口参考： http://wiki.open.qq.com/wiki/website/%E8%8E%B7%E5%8F%96%E7%94%A8%E6%88%B7OpenID_OAuth2.0
	host := "https://graph.qq.com/oauth2.0/"
	path := "me?access_token=" + utils.S(authData["access_token"])
	data, err := requestQQ(host+path, nil)
	if err != nil {
		return errs.E(errs.ObjectNotFound, "Failed to validate this access token with QQ.")
	}
	if data["uid"] != nil && utils.S(data["uid"]) == utils.S(authData["id"]) {
		return nil
	}
	return errs.E(errs.ObjectNotFound, "QQ auth is invalid for this user.")
}
