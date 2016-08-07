package rest

import (
	"reflect"
	"testing"
	"time"

	"github.com/lfq7413/tomato/cache"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

func Test_GetAuthForSessionToken(t *testing.T) {
	var schema types.M
	var object types.M
	var className string
	var sessionToken string
	var installationID string
	var result *Auth
	var err error
	var expect *Auth
	var expectErr error
	/********************************************************/
	cache.InitCache()
	initEnv()
	sessionToken = "abc"
	installationID = "111"
	_, err = GetAuthForSessionToken(sessionToken, installationID)
	expectErr = errs.E(errs.InvalidSessionToken, "invalid session token")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/********************************************************/
	cache.InitCache()
	initEnv()
	className = "_User"
	schema = types.M{
		"username": types.M{"type": "String"},
		"password": types.M{"type": "String"},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"username": "joe",
		"password": "123",
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "_Session"
	schema = types.M{
		"user":         types.M{"type": "Pointer", "targetClass": "_User"},
		"sessionToken": types.M{"type": "String"},
		"expiresAt":    types.M{"type": "Date"},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"user": types.M{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  "1001",
		},
		"sessionToken": "abc1001",
		"expiresAt":    utils.TimetoString(time.Now().UTC()),
	}
	orm.Adapter.CreateObject(className, schema, object)
	sessionToken = "abc"
	installationID = "111"
	_, err = GetAuthForSessionToken(sessionToken, installationID)
	expectErr = errs.E(errs.InvalidSessionToken, "invalid session token")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/********************************************************/
	cache.InitCache()
	initEnv()
	className = "_User"
	schema = types.M{
		"username": types.M{"type": "String"},
		"password": types.M{"type": "String"},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"username": "joe",
		"password": "123",
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "_Session"
	schema = types.M{
		"user":         types.M{"type": "Pointer", "targetClass": "_User"},
		"sessionToken": types.M{"type": "String"},
		"expiresAt":    types.M{"type": "Date"},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"user": types.M{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  "1001",
		},
		"sessionToken": "abc1001",
	}
	orm.Adapter.CreateObject(className, schema, object)
	sessionToken = "abc1001"
	installationID = "111"
	_, err = GetAuthForSessionToken(sessionToken, installationID)
	expectErr = errs.E(errs.InvalidSessionToken, "Session token is expired.")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/********************************************************/
	cache.InitCache()
	initEnv()
	className = "_User"
	schema = types.M{
		"username": types.M{"type": "String"},
		"password": types.M{"type": "String"},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"username": "joe",
		"password": "123",
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "_Session"
	schema = types.M{
		"user":         types.M{"type": "Pointer", "targetClass": "_User"},
		"sessionToken": types.M{"type": "String"},
		"expiresAt":    types.M{"type": "Date"},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"user": types.M{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  "1001",
		},
		"sessionToken": "abc1001",
		"expiresAt":    utils.TimetoString(time.Now().UTC()),
	}
	orm.Adapter.CreateObject(className, schema, object)
	sessionToken = "abc1001"
	installationID = "111"
	_, err = GetAuthForSessionToken(sessionToken, installationID)
	expectErr = errs.E(errs.InvalidSessionToken, "Session token is expired.")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/********************************************************/
	cache.InitCache()
	initEnv()
	className = "_User"
	schema = types.M{
		"username": types.M{"type": "String"},
		"password": types.M{"type": "String"},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"username": "joe",
		"password": "123",
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "_Session"
	schema = types.M{
		"user":         types.M{"type": "Pointer", "targetClass": "_User"},
		"sessionToken": types.M{"type": "String"},
		"expiresAt":    types.M{"type": "Date"},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"user": types.M{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  "1001",
		},
		"sessionToken": "abc1001",
		"expiresAt":    utils.TimetoString(time.Now().UTC().Add(5 * time.Second)),
	}
	orm.Adapter.CreateObject(className, schema, object)
	sessionToken = "abc1001"
	installationID = "111"
	result, err = GetAuthForSessionToken(sessionToken, installationID)
	expect = &Auth{
		IsMaster:       false,
		InstallationID: "111",
		User: types.M{
			"__type":       "Object",
			"className":    "_User",
			"objectId":     "1001",
			"username":     "joe",
			"sessionToken": "abc1001",
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_CouldUpdateUserID(t *testing.T) {
	var auth *Auth
	var result bool
	var expect bool
	/********************************************************/
	auth = &Auth{
		IsMaster: true,
	}
	result = auth.CouldUpdateUserID("1001")
	expect = true
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/********************************************************/
	auth = &Auth{
		IsMaster: false,
	}
	result = auth.CouldUpdateUserID("1001")
	expect = false
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/********************************************************/
	auth = &Auth{
		IsMaster: false,
		User:     types.M{},
	}
	result = auth.CouldUpdateUserID("1001")
	expect = false
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/********************************************************/
	auth = &Auth{
		IsMaster: false,
		User:     types.M{"objectId": "1002"},
	}
	result = auth.CouldUpdateUserID("1001")
	expect = false
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/********************************************************/
	auth = &Auth{
		IsMaster: false,
		User:     types.M{"objectId": "1001"},
	}
	result = auth.CouldUpdateUserID("1001")
	expect = true
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_GetUserRoles(t *testing.T) {
	// loadRoles
	// TODO
}

func Test_loadRoles(t *testing.T) {
	// getAllRoleNamesForID
	// TODO
}

func Test_getAllRoleNamesForID(t *testing.T) {
	// TODO
}
