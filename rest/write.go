package rest

import (
	"strings"
	"time"

	"github.com/lfq7413/tomato/authdatamanager"
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
		expiresAt = expiresAt.AddDate(1, 0, 0)
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
			// TODO auth 已经被使用
			return nil
		}
	}

	return nil
}

func (w *Write) handleAuthDataValidation(authData map[string]interface{}) error {
	for k, v := range authData {
		if v == nil {
			continue
		}
		result := authdatamanager.ValidateAuthData(k, utils.MapInterface(v))
		if result != nil {
			// 验证出现问题
			return nil
		}
	}

	return nil
}

func (w *Write) findUsersWithAuthData(authData map[string]interface{}) []interface{} {
	query := []interface{}{}
	for k, v := range authData {
		if v == nil {
			continue
		}
		key := "authData." + k + ".id"
		provider := utils.MapInterface(v)
		q := map[string]interface{}{
			key: provider["id"],
		}
		query = append(query, q)
	}

	findPromise := []interface{}{}
	if len(query) > 0 {
		where := map[string]interface{}{
			"$or": query,
		}
		findPromise = orm.Find(w.className, where, map[string]interface{}{})
	}

	return findPromise
}

func (w *Write) runBeforeTrigger() error {
	if w.response != nil {
		return nil
	}
	if TriggerExists(TypeBeforeSave, w.className) == false {
		return nil
	}

	updatedObject := map[string]interface{}{}
	if w.query != nil && w.query["objectId"] != nil {
		// 如果是更新，则把原始数据添加进来
		for k, v := range w.originalData {
			updatedObject[k] = v
		}
	}
	// 把需要更新的数据添加进来
	for k, v := range w.data {
		updatedObject[k] = v
	}

	response := RunTrigger(TypeBeforeSave, w.className, w.auth, updatedObject, w.originalData)
	if response != nil && response["object"] != nil {
		w.data = utils.MapInterface(response["object"])
		w.storage["changedByTrigger"] = true
		if w.query != nil && w.query["objectId"] != nil {
			delete(w.data, "objectId")
		}
	}

	return nil
}

func (w *Write) setRequiredFieldsIfNeeded() error {
	if w.data != nil {
		w.data["updatedAt"] = w.updatedAt
		if w.query == nil {
			w.data["createdAt"] = w.updatedAt

			if w.data["objectId"] == nil {
				w.data["objectId"] = utils.CreateObjectID()
			}
		}
	}

	return nil
}

func (w *Write) transformUser() error {
	if w.className != "_User" {
		return nil
	}

	// 如果是创建用户，则先创建 token
	if w.query == nil {
		token := "r:" + utils.CreateToken()
		w.storage["token"] = token
		expiresAt := time.Now().UTC()
		expiresAt = expiresAt.AddDate(1, 0, 0)
		user := map[string]interface{}{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  w.objectID(),
		}
		var authProvider interface{}
		if w.storage["authProvider"] != nil {
			authProvider = w.storage["authProvider"]
		} else {
			authProvider = "password"
		}
		createdWith := map[string]interface{}{
			"action":       "login",
			"authProvider": authProvider,
		}
		sessionData := map[string]interface{}{
			"sessionToken":   token,
			"user":           user,
			"createdWith":    createdWith,
			"restricted":     false,
			"installationId": w.data["installationId"],
			"expiresAt":      expiresAt,
		}
		if w.response != nil && w.response["response"] != nil {
			response := utils.MapInterface(w.response["response"])
			response["sessionToken"] = token
		}
		// TODO 处理创建结果
		NewWrite(Master(), "_Session", nil, sessionData, nil).Execute()
	}

	// 处理密码，计算 sha256
	if w.data["password"] == nil {

	} else {
		if w.query != nil && w.auth.IsMaster == false {
			w.storage["clearSessions"] = true
		}
		w.data["_hashed_password"] = utils.Hash(utils.String(w.data["password"]))
		delete(w.data, "password")
	}

	// 处理用户名，检测用户名是否唯一
	if w.data["username"] == nil {
		if w.query == nil {
			w.data["username"] = utils.CreateObjectID()
		}
	} else {
		objectID := map[string]interface{}{
			"$ne": w.objectID(),
		}
		where := map[string]interface{}{
			"username": w.data["username"],
			"objectId": objectID,
		}
		option := map[string]interface{}{
			"limit": 1,
		}
		results := orm.Find(w.className, where, option)
		if len(results) > 0 {
			// TODO 用户已经存在
			return nil
		}
	}

	// 处理 email ，检测合法性、检测是否唯一
	if w.data["email"] == nil {

	} else {
		if utils.IsEmail(utils.String(w.data["email"])) == false {
			// TODO email 不合法
			return nil
		}
		objectID := map[string]interface{}{
			"$ne": w.objectID(),
		}
		where := map[string]interface{}{
			"email":    w.data["email"],
			"objectId": objectID,
		}
		option := map[string]interface{}{
			"limit": 1,
		}
		results := orm.Find(w.className, where, option)
		if len(results) > 0 {
			// TODO email 已经存在
			return nil
		}

		// 更新 email ，需要发送验证邮件
		w.storage["sendVerificationEmail"] = true
		SetEmailVerifyToken(w.data)
	}

	return nil
}

