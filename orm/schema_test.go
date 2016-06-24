package orm

import (
	"testing"

	"github.com/lfq7413/tomato/types"
)

func Test_AddClassIfNotExists(t *testing.T) {
	// validateNewClass
	// convertSchemaToAdapterSchema
	// convertAdapterSchemaToParseSchema
	// TODO
}

func Test_UpdateClass(t *testing.T) {
	// GetOneSchema
	// buildMergedSchemaObject
	// validateSchemaData
	// deleteField
	// reloadData
	// enforceFieldExists
	// setPermissions
	// TODO
}

func Test_deleteField(t *testing.T) {
	// ClassNameIsValid
	// fieldNameIsValid
	// fieldNameIsValidForClass
	// GetOneSchema
	// TODO
}

func Test_validateObject(t *testing.T) {
	// EnforceClassExists
	// enforceFieldExists
	// thenValidateRequiredColumns
	// TODO
}

func Test_testBaseCLP(t *testing.T) {
	// TODO
}

func Test_validatePermission(t *testing.T) {
	// testBaseCLP
	// TODO
}

func Test_EnforceClassExists(t *testing.T) {
	// AddClassIfNotExists
	// TODO
}

func Test_validateNewClass(t *testing.T) {
	// InvalidClassNameMessage
	// validateSchemaData
	// TODO
}

func Test_validateSchemaData(t *testing.T) {
	// fieldNameIsValid
	// fieldNameIsValidForClass
	// fieldTypeIsInvalid
	// validateCLP
	// TODO
}

func Test_validateRequiredColumns(t *testing.T) {
	// TODO
}

func Test_enforceFieldExists(t *testing.T) {
	// fieldNameIsValid
	// getExpectedType
	// TODO
}

func Test_setPermissions(t *testing.T) {
	// validateCLP
	// reloadData
	// TODO
}

func Test_HasClass(t *testing.T) {
	// reloadData
	// TODO
}

func Test_getExpectedType(t *testing.T) {
	// TODO
}

func Test_reloadData(t *testing.T) {
	// GetAllClasses
	// injectDefaultSchema
	// TODO
}

func Test_GetAllClasses(t *testing.T) {
	// injectDefaultSchema
	// TODO
}

func Test_GetOneSchema(t *testing.T) {
	// injectDefaultSchema
	// TODO
}

////////////////////////////////////////////////////////////

func Test_thenValidateRequiredColumns(t *testing.T) {
	// validateRequiredColumns
	// TODO
}

func Test_getType(t *testing.T) {
	// getObjectType
	// TODO
}

func Test_getObjectType(t *testing.T) {
	// TODO
}

func Test_ClassNameIsValid(t *testing.T) {
	// joinClassIsValid
	// fieldNameIsValid
	// TODO
}

func Test_InvalidClassNameMessage(t *testing.T) {
	// TODO
}

func Test_joinClassIsValid(t *testing.T) {
	// TODO
}

func Test_fieldNameIsValid(t *testing.T) {
	// TODO
}

func Test_fieldNameIsValidForClass(t *testing.T) {
	// TODO
}

func Test_fieldTypeIsInvalid(t *testing.T) {
	// ClassNameIsValid
	// InvalidClassNameMessage
	// TODO
}

func Test_validateCLP(t *testing.T) {
	// verifyPermissionKey
	// TODO
}

func Test_verifyPermissionKey(t *testing.T) {
	// TODO
}

func Test_buildMergedSchemaObject(t *testing.T) {
	// TODO
}

func Test_injectDefaultSchema(t *testing.T) {
	// TODO
}

func Test_convertSchemaToAdapterSchema(t *testing.T) {
	// injectDefaultSchema
	// TODO
}

func Test_convertAdapterSchemaToParseSchema(t *testing.T) {
	// TODO
}

func Test_dbTypeMatchesObjectType(t *testing.T) {
	var dbType types.M
	var objectType types.M
	var ok bool
	var expect bool
	/************************************************************/
	dbType = nil
	objectType = nil
	ok = dbTypeMatchesObjectType(dbType, objectType)
	expect = true
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	dbType = types.M{}
	objectType = nil
	ok = dbTypeMatchesObjectType(dbType, objectType)
	expect = false
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	dbType = nil
	objectType = types.M{}
	ok = dbTypeMatchesObjectType(dbType, objectType)
	expect = false
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	dbType = types.M{"type": "String"}
	objectType = types.M{}
	ok = dbTypeMatchesObjectType(dbType, objectType)
	expect = false
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	dbType = types.M{"type": "String"}
	objectType = types.M{"type": "Date"}
	ok = dbTypeMatchesObjectType(dbType, objectType)
	expect = false
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	dbType = types.M{"type": "Pointer", "targetClass": "abc"}
	objectType = types.M{"type": "Pointer", "targetClass": "def"}
	ok = dbTypeMatchesObjectType(dbType, objectType)
	expect = false
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	dbType = types.M{"type": "Pointer", "targetClass": "abc"}
	objectType = types.M{"type": "Pointer", "targetClass": "abc"}
	ok = dbTypeMatchesObjectType(dbType, objectType)
	expect = true
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	dbType = types.M{"type": "String"}
	objectType = types.M{"type": "String"}
	ok = dbTypeMatchesObjectType(dbType, objectType)
	expect = true
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
}

func Test_Load(t *testing.T) {
	// reloadData
	// TODO
}
