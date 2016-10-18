package authdatamanager

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

type weixin struct{}

func (a weixin) ValidateAuthData(authData types.M, options types.M) error {
	// 具体接口参考： https://open.weixin.qq.com/cgi-bin/showdocument?action=dir_list&t=resource/res_list&verify=1&id=open1419317853&token=&lang=zh_CN
	host := "https://api.weixin.qq.com/sns/"
	path := "auth?access_token=" + utils.S(authData["access_token"]) + "&openid=" + utils.S(authData["id"])
	data, err := request(host+path, nil)
	if err != nil {
		return errs.E(errs.ObjectNotFound, "Failed to validate this access token with Weixin.")
	}
	if code, ok := data["errcode"].(float64); ok && code == 0 {
		if data["errmsg"] != nil && utils.S(data["errmsg"]) == "ok" {
			return nil
		}
	}
	return errs.E(errs.ObjectNotFound, "Weixin auth is invalid for this user.")
}
