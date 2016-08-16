package rest

import (
	"reflect"
	"testing"
	"time"

	"github.com/lfq7413/tomato/errs"
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
	// updateResponseWithData
	// location
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
	// findUsersWithAuthData
	// location
	// TODO
}

func Test_handleAuthDataValidation(t *testing.T) {
	// TODO
}

func Test_findUsersWithAuthData(t *testing.T) {
	// TODO
}

func Test_createSessionToken(t *testing.T) {
	// TODO
}

func Test_location(t *testing.T) {
	// TODO
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
		"_auth_data":       "fackbook",
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
	// TODO
}
