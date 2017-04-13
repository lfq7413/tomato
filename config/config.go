package config

import (
	"time"

	"log"

	"regexp"

	"strings"

	"github.com/astaxie/beego"
)

// Config ...
type Config struct {
	AppName                          string   // 应用名称，必填
	ServerURL                        string   // 服务对外地址，必填
	DatabaseType                     string   // 数据库类型，可选： MongoDB、PostgreSQL
	DatabaseURI                      string   // 数据库地址
	AppID                            string   // 必填
	MasterKey                        string   // 必填
	ClientKey                        string   // 选填
	JavaScriptKey                    string   // 选填
	DotNetKey                        string   // 选填
	RestAPIKey                       string   // 选填
	AllowClientClassCreation         bool     // 是否允许客户端操作不存在的 class ，默认为 fasle 不允许操作
	EnableAnonymousUsers             bool     // 是否支持匿名用户，默认为 true 支持匿名用户
	VerifyUserEmails                 bool     // 是否需要验证用户的 Email ，默认为 false 不需要验证
	EmailVerifyTokenValidityDuration int      // 邮箱验证 Token 有效期，单位为秒，取值大于等于 0 ，默认为 0 表示不设置 Token 有效期
	MailAdapter                      string   // 邮件发送模块，仅在 VerifyUserEmails=true 时需要配置，可选： smtp ，默认为 smtp
	SMTPServer                       string   // SMTP 邮箱服务器地址，仅在 MailAdapter=smtp 时需要配置
	MailUsername                     string   // SMTP 用户名，仅在 MailAdapter=smtp 时需要配置
	MailPassword                     string   // SMTP 密码，仅在 MailAdapter=smtp 时需要配置
	FileAdapter                      string   // 文件存储模块，可选： Disk、GridFS、Qiniu、Sina、Tencent， 默认为 Disk 本地磁盘存储
	FileDirectAccess                 bool     // 是否允许直接访问文件地址，默认为 true 允许直接访问而不是通过 tomato 中转
	QiniuBucket                      string   // 七牛云存储 Bucket ，仅在 FileAdapter=Qiniu 时需要配置
	QiniuDomain                      string   // 七牛云存储 Domain ，仅在 FileAdapter=Qiniu 时需要配置
	QiniuAccessKey                   string   // 七牛云存储 AccessKey ，仅在 FileAdapter=Qiniu 时需要配置
	QiniuSecretKey                   string   // 七牛云存储 SecretKey ，仅在 FileAdapter=Qiniu 时需要配置
	SinaBucket                       string   // 新浪云存储 Bucket ，仅在 FileAdapter=Sina 时需要配置
	SinaDomain                       string   // 新浪云存储 Domain ，仅在 FileAdapter=Sina 时需要配置
	SinaAccessKey                    string   // 新浪云存储 AccessKey ，仅在 FileAdapter=Sina 时需要配置
	SinaSecretKey                    string   // 新浪云存储 SecretKey ，仅在 FileAdapter=Sina 时需要配置
	TencentBucket                    string   // 腾讯云存储 Bucket ，仅在 FileAdapter=Tencent 时需要配置
	TencentAppID                     string   // 腾讯云存储 AppID ，仅在 FileAdapter=Tencent 时需要配置
	TencentSecretID                  string   // 腾讯云存储 SecretID ，仅在 FileAdapter=Tencent 时需要配置
	TencentSecretKey                 string   // 腾讯云存储 SecretKey ，仅在 FileAdapter=Tencent 时需要配置
	PushAdapter                      string   // 推送模块
	LiveQueryClasses                 string   // LiveQuery 支持的 classe ，多个 class 使用 | 隔开，如： classeA|classeB|classeC
	PublisherType                    string   // 发布者类型，可选：Redis ，默认使用自带的 EventEmitter
	PublisherURL                     string   // 发布者地址， PublisherType=Redis 时必填
	PublisherConfig                  string   // 发布者配置信息， PublisherType=Redis 时为 Redis 密码，选填
	SessionLength                    int      // Session 有效期，单位为秒，取值大于 0 ，默认为 31536000 秒，即 1 年
	RevokeSessionOnPasswordReset     bool     // 密码重置后是否清除 Session ，默认为 true 清除 Session
	PreventLoginWithUnverifiedEmail  bool     // 是否阻止未验证邮箱的用户登录，默认为 false 不阻止
	CacheAdapter                     string   // 缓存模块，可选： InMemory、Redis、Null， 默认为 InMemory 使用内存做缓存模块
	RedisAddress                     string   // Redis 地址， CacheAdapter=Redis 时必填
	RedisPassword                    string   // Redis 密码，选填
	SchemaCacheTTL                   int      // Schema 缓存有效期，单位为秒。取值： -1 表示永不过期，0 表示使用 CacheAdapter 自身的有效期，或者大于 0 ，默认为 5 秒
	EnableSingleSchemaCache          bool     // 是否允许缓存唯一一份 SchemaCache ，默认为 false 不允许
	WebhookKey                       string   // 用于云代码鉴权
	EnableAccountLockout             bool     // 是否启用账户锁定规则，默认为 false 不启用
	AccountLockoutThreshold          int      // 锁定账户需要的登录失败次数，取值范围： 1-999 ，默认为 3 次
	AccountLockoutDuration           int      // 锁定账户时长，单位为分钟，取值范围： 1-99999 ，默认为 10 分钟
	PasswordPolicy                   bool     // 是否启用密码规则，默认为 false 不启用
	ResetTokenValidityDuration       int      // 密码重置验证 Token 有效期，单位为秒，取值大于等于 0 ，默认为 0 表示不设置 Token 有效期
	ValidatorPattern                 string   // 校验密码规则的正则表达式
	DoNotAllowUsername               bool     // 是否启用密码中不允许包含用户名，默认为 false 不启用，密码中可包含用户名
	MaxPasswordAge                   int      // 密码的最长使用时间，单位为天，取值大于等于 0 ，默认为 0 表示不设置最长使用时间
	MaxPasswordHistory               int      // 最大密码历史个数，修改的密码不能与密码历史重复，取值范围： 0-20 ，默认为 0 表示不设置密码历史
	UserSensitiveFields              []string // 用户敏感字段，按需删除，多个字段使用 | 删除，如： email|password
	AnalyticsAdapter                 string   // 分析模块，可选：InfluxDB，默认使用空的分析模块
	InfluxDBURL                      string   // InfluxDB 地址，仅在 AnalyticsAdapter=InfluxDB 时需要配置
	InfluxDBUsername                 string   // InfluxDB 用户名，仅在 AnalyticsAdapter=InfluxDB 时需要配置
	InfluxDBPassword                 string   // InfluxDB 密码，仅在 AnalyticsAdapter=InfluxDB 时需要配置
	InfluxDBDatabaseName             string   // InfluxDB 数据库，仅在 AnalyticsAdapter=InfluxDB 时需要配置
	InvalidLink                      string   // 自定义页面地址，无效链接页面
	VerifyEmailSuccess               string   // 自定义页面地址，验证邮箱成功页面
	ChoosePassword                   string   // 自定义页面地址，修改密码页面
	PasswordResetSuccess             string   // 自定义页面地址，密码重置成功页面
	ParseFrameURL                    string   // 自定义页面地址，用于呈现验证 Email 页面和密码重置页面
}

