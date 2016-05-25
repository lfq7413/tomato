package rest

import (
	"regexp"
	"strings"
	"time"

	"github.com/lfq7413/tomato/authdatamanager"
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/files"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// Write ...
type Write struct {
	auth         *Auth
	className    string
	query        types.M
	data         types.M
	originalData types.M
	storage      types.M
	runOptions   types.M
	response     types.M
	updatedAt    string
}

// NewWrite 可用于 create 和 update ， create 时 	query 为 nil
// query 查询条件，当 update 请求时不为空
// data 写入数据
// originalData 原始对象数据，当 update 请求时不为空
func NewWrite(
	auth *Auth,
	className string,
	query types.M,
	data types.M,
	originalData types.M,
) (*Write, error) {
	// 当为 create 请求时，写入数据中不应该包含 objectId
	if query == nil && data["objectId"] != nil {
		return nil, errs.E(errs.InvalidKeyName, "objectId is an invalid field name.")
	}
	// query,data 可能会被修改，所以先复制出来
	// response 为最终返回的结果，其中包含三个字段：response、status、location
	write := &Write{
		auth:         auth,
		className:    className,
		query:        utils.CopyMap(query),
		data:         utils.CopyMap(data),
		originalData: originalData,
		storage:      types.M{},
		runOptions:   types.M{},
		response:     nil,
		updatedAt:    utils.TimetoString(time.Now().UTC()),
	}
	return write, nil
}

// Execute 执行写入操作，并返回结果
func (w *Write) Execute() (types.M, error) {
	err := w.getUserAndRoleACL()
	if err != nil {
		return nil, err
	}
	err = w.validateClientClassCreation()
	if err != nil {
		return nil, err
	}
	err = w.validateSchema()
	if err != nil {
		return nil, err
	}
	err = w.handleInstallation()
	if err != nil {
		return nil, err
	}
	err = w.handleSession()
	if err != nil {
		return nil, err
	}
	err = w.validateAuthData()
	if err != nil {
		return nil, err
	}
	err = w.runBeforeTrigger()
	if err != nil {
		return nil, err
	}
	err = w.setRequiredFieldsIfNeeded()
	if err != nil {
		return nil, err
	}
	err = w.transformUser()
	if err != nil {
		return nil, err
	}
	err = w.expandFilesForExistingObjects()
	if err != nil {
		return nil, err
	}
	err = w.runDatabaseOperation()
	if err != nil {
		return nil, err
	}
	err = w.handleFollowup()
	if err != nil {
		return nil, err
	}
	err = w.runAfterTrigger()
	if err != nil {
		return nil, err
	}

	w.cleanUserAuthData()

	return w.response, nil
}

// getUserAndRoleACL 获取用户角色信息，写入 acl 字段
func (w *Write) getUserAndRoleACL() error {
	if w.auth.IsMaster {
		return nil
	}
	acl := []string{"*"}
	if w.auth.User != nil {
		acl = append(acl, w.auth.User["objectId"].(string))
		acl = append(acl, w.auth.GetUserRoles()...)
	}
	w.runOptions["acl"] = acl
	return nil
}

// validateClientClassCreation 检测是否允许创建类
func (w *Write) validateClientClassCreation() error {
	sysClass := orm.SystemClasses
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
	// 允许操作已存在的表
	if orm.CollectionExists(w.className) {
		return nil
	}
	// 无法操作不存在的表
	return errs.E(errs.OperationForbidden, "This user is not allowed to access non-existent class: "+w.className)
}

// validateSchema 校验数据与权限是否允许进行当前操作
func (w *Write) validateSchema() error {
	return orm.ValidateObject(w.className, w.data, w.query, w.runOptions)
}

// handleInstallation 处理 _Installation 表的操作
func (w *Write) handleInstallation() error {
	if w.response != nil || w.className != "_Installation" {
		return nil
	}

	if w.query == nil && w.data["deviceToken"] == nil && w.data["installationId"] == nil {
		// create 操作时，设备 id 不能为空
		return errs.E(errs.MissingRequiredFieldError, "at least one ID field (deviceToken, installationId) must be specified in this operation")
	}

	if w.query == nil && w.data["deviceType"] == nil {
		// create 操作时，设备类型不能为空
		return errs.E(errs.MissingRequiredFieldError, "deviceType must be specified in this operation")
	}

	// 	如果 deviceToken 为 64 位，则认为是 iOS 设备
	if w.data["deviceToken"] != nil && len(utils.String(w.data["deviceToken"])) == 64 {
		w.data["deviceToken"] = strings.ToLower(utils.String(w.data["deviceToken"]))
	}

	if w.data["installationId"] != nil {
		w.data["installationId"] = strings.ToLower(utils.String(w.data["installationId"]))
	}

	var idMatch types.M
	var deviceTokenMatches types.S

	// 如果是 update 操作，并且 objectId 存在，
	// 校验是否能对 installationId、deviceToken、deviceType 进行修改
	if w.query != nil && w.query["objectId"] != nil {
		results, err := orm.Find("_Installation", types.M{"objectId": w.query["objectId"]}, types.M{})
		if err != nil {
			return err
		}
		if results == nil || len(results) == 0 {
			return errs.E(errs.ObjectNotFound, "Object not found for update.")
		}

		idMatch = utils.MapInterface(results[0])
		if w.data["installationId"] != nil && idMatch["installationId"] != nil &&
			w.data["installationId"] != idMatch["installationId"] {
			// installationId 不能修改
			return errs.E(errs.ChangedImmutableFieldError, "installationId may not be changed in this operation")
		}
		if w.data["deviceToken"] != nil && idMatch["deviceToken"] != nil &&
			w.data["deviceToken"] != idMatch["deviceToken"] &&
			w.data["installationId"] == nil && idMatch["installationId"] == nil {
			// deviceToken 不能修改
			return errs.E(errs.ChangedImmutableFieldError, "deviceToken may not be changed in this operation")
		}
		if w.data["deviceType"] != nil && idMatch["deviceType"] != nil &&
			w.data["deviceType"] != idMatch["deviceType"] {
			// deviceType 不能修改
			return errs.E(errs.ChangedImmutableFieldError, "deviceType may not be changed in this operation")
		}
	}

	// 检查是否已经存在 installationId、deviceToken
	idMatch = nil
	if w.data["installationId"] != nil {
		results, err := orm.Find("_Installation", types.M{"installationId": w.data["installationId"]}, types.M{})
		if err != nil {
			return err
		}
		if results != nil && len(results) > 0 {
			// 只取第一个结果
			idMatch = utils.MapInterface(results[0])
		}
	}
	if w.data["deviceToken"] != nil {
		results, err := orm.Find("_Installation", types.M{"deviceToken": w.data["deviceToken"]}, types.M{})
		if err != nil {
			return err
		}
		if results != nil {
			deviceTokenMatches = results
		}
	}

	var objID string
	// 要更新的 installationId 不存在
	if idMatch == nil {
		if deviceTokenMatches == nil || len(deviceTokenMatches) == 0 {
			// 要更新的 deviceToken 不存在
			objID = ""
		} else if len(deviceTokenMatches) == 1 &&
			(utils.MapInterface(deviceTokenMatches[0])["installationId"] == nil || w.data["installationId"] == nil) {
			// 要更新的 deviceToken 只存在一个，并且 installationId 不是同时都有
			objID = utils.String(utils.MapInterface(deviceTokenMatches[0])["objectId"])
		} else if w.data["installationId"] == nil {
			// 当有多个 deviceToken 时，必须指定 installationId
			return errs.E(errs.InvalidInstallationIDError, "Must specify installationId when deviceToken matches multiple Installation objects")
		} else {
			// 有多个 deviceToken 时，清理多余数据
			// 或者只有一个 deviceToken，但是同时存在两个 installationId 时，也要清理数据
			// 清理多余的 deviceToken，保留对应 installationId 的那个
			// 当前位置为 idMatch == nil ，所以不存在 installationId 对应的记录
			installationID := types.M{
				"$ne": w.data["installationId"],
			}
			delQuery := types.M{
				"deviceToken":    w.data["deviceToken"],
				"installationId": installationID,
			}
			if w.data["appIdentifier"] != nil {
				delQuery["appIdentifier"] = w.data["appIdentifier"]
			}
			err := orm.Destroy("_Installation", delQuery, types.M{})
			if err != nil {
				return err
			}
			objID = ""
		}
	} else {
		// 要更新的 installationId 存在
		if deviceTokenMatches != nil && len(deviceTokenMatches) == 1 &&
			utils.MapInterface(deviceTokenMatches[0])["installationId"] == nil {
			// deviceToken 存在，且只有一条，并且这条记录中的 installationId 为空
			// 首先清理 idMatch 对应的记录
			// 然后合并要更新的数据到 deviceToken 对应的记录
			delQuery := types.M{
				"objectId": idMatch["objectId"],
			}
			err := orm.Destroy("_Installation", delQuery, nil)
			if err != nil {
				return err
			}
			objID = utils.String(utils.MapInterface(deviceTokenMatches[0])["objectId"])
		} else {
			// deviceToken 不存在，或者有多条，或者存在 installationId 时
			if w.data["deviceToken"] != nil && idMatch["deviceToken"] != w.data["deviceToken"] {
				// deviceToken 有多条，并且与 idMatch 中的 deviceToken 不一致时
				// 清理多余数据，只保留 installationId 对应的数据
				// 合并要更新的数据到 installationId 对应的记录上
				delQuery := types.M{
					"deviceToken": w.data["deviceToken"],
				}
				// 当存在唯一 installationId 时，保护其不被删除
				if w.data["installationId"] != nil {
					delQuery["installationId"] = types.M{"$ne": w.data["installationId"]}
				} else if idMatch["objectId"] != nil && w.data["objectId"] != nil && idMatch["objectId"].(string) == w.data["objectId"].(string) {
					delQuery["objectId"] = types.M{"$ne": idMatch["objectId"]}
				} else {
					// 无需清理数据
					objID = utils.String(idMatch["objectId"])
				}
				// 需要清理数据
				if objID == "" {
					if w.data["appIdentifier"] != nil {
						delQuery["appIdentifier"] = w.data["appIdentifier"]
					}
					err := orm.Destroy("_Installation", delQuery, nil)
					if err != nil {
						return err
					}
				}
			}
			objID = utils.String(idMatch["objectId"])
		}
	}
	// objID 不为空时，转换当前请求为 update 请求
	if objID != "" {
		w.query = types.M{
			"objectId": objID,
		}
		delete(w.data, "objectId")
		delete(w.data, "createdAt")
	}
	// TODO Validate ops (add/remove on channels, $inc on badge, etc.)

	return nil
}

// handleSession 处理 _Session 表的操作
func (w *Write) handleSession() error {
	if w.response != nil || w.className != "_Session" {
		return nil
	}

	if w.auth.User == nil && w.auth.IsMaster == false {
		return errs.E(errs.InvalidSessionToken, "Session token required.")
	}

	if w.data["ACL"] != nil {
		return errs.E(errs.InvalidKeyName, "Cannot set ACL on a Session.")
	}

	// 当前为 create 请求，并且不是 Master 权限时
	if w.query == nil && w.auth.IsMaster == false {
		// 生成 token ，过期时间为 1 年
		token := "r:" + utils.CreateToken()
		expiresAt := config.GenerateSessionExpiresAt()
		user := types.M{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  w.auth.User["objectId"],
		}
		createdWith := types.M{
			"action": "create",
		}
		sessionData := types.M{
			"sessionToken": token,
			"user":         user,
			"createdWith":  createdWith,
			"restricted":   true,
			"expiresAt":    utils.TimetoString(expiresAt),
		}
		// 添加请求数据中的各字段
		for k, v := range w.data {
			if k == "objectId" {
				continue
			}
			sessionData[k] = v
		}
		// 以 Master 权限去创建 session
		write, err := NewWrite(Master(), "_Session", nil, sessionData, types.M{})
		if err != nil {
			return err
		}
		results, err := write.Execute()
		if err != nil {
			return err
		}
		if results["response"] == nil {
			return errs.E(errs.InternalServerError, "Error creating session.")
		}
		sessionData["objectId"] = utils.MapInterface(results["response"])["objectId"]
		w.response = types.M{
			"status":   201,
			"location": results["location"],
			"response": sessionData,
		}
	}

	return nil
}

// validateAuthData 校验用户登录数据，仅处理对 _User 表的操作
func (w *Write) validateAuthData() error {
	if w.className != "_User" {
		return nil
	}

	// 当前 create 请求，并且不存在第三方登录数据时
	if w.query == nil && w.data["authData"] == nil {
		if utils.String(w.data["username"]) == "" {
			return errs.E(errs.UsernameMissing, "bad or missing username")
		}
		if utils.String(w.data["password"]) == "" {
			return errs.E(errs.PasswordMissing, "password is required")
		}
	}

	// 不存在第三方登录数据时，直接返回
	if w.data["authData"] == nil || len(utils.MapInterface(w.data["authData"])) == 0 {
		return nil
	}

	authData := utils.MapInterface(w.data["authData"])
	canHandleAuthData := true

	if len(authData) > 0 {
		// authData 中包含 id 时，才需要进行处理
		for _, v := range authData {
			providerAuthData := utils.MapInterface(v)
			hasToken := (providerAuthData != nil && providerAuthData["id"] != nil)
			canHandleAuthData = (canHandleAuthData && (hasToken || providerAuthData == nil))
		}
		if canHandleAuthData {
			return w.handleAuthData(authData)
		}
	}
	// 这个 authentication 不支持
	return errs.E(errs.UnsupportedService, "This authentication method is unsupported.")
}

// handleAuthData 处理第三方登录数据
func (w *Write) handleAuthData(authData types.M) error {
	// 校验第三方数据
	err := w.handleAuthDataValidation(authData)
	if err != nil {
		return err
	}
	results, err := w.findUsersWithAuthData(authData)
	if err != nil {
		return err
	}
	if results != nil && len(results) > 1 {
		// auth 已经被多个用户使用
		return errs.E(errs.AccountAlreadyLinked, "this auth is already used")
	}

	// 保存登录方式
	keys := []string{}
	for k := range authData {
		keys = append(keys, k)
	}
	w.storage["authProvider"] = strings.Join(keys, ",")

	if results != nil || len(results) > 0 {
		if w.query == nil {
			// 存在一个用户，并且是 create 请求时，进行登录
			user := utils.MapInterface(results[0])
			delete(user, "password")
			// 在 location() 之前设置 objectId，否则 w.data["objectId"] 可能为空
			w.data["objectId"] = user["objectId"]
			w.response = types.M{
				"response": user,
				"location": w.location(),
			}
		} else if w.query != nil && w.query["objectId"] != nil {
			// 存在一个用户，并且当前为 update 请求，校验 objectId 是否一致
			user := utils.MapInterface(results[0])
			if utils.String(user["objectId"]) != utils.String(w.query["objectId"]) {
				// auth 已经被使用
				return errs.E(errs.AccountAlreadyLinked, "this auth is already used")
			}
		}
	}

	return nil
}

// handleAuthDataValidation 校验第三方登录数据
func (w *Write) handleAuthDataValidation(authData types.M) error {
	for k, v := range authData {
		if v == nil {
			continue
		}
		err := authdatamanager.ValidateAuthData(k, utils.MapInterface(v))
		if err != nil {
			// 验证出现问题
			return err
		}
	}

	return nil
}

// findUsersWithAuthData 查找第三方登录数据对应的用户
func (w *Write) findUsersWithAuthData(authData types.M) (types.S, error) {
	query := types.S{}
	for k, v := range authData {
		if v == nil {
			continue
		}
		key := "authData." + k + ".id"
		provider := utils.MapInterface(v)
		q := types.M{
			key: provider["id"],
		}
		query = append(query, q)
	}

	findPromise := types.S{}
	if len(query) > 0 {
		where := types.M{
			"$or": query,
		}
		var err error
		findPromise, err = orm.Find(w.className, where, types.M{})
		if err != nil {
			return nil, err
		}
	}

	return findPromise, nil
}

// runBeforeTrigger 运行数据修改前的回调函数
func (w *Write) runBeforeTrigger() error {
	if w.response != nil {
		return nil
	}
	if TriggerExists(TypeBeforeSave, w.className) == false {
		return nil
	}

	updatedObject := types.M{}
	if w.query != nil && w.query["objectId"] != nil {
		// 如果是更新，则把原始数据添加进来
		for k, v := range w.originalData {
			updatedObject[k] = v
		}
		updatedObject["objectId"] = w.query["objectId"]
	}
	// 把需要更新的数据添加进来
	for k, v := range w.sanitizedData() {
		updatedObject[k] = v
	}

	response := RunTrigger(TypeBeforeSave, w.className, w.auth, updatedObject, w.originalData)
	if response != nil && response["object"] != nil {
		// 运行完回调函数之后，把结果设置回 data ，并标识已被修改
		w.data = utils.MapInterface(response["object"])
		w.storage["changedByTrigger"] = true
		if w.query != nil && w.query["objectId"] != nil {
			delete(w.data, "objectId")
		}
		return w.validateSchema()
	}

	return nil
}

// setRequiredFieldsIfNeeded 设置必要的字段
func (w *Write) setRequiredFieldsIfNeeded() error {
	if w.data != nil {
		// 添加默认字段
		w.data["updatedAt"] = w.updatedAt
		if w.query == nil {
			// create 请求时，添加 createdAt，创建 objectId
			w.data["createdAt"] = w.updatedAt

			if w.data["objectId"] == nil {
				w.data["objectId"] = utils.CreateObjectID()
			}
		}
	}

	return nil
}

// transformUser 转换用户数据，仅处理 _User 表
func (w *Write) transformUser() error {
	if w.className != "_User" {
		return nil
	}

	// 如果是创建用户，则先创建 token
	if w.query == nil {
		token := "r:" + utils.CreateToken()
		w.storage["token"] = token
		expiresAt := config.GenerateSessionExpiresAt()
		user := types.M{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  w.objectID(),
		}
		// 确定登录方式
		var authProvider interface{}
		if w.storage["authProvider"] != nil {
			authProvider = w.storage["authProvider"]
		} else {
			authProvider = "password"
		}
		createdWith := types.M{
			"action":       "signup",
			"authProvider": authProvider,
		}
		sessionData := types.M{
			"sessionToken":   token,
			"user":           user,
			"createdWith":    createdWith,
			"restricted":     false,
			"installationId": w.auth.InstallationID,
			"expiresAt":      utils.TimetoString(expiresAt),
		}
		if w.response != nil && w.response["response"] != nil {
			// 此时为第三方登录时的情形，w.response 已经有值
			response := utils.MapInterface(w.response["response"])
			response["sessionToken"] = token
		}

		write, err := NewWrite(Master(), "_Session", nil, sessionData, nil)
		if err != nil {
			return err
		}
		_, err = write.Execute()
		if err != nil {
			return err
		}
	}

	// 如果是正在更新 _User ，则清除相应用户的 session 缓存
	if w.query != nil && w.auth.User != nil && w.auth.User["sessionToken"] != nil {
		usersCache.remove(w.auth.User["sessionToken"].(string))
	}

	// 处理密码，计算 sha256
	// TODO 后续需要加盐提高安全性
	if w.data["password"] == nil {

	} else {
		if w.query != nil && w.auth.IsMaster == false {
			// 如果是 update 请求时，标识出需要清理 Sessions
			w.storage["clearSessions"] = true
		}
		w.data["_hashed_password"] = utils.Hash(utils.String(w.data["password"]))
		delete(w.data, "password")
	}

	// 处理用户名，检测用户名是否唯一
	if w.data["username"] == nil {
		// 如果是 create 请求，则生成随机 ID
		if w.query == nil {
			w.data["username"] = utils.CreateObjectID()
		}
	} else {
		objectID := types.M{
			"$ne": w.objectID(),
		}
		where := types.M{
			"username": w.data["username"],
			"objectId": objectID,
		}
		option := types.M{
			"limit": 1,
		}
		results, err := orm.Find(w.className, where, option)
		if err != nil {
			return err
		}
		if len(results) > 0 {
			return errs.E(errs.UsernameTaken, "Account already exists for this username")
		}
	}

	// 处理 email ，检测合法性、检测是否唯一
	if w.data["email"] == nil {

	} else {
		if utils.IsEmail(utils.String(w.data["email"])) == false {
			return errs.E(errs.InvalidEmailAddress, "Email address format is invalid.")
		}
		objectID := types.M{
			"$ne": w.objectID(),
		}
		where := types.M{
			"email":    w.data["email"],
			"objectId": objectID,
		}
		option := types.M{
			"limit": 1,
		}
		results, err := orm.Find(w.className, where, option)
		if err != nil {
			return err
		}
		if len(results) > 0 {
			return errs.E(errs.EmailTaken, "Account already exists for this email address")
		}

		// 更新 email ，需要发送验证邮件
		w.storage["sendVerificationEmail"] = true
		SetEmailVerifyToken(w.data)
	}

	return nil
}

// expandFilesForExistingObjects 展开文件对象
func (w *Write) expandFilesForExistingObjects() error {
	if w.response != nil && w.response["response"] != nil {
		// 展开文件对象
		files.ExpandFilesInObject(w.response["response"])
	}

	return nil
}

// runDatabaseOperation 执行数据库操作
func (w *Write) runDatabaseOperation() error {
	if w.response != nil {
		return nil
	}

	if w.className == "_User" && w.query != nil &&
		w.auth.CouldUpdateUserID(utils.String(w.query["objectId"])) == false {
		// 不能更新该用户，Master 可以更新任意用户，普通用户仅可更新自身
		return errs.E(errs.SessionMissing, "cannot modify user "+utils.String(w.query["objectId"]))
	}

	if w.className == "_Product" && w.data["download"] != nil {
		download := utils.MapInterface(w.data["download"])
		w.data["downloadName"] = download["name"]
	}

	// TODO 确保不要出现用户无法访问自身数据的情况
	if w.data["ACL"] != nil && utils.MapInterface(w.data["ACL"])["*unresolved"] != nil {
		return errs.E(errs.InvalidAcl, "Invalid ACL.")
	}

	if w.query != nil {
		// 避免用户自身无法访问 _User 表
		if w.className == "_User" && w.data["ACL"] != nil {
			acl := w.data["ACL"].(map[string]interface{})
			acl[w.query["objectId"].(string)] = types.M{
				"read":  true,
				"write": true,
			}
			w.data["ACL"] = acl
		}
		// 执行更新
		resp, err := orm.Update(w.className, w.query, w.data, w.runOptions)
		if err != nil {
			return err
		}
		resp["updatedAt"] = w.updatedAt

		// 如果回调函数修改过数据，把 w.data 中存在但 resp 中不存在的字段复制到返回结果中
		if w.storage["changedByTrigger"] != nil {
			for k, v := range w.data {
				if resp[k] == nil {
					resp[k] = v
				}
			}
		}

		w.response = types.M{
			"response": resp,
		}
	} else {
		// 给新用户设置默认 ACL
		// TODO 为了用户信息安全性，应该禁止其他用户读取
		if w.className == "_User" {
			readwrite := types.M{
				"read":  true,
				"write": true,
			}
			onlyread := types.M{
				"read":  true,
				"write": false,
			}
			acl := w.data["ACL"].(map[string]interface{})
			if acl == nil {
				acl := types.M{}
				acl["*"] = onlyread
			}
			objectID := utils.String(w.data["objectId"])
			acl[objectID] = readwrite
			w.data["ACL"] = acl
		}

		// 创建对象
		err := orm.Create(w.className, w.data, w.runOptions)
		if err != nil {
			return err
		}
		resp := types.M{
			"objectId":  w.data["objectId"],
			"createdAt": w.data["createdAt"],
		}
		// 如果回调函数修改过数据，则将其复制到返回结果中
		if w.storage["changedByTrigger"] != nil {
			for k, v := range w.data {
				resp[k] = v
			}
		}
		// 如果新创建的用户包含 token，则复制到返回结果中
		if w.storage["token"] != nil {
			resp["sessionToken"] = w.storage["token"]
		}
		w.response = types.M{
			"status":   201,
			"response": resp,
			"location": w.location(),
		}
	}

	return nil
}

// handleFollowup 处理后续逻辑
func (w *Write) handleFollowup() error {
	if w.storage != nil && w.storage["clearSessions"] != nil {
		// 修改密码之后，清除 session
		user := types.M{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  w.objectID(),
		}
		sessionQuery := types.M{
			"user": user,
		}
		delete(w.storage, "clearSessions")
		err := orm.Destroy("_Session", sessionQuery, types.M{})
		if err != nil {
			return err
		}
	}

	if w.storage != nil && w.storage["sendVerificationEmail"] != nil {
		// 修改邮箱之后需要发送验证邮件
		delete(w.storage, "sendVerificationEmail")
		SendVerificationEmail(w.data)
	}

	return nil
}

// runAfterTrigger 运行数据修改后的回调函数
func (w *Write) runAfterTrigger() error {
	if w.response == nil || w.response["response"] == nil {
		return nil
	}

	hasAfterSaveHook := TriggerExists(TypeAfterSave, w.className)
	hasLiveQuery := config.TConfig.LiveQuery.HasLiveQuery(w.className)
	if hasAfterSaveHook == false && hasLiveQuery == false {
		return nil
	}

	updatedObject := types.M{"className": w.className}
	if w.query != nil && w.query["objectId"] != nil {
		// 如果是更新，则把原始数据添加进来
		for k, v := range w.originalData {
			updatedObject[k] = v
		}
		updatedObject["objectId"] = w.query["objectId"]
	}
	// 把需要更新的数据添加进来
	for k, v := range w.sanitizedData() {
		updatedObject[k] = v
	}

	// 尝试通知 LiveQueryServer
	config.TConfig.LiveQuery.OnAfterSave(w.className, updatedObject, w.originalData)

	RunTrigger(TypeAfterSave, w.className, w.auth, updatedObject, w.originalData)

	return nil
}

// location 获取对象路径
func (w *Write) location() string {
	var middle string
	if w.className == "_User" {
		middle = "/users/"
	} else {
		middle = "/classes/" + w.className + "/"
	}
	return config.TConfig.ServerURL + middle + utils.String(w.data["objectId"])
}

// objectID 从请求中获取 objectId
func (w *Write) objectID() interface{} {
	if w.data["objectId"] != nil {
		return w.data["objectId"]
	}
	return w.query["objectId"]
}

// sanitizedData 删除无效字段，如 _auth_data, _hashed_password...
func (w *Write) sanitizedData() types.M {
	data := utils.CopyMap(w.data)
	for k := range data {
		// 以字母开头，包含数字字母或下划线的为有效字段
		b, _ := regexp.MatchString("^[A-Za-z][0-9A-Za-z_]*$", k)
		if b == false {
			delete(data, k)
		}
	}
	return data
}

func (w *Write) cleanUserAuthData() {
	if w.response != nil && w.response["response"] != nil && w.className == "_User" {
		user := utils.MapInterface(w.response["response"])
		if user != nil && user["authData"] != nil {
			authData := utils.MapInterface(user["authData"])
			for provider, v := range authData {
				if v == nil {
					delete(authData, provider)
				}
			}
			if len(authData) == 0 {
				delete(user, "authData")
			}
		}
	}
}
