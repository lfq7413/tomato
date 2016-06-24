package orm

import (
	"reflect"
	"testing"

	"github.com/lfq7413/tomato/types"
)

func Test_AddClassIfNotExists(t *testing.T) {
	// validateNewClass
	// convertSchemaToAdapterSchema
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
	// TODO
}

func Test_GetAllClasses(t *testing.T) {
	// TODO
}

func Test_GetOneSchema(t *testing.T) {
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
	var schema types.M
	var result types.M
	var expect types.M
	/************************************************************/
	schema = nil
	result = injectDefaultSchema(schema)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/************************************************************/
	schema = types.M{
		"className": "user",
	}
	result = injectDefaultSchema(schema)
	expect = types.M{
		"className": "user",
		"fields": types.M{
			"objectId":  types.M{"type": "String"},
			"createdAt": types.M{"type": "Date"},
			"updatedAt": types.M{"type": "Date"},
			"ACL":       types.M{"type": "ACL"},
		},
		"classLevelPermissions": nil,
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/************************************************************/
	schema = types.M{
		"className": "user",
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	result = injectDefaultSchema(schema)
	expect = types.M{
		"className": "user",
		"fields": types.M{
			"objectId":  types.M{"type": "String"},
			"createdAt": types.M{"type": "Date"},
			"updatedAt": types.M{"type": "Date"},
			"ACL":       types.M{"type": "ACL"},
			"key":       types.M{"type": "String"},
		},
		"classLevelPermissions": nil,
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/************************************************************/
	schema = types.M{
		"className": "_User",
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	result = injectDefaultSchema(schema)
	expect = types.M{
		"className": "_User",
		"fields": types.M{
			"objectId":      types.M{"type": "String"},
			"createdAt":     types.M{"type": "Date"},
			"updatedAt":     types.M{"type": "Date"},
			"ACL":           types.M{"type": "ACL"},
			"key":           types.M{"type": "String"},
			"username":      types.M{"type": "String"},
			"password":      types.M{"type": "String"},
			"email":         types.M{"type": "String"},
			"emailVerified": types.M{"type": "Boolean"},
		},
		"classLevelPermissions": nil,
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/************************************************************/
	schema = types.M{
		"className": "_User",
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"find": types.M{"*": true},
		},
	}
	result = injectDefaultSchema(schema)
	expect = types.M{
		"className": "_User",
		"fields": types.M{
			"objectId":      types.M{"type": "String"},
			"createdAt":     types.M{"type": "Date"},
			"updatedAt":     types.M{"type": "Date"},
			"ACL":           types.M{"type": "ACL"},
			"key":           types.M{"type": "String"},
			"username":      types.M{"type": "String"},
			"password":      types.M{"type": "String"},
			"email":         types.M{"type": "String"},
			"emailVerified": types.M{"type": "Boolean"},
		},
		"classLevelPermissions": types.M{
			"find": types.M{"*": true},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_convertSchemaToAdapterSchema(t *testing.T) {
	// TODO
}

func Test_convertAdapterSchemaToParseSchema(t *testing.T) {
	var schema types.M
	var result types.M
	var expect types.M
	/************************************************************/
	schema = nil
	result = convertAdapterSchemaToParseSchema(schema)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/************************************************************/
	schema = types.M{}
	result = convertAdapterSchemaToParseSchema(schema)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/************************************************************/
	schema = types.M{
		"fields": types.M{
			"_rperm": types.M{"type": "Array"},
			"_wperm": types.M{"type": "Array"},
			"key":    types.M{"type": "String"},
		},
	}
	result = convertAdapterSchemaToParseSchema(schema)
	expect = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
			"ACL": types.M{"type": "ACL"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/************************************************************/
	schema = types.M{
		"className": "_User",
		"fields": types.M{
			"_rperm":           types.M{"type": "Array"},
			"_wperm":           types.M{"type": "Array"},
			"key":              types.M{"type": "String"},
			"authData":         types.M{"type": "String"},
			"_hashed_password": types.M{"type": "String"},
		},
	}
	result = convertAdapterSchemaToParseSchema(schema)
	expect = types.M{
		"className": "_User",
		"fields": types.M{
			"key":      types.M{"type": "String"},
			"ACL":      types.M{"type": "ACL"},
			"password": types.M{"type": "String"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/************************************************************/
	schema = types.M{
		"className": "other",
		"fields": types.M{
			"_rperm":           types.M{"type": "Array"},
			"_wperm":           types.M{"type": "Array"},
			"key":              types.M{"type": "String"},
			"authData":         types.M{"type": "String"},
			"_hashed_password": types.M{"type": "String"},
		},
	}
	result = convertAdapterSchemaToParseSchema(schema)
	expect = types.M{
		"className": "other",
		"fields": types.M{
			"key":              types.M{"type": "String"},
			"ACL":              types.M{"type": "ACL"},
			"authData":         types.M{"type": "String"},
			"_hashed_password": types.M{"type": "String"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
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
