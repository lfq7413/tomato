package rest

import (
	"reflect"
	"regexp"
	"strings"
	"time"

	"strconv"

	am "github.com/lfq7413/tomato/auth"
	"github.com/lfq7413/tomato/cache"
	"github.com/lfq7413/tomato/client"
	"github.com/lfq7413/tomato/cloud"
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/files"
	"github.com/lfq7413/tomato/livequery"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// Write ...
type Write struct {
	auth                       *Auth
	className                  string
	query                      types.M
	data                       types.M
	originalData               types.M
	storage                    types.M
	RunOptions                 types.M
	response                   types.M
	updatedAt                  string
	responseShouldHaveUsername bool
	clientSDK                  map[string]string
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
	clientSDK map[string]string,
) (*Write, error) {
	if auth == nil {
		auth = Nobody()
	}
	if data == nil {
		data = types.M{}
	}
	// 当为 create 请求时，写入数据中不应该包含 objectId
	if query == nil && data["objectId"] != nil {
		return nil, errs.E(errs.InvalidKeyName, "objectId is an invalid field name.")
	}
	var queryCopy types.M
	if query == nil {
		queryCopy = nil
	} else {
		queryCopy = utils.CopyMap(query)
	}
	// query,data 可能会被修改，所以先复制出来
	// response 为最终返回的结果，其中包含三个字段：response、status、location
	write := &Write{
		auth:                       auth,
		className:                  className,
		query:                      queryCopy,
		data:                       utils.CopyMap(data),
		originalData:               originalData,
		storage:                    types.M{},
		RunOptions:                 types.M{},
		response:                   nil,
		updatedAt:                  utils.TimetoString(time.Now().UTC()),
		responseShouldHaveUsername: false,
		clientSDK:                  clientSDK,
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
	err = w.validateSchema()
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
	err = w.createSessionTokenIfNeeded()
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
		acl = append(acl, utils.S(w.auth.User["objectId"]))
		acl = append(acl, w.auth.GetUserRoles()...)
	}
	w.RunOptions["acl"] = acl
	return nil
}

// validateClientClassCreation 检测是否允许创建类
func (w *Write) validateClientClassCreation() error {
	if config.TConfig.AllowClientClassCreation {
		return nil
	}
	if w.auth.IsMaster {
		return nil
	}
	for _, v := range orm.SystemClasses {
		if v == w.className {
			return nil
		}
	}
	// 允许操作已存在的表
	schema := orm.TomatoDBController.LoadSchema(nil)
	hasClass := schema.HasClass(w.className)
	if hasClass {
		return nil
	}
	// 无法操作不存在的表
	return errs.E(errs.OperationForbidden, "This user is not allowed to access non-existent class: "+w.className)
}

// validateSchema 校验数据与权限是否允许进行当前操作
func (w *Write) validateSchema() error {
	return orm.TomatoDBController.ValidateObject(w.className, w.data, w.query, w.RunOptions)
}

// handleInstallation 处理 _Installation 表的操作
// 新增安装记录时，必须要有设备标识： deviceToken 或者 installationId ，必须要有设备类型 deviceType
// 更新安装记录时，不能更新 installationId deviceToken deviceType 三个字段
func (w *Write) handleInstallation() error {
	if w.response != nil || w.className != "_Installation" {
		return nil
	}

	if w.query == nil && w.data["deviceToken"] == nil && w.data["installationId"] == nil && w.auth.InstallationID == "" {
		// create 操作时，设备 id 不能为空
		return errs.E(errs.MissingRequiredFieldError, "at least one ID field (deviceToken, installationId) must be specified in this operation")
	}

	// 	如果 deviceToken 为 64 位，则认为是 iOS 设备
	if w.data["deviceToken"] != nil && len(utils.S(w.data["deviceToken"])) == 64 {
		w.data["deviceToken"] = strings.ToLower(utils.S(w.data["deviceToken"]))
	}

	if w.data["installationId"] != nil {
		w.data["installationId"] = strings.ToLower(utils.S(w.data["installationId"]))
	}

	// 如果 w.data["installationId"] 不存在并且不是 Master 权限，则使用 w.auth.InstallationID
	var installationID string
	if w.data["installationId"] != nil {
		installationID = utils.S(w.data["installationId"])
	}
	if installationID == "" && w.auth.IsMaster == false {
		installationID = w.auth.InstallationID
	}
	installationID = strings.ToLower(installationID)

	// 不处理设备数据为空的情况
	if w.query != nil && w.data["deviceToken"] == nil && installationID == "" && w.data["deviceType"] == nil {
		return nil
	}

	var idMatch types.M
	var objectIDMatch types.M
	var installationIDMatch types.M
	deviceTokenMatches := []types.M{}

	orQueries := types.S{}
	if w.query != nil && w.query["objectId"] != nil {
		orQueries = append(orQueries, types.M{"objectId": w.query["objectId"]})
	}
	if installationID != "" {
		orQueries = append(orQueries, types.M{"installationId": installationID})
	}
	if w.data["deviceToken"] != nil {
		orQueries = append(orQueries, types.M{"deviceToken": w.data["deviceToken"]})
	}
	if len(orQueries) == 0 {
		return nil
	}

	// 查找跟提交的 objectId installationId deviceToken 相同的记录
	results, err := orm.TomatoDBController.Find("_Installation", types.M{"$or": orQueries}, types.M{})
	if err != nil {
		return err
	}

	for _, v := range results {
		if result := utils.M(v); result != nil {
			if w.query != nil && w.query["objectId"] != nil && utils.S(result["objectId"]) == utils.S(w.query["objectId"]) {
				objectIDMatch = result
			}
			if installationID != "" && utils.S(result["installationId"]) == installationID {
				installationIDMatch = result
			}
			if w.data["deviceToken"] != nil && utils.S(result["deviceToken"]) == utils.S(w.data["deviceToken"]) {
				deviceTokenMatches = append(deviceTokenMatches, result)
			}
		}
	}

	// 如果是 update 操作，并且 objectId 存在，
	// 校验是否能对 installationId、deviceToken、deviceType 进行修改
	if w.query != nil && w.query["objectId"] != nil {
		if objectIDMatch == nil {
			return errs.E(errs.ObjectNotFound, "Object not found for update.")
		}
		if w.data["installationId"] != nil && objectIDMatch["installationId"] != nil &&
			utils.S(w.data["installationId"]) != utils.S(objectIDMatch["installationId"]) {
			// installationId 不能修改
			return errs.E(errs.ChangedImmutableFieldError, "installationId may not be changed in this operation")
		}
		if w.data["deviceToken"] != nil && objectIDMatch["deviceToken"] != nil &&
			utils.S(w.data["deviceToken"]) != utils.S(objectIDMatch["deviceToken"]) &&
			w.data["installationId"] == nil && objectIDMatch["installationId"] == nil {
			// deviceToken 不能修改
			return errs.E(errs.ChangedImmutableFieldError, "deviceToken may not be changed in this operation")
		}
		if w.data["deviceType"] != nil && objectIDMatch["deviceType"] != nil &&
			utils.S(w.data["deviceType"]) != utils.S(objectIDMatch["deviceType"]) {
			// deviceType 不能修改
			return errs.E(errs.ChangedImmutableFieldError, "deviceType may not be changed in this operation")
		}
	}

	if w.query != nil && w.query["objectId"] != nil && objectIDMatch != nil {
		idMatch = objectIDMatch
	}

	if installationID != "" && installationIDMatch != nil {
		idMatch = installationIDMatch
	}

	if w.query == nil && w.data["deviceType"] == nil && idMatch == nil {
		// create 操作时，设备类型不能为空
		return errs.E(errs.MissingRequiredFieldError, "deviceType must be specified in this operation")
	}

	// 以下逻辑为：检测是否需要合并数据

	var objID string
	// 要更新的 installationId 不存在，存在一条 deviceToken 时，则更新这条记录，存在多条时则清理多余数据
	if idMatch == nil {
		if deviceTokenMatches == nil || len(deviceTokenMatches) == 0 {
			// 要更新的 deviceToken 不存在
			objID = ""
		} else if len(deviceTokenMatches) == 1 &&
			(deviceTokenMatches[0]["installationId"] == nil || installationID == "") {
			// 要更新的 deviceToken 只存在一个，并且 installationId 不是同时都有
			// 则更新这条记录
			objID = utils.S(deviceTokenMatches[0]["objectId"])
		} else if w.data["installationId"] == nil {
			// 当有多个 deviceToken 时，必须指定 installationId
			return errs.E(errs.InvalidInstallationIDError, "Must specify installationId when deviceToken matches multiple Installation objects")
		} else {
			// 有多个 deviceToken 时，清理多余数据
			// 或者只有一个 deviceToken，但是同时存在两个 installationId 时，也要清理数据
			// 清理多余的 deviceToken，保留对应 installationId 的那个
			// 当前位置为 idMatch == nil ，所以不存在 installationId 对应的记录
			installationID := types.M{
				"$ne": installationID,
			}
			delQuery := types.M{
				"deviceToken":    w.data["deviceToken"],
				"installationId": installationID,
			}
			if w.data["appIdentifier"] != nil {
				delQuery["appIdentifier"] = w.data["appIdentifier"]
			}
			err := orm.TomatoDBController.Destroy("_Installation", delQuery, types.M{})
			if err != nil {
				if errs.GetErrorCode(err) == errs.ObjectNotFound {

				} else {
					return err
				}
			}
			objID = ""
		}
	} else {
		// 要更新的 installationId 存在
		if deviceTokenMatches == nil || len(deviceTokenMatches) == 0 {
			// 不存在对应的 deviceToken ，无需清理
			objID = utils.S(idMatch["objectId"])
		} else if deviceTokenMatches != nil && len(deviceTokenMatches) == 1 &&
			deviceTokenMatches[0]["installationId"] == nil {
			// deviceToken 存在，且只有一条，并且这条记录中的 installationId 为空
			// 首先清理 idMatch 对应的记录
			// 然后合并要更新的数据到 deviceToken 对应的记录
			delQuery := types.M{
				"objectId": idMatch["objectId"],
			}
			err := orm.TomatoDBController.Destroy("_Installation", delQuery, nil)
			if err != nil {
				if errs.GetErrorCode(err) == errs.ObjectNotFound {

				} else {
					return err
				}
			}
			objID = utils.S(deviceTokenMatches[0]["objectId"])
		} else {
			// deviceToken 不存在，或者有多条，或者存在 installationId 时
			if w.data["deviceToken"] != nil && utils.S(idMatch["deviceToken"]) != utils.S(w.data["deviceToken"]) {
				// deviceToken 有多条，并且与 idMatch 中的 deviceToken 不一致时
				// 清理多余数据，只保留 installationId 对应的数据
				// 合并要更新的数据到 installationId 对应的记录上
				delQuery := types.M{
					"deviceToken": w.data["deviceToken"],
				}
				// 当存在唯一 installationId 时，保护其不被删除
				if w.data["installationId"] != nil {
					delQuery["installationId"] = types.M{"$ne": w.data["installationId"]}
				} else if idMatch["objectId"] != nil &&
					w.data["objectId"] != nil &&
					utils.S(idMatch["objectId"]) == utils.S(w.data["objectId"]) {
					delQuery["objectId"] = types.M{"$ne": idMatch["objectId"]}
				} else {
					// 无需清理数据
					objID = utils.S(idMatch["objectId"])
				}
				// 需要清理数据
				if objID == "" {
					if w.data["appIdentifier"] != nil {
						delQuery["appIdentifier"] = w.data["appIdentifier"]
					}
					err := orm.TomatoDBController.Destroy("_Installation", delQuery, nil)
					if err != nil {
						if errs.GetErrorCode(err) == errs.ObjectNotFound {

						} else {
							return err
						}
					}
				}
			}
			objID = utils.S(idMatch["objectId"])
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
			"expiresAt": types.M{
				"__type": "Date",
				"iso":    utils.TimetoString(expiresAt),
			},
		}
		// 添加请求数据中的各字段
		for k, v := range w.data {
			if k == "objectId" {
				continue
			}
			sessionData[k] = v
		}
		// 以 Master 权限去创建 session
		write, err := NewWrite(Master(), "_Session", nil, sessionData, types.M{}, w.clientSDK)
		if err != nil {
			return err
		}
		results, err := write.Execute()
		if err != nil {
			return err
		}
		sess := utils.M(results["response"])
		if sess == nil {
			return errs.E(errs.InternalServerError, "Error creating session.")
		}
		sessionData["objectId"] = sess["objectId"]
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
		if utils.S(w.data["username"]) == "" {
			return errs.E(errs.UsernameMissing, "bad or missing username")
		}
		if utils.S(w.data["password"]) == "" {
			return errs.E(errs.PasswordMissing, "password is required")
		}
	}

	// 不存在第三方登录数据时，直接返回
	if w.data["authData"] == nil || len(utils.M(w.data["authData"])) == 0 {
		return nil
	}

	authData := utils.M(w.data["authData"])
	canHandleAuthData := true

	if len(authData) > 0 {
		// authData 中包含 id 时，才需要进行处理
		for _, v := range authData {
			providerAuthData := utils.M(v)
			hasToken := (providerAuthData != nil && utils.S(providerAuthData["id"]) != "")
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

	if results != nil && len(results) > 0 {
		if w.query == nil {
			// 存在一个用户，并且是 create 请求时，进行登录
			userResult := utils.M(results[0])
			delete(userResult, "password")

			// 在 location() 之前设置 objectId，否则 w.data["objectId"] 可能为空
			w.data["objectId"] = userResult["objectId"]

			// 检测 authData 是否需要更新
			mutatedAuthData := types.M{}
			for provider, providerData := range authData {
				if auth := utils.M(userResult["authData"]); auth != nil {
					userAuthData := auth[provider]
					if reflect.DeepEqual(providerData, userAuthData) == false {
						mutatedAuthData[provider] = providerData
					}
				} else {
					mutatedAuthData[provider] = providerData
				}
			}

			w.response = types.M{
				"status":   200,
				"response": userResult,
				"location": w.location(),
			}

			// 未修改任何数据，直接返回
			if len(mutatedAuthData) == 0 {
				return nil
			}

			// 当第三方登录信息中的 token 刷新时，就需要更新 authData
			// 仅验证需要修改的部分
			err = w.handleAuthDataValidation(mutatedAuthData)
			if err != nil {
				return err
			}
			// 添加新的 authData 到返回数据中
			userAuthData := utils.M(userResult["authData"])
			if userAuthData == nil {
				userAuthData = types.M{}
			}
			for provider, providerData := range mutatedAuthData {
				userAuthData[provider] = providerData
			}
			userResult["authData"] = userAuthData
			w.response["response"] = userResult

			// 更新数据库中的 authData 字段
			_, err = orm.TomatoDBController.Update(w.className, types.M{"objectId": w.data["objectId"]}, types.M{"authData": mutatedAuthData}, types.M{}, false)
			return err
		} else if w.query != nil && w.query["objectId"] != nil {
			// 存在一个用户，并且当前为 update 请求，校验 objectId 是否一致
			user := utils.M(results[0])
			if utils.S(user["objectId"]) != utils.S(w.query["objectId"]) {
				// auth 已经被使用
				return errs.E(errs.AccountAlreadyLinked, "this auth is already used")
			}
		}
	}

	// 当前第三方数据未关联任何用户或者是 update 请求时，会来到这里
	return w.handleAuthDataValidation(authData)
}

// handleAuthDataValidation 校验第三方登录数据
func (w *Write) handleAuthDataValidation(authData types.M) error {
	for k, v := range authData {
		if v == nil {
			continue
		}
		err := am.ValidateAuthData(k, utils.M(v))
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
		provider := utils.M(v)
		if provider == nil {
			continue
		}
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
		findPromise, err = orm.TomatoDBController.Find(w.className, where, types.M{})
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
	if cloud.TriggerExists(cloud.TypeBeforeSave, w.className) == false {
		return nil
	}

	extraData := types.M{"className": w.className}
	if w.query != nil && w.query["objectId"] != nil {
		extraData["objectId"] = w.query["objectId"]
	}

	var originalObject types.M
	// 不添加原始对象
	updatedObject := inflate(extraData, nil)
	if w.query != nil && w.query["objectId"] != nil {
		// update 时初始化原始对象
		originalObject = inflate(extraData, w.originalData)
	}
	// 把需要更新的数据添加进来
	for k, v := range w.sanitizedData() {
		updatedObject[k] = v
	}

	response, err := maybeRunTrigger(cloud.TypeBeforeSave, w.auth, updatedObject, originalObject)
	if err != nil {
		return err
	}
	if response != nil && utils.M(response["object"]) != nil {
		object := utils.M(response["object"])
		// 删除额外添加的字段
		delete(object, "className")
		if w.query != nil && w.query["objectId"] != nil {
			delete(object, "objectId")
		}
		fields := []string{}
		for k, v := range object {
			if reflect.DeepEqual(w.data[k], v) == false {
				fields = append(fields, k)
			}
		}
		if len(fields) > 0 {
			w.storage["fieldsChangedByTrigger"] = fields
			w.data = object
		}
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

	if w.auth.IsMaster == false {
		if _, ok := w.data["emailVerified"]; ok {
			return errs.E(errs.OperationForbidden, "Clients aren't allowed to manually update email verification.")
		}
	}

	// 如果是正在更新 _User ，则清除相应用户的 session 缓存
	if w.query != nil {
		where := types.M{
			"user": types.M{
				"__type":    "Pointer",
				"className": "_User",
				"objectId":  w.objectID(),
			},
		}
		query, err := NewQuery(Master(), "_Session", where, types.M{}, w.clientSDK)
		if err != nil {
			return err
		}
		response, err := query.Execute()
		if err != nil {
			return err
		}

		if utils.HasResults(response) {
			results := utils.A(response["results"])
			for _, result := range results {
				session := utils.M(result)
				cache.User.Del(utils.S(session["sessionToken"]))
			}
		}
	}

	// 处理密码，计算 sha256
	// TODO 后续需要加盐提高安全性
	if w.data["password"] != nil {
		// 检测密码
		err := w.validatePasswordPolicy()
		if err != nil {
			return err
		}

		if w.query != nil {
			// 如果是 update 请求时，标识出需要清理 Sessions ，并生成新的 Session
			w.storage["clearSessions"] = true
			if w.auth.IsMaster == false {
				w.storage["generateNewSession"] = true
			}
		}
		w.data["_hashed_password"] = utils.Hash(utils.S(w.data["password"]))
		delete(w.data, "password")
	}

	// 处理用户名，检测用户名是否唯一
	err := w.validateUserName()
	if err != nil {
		return err
	}

	// 处理 email ，检测合法性、检测是否唯一
	err = w.validateEmail()
	if err != nil {
		return err
	}

	return nil
}

// validateUserName 处理用户名，检测用户名是否唯一
func (w *Write) validateUserName() error {
	if w.data["username"] == nil {
		// 如果是 create 请求，则生成随机 ID
		if w.query == nil {
			w.data["username"] = utils.CreateObjectID()
			w.responseShouldHaveUsername = true
		}
		return nil
	}
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
	results, err := orm.TomatoDBController.Find(w.className, where, option)
	if err != nil {
		return err
	}
	if len(results) > 0 {
		return errs.E(errs.UsernameTaken, "Account already exists for this username")
	}
	return nil
}

// validateEmail 处理 email ，检测合法性、检测是否唯一
func (w *Write) validateEmail() error {
	if w.data["email"] == nil {
		return nil
	}

	if p := utils.M(w.data["email"]); p != nil {
		if utils.S(p["__op"]) == "Delete" {
			return nil
		}
	}

	if utils.IsEmail(utils.S(w.data["email"])) == false {
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
	results, err := orm.TomatoDBController.Find(w.className, where, option)
	if err != nil {
		return err
	}
	if len(results) > 0 {
		return errs.E(errs.EmailTaken, "Account already exists for this email address")
	}

	// 更新 email ，需要发送验证邮件
	w.storage["sendVerificationEmail"] = true
	SetEmailVerifyToken(w.data)

	return nil
}

// validatePasswordPolicy 校验密码合法性
func (w *Write) validatePasswordPolicy() error {
	if config.TConfig.PasswordPolicy == false {
		return nil
	}
	err := w.validatePasswordRequirements()
	if err != nil {
		return err
	}
	return w.validatePasswordHistory()
}

// validatePasswordRequirements 检测密码是否符合设定的密码规则。 go 中的 regexp 不支持 backtracking ，无法使用 (?= 表达式
func (w *Write) validatePasswordRequirements() error {
	policyError := "Password does not meet the Password Policy requirements."
	password := utils.S(w.data["password"])
	// 检测密码是否符合设定的正则表达式
	if config.TConfig.ValidatorPattern != "" {
		b, _ := regexp.MatchString(config.TConfig.ValidatorPattern, password)
		if b == false {
			return errs.E(errs.ValidationError, policyError)
		}
	}
	// 检测密码是否包含用户名
	if config.TConfig.DoNotAllowUsername {
		if username := utils.S(w.data["username"]); username != "" {
			if strings.Index(password, username) >= 0 {
				return errs.E(errs.ValidationError, policyError)
			}
		} else {
			// username 不存在时，从数据库中取出再去检测
			query := types.M{"objectId": w.objectID()}
			results, err := orm.TomatoDBController.Find("_User", query, types.M{})
			if err != nil {
				return err
			}
			if len(results) != 1 {
				return errs.E(errs.ValidationError, policyError)
			}
			result := utils.M(results[0])
			if result == nil || strings.Index(password, utils.S(result["username"])) >= 0 {
				return errs.E(errs.ValidationError, policyError)
			}
		}
	}
	return nil
}

// validatePasswordHistory 校验密码历史
func (w *Write) validatePasswordHistory() error {
	if w.query == nil || config.TConfig.MaxPasswordHistory == 0 {
		return nil
	}
	query := types.M{
		"objectId": w.objectID(),
	}
	options := types.M{
		"keys": []string{"_password_history", "_hashed_password"},
	}
	results, err := orm.TomatoDBController.Find("_User", query, options)
	if err != nil {
		return err
	}
	if len(results) != 1 {
		return errs.E(errs.ObjectNotFound, "User not found for update.")
	}
	user := utils.M(results[0])
	oldPasswords := []string{}
	if h, ok := user["_password_history"].([]interface{}); ok {
		history := getLastItems(h, config.TConfig.MaxPasswordHistory-1)
		for _, pw := range history {
			if s, ok := pw.(string); ok {
				oldPasswords = append(oldPasswords, s)
			}
		}
	}
	oldPasswords = append(oldPasswords, utils.S(user["password"]))
	newPassword := utils.S(w.data["password"])
	for _, hash := range oldPasswords {
		if utils.Compare(newPassword, hash) {
			return errs.E(errs.ValidationError, "New password should not be the same as last "+strconv.Itoa(config.TConfig.MaxPasswordHistory)+" passwords.")
		}
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

	if w.className == "_Role" {
		cache.Role.Clear()
	}

	if w.className == "_User" && w.query != nil &&
		w.auth.CouldUpdateUserID(utils.S(w.query["objectId"])) == false {
		// 不能更新该用户，Master 可以更新任意用户，普通用户仅可更新自身
		return errs.E(errs.SessionMissing, "cannot modify user "+utils.S(w.query["objectId"]))
	}

	if w.className == "_Product" && w.data["download"] != nil {
		if download := utils.M(w.data["download"]); download != nil {
			w.data["downloadName"] = download["name"]
		}
	}

	// TODO 确保不要出现用户无法访问自身数据的情况
	if acl := utils.M(w.data["ACL"]); acl != nil && acl["*unresolved"] != nil {
		return errs.E(errs.InvalidACL, "Invalid ACL.")
	}

	if w.query != nil {
		// 避免用户自身无法访问 _User 表
		if w.className == "_User" {
			if acl := utils.M(w.data["ACL"]); acl != nil {
				acl[utils.S(w.query["objectId"])] = types.M{
					"read":  true,
					"write": true,
				}
				w.data["ACL"] = acl
			}
		}
		// 更新密码时，同时更新密码重置时间戳
		if w.className == "_User" && w.data["_hashed_password"] != nil && config.TConfig.PasswordPolicy && config.TConfig.MaxPasswordAge > 0 {
			w.data["_password_changed_at"] = utils.TimetoString(time.Now().UTC())
		}
		// 更新时忽略 createdAt 字段
		delete(w.data, "createdAt")
		// 密码历史功能开启时，保存当前密码到历史中
		if w.className == "_User" && w.data["_hashed_password"] != nil && config.TConfig.PasswordPolicy && config.TConfig.MaxPasswordHistory > 0 {
			query := types.M{
				"objectId": w.objectID(),
			}
			options := types.M{
				"keys": []string{"_password_history", "_hashed_password"},
			}
			results, err := orm.TomatoDBController.Find("_User", query, options)
			if err != nil {
				return err
			}
			if len(results) != 1 {
				return errs.E(errs.ObjectNotFound, "User not found for update.")
			}
			user := utils.M(results[0])
			oldPasswords := []interface{}{}
			if h, ok := user["_password_history"].([]interface{}); ok {
				oldPasswords = getLastItems(h, config.TConfig.MaxPasswordHistory-2)
			}
			// _password_history 中保存的密码加上 _hashed_password 密码的数量等于 MaxPasswordHistory
			// 因此当 MaxPasswordHistory = 1 时，只在 _hashed_password 中保存密码
			// 当 MaxPasswordHistory > 1 时，才在 _password_history 中保存历史密码
			if config.TConfig.MaxPasswordHistory > 1 {
				oldPasswords = append(oldPasswords, user["password"])
			}
			w.data["_password_history"] = oldPasswords
		}
		// 执行更新
		response, err := orm.TomatoDBController.Update(w.className, w.query, w.data, w.RunOptions, false)
		if err != nil {
			return err
		}
		response["updatedAt"] = w.updatedAt

		// 如果回调函数修改过数据，把 w.data 中存在但 response 中不存在的字段复制到返回结果中
		w.updateResponseWithData(response, w.data)

		w.response = types.M{
			"response": response,
		}
	} else {
		// 给新用户设置默认 ACL ，密码过期时间戳
		// TODO 为了用户信息安全性，应该禁止其他用户读取
		if w.className == "_User" {
			readwrite := types.M{
				"read":  true,
				"write": true,
			}
			onlyread := types.M{
				"read": true,
			}
			acl := utils.M(w.data["ACL"])
			if acl == nil {
				acl = types.M{}
				acl["*"] = onlyread
			}
			objectID := utils.S(w.data["objectId"])
			acl[objectID] = readwrite
			w.data["ACL"] = acl

			if config.TConfig.PasswordPolicy && config.TConfig.MaxPasswordAge > 0 {
				w.data["_password_changed_at"] = utils.TimetoString(time.Now().UTC())
			}
		}

		// 创建对象
		err := orm.TomatoDBController.Create(w.className, w.data, w.RunOptions)
		if err != nil {
			if w.className != "_User" {
				return err
			}
			if errs.GetErrorCode(err) != errs.DuplicateValue {
				return err
			}

			if w.data["username"] != nil {
				where := types.M{
					"username": w.data["username"],
					"objectId": types.M{"$ne": w.objectID()},
				}
				results, err := orm.TomatoDBController.Find(w.className, where, types.M{"limit": 1})
				if err != nil {
					return err
				}
				if len(results) > 0 {
					return errs.E(errs.UsernameTaken, "Account already exists for this username.")
				}
			}

			if w.data["email"] != nil {
				where := types.M{
					"email":    w.data["email"],
					"objectId": types.M{"$ne": w.objectID()},
				}
				results, err := orm.TomatoDBController.Find(w.className, where, types.M{"limit": 1})
				if err != nil {
					return err
				}
				if len(results) > 0 {
					return errs.E(errs.EmailTaken, "Account already exists for this email address.")
				}
			}

			return errs.E(errs.DuplicateValue, "A duplicate value for a field with unique values was provided")
		}
		response := types.M{
			"objectId":  w.data["objectId"],
			"createdAt": w.data["createdAt"],
		}
		if w.responseShouldHaveUsername {
			response["username"] = w.data["username"]
		}
		// 如果回调函数修改过数据，则将其复制到返回结果中
		w.updateResponseWithData(response, w.data)
		w.response = types.M{
			"status":   201,
			"response": response,
			"location": w.location(),
		}
	}

	return nil
}

// createSessionTokenIfNeeded 创建 Token
func (w *Write) createSessionTokenIfNeeded() error {
	if w.className != "_User" {
		return nil
	}
	if w.query != nil {
		return nil
	}
	return w.createSessionToken()
}

// createSessionToken 创建 Token
func (w *Write) createSessionToken() error {
	token := "r:" + utils.CreateToken()
	expiresAt := config.GenerateSessionExpiresAt()
	user := types.M{
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
		"expiresAt": types.M{
			"__type": "Date",
			"iso":    utils.TimetoString(expiresAt),
		},
	}
	if w.response != nil {
		if r := utils.M(w.response["response"]); r != nil {
			r["sessionToken"] = token
		}
	}

	create, err := NewWrite(Master(), "_Session", nil, sessionData, types.M{}, w.clientSDK)
	if err != nil {
		return err
	}
	_, err = create.Execute()

	return err
}

// handleFollowup 处理后续逻辑
func (w *Write) handleFollowup() error {
	if w.storage != nil && w.storage["clearSessions"] != nil && config.TConfig.RevokeSessionOnPasswordReset {
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
		err := orm.TomatoDBController.Destroy("_Session", sessionQuery, types.M{})
		if err != nil {
			return err
		}
	}

	if w.storage != nil && w.storage["generateNewSession"] != nil {
		delete(w.storage, "generateNewSession")
		err := w.createSessionToken()
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

	hasAfterSaveHook := cloud.TriggerExists(cloud.TypeAfterSave, w.className)
	hasLiveQuery := false
	if livequery.TLiveQuery != nil {
		hasLiveQuery = livequery.TLiveQuery.HasLiveQuery(w.className)
	}
	if hasAfterSaveHook == false && hasLiveQuery == false {
		return nil
	}

	extraData := types.M{"className": w.className}
	if w.query != nil && w.query["objectId"] != nil {
		extraData["objectId"] = w.query["objectId"]
	}

	var originalObject types.M
	if w.query != nil && w.query["objectId"] != nil {
		originalObject = inflate(extraData, w.originalData)
	}

	updatedObject := inflate(extraData, w.originalData)
	// 把需要更新的数据添加进来
	for k, v := range w.sanitizedData() {
		updatedObject[k] = v
	}

	if hasLiveQuery {
		// 尝试通知 LiveQueryServer
		livequery.TLiveQuery.OnAfterSave(w.className, updatedObject, originalObject)
	}

	if hasAfterSaveHook {
		// TODO 不等待回调返回
		maybeRunTrigger(cloud.TypeAfterSave, w.auth, updatedObject, originalObject)
	}

	return nil
}

// location 获取对象路径
func (w *Write) location() string {
	var middle string
	if w.className == "_User" {
		middle = "/users/"
	} else if w.className == "_Session" {
		middle = "/sessions/"
	} else if w.className == "_Role" {
		middle = "/roles/"
	} else if w.className == "_Installation" {
		middle = "/installations/"
	} else {
		middle = "/classes/" + w.className + "/"
	}
	return config.TConfig.ServerURL + middle + utils.S(w.data["objectId"])
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
		if user := utils.M(w.response["response"]); user != nil {
			if authData := utils.M(user["authData"]); authData != nil {
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
}

func (w *Write) updateResponseWithData(response, data types.M) types.M {
	if w.storage["fieldsChangedByTrigger"] == nil {
		return response
	}
	fields := []string{}
	if v, ok := w.storage["fieldsChangedByTrigger"].([]string); ok {
		fields = v
	}
	clientSupportsDelete := client.SupportsForwardDelete(w.clientSDK)
	for _, fieldName := range fields {
		dataValue := data[fieldName]
		if response[fieldName] == nil {
			response[fieldName] = dataValue
		}

		// 删除 __op 操作符
		value := utils.M(response[fieldName])
		if value != nil && value["__op"] != nil {
			delete(response, fieldName)
			if v := utils.M(dataValue); v != nil {
				if clientSupportsDelete && utils.S(v["__op"]) == "Delete" {
					response[fieldName] = dataValue
				}
			}
		}

	}

	return response
}

// getLastItems 获取最后 n 个元素
func getLastItems(items []interface{}, n int) []interface{} {
	if items == nil {
		return items
	}
	if n <= 0 {
		return []interface{}{}
	}
	l := len(items)
	if l < n {
		return items
	}
	return items[(l - n):]
}
