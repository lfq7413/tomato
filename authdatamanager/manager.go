package authdatamanager

import (
	"github.com/lfq7413/tomato/config"
)

var providers map[string]Provider

func init() {
	providers = map[string]Provider{
		"anonymous": anonymous{},
	}
}

// ValidateAuthData ...
func ValidateAuthData(provider string, authData map[string]interface{}) error {
	if provider == "anonymous" && config.TConfig.EnableAnonymousUsers == false {
		// TODO 不支持 anonymous
		return nil
	}
	defaultProvider := providers[provider]
	if defaultProvider == nil {
		// TODO 不支持该方式
		return nil
	}

	return defaultProvider.ValidateAuthData(authData)
}

type anonymous struct{}

func (a anonymous) ValidateAuthData(authData map[string]interface{}) error {
	return nil
}

// Provider ...
type Provider interface {
	ValidateAuthData(map[string]interface{}) error
}
