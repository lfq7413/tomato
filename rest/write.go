package rest

import (
	"strings"
	"time"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/utils"
)

// Write ...
type Write struct {
	auth         *Auth
	className    string
	query        map[string]interface{}
	data         map[string]interface{}
	originalData map[string]interface{}
	storage      map[string]interface{}
	runOptions   map[string]interface{}
	response     map[string]interface{}
	updatedAt    time.Time
}

// NewWrite 可用于 create 和 update ， create 时 query 为 nil
func NewWrite(
	auth *Auth,
	className string,
	query map[string]interface{},
	data map[string]interface{},
	originalData map[string]interface{},
) *Write {
	if query == nil && data["objectId"] != nil {
		// TODO objectId 无效
	}
	write := &Write{
		auth:         auth,
		className:    className,
		query:        query,
		data:         data,
		originalData: originalData,
		storage:      map[string]interface{}{},
		runOptions:   map[string]interface{}{},
		response:     nil,
		updatedAt:    time.Now().UTC(),
	}
	return write
}

// Execute ...
func (w *Write) Execute() map[string]interface{} {
	w.getUserAndRoleACL()
	w.validateClientClassCreation()
	w.validateSchema()
	w.handleInstallation()
	w.handleSession()
	w.validateAuthData()
	w.runBeforeTrigger()
	w.setRequiredFieldsIfNeeded()
	w.transformUser()
	w.expandFilesForExistingObjects()
	w.runDatabaseOperation()
	w.handleFollowup()
	w.runAfterTrigger()
	return w.response
}

func (w *Write) getUserAndRoleACL() error {
	if w.auth.IsMaster {
		return nil
	}
	w.runOptions["acl"] = []string{"*"}
	if w.auth.User != nil {
		roles := w.auth.GetUserRoles()
		roles = append(roles, w.auth.User.ID)
		if v, ok := w.runOptions["acl"].([]string); ok {
			v = utils.AppendString(v, roles)
		}
	}
	return nil
}

func (w *Write) validateClientClassCreation() error {
	sysClass := []string{"_User", "_Installation", "_Role", "_Session", "_Product"}
	if config.TConfig.AllowClientClassCreation {
		return nil
	}
	if w.auth.IsMaster {
		return nil
	}
	for _, v := range sysClass {
		if v == w.className {
			return nil
		}
	}
	if orm.CollectionExists(w.className) {
		return nil
	}
	// TODO 无法操作不存在的表
	return nil
}

func (w *Write) validateSchema() error {
	// TODO 判断是否可以进行操作
	orm.ValidateObject(w.className, w.data, w.query, w.runOptions)
	return nil
}

