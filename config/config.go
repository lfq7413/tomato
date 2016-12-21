package config

import (
	"time"

	"log"

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
	PublisherConfig                  string
	SessionLength                    int
	RevokeSessionOnPasswordReset     bool
	PreventLoginWithUnverifiedEmail  bool
	EmailVerifyTokenValidityDuration int
	SchemaCacheTTL                   int
	SMTPServer                       string
	MailUsername                     string
	MailPassword                     string
	WebhookKey                       string
	EnableAccountLockout             bool
	AccountLockoutThreshold          int
	AccountLockoutDuration           int
	CacheAdapter                     string
	RedisAddress                     string
	RedisPassword                    string
	EnableSingleSchemaCache          bool
	QiniuBucket                      string
	QiniuDomain                      string
	QiniuAccessKey                   string
	QiniuSecretKey                   string
	FileDirectAccess                 bool
	SinaBucket                       string
	SinaDomain                       string
	SinaAccessKey                    string
	SinaSecretKey                    string
	TencentBucket                    string
	TencentAppID                     string
	TencentSecretID                  string
	TencentSecretKey                 string
	PasswordPolicy                   bool
	ResetTokenValidityDuration       int
	ValidatorPattern                 string
	DoNotAllowUsername               bool
	MaxPasswordAge                   int
	MaxPasswordHistory               int
}

var (
	// TConfig ...
	TConfig *Config
)