func (w *Write) expandFilesForExistingObjects() error {
	if w.response != nil && w.response["response"] != nil {
		// TODO 展开文件对象
	}

	return nil
}

func (w *Write) runDatabaseOperation() error {
	if w.response != nil {
		return nil
	}

	if w.className == "_User" && w.query != nil &&
		w.auth.CouldUpdateUserID(utils.String(w.query["objectId"])) == false {
		//TODO 不能更新当前用户
		return nil
	}

	if w.className == "_Product" && w.data["download"] != nil {
		download := utils.MapInterface(w.data["download"])
		w.data["downloadName"] = download["name"]
	}

	if w.data["ACL"] != nil && utils.MapInterface(w.data["ACL"])["*unresolved"] != nil {
		// TODO 无效的 ACL
		return nil
	}

	if w.query != nil {
		orm.Update(w.className, w.query, w.data, w.runOptions)
		// TODO 处理错误
		resp := map[string]interface{}{
			"updatedAt": w.updatedAt,
		}
		w.response = map[string]interface{}{
			"response": resp,
		}
	} else {
		// 给新用户设置默认 ACL
		if w.data["ACL"] == nil && w.className == "_User" {
			readwrite := map[string]interface{}{
				"read":  true,
				"write": true,
			}
			onlyread := map[string]interface{}{
				"read":  true,
				"write": false,
			}
			objectID := utils.String(w.data["objectId"])
			w.data["ACL"] = map[string]interface{}{
				objectID: readwrite,
				"*":      onlyread,
			}
		}
		// 创建对象
		orm.Create(w.className, w.data, w.runOptions)
		resp := map[string]interface{}{
			"objectId":  w.data["objectId"],
			"createdAt": w.data["createdAt"],
		}
		if w.storage["changedByTrigger"] != nil {
			for k, v := range w.data {
				resp[k] = v
			}
		}
		if w.storage["token"] != nil {
			resp["sessionToken"] = w.storage["token"]
		}
		w.response = map[string]interface{}{
			"status":   201,
			"response": resp,
			"location": w.location(),
		}
	}

	return nil
}

func (w *Write) handleFollowup() error {
	if w.storage != nil && w.storage["clearSessions"] != nil {
		// 修改密码之后，清除 session
		user := map[string]interface{}{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  w.objectID(),
		}
		sessionQuery := map[string]interface{}{
			"user": user,
		}
		delete(w.storage, "clearSessions")
		orm.Destroy("_Session", sessionQuery)
	}

	if w.storage != nil && w.storage["sendVerificationEmail"] != nil {
		// 修改邮箱之后需要发送验证邮件
		delete(w.storage, "sendVerificationEmail")
		SendVerificationEmail(w.data)
	}

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

func (w *Write) objectID() interface{} {
	if w.data["objectId"] != nil {
		return w.data["objectId"]
	}
	return w.query["objectId"]
}
