package rest

import (
	"reflect"
	"testing"
	"time"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
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
	// TODO
}

func Test_validateClientClassCreation_Write(t *testing.T) {
	// TODO
}

func Test_validateSchema(t *testing.T) {
	// TODO
}

func Test_handleInstallation(t *testing.T) {
	// TODO
}

func Test_handleSession(t *testing.T) {
	// Execute
	// TODO
}

func Test_validateAuthData(t *testing.T) {
	// handleAuthData
	// TODO
}

func Test_runBeforeTrigger(t *testing.T) {
	// TODO
}

func Test_setRequiredFieldsIfNeeded(t *testing.T) {
	// TODO
}

func Test_transformUser(t *testing.T) {
	// TODO
}

func Test_expandFilesForExistingObjects(t *testing.T) {
	// TODO
}

func Test_runDatabaseOperation(t *testing.T) {
	// TODO
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
	// TODO
}

func Test_cleanUserAuthData(t *testing.T) {
	// TODO
}

/////////////////////////////////////////////////////////////

func Test_handleAuthData(t *testing.T) {
	// handleAuthDataValidation
	// TODO
}

func Test_handleAuthDataValidation(t *testing.T) {
	// TODO
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
