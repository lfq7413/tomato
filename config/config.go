package config

import "github.com/astaxie/beego"

// Config ...
type Config struct {
	ServerURL   string
	DatabaseURI string
	AppID       string
	MasterKey   string
	ClientKey   string
}

var (
	// TConfig ...
	TConfig *Config
)

func init() {
	TConfig = &Config{
		ServerURL:   "http://127.0.0.1:8080/v1",
		DatabaseURI: "192.168.99.100:27017/test",
		AppID:       "",
		MasterKey:   "",
		ClientKey:   "",
	}

	ParseConfig()
}

// ParseConfig ...
func ParseConfig() {
	TConfig.ServerURL = beego.AppConfig.String("serverurl")
	TConfig.DatabaseURI = beego.AppConfig.String("databaseuri")
	TConfig.AppID = beego.AppConfig.String("appid")
	TConfig.MasterKey = beego.AppConfig.String("masterkey")
	TConfig.ClientKey = beego.AppConfig.String("clientkey")
}
