package mongo

import (
	"reflect"
	"testing"

	"gopkg.in/mgo.v2"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/storage"
	"github.com/lfq7413/tomato/types"
)

func Test_ClassExists(t *testing.T) {
	adapter := getAdapter()
	var name string
	var exist bool
	/*****************************************************/
	adapter.adaptiveCollection("user").insertOne(types.M{"_id": "01"})
	adapter.adaptiveCollection("user2").insertOne(types.M{"_id": "01"})
	name = "user"
	exist = adapter.ClassExists(name)
	if exist == false {
		t.Error("expect:", true, "result:", exist)
	}
	adapter.DeleteAllClasses()
	/*****************************************************/
	adapter.adaptiveCollection("user").insertOne(types.M{"_id": "01"})
	adapter.adaptiveCollection("user2").insertOne(types.M{"_id": "01"})
	name = "user3"
	exist = adapter.ClassExists(name)
	if exist == true {
		t.Error("expect:", false, "result:", exist)
	}
	adapter.DeleteAllClasses()
	/*****************************************************/
	adapter.adaptiveCollection("user").insertOne(types.M{"_id": "01"})
	adapter.adaptiveCollection("user2").insertOne(types.M{"_id": "01"})
	name = "user3"
	exist = adapter.ClassExists(name)
	if exist == true {
		t.Error("expect:", false, "result:", exist)
	}
	adapter.adaptiveCollection("user3").insertOne(types.M{"_id": "01"})
	name = "user3"
	exist = adapter.ClassExists(name)
	if exist == false {
		t.Error("expect:", true, "result:", exist)
	}
	adapter.DeleteAllClasses()
}

