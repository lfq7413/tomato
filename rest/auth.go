package rest

import (
	"time"

	"github.com/lfq7413/tomato/cache"
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
	cachedUser := cache.User.Get(sessionToken)
	if u := utils.M(cachedUser); u != nil {
		return &Auth{
			IsMaster:       false,
			InstallationID: installationID,
			User:           u,
		}, nil
	}
	// 缓存中不存在时，从数据库中查询
	restOptions := types.M{
		"limit":   1,
		"include": "user",
	}
	restWhere := types.M{
		"sessionToken": sessionToken,
	}

	sessionErr := errs.E(errs.InvalidSessionToken, "invalid session token")
	query, err := NewQuery(Master(), "_Session", restWhere, restOptions, nil)
	if err != nil {
		return nil, sessionErr
	}
	response, err := query.Execute()
	if err != nil {
		return nil, sessionErr
	}

	if response == nil || response["results"] == nil {
		return nil, sessionErr
	}
	results := utils.A(response["results"])
	if results == nil || len(results) != 1 {
		return nil, sessionErr
	}
	result := utils.M(results[0])
	if result == nil || result["user"] == nil {
		return nil, sessionErr
	}

	now := time.Now().UTC()
	if result["expiresAt"] == nil {
		return nil, errs.E(errs.InvalidSessionToken, "Session token is expired.")
	}
	expiresAtString := utils.S(result["expiresAt"])
	expiresAt, err := utils.StringtoTime(expiresAtString)
	if err != nil {
		return nil, errs.E(errs.InvalidSessionToken, "Session token is expired.")
	}
	if expiresAt.UnixNano() < now.UnixNano() {
		return nil, errs.E(errs.InvalidSessionToken, "Session token is expired.")
	}

	user := utils.M(result["user"])
	delete(user, "password")
	user["className"] = "_User"
	user["sessionToken"] = sessionToken
	// 写入缓存
	cache.User.Put(sessionToken, user, 0)

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
	if a.User != nil && utils.S(a.User["objectId"]) == objectID {
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
	cachedRoles := cache.Role.Get(a.User["objectId"].(string))
	if cachedRoles != nil {
		a.FetchedRoles = true
		a.UserRoles = cachedRoles.([]string)
		return cachedRoles.([]string)
	}

	users := types.M{
		"__type":    "Pointer",
		"className": "_User",
		"objectId":  a.User["objectId"],
	}
	restWhere := types.M{
		"users": users,
	}
	// 取出当前用户直接对应的所有角色
	// TODO 处理错误，处理结果大于100的情况
	query, err := NewQuery(Master(), "_Role", restWhere, types.M{}, nil)
	if err != nil {
		a.UserRoles = []string{}
		a.FetchedRoles = true
		a.RolePromise = nil
		cache.Role.Put(a.User["objectId"].(string), a.UserRoles, 0)
		return a.UserRoles
	}

	response, err := query.Execute()
	if err != nil || utils.HasResults(response) == false {
		a.UserRoles = []string{}
		a.FetchedRoles = true
		a.RolePromise = nil
		cache.Role.Put(a.User["objectId"].(string), a.UserRoles, 0)
		return a.UserRoles
	}

	results := utils.A(response["results"])
	ids := []string{}
	names := []string{}
	for _, v := range results {
		roleObj := utils.M(v)
		ids = append(ids, utils.S(roleObj["objectId"]))
		names = append(names, utils.S(roleObj["name"]))
	}

	queriedRoles := map[string]bool{} // 记录查询过的 role ，避免多次查询
	roleNames := a.getAllRolesNamesForRoleIds(ids, names, queriedRoles)

	a.UserRoles = []string{}
	for _, v := range roleNames {
		a.UserRoles = append(a.UserRoles, "role:"+v)
	}
	a.FetchedRoles = true
	a.RolePromise = nil

	cache.Role.Put(a.User["objectId"].(string), a.UserRoles, 0)
	return a.UserRoles
}

// getAllRolesNamesForRoleIds 取出角色 id 对应的父角色
func (a *Auth) getAllRolesNamesForRoleIds(roleIDs, names []string, queriedRoles map[string]bool) []string {
	if names == nil {
		names = []string{}
	}
	if queriedRoles == nil {
		queriedRoles = map[string]bool{}
	}
	ins := types.S{}
	for _, roleID := range roleIDs {
		if _, ok := queriedRoles[roleID]; ok {
			continue
		}
		// 标记该 roleID 已经获取过一次父角色了
		queriedRoles[roleID] = true
		object := types.M{
			"__type":    "Pointer",
			"className": "_Role",
			"objectId":  roleID,
		}
		ins = append(ins, object)
	}

	// 已经没有待获取父角色的 roleID，返回 names
	if len(ins) == 0 {
		return names
	}

	restWhere := types.M{}
	if len(ins) == 1 {
		restWhere["roles"] = ins[0]
	} else {
		restWhere["roles"] = types.M{"$in": ins}
	}

	query, err := NewQuery(Master(), "_Role", restWhere, types.M{}, nil)
	if err != nil {
		return names
	}

	// 未找到角色
	response, err := query.Execute()
	if err != nil || utils.HasResults(response) == false {
		return names
	}

	results := utils.A(response["results"])
	ids := []string{}
	pnames := []string{}
	for _, v := range results {
		roleObj := utils.M(v)
		if roleObj == nil {
			continue
		}
		ids = append(ids, utils.S(roleObj["objectId"]))
		pnames = append(pnames, utils.S(roleObj["name"]))
	}

	// 存储找到的角色名
	names = append(names, pnames...)

	// 继续查找最新角色的父角色
	return a.getAllRolesNamesForRoleIds(ids, names, queriedRoles)
}
