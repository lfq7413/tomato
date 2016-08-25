package rest

import (
	"reflect"
	"testing"
	"time"

	"github.com/lfq7413/tomato/cache"
	"github.com/lfq7413/tomato/cloud"
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/livequery"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

func Test_NewWrite(t *testing.T) {
	var auth *Auth
	var className string
	var query types.M
	var data types.M
	var originalData types.M
	var clientSDK map[string]string
	var result *Write
	var err error
	var expect *Write
	var expectErr error
	/***************************************************************/
	auth = nil
	className = "user"
	query = nil
	data = types.M{
		"objectId": "1001",
	}
	originalData = nil
	clientSDK = nil
	_, err = NewWrite(auth, className, query, data, originalData, clientSDK)
	expectErr = errs.E(errs.InvalidKeyName, "objectId is an invalid field name.")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/***************************************************************/
	auth = nil
	className = "user"
	query = nil
	data = types.M{
		"key": "hello",
	}
	originalData = nil
	clientSDK = nil
	result, err = NewWrite(auth, className, query, data, originalData, clientSDK)
	expect = &Write{
		auth:                       Nobody(),
		className:                  "user",
		query:                      nil,
		data:                       types.M{"key": "hello"},
		originalData:               nil,
		storage:                    types.M{},
		RunOptions:                 types.M{},
		response:                   nil,
		updatedAt:                  utils.TimetoString(time.Now().UTC()),
		responseShouldHaveUsername: false,
		clientSDK:                  nil,
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/***************************************************************/
	auth = nil
	className = "user"
	query = types.M{
		"objectId": "1001",
	}
	data = types.M{
		"key": "hello",
	}
	originalData = types.M{
		"key": "hi",
	}
	clientSDK = nil
	result, err = NewWrite(auth, className, query, data, originalData, clientSDK)
	expect = &Write{
		auth:                       Nobody(),
		className:                  "user",
		query:                      types.M{"objectId": "1001"},
		data:                       types.M{"key": "hello"},
		originalData:               types.M{"key": "hi"},
		storage:                    types.M{},
		RunOptions:                 types.M{},
		response:                   nil,
		updatedAt:                  utils.TimetoString(time.Now().UTC()),
		responseShouldHaveUsername: false,
		clientSDK:                  nil,
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_Execute_Write(t *testing.T) {
	// TODO
}

func Test_getUserAndRoleACL_Write(t *testing.T) {
	var schema types.M
	var object types.M
	var className string
	var w *Write
	var auth *Auth
	var query types.M
	var data types.M
	var originalData types.M
	var expect []string
	/***************************************************************/
	auth = Master()
	query = nil
	data = types.M{}
	originalData = nil
	w, _ = NewWrite(auth, "user", query, data, originalData, nil)
	w.getUserAndRoleACL()
	if _, ok := w.RunOptions["acl"]; ok {
		t.Error("findOptions[acl] exist")
	}
	/***************************************************************/
	auth = Nobody()
	query = nil
	data = types.M{}
	originalData = nil
	w, _ = NewWrite(auth, "user", query, data, originalData, nil)
	w.getUserAndRoleACL()
	expect = []string{"*"}
	if reflect.DeepEqual(expect, w.RunOptions["acl"]) == false {
		t.Error("expect:", expect, "result:", w.RunOptions["acl"])
	}
	/***************************************************************/
	cache.InitCache()
	initEnv()
	className = "_Role"
	schema = types.M{
		"fields": types.M{
			"name":  types.M{"type": "String"},
			"users": types.M{"type": "Relation", "targetClass": "_User"},
			"roles": types.M{"type": "Relation", "targetClass": "_Role"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"name":     "role1001",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "1002",
		"name":     "role1002",
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "_Join:roles:_Role"
	schema = types.M{
		"fields": types.M{
			"relatedId": types.M{"type": "String"},
			"owningId":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId":  "5001",
		"owningId":  "1002",
		"relatedId": "1001",
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "_Join:users:_Role"
	schema = types.M{
		"fields": types.M{
			"relatedId": types.M{"type": "String"},
			"owningId":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId":  "5002",
		"owningId":  "1001",
		"relatedId": "9001",
	}
	orm.Adapter.CreateObject(className, schema, object)
	auth = &Auth{
		IsMaster: false,
		User: types.M{
			"objectId": "9001",
		},
		FetchedRoles: false,
		RolePromise:  nil,
	}
	w, _ = NewWrite(auth, "user", query, data, originalData, nil)
	w.getUserAndRoleACL()
	expect = []string{"*", "9001", "role:role1001", "role:role1002"}
	if reflect.DeepEqual(expect, w.RunOptions["acl"]) == false {
		t.Error("expect:", expect, "result:", w.RunOptions["acl"])
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_validateClientClassCreation_Write(t *testing.T) {
	// 测试用例与 query.validateClientClassCreation 相同
}

func Test_validateSchema(t *testing.T) {
	// 测试用例与 DBController.ValidateObject 相同
}

func Test_handleInstallation(t *testing.T) {
	// TODO
}

func Test_handleSession(t *testing.T) {
	// Execute
	// TODO
}

func Test_validateAuthData(t *testing.T) {
	var className string
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	var result error
	var expect error
	/***************************************************************/
	initEnv()
	className = "user"
	query = nil
	data = types.M{}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.validateAuthData()
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	className = "_User"
	query = nil
	data = types.M{
		"key": "hello",
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.validateAuthData()
	expect = errs.E(errs.UsernameMissing, "bad or missing username")
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	className = "_User"
	query = nil
	data = types.M{
		"username": "joe",
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.validateAuthData()
	expect = errs.E(errs.PasswordMissing, "password is required")
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	className = "_User"
	query = nil
	data = types.M{
		"username": "joe",
		"password": "123",
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.validateAuthData()
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	className = "_User"
	query = nil
	data = types.M{
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.validateAuthData()
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	className = "_User"
	query = nil
	data = types.M{
		"authData": types.M{
			"facebook": types.M{
				"key": "1001",
			},
		},
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.validateAuthData()
	expect = errs.E(errs.UnsupportedService, "This authentication method is unsupported.")
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_runBeforeTrigger(t *testing.T) {
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	var result error
	var expect error
	var expectData types.M
	/***************************************************************/
	query = nil
	data = types.M{}
	originalData = nil
	w, _ = NewWrite(Master(), "user", query, data, originalData, nil)
	result = w.runBeforeTrigger()
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/***************************************************************/
	cloud.BeforeSave("user", func(request cloud.TriggerRequest, response cloud.Response) {
		object := request.Object
		if username := utils.S(object["username"]); username != "" {
			object["username"] = username + "_tomato"
			response.Success(nil)
		} else {
			response.Error(1, "need a username")
		}
	})
	query = nil
	data = types.M{
		"key": "hello",
	}
	originalData = nil
	w, _ = NewWrite(Master(), "user", query, data, originalData, nil)
	result = w.runBeforeTrigger()
	expect = errs.E(1, "need a username")
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	cloud.UnregisterAll()
	/***************************************************************/
	cloud.BeforeSave("user", func(request cloud.TriggerRequest, response cloud.Response) {
		object := request.Object
		if username := utils.S(object["username"]); username != "" {
			object["username"] = username + "_tomato"
			response.Success(nil)
		} else {
			response.Error(1, "need a username")
		}
	})
	query = nil
	data = types.M{
		"username": "joe",
		"key":      "hello",
	}
	originalData = nil
	w, _ = NewWrite(Master(), "user", query, data, originalData, nil)
	result = w.runBeforeTrigger()
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	expectData = types.M{
		"username": "joe_tomato",
		"key":      "hello",
	}
	if reflect.DeepEqual(expectData, w.data) == false {
		t.Error("expect:", expectData, "result:", w.data)
	}
	if reflect.DeepEqual(true, w.storage["changedByTrigger"]) == false {
		t.Error("expect:", true, "result:", w.storage["changedByTrigger"])
	}
	cloud.UnregisterAll()
	/***************************************************************/
	cloud.BeforeSave("user", func(request cloud.TriggerRequest, response cloud.Response) {
		object := request.Object
		if username := utils.S(object["username"]); username != "" {
			response.Success(nil)
		} else {
			response.Error(1, "need a username")
		}
	})
	query = nil
	data = types.M{
		"username": "joe",
		"key":      "hello",
	}
	originalData = nil
	w, _ = NewWrite(Master(), "user", query, data, originalData, nil)
	result = w.runBeforeTrigger()
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	expectData = types.M{
		"username": "joe",
		"key":      "hello",
	}
	if reflect.DeepEqual(expectData, w.data) == false {
		t.Error("expect:", expectData, "result:", w.data)
	}
	if reflect.DeepEqual(nil, w.storage["changedByTrigger"]) == false {
		t.Error("expect:", nil, "result:", w.storage["changedByTrigger"])
	}
	cloud.UnregisterAll()
	/***************************************************************/
	cloud.BeforeSave("user", func(request cloud.TriggerRequest, response cloud.Response) {
		object := request.Object
		if username := utils.S(object["username"]); username != "" {
			object["username"] = username + "_tomato"
			response.Success(nil)
		} else {
			response.Error(1, "need a username")
		}
	})
	query = types.M{"objectId": "1001"}
	data = types.M{
		"key": "hello",
	}
	originalData = nil
	w, _ = NewWrite(Master(), "user", query, data, originalData, nil)
	result = w.runBeforeTrigger()
	expect = errs.E(1, "need a username")
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	cloud.UnregisterAll()
	/***************************************************************/
	cloud.BeforeSave("user", func(request cloud.TriggerRequest, response cloud.Response) {
		object := request.Object
		if username := utils.S(object["username"]); username != "" {
			object["username"] = username + "_tomato"
			response.Success(nil)
		} else {
			response.Error(1, "need a username")
		}
	})
	query = types.M{"objectId": "1001"}
	data = types.M{
		"username": "joe",
		"key":      "hello",
	}
	originalData = nil
	w, _ = NewWrite(Master(), "user", query, data, originalData, nil)
	result = w.runBeforeTrigger()
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	expectData = types.M{
		"username": "joe_tomato",
		"key":      "hello",
	}
	if reflect.DeepEqual(expectData, w.data) == false {
		t.Error("expect:", expectData, "result:", w.data)
	}
	if reflect.DeepEqual(true, w.storage["changedByTrigger"]) == false {
		t.Error("expect:", true, "result:", w.storage["changedByTrigger"])
	}
	cloud.UnregisterAll()
	/***************************************************************/
	cloud.BeforeSave("user", func(request cloud.TriggerRequest, response cloud.Response) {
		object := request.Object
		if username := utils.S(object["username"]); username != "" {
			response.Success(nil)
		} else {
			response.Error(1, "need a username")
		}
	})
	query = types.M{"objectId": "1001"}
	data = types.M{
		"username": "joe",
		"key":      "hello",
	}
	originalData = nil
	w, _ = NewWrite(Master(), "user", query, data, originalData, nil)
	result = w.runBeforeTrigger()
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	expectData = types.M{
		"username": "joe",
		"key":      "hello",
	}
	if reflect.DeepEqual(expectData, w.data) == false {
		t.Error("expect:", expectData, "result:", w.data)
	}
	if reflect.DeepEqual(nil, w.storage["changedByTrigger"]) == false {
		t.Error("expect:", nil, "result:", w.storage["changedByTrigger"])
	}
	cloud.UnregisterAll()
}

func Test_setRequiredFieldsIfNeeded(t *testing.T) {
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	var expect types.M
	timeStr := utils.TimetoString(time.Now().UTC())
	/***************************************************************/
	query = types.M{"objectId": "1001"}
	data = types.M{"key": "hello"}
	originalData = types.M{}
	w, _ = NewWrite(Master(), "user", query, data, originalData, nil)
	w.updatedAt = timeStr
	w.setRequiredFieldsIfNeeded()
	expect = types.M{
		"key":       "hello",
		"updatedAt": timeStr,
	}
	if reflect.DeepEqual(expect, w.data) == false {
		t.Error("expect:", expect, "result:", w.data)
	}
	/***************************************************************/
	query = nil
	data = types.M{
		"key": "hello",
	}
	originalData = nil
	w, _ = NewWrite(Master(), "user", query, data, originalData, nil)
	w.updatedAt = timeStr
	w.setRequiredFieldsIfNeeded()
	expect = types.M{
		"key":       "hello",
		"updatedAt": timeStr,
		"createdAt": timeStr,
	}
	if w.data["objectId"] == nil {
		t.Error("expect:", "objectId", "result:", nil)
	}
	delete(w.data, "objectId")
	if reflect.DeepEqual(expect, w.data) == false {
		t.Error("expect:", expect, "result:", w.data)
	}
}

func Test_transformUser(t *testing.T) {
	var schema, object types.M
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	var expect types.M
	var err, expectErr error
	/***************************************************************/
	query = nil
	data = types.M{}
	originalData = nil
	w, _ = NewWrite(Master(), "_User", query, data, originalData, nil)
	w.transformUser()
	if v, ok := w.data["username"]; ok == false {
		t.Error("expect:", "username", "result:", v)
	}
	/***************************************************************/
	query = nil
	data = types.M{
		"password": "123456",
	}
	originalData = nil
	w, _ = NewWrite(Master(), "_User", query, data, originalData, nil)
	w.transformUser()
	expect = types.M{
		"_hashed_password": utils.Hash("123456"),
	}
	if v, ok := w.data["username"]; ok == false {
		t.Error("expect:", "username", "result:", v)
	}
	delete(w.data, "username")
	if reflect.DeepEqual(expect, w.data) == false {
		t.Error("expect:", expect, "result:", w.data)
	}
	/***************************************************************/
	initEnv()
	query = nil
	data = types.M{
		"username": "joe",
		"password": "123456",
	}
	originalData = nil
	w, _ = NewWrite(Master(), "_User", query, data, originalData, nil)
	w.data["objectId"] = "1001"
	w.transformUser()
	expect = types.M{
		"objectId":         "1001",
		"username":         "joe",
		"_hashed_password": utils.Hash("123456"),
	}
	if reflect.DeepEqual(expect, w.data) == false {
		t.Error("expect:", expect, "result:", w.data)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId": "1002",
		"username": "joe",
	}
	orm.Adapter.CreateObject("_User", schema, object)
	query = nil
	data = types.M{
		"username": "joe",
		"password": "123456",
	}
	originalData = nil
	w, _ = NewWrite(Master(), "_User", query, data, originalData, nil)
	w.data["objectId"] = "1001"
	err = w.transformUser()
	expectErr = errs.E(errs.UsernameTaken, "Account already exists for this username")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	query = nil
	data = types.M{
		"username": "joe",
		"password": "123456",
		"email":    "hello",
	}
	originalData = nil
	w, _ = NewWrite(Master(), "_User", query, data, originalData, nil)
	w.data["objectId"] = "1001"
	err = w.transformUser()
	expectErr = errs.E(errs.InvalidEmailAddress, "Email address format is invalid.")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
			"email":    types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId": "1002",
		"username": "jack",
		"email":    "a@g.cn",
	}
	orm.Adapter.CreateObject("_User", schema, object)
	query = nil
	data = types.M{
		"username": "joe",
		"password": "123456",
		"email":    "a@g.cn",
	}
	originalData = nil
	w, _ = NewWrite(Master(), "_User", query, data, originalData, nil)
	w.data["objectId"] = "1001"
	err = w.transformUser()
	expectErr = errs.E(errs.EmailTaken, "Account already exists for this email address")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	config.TConfig.VerifyUserEmails = false
	query = nil
	data = types.M{
		"username": "joe",
		"password": "123456",
		"email":    "a@g.cn",
	}
	originalData = nil
	w, _ = NewWrite(Master(), "_User", query, data, originalData, nil)
	w.data["objectId"] = "1001"
	w.transformUser()
	expect = types.M{
		"objectId":         "1001",
		"username":         "joe",
		"_hashed_password": utils.Hash("123456"),
		"email":            "a@g.cn",
	}
	if reflect.DeepEqual(true, w.storage["sendVerificationEmail"]) == false {
		t.Error("expect:", true, "result:", w.storage["sendVerificationEmail"])
	}
	if reflect.DeepEqual(expect, w.data) == false {
		t.Error("expect:", expect, "result:", w.data)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	config.TConfig.VerifyUserEmails = true
	config.TConfig.EmailVerifyTokenValidityDuration = 180
	query = nil
	data = types.M{
		"username": "joe",
		"password": "123456",
		"email":    "a@g.cn",
	}
	originalData = nil
	w, _ = NewWrite(Master(), "_User", query, data, originalData, nil)
	w.data["objectId"] = "1001"
	w.transformUser()
	expect = types.M{
		"objectId":                       "1001",
		"username":                       "joe",
		"_hashed_password":               utils.Hash("123456"),
		"email":                          "a@g.cn",
		"emailVerified":                  false,
		"_email_verify_token_expires_at": utils.TimetoString(time.Now().UTC().Add(180 * time.Second)),
	}
	if reflect.DeepEqual(true, w.storage["sendVerificationEmail"]) == false {
		t.Error("expect:", true, "result:", w.storage["sendVerificationEmail"])
	}
	if v, ok := w.data["_email_verify_token"]; ok == false {
		t.Error("expect:", "need _email_verify_token", "result:", v)
	}
	delete(w.data, "_email_verify_token")
	if reflect.DeepEqual(expect, w.data) == false {
		t.Error("expect:", expect, "result:", w.data)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	schema = types.M{
		"fields": types.M{
			"user":         types.M{"type": "Pointer", "targetClass": "_User"},
			"sessionToken": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_Session", schema)
	object = types.M{
		"objectId": "2001",
		"user": types.M{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  "1001",
		},
		"sessionToken": "aaaaa",
	}
	orm.Adapter.CreateObject("_Session", schema, object)
	cache.User.Put("aaaaa", "user1", 0)
	query = types.M{"objectId": "1001"}
	data = types.M{
		"password": "123456",
	}
	originalData = types.M{}
	w, _ = NewWrite(&Auth{IsMaster: false, User: types.M{"objectId": "1001"}}, "_User", query, data, originalData, nil)
	err = w.transformUser()
	expect = types.M{
		"_hashed_password": utils.Hash("123456"),
	}
	if cache.User.Get("aaaaa") != nil {
		t.Error("expect:", nil, "result:", cache.User.Get("aaaaa"))
	}
	if reflect.DeepEqual(true, w.storage["clearSessions"]) == false {
		t.Error("expect:", true, "result:", w.storage["clearSessions"])
	}
	if reflect.DeepEqual(true, w.storage["generateNewSession"]) == false {
		t.Error("expect:", true, "result:", w.storage["generateNewSession"])
	}
	if reflect.DeepEqual(expect, w.data) == false {
		t.Error("expect:", expect, "result:", w.data, err)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_expandFilesForExistingObjects(t *testing.T) {
	// TODO 展开文件类型
}

func Test_runDatabaseOperation(t *testing.T) {
	var schema, object types.M
	var w *Write
	var className string
	var auth *Auth
	var query, data, originalData types.M
	var expect types.M
	var err, expectErr error
	var timeStr = utils.TimetoString(time.Now().UTC())
	var results types.S
	/***************************************************************/
	className = "_User"
	auth = Nobody()
	query = types.M{"objectId": "1001"}
	data = types.M{"key": "hello"}
	originalData = types.M{}
	w, _ = NewWrite(auth, className, query, data, originalData, nil)
	err = w.runDatabaseOperation()
	expectErr = errs.E(errs.SessionMissing, "cannot modify user 1001")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/***************************************************************/
	className = "_User"
	auth = &Auth{
		IsMaster: false,
		User:     types.M{"objectId": "1002"},
	}
	query = types.M{"objectId": "1001"}
	data = types.M{"key": "hello"}
	originalData = types.M{}
	w, _ = NewWrite(auth, className, query, data, originalData, nil)
	err = w.runDatabaseOperation()
	expectErr = errs.E(errs.SessionMissing, "cannot modify user 1001")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/***************************************************************/
	className = "_User"
	auth = &Auth{
		IsMaster: false,
		User:     types.M{"objectId": "1001"},
	}
	query = types.M{"objectId": "1001"}
	data = types.M{
		"key": "hello",
		"ACL": types.M{"*unresolved": "hello"},
	}
	originalData = types.M{}
	w, _ = NewWrite(auth, className, query, data, originalData, nil)
	err = w.runDatabaseOperation()
	expectErr = errs.E(errs.InvalidAcl, "Invalid ACL.")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/***************************************************************/
	className = "_User"
	auth = Master()
	query = types.M{"objectId": "1001"}
	data = types.M{
		"key": "hello",
		"ACL": types.M{"*unresolved": "hello"},
	}
	originalData = types.M{}
	w, _ = NewWrite(auth, className, query, data, originalData, nil)
	err = w.runDatabaseOperation()
	expectErr = errs.E(errs.InvalidAcl, "Invalid ACL.")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/***************************************************************/
	initEnv()
	className = "_User"
	auth = Master()
	query = types.M{"objectId": "1001"}
	data = types.M{
		"key": "hello",
		"ACL": types.M{},
	}
	originalData = types.M{}
	w, _ = NewWrite(auth, className, query, data, originalData, nil)
	err = w.runDatabaseOperation()
	expectErr = errs.E(errs.ObjectNotFound, "Object not found.")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId": "1001",
		"username": "joe",
	}
	orm.Adapter.CreateObject("_User", schema, object)
	className = "_User"
	auth = Master()
	query = types.M{"objectId": "1001"}
	data = types.M{
		"key": "hello",
		"ACL": types.M{},
	}
	originalData = types.M{}
	w, _ = NewWrite(auth, className, query, data, originalData, nil)
	w.updatedAt = timeStr
	err = w.runDatabaseOperation()
	expect = types.M{"updatedAt": timeStr}
	if err != nil || reflect.DeepEqual(expect, w.response["response"]) == false {
		t.Error("expect:", expect, "result:", w.response["response"], err)
	}
	results, _ = orm.TomatoDBController.Find(className, types.M{}, types.M{})
	if len(results) != 1 {
		t.Error("expect:", "1", "result:", results)
	} else {
		expect = types.M{
			"objectId": "1001",
			"username": "joe",
			"key":      "hello",
			"ACL": types.M{
				"1001": types.M{
					"read":  true,
					"write": true,
				},
			},
		}
		if reflect.DeepEqual(expect, results[0]) == false {
			t.Error("expect:", expect, "result:", results[0])
		}
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("user", schema)
	object = types.M{
		"objectId": "1001",
		"username": "joe",
	}
	orm.Adapter.CreateObject("user", schema, object)
	className = "user"
	auth = Master()
	query = types.M{"objectId": "1001"}
	data = types.M{
		"key": "hello",
	}
	originalData = types.M{}
	w, _ = NewWrite(auth, className, query, data, originalData, nil)
	w.updatedAt = timeStr
	w.storage["changedByTrigger"] = true
	err = w.runDatabaseOperation()
	expect = types.M{
		"key":       "hello",
		"updatedAt": timeStr,
	}
	if err != nil || reflect.DeepEqual(expect, w.response["response"]) == false {
		t.Error("expect:", expect, "result:", w.response["response"], err)
	}
	results, _ = orm.TomatoDBController.Find(className, types.M{}, types.M{})
	if len(results) != 1 {
		t.Error("expect:", "1", "result:", results)
	} else {
		expect = types.M{
			"objectId": "1001",
			"username": "joe",
			"key":      "hello",
		}
		if reflect.DeepEqual(expect, results[0]) == false {
			t.Error("expect:", expect, "result:", results[0])
		}
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	className = "_User"
	auth = Master()
	query = nil
	data = types.M{
		"username": "joe",
	}
	originalData = nil
	w, _ = NewWrite(auth, className, query, data, originalData, nil)
	w.data["objectId"] = "1001"
	w.data["createdAt"] = timeStr
	config.TConfig.ServerURL = "http://127.0.0.1/v1"
	err = w.runDatabaseOperation()
	expect = types.M{
		"status": 201,
		"response": types.M{
			"objectId":  "1001",
			"createdAt": timeStr,
		},
		"location": "http://127.0.0.1/v1/users/1001",
	}
	if reflect.DeepEqual(expect, w.response) == false {
		t.Error("expect:", expect, "result:", w.response)
	}
	results, _ = orm.TomatoDBController.Find(className, types.M{}, types.M{})
	if len(results) != 1 {
		t.Error("expect:", "1", "result:", results)
	} else {
		expect = types.M{
			"objectId":  "1001",
			"username":  "joe",
			"createdAt": timeStr,
			"ACL": types.M{
				"1001": types.M{
					"read":  true,
					"write": true,
				},
				"*": types.M{
					"read": true,
				},
			},
		}
		if reflect.DeepEqual(expect, results[0]) == false {
			t.Error("expect:", expect, "result:", results[0])
		}
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	className = "user"
	auth = Master()
	query = nil
	data = types.M{
		"username": "joe",
	}
	originalData = nil
	w, _ = NewWrite(auth, className, query, data, originalData, nil)
	w.data["objectId"] = "1001"
	w.data["createdAt"] = timeStr
	config.TConfig.ServerURL = "http://127.0.0.1/v1"
	w.storage["changedByTrigger"] = true
	err = w.runDatabaseOperation()
	expect = types.M{
		"status": 201,
		"response": types.M{
			"objectId":  "1001",
			"username":  "joe",
			"createdAt": timeStr,
		},
		"location": "http://127.0.0.1/v1/classes/user/1001",
	}
	if reflect.DeepEqual(expect, w.response) == false {
		t.Error("expect:", expect, "result:", w.response)
	}
	results, _ = orm.TomatoDBController.Find(className, types.M{}, types.M{})
	if len(results) != 1 {
		t.Error("expect:", "1", "result:", results)
	} else {
		expect = types.M{
			"objectId":  "1001",
			"username":  "joe",
			"createdAt": timeStr,
		}
		if reflect.DeepEqual(expect, results[0]) == false {
			t.Error("expect:", expect, "result:", results[0])
		}
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId": "1001",
		"username": "joe",
	}
	orm.Adapter.CreateObject("_User", schema, object)
	className = "_User"
	auth = Master()
	query = nil
	data = types.M{
		"username": "joe",
	}
	originalData = nil
	w, _ = NewWrite(auth, className, query, data, originalData, nil)
	w.data["objectId"] = "1001"
	w.data["createdAt"] = timeStr
	config.TConfig.ServerURL = "http://127.0.0.1/v1"
	err = w.runDatabaseOperation()
	expectErr = errs.E(errs.DuplicateValue, "A duplicate value for a field with unique values was provided")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_createSessionTokenIfNeeded(t *testing.T) {
	// createSessionToken
	// TODO
}

func Test_handleFollowup(t *testing.T) {
	// createSessionToken
	// TODO
}

func Test_runAfterTrigger(t *testing.T) {
	var className string
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	/***************************************************************/
	config.TConfig.LiveQuery = livequery.NewLiveQuery([]string{}, "", "")
	className = "user"
	query = nil
	data = types.M{"key": "hello"}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	w.response = types.M{"response": "hello"}
	w.runAfterTrigger()
	/***************************************************************/
	config.TConfig.LiveQuery = livequery.NewLiveQuery([]string{}, "", "")
	cloud.AfterSave("user", func(res cloud.TriggerRequest, resp cloud.Response) {
	})
	className = "user"
	query = nil
	data = types.M{"key": "hello"}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	w.response = types.M{"response": "hello"}
	w.runAfterTrigger()
	cloud.UnregisterAll()
}

func Test_cleanUserAuthData(t *testing.T) {
	// TODO
}

/////////////////////////////////////////////////////////////

func Test_handleAuthData(t *testing.T) {
	var className string
	var schema types.M
	var object types.M
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	var result error
	var expect error
	var response types.M
	var location string
	/***************************************************************/
	className = "_User"
	query = nil
	data = types.M{
		"authData": types.M{
			"other": types.M{
				"id": "1001",
			},
		},
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.handleAuthData(utils.M(w.data["authData"]))
	expect = errs.E(errs.UnsupportedService, "This authentication method is unsupported.")
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/***************************************************************/
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	object = types.M{
		"objectId": "102",
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	className = "_User"
	query = nil
	data = types.M{
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.handleAuthData(utils.M(w.data["authData"]))
	expect = errs.E(errs.AccountAlreadyLinked, "this auth is already used")
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	className = "_User"
	query = nil
	data = types.M{
		"authData": types.M{
			"facebook": types.M{
				"id": "1002",
			},
		},
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.handleAuthData(utils.M(w.data["authData"]))
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	if reflect.DeepEqual("facebook", w.storage["authProvider"]) == false {
		t.Error("expect:", "facebook", "result:", w.storage["authProvider"])
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	config.TConfig = &config.Config{
		ServerURL: "http://www.g.cn",
	}
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	className = "_User"
	query = nil
	data = types.M{
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.handleAuthData(utils.M(w.data["authData"]))
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	response = types.M{
		"objectId": "101",
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	if reflect.DeepEqual(utils.M(response), w.response["response"]) == false {
		t.Error("expect:", response, "result:", w.response["response"])
	}
	location = "http://www.g.cn/users/101"
	if reflect.DeepEqual(location, w.response["location"]) == false {
		t.Error("expect:", location, "result:", w.response["location"])
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	config.TConfig = &config.Config{
		ServerURL: "http://www.g.cn",
	}
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"facebook": types.M{
				"id":    "1001",
				"token": "aaa",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	className = "_User"
	query = nil
	data = types.M{
		"authData": types.M{
			"facebook": types.M{
				"id":    "1001",
				"token": "abc",
			},
		},
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.handleAuthData(utils.M(w.data["authData"]))
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	response = types.M{
		"objectId": "101",
		"authData": map[string]interface{}{
			"facebook": types.M{
				"id":    "1001",
				"token": "abc",
			},
		},
	}
	if reflect.DeepEqual(utils.M(response), w.response["response"]) == false {
		t.Error("expect:", response, "result:", w.response["response"])
	}
	location = "http://www.g.cn/users/101"
	if reflect.DeepEqual(location, w.response["location"]) == false {
		t.Error("expect:", location, "result:", w.response["location"])
	}
	r, _ := orm.TomatoDBController.Find(className, types.M{"objectId": "101"}, types.M{})
	response = types.M{
		"objectId": "101",
		"authData": types.M{
			"facebook": types.M{
				"id":    "1001",
				"token": "abc",
			},
		},
	}
	if reflect.DeepEqual(response, r[0]) == false {
		t.Error("expect:", response, "result:", r[0])
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	config.TConfig = &config.Config{
		ServerURL: "http://www.g.cn",
	}
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"facebook": types.M{
				"id":    "1001",
				"token": "aaa",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	className = "_User"
	query = types.M{"objectId": "101"}
	data = types.M{
		"authData": types.M{
			"facebook": types.M{
				"id":    "1001",
				"token": "aaa",
			},
		},
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.handleAuthData(utils.M(w.data["authData"]))
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	config.TConfig = &config.Config{
		ServerURL: "http://www.g.cn",
	}
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"facebook": types.M{
				"id":    "1001",
				"token": "aaa",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	className = "_User"
	query = types.M{"objectId": "102"}
	data = types.M{
		"authData": types.M{
			"facebook": types.M{
				"id":    "1001",
				"token": "aaa",
			},
		},
	}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	result = w.handleAuthData(utils.M(w.data["authData"]))
	expect = errs.E(errs.AccountAlreadyLinked, "this auth is already used")
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_handleAuthDataValidation(t *testing.T) {
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	var authData types.M
	var result error
	var expect error
	/***************************************************************/
	query = nil
	data = types.M{}
	originalData = nil
	w, _ = NewWrite(Master(), "user", query, data, originalData, nil)
	authData = types.M{
		"facebook": types.M{
			"id": "1001",
		},
	}
	result = w.handleAuthDataValidation(authData)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/***************************************************************/
	query = nil
	data = types.M{}
	originalData = nil
	w, _ = NewWrite(Master(), "user", query, data, originalData, nil)
	authData = types.M{
		"other": types.M{
			"id": "1001",
		},
	}
	result = w.handleAuthDataValidation(authData)
	expect = errs.E(errs.UnsupportedService, "This authentication method is unsupported.")
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_findUsersWithAuthData(t *testing.T) {
	var schema types.M
	var object types.M
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	var className string
	var authData types.M
	var result types.S
	var err error
	var expect types.S
	/***************************************************************/
	initEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	query = nil
	data = types.M{}
	originalData = nil
	className = "user"
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	authData = types.M{
		"facebook": types.M{
			"id": "1002",
		},
	}
	result, err = w.findUsersWithAuthData(authData)
	expect = types.S{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	query = nil
	data = types.M{}
	originalData = nil
	className = "_User"
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	authData = types.M{
		"facebook": types.M{
			"id": "1001",
		},
	}
	result, err = w.findUsersWithAuthData(authData)
	expect = types.S{
		types.M{
			"objectId": "101",
			"authData": types.M{
				"facebook": types.M{
					"id": "1001",
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	object = types.M{
		"objectId": "102",
		"authData": types.M{
			"twitter": types.M{
				"id": "1001",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	query = nil
	data = types.M{}
	originalData = nil
	className = "_User"
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	authData = types.M{
		"facebook": types.M{
			"id": "1001",
		},
		"twitter": types.M{
			"id": "1001",
		},
	}
	result, err = w.findUsersWithAuthData(authData)
	expect = types.S{
		types.M{
			"objectId": "101",
			"authData": types.M{
				"facebook": types.M{
					"id": "1001",
				},
			},
		},
		types.M{
			"objectId": "102",
			"authData": types.M{
				"twitter": types.M{
					"id": "1001",
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_createSessionToken(t *testing.T) {
	// Execute
	// TODO
}

func Test_location(t *testing.T) {
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	var className string
	var result string
	var expect string
	/***************************************************************/
	config.TConfig = &config.Config{
		ServerURL: "http://www.g.cn",
	}
	query = nil
	data = types.M{}
	originalData = nil
	className = "post"
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	w.data["objectId"] = "1001"
	result = w.location()
	expect = "http://www.g.cn/classes/post/1001"
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/***************************************************************/
	config.TConfig = &config.Config{
		ServerURL: "http://www.g.cn",
	}
	query = nil
	data = types.M{}
	originalData = nil
	className = "_User"
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	w.data["objectId"] = "1001"
	result = w.location()
	expect = "http://www.g.cn/users/1001"
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_objectID(t *testing.T) {
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	var result interface{}
	var expect interface{}
	/***************************************************************/
	query = nil
	data = types.M{"key": "hello"}
	originalData = nil
	w, _ = NewWrite(Master(), "user", query, data, originalData, nil)
	w.data["objectId"] = "1001"
	result = w.objectID()
	expect = "1001"
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/***************************************************************/
	query = types.M{"objectId": "1001"}
	data = types.M{"key": "hello"}
	originalData = types.M{"key": "hi"}
	w, _ = NewWrite(Master(), "user", query, data, originalData, nil)
	result = w.objectID()
	expect = "1001"
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_sanitizedData(t *testing.T) {
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	var result types.M
	var expect types.M
	/***************************************************************/
	query = nil
	data = types.M{
		"key":              "hello",
		"_auth_data":       "facebook",
		"_hashed_password": "123456",
	}
	originalData = nil
	w, _ = NewWrite(Master(), "user", query, data, originalData, nil)
	result = w.sanitizedData()
	expect = types.M{
		"key": "hello",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_updateResponseWithData(t *testing.T) {
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	var response, updateData types.M
	var result types.M
	var expect types.M
	/***************************************************************/
	query = nil
	data = types.M{}
	originalData = nil
	w, _ = NewWrite(Master(), "user", query, data, originalData, nil)
	response = types.M{
		"key":  "hello",
		"key1": 10,
	}
	updateData = types.M{
		"key1": types.M{
			"__op": "Increment",
		},
		"key2": "world",
		"key3": types.M{
			"__op": "Delete",
		},
	}
	result = w.updateResponseWithData(response, updateData)
	expect = types.M{
		"key":  "hello",
		"key1": 10,
		"key2": "world",
		"key3": types.M{
			"__op": "Delete",
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}
