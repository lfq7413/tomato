package orm

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/storage"
	"github.com/lfq7413/tomato/storage/postgres"
	"github.com/lfq7413/tomato/test"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

func TestPostgres_AddClassIfNotExists(t *testing.T) {
	adapter := getPostgresAdapter()
	schama := getPostgresSchema()
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

func TestPostgres_UpdateClass(t *testing.T) {
	adapter := getPostgresAdapter()
	schama := getPostgresSchema()
	var class types.M
	var className string
	var submittedFields types.M
	var classLevelPermissions types.M
	var result types.M
	var err error
	var expect interface{}
	/************************************************************/
	className = "user"
	submittedFields = nil
	classLevelPermissions = nil
	result, err = schama.UpdateClass(className, submittedFields, classLevelPermissions)
	expect = errs.E(errs.InvalidClassName, "Class "+className+" does not exist.")
	if err == nil || reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	schama.data = nil
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	adapter.CreateClass(className, class)
	className = "user"
	submittedFields = types.M{
		"key": types.M{"type": "String"},
	}
	classLevelPermissions = nil
	result, err = schama.UpdateClass(className, submittedFields, classLevelPermissions)
	expect = errs.E(errs.ClassNotEmpty, "Field key exists, cannot update.")
	if err == nil || reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	schama.data = nil
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	adapter.CreateClass(className, class)
	className = "user"
	submittedFields = types.M{
		"key1": types.M{"__op": "Delete"},
	}
	classLevelPermissions = nil
	result, err = schama.UpdateClass(className, submittedFields, classLevelPermissions)
	expect = errs.E(errs.ClassNotEmpty, "Field key1 does not exist, cannot delete.")
	if err == nil || reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	schama.data = nil
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	adapter.CreateClass(className, class)
	className = "user"
	submittedFields = types.M{
		"key1": types.M{"type": "String"},
	}
	classLevelPermissions = nil
	result, err = schama.UpdateClass(className, submittedFields, classLevelPermissions)
	expect = types.M{
		"className": className,
		"fields": types.M{
			"key":       map[string]interface{}{"type": "String"},
			"key1":      map[string]interface{}{"type": "String"},
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
	schama.data = nil
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	adapter.CreateClass(className, class)
	className = "user"
	submittedFields = types.M{
		"key1": types.M{"type": "String"},
		"key":  types.M{"__op": "Delete"},
	}
	classLevelPermissions = nil
	result, err = schama.UpdateClass(className, submittedFields, classLevelPermissions)
	expect = types.M{
		"className": className,
		"fields": types.M{
			"key1":      map[string]interface{}{"type": "String"},
			"objectId":  types.M{"type": "String"},
			"updatedAt": types.M{"type": "Date"},
			"createdAt": types.M{"type": "Date"},
			"ACL":       types.M{"type": "ACL"},
		},
		"classLevelPermissions": types.M{
			"find":     map[string]interface{}{"*": true},
			"get":      map[string]interface{}{"*": true},
			"create":   map[string]interface{}{"*": true},
			"update":   map[string]interface{}{"*": true},
			"delete":   map[string]interface{}{"*": true},
			"addField": map[string]interface{}{"*": true},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	schama.data = nil
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	adapter.CreateClass(className, class)
	className = "user"
	submittedFields = types.M{
		"key1": types.M{"type": "String"},
		"key":  types.M{"__op": "Delete"},
	}
	classLevelPermissions = types.M{
		"find":   types.M{"*": true},
		"get":    types.M{"*": true},
		"create": types.M{"*": true},
		"update": types.M{"*": true},
		"delete": types.M{"*": true},
	}
	result, err = schama.UpdateClass(className, submittedFields, classLevelPermissions)
	expect = types.M{
		"className": className,
		"fields": types.M{
			"key1":      map[string]interface{}{"type": "String"},
			"objectId":  types.M{"type": "String"},
			"updatedAt": types.M{"type": "Date"},
			"createdAt": types.M{"type": "Date"},
			"ACL":       types.M{"type": "ACL"},
		},
		"classLevelPermissions": types.M{
			"find":     map[string]interface{}{"*": true},
			"get":      map[string]interface{}{"*": true},
			"create":   map[string]interface{}{"*": true},
			"update":   map[string]interface{}{"*": true},
			"delete":   map[string]interface{}{"*": true},
			"addField": types.M{"*": true},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	schama.data = nil
	adapter.DeleteAllClasses()
}

func TestPostgres_deleteField(t *testing.T) {
	adapter := getPostgresAdapter()
	schama := getPostgresSchema()
	var classSchama types.M
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
	classSchama = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass(className, classSchama)
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
	classSchama = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"key1":     types.M{"type": "String"},
		},
	}
	adapter.CreateClass(className, classSchama)
	class = types.M{
		"objectId": "1024",
		"key":      "hello",
		"key1":     "world",
	}
	adapter.CreateObject(className, classSchama, class)
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
		"key1":      map[string]interface{}{"type": "String"},
		"objectId":  map[string]interface{}{"type": "String"},
		"updatedAt": map[string]interface{}{"type": "Date"},
		"createdAt": map[string]interface{}{"type": "Date"},
		"ACL":       map[string]interface{}{"type": "ACL"},
	}
	if reflect.DeepEqual(class, r1["fields"]) == false {
		t.Error("expect:", class, "result:", r1["fields"])
	}
	// 检查数据
	r2, _ = adapter.Find(className, classSchama, types.M{}, types.M{})
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
	classSchama = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "Relation", "targetClass": "user"},
			"key1":     types.M{"type": "String"},
		},
	}
	adapter.CreateClass(className, classSchama)
	className = "abc"
	class = types.M{
		"objectId": "1024",
		"key1":     "world",
	}
	adapter.CreateObject(className, classSchama, class)
	className = "_Join:key:abc"
	classSchama = types.M{
		"fields": types.M{
			"relatedId": types.M{"type": "String"},
			"owningId":  types.M{"type": "String"},
		},
	}
	adapter.CreateClass(className, classSchama)
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
		"key1":      map[string]interface{}{"type": "String"},
		"objectId":  map[string]interface{}{"type": "String"},
		"updatedAt": map[string]interface{}{"type": "Date"},
		"createdAt": map[string]interface{}{"type": "Date"},
		"ACL":       map[string]interface{}{"type": "ACL"},
	}
	if reflect.DeepEqual(class, r1["fields"]) == false {
		t.Error("expect:", class, "result:", r1["fields"])
	}
	// 检查 schema
	className = "_Join:key:abc"
	r1, _ = adapter.GetClass(className)
	class = types.M{}
	if reflect.DeepEqual(class, r1) == false {
		t.Error("expect:", class, "result:", r1)
	}
	// 检查数据
	classSchama = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key1":     types.M{"type": "String"},
		},
	}
	className = "abc"
	r2, _ = adapter.Find(className, classSchama, types.M{}, types.M{})
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

