package controllers

import "github.com/astaxie/beego"

// PublicController 处理密码修改与邮箱验证请求
type PublicController struct {
	beego.Controller
}

// VerifyEmail ...
// @router /verify_email [get]
func (p *PublicController) VerifyEmail() {

}

// ChangePassword ...
// @router /choose_password [get]
func (p *PublicController) ChangePassword() {

}

// ResetPassword ...
// @router /request_password_reset [post]
func (p *PublicController) ResetPassword() {

}

// RequestResetPassword ...
// @router /request_password_reset [get]
func (p *PublicController) RequestResetPassword() {

}
