package config

import (
	"time"

	"github.com/astaxie/beego"
)

// Config ...
type Config struct {
	AppName                          string
	ServerURL                        string
	DatabaseURI                      string
	AppID                            string
	MasterKey                        string
	ClientKey                        string
	JavascriptKey                    string
	DotNetKey                        string
	RestAPIKey                       string
	AllowClientClassCreation         bool
	EnableAnonymousUsers             bool
	VerifyUserEmails                 bool
	FileAdapter                      string
	PushAdapter                      string
	MailAdapter                      string
	LiveQueryClasses                 string
	PublisherType                    string
	PublisherURL                     string
	SessionLength                    int
	RevokeSessionOnPasswordReset     bool
	PreventLoginWithUnverifiedEmail  bool
	EmailVerifyTokenValidityDuration int
	SchemaCacheTTL                   int
	SMTPServer                       string
	MailUsername                     string
	MailPassword                     string
	WebhookKey                       string
}

var (
	// TConfig ...
	TConfig *Config
)

func init() {
	TConfig = &Config{
		AppName:                          "",
		ServerURL:                        "http://127.0.0.1:8080/v1",
		DatabaseURI:                      "192.168.99.100:27017/test",
		AppID:                            "",
		MasterKey:                        "",
		ClientKey:                        "",
		AllowClientClassCreation:         false,
		EnableAnonymousUsers:             true,
		VerifyUserEmails:                 false,
		FileAdapter:                      "disk",
		PushAdapter:                      "tomato",
		MailAdapter:                      "smtp",
		SessionLength:                    31536000,
		RevokeSessionOnPasswordReset:     true,
		PreventLoginWithUnverifiedEmail:  false,
		EmailVerifyTokenValidityDuration: -1,
		SchemaCacheTTL:                   5,
	}

	parseConfig()
}

func parseConfig() {
	TConfig.AppName = beego.AppConfig.String("appname")
	TConfig.ServerURL = beego.AppConfig.String("ServerURL")
	TConfig.DatabaseURI = beego.AppConfig.String("DatabaseURI")
	TConfig.AppID = beego.AppConfig.String("AppID")
	TConfig.MasterKey = beego.AppConfig.String("MasterKey")
	TConfig.ClientKey = beego.AppConfig.String("ClientKey")
	TConfig.JavascriptKey = beego.AppConfig.String("JavascriptKey")
	TConfig.DotNetKey = beego.AppConfig.String("DotNetKey")
	TConfig.RestAPIKey = beego.AppConfig.String("RestAPIKey")
	TConfig.AllowClientClassCreation = beego.AppConfig.DefaultBool("AllowClientClassCreation", false)
	TConfig.EnableAnonymousUsers = beego.AppConfig.DefaultBool("EnableAnonymousUsers", true)
	TConfig.VerifyUserEmails = beego.AppConfig.DefaultBool("VerifyUserEmails", false)
	TConfig.FileAdapter = beego.AppConfig.DefaultString("FileAdapter", "disk")
	TConfig.PushAdapter = beego.AppConfig.DefaultString("PushAdapter", "tomato")
	TConfig.MailAdapter = beego.AppConfig.DefaultString("MailAdapter", "smtp")

	// LiveQueryClasses 支持的类列表，格式： classeA|classeB|classeC
	TConfig.LiveQueryClasses = beego.AppConfig.String("LiveQueryClasses")
	TConfig.PublisherType = beego.AppConfig.String("PublisherType")
	TConfig.PublisherURL = beego.AppConfig.String("PublisherURL")

	TConfig.SessionLength = beego.AppConfig.DefaultInt("SessionLength", 31536000)
	TConfig.RevokeSessionOnPasswordReset = beego.AppConfig.DefaultBool("RevokeSessionOnPasswordReset", true)
	TConfig.PreventLoginWithUnverifiedEmail = beego.AppConfig.DefaultBool("PreventLoginWithUnverifiedEmail", false)
	TConfig.EmailVerifyTokenValidityDuration = beego.AppConfig.DefaultInt("EmailVerifyTokenValidityDuration", -1)
	TConfig.SchemaCacheTTL = beego.AppConfig.DefaultInt("SchemaCacheTTL", 5)

	TConfig.SMTPServer = beego.AppConfig.String("SMTPServer")
	TConfig.MailUsername = beego.AppConfig.String("MailUsername")
	TConfig.MailPassword = beego.AppConfig.String("MailPassword")
	TConfig.WebhookKey = beego.AppConfig.String("WebhookKey")
}

// GenerateSessionExpiresAt 获取 Session 过期时间
func GenerateSessionExpiresAt() time.Time {
	expiresAt := time.Now().UTC()
	expiresAt = expiresAt.Add(time.Duration(TConfig.SessionLength) * time.Second)
	return expiresAt
}

// GenerateEmailVerifyTokenExpiresAt 获取 Email 验证 Token 过期时间
func GenerateEmailVerifyTokenExpiresAt() time.Time {
	if TConfig.VerifyUserEmails == false || TConfig.EmailVerifyTokenValidityDuration == -1 {
		return time.Time{}
	}
	expiresAt := time.Now().UTC()
	expiresAt = expiresAt.Add(time.Duration(TConfig.EmailVerifyTokenValidityDuration) * time.Second)
	return expiresAt
}
