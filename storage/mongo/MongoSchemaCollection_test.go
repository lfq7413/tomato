package mongo

import (
	"reflect"
	"testing"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"

	"gopkg.in/mgo.v2"
)

func Test_getAllSchemas(t *testing.T) {
	// TODO
}

func Test_findSchema(t *testing.T) {
	// TODO
}

func Test_findAndDeleteSchema(t *testing.T) {
	// TODO
}

func Test_addSchema(t *testing.T) {
	db := openDB()
	defer db.Session.Close()
	msc := getSchemaCollection(db)
	var name string
	var fields types.M
	var classLevelPermissions types.M
	var result types.M
	var results []types.M
	var err error
	var expect types.M
	/*****************************************************/
	name = "user"
	fields = nil
	classLevelPermissions = nil
	result, err = msc.addSchema(name, fields, classLevelPermissions)
	expect = types.M{
		"className": "user",
		"fields": types.M{
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
	results, err = msc.collection.find(types.M{"_id": "user"}, types.M{})
	expect = types.M{
		"_id":       "user",
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
	}
	if err != nil || len(results) != 1 {
		t.Error("expect:", expect, "result:", results, err)
	}
	if len(results) == 1 && reflect.DeepEqual(expect, results[0]) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	msc.collection.drop()
	/*****************************************************/
	name = "user"
	fields = types.M{}
	classLevelPermissions = nil
	result, err = msc.addSchema(name, fields, classLevelPermissions)
	expect = types.M{
		"className": "user",
		"fields": types.M{
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
	results, err = msc.collection.find(types.M{"_id": "user"}, types.M{})
	expect = types.M{
		"_id":       "user",
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
	}
	if err != nil || len(results) != 1 {
		t.Error("expect:", expect, "result:", results, err)
	}
	if len(results) == 1 && reflect.DeepEqual(expect, results[0]) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	msc.collection.drop()
	/*****************************************************/
	name = "user"
	fields = types.M{
		"k1": types.M{
			"type":        "Pointer",
			"targetClass": "user",
		},
		"k2": types.M{
			"type":        "Relation",
			"targetClass": "user",
		},
		"k3": types.M{
			"type": "String",
		},
	}
	classLevelPermissions = nil
	result, err = msc.addSchema(name, fields, classLevelPermissions)
	expect = types.M{
		"className": "user",
		"fields": types.M{
			"k1": types.M{
				"type":        "Pointer",
				"targetClass": "user",
			},
			"k2": types.M{
				"type":        "Relation",
				"targetClass": "user",
			},
			"k3":        types.M{"type": "String"},
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
	results, err = msc.collection.find(types.M{"_id": "user"}, types.M{})
	expect = types.M{
		"_id":       "user",
		"k1":        "*user",
		"k2":        "relation<user>",
		"k3":        "string",
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
	}
	if err != nil || len(results) != 1 {
		t.Error("expect:", expect, "result:", results, err)
	}
	if len(results) == 1 && reflect.DeepEqual(expect, results[0]) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	msc.collection.drop()
	/*****************************************************/
	name = "user"
	fields = types.M{
		"k1": types.M{
			"type":        "Pointer",
			"targetClass": "user",
		},
		"k2": types.M{
			"type":        "Relation",
			"targetClass": "user",
		},
		"k3": types.M{
			"type": "String",
		},
	}
	classLevelPermissions = types.M{
		"find":   types.M{"*": true},
		"get":    types.M{"*": true},
		"create": types.M{"*": true},
		"update": types.M{"*": true},
	}
	result, err = msc.addSchema(name, fields, classLevelPermissions)
	expect = types.M{
		"className": "user",
		"fields": types.M{
			"k1": types.M{
				"type":        "Pointer",
				"targetClass": "user",
			},
			"k2": types.M{
				"type":        "Relation",
				"targetClass": "user",
			},
			"k3":        types.M{"type": "String"},
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
			"delete":   types.M{},
			"addField": types.M{},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	results, err = msc.collection.find(types.M{"_id": "user"}, types.M{})
	expect = types.M{
		"_id":       "user",
		"k1":        "*user",
		"k2":        "relation<user>",
		"k3":        "string",
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
		"_metadata": types.M{
			"class_permissions": types.M{
				"find":   types.M{"*": true},
				"get":    types.M{"*": true},
				"create": types.M{"*": true},
				"update": types.M{"*": true},
			},
		},
	}
	if err != nil || len(results) != 1 {
		t.Error("expect:", expect, "result:", results, err)
	}
	if len(results) == 1 && reflect.DeepEqual(expect, results[0]) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	msc.collection.drop()
	/*****************************************************/
	name = "user"
	fields = types.M{}
	classLevelPermissions = nil
	result, err = msc.addSchema(name, fields, classLevelPermissions)
	result, err = msc.addSchema(name, fields, classLevelPermissions)
	expectErr := errs.E(errs.DuplicateValue, "Class already exists.")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	msc.collection.drop()
}

func Test_updateSchema(t *testing.T) {
	db := openDB()
	defer db.Session.Close()
	msc := getSchemaCollection(db)
	var name string
	var fields types.M
	var classLevelPermissions types.M
	var results []types.M
	var update types.M
	var err error
	var expect types.M
	/*****************************************************/
	name = "user"
	fields = types.M{
		"key": types.M{
			"type": "String",
		},
	}
	classLevelPermissions = nil
	msc.addSchema(name, fields, classLevelPermissions)
	update = types.M{
		"$set": types.M{
			"key1": "string",
		},
	}
	err = msc.updateSchema(name, update)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	results, err = msc.collection.find(types.M{"_id": "user"}, types.M{})
	expect = types.M{
		"_id":       "user",
		"key":       "string",
		"key1":      "string",
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
	}
	if err != nil || len(results) != 1 {
		t.Error("expect:", expect, "result:", results, err)
	}
	if len(results) == 1 && reflect.DeepEqual(expect, results[0]) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	msc.collection.drop()
	/*****************************************************/
	name = "user"
	fields = types.M{
		"key": types.M{
			"type": "String",
		},
	}
	classLevelPermissions = nil
	msc.addSchema(name, fields, classLevelPermissions)
	update = types.M{
		"$unset": types.M{
			"key": nil,
		},
	}
	err = msc.updateSchema(name, update)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	results, err = msc.collection.find(types.M{"_id": "user"}, types.M{})
	expect = types.M{
		"_id":       "user",
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
	}
	if err != nil || len(results) != 1 {
		t.Error("expect:", expect, "result:", results, err)
	}
	if len(results) == 1 && reflect.DeepEqual(expect, results[0]) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	msc.collection.drop()
	/*****************************************************/
	name = "user"
	fields = types.M{
		"key": types.M{
			"type": "String",
		},
	}
	classLevelPermissions = types.M{
		"find": types.M{"*": true},
	}
	msc.addSchema(name, fields, classLevelPermissions)
	update = types.M{
		"$set": types.M{
			"_metadata": types.M{
				"class_permissions": types.M{
					"get": types.M{"*": true},
				},
			},
		},
	}
	err = msc.updateSchema(name, update)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	results, err = msc.collection.find(types.M{"_id": "user"}, types.M{})
	expect = types.M{
		"_id":       "user",
		"key":       "string",
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
		"_metadata": types.M{
			"class_permissions": types.M{
				"get": types.M{"*": true},
			},
		},
	}
	if err != nil || len(results) != 1 {
		t.Error("expect:", expect, "result:", results, err)
	}
	if len(results) == 1 && reflect.DeepEqual(expect, results[0]) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	msc.collection.drop()
	/*****************************************************/
	name = "user"
	fields = types.M{
		"key": types.M{
			"type": "String",
		},
	}
	classLevelPermissions = nil
	msc.addSchema(name, fields, classLevelPermissions)
	name = "user2"
	update = types.M{
		"$set": types.M{
			"key1": "string",
		},
	}
	err = msc.updateSchema(name, update)
	if err == nil || err.Error() != "not found" {
		t.Error("expect:", "not found", "result:", err)
	}
	msc.collection.drop()
}

func Test_upsertSchema(t *testing.T) {
	// TODO
}

func Test_addFieldIfNotExists(t *testing.T) {
	// findSchema
	// upsertSchema
	// TODO
}

func Test_mongoSchemaQueryFromNameQuery(t *testing.T) {
	// 测试用例与 mongoSchemaObjectFromNameFields 相同
}

func Test_mongoSchemaObjectFromNameFields(t *testing.T) {
	var name string
	var fields types.M
	var result types.M
	var expect types.M
	/*****************************************************/
	name = "user"
	fields = nil
	result = mongoSchemaObjectFromNameFields(name, fields)
	expect = types.M{
		"_id": name,
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	name = "user"
	fields = types.M{
		"key":  "string",
		"key1": "*_User",
	}
	result = mongoSchemaObjectFromNameFields(name, fields)
	expect = types.M{
		"_id":  name,
		"key":  "string",
		"key1": "*_User",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_mongoFieldToParseSchemaField(t *testing.T) {
	var ty string
	var result types.M
	var expect types.M
	/*****************************************************/
	ty = ""
	result = mongoFieldToParseSchemaField(ty)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	ty = "*user"
	result = mongoFieldToParseSchemaField(ty)
	expect = types.M{
		"type":        "Pointer",
		"targetClass": "user",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	ty = "relation<user>"
	result = mongoFieldToParseSchemaField(ty)
	expect = types.M{
		"type":        "Relation",
		"targetClass": "user",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	ty = "number"
	result = mongoFieldToParseSchemaField(ty)
	expect = types.M{
		"type": "Number",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	ty = "string"
	result = mongoFieldToParseSchemaField(ty)
	expect = types.M{
		"type": "String",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	ty = "boolean"
	result = mongoFieldToParseSchemaField(ty)
	expect = types.M{
		"type": "Boolean",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	ty = "date"
	result = mongoFieldToParseSchemaField(ty)
	expect = types.M{
		"type": "Date",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	ty = "map"
	result = mongoFieldToParseSchemaField(ty)
	expect = types.M{
		"type": "Object",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	ty = "object"
	result = mongoFieldToParseSchemaField(ty)
	expect = types.M{
		"type": "Object",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	ty = "array"
	result = mongoFieldToParseSchemaField(ty)
	expect = types.M{
		"type": "Array",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	ty = "geopoint"
	result = mongoFieldToParseSchemaField(ty)
	expect = types.M{
		"type": "GeoPoint",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	ty = "file"
	result = mongoFieldToParseSchemaField(ty)
	expect = types.M{
		"type": "File",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	ty = "bytes"
	result = mongoFieldToParseSchemaField(ty)
	expect = types.M{
		"type": "Bytes",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	ty = "other"
	result = mongoFieldToParseSchemaField(ty)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_mongoSchemaFieldsToParseSchemaFields(t *testing.T) {
	var schema types.M
	var result types.M
	var expect types.M
	/*****************************************************/
	schema = nil
	result = mongoSchemaFieldsToParseSchemaFields(schema)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	schema = types.M{}
	result = mongoSchemaFieldsToParseSchemaFields(schema)
	expect = types.M{
		"ACL":       types.M{"type": "ACL"},
		"createdAt": types.M{"type": "Date"},
		"updatedAt": types.M{"type": "Date"},
		"objectId":  types.M{"type": "String"},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	schema = types.M{
		"_id":                 "string",
		"_metadata":           "object",
		"_client_permissions": "object",
	}
	result = mongoSchemaFieldsToParseSchemaFields(schema)
	expect = types.M{
		"ACL":       types.M{"type": "ACL"},
		"createdAt": types.M{"type": "Date"},
		"updatedAt": types.M{"type": "Date"},
		"objectId":  types.M{"type": "String"},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	schema = types.M{
		"_id":                 "string",
		"_metadata":           "object",
		"_client_permissions": "object",
		"key1":                "*user",
		"key2":                "relation<user>",
		"key3":                "string",
	}
	result = mongoSchemaFieldsToParseSchemaFields(schema)
	expect = types.M{
		"key1": types.M{
			"type":        "Pointer",
			"targetClass": "user",
		},
		"key2": types.M{
			"type":        "Relation",
			"targetClass": "user",
		},
		"key3": types.M{
			"type": "String",
		},
		"ACL":       types.M{"type": "ACL"},
		"createdAt": types.M{"type": "Date"},
		"updatedAt": types.M{"type": "Date"},
		"objectId":  types.M{"type": "String"},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_mongoSchemaToParseSchema(t *testing.T) {
	var schema types.M
	var result types.M
	var expect types.M
	/*****************************************************/
	schema = nil
	result = mongoSchemaToParseSchema(schema)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	schema = types.M{
		"_id": "user",
	}
	result = mongoSchemaToParseSchema(schema)
	expect = types.M{
		"className": "user",
		"fields": types.M{
			"ACL":       types.M{"type": "ACL"},
			"createdAt": types.M{"type": "Date"},
			"updatedAt": types.M{"type": "Date"},
			"objectId":  types.M{"type": "String"},
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
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	schema = types.M{
		"_id": "user",
		"_metadata": types.M{
			"class_permissions": types.M{
				"find":   types.M{"*": true},
				"get":    types.M{"*": true},
				"create": types.M{"*": true},
				"update": types.M{"*": true},
			},
		},
	}
	result = mongoSchemaToParseSchema(schema)
	expect = types.M{
		"className": "user",
		"fields": types.M{
			"ACL":       types.M{"type": "ACL"},
			"createdAt": types.M{"type": "Date"},
			"updatedAt": types.M{"type": "Date"},
			"objectId":  types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"find":     types.M{"*": true},
			"get":      types.M{"*": true},
			"create":   types.M{"*": true},
			"update":   types.M{"*": true},
			"delete":   types.M{},
			"addField": types.M{},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	schema = types.M{
		"_id":  "user",
		"key1": "*user",
		"key2": "relation<user>",
		"key3": "string",
		"_metadata": types.M{
			"class_permissions": types.M{
				"find":   types.M{"*": true},
				"get":    types.M{"*": true},
				"create": types.M{"*": true},
				"update": types.M{"*": true},
			},
		},
	}
	result = mongoSchemaToParseSchema(schema)
	expect = types.M{
		"className": "user",
		"fields": types.M{
			"key1": types.M{
				"type":        "Pointer",
				"targetClass": "user",
			},
			"key2": types.M{
				"type":        "Relation",
				"targetClass": "user",
			},
			"key3": types.M{
				"type": "String",
			},
			"ACL":       types.M{"type": "ACL"},
			"createdAt": types.M{"type": "Date"},
			"updatedAt": types.M{"type": "Date"},
			"objectId":  types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"find":     types.M{"*": true},
			"get":      types.M{"*": true},
			"create":   types.M{"*": true},
			"update":   types.M{"*": true},
			"delete":   types.M{},
			"addField": types.M{},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_parseFieldTypeToMongoFieldType(t *testing.T) {
	var fieldType types.M
	var result string
	var expect string
	/*****************************************************/
	fieldType = nil
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = ""
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = ""
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type": "Pointer",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "*"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type": "Relation",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "relation<>"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type":        "Pointer",
		"targetClass": "user",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "*user"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type":        "Relation",
		"targetClass": "user",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "relation<user>"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type": "Number",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "number"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type": "String",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "string"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type": "Boolean",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "boolean"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type": "Date",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "date"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type": "Object",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "object"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type": "Array",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "array"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type": "GeoPoint",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "geopoint"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type": "File",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "file"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type": "Other",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = ""
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_mongoSchemaFromFieldsAndClassNameAndCLP(t *testing.T) {
	var fields types.M
	var className string
	var classLevelPermissions types.M
	var result types.M
	var expect types.M
	/*****************************************************/
	fields = nil
	className = "user"
	classLevelPermissions = nil
	result = mongoSchemaFromFieldsAndClassNameAndCLP(fields, className, classLevelPermissions)
	expect = types.M{
		"_id":       className,
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fields = types.M{
		"key1": types.M{
			"type":        "Pointer",
			"targetClass": "_User",
		},
		"key2": types.M{
			"type":        "Relation",
			"targetClass": "_User",
		},
		"loc": types.M{
			"type": "GeoPoint",
		},
	}
	className = "user"
	classLevelPermissions = nil
	result = mongoSchemaFromFieldsAndClassNameAndCLP(fields, className, classLevelPermissions)
	expect = types.M{
		"_id":       className,
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
		"key1":      "*_User",
		"key2":      "relation<_User>",
		"loc":       "geopoint",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fields = types.M{
		"key1": types.M{
			"type":        "Pointer",
			"targetClass": "_User",
		},
		"key2": types.M{
			"type":        "Relation",
			"targetClass": "_User",
		},
		"loc": types.M{
			"type": "GeoPoint",
		},
	}
	className = "user"
	classLevelPermissions = types.M{
		"find":     types.M{"*": true},
		"get":      types.M{"*": true},
		"create":   types.M{"*": true},
		"update":   types.M{"*": true},
		"delete":   types.M{"*": true},
		"addField": types.M{"*": true},
	}
	result = mongoSchemaFromFieldsAndClassNameAndCLP(fields, className, classLevelPermissions)
	expect = types.M{
		"_id":       className,
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
		"key1":      "*_User",
		"key2":      "relation<_User>",
		"loc":       "geopoint",
		"_metadata": types.M{
			"class_permissions": types.M{
				"find":     types.M{"*": true},
				"get":      types.M{"*": true},
				"create":   types.M{"*": true},
				"update":   types.M{"*": true},
				"delete":   types.M{"*": true},
				"addField": types.M{"*": true},
			},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func getSchemaCollection(db *mgo.Database) *MongoSchemaCollection {
	mc := newMongoCollection(db.C("SCHEMA"))
	msc := newMongoSchemaCollection(mc)
	return msc
}
