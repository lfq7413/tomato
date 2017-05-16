package controllers

import (
	"strings"

	"github.com/astaxie/beego"
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/publichtml"
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
		p.Ctx.Output.Header("location", config.VerifyEmailSuccessURL()+"?username="+username)
	} else {
		p.invalidVerification()
	}
}

// ResendVerificationEmail 处理重新发送验证邮件请求
// @router /resend_verification_email [post]
func (p *PublicController) ResendVerificationEmail() {
	username := p.GetString("username")
	if config.TConfig.ServerURL == "" {
		p.missingPublicServerURL()
		return
	}
	if username == "" {
		p.invalid()
		return
	}
	err := rest.ResendVerificationEmail(username)
	if err != nil {
		p.Ctx.Output.SetStatus(302)
		p.Ctx.Output.Header("location", config.LinkSendFailURL())
	} else {
		p.Ctx.Output.SetStatus(302)
		p.Ctx.Output.Header("location", config.LinkSendSuccessURL())
	}
}

// ChangePassword 修改密码页面
// @router /choose_password [get]
func (p *PublicController) ChangePassword() {
	if config.TConfig.ServerURL == "" {
		p.missingPublicServerURL()
		return
	}

	data := strings.Replace(publichtml.ChoosePasswordPage, "PARSE_SERVER_URL", `"`+config.TConfig.ServerURL+`"`, -1)
	p.Ctx.Output.Header("Content-Type", "text/html")
	p.Ctx.Output.Body([]byte(data))
}

// ResetPassword 处理实际的重置密码请求
// @router /request_password_reset [post]
func (p *PublicController) ResetPassword() {
	if config.TConfig.ServerURL == "" {
		p.missingPublicServerURL()
		return
	}

	username := p.GetString("username")
	token := p.GetString("token")
	newPassword := p.GetString("new_password")

	if token == "" || username == "" || newPassword == "" {
		p.invalid()
		return
	}

	err := rest.UpdatePassword(username, token, newPassword)
	if err == nil {
		p.Ctx.Output.SetStatus(302)
		p.Ctx.Output.Header("location", config.PasswordResetSuccessURL()+"?username="+username)
	} else {
		p.Ctx.Output.SetStatus(302)
		location := config.ChoosePasswordURL()
		location += "?token=" + token
		location += "&id=" + config.TConfig.AppID
		location += "&username=" + username
		location += "&error=" + err.Error()
		location += "&app=" + config.TConfig.AppName
		p.Ctx.Output.Header("location", location)
	}
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

	user := rest.CheckResetTokenValidity(username, token)
	if user != nil {
		p.Ctx.Output.SetStatus(302)
		location := config.ChoosePasswordURL()
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
	p.Ctx.Output.Header("Content-Type", "text/html")
	p.Ctx.Output.Body([]byte(publichtml.InvalidLinkPage))
}

// InvalidVerificationLink 无效验证链接页面
// @router /invalid_verification_link [get]
func (p *PublicController) InvalidVerificationLink() {
	data := strings.Replace(publichtml.InvalidVerificationLink, "RESEND_VERIFICATION_URL", config.TConfig.ServerURL+"/apps/resend_verification_email", -1)
	p.Ctx.Output.Header("Content-Type", "text/html")
	p.Ctx.Output.Body([]byte(data))
}

// LinkSendSuccess 发送成功页面
// @router /link_send_success [get]
func (p *PublicController) LinkSendSuccess() {
	p.Ctx.Output.Header("Content-Type", "text/html")
	p.Ctx.Output.Body([]byte(publichtml.LinkSendSuccess))
}

// LinkSendFail 发送失败页面
// @router /link_send_fail [get]
func (p *PublicController) LinkSendFail() {
	p.Ctx.Output.Header("Content-Type", "text/html")
	p.Ctx.Output.Body([]byte(publichtml.LinkSendFail))
}

// PasswordResetSuccess 密码重置成功页面
// @router /password_reset_success [get]
func (p *PublicController) PasswordResetSuccess() {
	p.Ctx.Output.Header("Content-Type", "text/html")
	p.Ctx.Output.Body([]byte(publichtml.PasswordResetSuccessPage))
}

// VerifyEmailSuccess 验证邮箱成功页面
// @router /verify_email_success [get]
func (p *PublicController) VerifyEmailSuccess() {
	p.Ctx.Output.Header("Content-Type", "text/html")
	p.Ctx.Output.Body([]byte(publichtml.VerifyEmailSuccessPage))
}

func (p *PublicController) invalid() {
	p.Ctx.Output.SetStatus(302)
	p.Ctx.Output.Header("location", config.InvalidLinkURL())
}

func (p *PublicController) invalidVerification() {
	username := p.GetString("username")
	if username != "" {
		p.Ctx.Output.SetStatus(302)
		p.Ctx.Output.Header("location", config.InvalidVerificationLinkURL()+"?username="+username)
	} else {
		p.invalid()
	}
}

func (p *PublicController) missingPublicServerURL() {
	p.Ctx.Output.SetStatus(404)
	p.Ctx.Output.Body([]byte("Not found."))
}
