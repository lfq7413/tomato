package rest

import (
	"testing"

	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/types"
)

func Test_HandleLoginAttempt(t *testing.T) {
	// TODO
	// notLocked
	// setFailedLoginCount
	// handleFailedLoginAttempt
}

func Test_notLocked(t *testing.T) {
	// TODO
}

func Test_setFailedLoginCount(t *testing.T) {
	// TODO
}

func Test_handleFailedLoginAttempt(t *testing.T) {
	// TODO
	// initFailedLoginCount
	// incrementFailedLoginCount
	// setLockoutExpiration
}

func Test_initFailedLoginCount(t *testing.T) {
	// TODO
	// isFailedLoginCountSet
	// setFailedLoginCount
}

func Test_incrementFailedLoginCount(t *testing.T) {
	// TODO
}

func Test_setLockoutExpiration(t *testing.T) {
	// TODO
}

func Test_isFailedLoginCountSet(t *testing.T) {
	var username string
	var object, schema types.M
	var accountLockout *AccountLockout
	var isSet bool
	var err error
	/*****************************************************************/
	initEnv()
	username = "joe"
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId": "01",
		"username": username,
	}
	orm.Adapter.CreateObject("_User", schema, object)
	accountLockout = NewAccountLockout(username)
	isSet, err = accountLockout.isFailedLoginCountSet()
	if err != nil || isSet != false {
		t.Error("expect:", false, "result:", isSet, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/*****************************************************************/
	initEnv()
	username = "joe"
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId":            "01",
		"username":            username,
		"_failed_login_count": 3,
	}
	orm.Adapter.CreateObject("_User", schema, object)
	accountLockout = NewAccountLockout(username)
	isSet, err = accountLockout.isFailedLoginCountSet()
	if err != nil || isSet != true {
		t.Error("expect:", true, "result:", isSet, err)
	}
	orm.TomatoDBController.DeleteEverything()
}
