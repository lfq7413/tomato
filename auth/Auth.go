package auth

import "github.com/lfq7413/tomato/rest"
import "github.com/lfq7413/tomato/conv"

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
	if response == nil ||
		response["results"] == nil ||
		conv.SliceInterface(response["results"]) == nil ||
		len(conv.SliceInterface(response["results"])) == 0 {
		a.UserRoles = []string{}
		a.FetchedRoles = true
		a.RolePromise = nil
		return a.UserRoles
	}
	results := conv.SliceInterface(response["results"])

	return []string{}
}
