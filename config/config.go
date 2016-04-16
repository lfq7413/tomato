package config

import (
	"github.com/astaxie/beego"
)

// Config ...
type Config struct {
	AppName                  string
	ServerURL                string
	DatabaseURI              string
	AppID                    string
	MasterKey                string
	ClientKey                string
	AllowClientClassCreation bool
	EnableAnonymousUsers     bool
	VerifyUserEmails         bool
	FileAdapter              string
	FileDir                  string
	PushAdapter              string
	MailAdapter              string
	VerifyEmailURL           string
}

var (
	// TConfig ...
	TConfig *Config
)

func init() {
	TConfig = &Config{
		AppName:                  "",
		ServerURL:                "http://127.0.0.1:8080/v1",
		DatabaseURI:              "192.168.99.100:27017/test",
		AppID:                    "",
		MasterKey:                "",
		ClientKey:                "",
		AllowClientClassCreation: false,
		EnableAnonymousUsers:     true,
		VerifyUserEmails:         false,
		FileAdapter:              "disk",
		FileDir:                  "/Users",
		PushAdapter:              "tomato",
		MailAdapter:              "smtp",
		VerifyEmailURL:           "http://127.0.0.1:8080/v1/app/verify_email",
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
	TConfig.VerifyEmailURL = beego.AppConfig.String("VerifyEmailURL")
}
