package config

import "github.com/astaxie/beego"

// Config ...
type Config struct {
	URL string
}

var (
	// TConfig ...
	TConfig *Config
)

func init() {
	TConfig = &Config{
		URL: "http://127.0.0.1:8080/v1/",
	}

	ParseConfig()
}

// ParseConfig ...
func ParseConfig() {
	TConfig.URL = beego.AppConfig.String("myhttpurl")
}
