package rest

import (
	"github.com/lfq7413/tomato/utils"
	"gopkg.in/mgo.v2/bson"
)

// Auth 保存当前请求的用户权限信息
type Auth struct {
	IsMaster       bool
	InstallationID string
	User           map[string]interface{}
	UserRoles      []string
	FetchedRoles   bool
	RolePromise    []string
}

// Master 生成 Master 级别用户
func Master() *Auth {
	return &Auth{IsMaster: true}
}

// Nobody 生成空用户
func Nobody() *Auth {
	return &Auth{IsMaster: false}
}

// GetAuthForSessionToken 返回 sessionToken 对应的用户权限信息
func GetAuthForSessionToken(sessionToken string, installationID string) *Auth {
	// 从缓存获取用户信息
	cachedUser := usersCache.get(sessionToken)
	if cachedUser != nil {
		return &Auth{
			IsMaster:       false,
			InstallationID: installationID,
			User:           cachedUser.(map[string]interface{}),
		}
	}
	// 缓存中不存在时，从数据库中查询
	restOptions := bson.M{
		"limit":   1,
		"include": "user",
	}
	restWhere := bson.M{
		"_session_token": sessionToken,
	}
	response := NewQuery(Master(), "_Session", restWhere, restOptions).Execute()

	if response == nil || response["results"] == nil {
		return Nobody()
	}
	results := utils.SliceInterface(response["results"])
	if results == nil || len(results) != 1 {
		return Nobody()
	}
	result := utils.MapInterface(results[0])
	if result == nil || result["user"] == nil {
		return Nobody()
	}

	user := utils.MapInterface(result["user"])
	delete(user, "password")
	user["className"] = "_User"
	user["sessionToken"] = sessionToken
	// 写入缓存
	usersCache.set(sessionToken, user)

	return &Auth{
		IsMaster:       false,
		InstallationID: installationID,
		User:           user,
	}
}

// CouldUpdateUserID Master 与当前用户可进行修改
func (a *Auth) CouldUpdateUserID(objectID string) bool {
	if a.IsMaster {
		return true
	}
	if a.User != nil && a.User["objectId"].(string) == objectID {
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

	users := map[string]interface{}{
		"__type":    "Pointer",
		"className": "_User",
		"objectId":  a.User["objectId"],
	}
	where := map[string]interface{}{
		"users": users,
	}
	// 取出当前用户直接对应的所有角色
	response := Find(Master(), "_Role", where, map[string]interface{}{})
	if utils.HasResults(response) == false {
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
		roleIDs = append(roleIDs, a.getAllRoleNamesForID(v)...)
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
	rolePointer := map[string]interface{}{
		"__type":    "Pointer",
		"className": "_Role",
		"objectId":  roleID,
	}
	where := map[string]interface{}{
		"roles": rolePointer,
	}
	// 取出当前角色对应的直接父角色
	response := Find(Master(), "_Role", where, map[string]interface{}{})
	if utils.HasResults(response) == false {
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
		roleIDs = append(roleIDs, a.getAllRoleNamesForID(v)...)
	}
	return roleIDs
}