var (
	// TConfig ...
	TConfig *Config
)

func init() {
	TConfig = &Config{
		DatabaseURI:         "192.168.99.100:27017/test",
		UserSensitiveFields: []string{"email"},
	}

	parseConfig()
}

func parseConfig() {
	TConfig.AppName = beego.AppConfig.String("appname")
	TConfig.ServerURL = beego.AppConfig.String("ServerURL")
	TConfig.DatabaseType = beego.AppConfig.String("DatabaseType")
	TConfig.DatabaseURI = beego.AppConfig.String("DatabaseURI")
	TConfig.AppID = beego.AppConfig.String("AppID")
	TConfig.MasterKey = beego.AppConfig.String("MasterKey")
	TConfig.ClientKey = beego.AppConfig.String("ClientKey")
	TConfig.JavaScriptKey = beego.AppConfig.String("JavaScriptKey")
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
	TConfig.AccountLockoutThreshold = beego.AppConfig.DefaultInt("AccountLockoutThreshold", 3)
	TConfig.AccountLockoutDuration = beego.AppConfig.DefaultInt("AccountLockoutDuration", 10)

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

	for _, field := range strings.Split(beego.AppConfig.String("UserSensitiveFields"), "|") {
		TConfig.UserSensitiveFields = append(TConfig.UserSensitiveFields, field)
	}

	TConfig.AnalyticsAdapter = beego.AppConfig.String("AnalyticsAdapter")
	TConfig.InfluxDBURL = beego.AppConfig.String("InfluxDBURL")
	TConfig.InfluxDBUsername = beego.AppConfig.String("InfluxDBUsername")
	TConfig.InfluxDBPassword = beego.AppConfig.String("InfluxDBPassword")
	TConfig.InfluxDBDatabaseName = beego.AppConfig.String("InfluxDBDatabaseName")

	TConfig.InvalidLink = beego.AppConfig.String("InvalidLink")
	TConfig.VerifyEmailSuccess = beego.AppConfig.String("VerifyEmailSuccess")
	TConfig.ChoosePassword = beego.AppConfig.String("ChoosePassword")
	TConfig.PasswordResetSuccess = beego.AppConfig.String("PasswordResetSuccess")
	TConfig.ParseFrameURL = beego.AppConfig.String("ParseFrameURL")
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
	validatePasswordPolicy()
	validateCacheConfiguration()
	validateAnalyticsConfiguration()
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
	if TConfig.ClientKey == "" && TConfig.JavaScriptKey == "" && TConfig.DotNetKey == "" && TConfig.RestAPIKey == "" {
		log.Fatalln("ClientKey or JavaScriptKey or DotNetKey or RestAPIKey is required")
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

// validatePasswordPolicy 校验密码规则
func validatePasswordPolicy() {
	if TConfig.PasswordPolicy == false {
		return
	}
	if TConfig.ResetTokenValidityDuration < 0 {
		log.Fatalln("ResetTokenValidityDuration must be a positive number")
	}
	if TConfig.ValidatorPattern != "" {
		_, err := regexp.Compile(TConfig.ValidatorPattern)
		if err != nil {
			log.Fatalln("ValidatorPattern must be a RegExp")
		}
	}
	if TConfig.MaxPasswordAge < 0 {
		log.Fatalln("MaxPasswordAge must be a positive number")
	}
	if TConfig.MaxPasswordHistory < 0 || TConfig.MaxPasswordHistory > 20 {
		log.Fatalln("MaxPasswordHistory must be an integer ranging 0 - 20")
	}
}

// validateCacheConfiguration 校验缓存相关参数
func validateCacheConfiguration() {
	adapter := TConfig.CacheAdapter
	switch adapter {
	case "", "InMemory", "Null":
	case "Redis":
		if TConfig.RedisAddress == "" {
			log.Fatalln("RedisAddress is required")
		}
	default:
		log.Fatalln("Unsupported CacheAdapter")
	}
	if TConfig.SchemaCacheTTL < -1 {
		log.Fatalln("SchemaCacheTTL should be -1 or 0 or an integer greater than 0")
	}
}

// validateAnalyticsConfiguration 校验分析模块相关参数
func validateAnalyticsConfiguration() {
	adapter := TConfig.AnalyticsAdapter
	switch adapter {
	case "InfluxDB":
		if TConfig.InfluxDBURL == "" {
			log.Fatalln("InfluxDBURL is required")
		}
		if TConfig.InfluxDBUsername == "" {
			log.Fatalln("InfluxDBUsername is required")
		}
		if TConfig.InfluxDBPassword == "" {
			log.Fatalln("InfluxDBPassword is required")
		}
		if TConfig.InfluxDBDatabaseName == "" {
			log.Fatalln("InfluxDBDatabaseName is required")
		}
	case "":
		// 默认使用空实现
	default:
		log.Fatalln("Unsupported AnalyticsAdapter")
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

// InvalidLinkURL ...
func InvalidLinkURL() string {
	if TConfig.InvalidLink != "" {
		return TConfig.InvalidLink
	}
	return TConfig.ServerURL + `/apps/invalid_link`
}

// VerifyEmailSuccessURL ...
func VerifyEmailSuccessURL() string {
	if TConfig.VerifyEmailSuccess != "" {
		return TConfig.VerifyEmailSuccess
	}
	return TConfig.ServerURL + `/apps/verify_email_success`
}

// ChoosePasswordURL ...
func ChoosePasswordURL() string {
	if TConfig.ChoosePassword != "" {
		return TConfig.ChoosePassword
	}
	return TConfig.ServerURL + `/apps/choose_password`
}

// RequestResetPasswordURL ...
func RequestResetPasswordURL() string {
	return TConfig.ServerURL + `/apps/request_password_reset`
}

// PasswordResetSuccessURL ...
func PasswordResetSuccessURL() string {
	if TConfig.PasswordResetSuccess != "" {
		return TConfig.PasswordResetSuccess
	}
	return TConfig.ServerURL + `/apps/password_reset_success`
}

// ParseFrameURL ...
func ParseFrameURL() string {
	return TConfig.ParseFrameURL
}

// VerifyEmailURL ...
func VerifyEmailURL() string {
	return TConfig.ServerURL + `/apps/verify_email`
}
