package auth

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

type baidu struct{}

func (a baidu) ValidateAuthData(authData types.M, options types.M) error {
	// http://developer.baidu.com/ms/oauth
	// 具体接口参考： https://open.weixin.qq.com/cgi-bin/showdocument?action=dir_list&t=resource/res_list&verify=1&id=open1419317853&token=&lang=zh_CN
	host := "https://openapi.baidu.com/rest/2.0/passport/users/"
	path := "isAppUser?access_token=" + utils.S(authData["access_token"]) + "&uid=" + utils.S(authData["id"])
	data, err := request(host+path, nil)
	if err != nil {
		return errs.E(errs.ObjectNotFound, "Failed to validate this access token with Baidu.")
	}
	if result, ok := data["result"].(float64); ok && result == 1 {
		return nil
	}
	return errs.E(errs.ObjectNotFound, "Baidu auth is invalid for this user.")
}