func (w *Write) handleInstallation() error {
	if w.response != nil || w.className != "_Installation" {
		return nil
	}

	if w.query == nil && w.data["deviceToken"] == nil && w.data["installationId"] == nil {
		// TODO 设备 id 不能为空
		return nil
	}

	if w.query == nil && w.data["deviceType"] == nil {
		// TODO 设备类型不能为空
		return nil
	}

	if w.data["deviceToken"] != nil && len(utils.String(w.data["deviceToken"])) == 64 {
		w.data["deviceToken"] = strings.ToLower(utils.String(w.data["deviceToken"]))
	}

	if w.data["installationId"] != nil {
		w.data["installationId"] = strings.ToLower(utils.String(w.data["installationId"]))
	}

	var idMatch map[string]interface{}
	var deviceTokenMatches []interface{}

	if w.query != nil && w.query["objectId"] != nil {
		results := orm.Find("_Installation", map[string]interface{}{"objectId": w.query["objectId"]}, map[string]interface{}{})
		if results == nil || len(results) == 0 {
			// TODO 更新对象未找到
			return nil
		}
		idMatch = utils.MapInterface(results[0])
		if w.data["installationId"] != nil && idMatch["installationId"] != nil &&
			w.data["installationId"] != idMatch["installationId"] {
			//TODO installationId 不能修改
			return nil
		}
		if w.data["deviceToken"] != nil && idMatch["deviceToken"] != nil &&
			w.data["deviceToken"] != idMatch["deviceToken"] &&
			w.data["installationId"] == nil && idMatch["installationId"] == nil {
			//TODO deviceToken 不能修改
			return nil
		}
		if w.data["deviceType"] != nil && idMatch["deviceType"] != nil &&
			w.data["deviceType"] != idMatch["deviceType"] {
			//TODO deviceType 不能修改
			return nil
		}
	}

	idMatch = nil
	if w.data["installationId"] != nil {
		results := orm.Find("_Installation", map[string]interface{}{"installationId": w.data["installationId"]}, map[string]interface{}{})
		if results != nil && len(results) > 0 {
			idMatch = utils.MapInterface(results[0])
		}
	}
	if w.data["deviceToken"] != nil {
		results := orm.Find("_Installation", map[string]interface{}{"deviceToken": w.data["deviceToken"]}, map[string]interface{}{})
		if results != nil {
			deviceTokenMatches = results
		}
	}

	var objID string
	if idMatch == nil {
		if deviceTokenMatches == nil || len(deviceTokenMatches) == 0 {
			objID = ""
		} else if len(deviceTokenMatches) == 1 &&
			(utils.MapInterface(deviceTokenMatches[0])["installationId"] == nil || w.data["installationId"] == nil) {
			objID = utils.String(utils.MapInterface(deviceTokenMatches[0])["objectId"])
		} else if w.data["installationId"] == nil {
			// TODO 当有多个 deviceToken 时，必须指定 installationId
			return nil
		} else {
			// 清理多余数据
			installationID := map[string]interface{}{
				"$ne": w.data["installationId"],
			}
			delQuery := map[string]interface{}{
				"deviceToken":    w.data["deviceToken"],
				"installationId": installationID,
			}
			if w.data["appIdentifier"] != nil {
				delQuery["appIdentifier"] = w.data["appIdentifier"]
			}
			orm.Destroy("_Installation", delQuery)
			objID = ""
		}
	} else {
		if deviceTokenMatches != nil && len(deviceTokenMatches) == 1 &&
			utils.MapInterface(deviceTokenMatches[0])["installationId"] == nil {
			// 合并
			delQuery := map[string]interface{}{
				"objectId": idMatch["objectId"],
			}
			orm.Destroy("_Installation", delQuery)
			objID = utils.String(utils.MapInterface(deviceTokenMatches[0])["objectId"])
		} else {
			if w.data["deviceToken"] != nil && idMatch["deviceToken"] != w.data["deviceToken"] {
				// 清理多余数据
				installationID := map[string]interface{}{
					"$ne": w.data["installationId"],
				}
				delQuery := map[string]interface{}{
					"deviceToken":    w.data["deviceToken"],
					"installationId": installationID,
				}
				if w.data["appIdentifier"] != nil {
					delQuery["appIdentifier"] = w.data["appIdentifier"]
				}
				orm.Destroy("_Installation", delQuery)
			}
			objID = utils.String(idMatch["objectId"])
		}
	}
	if objID != "" {
		w.query = map[string]interface{}{
			"objectId": objID,
		}
		delete(w.data, "objectId")
		delete(w.data, "createdAt")
	}
	// TODO Validate ops (add/remove on channels, $inc on badge, etc.)

	return nil
}

