package auth

import "github.com/lfq7413/tomato/rest"

// Auth ...
type Auth struct {
	IsMaster       bool
	InstallationID string
	User           *UserInfo
	UserRoles      []string
	FetchedRoles   bool
	RolePromise    []string
}

// UserInfo ...
type UserInfo struct {
	ID string
}

// Master ...
func Master() *Auth {
	return &Auth{IsMaster: true}
}

// Nobody ...
func Nobody() *Auth {
	return &Auth{IsMaster: false}
}

// GetAuthForSessionToken ...
func GetAuthForSessionToken(sessionToken string, installationID string) *Auth {
	return &Auth{IsMaster: false, InstallationID: installationID}
}

// CouldUpdateUserID ...
func (a *Auth) CouldUpdateUserID(objectID string) bool {
	if a.IsMaster {
		return true
	}
	if a.User != nil && a.User.ID == objectID {
		return true
	}
	return false
}

// GetUserRoles ...
func (a *Auth) GetUserRoles() []string {
	if a.IsMaster || a.User == nil {
		return []string{}
	}
	if a.FetchedRoles {
		return a.UserRoles
	}
	if a.RolePromise != nil {
		return a.RolePromise
	}
	a.RolePromise = a.loadRoles()
	return a.RolePromise
}

func (a *Auth) loadRoles() []string {

	users := map[string]string{"__type": "Pointer", "className": "_User", "objectId": a.User.ID}
	where := map[string]interface{}{"users": users}

	response := rest.Find(Master(), "_Role", where, map[string]interface{}{})
	if response == nil {
		return []string{}
	}
	if response["results"] == nil {
		return []string{}
	}
	var results []interface{}
	if v, ok := response["results"].([]interface{}); ok {
		results = v
	} else {
		return []string{}
	}
	if len(results) == 0 {
		return []string{}
	}

	return []string{}
}