func Test_SetClassLevelPermissions(t *testing.T) {
	adapter := getAdapter()
	var className string
	var clps types.M
	var err error
	var result []types.M
	var expect types.M
	/*****************************************************/
	className = "user"
	clps = nil
	adapter.CreateClass(className, nil)
	err = adapter.SetClassLevelPermissions(className, clps)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	result, err = adapter.schemaCollection().collection.find(types.M{"_id": className}, types.M{})
	expect = types.M{
		"_id":       className,
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
		"_metadata": types.M{
			"class_permissions": types.M{},
		},
	}
	if err != nil || result == nil || reflect.DeepEqual(expect, result[0]) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
	/*****************************************************/
	className = "user"
	clps = types.M{
		"find":   types.M{"*": true},
		"get":    types.M{"*": true},
		"create": types.M{"*": true},
		"update": types.M{"*": true},
	}
	adapter.CreateClass(className, nil)
	err = adapter.SetClassLevelPermissions(className, clps)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	result, err = adapter.schemaCollection().collection.find(types.M{"_id": className}, types.M{})
	expect = types.M{
		"_id":       className,
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
	if err != nil || result == nil || reflect.DeepEqual(expect, result[0]) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
}

func Test_CreateClass(t *testing.T) {
	adapter := getAdapter()
	var className string
	var schema types.M
	var result types.M
	var err error
	var expect types.M
	var results []types.M
	/*****************************************************/
	className = "user"
	schema = nil
	result, err = adapter.CreateClass(className, schema)
	expect = types.M{
		"className": className,
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
	} else {
		results, err = adapter.schemaCollection().collection.find(types.M{"_id": className}, types.M{})
		expect = types.M{
			"_id":       className,
			"objectId":  "string",
			"updatedAt": "string",
			"createdAt": "string",
		}
		if results == nil || len(results) != 1 {
			t.Error("expect:", expect, "result:", results, err)
		} else {
			if reflect.DeepEqual(expect, results[0]) == false {
				t.Error("expect:", expect, "result:", results[0], err)
			}
		}
	}
	adapter.DeleteAllClasses()
	/*****************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{},
	}
	result, err = adapter.CreateClass(className, schema)
	expect = types.M{
		"className": className,
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
	} else {
		results, err = adapter.schemaCollection().collection.find(types.M{"_id": className}, types.M{})
		expect = types.M{
			"_id":       className,
			"objectId":  "string",
			"updatedAt": "string",
			"createdAt": "string",
		}
		if results == nil || len(results) != 1 {
			t.Error("expect:", expect, "result:", results, err)
		} else {
			if reflect.DeepEqual(expect, results[0]) == false {
				t.Error("expect:", expect, "result:", results[0], err)
			}
		}
	}
	adapter.DeleteAllClasses()
	/*****************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
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
		},
	}
	result, err = adapter.CreateClass(className, schema)
	expect = types.M{
		"className": className,
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
	} else {
		results, err = adapter.schemaCollection().collection.find(types.M{"_id": className}, types.M{})
		expect = types.M{
			"_id":       className,
			"k1":        "*user",
			"k2":        "relation<user>",
			"k3":        "string",
			"objectId":  "string",
			"updatedAt": "string",
			"createdAt": "string",
		}
		if results == nil || len(results) != 1 {
			t.Error("expect:", expect, "result:", results, err)
		} else {
			if reflect.DeepEqual(expect, results[0]) == false {
				t.Error("expect:", expect, "result:", results[0], err)
			}
		}
	}
	adapter.DeleteAllClasses()
	/*****************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
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
		},
		"classLevelPermissions": types.M{
			"find":   types.M{"*": true},
			"get":    types.M{"*": true},
			"create": types.M{"*": true},
			"update": types.M{"*": true},
		},
	}
	result, err = adapter.CreateClass(className, schema)
	expect = types.M{
		"className": className,
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
	} else {
		results, err = adapter.schemaCollection().collection.find(types.M{"_id": className}, types.M{})
		expect = types.M{
			"_id":       className,
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
		if results == nil || len(results) != 1 {
			t.Error("expect:", expect, "result:", results, err)
		} else {
			if reflect.DeepEqual(expect, results[0]) == false {
				t.Error("expect:", expect, "result:", results[0], err)
			}
		}
	}
	adapter.DeleteAllClasses()
	/*****************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{},
	}
	result, err = adapter.CreateClass(className, schema)
	result, err = adapter.CreateClass(className, schema)
	expectErr := errs.E(errs.DuplicateValue, "Class already exists.")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	adapter.DeleteAllClasses()
}

func Test_AddFieldIfNotExists(t *testing.T) {
	// 测试用例与 MongoSchemaCollection.addFieldIfNotExists 相同
}

func Test_DeleteClass(t *testing.T) {
	adapter := getAdapter()
	var className string
	var result types.M
	var err error
	var expect types.M
	/*****************************************************/
	className = "user"
	result, err = adapter.DeleteClass(className)
	expect = nil
	if err != nil || result != nil {
		t.Error("expect:", expect, "result:", result, err)
	}
	/*****************************************************/
	className = "user"
	adapter.adaptiveCollection(className).insertOne(types.M{"_id": "1024"})
	adapter.CreateClass(className, nil)
	result, err = adapter.DeleteClass(className)
	expect = types.M{
		"_id":       className,
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	} else {
		results, err := adapter.adaptiveCollection(className).find(types.M{}, types.M{})
		if results != nil && len(results) > 0 {
			t.Error("expect:", 0, "result:", results, err)
		}
		results, err = adapter.schemaCollection().collection.find(types.M{"_id": className}, types.M{})
		if results != nil && len(results) > 0 {
			t.Error("expect:", 0, "result:", results, err)
		}
	}
	adapter.DeleteAllClasses()
}

func Test_DeleteAllClasses(t *testing.T) {
	adapter := getAdapter()
	var err error
	var names []string
	/*****************************************************/
	err = adapter.DeleteAllClasses()
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	names = adapter.getCollectionNames()
	if names != nil && len(names) != 0 {
		t.Error("expect:", 0, "result:", len(names))
	}
	/*****************************************************/
	adapter.adaptiveCollection("user").insertOne(types.M{"_id": "01"})
	adapter.adaptiveCollection("user1").insertOne(types.M{"_id": "01"})
	err = adapter.DeleteAllClasses()
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	names = adapter.getCollectionNames()
	if names != nil && len(names) != 0 {
		t.Error("expect:", 0, "result:", len(names))
	}
}

func Test_DeleteFields(t *testing.T) {
	adapter := getAdapter()
	var className string
	var schema types.M
	var fieldNames []string
	var err error
	var results []types.M
	var expect types.M
	/*****************************************************/
	className = "user"
	schema = nil
	fieldNames = nil
	err = adapter.DeleteFields(className, schema, fieldNames)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	/*****************************************************/
	className = "user"
	schema = nil
	fieldNames = []string{}
	err = adapter.DeleteFields(className, schema, fieldNames)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	/*****************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	fieldNames = []string{"key"}
	adapter.adaptiveCollection(className).insertOne(types.M{"_id": "01", "key": "hello"})
	adapter.CreateClass(className, schema)
	err = adapter.DeleteFields(className, schema, fieldNames)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	results, err = adapter.adaptiveCollection(className).find(types.M{"_id": "01"}, types.M{})
	expect = types.M{
		"_id": "01",
	}
	if err != nil || results == nil || len(results) != 1 || reflect.DeepEqual(expect, results[0]) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	results, err = adapter.schemaCollection().collection.find(types.M{"_id": className}, types.M{})
	expect = types.M{
		"_id":       className,
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
	}
	if err != nil || results == nil || len(results) != 1 || reflect.DeepEqual(expect, results[0]) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	adapter.DeleteAllClasses()
	/*****************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "Pointer"},
		},
	}
	fieldNames = []string{"key"}
	adapter.adaptiveCollection(className).insertOne(types.M{"_id": "01", "_p_key": "hello"})
	adapter.CreateClass(className, schema)
	err = adapter.DeleteFields(className, schema, fieldNames)
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	results, err = adapter.adaptiveCollection(className).find(types.M{"_id": "01"}, types.M{})
	expect = types.M{
		"_id": "01",
	}
	if err != nil || results == nil || len(results) != 1 || reflect.DeepEqual(expect, results[0]) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	results, err = adapter.schemaCollection().collection.find(types.M{"_id": className}, types.M{})
	expect = types.M{
		"_id":       className,
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
	}
	if err != nil || results == nil || len(results) != 1 || reflect.DeepEqual(expect, results[0]) == false {
		t.Error("expect:", expect, "result:", results, err)
	}
	adapter.DeleteAllClasses()
}

func Test_CreateObject(t *testing.T) {
	// TODO
}

func Test_GetClass(t *testing.T) {
	adapter := getAdapter()
	var className string
	var result types.M
	var err error
	var expect types.M
	/*****************************************************/
	className = "user"
	result, err = adapter.GetClass(className)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
	/*****************************************************/
	className = "user"
	adapter.CreateClass(className, nil)
	result, err = adapter.GetClass(className)
	expect = types.M{
		"className": className,
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
	adapter.DeleteAllClasses()
}

func Test_GetAllClasses(t *testing.T) {
	adapter := getAdapter()
	var className string
	var result []types.M
	var err error
	var expect []types.M
	/*****************************************************/
	result, err = adapter.GetAllClasses()
	expect = []types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
	/*****************************************************/
	className = "user"
	adapter.CreateClass(className, nil)
	className = "user1"
	adapter.CreateClass(className, nil)
	result, err = adapter.GetAllClasses()
	expect = []types.M{
		types.M{
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
		},
		types.M{
			"className": "user1",
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
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	adapter.DeleteAllClasses()
}

func Test_getCollectionNames(t *testing.T) {
	adapter := getAdapter()
	var names []string
	/*****************************************************/
	names = adapter.getCollectionNames()
	if names != nil && len(names) > 0 {
		t.Error("expect:", 0, "result:", len(names))
	}
	/*****************************************************/
	adapter.adaptiveCollection("user").insertOne(types.M{"_id": "01"})
	adapter.adaptiveCollection("user1").insertOne(types.M{"_id": "01"})
	names = adapter.getCollectionNames()
	if names == nil || len(names) != 2 {
		t.Error("expect:", 2, "result:", len(names))
	} else {
		expect := []string{"tomatouser", "tomatouser1"}
		if reflect.DeepEqual(expect, names) == false {
			t.Error("expect:", expect, "result:", names)
		}
	}
	adapter.adaptiveCollection("user").drop()
	adapter.adaptiveCollection("user1").drop()
}

func Test_DeleteObjectsByQuery(t *testing.T) {
	// TODO
}

func Test_UpdateObjectsByQuery(t *testing.T) {
	// TODO
}

func Test_FindOneAndUpdate(t *testing.T) {
	// TODO
}

func Test_UpsertOneObject(t *testing.T) {
	// TODO
}

func Test_Find(t *testing.T) {
	// TODO
}

func Test_AdapterRawFind(t *testing.T) {
	// TODO
}

func Test_Count(t *testing.T) {
	// TODO
}

func Test_EnsureUniqueness(t *testing.T) {
	// TODO
}

func Test_storageAdapterAllCollections(t *testing.T) {
	adapter := getAdapter()
	var result []*MongoCollection
	var expect []*MongoCollection
	/*****************************************************/
	result = storageAdapterAllCollections(adapter)
	expect = []*MongoCollection{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	adapter.adaptiveCollection("user").insertOne(types.M{"_id": "01"})
	adapter.adaptiveCollection("user1").insertOne(types.M{"_id": "01"})
	result = storageAdapterAllCollections(adapter)
	if result == nil || len(result) != 2 {
		t.Error("expect:", 2, "result:", len(result))
	} else {
		expect := []string{"tomatouser", "tomatouser1"}
		names := []string{result[0].collection.Name, result[1].collection.Name}
		if reflect.DeepEqual(expect, names) == false {
			t.Error("expect:", expect, "result:", names)
		}
	}
	adapter.adaptiveCollection("user").drop()
	adapter.adaptiveCollection("user1").drop()
	/*****************************************************/
	adapter.adaptiveCollection("user").insertOne(types.M{"_id": "01"})
	adapter.adaptiveCollection("user1").insertOne(types.M{"_id": "01"})
	adapter.adaptiveCollection("user.system.id").insertOne(types.M{"_id": "01"})
	result = storageAdapterAllCollections(adapter)
	if result == nil || len(result) != 2 {
		t.Error("expect:", 2, "result:", len(result))
	} else {
		expect := []string{"tomatouser", "tomatouser1"}
		names := []string{result[0].collection.Name, result[1].collection.Name}
		if reflect.DeepEqual(expect, names) == false {
			t.Error("expect:", expect, "result:", names)
		}
	}
	adapter.adaptiveCollection("user").drop()
	adapter.adaptiveCollection("user1").drop()
	adapter.adaptiveCollection("ser.system.id").drop()
}

func Test_convertParseSchemaToMongoSchema(t *testing.T) {
	var schema types.M
	var result types.M
	var expect types.M
	/*****************************************************/
	schema = nil
	result = convertParseSchemaToMongoSchema(schema)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	schema = types.M{}
	result = convertParseSchemaToMongoSchema(schema)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	schema = types.M{
		"className": "user",
		"fields": types.M{
			"name":             types.M{"type": "string"},
			"_rperm":           types.M{"type": "array"},
			"_wperm":           types.M{"type": "array"},
			"_hashed_password": types.M{"type": "string"},
		},
	}
	result = convertParseSchemaToMongoSchema(schema)
	expect = types.M{
		"className": "user",
		"fields": types.M{
			"name":             types.M{"type": "string"},
			"_hashed_password": types.M{"type": "string"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	schema = types.M{
		"className": "_User",
		"fields": types.M{
			"name":             types.M{"type": "string"},
			"_rperm":           types.M{"type": "array"},
			"_wperm":           types.M{"type": "array"},
			"_hashed_password": types.M{"type": "string"},
		},
	}
	result = convertParseSchemaToMongoSchema(schema)
	expect = types.M{
		"className": "_User",
		"fields": types.M{
			"name": types.M{"type": "string"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
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

func getAdapter() *MongoAdapter {
	storage.TomatoDB = newMongoDB("192.168.99.100:27017/test")
	return NewMongoAdapter("tomato")
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
