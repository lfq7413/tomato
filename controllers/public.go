package controllers

import (
	"github.com/astaxie/beego"
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/rest"
)

// PublicController 处理密码修改与邮箱验证请求
type PublicController struct {
	beego.Controller
}

// VerifyEmail 处理验证邮箱请求
// 该接口从验证邮件内部发起请求，见 rest.SendVerificationEmail()
// @router /verify_email [get]
func (p *PublicController) VerifyEmail() {
	token := p.GetString("token")
	username := p.GetString("username")

	if config.TConfig.ServerURL == "" {
		p.missingPublicServerURL()
		return
	}

	if token == "" || username == "" {
		p.invalid()
		return
	}

	ok := rest.VerifyEmail(username, token)
	if ok {
		p.Ctx.Output.SetStatus(302)
		p.Ctx.Output.Header("location", config.TConfig.ServerURL+"app/verify_email_success")
	} else {
		p.invalid()
	}
}

// ChangePassword ...
// @router /choose_password [get]
func (p *PublicController) ChangePassword() {

}

// ResetPassword ...
// @router /request_password_reset [post]
func (p *PublicController) ResetPassword() {

}

// RequestResetPassword 处理重置密码请求
// 该接口从重置密码邮件内部发起请求，见 rest.SendPasswordResetEmail()
// @router /request_password_reset [get]
func (p *PublicController) RequestResetPassword() {
	token := p.GetString("token")
	username := p.GetString("username")

	if config.TConfig.ServerURL == "" {
		p.missingPublicServerURL()
		return
	}

	if token == "" || username == "" {
		p.invalid()
		return
	}

	ok := rest.CheckResetTokenValidity(username, token)
	if ok {
		p.Ctx.Output.SetStatus(302)
		location := config.TConfig.ServerURL + "app/choose_password"
		location += "?token=" + token
		location += "&id=" + config.TConfig.AppID
		location += "&username=" + username
		location += "&app=" + config.TConfig.AppName
		p.Ctx.Output.Header("location", location)
	} else {
		p.invalid()
	}
}

// InvalidLink 无效链接页面
// @router /invalid_link [get]
func (p *PublicController) InvalidLink() {

}

// PasswordResetSuccess 密码重置成功页面
// @router /password_reset_success [get]
func (p *PublicController) PasswordResetSuccess() {

}

// VerifyEmailSuccess 验证邮箱成功页面
// @router /verify_email_success [get]
func (p *PublicController) VerifyEmailSuccess() {

}

func (p *PublicController) invalid() {
	p.Ctx.Output.SetStatus(302)
	p.Ctx.Output.Header("location", config.TConfig.ServerURL+"app/invalid_link")
}

func (p *PublicController) missingPublicServerURL() {
	p.Ctx.Output.SetStatus(404)
	p.Ctx.Output.Body([]byte("Not found."))
}