func init() {
	TConfig = &Config{
		DatabaseURI: "192.168.99.100:27017/test",
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
	TConfig.FileAdapter = beego.AppConfig.DefaultString("FileAdapter", "Disk")
	TConfig.PushAdapter = beego.AppConfig.DefaultString("PushAdapter", "tomato")
	TConfig.MailAdapter = beego.AppConfig.DefaultString("MailAdapter", "smtp")

	// LiveQueryClasses 支持的类列表，格式： classeA|classeB|classeC
	TConfig.LiveQueryClasses = beego.AppConfig.String("LiveQueryClasses")
	TConfig.PublisherType = beego.AppConfig.String("PublisherType")
	TConfig.PublisherURL = beego.AppConfig.String("PublisherURL")
	TConfig.PublisherConfig = beego.AppConfig.String("PublisherConfig")

	TConfig.SessionLength = beego.AppConfig.DefaultInt("SessionLength", 31536000)
	TConfig.RevokeSessionOnPasswordReset = beego.AppConfig.DefaultBool("RevokeSessionOnPasswordReset", true)
	TConfig.PreventLoginWithUnverifiedEmail = beego.AppConfig.DefaultBool("PreventLoginWithUnverifiedEmail", false)
	TConfig.EmailVerifyTokenValidityDuration = beego.AppConfig.DefaultInt("EmailVerifyTokenValidityDuration", 0)
	TConfig.SchemaCacheTTL = beego.AppConfig.DefaultInt("SchemaCacheTTL", 5)

	TConfig.SMTPServer = beego.AppConfig.String("SMTPServer")
	TConfig.MailUsername = beego.AppConfig.String("MailUsername")
	TConfig.MailPassword = beego.AppConfig.String("MailPassword")
	TConfig.WebhookKey = beego.AppConfig.String("WebhookKey")

	TConfig.EnableAccountLockout = beego.AppConfig.DefaultBool("EnableAccountLockout", false)
	TConfig.AccountLockoutThreshold = beego.AppConfig.DefaultInt("AccountLockoutThreshold", 10)
	TConfig.AccountLockoutDuration = beego.AppConfig.DefaultInt("AccountLockoutDuration", 3)

	TConfig.CacheAdapter = beego.AppConfig.DefaultString("CacheAdapter", "InMemory")
	TConfig.RedisAddress = beego.AppConfig.String("RedisAddress")
	TConfig.RedisPassword = beego.AppConfig.String("RedisPassword")

	TConfig.EnableSingleSchemaCache = beego.AppConfig.DefaultBool("EnableSingleSchemaCache", false)

	TConfig.QiniuBucket = beego.AppConfig.String("QiniuBucket")
	TConfig.QiniuDomain = beego.AppConfig.String("QiniuDomain")
	TConfig.QiniuAccessKey = beego.AppConfig.String("QiniuAccessKey")
	TConfig.QiniuSecretKey = beego.AppConfig.String("QiniuSecretKey")
	TConfig.FileDirectAccess = beego.AppConfig.DefaultBool("FileDirectAccess", true)

	TConfig.SinaBucket = beego.AppConfig.String("SinaBucket")
	TConfig.SinaDomain = beego.AppConfig.String("SinaDomain")
	TConfig.SinaAccessKey = beego.AppConfig.String("SinaAccessKey")
	TConfig.SinaSecretKey = beego.AppConfig.String("SinaSecretKey")

	TConfig.TencentAppID = beego.AppConfig.String("TencentAppID")
	TConfig.TencentBucket = beego.AppConfig.String("TencentBucket")
	TConfig.TencentSecretID = beego.AppConfig.String("TencentSecretID")
	TConfig.TencentSecretKey = beego.AppConfig.String("TencentSecretKey")

	TConfig.PasswordPolicy = beego.AppConfig.DefaultBool("PasswordPolicy", false)
	TConfig.ResetTokenValidityDuration = beego.AppConfig.DefaultInt("ResetTokenValidityDuration", 0)
	TConfig.ValidatorPattern = beego.AppConfig.String("ValidatorPattern")
	TConfig.DoNotAllowUsername = beego.AppConfig.DefaultBool("DoNotAllowUsername", false)
	TConfig.MaxPasswordAge = beego.AppConfig.DefaultInt("MaxPasswordAge", 0)
	TConfig.MaxPasswordHistory = beego.AppConfig.DefaultInt("MaxPasswordHistory", 0)
}

// Validate 校验用户参数合法性
func Validate() {
	validateApplicationConfiguration()
	validateFileConfiguration()
	validatePushConfiguration()
	validateMailConfiguration()
	validateLiveQueryConfiguration()
	validateSessionConfiguration()
	validateAccountLockoutPolicy()
}

// validateApplicationConfiguration 校验应用相关参数
func validateApplicationConfiguration() {
	if TConfig.AppName == "" {
		log.Fatalln("AppName is required")
	}
	if TConfig.ServerURL == "" {
		log.Fatalln("ServerURL is required")
	}
	if TConfig.AppID == "" {
		log.Fatalln("AppID is required")
	}
	if TConfig.MasterKey == "" {
		log.Fatalln("MasterKey is required")
	}
	if TConfig.ClientKey == "" && TConfig.JavascriptKey == "" && TConfig.DotNetKey == "" && TConfig.RestAPIKey == "" {
		log.Fatalln("ClientKey or JavascriptKey or DotNetKey or RestAPIKey is required")
	}
}

// validateFileConfiguration 校验文件存储相关参数
func validateFileConfiguration() {
	adapter := TConfig.FileAdapter
	switch adapter {
	case "", "Disk":
	case "GridFS":
	// TODO 校验 MongoDB 配置
	case "Qiniu":
		if TConfig.QiniuDomain == "" && TConfig.QiniuBucket == "" && TConfig.QiniuAccessKey == "" && TConfig.QiniuSecretKey == "" {
			log.Fatalln("QiniuDomain, QiniuBucket, QiniuAccessKey, QiniuSecretKey is required")
		}
	case "Sina":
		if TConfig.SinaDomain == "" && TConfig.SinaBucket == "" && TConfig.SinaAccessKey == "" && TConfig.SinaSecretKey == "" {
			log.Fatalln("SinaDomain, SinaBucket, SinaAccessKey, SinaSecretKey is required")
		}
	case "Tencent":
		if TConfig.TencentAppID == "" && TConfig.TencentBucket == "" && TConfig.TencentSecretID == "" && TConfig.TencentSecretKey == "" {
			log.Fatalln("TencentAppID, TencentBucket, TencentSecretID, TencentSecretKey is required")
		}
	default:
		log.Fatalln("Unsupported FileAdapter")
	}
}

// validatePushConfiguration 校验推送相关参数
func validatePushConfiguration() {
	// TODO
}

// validateMailConfiguration 校验发送邮箱相关参数
func validateMailConfiguration() {
	if TConfig.VerifyUserEmails == false {
		return
	}
	adapter := TConfig.MailAdapter
	switch adapter {
	case "", "smtp":
		if TConfig.SMTPServer == "" {
			log.Fatalln("SMTPServer is required")
		}
		if TConfig.MailUsername == "" {
			log.Fatalln("MailUsername is required")
		}
		if TConfig.MailPassword == "" {
			log.Fatalln("MailPassword is required")
		}
	default:
		log.Fatalln("Unsupported MailAdapter")
	}
	if TConfig.EmailVerifyTokenValidityDuration < 0 {
		log.Fatalln("Email verify token validity duration must be a value greater than 0")
	}
}

// validateLiveQueryConfiguration 校验 LiveQuery 相关参数
func validateLiveQueryConfiguration() {
	t := TConfig.PublisherType
	switch t {
	case "": // 默认为 EventEmitter
	case "Redis":
		if TConfig.PublisherURL == "" {
			log.Fatalln("Redis PublisherURL is required")
		}
	default:
		log.Fatalln("Unsupported LiveQuery PublisherType")
	}
}

// validateSessionConfiguration 校验 Session 有效期
func validateSessionConfiguration() {
	if TConfig.SessionLength <= 0 {
		log.Fatalln("Session length must be a value greater than 0")
	}
}

// validateAccountLockoutPolicy 校验账户锁定规则
func validateAccountLockoutPolicy() {
	if TConfig.EnableAccountLockout == false {
		return
	}
	if TConfig.AccountLockoutDuration < 1 || TConfig.AccountLockoutDuration > 99999 {
		log.Fatalln("Account lockout duration should be greater than 0 and less than 100000")
	}
	if TConfig.AccountLockoutThreshold < 1 || TConfig.AccountLockoutThreshold > 999 {
		log.Fatalln("Account lockout threshold should be an integer greater than 0 and less than 1000")
	}
}

// GenerateSessionExpiresAt 获取 Session 过期时间
func GenerateSessionExpiresAt() time.Time {
	expiresAt := time.Now().UTC()
	expiresAt = expiresAt.Add(time.Duration(TConfig.SessionLength) * time.Second)
	return expiresAt
}

// GenerateEmailVerifyTokenExpiresAt 获取 Email 验证 Token 过期时间
func GenerateEmailVerifyTokenExpiresAt() time.Time {
	if TConfig.VerifyUserEmails == false || TConfig.EmailVerifyTokenValidityDuration <= 0 {
		return time.Time{}
	}
	expiresAt := time.Now().UTC()
	expiresAt = expiresAt.Add(time.Duration(TConfig.EmailVerifyTokenValidityDuration) * time.Second)
	return expiresAt
}

// GeneratePasswordResetTokenExpiresAt 获取 重置密码 验证 Token 过期时间
func GeneratePasswordResetTokenExpiresAt() time.Time {
	if TConfig.PasswordPolicy == false || TConfig.ResetTokenValidityDuration == 0 {
		return time.Time{}
	}
	expiresAt := time.Now().UTC()
	expiresAt = expiresAt.Add(time.Duration(TConfig.ResetTokenValidityDuration) * time.Second)
	return expiresAt
}
