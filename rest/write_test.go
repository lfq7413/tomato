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
	var className string
	var w *Write
	var auth *Auth
	var query types.M
	var data types.M
	var originalData types.M
	var err error
	var result types.M
	var results types.S
	/***************************************************************/
	initEnv()
	className = "user"
	auth = Master()
	query = nil
	data = types.M{"username": "joe"}
	originalData = nil
	w, err = NewWrite(auth, className, query, data, originalData, nil)
	result, err = w.Execute()
	if err != nil || result == nil {
		t.Error("expect:", nil, "result:", result, err)
	}
	results, _ = orm.TomatoDBController.Find(className, types.M{}, types.M{})
	if len(results) != 1 {
		t.Error("expect:", "len 1", "result:", results)
	}
	id := utils.M(result["response"])["objectId"]
	auth = Master()
	query = types.M{"objectId": id}
	data = types.M{"username": "jack"}
	originalData = types.M{"username": "joe"}
	w, err = NewWrite(auth, className, query, data, originalData, nil)
	result, err = w.Execute()
	if err != nil || result == nil {
		t.Error("expect:", nil, "result:", result, err)
	}
	results, _ = orm.TomatoDBController.Find(className, types.M{}, types.M{})
	if len(results) != 1 {
		t.Error("expect:", "len 1", "result:", results)
	}
	orm.TomatoDBController.DeleteEverything()
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
	var schema, object types.M
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	var err, expectErr error
	var expectQuery types.M
	var results, expect types.S
	/***************************************************************/
	query = nil
	data = types.M{}
	originalData = nil
	w, _ = NewWrite(Master(), "_Installation", query, data, originalData, nil)
	err = w.handleInstallation()
	expectErr = errs.E(errs.MissingRequiredFieldError, "at least one ID field (deviceToken, installationId) must be specified in this operation")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/***************************************************************/
	query = types.M{"objectId": "123"}
	data = types.M{}
	originalData = nil
	w, _ = NewWrite(Master(), "_Installation", query, data, originalData, nil)
	err = w.handleInstallation()
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/***************************************************************/
	query = nil
	data = types.M{"deviceToken": "abc"}
	originalData = nil
	w, _ = NewWrite(Master(), "_Installation", query, data, originalData, nil)
	err = w.handleInstallation()
	expectErr = errs.E(errs.MissingRequiredFieldError, "deviceType must be specified in this operation")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/***************************************************************/
	query = nil
	data = types.M{"installationId": "abc"}
	originalData = nil
	w, _ = NewWrite(Master(), "_Installation", query, data, originalData, nil)
	err = w.handleInstallation()
	expectErr = errs.E(errs.MissingRequiredFieldError, "deviceType must be specified in this operation")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/***************************************************************/
	initEnv()
	schema = types.M{
		"fields": types.M{
			"installationId": types.M{"type": "String"},
			"deviceToken":    types.M{"type": "String"},
			"deviceType":     types.M{"type": "String"},
			"appIdentifier":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_Installation", schema)
	object = types.M{
		"objectId":       "1001",
		"installationId": "123",
		"deviceToken":    "abc",
		"deviceType":     "ios",
	}
	orm.Adapter.CreateObject("_Installation", schema, object)
	query = types.M{"objectId": "1002"}
	data = types.M{"installationId": "abc"}
	originalData = types.M{}
	w, _ = NewWrite(Master(), "_Installation", query, data, originalData, nil)
	err = w.handleInstallation()
	expectErr = errs.E(errs.ObjectNotFound, "Object not found for update.")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	schema = types.M{
		"fields": types.M{
			"installationId": types.M{"type": "String"},
			"deviceToken":    types.M{"type": "String"},
			"deviceType":     types.M{"type": "String"},
			"appIdentifier":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_Installation", schema)
	object = types.M{
		"objectId":       "1001",
		"installationId": "123",
		"deviceToken":    "abc",
		"deviceType":     "ios",
	}
	orm.Adapter.CreateObject("_Installation", schema, object)
	query = types.M{"objectId": "1001"}
	data = types.M{"installationId": "abc"}
	originalData = types.M{}
	w, _ = NewWrite(Master(), "_Installation", query, data, originalData, nil)
	err = w.handleInstallation()
	expectErr = errs.E(errs.ChangedImmutableFieldError, "installationId may not be changed in this operation")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	schema = types.M{
		"fields": types.M{
			"installationId": types.M{"type": "String"},
			"deviceToken":    types.M{"type": "String"},
			"deviceType":     types.M{"type": "String"},
			"appIdentifier":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_Installation", schema)
	object = types.M{
		"objectId":    "1001",
		"deviceToken": "abc",
		"deviceType":  "ios",
	}
	orm.Adapter.CreateObject("_Installation", schema, object)
	query = types.M{"objectId": "1001"}
	data = types.M{"deviceToken": "123"}
	originalData = types.M{}
	w, _ = NewWrite(Master(), "_Installation", query, data, originalData, nil)
	err = w.handleInstallation()
	expectErr = errs.E(errs.ChangedImmutableFieldError, "deviceToken may not be changed in this operation")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	schema = types.M{
		"fields": types.M{
			"installationId": types.M{"type": "String"},
			"deviceToken":    types.M{"type": "String"},
			"deviceType":     types.M{"type": "String"},
			"appIdentifier":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_Installation", schema)
	object = types.M{
		"objectId":    "1001",
		"deviceToken": "abc",
		"deviceType":  "ios",
	}
	orm.Adapter.CreateObject("_Installation", schema, object)
	query = types.M{"objectId": "1001"}
	data = types.M{"deviceType": "android"}
	originalData = types.M{}
	w, _ = NewWrite(Master(), "_Installation", query, data, originalData, nil)
	err = w.handleInstallation()
	expectErr = errs.E(errs.ChangedImmutableFieldError, "deviceType may not be changed in this operation")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	schema = types.M{
		"fields": types.M{
			"installationId": types.M{"type": "String"},
			"deviceToken":    types.M{"type": "String"},
			"deviceType":     types.M{"type": "String"},
			"appIdentifier":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_Installation", schema)
	object = types.M{
		"objectId":       "1001",
		"installationId": "111",
		"deviceToken":    "aaa",
		"deviceType":     "ios",
	}
	orm.Adapter.CreateObject("_Installation", schema, object)
	query = nil
	data = types.M{"installationId": "222", "deviceType": "android"}
	originalData = nil
	w, _ = NewWrite(Master(), "_Installation", query, data, originalData, nil)
	err = w.handleInstallation()
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	expectQuery = nil
	if reflect.DeepEqual(expectQuery, w.query) == false {
		t.Error("expect:", expectQuery, "result:", w.query)
	}
	results, err = orm.TomatoDBController.Find("_Installation", types.M{}, types.M{})
	expect = types.S{
		types.M{
			"objectId":       "1001",
			"installationId": "111",
			"deviceToken":    "aaa",
			"deviceType":     "ios",
		},
	}
	if err != nil || reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	schema = types.M{
		"fields": types.M{
			"installationId": types.M{"type": "String"},
			"deviceToken":    types.M{"type": "String"},
			"deviceType":     types.M{"type": "String"},
			"appIdentifier":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_Installation", schema)
	object = types.M{
		"objectId":       "1001",
		"installationId": "111",
		"deviceToken":    "aaa",
		"deviceType":     "ios",
	}
	orm.Adapter.CreateObject("_Installation", schema, object)
	query = nil
	data = types.M{"deviceToken": "aaa", "deviceType": "android"}
	originalData = nil
	w, _ = NewWrite(Master(), "_Installation", query, data, originalData, nil)
	err = w.handleInstallation()
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	expectQuery = types.M{"objectId": "1001"}
	if reflect.DeepEqual(expectQuery, w.query) == false {
		t.Error("expect:", expectQuery, "result:", w.query)
	}
	results, err = orm.TomatoDBController.Find("_Installation", types.M{}, types.M{})
	expect = types.S{
		types.M{
			"objectId":       "1001",
			"installationId": "111",
			"deviceToken":    "aaa",
			"deviceType":     "ios",
		},
	}
	if err != nil || reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	schema = types.M{
		"fields": types.M{
			"installationId": types.M{"type": "String"},
			"deviceToken":    types.M{"type": "String"},
			"deviceType":     types.M{"type": "String"},
			"appIdentifier":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_Installation", schema)
	object = types.M{
		"objectId":       "1001",
		"installationId": "111",
		"deviceToken":    "aaa",
		"deviceType":     "ios",
	}
	orm.Adapter.CreateObject("_Installation", schema, object)
	object = types.M{
		"objectId":       "1002",
		"installationId": "222",
		"deviceToken":    "aaa",
		"deviceType":     "ios",
	}
	orm.Adapter.CreateObject("_Installation", schema, object)
	query = nil
	data = types.M{"deviceToken": "aaa", "deviceType": "android"}
	originalData = nil
	w, _ = NewWrite(Master(), "_Installation", query, data, originalData, nil)
	err = w.handleInstallation()
	expectErr = errs.E(errs.InvalidInstallationIDError, "Must specify installationId when deviceToken matches multiple Installation objects")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	schema = types.M{
		"fields": types.M{
			"installationId": types.M{"type": "String"},
			"deviceToken":    types.M{"type": "String"},
			"deviceType":     types.M{"type": "String"},
			"appIdentifier":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_Installation", schema)
	object = types.M{
		"objectId":       "1001",
		"installationId": "111",
		"deviceToken":    "aaa",
		"deviceType":     "ios",
	}
	orm.Adapter.CreateObject("_Installation", schema, object)
	object = types.M{
		"objectId":       "1002",
		"installationId": "222",
		"deviceToken":    "aaa",
		"deviceType":     "ios",
	}
	orm.Adapter.CreateObject("_Installation", schema, object)
	query = nil
	data = types.M{"deviceToken": "aaa", "installationId": "333", "deviceType": "android"}
	originalData = nil
	w, _ = NewWrite(Master(), "_Installation", query, data, originalData, nil)
	err = w.handleInstallation()
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	expectQuery = nil
	if reflect.DeepEqual(expectQuery, w.query) == false {
		t.Error("expect:", expectQuery, "result:", w.query)
	}
	results, err = orm.TomatoDBController.Find("_Installation", types.M{}, types.M{})
	expect = types.S{}
	if err != nil || reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	schema = types.M{
		"fields": types.M{
			"installationId": types.M{"type": "String"},
			"deviceToken":    types.M{"type": "String"},
			"deviceType":     types.M{"type": "String"},
			"appIdentifier":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_Installation", schema)
	object = types.M{
		"objectId":       "1001",
		"installationId": "111",
		"deviceToken":    "bbb",
		"deviceType":     "ios",
	}
	orm.Adapter.CreateObject("_Installation", schema, object)
	object = types.M{
		"objectId":    "1002",
		"deviceToken": "aaa",
		"deviceType":  "ios",
	}
	orm.Adapter.CreateObject("_Installation", schema, object)
	query = nil
	data = types.M{"deviceToken": "aaa", "installationId": "111", "deviceType": "android"}
	originalData = nil
	w, _ = NewWrite(Master(), "_Installation", query, data, originalData, nil)
	err = w.handleInstallation()
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	expectQuery = types.M{"objectId": "1002"}
	if reflect.DeepEqual(expectQuery, w.query) == false {
		t.Error("expect:", expectQuery, "result:", w.query)
	}
	results, err = orm.TomatoDBController.Find("_Installation", types.M{}, types.M{})
	expect = types.S{
		types.M{
			"objectId":    "1002",
			"deviceToken": "aaa",
			"deviceType":  "ios",
		},
	}
	if err != nil || reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	schema = types.M{
		"fields": types.M{
			"installationId": types.M{"type": "String"},
			"deviceToken":    types.M{"type": "String"},
			"deviceType":     types.M{"type": "String"},
			"appIdentifier":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_Installation", schema)
	object = types.M{
		"objectId":       "1001",
		"installationId": "111",
		"deviceToken":    "bbb",
		"deviceType":     "ios",
	}
	orm.Adapter.CreateObject("_Installation", schema, object)
	object = types.M{
		"objectId":    "1002",
		"deviceToken": "aaa",
		"deviceType":  "ios",
	}
	orm.Adapter.CreateObject("_Installation", schema, object)
	query = nil
	data = types.M{"installationId": "111", "deviceType": "android"}
	originalData = nil
	w, _ = NewWrite(Master(), "_Installation", query, data, originalData, nil)
	err = w.handleInstallation()
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	expectQuery = types.M{"objectId": "1001"}
	if reflect.DeepEqual(expectQuery, w.query) == false {
		t.Error("expect:", expectQuery, "result:", w.query)
	}
	results, err = orm.TomatoDBController.Find("_Installation", types.M{}, types.M{})
	expect = types.S{
		types.M{
			"objectId":       "1001",
			"installationId": "111",
			"deviceToken":    "bbb",
			"deviceType":     "ios",
		},
		types.M{
			"objectId":    "1002",
			"deviceToken": "aaa",
			"deviceType":  "ios",
		},
	}
	if err != nil || reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	schema = types.M{
		"fields": types.M{
			"installationId": types.M{"type": "String"},
			"deviceToken":    types.M{"type": "String"},
			"deviceType":     types.M{"type": "String"},
			"appIdentifier":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_Installation", schema)
	object = types.M{
		"objectId":       "1001",
		"installationId": "111",
		"deviceToken":    "aaa",
		"deviceType":     "ios",
	}
	orm.Adapter.CreateObject("_Installation", schema, object)
	object = types.M{
		"objectId":    "1002",
		"deviceToken": "aaa",
		"deviceType":  "ios",
	}
	orm.Adapter.CreateObject("_Installation", schema, object)
	query = nil
	data = types.M{"deviceToken": "aaa", "installationId": "111", "deviceType": "android"}
	originalData = nil
	w, _ = NewWrite(Master(), "_Installation", query, data, originalData, nil)
	err = w.handleInstallation()
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	expectQuery = types.M{"objectId": "1001"}
	if reflect.DeepEqual(expectQuery, w.query) == false {
		t.Error("expect:", expectQuery, "result:", w.query)
	}
	results, err = orm.TomatoDBController.Find("_Installation", types.M{}, types.M{})
	expect = types.S{
		types.M{
			"objectId":       "1001",
			"installationId": "111",
			"deviceToken":    "aaa",
			"deviceType":     "ios",
		},
		types.M{
			"objectId":    "1002",
			"deviceToken": "aaa",
			"deviceType":  "ios",
		},
	}
	if err != nil || reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	schema = types.M{
		"fields": types.M{
			"installationId": types.M{"type": "String"},
			"deviceToken":    types.M{"type": "String"},
			"deviceType":     types.M{"type": "String"},
			"appIdentifier":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_Installation", schema)
	object = types.M{
		"objectId":       "1001",
		"installationId": "111",
		"deviceToken":    "bbb",
		"deviceType":     "ios",
	}
	orm.Adapter.CreateObject("_Installation", schema, object)
	object = types.M{
		"objectId":    "1002",
		"deviceToken": "aaa",
		"deviceType":  "ios",
	}
	orm.Adapter.CreateObject("_Installation", schema, object)
	query = nil
	data = types.M{"deviceToken": "ccc", "installationId": "111", "deviceType": "android"}
	originalData = nil
	w, _ = NewWrite(Master(), "_Installation", query, data, originalData, nil)
	err = w.handleInstallation()
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	expectQuery = types.M{"objectId": "1001"}
	if reflect.DeepEqual(expectQuery, w.query) == false {
		t.Error("expect:", expectQuery, "result:", w.query)
	}
	results, err = orm.TomatoDBController.Find("_Installation", types.M{}, types.M{})
	expect = types.S{
		types.M{
			"objectId":       "1001",
			"installationId": "111",
			"deviceToken":    "bbb",
			"deviceType":     "ios",
		},
		types.M{
			"objectId":    "1002",
			"deviceToken": "aaa",
			"deviceType":  "ios",
		},
	}
	if err != nil || reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	schema = types.M{
		"fields": types.M{
			"installationId": types.M{"type": "String"},
			"deviceToken":    types.M{"type": "String"},
			"deviceType":     types.M{"type": "String"},
			"appIdentifier":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_Installation", schema)
	object = types.M{
		"objectId":       "1001",
		"installationId": "111",
		"deviceToken":    "bbb",
		"deviceType":     "ios",
	}
	orm.Adapter.CreateObject("_Installation", schema, object)
	object = types.M{
		"objectId":       "1002",
		"installationId": "222",
		"deviceToken":    "aaa",
		"deviceType":     "ios",
	}
	orm.Adapter.CreateObject("_Installation", schema, object)
	query = nil
	data = types.M{"deviceToken": "aaa", "installationId": "111", "deviceType": "android"}
	originalData = nil
	w, _ = NewWrite(Master(), "_Installation", query, data, originalData, nil)
	err = w.handleInstallation()
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	expectQuery = types.M{"objectId": "1001"}
	if reflect.DeepEqual(expectQuery, w.query) == false {
		t.Error("expect:", expectQuery, "result:", w.query)
	}
	results, err = orm.TomatoDBController.Find("_Installation", types.M{}, types.M{})
	expect = types.S{
		types.M{
			"objectId":       "1001",
			"installationId": "111",
			"deviceToken":    "bbb",
			"deviceType":     "ios",
		},
	}
	if err != nil || reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_handleSession(t *testing.T) {
	var w *Write
	var auth *Auth
	var query types.M
	var data types.M
	var originalData types.M
	var err, expectErr error
	var results types.S
	/***************************************************************/
	auth = Nobody()
	query = nil
	data = types.M{}
	originalData = nil
	w, _ = NewWrite(auth, "_Session", query, data, originalData, nil)
	err = w.handleSession()
	expectErr = errs.E(errs.InvalidSessionToken, "Session token required.")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/***************************************************************/
	auth = &Auth{
		IsMaster: false,
		User:     types.M{"objectId": "1001"},
	}
	query = nil
	data = types.M{"ACL": "hello"}
	originalData = nil
	w, _ = NewWrite(auth, "_Session", query, data, originalData, nil)
	err = w.handleSession()
	expectErr = errs.E(errs.InvalidKeyName, "Cannot set ACL on a Session.")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/***************************************************************/
	auth = Master()
	query = nil
	data = types.M{"ACL": "hello"}
	originalData = nil
	w, _ = NewWrite(auth, "_Session", query, data, originalData, nil)
	err = w.handleSession()
	expectErr = errs.E(errs.InvalidKeyName, "Cannot set ACL on a Session.")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/***************************************************************/
	initEnv()
	auth = &Auth{
		IsMaster: false,
		User:     types.M{"objectId": "1001"},
	}
	query = nil
	data = types.M{}
	originalData = nil
	config.TConfig.SessionLength = 31536000
	livequery.TLiveQuery = livequery.NewLiveQuery([]string{}, "", "", "")
	w, _ = NewWrite(auth, "_Session", query, data, originalData, nil)
	err = w.handleSession()
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	results, _ = orm.TomatoDBController.Find("_Session", types.M{}, types.M{})
	if len(results) != 1 {
		t.Error("expect:", "len 1", "result:", results)
	}
	orm.TomatoDBController.DeleteEverything()
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
	config.TConfig.EnableAnonymousUsers = true
	className = "_User"
	query = nil
	data = types.M{
		"authData": types.M{
			"anonymous": types.M{
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
	if reflect.DeepEqual([]string{"username"}, w.storage["fieldsChangedByTrigger"]) == false {
		t.Error("expect:", []string{"username"}, "result:", w.storage["fieldsChangedByTrigger"])
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
	if reflect.DeepEqual(nil, w.storage["fieldsChangedByTrigger"]) == false {
		t.Error("expect:", nil, "result:", w.storage["fieldsChangedByTrigger"])
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
	if reflect.DeepEqual([]string{"username"}, w.storage["fieldsChangedByTrigger"]) == false {
		t.Error("expect:", []string{"username"}, "result:", w.storage["fieldsChangedByTrigger"])
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
	if reflect.DeepEqual(nil, w.storage["fieldsChangedByTrigger"]) == false {
		t.Error("expect:", nil, "result:", w.storage["fieldsChangedByTrigger"])
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
	policyError := "Password does not meet the Password Policy requirements."
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
		"emailVerified": true,
	}
	originalData = nil
	w, _ = NewWrite(Nobody(), "_User", query, data, originalData, nil)
	err = w.transformUser()
	expectErr = errs.E(errs.OperationForbidden, "Clients aren't allowed to manually update email verification.")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
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
	config.TConfig.PasswordPolicy = true
	config.TConfig.DoNotAllowUsername = true
	query = nil
	data = types.M{
		"username": "joe",
		"password": "joe123456",
	}
	originalData = nil
	w, _ = NewWrite(Master(), "_User", query, data, originalData, nil)
	w.data["objectId"] = "1001"
	err = w.transformUser()
	expectErr = errs.E(errs.ValidationError, policyError)
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	config.TConfig.DoNotAllowUsername = false
	config.TConfig.PasswordPolicy = false
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	config.TConfig.PasswordPolicy = true
	config.TConfig.DoNotAllowUsername = true
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
	query = types.M{"objectId": "1002"}
	data = types.M{
		"password": "joe123456",
	}
	originalData = nil
	w, _ = NewWrite(Master(), "_User", query, data, originalData, nil)
	w.data["objectId"] = "1002"
	err = w.transformUser()
	expectErr = errs.E(errs.ValidationError, policyError)
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	config.TConfig.DoNotAllowUsername = false
	config.TConfig.PasswordPolicy = false
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	config.TConfig.PasswordPolicy = true
	config.TConfig.MaxPasswordHistory = 3
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId":         "1002",
		"username":         "joe",
		"_hashed_password": "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3",
		"_password_history": []interface{}{
			"b3a8e0e1f9ab1bfe3a36f231f676f78bb30a519d2b21e6c530c0eee8ebb4a5d0",
			"35a9e381b1a27567549b5f8a6f783c167ebf809f1c4d6a9e367240484d8ce281",
		},
	}
	err = orm.Adapter.CreateObject("_User", schema, object)
	query = types.M{"objectId": "1002"}
	data = types.M{
		"password": "123456",
	}
	originalData = nil
	w, _ = NewWrite(Master(), "_User", query, data, originalData, nil)
	w.data["objectId"] = "1002"
	err = w.transformUser()
	expect = types.M{
		"objectId":         "1002",
		"_hashed_password": utils.Hash("123456"),
	}
	if err != nil || reflect.DeepEqual(expect, w.data) == false {
		t.Error("expect:", expect, "result:", w.data, "err:", err)
	}
	config.TConfig.MaxPasswordHistory = 0
	config.TConfig.PasswordPolicy = false
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	initEnv()
	config.TConfig.PasswordPolicy = true
	config.TConfig.MaxPasswordHistory = 3
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId":         "1002",
		"username":         "joe",
		"_hashed_password": "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3",
		"_password_history": []interface{}{
			"b3a8e0e1f9ab1bfe3a36f231f676f78bb30a519d2b21e6c530c0eee8ebb4a5d0",
			"35a9e381b1a27567549b5f8a6f783c167ebf809f1c4d6a9e367240484d8ce281",
		},
	}
	err = orm.Adapter.CreateObject("_User", schema, object)
	query = types.M{"objectId": "1002"}
	data = types.M{
		"password": "123",
	}
	originalData = nil
	w, _ = NewWrite(Master(), "_User", query, data, originalData, nil)
	w.data["objectId"] = "1002"
	err = w.transformUser()
	expectErr = errs.E(errs.ValidationError, "New password should not be the same as last 3 passwords.")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	config.TConfig.MaxPasswordHistory = 0
	config.TConfig.PasswordPolicy = false
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
	config.TConfig.ServerURL = "http://127.0.0.1"
	config.TConfig.AppID = "1001"
	w, _ := NewWrite(Master(), "user", nil, types.M{}, nil, nil)
	w.response = types.M{
		"response": types.M{
			"file": types.M{
				"__type": "File",
				"name":   "hello.jpg",
			},
		},
	}
	w.expandFilesForExistingObjects()
	expect := types.M{
		"response": types.M{
			"file": types.M{
				"__type": "File",
				"name":   "hello.jpg",
				"url":    "http://127.0.0.1/files/1001/hello.jpg",
			},
		},
	}
	if reflect.DeepEqual(expect, w.response) == false {
		t.Error("expect:", expect, "result:", w.response)
	}
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
	expectErr = errs.E(errs.InvalidACL, "Invalid ACL.")
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
	expectErr = errs.E(errs.InvalidACL, "Invalid ACL.")
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
	w.storage["fieldsChangedByTrigger"] = []string{"key"}
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
	w.storage["fieldsChangedByTrigger"] = []string{"username"}
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
	// 测试用例与 createSessionToken 相同
}

func Test_handleFollowup(t *testing.T) {
	var schema, object types.M
	var className string
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	var err error
	var results types.S
	/***************************************************************/
	initEnv()
	className = "_Session"
	schema = types.M{
		"fields": types.M{
			"restricted":     types.M{"type": "Boolean"},
			"user":           types.M{"type": "Pointer", "targetClass": "_User"},
			"installationId": types.M{"type": "String"},
			"sessionToken":   types.M{"type": "String"},
			"expiresAt":      types.M{"type": "Date"},
			"createdWith":    types.M{"type": "Object"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"user": types.M{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  "1001",
		},
		"sessionToken": "r:aaa",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2002",
		"user": types.M{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  "1002",
		},
		"sessionToken": "r:bbb",
	}
	orm.Adapter.CreateObject(className, schema, object)
	config.TConfig.RevokeSessionOnPasswordReset = true
	className = "_User"
	query = types.M{"objectId": "1001"}
	data = types.M{}
	originalData = types.M{}
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	w.storage = types.M{
		"clearSessions": true,
	}
	err = w.handleFollowup()
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	results, _ = orm.TomatoDBController.Find("_Session", types.M{}, types.M{})
	if len(results) != 1 {
		t.Error("expect:", "len 1", "result:", results)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_runAfterTrigger(t *testing.T) {
	var className string
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	/***************************************************************/
	livequery.TLiveQuery = livequery.NewLiveQuery([]string{}, "", "", "")
	className = "user"
	query = nil
	data = types.M{"key": "hello"}
	originalData = nil
	w, _ = NewWrite(Master(), className, query, data, originalData, nil)
	w.response = types.M{"response": "hello"}
	w.runAfterTrigger()
	/***************************************************************/
	livequery.TLiveQuery = livequery.NewLiveQuery([]string{}, "", "", "")
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
	var w *Write
	var expect types.M
	/***************************************************************/
	w, _ = NewWrite(Master(), "_User", nil, types.M{}, nil, nil)
	w.response = types.M{
		"response": types.M{
			"username": "joe",
			"authData": types.M{
				"facebook": types.M{"id": "1001"},
				"weibo":    nil,
			},
		},
	}
	w.cleanUserAuthData()
	expect = types.M{
		"response": types.M{
			"username": "joe",
			"authData": types.M{
				"facebook": types.M{"id": "1001"},
			},
		},
	}
	if reflect.DeepEqual(expect, w.response) == false {
		t.Error("expect:", expect, "result:", w.response)
	}
	/***************************************************************/
	w, _ = NewWrite(Master(), "_User", nil, types.M{}, nil, nil)
	w.response = types.M{
		"response": types.M{
			"username": "joe",
			"authData": types.M{
				"facebook": nil,
				"weibo":    nil,
			},
		},
	}
	w.cleanUserAuthData()
	expect = types.M{
		"response": types.M{
			"username": "joe",
		},
	}
	if reflect.DeepEqual(expect, w.response) == false {
		t.Error("expect:", expect, "result:", w.response)
	}
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
	config.TConfig.EnableAnonymousUsers = true
	className = "_User"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"anonymous": types.M{
				"id": "1001",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	object = types.M{
		"objectId": "102",
		"authData": types.M{
			"anonymous": types.M{
				"id": "1001",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	className = "_User"
	query = nil
	data = types.M{
		"authData": types.M{
			"anonymous": types.M{
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
	config.TConfig.EnableAnonymousUsers = true
	className = "_User"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"anonymous": types.M{
				"id": "1001",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	className = "_User"
	query = nil
	data = types.M{
		"authData": types.M{
			"anonymous": types.M{
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
	if reflect.DeepEqual("anonymous", w.storage["authProvider"]) == false {
		t.Error("expect:", "anonymous", "result:", w.storage["authProvider"])
	}
	orm.TomatoDBController.DeleteEverything()
	/***************************************************************/
	config.TConfig = &config.Config{
		ServerURL: "http://www.g.cn",
	}
	config.TConfig.EnableAnonymousUsers = true
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"anonymous": types.M{
				"id": "1001",
			},
		},
	}
	orm.TomatoDBController.Create(className, object, nil)
	className = "_User"
	query = nil
	data = types.M{
		"authData": types.M{
			"anonymous": types.M{
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
			"anonymous": types.M{
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
	config.TConfig.EnableAnonymousUsers = true
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"anonymous": types.M{
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
			"anonymous": types.M{
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
			"anonymous": types.M{
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
			"anonymous": types.M{
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
	config.TConfig.EnableAnonymousUsers = true
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"anonymous": types.M{
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
			"anonymous": types.M{
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
	config.TConfig.EnableAnonymousUsers = true
	initEnv()
	className = "_User"
	schema = types.M{
		"fields": types.M{},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "101",
		"authData": types.M{
			"anonymous": types.M{
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
			"anonymous": types.M{
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
	config.TConfig.EnableAnonymousUsers = true
	query = nil
	data = types.M{}
	originalData = nil
	w, _ = NewWrite(Master(), "user", query, data, originalData, nil)
	authData = types.M{
		"anonymous": types.M{
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
	var w *Write
	var query types.M
	var data types.M
	var originalData types.M
	var results types.S
	/***************************************************************/
	initEnv()
	livequery.TLiveQuery = livequery.NewLiveQuery([]string{}, "", "", "")
	config.TConfig.SessionLength = 31536000
	query = nil
	data = types.M{
		"username": "joe",
	}
	originalData = nil
	w, _ = NewWrite(Master(), "_User", query, data, originalData, nil)
	w.data["objectId"] = "1001"
	w.response = types.M{
		"response": types.M{
			"objectId": "1001",
			"username": "joe",
		},
	}
	w.createSessionToken()
	results, _ = orm.TomatoDBController.Find("_Session", types.M{}, types.M{})
	if len(results) != 1 {
		t.Error("expect:", "len 1", "result:", results)
	}
	if r := utils.M(w.response["response"]); r != nil {
		if r["sessionToken"] == nil {
			t.Error("expect:", "need sessionToken", "result:", r["sessionToken"])
		}
	}
	orm.TomatoDBController.DeleteEverything()
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
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
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
	w.storage["fieldsChangedByTrigger"] = []string{"key3"}
	result = w.updateResponseWithData(response, updateData)
	expect = types.M{
		"key":  "hello",
		"key1": 10,
		"key3": types.M{
			"__op": "Delete",
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_getLastItems(t *testing.T) {
	type fields struct {
		items []interface{}
		n     int
	}
	tests := []struct {
		name   string
		fields fields
		want   []interface{}
	}{
		{
			name: "1",
			fields: fields{
				items: nil,
				n:     0,
			},
			want: nil,
		},
		{
			name: "2",
			fields: fields{
				items: []interface{}{"abc"},
				n:     0,
			},
			want: []interface{}{},
		},
		{
			name: "2",
			fields: fields{
				items: []interface{}{"abc"},
				n:     -1,
			},
			want: []interface{}{},
		},
		{
			name: "3",
			fields: fields{
				items: []interface{}{"abc", "def"},
				n:     3,
			},
			want: []interface{}{"abc", "def"},
		},
		{
			name: "4",
			fields: fields{
				items: []interface{}{"abc", "def", "hij"},
				n:     3,
			},
			want: []interface{}{"abc", "def", "hij"},
		},
		{
			name: "5",
			fields: fields{
				items: []interface{}{"abc", "def", "hij", "klm"},
				n:     3,
			},
			want: []interface{}{"def", "hij", "klm"},
		},
	}
	for _, tt := range tests {
		w := getLastItems(tt.fields.items, tt.fields.n)
		if reflect.DeepEqual(w, tt.want) == false {
			t.Errorf("%q. getLastItems() result = %v, want %v", tt.name, w, tt.want)
		}
	}
}
