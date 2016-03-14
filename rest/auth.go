package rest

import (
	"github.com/lfq7413/tomato/utils"
)

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

	users := map[string]string{
		"__type":    "Pointer",
		"className": "_User",
		"objectId":  a.User.ID,
	}
	where := map[string]interface{}{
		"users": users,
	}
	// 取出当前用户直接对应的所有角色
	response := Find(Master(), "_Role", where, map[string]interface{}{})
	if response == nil ||
		response["results"] == nil ||
		utils.SliceInterface(response["results"]) == nil ||
		len(utils.SliceInterface(response["results"])) == 0 {
		a.UserRoles = []string{}
		a.FetchedRoles = true
		a.RolePromise = nil
		return a.UserRoles
	}
	results := utils.SliceInterface(response["results"])
	roleIDs := []string{}
	for _, v := range results {
		roleObj := utils.MapInterface(v)
		roleIDs = append(roleIDs, utils.String(roleObj["objectId"]))
	}
	// 取出角色对应的父角色
	for _, v := range roleIDs {
		roleIDs = utils.AppendString(roleIDs, a.getAllRoleNamesForID(v))
	}

	objectID := map[string]interface{}{
		"$in": roleIDs,
	}
	where = map[string]interface{}{
		"objectId": objectID,
	}
	// 取出所有角色名称
	response = Find(Master(), "_Role", where, map[string]interface{}{})
	results = utils.SliceInterface(response["results"])
	a.UserRoles = []string{}
	for _, v := range results {
		roleObj := utils.MapInterface(v)
		a.UserRoles = append(a.UserRoles, "role:"+utils.String(roleObj["name"]))
	}
	a.FetchedRoles = true
	a.RolePromise = nil

	return a.UserRoles
}

func (a *Auth) getAllRoleNamesForID(roleID string) []string {
	rolePointer := map[string]string{
		"__type":    "Pointer",
		"className": "_Role",
		"objectId":  roleID,
	}
	where := map[string]interface{}{
		"roles": rolePointer,
	}
	// 取出当前角色对应的直接父角色
	response := Find(Master(), "_Role", where, map[string]interface{}{})
	if response == nil ||
		response["results"] == nil ||
		utils.SliceInterface(response["results"]) == nil ||
		len(utils.SliceInterface(response["results"])) == 0 {
		return []string{}
	}
	results := utils.SliceInterface(response["results"])
	roleIDs := []string{}
	for _, v := range results {
		roleObj := utils.MapInterface(v)
		roleIDs = append(roleIDs, utils.String(roleObj["objectId"]))
	}
	// 递归取出角色对应的父角色
	for _, v := range roleIDs {
		roleIDs = utils.AppendString(roleIDs, a.getAllRoleNamesForID(v))
	}
	return roleIDs
}
