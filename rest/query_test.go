package rest

import (
	"reflect"
	"testing"

	"github.com/lfq7413/tomato/types"
)

func Test_Execute(t *testing.T) {
	// BuildRestWhere
	// runFind
	// runCount
	// handleInclude
	// TODO
}

func Test_BuildRestWhere(t *testing.T) {
	// getUserAndRoleACL
	// redirectClassNameForKey
	// validateClientClassCreation
	// replaceSelect
	// replaceDontSelect
	// replaceInQuery
	// replaceNotInQuery
	// TODO
}

func Test_getUserAndRoleACL(t *testing.T) {
	// TODO
}

func Test_redirectClassNameForKey(t *testing.T) {
	// TODO
}

func Test_validateClientClassCreation(t *testing.T) {
	// TODO
}

func Test_replaceSelect(t *testing.T) {
	// findObjectWithKey
	// NewQuery
	// Execute
	// transformSelect
	// TODO
}

func Test_replaceDontSelect(t *testing.T) {
	// findObjectWithKey
	// NewQuery
	// Execute
	// transformDontSelect
	// TODO
}

func Test_replaceInQuery(t *testing.T) {
	// findObjectWithKey
	// NewQuery
	// Execute
	// transformInQuery
	// TODO
}

func Test_replaceNotInQuery(t *testing.T) {
	// findObjectWithKey
	// NewQuery
	// Execute
	// transformNotInQuery
	// TODO
}

func Test_runFind(t *testing.T) {
	// TODO
}

func Test_runCount(t *testing.T) {
	// TODO
}

func Test_handleInclude(t *testing.T) {
	// includePath
	// TODO
}

/////////////////////////////////////////////////////////////////

func Test_NewQuery(t *testing.T) {
	// TODO
}

func Test_includePath(t *testing.T) {
	// findPointers
	// NewQuery
	// Execute
	// replacePointers
	// TODO
}

func Test_findPointers(t *testing.T) {
	// TODO
}

func Test_replacePointers(t *testing.T) {
	// TODO
}

func Test_findObjectWithKey(t *testing.T) {
	// TODO
}

func Test_transformSelect(t *testing.T) {
	// TODO
}

func Test_transformDontSelect(t *testing.T) {
	// TODO
}

func Test_transformInQuery(t *testing.T) {
	// TODO
}

func Test_transformNotInQuery(t *testing.T) {
	var notInQueryObject types.M
	var className string
	var results []types.M
	var expect types.M
	/**********************************************************/
	notInQueryObject = nil
	className = "user"
	results = nil
	transformNotInQuery(notInQueryObject, className, results)
	expect = nil
	if reflect.DeepEqual(expect, notInQueryObject) == false {
		t.Error("expect:", expect, "result:", notInQueryObject)
	}
	/**********************************************************/
	notInQueryObject = types.M{}
	className = "user"
	results = nil
	transformNotInQuery(notInQueryObject, className, results)
	expect = types.M{}
	if reflect.DeepEqual(expect, notInQueryObject) == false {
		t.Error("expect:", expect, "result:", notInQueryObject)
	}
	// TODO
}
