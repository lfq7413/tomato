package config

import (
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/lfq7413/tomato/livequery"
)

// Config ...
type Config struct {
	AppName                          string
	ServerURL                        string
	DatabaseURI                      string
	AppID                            string
	MasterKey                        string
	ClientKey                        string
	AllowClientClassCreation         bool
	EnableAnonymousUsers             bool
	VerifyUserEmails                 bool
	FileAdapter                      string
	FileDir                          string
	PushAdapter                      string
	MailAdapter                      string
	LiveQuery                        *livequery.LiveQuery
	SessionLength                    int
	RevokeSessionOnPasswordReset     bool
	PreventLoginWithUnverifiedEmail  bool
	EmailVerifyTokenValidityDuration int
	SchemaCacheTTL                   int
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
		FileDir:                          "/Users",
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
	TConfig.ServerURL = beego.AppConfig.String("serverurl")
	TConfig.DatabaseURI = beego.AppConfig.String("databaseuri")
	TConfig.AppID = beego.AppConfig.String("appid")
	TConfig.MasterKey = beego.AppConfig.String("masterkey")
	TConfig.ClientKey = beego.AppConfig.String("clientkey")
	TConfig.AllowClientClassCreation = beego.AppConfig.DefaultBool("allowclientclasscreation", false)
	TConfig.EnableAnonymousUsers = beego.AppConfig.DefaultBool("EnableAnonymousUsers", true)
	TConfig.VerifyUserEmails = beego.AppConfig.DefaultBool("VerifyUserEmails", false)
	TConfig.FileAdapter = beego.AppConfig.DefaultString("FileAdapter", "disk")
	TConfig.FileDir = beego.AppConfig.DefaultString("FileDir", "/Users")
	TConfig.PushAdapter = beego.AppConfig.DefaultString("PushAdapter", "tomato")
	TConfig.MailAdapter = beego.AppConfig.DefaultString("MailAdapter", "smtp")

	// LiveQuery 支持的类列表，格式： classeA|classeB|classeC
	classeNames := beego.AppConfig.DefaultString("LiveQuery", "")
	pubType := beego.AppConfig.DefaultString("PubType", "")
	pubURL := beego.AppConfig.DefaultString("pubURL", "")
	liveQuery := strings.Split(classeNames, "|")
	TConfig.LiveQuery = livequery.NewLiveQuery(liveQuery, pubType, pubURL)

	TConfig.SessionLength = beego.AppConfig.DefaultInt("SessionLength", 31536000)
	TConfig.RevokeSessionOnPasswordReset = beego.AppConfig.DefaultBool("RevokeSessionOnPasswordReset", true)
	TConfig.PreventLoginWithUnverifiedEmail = beego.AppConfig.DefaultBool("PreventLoginWithUnverifiedEmail", false)
	TConfig.EmailVerifyTokenValidityDuration = beego.AppConfig.DefaultInt("EmailVerifyTokenValidityDuration", -1)
	TConfig.SchemaCacheTTL = beego.AppConfig.DefaultInt("SchemaCacheTTL", 5)
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
