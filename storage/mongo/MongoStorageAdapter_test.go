package mongo

import (
	"reflect"
	"testing"

	"gopkg.in/mgo.v2"

	"github.com/lfq7413/tomato/storage"
	"github.com/lfq7413/tomato/types"
)

func Test_ClassExists(t *testing.T) {
	// TODO
}

func Test_SetClassLevelPermissions(t *testing.T) {
	// TODO
}

func Test_CreateClass(t *testing.T) {
	// TODO
}

func Test_AddFieldIfNotExists(t *testing.T) {
	// TODO
}

func Test_DeleteClass(t *testing.T) {
	// TODO
}

func Test_DeleteAllClasses(t *testing.T) {
	// TODO
}

func Test_DeleteFields(t *testing.T) {
	// TODO
}

func Test_CreateObject(t *testing.T) {
	// TODO
}

func Test_GetClass(t *testing.T) {
	// TODO
}

func Test_GetAllClasses(t *testing.T) {
	// TODO
}

func Test_getCollectionNames(t *testing.T) {
	adapter := getAdapter()
	var names []string
	/*****************************************************************/
	names = adapter.getCollectionNames()
	if names != nil && len(names) > 0 {
		t.Error("expect:", 0, "result:", len(names))
	}
	/*****************************************************************/
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

	// TODO
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
	// TODO
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
