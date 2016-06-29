package orm

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/storage"
	"github.com/lfq7413/tomato/storage/mongo"
	"github.com/lfq7413/tomato/types"
)

func Test_AddClassIfNotExists(t *testing.T) {
	adapter := getAdapter()
	schama := getSchema()
	var class types.M
	var className string
	var fields types.M
	var classLevelPermissions types.M
	var result types.M
	var err error
	var expect interface{}
	/************************************************************/
	className = "post"
	fields = types.M{
		"key": types.M{"type": "String"},
	}
	classLevelPermissions = nil
	result, err = schama.AddClassIfNotExists(className, fields, classLevelPermissions)
	expect = types.M{
		"className": className,
		"fields": types.M{
			"key":       types.M{"type": "String"},
			"objectId":  types.M{"type": "String"},
			"updatedAt": types.M{"type": "Date"},
			"createdAt": types.M{"type": "Date"},
			"ACL":       types.M{"type": "ACL"},
		},
		"classLevelPermissions": types.M{
			"find":     types.M{"*": true},
			"get":      types.M{"*": true},
			"create":   types.M{"*": true},
			"update":   types.M{"*": true},
			"delete":   types.M{"*": true},
			"addField": types.M{"*": true},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
	/************************************************************/
	className = "post"
	class = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	adapter.CreateClass(className, class)
	className = "post"
	fields = types.M{
		"key": types.M{"type": "String"},
	}
	classLevelPermissions = nil
	result, err = schama.AddClassIfNotExists(className, fields, classLevelPermissions)
	expect = errs.E(errs.InvalidClassName, "Class "+className+" already exists.")
	if err == nil || reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
}

func Test_UpdateClass(t *testing.T) {
	// TODO
}

func Test_deleteField(t *testing.T) {
	adapter := getAdapter()
	schama := getSchema()
	var class types.M
	var fieldName string
	var className string
	var err error
	var expect error
	var r1 types.M
	var r2 []types.M
	/************************************************************/
	fieldName = "abc"
	className = "@abc"
	err = schama.deleteField(fieldName, className)
	expect = errs.E(errs.InvalidClassName, InvalidClassNameMessage(className))
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	fieldName = "@abc"
	className = "abc"
	err = schama.deleteField(fieldName, className)
	expect = errs.E(errs.InvalidKeyName, "invalid field name: @abc")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	fieldName = "objectId"
	className = "abc"
	err = schama.deleteField(fieldName, className)
	expect = errs.E(errs.ChangedImmutableFieldError, "field objectId cannot be changed")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	fieldName = "key"
	className = "abc"
	err = schama.deleteField(fieldName, className)
	expect = errs.E(errs.InvalidClassName, "Class abc does not exist.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	className = "abc"
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass(className, class)
	fieldName = "key"
	className = "abc"
	err = schama.deleteField(fieldName, className)
	expect = errs.E(errs.ClassNotEmpty, "Field key does not exist, cannot delete.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	adapter.DeleteAllClasses()
	/************************************************************/
	className = "abc"
	class = types.M{
		"fields": types.M{
			"key":  types.M{"type": "String"},
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass(className, class)
	class = types.M{
		"objectId": "1024",
		"key":      "hello",
		"key1":     "world",
	}
	adapter.CreateObject(className, types.M{}, class)

	fieldName = "key"
	className = "abc"
	err = schama.deleteField(fieldName, className)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	// 检查 schema
	r1, _ = adapter.GetClass(className)
	class = types.M{
		"key1":      types.M{"type": "String"},
		"objectId":  types.M{"type": "String"},
		"updatedAt": types.M{"type": "Date"},
		"createdAt": types.M{"type": "Date"},
		"ACL":       types.M{"type": "ACL"},
	}
	if reflect.DeepEqual(class, r1["fields"]) == false {
		t.Error("expect:", class, "result:", r1)
	}
	// 检查数据
	r2, _ = adapter.Find(className, types.M{}, types.M{}, types.M{})
	class = types.M{
		"objectId": "1024",
		"key1":     "world",
	}
	if r2 == nil || len(r2) == 0 || reflect.DeepEqual(class, r2[0]) == false {
		t.Error("expect:", class, "result:", r2)
	}
	adapter.DeleteAllClasses()
	/************************************************************/
	className = "abc"
	class = types.M{
		"fields": types.M{
			"key":  types.M{"type": "Relation", "targetClass": "user"},
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass(className, class)
	className = "_Join:key:abc"
	class = types.M{
		"fields": types.M{
			"relatedId": types.M{"type": "String"},
			"owningId":  types.M{"type": "String"},
		},
	}
	adapter.CreateClass(className, class)
	className = "abc"
	class = types.M{
		"objectId": "1024",
		"key":      "hello",
		"key1":     "world",
	}
	adapter.CreateObject(className, types.M{}, class)
	className = "_Join:key:abc"
	class = types.M{
		"objectId":  "1024",
		"relatedId": "123",
		"owningId":  "456",
	}
	adapter.CreateObject(className, types.M{}, class)

	fieldName = "key"
	className = "abc"
	err = schama.deleteField(fieldName, className)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	// 检查 schema
	className = "abc"
	r1, _ = adapter.GetClass(className)
	class = types.M{
		"key1":      types.M{"type": "String"},
		"objectId":  types.M{"type": "String"},
		"updatedAt": types.M{"type": "Date"},
		"createdAt": types.M{"type": "Date"},
		"ACL":       types.M{"type": "ACL"},
	}
	if reflect.DeepEqual(class, r1["fields"]) == false {
		t.Error("expect:", class, "result:", r1)
	}
	// 检查 schema
	className = "_Join:key:abc"
	r1, _ = adapter.GetClass(className)
	class = types.M{}
	if reflect.DeepEqual(class, r1) == false {
		t.Error("expect:", class, "result:", r1)
	}
	// 检查数据
	className = "abc"
	r2, _ = adapter.Find(className, types.M{}, types.M{}, types.M{})
	class = types.M{
		"objectId": "1024",
		"key1":     "world",
	}
	if r2 == nil || len(r2) == 0 || reflect.DeepEqual(class, r2[0]) == false {
		t.Error("expect:", class, "result:", r2)
	}
	// 检查 Join 数据
	className = "_Join:key:abc"
	r2, _ = adapter.Find(className, types.M{}, types.M{}, types.M{})
	if r2 != nil && reflect.DeepEqual([]types.M{}, r2) == false {
		t.Error("expect:", class, "result:", r2)
	}
	adapter.DeleteAllClasses()
}

func Test_validateObject(t *testing.T) {
	// EnforceClassExists
	// TODO
}

func Test_testBaseCLP(t *testing.T) {
	schama := getSchema()
	var className string
	var aclGroup []string
	var operation string
	var ok bool
	var expect bool
	/************************************************************/
	schama.perms = nil
	className = "post"
	aclGroup = nil
	operation = "get"
	ok = schama.testBaseCLP(className, aclGroup, operation)
	expect = true
	if reflect.DeepEqual(expect, ok) == false {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	schama.perms = types.M{}
	className = "post"
	aclGroup = nil
	operation = "get"
	ok = schama.testBaseCLP(className, aclGroup, operation)
	expect = true
	if reflect.DeepEqual(expect, ok) == false {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	schama.perms = types.M{
		"post": types.M{},
	}
	className = "post"
	aclGroup = nil
	operation = "get"
	ok = schama.testBaseCLP(className, aclGroup, operation)
	expect = true
	if reflect.DeepEqual(expect, ok) == false {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	schama.perms = types.M{
		"post": types.M{
			"get": types.M{"*": true},
		},
	}
	className = "post"
	aclGroup = nil
	operation = "get"
	ok = schama.testBaseCLP(className, aclGroup, operation)
	expect = true
	if reflect.DeepEqual(expect, ok) == false {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	schama.perms = types.M{
		"post": types.M{
			"get": types.M{},
		},
	}
	className = "post"
	aclGroup = nil
	operation = "get"
	ok = schama.testBaseCLP(className, aclGroup, operation)
	expect = false
	if reflect.DeepEqual(expect, ok) == false {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	schama.perms = types.M{
		"post": types.M{
			"get": types.M{},
		},
	}
	className = "post"
	aclGroup = []string{"role:1024"}
	operation = "get"
	ok = schama.testBaseCLP(className, aclGroup, operation)
	expect = false
	if reflect.DeepEqual(expect, ok) == false {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	schama.perms = types.M{
		"post": types.M{
			"get": types.M{"role:1024": true},
		},
	}
	className = "post"
	aclGroup = []string{"role:1024"}
	operation = "get"
	ok = schama.testBaseCLP(className, aclGroup, operation)
	expect = true
	if reflect.DeepEqual(expect, ok) == false {
		t.Error("expect:", expect, "result:", ok)
	}
}

func Test_validatePermission(t *testing.T) {
	schama := getSchema()
	var className string
	var aclGroup []string
	var operation string
	var err error
	var expect error
	/************************************************************/
	schama.perms = types.M{
		"post": types.M{
			"create": types.M{"role:1024": true},
		},
	}
	className = "post"
	aclGroup = []string{"role:abc"}
	operation = "create"
	err = schama.validatePermission(className, aclGroup, operation)
	expect = errs.E(errs.OperationForbidden, "Permission denied for this action.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	schama.perms = types.M{
		"post": types.M{
			"get": types.M{"role:1024": true},
		},
	}
	className = "post"
	aclGroup = []string{"role:abc"}
	operation = "get"
	err = schama.validatePermission(className, aclGroup, operation)
	expect = errs.E(errs.OperationForbidden, "Permission denied for this action.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	schama.perms = types.M{
		"post": types.M{
			"get":            types.M{"role:1024": true},
			"readUserFields": types.S{"key"},
		},
	}
	className = "post"
	aclGroup = []string{"role:abc"}
	operation = "get"
	err = schama.validatePermission(className, aclGroup, operation)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
}

func Test_EnforceClassExists(t *testing.T) {
	// TODO
}

func Test_validateNewClass(t *testing.T) {
	adapter := getAdapter()
	schama := getSchema()
	var class types.M
	var className string
	var fields types.M
	var classLevelPermissions types.M
	var err error
	var expect error
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	schama.reloadData()
	className = "post"
	fields = nil
	classLevelPermissions = nil
	err = schama.validateNewClass(className, fields, classLevelPermissions)
	expect = errs.E(errs.InvalidClassName, "Class post already exists.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	schama.reloadData()
	className = "@post"
	fields = nil
	classLevelPermissions = nil
	err = schama.validateNewClass(className, fields, classLevelPermissions)
	expect = errs.E(errs.InvalidClassName, InvalidClassNameMessage("@post"))
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	schama.reloadData()
	className = "user"
	fields = types.M{
		"key": types.M{"type": "String"},
	}
	classLevelPermissions = nil
	err = schama.validateNewClass(className, fields, classLevelPermissions)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	adapter.DeleteAllClasses()
}

func Test_validateSchemaData(t *testing.T) {
	schama := getSchema()
	var className string
	var fields types.M
	var classLevelPermissions types.M
	var existingFieldNames []string
	var err error
	var expect error
	/************************************************************/
	className = "post"
	fields = nil
	classLevelPermissions = nil
	existingFieldNames = nil
	err = schama.validateSchemaData(className, fields, classLevelPermissions, existingFieldNames)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	className = "post"
	fields = types.M{
		"key": types.M{"type": "String"},
	}
	classLevelPermissions = nil
	existingFieldNames = nil
	err = schama.validateSchemaData(className, fields, classLevelPermissions, existingFieldNames)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	className = "post"
	fields = types.M{
		"key":  types.M{"type": "String"},
		"key2": types.M{"type": "String"},
	}
	classLevelPermissions = nil
	existingFieldNames = []string{"key"}
	err = schama.validateSchemaData(className, fields, classLevelPermissions, existingFieldNames)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	className = "post"
	fields = types.M{
		"key":      types.M{"type": "String"},
		"objectId": types.M{"type": "String"},
	}
	classLevelPermissions = nil
	existingFieldNames = []string{"key"}
	err = schama.validateSchemaData(className, fields, classLevelPermissions, existingFieldNames)
	expect = errs.E(errs.ChangedImmutableFieldError, "field objectId cannot be added")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	className = "post"
	fields = types.M{
		"key":  types.M{"type": "String"},
		"key2": types.M{"type": "Other"},
	}
	classLevelPermissions = nil
	existingFieldNames = []string{"key"}
	err = schama.validateSchemaData(className, fields, classLevelPermissions, existingFieldNames)
	expect = errs.E(errs.IncorrectType, "invalid field type: Other")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	className = "_User"
	fields = types.M{
		"key":  types.M{"type": "String"},
		"key2": types.M{"type": "String"},
	}
	classLevelPermissions = nil
	existingFieldNames = []string{"key"}
	err = schama.validateSchemaData(className, fields, classLevelPermissions, existingFieldNames)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	className = "_User"
	fields = types.M{
		"key":  types.M{"type": "String"},
		"key2": types.M{"type": "String"},
		"loc":  types.M{"type": "GeoPoint"},
	}
	classLevelPermissions = nil
	existingFieldNames = []string{"key"}
	err = schama.validateSchemaData(className, fields, classLevelPermissions, existingFieldNames)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	className = "_User"
	fields = types.M{
		"key":  types.M{"type": "String"},
		"key2": types.M{"type": "String"},
		"loc":  types.M{"type": "GeoPoint"},
		"loc2": types.M{"type": "GeoPoint"},
	}
	classLevelPermissions = nil
	existingFieldNames = []string{"key"}
	err = schama.validateSchemaData(className, fields, classLevelPermissions, existingFieldNames)
	expect = errs.E(errs.IncorrectType, "currently, only one GeoPoint field may exist in an object. Adding loc when loc2 already exists.")
	expect2 := errs.E(errs.IncorrectType, "currently, only one GeoPoint field may exist in an object. Adding loc2 when loc already exists.")
	if reflect.DeepEqual(expect, err) == false && reflect.DeepEqual(expect2, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
}

func Test_validateRequiredColumns(t *testing.T) {
	schama := getSchema()
	var className string
	var object types.M
	var query types.M
	var err error
	var expect error
	/************************************************************/
	className = "user"
	object = nil
	query = nil
	err = schama.validateRequiredColumns(className, object, query)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	className = "_Role"
	object = types.M{
		"name": "joe",
	}
	query = nil
	err = schama.validateRequiredColumns(className, object, query)
	expect = errs.E(errs.IncorrectType, "ACL is required.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	className = "_Role"
	object = types.M{
		"name": "joe",
		"ACL": types.M{
			"__op": "Delete",
		},
	}
	query = types.M{
		"objectId": "1024",
	}
	err = schama.validateRequiredColumns(className, object, query)
	expect = errs.E(errs.IncorrectType, "ACL is required.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	className = "_Product"
	object = types.M{
		"productIdentifier": "1024",
		"icon":              "a.jpg",
		"order":             "name",
		"title":             "tomato",
	}
	query = nil
	err = schama.validateRequiredColumns(className, object, query)
	expect = errs.E(errs.IncorrectType, "subtitle is required.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	className = "_Product"
	object = types.M{
		"productIdentifier": "1024",
		"icon":              "a.jpg",
		"order":             "name",
		"title":             "tomato",
		"subtitle": types.M{
			"__op": "Delete",
		},
	}
	query = types.M{
		"objectId": "1024",
	}
	err = schama.validateRequiredColumns(className, object, query)
	expect = errs.E(errs.IncorrectType, "subtitle is required.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
}

func Test_enforceFieldExists(t *testing.T) {
	adapter := getAdapter()
	schama := getSchema()
	var class types.M
	var className string
	var fieldName string
	var fieldtype types.M
	var err error
	var expect error
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	className = "post"
	fieldName = "key2"
	fieldtype = types.M{
		"type": "String",
	}
	err = schama.enforceFieldExists(className, fieldName, fieldtype, true)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	className = "post"
	fieldName = "key2.key"
	fieldtype = types.M{
		"type": "String",
	}
	err = schama.enforceFieldExists(className, fieldName, fieldtype, true)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	className = "post"
	fieldName = "@key2"
	fieldtype = types.M{
		"type": "String",
	}
	err = schama.enforceFieldExists(className, fieldName, fieldtype, true)
	expect = errs.E(errs.InvalidKeyName, "Invalid field name: "+fieldName)
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	className = "post"
	fieldName = "key2"
	fieldtype = nil
	err = schama.enforceFieldExists(className, fieldName, fieldtype, true)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	className = "post"
	fieldName = "key2"
	fieldtype = types.M{}
	err = schama.enforceFieldExists(className, fieldName, fieldtype, true)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	className = "post"
	fieldName = "key1"
	fieldtype = types.M{
		"type": "Number",
	}
	err = schama.enforceFieldExists(className, fieldName, fieldtype, true)
	expect = errs.E(errs.IncorrectType, "schema mismatch for post.key1; expected String but got Number")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	className = "post"
	fieldName = "key1"
	fieldtype = types.M{
		"type": "String",
	}
	err = schama.enforceFieldExists(className, fieldName, fieldtype, true)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	adapter.DeleteAllClasses()
}

func Test_setPermissions(t *testing.T) {
	adapter := getAdapter()
	schama := getSchema()
	var class types.M
	var className string
	var perms types.M
	var newSchema types.M
	var err error
	var expect interface{}
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("user", class)
	className = "class"
	perms = types.M{
		"get": types.M{"*": true},
	}
	newSchema = nil
	err = schama.setPermissions(className, perms, newSchema)
	expect = errors.New("not found")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("user", class)
	className = "post"
	perms = types.M{
		"get": types.M{"*": true},
	}
	newSchema = nil
	err = schama.setPermissions(className, perms, newSchema)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	expect = types.M{
		"get":      types.M{"*": true},
		"create":   types.M{},
		"find":     types.M{},
		"update":   types.M{},
		"delete":   types.M{},
		"addField": types.M{},
	}
	if reflect.DeepEqual(expect, schama.perms[className]) == false {
		t.Error("expect:", expect, "result:", schama.perms[className])
	}
	adapter.DeleteAllClasses()
}

func Test_HasClass(t *testing.T) {
	adapter := getAdapter()
	schama := getSchema()
	var class types.M
	var className string
	var result bool
	var expect bool
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("user", class)
	/************************************************************/
	className = "class"
	result = schama.HasClass(className)
	expect = false
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/************************************************************/
	className = "post"
	result = schama.HasClass(className)
	expect = true
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}

	adapter.DeleteAllClasses()
}

func Test_getExpectedType(t *testing.T) {
	adapter := getAdapter()
	schama := getSchema()
	var class types.M
	var className string
	var fieldName string
	var result types.M
	var expect types.M
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("user", class)
	schama.reloadData()
	className = "class"
	fieldName = "field"
	result = schama.getExpectedType(className, fieldName)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("user", class)
	schama.reloadData()
	className = "post"
	fieldName = "field"
	result = schama.getExpectedType(className, fieldName)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("user", class)
	schama.reloadData()
	className = "post"
	fieldName = "key1"
	result = schama.getExpectedType(className, fieldName)
	expect = types.M{"type": "String"}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	adapter.DeleteAllClasses()
}

func Test_reloadData(t *testing.T) {
	adapter := getAdapter()
	schama := getSchema()
	var class types.M
	var expect types.M
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("_User", class)
	schama.reloadData()
	expect = types.M{
		"post": types.M{
			"key1":      types.M{"type": "String"},
			"objectId":  types.M{"type": "String"},
			"updatedAt": types.M{"type": "Date"},
			"createdAt": types.M{"type": "Date"},
			"ACL":       types.M{"type": "ACL"},
		},
		"_User": types.M{
			"key1":          types.M{"type": "String"},
			"objectId":      types.M{"type": "String"},
			"updatedAt":     types.M{"type": "Date"},
			"createdAt":     types.M{"type": "Date"},
			"ACL":           types.M{"type": "ACL"},
			"username":      types.M{"type": "String"},
			"password":      types.M{"type": "String"},
			"email":         types.M{"type": "String"},
			"emailVerified": types.M{"type": "Boolean"},
		},
		"_PushStatus": types.M{
			"className": "_PushStatus",
			"fields": types.M{
				"objectId":      types.M{"type": "String"},
				"updatedAt":     types.M{"type": "Date"},
				"createdAt":     types.M{"type": "Date"},
				"ACL":           types.M{"type": "ACL"},
				"pushTime":      types.M{"type": "String"},
				"source":        types.M{"type": "String"},
				"query":         types.M{"type": "String"},
				"payload":       types.M{"type": "Object"},
				"title":         types.M{"type": "String"},
				"expiry":        types.M{"type": "Number"},
				"status":        types.M{"type": "String"},
				"numSent":       types.M{"type": "Number"},
				"numFailed":     types.M{"type": "Number"},
				"pushHash":      types.M{"type": "String"},
				"errorMessage":  types.M{"type": "Object"},
				"sentPerType":   types.M{"type": "Object"},
				"failedPerType": types.M{"type": "Object"},
			},
			"classLevelPermissions": types.M{},
		},
		"_Hooks": types.M{
			"className": "_Hooks",
			"fields": types.M{
				"objectId":  types.M{"type": "String"},
				"updatedAt": types.M{"type": "Date"},
				"createdAt": types.M{"type": "Date"},
				"ACL":       types.M{"type": "ACL"},
			},
			"classLevelPermissions": types.M{},
		},
		"_GlobalConfig": types.M{
			"className": "_GlobalConfig",
			"fields": types.M{
				"objectId":  types.M{"type": "String"},
				"updatedAt": types.M{"type": "Date"},
				"createdAt": types.M{"type": "Date"},
				"ACL":       types.M{"type": "ACL"},
			},
			"classLevelPermissions": types.M{},
		},
	}
	if reflect.DeepEqual(expect, schama.data) == false {
		t.Error("expect:", expect, "result:", schama.data)
	}
	expect = types.M{
		"post": types.M{
			"find":     types.M{"*": true},
			"get":      types.M{"*": true},
			"create":   types.M{"*": true},
			"update":   types.M{"*": true},
			"delete":   types.M{"*": true},
			"addField": types.M{"*": true},
		},
		"_User": types.M{
			"find":     types.M{"*": true},
			"get":      types.M{"*": true},
			"create":   types.M{"*": true},
			"update":   types.M{"*": true},
			"delete":   types.M{"*": true},
			"addField": types.M{"*": true},
		},
	}
	if reflect.DeepEqual(expect, schama.perms) == false {
		t.Error("expect:", expect, "result:", schama.perms)
	}
	adapter.DeleteAllClasses()
}

func Test_GetAllClasses(t *testing.T) {
	adapter := getAdapter()
	schama := getSchema()
	var class types.M
	var result []types.M
	var err error
	var expect []types.M
	/************************************************************/
	result, err = schama.GetAllClasses()
	expect = []types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	result, err = schama.GetAllClasses()
	expect = []types.M{
		types.M{
			"className": "post",
			"fields": types.M{
				"key1":      types.M{"type": "String"},
				"objectId":  types.M{"type": "String"},
				"updatedAt": types.M{"type": "Date"},
				"createdAt": types.M{"type": "Date"},
				"ACL":       types.M{"type": "ACL"},
			},
			"classLevelPermissions": types.M{
				"find":     types.M{"*": true},
				"get":      types.M{"*": true},
				"create":   types.M{"*": true},
				"update":   types.M{"*": true},
				"delete":   types.M{"*": true},
				"addField": types.M{"*": true},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("_User", class)
	result, err = schama.GetAllClasses()
	expect = []types.M{
		types.M{
			"className": "_User",
			"fields": types.M{
				"key1":          types.M{"type": "String"},
				"objectId":      types.M{"type": "String"},
				"updatedAt":     types.M{"type": "Date"},
				"createdAt":     types.M{"type": "Date"},
				"ACL":           types.M{"type": "ACL"},
				"username":      types.M{"type": "String"},
				"password":      types.M{"type": "String"},
				"email":         types.M{"type": "String"},
				"emailVerified": types.M{"type": "Boolean"},
			},
			"classLevelPermissions": types.M{
				"find":     types.M{"*": true},
				"get":      types.M{"*": true},
				"create":   types.M{"*": true},
				"update":   types.M{"*": true},
				"delete":   types.M{"*": true},
				"addField": types.M{"*": true},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("user", class)
	result, err = schama.GetAllClasses()
	expect = []types.M{
		types.M{
			"className": "post",
			"fields": types.M{
				"key1":      types.M{"type": "String"},
				"objectId":  types.M{"type": "String"},
				"updatedAt": types.M{"type": "Date"},
				"createdAt": types.M{"type": "Date"},
				"ACL":       types.M{"type": "ACL"},
			},
			"classLevelPermissions": types.M{
				"find":     types.M{"*": true},
				"get":      types.M{"*": true},
				"create":   types.M{"*": true},
				"update":   types.M{"*": true},
				"delete":   types.M{"*": true},
				"addField": types.M{"*": true},
			},
		},
		types.M{
			"className": "user",
			"fields": types.M{
				"key1":      types.M{"type": "String"},
				"objectId":  types.M{"type": "String"},
				"updatedAt": types.M{"type": "Date"},
				"createdAt": types.M{"type": "Date"},
				"ACL":       types.M{"type": "ACL"},
			},
			"classLevelPermissions": types.M{
				"find":     types.M{"*": true},
				"get":      types.M{"*": true},
				"create":   types.M{"*": true},
				"update":   types.M{"*": true},
				"delete":   types.M{"*": true},
				"addField": types.M{"*": true},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
}

func Test_GetOneSchema(t *testing.T) {
	adapter := getAdapter()
	schama := getSchema()
	var class types.M
	var className string
	var allowVolatileClasses bool
	var result types.M
	var err error
	var expect types.M
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("user", class)
	className = "class"
	allowVolatileClasses = false
	result, err = schama.GetOneSchema(className, allowVolatileClasses)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("user", class)
	className = "post"
	allowVolatileClasses = false
	result, err = schama.GetOneSchema(className, allowVolatileClasses)
	expect = types.M{
		"className": "post",
		"fields": types.M{
			"key1":      types.M{"type": "String"},
			"objectId":  types.M{"type": "String"},
			"updatedAt": types.M{"type": "Date"},
			"createdAt": types.M{"type": "Date"},
			"ACL":       types.M{"type": "ACL"},
		},
		"classLevelPermissions": types.M{
			"find":     types.M{"*": true},
			"get":      types.M{"*": true},
			"create":   types.M{"*": true},
			"update":   types.M{"*": true},
			"delete":   types.M{"*": true},
			"addField": types.M{"*": true},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("user", class)
	className = "post"
	allowVolatileClasses = true
	result, err = schama.GetOneSchema(className, allowVolatileClasses)
	expect = types.M{
		"className": "post",
		"fields": types.M{
			"key1":      types.M{"type": "String"},
			"objectId":  types.M{"type": "String"},
			"updatedAt": types.M{"type": "Date"},
			"createdAt": types.M{"type": "Date"},
			"ACL":       types.M{"type": "ACL"},
		},
		"classLevelPermissions": types.M{
			"find":     types.M{"*": true},
			"get":      types.M{"*": true},
			"create":   types.M{"*": true},
			"update":   types.M{"*": true},
			"delete":   types.M{"*": true},
			"addField": types.M{"*": true},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	className = "_PushStatus"
	allowVolatileClasses = true
	result, err = schama.GetOneSchema(className, allowVolatileClasses)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	schama.data = types.M{
		"_PushStatus": types.M{
			"fields": types.M{
				"key1": types.M{"type": "String"},
			},
		},
	}
	className = "_PushStatus"
	allowVolatileClasses = true
	result, err = schama.GetOneSchema(className, allowVolatileClasses)
	expect = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
}

////////////////////////////////////////////////////////////

func Test_thenValidateRequiredColumns(t *testing.T) {
	// 测试用例与 validateRequiredColumns 相同
}

func Test_getType(t *testing.T) {
	var object interface{}
	var result types.M
	var err error
	var expect interface{}
	/************************************************************/
	object = true
	result, err = getType(object)
	expect = types.M{"type": "Boolean"}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = "hello"
	result, err = getType(object)
	expect = types.M{"type": "String"}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = 1024
	result, err = getType(object)
	expect = types.M{"type": "Number"}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = 10.24
	result, err = getType(object)
	expect = types.M{"type": "Number"}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = types.M{
		"__type": "Date",
		"iso":    "abc",
	}
	result, err = getType(object)
	expect = types.M{"type": "Date"}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = map[string]interface{}{
		"__type": "File",
		"name":   "abc",
	}
	result, err = getType(object)
	expect = types.M{"type": "File"}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = types.S{1, 2, 3}
	result, err = getType(object)
	expect = types.M{"type": "Array"}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = []interface{}{1, 2, 3}
	result, err = getType(object)
	expect = types.M{"type": "Array"}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = time.Now()
	result, err = getType(object)
	expect = errs.E(errs.IncorrectType, "bad obj. can not get type")
	if err == nil || reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
}

func Test_getObjectType(t *testing.T) {
	var object interface{}
	var result types.M
	var err error
	var expect interface{}
	/************************************************************/
	object = []interface{}{1, 2, 3}
	result, err = getObjectType(object)
	expect = types.M{"type": "Array"}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = types.M{
		"__type":    "Pointer",
		"className": "abc",
	}
	result, err = getObjectType(object)
	expect = types.M{
		"type":        "Pointer",
		"targetClass": "abc",
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = types.M{
		"__type": "File",
		"name":   "abc",
	}
	result, err = getObjectType(object)
	expect = types.M{"type": "File"}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = types.M{
		"__type": "Date",
		"iso":    "abc",
	}
	result, err = getObjectType(object)
	expect = types.M{"type": "Date"}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = types.M{
		"__type":    "GeoPoint",
		"latitude":  10,
		"longitude": 10,
	}
	result, err = getObjectType(object)
	expect = types.M{"type": "GeoPoint"}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = types.M{
		"__type": "Bytes",
		"base64": "abc",
	}
	result, err = getObjectType(object)
	expect = types.M{"type": "Bytes"}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = types.M{
		"__type": "Other",
	}
	result, err = getObjectType(object)
	expect = errs.E(errs.IncorrectType, "This is not a valid Other")
	if err == nil || reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = types.M{
		"$ne": types.M{
			"__type":    "Pointer",
			"className": "abc",
		},
	}
	result, err = getObjectType(object)
	expect = types.M{
		"type":        "Pointer",
		"targetClass": "abc",
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = types.M{
		"__op": "Increment",
	}
	result, err = getObjectType(object)
	expect = types.M{"type": "Number"}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = types.M{
		"__op": "Delete",
	}
	result, err = getObjectType(object)
	expect = nil
	if err != nil || result != nil {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = types.M{
		"__op": "Add",
	}
	result, err = getObjectType(object)
	expect = types.M{"type": "Array"}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = types.M{
		"__op": "AddUnique",
	}
	result, err = getObjectType(object)
	expect = types.M{"type": "Array"}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = types.M{
		"__op": "Remove",
	}
	result, err = getObjectType(object)
	expect = types.M{"type": "Array"}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = types.M{
		"__op": "AddRelation",
		"objects": types.S{
			types.M{
				"className": "abc",
			},
		},
	}
	result, err = getObjectType(object)
	expect = types.M{
		"type":        "Relation",
		"targetClass": "abc",
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = types.M{
		"__op": "RemoveRelation",
		"objects": types.S{
			types.M{
				"className": "abc",
			},
		},
	}
	result, err = getObjectType(object)
	expect = types.M{
		"type":        "Relation",
		"targetClass": "abc",
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = types.M{
		"__op": "Batch",
		"ops": types.S{
			types.M{
				"__type": "File",
				"name":   "abc",
			},
		},
	}
	result, err = getObjectType(object)
	expect = types.M{
		"type": "File",
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = types.M{"__op": "Other"}
	result, err = getObjectType(object)
	expect = errs.E(errs.IncorrectType, "unexpected op: Other")
	if err == nil || reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/************************************************************/
	object = types.M{"key": "value"}
	result, err = getObjectType(object)
	expect = types.M{"type": "object"}
	if err != nil || reflect.DeepEqual(expect, expect) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
}

func Test_ClassNameIsValid(t *testing.T) {
	var className string
	var ok bool
	var expect bool
	/************************************************************/
	className = "_User"
	ok = ClassNameIsValid(className)
	expect = true
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	className = "_Installation"
	ok = ClassNameIsValid(className)
	expect = true
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	className = "_Role"
	ok = ClassNameIsValid(className)
	expect = true
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	className = "_Session"
	ok = ClassNameIsValid(className)
	expect = true
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	className = "_Join:abc:123"
	ok = ClassNameIsValid(className)
	expect = true
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	className = "abc"
	ok = ClassNameIsValid(className)
	expect = true
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
}

func Test_InvalidClassNameMessage(t *testing.T) {
	var className string
	var result string
	var expect string
	/************************************************************/
	className = "abc"
	result = InvalidClassNameMessage(className)
	expect = "Invalid classname: abc, classnames can only have alphanumeric characters and _, and must start with an alpha character "
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_joinClassIsValid(t *testing.T) {
	var className string
	var ok bool
	var expect bool
	/************************************************************/
	className = "_Join:abc:def"
	ok = joinClassIsValid(className)
	expect = true
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	className = "_Join:abc123:def123"
	ok = joinClassIsValid(className)
	expect = true
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	className = "_Join:_abc123:def_123"
	ok = joinClassIsValid(className)
	expect = true
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	className = "abc"
	ok = joinClassIsValid(className)
	expect = false
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	className = "_Join:@123:!def"
	ok = joinClassIsValid(className)
	expect = false
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
}

func Test_fieldNameIsValid(t *testing.T) {
	var fieldName string
	var ok bool
	var expect bool
	/************************************************************/
	fieldName = "abc_123"
	ok = fieldNameIsValid(fieldName)
	expect = true
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	fieldName = "abc123"
	ok = fieldNameIsValid(fieldName)
	expect = true
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	fieldName = "123abc"
	ok = fieldNameIsValid(fieldName)
	expect = false
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	fieldName = "*abc"
	ok = fieldNameIsValid(fieldName)
	expect = false
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	fieldName = "abc@123"
	ok = fieldNameIsValid(fieldName)
	expect = false
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
}

func Test_fieldNameIsValidForClass(t *testing.T) {
	var fieldName string
	var className string
	var ok bool
	var expect bool
	/************************************************************/
	fieldName = ""
	className = ""
	ok = fieldNameIsValidForClass(fieldName, className)
	expect = false
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	fieldName = "abc"
	className = ""
	ok = fieldNameIsValidForClass(fieldName, className)
	expect = true
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	fieldName = "objectId"
	className = ""
	ok = fieldNameIsValidForClass(fieldName, className)
	expect = false
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	fieldName = "abc"
	className = "_User"
	ok = fieldNameIsValidForClass(fieldName, className)
	expect = true
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	fieldName = "username"
	className = "_User"
	ok = fieldNameIsValidForClass(fieldName, className)
	expect = false
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
	/************************************************************/
	fieldName = "key"
	className = "class"
	ok = fieldNameIsValidForClass(fieldName, className)
	expect = true
	if ok != expect {
		t.Error("expect:", expect, "result:", ok)
	}
}

func Test_fieldTypeIsInvalid(t *testing.T) {
	var tp types.M
	var err error
	var expect error
	/************************************************************/
	tp = nil
	err = fieldTypeIsInvalid(tp)
	expect = errs.E(errs.InvalidJSON, "invalid JSON")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	tp = types.M{}
	err = fieldTypeIsInvalid(tp)
	expect = errs.E(errs.InvalidJSON, "invalid JSON")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	tp = types.M{"type": 1024}
	err = fieldTypeIsInvalid(tp)
	expect = errs.E(errs.InvalidJSON, "invalid JSON")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	tp = types.M{"type": "Pointer"}
	err = fieldTypeIsInvalid(tp)
	expect = errs.E(errs.MissingRequiredFieldError, "type Pointer needs a class name")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	tp = types.M{
		"type":        "Pointer",
		"targetClass": 1024,
	}
	err = fieldTypeIsInvalid(tp)
	expect = errs.E(errs.InvalidJSON, "invalid JSON")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	tp = types.M{
		"type":        "Pointer",
		"targetClass": "@abc",
	}
	err = fieldTypeIsInvalid(tp)
	expect = errs.E(errs.InvalidClassName, "Invalid classname: @abc, classnames can only have alphanumeric characters and _, and must start with an alpha character ")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	tp = types.M{
		"type":        "Pointer",
		"targetClass": "abc",
	}
	err = fieldTypeIsInvalid(tp)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	tp = types.M{"type": "Relation"}
	err = fieldTypeIsInvalid(tp)
	expect = errs.E(errs.MissingRequiredFieldError, "type Relation needs a class name")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	tp = types.M{
		"type":        "Relation",
		"targetClass": 1024,
	}
	err = fieldTypeIsInvalid(tp)
	expect = errs.E(errs.InvalidJSON, "invalid JSON")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	tp = types.M{
		"type":        "Relation",
		"targetClass": "@abc",
	}
	err = fieldTypeIsInvalid(tp)
	expect = errs.E(errs.InvalidClassName, "Invalid classname: @abc, classnames can only have alphanumeric characters and _, and must start with an alpha character ")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	tp = types.M{
		"type":        "Relation",
		"targetClass": "abc",
	}
	err = fieldTypeIsInvalid(tp)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	tp = types.M{
		"type": "Number",
	}
	err = fieldTypeIsInvalid(tp)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	tp = types.M{
		"type": "String",
	}
	err = fieldTypeIsInvalid(tp)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	tp = types.M{
		"type": "Boolean",
	}
	err = fieldTypeIsInvalid(tp)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	tp = types.M{
		"type": "Date",
	}
	err = fieldTypeIsInvalid(tp)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	tp = types.M{
		"type": "Object",
	}
	err = fieldTypeIsInvalid(tp)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	tp = types.M{
		"type": "Array",
	}
	err = fieldTypeIsInvalid(tp)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	tp = types.M{
		"type": "GeoPoint",
	}
	err = fieldTypeIsInvalid(tp)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	tp = types.M{
		"type": "File",
	}
	err = fieldTypeIsInvalid(tp)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	tp = types.M{
		"type": "Other",
	}
	err = fieldTypeIsInvalid(tp)
	expect = errs.E(errs.IncorrectType, "invalid field type: Other")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
}

func Test_validateCLP(t *testing.T) {
	var perms types.M
	var fields types.M
	var err error
	var expect error
	/************************************************************/
	perms = nil
	fields = nil
	err = validateCLP(perms, fields)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	perms = types.M{}
	fields = nil
	err = validateCLP(perms, fields)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	perms = types.M{
		"get": types.M{"012345678901234567890123": true},
	}
	fields = nil
	err = validateCLP(perms, fields)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	perms = types.M{
		"find": types.M{"012345678901234567890123": true},
	}
	fields = nil
	err = validateCLP(perms, fields)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	perms = types.M{
		"create": types.M{"012345678901234567890123": true},
	}
	fields = nil
	err = validateCLP(perms, fields)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	perms = types.M{
		"update": types.M{"012345678901234567890123": true},
	}
	fields = nil
	err = validateCLP(perms, fields)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	perms = types.M{
		"delete": types.M{"012345678901234567890123": true},
	}
	fields = nil
	err = validateCLP(perms, fields)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	perms = types.M{
		"addField": types.M{"012345678901234567890123": true},
	}
	fields = nil
	err = validateCLP(perms, fields)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	perms = types.M{
		"other": types.M{"012345678901234567890123": true},
	}
	fields = nil
	err = validateCLP(perms, fields)
	expect = errs.E(errs.InvalidJSON, "other is not a valid operation for class level permissions")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	perms = types.M{
		"readUserFields": types.S{"key1", "key2"},
	}
	fields = types.M{
		"key1": types.M{
			"type":        "Pointer",
			"targetClass": "_User",
		},
		"key2": types.M{
			"type":        "Pointer",
			"targetClass": "_User",
		},
	}
	err = validateCLP(perms, fields)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	perms = types.M{
		"writeUserFields": types.S{"key1", "key2"},
	}
	fields = types.M{
		"key1": types.M{
			"type":        "Pointer",
			"targetClass": "_User",
		},
		"key2": types.M{
			"type":        "Pointer",
			"targetClass": "_User",
		},
	}
	err = validateCLP(perms, fields)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	perms = types.M{
		"readUserFields": "hello",
	}
	fields = nil
	err = validateCLP(perms, fields)
	expect = errs.E(errs.InvalidJSON, "this perms[operation] is not a valid value for class level permissions readUserFields")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	perms = types.M{
		"readUserFields": types.S{"key1", "key2"},
	}
	fields = nil
	err = validateCLP(perms, fields)
	expect = errs.E(errs.InvalidJSON, "key1 is not a valid column for class level pointer permissions readUserFields")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	perms = types.M{
		"readUserFields": types.S{"key1", "key2"},
	}
	fields = types.M{}
	err = validateCLP(perms, fields)
	expect = errs.E(errs.InvalidJSON, "key1 is not a valid column for class level pointer permissions readUserFields")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	perms = types.M{
		"readUserFields": types.S{"key1", "key2"},
	}
	fields = types.M{"key1": 1024}
	err = validateCLP(perms, fields)
	expect = errs.E(errs.InvalidJSON, "key1 is not a valid column for class level pointer permissions readUserFields")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	perms = types.M{
		"readUserFields": types.S{"key1", "key2"},
	}
	fields = types.M{
		"key1": types.M{
			"type": "Other",
		},
	}
	err = validateCLP(perms, fields)
	expect = errs.E(errs.InvalidJSON, "key1 is not a valid column for class level pointer permissions readUserFields")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	perms = types.M{
		"get": types.M{"abc": true},
	}
	fields = nil
	err = validateCLP(perms, fields)
	expect = errs.E(errs.InvalidJSON, "abc is not a valid key for class level permissions")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	perms = types.M{
		"get": types.M{"role:abc": false},
	}
	fields = nil
	err = validateCLP(perms, fields)
	expect = errs.E(errs.InvalidJSON, "false is not a valid value for class level permissions get:role:abc:false")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	perms = types.M{
		"get": types.M{"role:abc": "hello"},
	}
	fields = nil
	err = validateCLP(perms, fields)
	expect = errs.E(errs.InvalidJSON, "this perm is not a valid value for class level permissions get:role:abc:perm")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
}

func Test_verifyPermissionKey(t *testing.T) {
	var key string
	var err error
	var expect error
	/************************************************************/
	key = "0123456789abcdefghij0123"
	err = verifyPermissionKey(key)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	key = "role:1024"
	err = verifyPermissionKey(key)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	key = "*"
	err = verifyPermissionKey(key)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	key = "abcd"
	err = verifyPermissionKey(key)
	expect = errs.E(errs.InvalidJSON, key+" is not a valid key for class level permissions")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	key = "*abc"
	err = verifyPermissionKey(key)
	expect = errs.E(errs.InvalidJSON, key+" is not a valid key for class level permissions")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	key = "role:*abc"
	err = verifyPermissionKey(key)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	key = "@mail"
	err = verifyPermissionKey(key)
	expect = errs.E(errs.InvalidJSON, key+" is not a valid key for class level permissions")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
}

func Test_buildMergedSchemaObject(t *testing.T) {
	var existingFields types.M
	var putRequest types.M
	var result types.M
	var expect types.M
	/************************************************************/
	existingFields = nil
	putRequest = nil
	result = buildMergedSchemaObject(existingFields, putRequest)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/************************************************************/
	existingFields = types.M{}
	putRequest = types.M{}
	result = buildMergedSchemaObject(existingFields, putRequest)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/************************************************************/
	existingFields = types.M{
		"_id":           "_User",
		"objectId":      types.M{"type": "String"},
		"createdAt":     types.M{"type": "Date"},
		"updatedAt":     types.M{"type": "Date"},
		"ACL":           types.M{"type": "ACL"},
		"username":      types.M{"type": "String"},
		"password":      types.M{"type": "String"},
		"email":         types.M{"type": "String"},
		"emailVerified": types.M{"type": "Boolean"},
		"name":          types.M{"type": "String"},
		"skill":         types.M{"type": "Array"},
	}
	putRequest = types.M{
		"age":   types.M{"type": "Number"},
		"skill": types.M{"__op": "Delete"},
	}
	result = buildMergedSchemaObject(existingFields, putRequest)
	expect = types.M{
		"name": types.M{"type": "String"},
		"age":  types.M{"type": "Number"},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/************************************************************/
	existingFields = types.M{
		"_id":           "user",
		"objectId":      types.M{"type": "String"},
		"createdAt":     types.M{"type": "Date"},
		"updatedAt":     types.M{"type": "Date"},
		"ACL":           types.M{"type": "ACL"},
		"username":      types.M{"type": "String"},
		"password":      types.M{"type": "String"},
		"email":         types.M{"type": "String"},
		"emailVerified": types.M{"type": "Boolean"},
		"name":          types.M{"type": "String"},
		"skill":         types.M{"type": "Array"},
	}
	putRequest = types.M{
		"age":   types.M{"type": "Number"},
		"skill": types.M{"__op": "Delete"},
	}
	result = buildMergedSchemaObject(existingFields, putRequest)
	expect = types.M{
		"username":      types.M{"type": "String"},
		"password":      types.M{"type": "String"},
		"email":         types.M{"type": "String"},
		"emailVerified": types.M{"type": "Boolean"},
		"name":          types.M{"type": "String"},
		"age":           types.M{"type": "Number"},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
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
	var schema types.M
	var result types.M
	var expect types.M
	/************************************************************/
	schema = nil
	result = convertSchemaToAdapterSchema(schema)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/************************************************************/
	schema = types.M{
		"className": "user",
	}
	result = convertSchemaToAdapterSchema(schema)
	expect = types.M{
		"className": "user",
		"fields": types.M{
			"objectId":  types.M{"type": "String"},
			"createdAt": types.M{"type": "Date"},
			"updatedAt": types.M{"type": "Date"},
			"_rperm":    types.M{"type": "Array"},
			"_wperm":    types.M{"type": "Array"},
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
	result = convertSchemaToAdapterSchema(schema)
	expect = types.M{
		"className": "_User",
		"fields": types.M{
			"objectId":         types.M{"type": "String"},
			"createdAt":        types.M{"type": "Date"},
			"updatedAt":        types.M{"type": "Date"},
			"key":              types.M{"type": "String"},
			"username":         types.M{"type": "String"},
			"_hashed_password": types.M{"type": "String"},
			"email":            types.M{"type": "String"},
			"emailVerified":    types.M{"type": "Boolean"},
			"_rperm":           types.M{"type": "Array"},
			"_wperm":           types.M{"type": "Array"},
		},
		"classLevelPermissions": nil,
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
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
	// TODO
}

func getSchema() *Schema {
	return &Schema{
		dbAdapter: getAdapter(),
	}
}

func getAdapter() *mongo.MongoAdapter {
	storage.TomatoDB = newMongoDB("192.168.99.100:27017/test")
	return mongo.NewMongoAdapter("tomato")
}

func newMongoDB(url string) *storage.Database {
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	database := session.DB("")
	db := &storage.Database{
		MongoSession:  session,
		MongoDatabase: database,
	}
	return db
}
