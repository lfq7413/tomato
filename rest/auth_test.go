package rest

import (
	"testing"

	"github.com/lfq7413/tomato/types"
)

func Test_GetAuthForSessionToken(t *testing.T) {
	// TODO
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