func (w *Write) handleSession() error {
	if w.response != nil || w.className != "_Session" {
		return nil
	}

	if w.auth.User == nil && w.auth.IsMaster == false {
		// TODO 需要 Session token
		return nil
	}

	if w.data["ACL"] != nil {
		// TODO Session 不能设置 ACL
		return nil
	}

	if w.query == nil && w.auth.IsMaster == false {
		token := "r:" + utils.CreateToken()
		expiresAt := time.Now().UTC()
		expiresAt.AddDate(1, 0, 0)
		user := map[string]interface{}{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  w.auth.User.ID,
		}
		createdWith := map[string]interface{}{
			"action": "create",
		}
		sessionData := map[string]interface{}{
			"sessionToken": token,
			"user":         user,
			"createdWith":  createdWith,
			"restricted":   true,
			"expiresAt":    expiresAt,
		}
		for k, v := range w.data {
			if k == "objectId" {
				continue
			}
			sessionData[k] = v
		}
		results := NewWrite(Master(), "_Session", nil, sessionData, map[string]interface{}{}).Execute()
		if results["response"] == nil {
			// TODO 创建 Session 失败
			return nil
		}
		sessionData["objectId"] = utils.MapInterface(results["response"])["objectId"]
		w.response = map[string]interface{}{
			"status":   201,
			"location": results["location"],
			"response": sessionData,
		}
	}

	return nil
}

func (w *Write) validateAuthData() error {
	if w.className != "_User" {
		return nil
	}

	if w.query == nil && w.data["authData"] == nil {
		if utils.String(w.data["username"]) == "" {
			// TODO 没有设置 username
			return nil
		}
		if utils.String(w.data["password"]) == "" {
			// TODO 没有设置 password
			return nil
		}
	}

	if w.data["authData"] == nil || len(utils.MapInterface(w.data["authData"])) == 0 {
		return nil
	}

	authData := utils.MapInterface(w.data["authData"])
	canHandleAuthData := true

	for _, v := range authData {
		providerAuthData := utils.MapInterface(v)
		hasToken := (providerAuthData != nil && providerAuthData["id"] != nil)
		canHandleAuthData = (canHandleAuthData && (hasToken || providerAuthData == nil))
	}
	if canHandleAuthData {
		return w.handleAuthData(authData)
	}
	// TODO 这个 authentication 不支持
	return nil
}

func (w *Write) handleAuthData(authData map[string]interface{}) error {
	w.handleAuthDataValidation(authData)
	results := w.findUsersWithAuthData(authData)
	if results != nil && len(results) > 1 {
		// TODO auth 已经被使用
		return nil
	}

	keys := []string{}
	for k := range authData {
		keys = append(keys, k)
	}
	w.storage["authProvider"] = strings.Join(keys, ",")

	if results == nil || len(results) == 0 {
		w.data["username"] = utils.CreateToken()
	} else if w.query == nil {
		// 登录
		user := utils.MapInterface(results[0])
		delete(user, "password")
		w.response = map[string]interface{}{
			"response": user,
			"location": w.location(),
		}
		w.data["objectId"] = user["objectId"]
	} else if w.query != nil && w.query["objectId"] != nil {
		// 更新
		user := utils.MapInterface(results[0])
		if utils.String(user["objectId"]) != utils.String(w.query["objectId"]) {
			// auth 已经被使用
			return nil
		}
	}

	return nil
}

func (w *Write) handleAuthDataValidation(authData map[string]interface{}) error {
	return nil
}

func (w *Write) findUsersWithAuthData(authData map[string]interface{}) []interface{} {
	return nil
}

func (w *Write) runBeforeTrigger() error {
	return nil
}

func (w *Write) setRequiredFieldsIfNeeded() error {
	return nil
}

func (w *Write) transformUser() error {
	return nil
}

func (w *Write) expandFilesForExistingObjects() error {
	return nil
}

func (w *Write) runDatabaseOperation() error {
	return nil
}

func (w *Write) handleFollowup() error {
	return nil
}

func (w *Write) runAfterTrigger() error {
	return nil
}

func (w *Write) location() string {
	var middle string
	if w.className == "_User" {
		middle = "/users/"
	} else {
		middle = "/classes/" + w.className + "/"
	}
	return config.TConfig.ServerURL + middle + utils.String(w.data["objectId"])
}