func TestPostgres_validateObject(t *testing.T) {
	adapter := getPostgresAdapter()
	schama := getPostgresSchema()
	var className string
	var object types.M
	var query types.M
	var err error
	var expect error
	/************************************************************/
	className = "user"
	object = types.M{
		"key": "hello",
	}
	query = types.M{}
	err = schama.validateObject(className, object, query)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	schama.data = nil
	adapter.DeleteAllClasses()
	/************************************************************/
	className = "user"
	object = types.M{
		"key": time.Now(),
	}
	query = types.M{}
	err = schama.validateObject(className, object, query)
	expect = errs.E(errs.IncorrectType, "bad obj. can not get type")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	schama.data = nil
	adapter.DeleteAllClasses()
	/************************************************************/
	className = "user"
	object = types.M{
		"key": types.M{
			"__type":    "GeoPoint",
			"latitude":  20,
			"longitude": 20,
		},
		"key1": types.M{
			"__type":    "GeoPoint",
			"latitude":  20,
			"longitude": 20,
		},
	}
	query = types.M{}
	err = schama.validateObject(className, object, query)
	expect = errs.E(errs.IncorrectType, "there can only be one geopoint field in a class")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	schama.data = nil
	adapter.DeleteAllClasses()
}

func TestPostgres_testBaseCLP(t *testing.T) {
	schama := getPostgresSchema()
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

func TestPostgres_validatePermission(t *testing.T) {
	schama := getPostgresSchema()
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
	expect = errs.E(errs.OperationForbidden, "Permission denied for action create on class post.")
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
	expect = errs.E(errs.OperationForbidden, "Permission denied for action get on class post.")
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
	/************************************************************/
	schama.perms = types.M{
		"post": types.M{
			"get": types.M{"requiresAuthentication": true},
		},
	}
	className = "post"
	aclGroup = []string{}
	operation = "get"
	err = schama.validatePermission(className, aclGroup, operation)
	expect = errs.E(errs.ObjectNotFound, "Permission denied, user needs to be authenticated.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	schama.perms = types.M{
		"post": types.M{
			"get": types.M{"requiresAuthentication": true},
		},
	}
	className = "post"
	aclGroup = []string{"*"}
	operation = "get"
	err = schama.validatePermission(className, aclGroup, operation)
	expect = errs.E(errs.ObjectNotFound, "Permission denied, user needs to be authenticated.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/************************************************************/
	schama.perms = types.M{
		"post": types.M{
			"get": types.M{"requiresAuthentication": true},
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
	/************************************************************/
	schama.perms = types.M{
		"post": types.M{
			"get":            types.M{"role:1024": true, "requiresAuthentication": true},
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

func TestPostgres_EnforceClassExists(t *testing.T) {
	adapter := getPostgresAdapter()
	schama := getPostgresSchema()
	var class types.M
	var className string
	var err error
	/************************************************************/
	className = "post"
	err = schama.EnforceClassExists(className)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	className = "user"
	adapter.CreateClass(className, class)
	className = "user"
	err = schama.EnforceClassExists(className)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	className = "skill"
	adapter.CreateClass(className, class)
	schama.reloadData(nil)
	className = "skill"
	err = schama.EnforceClassExists(className)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	adapter.DeleteAllClasses()
}

func TestPostgres_validateNewClass(t *testing.T) {
	adapter := getPostgresAdapter()
	schama := getPostgresSchema()
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
	schama.reloadData(nil)
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
	schama.reloadData(nil)
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
	schama.reloadData(nil)
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

func TestPostgres_validateSchemaData(t *testing.T) {
	schama := getPostgresSchema()
	var className string
	var fields types.M
	var classLevelPermissions types.M
	var existingFieldNames map[string]bool
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
	existingFieldNames = map[string]bool{"key": true}
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
	existingFieldNames = map[string]bool{"key": true}
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
	existingFieldNames = map[string]bool{"key": true}
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
	existingFieldNames = map[string]bool{"key": true}
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
	existingFieldNames = map[string]bool{"key": true}
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
	existingFieldNames = map[string]bool{"key": true}
	err = schama.validateSchemaData(className, fields, classLevelPermissions, existingFieldNames)
	expect = errs.E(errs.IncorrectType, "currently, only one GeoPoint field may exist in an object. Adding loc when loc2 already exists.")
	expect2 := errs.E(errs.IncorrectType, "currently, only one GeoPoint field may exist in an object. Adding loc2 when loc already exists.")
	if reflect.DeepEqual(expect, err) == false && reflect.DeepEqual(expect2, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
}

func TestPostgres_validateRequiredColumns(t *testing.T) {
	schama := getPostgresSchema()
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

func TestPostgres_enforceFieldExists(t *testing.T) {
	adapter := getPostgresAdapter()
	schama := getPostgresSchema()
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
	err = schama.enforceFieldExists(className, fieldName, fieldtype)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	schama.reloadDataPromise = nil
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
	err = schama.enforceFieldExists(className, fieldName, fieldtype)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	schama.reloadDataPromise = nil
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
	err = schama.enforceFieldExists(className, fieldName, fieldtype)
	expect = errs.E(errs.InvalidKeyName, "Invalid field name: "+fieldName)
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	schama.reloadDataPromise = nil
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
	err = schama.enforceFieldExists(className, fieldName, fieldtype)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	schama.reloadDataPromise = nil
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
	err = schama.enforceFieldExists(className, fieldName, fieldtype)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	schama.reloadDataPromise = nil
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
	err = schama.enforceFieldExists(className, fieldName, fieldtype)
	expect = errs.E(errs.IncorrectType, "schema mismatch for post.key1; expected String but got Number")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	schama.reloadDataPromise = nil
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
	err = schama.enforceFieldExists(className, fieldName, fieldtype)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	schama.reloadDataPromise = nil
	adapter.DeleteAllClasses()
}

func TestPostgres_setPermissions(t *testing.T) {
	adapter := getPostgresAdapter()
	schama := getPostgresSchema()
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
	expect = error(nil)
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
		"get":      map[string]interface{}{"*": true},
		"create":   types.M{"*": true},
		"find":     types.M{"*": true},
		"update":   types.M{"*": true},
		"delete":   types.M{"*": true},
		"addField": types.M{"*": true},
	}
	if reflect.DeepEqual(expect, schama.perms[className]) == false {
		t.Error("expect:", expect, "result:", schama.perms[className])
	}
	adapter.DeleteAllClasses()
}

func TestPostgres_HasClass(t *testing.T) {
	adapter := getPostgresAdapter()
	schama := getPostgresSchema()
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

func TestPostgres_getExpectedType(t *testing.T) {
	adapter := getPostgresAdapter()
	schama := getPostgresSchema()
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
	schama.reloadData(nil)
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
	schama.reloadData(nil)
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
	schama.reloadData(nil)
	className = "post"
	fieldName = "key1"
	result = schama.getExpectedType(className, fieldName)
	expect = types.M{"type": "String"}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	adapter.DeleteAllClasses()
}

func TestPostgres_reloadData(t *testing.T) {
	adapter := getPostgresAdapter()
	schama := getPostgresSchema()
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
	schama.reloadData(nil)
	expect = types.M{
		"post": types.M{
			"key1":      map[string]interface{}{"type": "String"},
			"objectId":  types.M{"type": "String"},
			"updatedAt": types.M{"type": "Date"},
			"createdAt": types.M{"type": "Date"},
			"ACL":       types.M{"type": "ACL"},
		},
		"_User": types.M{
			"key1":          map[string]interface{}{"type": "String"},
			"objectId":      types.M{"type": "String"},
			"updatedAt":     types.M{"type": "Date"},
			"createdAt":     types.M{"type": "Date"},
			"ACL":           types.M{"type": "ACL"},
			"username":      types.M{"type": "String"},
			"password":      types.M{"type": "String"},
			"email":         types.M{"type": "String"},
			"emailVerified": types.M{"type": "Boolean"},
			"authData":      types.M{"type": "Object"},
		},
		"_PushStatus": types.M{
			"objectId":      types.M{"type": "String"},
			"updatedAt":     types.M{"type": "Date"},
			"createdAt":     types.M{"type": "Date"},
			"ACL":           types.M{"type": "ACL"},
			"pushTime":      types.M{"type": "String"},
			"source":        types.M{"type": "String"},
			"query":         types.M{"type": "String"},
			"payload":       types.M{"type": "String"},
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
		"_JobStatus": types.M{
			"objectId":   types.M{"type": "String"},
			"updatedAt":  types.M{"type": "Date"},
			"createdAt":  types.M{"type": "Date"},
			"ACL":        types.M{"type": "ACL"},
			"jobName":    types.M{"type": "String"},
			"source":     types.M{"type": "String"},
			"status":     types.M{"type": "String"},
			"message":    types.M{"type": "String"},
			"params":     types.M{"type": "Object"},
			"finishedAt": types.M{"type": "Date"},
		},
		"_Hooks": types.M{
			"objectId":     types.M{"type": "String"},
			"updatedAt":    types.M{"type": "Date"},
			"createdAt":    types.M{"type": "Date"},
			"ACL":          types.M{"type": "ACL"},
			"functionName": types.M{"type": "String"},
			"className":    types.M{"type": "String"},
			"triggerName":  types.M{"type": "String"},
			"url":          types.M{"type": "String"},
		},
		"_GlobalConfig": types.M{
			"objectId":  types.M{"type": "String"},
			"updatedAt": types.M{"type": "Date"},
			"createdAt": types.M{"type": "Date"},
			"ACL":       types.M{"type": "ACL"},
			"params":    types.M{"type": "Object"},
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
		"_PushStatus":   types.M{},
		"_JobStatus":    types.M{},
		"_Hooks":        types.M{},
		"_GlobalConfig": types.M{},
	}
	if reflect.DeepEqual(expect, schama.perms) == false {
		t.Error("expect:", expect, "result:", schama.perms)
	}
	adapter.DeleteAllClasses()
}

func TestPostgres_GetAllClasses(t *testing.T) {
	adapter := getPostgresAdapter()
	schama := getPostgresSchema()
	var class types.M
	var result []types.M
	var err error
	var expect []types.M
	/************************************************************/
	result, err = schama.GetAllClasses(types.M{"clearCache": true})
	expect = []types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	schama.reloadDataPromise = nil
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("post", class)
	result, err = schama.GetAllClasses(types.M{"clearCache": true})
	expect = []types.M{
		types.M{
			"className": "post",
			"fields": types.M{
				"key1":      map[string]interface{}{"type": "String"},
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
	schama.reloadDataPromise = nil
	adapter.DeleteAllClasses()
	/************************************************************/
	class = types.M{
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
	}
	adapter.CreateClass("_User", class)
	result, err = schama.GetAllClasses(types.M{"clearCache": true})
	expect = []types.M{
		types.M{
			"className": "_User",
			"fields": types.M{
				"key1":          map[string]interface{}{"type": "String"},
				"objectId":      types.M{"type": "String"},
				"updatedAt":     types.M{"type": "Date"},
				"createdAt":     types.M{"type": "Date"},
				"ACL":           types.M{"type": "ACL"},
				"username":      types.M{"type": "String"},
				"password":      types.M{"type": "String"},
				"email":         types.M{"type": "String"},
				"emailVerified": types.M{"type": "Boolean"},
				"authData":      types.M{"type": "Object"},
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
	schama.reloadDataPromise = nil
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
	result, err = schama.GetAllClasses(types.M{"clearCache": true})
	expect = []types.M{
		types.M{
			"className": "post",
			"fields": types.M{
				"key1":      map[string]interface{}{"type": "String"},
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
				"key1":      map[string]interface{}{"type": "String"},
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
	schama.reloadDataPromise = nil
	adapter.DeleteAllClasses()
}

func TestPostgres_GetOneSchema(t *testing.T) {
	adapter := getPostgresAdapter()
	schama := getPostgresSchema()
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
	result, err = schama.GetOneSchema(className, allowVolatileClasses, nil)
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
	result, err = schama.GetOneSchema(className, allowVolatileClasses, nil)
	expect = types.M{
		"className": "post",
		"fields": types.M{
			"key1":      map[string]interface{}{"type": "String"},
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
	result, err = schama.GetOneSchema(className, allowVolatileClasses, nil)
	expect = types.M{
		"className": "post",
		"fields": types.M{
			"key1":      map[string]interface{}{"type": "String"},
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
	result, err = schama.GetOneSchema(className, allowVolatileClasses, nil)
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
			"key1": types.M{"type": "String"},
		},
	}
	className = "_PushStatus"
	allowVolatileClasses = true
	result, err = schama.GetOneSchema(className, allowVolatileClasses, nil)
	expect = types.M{
		"className": "_PushStatus",
		"fields": types.M{
			"key1": types.M{"type": "String"},
		},
		"classLevelPermissions": nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
}

func getPostgresSchema() *Schema {
	return &Schema{
		dbAdapter: getPostgresAdapter(),
		cache:     getSchemaCache(),
	}
}

func getPostgresAdapter() storage.Adapter {
	return postgres.NewPostgresAdapter("tomato", test.OpenPostgreSQForTest())
}

func objectType(n string, x interface{}) {
	v := reflect.ValueOf(x)
	fmt.Println(n, "type:", v.Type())
	if m := utils.M(x); m != nil {
		for key, value := range m {
			objectType(key, value)
		}
	}
}
