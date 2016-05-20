package rest

import (
	"time"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// Auth 保存当前请求的用户权限信息
type Auth struct {
	IsMaster       bool
	InstallationID string
	User           types.M
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
func GetAuthForSessionToken(sessionToken string, installationID string) (*Auth, error) {
	// 从缓存获取用户信息
	cachedUser := usersCache.get(sessionToken)
	if cachedUser != nil {
		return &Auth{
			IsMaster:       false,
			InstallationID: installationID,
			User:           cachedUser.(map[string]interface{}),
		}, nil
	}
	// 缓存中不存在时，从数据库中查询
	restOptions := types.M{
		"limit":   1,
		"include": "user",
	}
	restWhere := types.M{
		"_session_token": sessionToken,
	}

	query, err := NewQuery(Master(), "_Session", restWhere, restOptions)
	if err != nil {
		return Nobody(), nil
	}
	response, err := query.Execute()
	if err != nil {
		return Nobody(), nil
	}

	if response == nil || response["results"] == nil {
		return Nobody(), nil
	}
	results := utils.SliceInterface(response["results"])
	if results == nil || len(results) != 1 {
		return Nobody(), nil
	}
	result := utils.MapInterface(results[0])
	if result == nil || result["user"] == nil {
		return Nobody(), nil
	}

	now := time.Now().UTC()
	expiresAtString := utils.MapInterface(result["expiresAt"])["iso"].(string)
	expiresAt, _ := utils.StringtoTime(expiresAtString)
	if expiresAt.UnixNano() < now.UnixNano() {
		return nil, errs.E(errs.InvalidSessionToken, "Session token is expired.")
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
	}, nil
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

// GetUserRoles 获取用户所属的角色列表
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

// loadRoles 从数据库加载用户角色列表
func (a *Auth) loadRoles() []string {
	// _Role 表中的 users 字段为 Relation 类型，应该使用 $relatedTo 去查询
	users := types.M{
		"__type":    "Pointer",
		"className": "_User",
		"objectId":  a.User["objectId"],
	}
	relatedTo := types.M{
		"object": users,
		"key":    "users",
	}
	where := types.M{
		"$relatedTo": relatedTo,
	}
	// 取出当前用户直接对应的所有角色
	// TODO 处理错误，处理结果大于100的情况
	response, _ := Find(Master(), "_Role", where, types.M{})
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
	queriedRoles := map[string]bool{} // 记录查询过的 role ，避免多次查询
	for _, v := range roleIDs {
		roleIDs = append(roleIDs, a.getAllRoleNamesForID(v, queriedRoles)...)
	}

	objectID := types.M{
		"$in": roleIDs,
	}
	where = types.M{
		"objectId": objectID,
	}
	// 取出所有角色名称
	// TODO 处理错误，处理结果大于100的情况
	response, _ = Find(Master(), "_Role", where, types.M{})
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

// getAllRoleNamesForID 取出角色 id 对应的父角色
func (a *Auth) getAllRoleNamesForID(roleID string, queriedRoles map[string]bool) []string {
	if _, ok := queriedRoles[roleID]; ok {
		return []string{}
	}
	queriedRoles[roleID] = true

	// _Role 表中的 roles 字段为 Relation 类型，应该使用 $relatedTo 去查询
	rolePointer := types.M{
		"__type":    "Pointer",
		"className": "_Role",
		"objectId":  roleID,
	}
	relatedTo := types.M{
		"object": rolePointer,
		"key":    "roles",
	}
	where := types.M{
		"$relatedTo": relatedTo,
	}
	// 取出当前角色对应的直接父角色
	// TODO 处理错误，处理结果大于100的情况
	response, _ := Find(Master(), "_Role", where, types.M{})
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
		roleIDs = append(roleIDs, a.getAllRoleNamesForID(v, queriedRoles)...)
	}
	return roleIDs
}
