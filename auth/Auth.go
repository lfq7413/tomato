package auth

// Auth ...
type Auth struct {
	IsMaster       bool
	InstallationID string
	User           *UserInfo
}

// UserInfo ...
type UserInfo struct{}

// Master ...
func Master() *Auth {
	return &Auth{IsMaster: true}
}

// Nobody ...
func Nobody() *Auth {
	return &Auth{IsMaster: false}
}
