package orm

import (
	"reflect"
	"testing"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

func Test_CollectionExists(t *testing.T) {
	initEnv()
	var object types.M
	var className string
	var result bool
	var expect bool
	/*************************************************/
	className = "user"
	result = TomatoDBController.CollectionExists(className)
	expect = false
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	className = "user"
	object = types.M{
		"key": "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	className = "user"
	result = TomatoDBController.CollectionExists(className)
	expect = true
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
}

func Test_PurgeCollection(t *testing.T) {
	// TODO
}

func Test_Find(t *testing.T) {
	// TODO
}

func Test_Destroy(t *testing.T) {
	// TODO
}

func Test_Update(t *testing.T) {
	// handleRelationUpdates
	// TODO
}

func Test_Create(t *testing.T) {
	// validateClassName
	// handleRelationUpdates
	// TODO
}

func Test_validateClassName(t *testing.T) {
	// TODO
}

func Test_handleRelationUpdates(t *testing.T) {
	// removeRelation
	// TODO
}

func Test_addRelation(t *testing.T) {
	initEnv()
	var object types.M
	var key, fromClassName, fromID, toID string
	var err error
	var expect error
	var results []types.M
	var expectRes types.M
	/*************************************************/
	key = "post"
	fromClassName = "user"
	fromID = "1001"
	toID = "2001"
	err = TomatoDBController.addRelation(key, fromClassName, fromID, toID)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	results, err = Adapter.Find("_Join:post:user", relationSchema, types.M{}, types.M{})
	expectRes = types.M{
		"relatedId": "2001",
		"owningId":  "1001",
	}
	delete(results[0], "objectId")
	if len(results) != 1 || reflect.DeepEqual(expectRes, results[0]) == false {
		t.Error("expect:", expectRes, "result:", results[0])
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"objectId":  "01",
		"relatedId": "2001",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:post:user", relationSchema, object)
	key = "post"
	fromClassName = "user"
	fromID = "1001"
	toID = "2001"
	err = TomatoDBController.addRelation(key, fromClassName, fromID, toID)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	results, err = Adapter.Find("_Join:post:user", relationSchema, types.M{}, types.M{})
	expectRes = types.M{
		"objectId":  "01",
		"relatedId": "2001",
		"owningId":  "1001",
	}
	if len(results) != 1 || reflect.DeepEqual(expectRes, results[0]) == false {
		t.Error("expect:", expectRes, "result:", results[0])
	}
	Adapter.DeleteAllClasses()
}

func Test_removeRelation(t *testing.T) {
	// TODO
}

func Test_ValidateObject(t *testing.T) {
	// TODO
}

func Test_LoadSchema(t *testing.T) {
	initEnv()
	var object types.M
	var className string
	var result *Schema
	var expect types.M
	/*************************************************/
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"get": types.M{"*": true},
		},
	}
	className = "user"
	Adapter.CreateClass(className, object)
	result = TomatoDBController.LoadSchema()
	expect = types.M{
		"key":       types.M{"type": "String"},
		"objectId":  types.M{"type": "String"},
		"createdAt": types.M{"type": "Date"},
		"updatedAt": types.M{"type": "Date"},
		"ACL":       types.M{"type": "ACL"},
	}
	if reflect.DeepEqual(expect, result.data["user"]) == false {
		t.Error("expect:", expect, "result:", result.data["user"])
	}
	expect = types.M{
		"get":      types.M{"*": true},
		"find":     types.M{},
		"create":   types.M{},
		"update":   types.M{},
		"delete":   types.M{},
		"addField": types.M{},
	}
	if reflect.DeepEqual(expect, result.perms["user"]) == false {
		t.Error("expect:", expect, "result:", result.perms["user"])
	}
	Adapter.DeleteAllClasses()
}

func Test_DeleteEverything(t *testing.T) {
	// TODO
}

func Test_RedirectClassNameForKey(t *testing.T) {
	initEnv()
	var object types.M
	var className string
	var key string
	var result string
	var expect string
	/*************************************************/
	className = "user"
	key = "name"
	result = TomatoDBController.RedirectClassNameForKey(className, key)
	expect = "user"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	object = types.M{
		"fields": types.M{
			"name": types.M{"type": "String"},
		},
	}
	className = "user"
	Adapter.CreateClass(className, object)
	className = "user"
	key = "name"
	result = TomatoDBController.RedirectClassNameForKey(className, key)
	expect = "user"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	object = types.M{
		"fields": types.M{
			"name": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
		},
	}
	className = "user"
	Adapter.CreateClass(className, object)
	className = "user"
	key = "name"
	result = TomatoDBController.RedirectClassNameForKey(className, key)
	expect = "post"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	TomatoDBController.DeleteEverything()
}

func Test_canAddField(t *testing.T) {
	initEnv()
	var schema *Schema
	var className string
	var object types.M
	var acl []string
	var err error
	var expect error
	/*************************************************/
	schema = nil
	className = "user"
	object = nil
	acl = nil
	err = TomatoDBController.canAddField(schema, className, object, acl)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	schema = getSchema()
	schema.reloadData()
	className = "user"
	object = nil
	acl = nil
	err = TomatoDBController.canAddField(schema, className, object, acl)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	Adapter.CreateClass("user", object)
	schema = getSchema()
	schema.reloadData()
	className = "user"
	object = nil
	acl = nil
	err = TomatoDBController.canAddField(schema, className, object, acl)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	Adapter.CreateClass("user", object)
	schema = getSchema()
	schema.reloadData()
	className = "user"
	object = types.M{
		"key":  "hello",
		"key1": "hello",
	}
	acl = nil
	err = TomatoDBController.canAddField(schema, className, object, acl)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"get": types.M{"*": true},
		},
	}
	Adapter.CreateClass("user", object)
	schema = getSchema()
	schema.reloadData()
	className = "user"
	object = types.M{
		"key":  "hello",
		"key1": "hello",
	}
	acl = nil
	err = TomatoDBController.canAddField(schema, className, object, acl)
	expect = errs.E(errs.OperationForbidden, "Permission denied for this action.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"addField": types.M{"*": true},
		},
	}
	Adapter.CreateClass("user", object)
	schema = getSchema()
	schema.reloadData()
	className = "user"
	object = types.M{
		"key":  "hello",
		"key1": "hello",
	}
	acl = nil
	err = TomatoDBController.canAddField(schema, className, object, acl)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"addField": types.M{"role:2048": true},
		},
	}
	Adapter.CreateClass("user", object)
	schema = getSchema()
	schema.reloadData()
	className = "user"
	object = types.M{
		"key":  "hello",
		"key1": "hello",
	}
	acl = []string{"role:1024"}
	err = TomatoDBController.canAddField(schema, className, object, acl)
	expect = errs.E(errs.OperationForbidden, "Permission denied for this action.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"addField": types.M{"role:2048": true},
		},
	}
	Adapter.CreateClass("user", object)
	schema = getSchema()
	schema.reloadData()
	className = "user"
	object = types.M{
		"key":  "hello",
		"key1": "hello",
	}
	acl = []string{"role:2048"}
	err = TomatoDBController.canAddField(schema, className, object, acl)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	Adapter.DeleteAllClasses()
}

func Test_reduceRelationKeys(t *testing.T) {
	initEnv()
	var object types.M
	var className string
	var query types.M
	var result types.M
	var expect types.M
	/*************************************************/
	className = "user"
	query = nil
	result = TomatoDBController.reduceRelationKeys(className, query)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	className = "user"
	query = types.M{}
	result = TomatoDBController.reduceRelationKeys(className, query)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	className = "user"
	query = types.M{"$relatedTo": "1024"}
	result = TomatoDBController.reduceRelationKeys(className, query)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	className = "user"
	query = types.M{"$relatedTo": types.M{}}
	result = TomatoDBController.reduceRelationKeys(className, query)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	className = "user"
	query = types.M{
		"$relatedTo": types.M{
			"object": types.M{
				"__type":    "Pointer",
				"className": "Post",
				"objectId":  "1001",
			},
			"key": "key",
		},
	}
	result = TomatoDBController.reduceRelationKeys(className, query)
	expect = types.M{
		"objectId": types.M{"$in": types.S{}},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"_id":       "1",
		"relatedId": "2001",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:key:Post", relationSchema, object)
	className = "user"
	query = types.M{
		"$relatedTo": types.M{
			"object": types.M{
				"__type":    "Pointer",
				"className": "Post",
				"objectId":  "1001",
			},
			"key": "key",
		},
	}
	result = TomatoDBController.reduceRelationKeys(className, query)
	expect = types.M{
		"objectId": types.M{"$in": types.S{"2001"}},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"_id":       "1",
		"relatedId": "2001",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:key:Post", relationSchema, object)
	object = types.M{
		"_id":       "2",
		"relatedId": "2002",
		"owningId":  "1002",
	}
	Adapter.CreateObject("_Join:key:Post", relationSchema, object)
	className = "user"
	query = types.M{
		"$or": types.S{
			types.M{
				"$relatedTo": types.M{
					"object": types.M{
						"__type":    "Pointer",
						"className": "Post",
						"objectId":  "1001",
					},
					"key": "key",
				},
			},
			types.M{
				"$relatedTo": types.M{
					"object": types.M{
						"__type":    "Pointer",
						"className": "Post",
						"objectId":  "1002",
					},
					"key": "key",
				},
			},
		},
	}
	result = TomatoDBController.reduceRelationKeys(className, query)
	expect = types.M{
		"$or": types.S{
			types.M{
				"objectId": types.M{"$in": types.S{"2001"}},
			},
			types.M{
				"objectId": types.M{"$in": types.S{"2002"}},
			},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
}

func Test_relatedIds(t *testing.T) {
	initEnv()
	var object types.M
	var className string
	var key string
	var owningID string
	var result types.S
	var expect types.S
	/*************************************************/
	className = "user"
	key = "name"
	owningID = "1001"
	result = TomatoDBController.relatedIds(className, key, owningID)
	expect = types.S{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"_id":       "1",
		"relatedId": "01",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:name:user", relationSchema, object)
	object = types.M{
		"_id":       "2",
		"relatedId": "02",
		"owningId":  "1002",
	}
	Adapter.CreateObject("_Join:name:user", relationSchema, object)
	object = types.M{
		"_id":       "3",
		"relatedId": "03",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:name:user", relationSchema, object)
	className = "user"
	key = "name"
	owningID = "1001"
	result = TomatoDBController.relatedIds(className, key, owningID)
	expect = types.S{"01", "03"}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
}

func Test_addInObjectIdsIds(t *testing.T) {
	initEnv()
	var ids types.S
	var query types.M
	var result types.M
	var expect types.M
	/*************************************************/
	ids = nil
	query = nil
	result = TomatoDBController.addInObjectIdsIds(ids, query)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	ids = types.S{}
	query = nil
	result = TomatoDBController.addInObjectIdsIds(ids, query)
	expect = types.M{
		"objectId": types.M{"$in": types.S{}},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	ids = nil
	query = types.M{}
	result = TomatoDBController.addInObjectIdsIds(ids, query)
	expect = types.M{
		"objectId": types.M{"$in": types.S{}},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	ids = types.S{}
	query = types.M{}
	result = TomatoDBController.addInObjectIdsIds(ids, query)
	expect = types.M{
		"objectId": types.M{"$in": types.S{}},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	ids = types.S{"1024"}
	query = types.M{"objectId": "1024"}
	result = TomatoDBController.addInObjectIdsIds(ids, query)
	expect = types.M{
		"objectId": types.M{"$in": types.S{"1024"}},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	ids = types.S{"1024"}
	query = types.M{"objectId": types.M{"$eq": "1024"}}
	result = TomatoDBController.addInObjectIdsIds(ids, query)
	expect = types.M{
		"objectId": types.M{
			"$eq": "1024",
			"$in": types.S{"1024"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	ids = types.S{"1024"}
	query = types.M{"objectId": types.M{"$in": types.S{"1024"}}}
	result = TomatoDBController.addInObjectIdsIds(ids, query)
	expect = types.M{
		"objectId": types.M{
			"$in": types.S{"1024"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	ids = types.S{"1024"}
	query = types.M{
		"objectId": types.M{
			"$in": types.S{"1024"},
			"$lt": "100",
		},
		"key": "value",
	}
	result = TomatoDBController.addInObjectIdsIds(ids, query)
	expect = types.M{
		"objectId": types.M{
			"$in": types.S{"1024"},
			"$lt": "100",
		},
		"key": "value",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	ids = types.S{"1024", "2048"}
	query = types.M{"objectId": types.M{"$in": types.S{"1024"}}}
	result = TomatoDBController.addInObjectIdsIds(ids, query)
	expect = types.M{
		"objectId": types.M{
			"$in": types.S{"1024"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	ids = types.S{"1024"}
	query = types.M{"objectId": types.M{"$in": types.S{"1024", "2048"}}}
	result = TomatoDBController.addInObjectIdsIds(ids, query)
	expect = types.M{
		"objectId": types.M{
			"$in": types.S{"1024"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	ids = types.S{"1024", "2048"}
	query = types.M{"objectId": types.M{"$in": types.S{"1024", "2048"}}}
	result = TomatoDBController.addInObjectIdsIds(ids, query)
	expect = types.M{
		"objectId": types.M{
			"$in": types.S{"1024", "2048"},
		},
	}
	expect1 := types.M{
		"objectId": types.M{
			"$in": types.S{"2048", "1024"},
		},
	}
	if reflect.DeepEqual(expect, result) == false &&
		reflect.DeepEqual(expect1, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	ids = types.S{"1024"}
	query = types.M{"objectId": types.M{"$in": types.S{"1024", "2048", "2048"}}}
	result = TomatoDBController.addInObjectIdsIds(ids, query)
	expect = types.M{
		"objectId": types.M{
			"$in": types.S{"1024"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	ids = types.S{"1024", "2048", "2048"}
	query = types.M{"objectId": types.M{"$in": types.S{"1024"}}}
	result = TomatoDBController.addInObjectIdsIds(ids, query)
	expect = types.M{
		"objectId": types.M{
			"$in": types.S{"1024"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_addNotInObjectIdsIds(t *testing.T) {
	initEnv()
	var ids types.S
	var query types.M
	var result types.M
	var expect types.M
	var expect1 types.M
	/*************************************************/
	ids = nil
	query = nil
	result = TomatoDBController.addNotInObjectIdsIds(ids, query)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	ids = types.S{}
	query = nil
	result = TomatoDBController.addNotInObjectIdsIds(ids, query)
	expect = types.M{
		"objectId": types.M{"$nin": types.S{}},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	ids = nil
	query = types.M{}
	result = TomatoDBController.addNotInObjectIdsIds(ids, query)
	expect = types.M{
		"objectId": types.M{"$nin": types.S{}},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	ids = types.S{}
	query = types.M{}
	result = TomatoDBController.addNotInObjectIdsIds(ids, query)
	expect = types.M{
		"objectId": types.M{"$nin": types.S{}},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	ids = types.S{"1024"}
	query = types.M{"objectId": types.M{"$nin": types.S{"1024"}}}
	result = TomatoDBController.addNotInObjectIdsIds(ids, query)
	expect = types.M{
		"objectId": types.M{
			"$nin": types.S{"1024"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	ids = types.S{"1024"}
	query = types.M{
		"objectId": types.M{
			"$nin": types.S{"1024"},
			"$lt":  "100",
		},
		"key": "value",
	}
	result = TomatoDBController.addNotInObjectIdsIds(ids, query)
	expect = types.M{
		"objectId": types.M{
			"$nin": types.S{"1024"},
			"$lt":  "100",
		},
		"key": "value",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	ids = types.S{"1024", "2048"}
	query = types.M{"objectId": types.M{"$nin": types.S{"1024"}}}
	result = TomatoDBController.addNotInObjectIdsIds(ids, query)
	expect = types.M{
		"objectId": types.M{
			"$nin": types.S{"1024", "2048"},
		},
	}
	expect1 = types.M{
		"objectId": types.M{
			"$nin": types.S{"2048", "1024"},
		},
	}
	if reflect.DeepEqual(expect, result) == false &&
		reflect.DeepEqual(expect1, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	ids = types.S{"1024"}
	query = types.M{"objectId": types.M{"$nin": types.S{"1024", "2048"}}}
	result = TomatoDBController.addNotInObjectIdsIds(ids, query)
	expect = types.M{
		"objectId": types.M{
			"$nin": types.S{"1024", "2048"},
		},
	}
	expect1 = types.M{
		"objectId": types.M{
			"$nin": types.S{"2048", "1024"},
		},
	}
	if reflect.DeepEqual(expect, result) == false &&
		reflect.DeepEqual(expect1, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	ids = types.S{"1024", "2048"}
	query = types.M{"objectId": types.M{"$nin": types.S{"1024", "2048"}}}
	result = TomatoDBController.addNotInObjectIdsIds(ids, query)
	expect = types.M{
		"objectId": types.M{
			"$nin": types.S{"1024", "2048"},
		},
	}
	expect1 = types.M{
		"objectId": types.M{
			"$nin": types.S{"2048", "1024"},
		},
	}
	if reflect.DeepEqual(expect, result) == false &&
		reflect.DeepEqual(expect1, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	ids = types.S{"1024"}
	query = types.M{"objectId": types.M{"$nin": types.S{"1024", "2048", "2048"}}}
	result = TomatoDBController.addNotInObjectIdsIds(ids, query)
	expect = types.M{
		"objectId": types.M{
			"$nin": types.S{"1024", "2048"},
		},
	}
	expect1 = types.M{
		"objectId": types.M{
			"$nin": types.S{"2048", "1024"},
		},
	}
	if reflect.DeepEqual(expect, result) == false &&
		reflect.DeepEqual(expect1, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	ids = types.S{"1024", "2048", "2048"}
	query = types.M{"objectId": types.M{"$nin": types.S{"1024"}}}
	result = TomatoDBController.addNotInObjectIdsIds(ids, query)
	expect = types.M{
		"objectId": types.M{
			"$nin": types.S{"1024", "2048"},
		},
	}
	expect1 = types.M{
		"objectId": types.M{
			"$nin": types.S{"2048", "1024"},
		},
	}
	if reflect.DeepEqual(expect, result) == false &&
		reflect.DeepEqual(expect1, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_reduceInRelation(t *testing.T) {
	initEnv()
	var object types.M
	var className string
	var query types.M
	var schema *Schema
	var result types.M
	var expect types.M
	/*************************************************/
	className = "user"
	query = nil
	schema = nil
	result = TomatoDBController.reduceInRelation(className, query, schema)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	className = "user"
	query = types.M{}
	schema = nil
	result = TomatoDBController.reduceInRelation(className, query, schema)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	className = "user"
	query = types.M{"key": "hello"}
	schema = nil
	result = TomatoDBController.reduceInRelation(className, query, schema)
	expect = types.M{"key": "hello"}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	className = "user"
	query = types.M{"key": types.M{"k": "v"}}
	schema = nil
	result = TomatoDBController.reduceInRelation(className, query, schema)
	expect = types.M{"key": types.M{"k": "v"}}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	className = "user"
	query = types.M{"key": types.M{"$in": "v"}}
	schema = nil
	result = TomatoDBController.reduceInRelation(className, query, schema)
	expect = types.M{"key": types.M{"$in": "v"}}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"className": "user",
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	Adapter.CreateClass("user", object)
	className = "user"
	query = types.M{"key": types.M{"$in": "v"}}
	schema = getSchema()
	schema.reloadData()
	result = TomatoDBController.reduceInRelation(className, query, schema)
	expect = types.M{"key": types.M{"$in": "v"}}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"className": "user",
		"fields": types.M{
			"key": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
		},
	}
	Adapter.CreateClass("user", object)
	object = types.M{
		"_id":       "1",
		"relatedId": "2001",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:key:user", relationSchema, object)
	className = "user"
	query = types.M{
		"key": types.M{
			"__type":   "Pointer",
			"objectId": "2001",
		},
	}
	schema = getSchema()
	schema.reloadData()
	result = TomatoDBController.reduceInRelation(className, query, schema)
	expect = types.M{
		"objectId": types.M{
			"$in": types.S{"1001"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"className": "user",
		"fields": types.M{
			"key": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
		},
	}
	Adapter.CreateClass("user", object)
	object = types.M{
		"_id":       "1",
		"relatedId": "2001",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:key:user", relationSchema, object)
	className = "user"
	query = types.M{
		"key": types.M{
			"$in": types.S{
				types.M{
					"__type":   "Pointer",
					"objectId": "2001",
				},
			},
		},
	}
	schema = getSchema()
	schema.reloadData()
	result = TomatoDBController.reduceInRelation(className, query, schema)
	expect = types.M{
		"objectId": types.M{
			"$in": types.S{"1001"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"className": "user",
		"fields": types.M{
			"key": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
		},
	}
	Adapter.CreateClass("user", object)
	object = types.M{
		"_id":       "1",
		"relatedId": "2001",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:key:user", relationSchema, object)
	className = "user"
	query = types.M{
		"key": types.M{
			"$nin": types.S{
				types.M{
					"__type":   "Pointer",
					"objectId": "2001",
				},
			},
		},
	}
	schema = getSchema()
	schema.reloadData()
	result = TomatoDBController.reduceInRelation(className, query, schema)
	expect = types.M{
		"objectId": types.M{
			"$nin": types.S{"1001"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"className": "user",
		"fields": types.M{
			"key": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
		},
	}
	Adapter.CreateClass("user", object)
	object = types.M{
		"_id":       "1",
		"relatedId": "2001",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:key:user", relationSchema, object)
	className = "user"
	query = types.M{
		"key": types.M{
			"$ne": types.M{
				"__type":   "Pointer",
				"objectId": "2001",
			},
		},
	}
	schema = getSchema()
	schema.reloadData()
	result = TomatoDBController.reduceInRelation(className, query, schema)
	expect = types.M{
		"objectId": types.M{
			"$nin": types.S{"1001"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"className": "user",
		"fields": types.M{
			"key": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
		},
	}
	Adapter.CreateClass("user", object)
	object = types.M{
		"_id":       "1",
		"relatedId": "2001",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:key:user", relationSchema, object)
	className = "user"
	query = types.M{
		"key": types.M{
			"$eq": types.M{
				"__type":   "Pointer",
				"objectId": "2001",
			},
		},
	}
	schema = getSchema()
	schema.reloadData()
	result = TomatoDBController.reduceInRelation(className, query, schema)
	expect = types.M{
		"objectId": types.M{
			"$in": types.S{"1001"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"className": "user",
		"fields": types.M{
			"key": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
		},
	}
	Adapter.CreateClass("user", object)
	object = types.M{
		"_id":       "1",
		"relatedId": "2001",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:key:user", relationSchema, object)
	className = "user"
	query = types.M{
		"key": types.M{
			"$eq": types.M{
				"__type":   "Pointer",
				"objectId": "2001",
			},
		},
		"key2": "hello",
	}
	schema = getSchema()
	schema.reloadData()
	result = TomatoDBController.reduceInRelation(className, query, schema)
	expect = types.M{
		"objectId": types.M{
			"$in": types.S{"1001"},
		},
		"key2": "hello",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"className": "user",
		"fields": types.M{
			"key": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
		},
	}
	Adapter.CreateClass("user", object)
	object = types.M{
		"_id":       "1",
		"relatedId": "2001",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:key:user", relationSchema, object)
	className = "user"
	query = types.M{
		"key": types.M{
			"$eq": "hello",
		},
		"key2": "hello",
	}
	schema = getSchema()
	schema.reloadData()
	result = TomatoDBController.reduceInRelation(className, query, schema)
	expect = types.M{
		"key2": "hello",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"className": "user",
		"fields": types.M{
			"key": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
			"key1": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
		},
	}
	Adapter.CreateClass("user", object)
	object = types.M{
		"_id":       "1",
		"relatedId": "2001",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:key:user", relationSchema, object)
	object = types.M{
		"_id":       "2",
		"relatedId": "2002",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:key1:user", relationSchema, object)
	className = "user"
	query = types.M{
		"key": types.M{
			"$eq": types.M{
				"__type":   "Pointer",
				"objectId": "2001",
			},
		},
		"key1": types.M{
			"$eq": types.M{
				"__type":   "Pointer",
				"objectId": "2002",
			},
		},
		"key2": "hello",
	}
	schema = getSchema()
	schema.reloadData()
	result = TomatoDBController.reduceInRelation(className, query, schema)
	expect = types.M{
		"objectId": types.M{
			"$in": types.S{"1001"},
		},
		"key2": "hello",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"className": "user",
		"fields": types.M{
			"key": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
			"key1": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
		},
	}
	Adapter.CreateClass("user", object)
	object = types.M{
		"_id":       "1",
		"relatedId": "2001",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:key:user", relationSchema, object)
	object = types.M{
		"_id":       "2",
		"relatedId": "2002",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:key1:user", relationSchema, object)
	className = "user"
	query = types.M{
		"$or": types.S{
			types.M{
				"key": types.M{
					"$eq": types.M{
						"__type":   "Pointer",
						"objectId": "2001",
					},
				},
			},
			types.M{
				"key1": types.M{
					"$eq": types.M{
						"__type":   "Pointer",
						"objectId": "2002",
					},
				},
			},
		},
	}
	schema = getSchema()
	schema.reloadData()
	result = TomatoDBController.reduceInRelation(className, query, schema)
	expect = types.M{
		"$or": types.S{
			types.M{
				"objectId": types.M{
					"$in": types.S{"1001"},
				},
			},
			types.M{
				"objectId": types.M{
					"$in": types.S{"1001"},
				},
			},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
}

func Test_owningIds(t *testing.T) {
	initEnv()
	var object types.M
	var className string
	var key string
	var relatedIds types.S
	var result types.S
	var expect types.S
	/*************************************************/
	className = "user"
	key = "name"
	relatedIds = types.S{"01", "02"}
	result = TomatoDBController.owningIds(className, key, relatedIds)
	expect = types.S{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"_id":       "1",
		"relatedId": "01",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:name:user", relationSchema, object)
	object = types.M{
		"_id":       "2",
		"relatedId": "02",
		"owningId":  "1002",
	}
	Adapter.CreateObject("_Join:name:user", relationSchema, object)
	object = types.M{
		"_id":       "3",
		"relatedId": "03",
		"owningId":  "1003",
	}
	Adapter.CreateObject("_Join:name:user", relationSchema, object)
	className = "user"
	key = "name"
	relatedIds = types.S{"01", "02"}
	result = TomatoDBController.owningIds(className, key, relatedIds)
	expect = types.S{"1001", "1002"}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
}

func Test_DeleteSchema(t *testing.T) {
	// TODO
}

func Test_addPointerPermissions(t *testing.T) {
	initEnv()
	var object types.M
	var schema *Schema
	var className string
	var operation string
	var query types.M
	var aclGroup []string
	var result types.M
	var expect types.M
	/*************************************************/
	schema = nil
	className = "user"
	operation = "get"
	query = nil
	aclGroup = nil
	result = TomatoDBController.addPointerPermissions(schema, className, operation, query, aclGroup)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	schema = getSchema()
	schema.reloadData()
	className = "user"
	operation = "get"
	query = nil
	aclGroup = nil
	result = TomatoDBController.addPointerPermissions(schema, className, operation, query, aclGroup)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	className = "user"
	object = types.M{
		"className": className,
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"get": types.M{"*": true},
		},
	}
	Adapter.CreateClass(className, object)
	schema = getSchema()
	schema.reloadData()
	className = "user"
	operation = "get"
	query = nil
	aclGroup = nil
	result = TomatoDBController.addPointerPermissions(schema, className, operation, query, aclGroup)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	className = "user"
	object = types.M{
		"className": className,
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"get": types.M{"role:1024": true},
		},
	}
	Adapter.CreateClass(className, object)
	schema = getSchema()
	schema.reloadData()
	className = "user"
	operation = "get"
	query = types.M{}
	aclGroup = []string{"123456789012345678901234"}
	result = TomatoDBController.addPointerPermissions(schema, className, operation, query, aclGroup)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	className = "user"
	object = types.M{
		"className": className,
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"get":            types.M{"role:1024": true},
			"readUserFields": types.S{"key2"},
		},
	}
	Adapter.CreateClass(className, object)
	schema = getSchema()
	schema.reloadData()
	className = "user"
	operation = "get"
	query = types.M{
		"key": "hello",
	}
	aclGroup = []string{"123456789012345678901234"}
	result = TomatoDBController.addPointerPermissions(schema, className, operation, query, aclGroup)
	expect = types.M{
		"$and": types.S{
			types.M{
				"key2": types.M{
					"__type":    "Pointer",
					"className": "_User",
					"objectId":  "123456789012345678901234",
				},
			},
			types.M{
				"key": "hello",
			},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	className = "user"
	object = types.M{
		"className": className,
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"update":          types.M{"role:1024": true},
			"writeUserFields": types.S{"key2"},
		},
	}
	Adapter.CreateClass(className, object)
	schema = getSchema()
	schema.reloadData()
	className = "user"
	operation = "update"
	query = types.M{
		"key": "hello",
	}
	aclGroup = []string{"123456789012345678901234"}
	result = TomatoDBController.addPointerPermissions(schema, className, operation, query, aclGroup)
	expect = types.M{
		"$and": types.S{
			types.M{
				"key2": types.M{
					"__type":    "Pointer",
					"className": "_User",
					"objectId":  "123456789012345678901234",
				},
			},
			types.M{
				"key": "hello",
			},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	className = "user"
	object = types.M{
		"className": className,
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"get":            types.M{"role:1024": true},
			"readUserFields": types.S{"key2"},
		},
	}
	Adapter.CreateClass(className, object)
	schema = getSchema()
	schema.reloadData()
	className = "user"
	operation = "get"
	query = types.M{
		"key": "hello",
	}
	aclGroup = []string{"role:2048"}
	result = TomatoDBController.addPointerPermissions(schema, className, operation, query, aclGroup)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	className = "user"
	object = types.M{
		"className": className,
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"get":            types.M{"role:1024": true},
			"readUserFields": types.S{"key2", "key3"},
		},
	}
	Adapter.CreateClass(className, object)
	schema = getSchema()
	schema.reloadData()
	className = "user"
	operation = "get"
	query = types.M{
		"key": "hello",
	}
	aclGroup = []string{"123456789012345678901234"}
	result = TomatoDBController.addPointerPermissions(schema, className, operation, query, aclGroup)
	expect = types.M{
		"$or": []types.M{
			types.M{
				"$and": types.S{
					types.M{
						"key2": types.M{
							"__type":    "Pointer",
							"className": "_User",
							"objectId":  "123456789012345678901234",
						},
					},
					types.M{"key": "hello"},
				},
			},
			types.M{
				"$and": types.S{
					types.M{
						"key3": types.M{
							"__type":    "Pointer",
							"className": "_User",
							"objectId":  "123456789012345678901234",
						},
					},
					types.M{"key": "hello"},
				},
			},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
}

//////////////////////////////////////////////////////

func Test_sanitizeDatabaseResult(t *testing.T) {
	var originalObject types.M
	var object types.M
	var result types.M
	var expect types.M
	/*************************************************/
	originalObject = nil
	object = nil
	result = sanitizeDatabaseResult(originalObject, object)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	originalObject = types.M{
		"key": types.M{
			"__op":    "Add",
			"objects": types.S{"hello", "world"},
		},
		"key2": "hello",
	}
	object = types.M{
		"key":  types.S{"hello", "world"},
		"key2": "hello",
	}
	result = sanitizeDatabaseResult(originalObject, object)
	expect = types.M{
		"key": types.S{"hello", "world"},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	originalObject = types.M{
		"key": types.M{
			"__op":    "AddUnique",
			"objects": types.S{"hello", "world"},
		},
		"key2": "hello",
	}
	object = types.M{
		"key":  types.S{"hello", "world"},
		"key2": "hello",
	}
	result = sanitizeDatabaseResult(originalObject, object)
	expect = types.M{
		"key": types.S{"hello", "world"},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	originalObject = types.M{
		"key": types.M{
			"__op":    "Remove",
			"objects": types.S{"hello", "world"},
		},
		"key2": "hello",
	}
	object = types.M{
		"key":  types.S{"value"},
		"key2": "hello",
	}
	result = sanitizeDatabaseResult(originalObject, object)
	expect = types.M{
		"key": types.S{"value"},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	originalObject = types.M{
		"key": types.M{
			"__op":   "Increment",
			"amount": 10,
		},
		"key2": "hello",
	}
	object = types.M{
		"key":  20,
		"key2": "hello",
	}
	result = sanitizeDatabaseResult(originalObject, object)
	expect = types.M{
		"key": 20,
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_keysForQuery(t *testing.T) {
	var query types.M
	var result []string
	var expect []string
	/*************************************************/
	query = nil
	result = keysForQuery(query)
	expect = []string{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	query = types.M{
		"$and": 1024,
	}
	result = keysForQuery(query)
	expect = []string{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	query = types.M{
		"$and": types.S{"hello"},
	}
	result = keysForQuery(query)
	expect = []string{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	query = types.M{
		"$and": types.S{
			types.M{"key": "hello"},
			types.M{"key2": "hello"},
		},
	}
	result = keysForQuery(query)
	expect = []string{"key", "key2"}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	query = types.M{
		"$or": 1024,
	}
	result = keysForQuery(query)
	expect = []string{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	query = types.M{
		"$or": types.S{"hello"},
	}
	result = keysForQuery(query)
	expect = []string{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	query = types.M{
		"$or": types.S{
			types.M{"key": "hello"},
			types.M{"key2": "hello"},
		},
	}
	result = keysForQuery(query)
	expect = []string{"key", "key2"}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	query = types.M{
		"key": "hello",
	}
	result = keysForQuery(query)
	expect = []string{"key"}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_joinTableName(t *testing.T) {
	var className string
	var key string
	var result string
	var expect string
	/*************************************************/
	className = "user"
	key = "name"
	result = joinTableName(className, key)
	expect = "_Join:name:user"
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_filterSensitiveData(t *testing.T) {
	var isMaster bool
	var aclGroup []string
	var className string
	var object types.M
	var result types.M
	var expect types.M
	/*************************************************/
	className = "other"
	isMaster = false
	aclGroup = nil
	object = types.M{"key": "value"}
	result = filterSensitiveData(isMaster, aclGroup, className, object)
	expect = types.M{"key": "value"}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	className = "_User"
	isMaster = false
	aclGroup = nil
	object = nil
	result = filterSensitiveData(isMaster, aclGroup, className, object)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	className = "_User"
	isMaster = false
	aclGroup = nil
	object = types.M{"_hashed_password": "1024"}
	result = filterSensitiveData(isMaster, aclGroup, className, object)
	expect = types.M{"password": "1024"}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	className = "_User"
	isMaster = false
	aclGroup = nil
	object = types.M{
		"_hashed_password": "1024",
		"sessionToken":     "abc",
	}
	result = filterSensitiveData(isMaster, aclGroup, className, object)
	expect = types.M{"password": "1024"}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	className = "_User"
	isMaster = false
	aclGroup = nil
	object = types.M{
		"_hashed_password": "1024",
		"sessionToken":     "abc",
		"authData": types.M{
			"facebook": types.M{"id": "1024"},
		},
	}
	result = filterSensitiveData(isMaster, aclGroup, className, object)
	expect = types.M{"password": "1024"}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	className = "_User"
	isMaster = true
	aclGroup = nil
	object = types.M{
		"_hashed_password": "1024",
		"sessionToken":     "abc",
		"authData": types.M{
			"facebook": types.M{"id": "1024"},
		},
	}
	result = filterSensitiveData(isMaster, aclGroup, className, object)
	expect = types.M{
		"password": "1024",
		"authData": types.M{
			"facebook": types.M{"id": "1024"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	className = "_User"
	isMaster = false
	aclGroup = []string{"1024"}
	object = types.M{
		"objectId":         "1024",
		"_hashed_password": "1024",
		"sessionToken":     "abc",
		"authData": types.M{
			"facebook": types.M{"id": "1024"},
		},
	}
	result = filterSensitiveData(isMaster, aclGroup, className, object)
	expect = types.M{
		"objectId": "1024",
		"password": "1024",
		"authData": types.M{
			"facebook": types.M{"id": "1024"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_addWriteACL(t *testing.T) {
	var query types.M
	var acl []string
	var result types.M
	var expect types.M
	/*************************************************/
	query = nil
	acl = nil
	result = addWriteACL(query, acl)
	expect = types.M{
		"_wperm": types.M{
			"$in": types.S{nil},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	query = types.M{"key": "hello"}
	acl = nil
	result = addWriteACL(query, acl)
	expect = types.M{
		"key": "hello",
		"_wperm": types.M{
			"$in": types.S{nil},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	query = types.M{"key": "hello"}
	acl = []string{"role:1024"}
	result = addWriteACL(query, acl)
	expect = types.M{
		"key": "hello",
		"_wperm": types.M{
			"$in": types.S{nil, "role:1024"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_addReadACL(t *testing.T) {
	var query types.M
	var acl []string
	var result types.M
	var expect types.M
	/*************************************************/
	query = nil
	acl = nil
	result = addReadACL(query, acl)
	expect = types.M{
		"_rperm": types.M{
			"$in": types.S{nil, "*"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	query = types.M{"key": "hello"}
	acl = nil
	result = addReadACL(query, acl)
	expect = types.M{
		"key": "hello",
		"_rperm": types.M{
			"$in": types.S{nil, "*"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	query = types.M{"key": "hello"}
	acl = []string{"role:1024"}
	result = addReadACL(query, acl)
	expect = types.M{
		"key": "hello",
		"_rperm": types.M{
			"$in": types.S{nil, "*", "role:1024"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_validateQuery(t *testing.T) {
	var query types.M
	var err error
	var expect error
	/*************************************************/
	query = nil
	err = validateQuery(query)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/*************************************************/
	query = types.M{"ACL": "ACL"}
	err = validateQuery(query)
	expect = errs.E(errs.InvalidQuery, "Cannot query on ACL.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/*************************************************/
	query = types.M{
		"key": types.M{
			"$regex":   "hello",
			"$options": "imxs",
		},
	}
	err = validateQuery(query)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/*************************************************/
	query = types.M{
		"key": types.M{
			"$regex":   "hello",
			"$options": "abc",
		},
	}
	err = validateQuery(query)
	expect = errs.E(errs.InvalidQuery, "Bad $options value for query: abc")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/*************************************************/
	query = types.M{
		"_rperm":              "hello",
		"_wperm":              "hello",
		"_perishable_token":   "hello",
		"_email_verify_token": "hello",
	}
	err = validateQuery(query)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/*************************************************/
	query = types.M{
		"_other": "hello",
	}
	err = validateQuery(query)
	expect = errs.E(errs.InvalidKeyName, "Invalid key name: _other")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/*************************************************/
	query = types.M{
		"$or": "hello",
	}
	err = validateQuery(query)
	expect = errs.E(errs.InvalidQuery, "Bad $or format - use an array value.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/*************************************************/
	query = types.M{
		"$or": types.S{"hello", "world"},
	}
	err = validateQuery(query)
	expect = errs.E(errs.InvalidQuery, "Bad $or format - invalid sub query.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/*************************************************/
	query = types.M{
		"$or": types.S{
			types.M{"key": "value"},
			types.M{"key": "value"},
		},
	}
	err = validateQuery(query)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/*************************************************/
	query = types.M{
		"$and": "hello",
	}
	err = validateQuery(query)
	expect = errs.E(errs.InvalidQuery, "Bad $and format - use an array value.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/*************************************************/
	query = types.M{
		"$and": types.S{"hello", "world"},
	}
	err = validateQuery(query)
	expect = errs.E(errs.InvalidQuery, "Bad $and format - invalid sub query.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/*************************************************/
	query = types.M{
		"$and": types.S{
			types.M{"key": "value"},
			types.M{"key": "value"},
		},
	}
	err = validateQuery(query)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
}

func Test_transformObjectACL(t *testing.T) {
	var object types.M
	var result types.M
	var expect types.M
	/*************************************************/
	object = nil
	result = transformObjectACL(object)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	object = types.M{}
	result = transformObjectACL(object)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	object = types.M{"ACL": "hello"}
	result = transformObjectACL(object)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	object = types.M{
		"ACL": types.M{
			"userid": types.M{
				"read":  true,
				"write": true,
			},
			"role:xxx": types.M{
				"read":  true,
				"write": true,
			},
			"*": types.M{
				"read": true,
			},
		},
	}
	result = transformObjectACL(object)
	expect = types.M{
		"_rperm": types.S{"userid", "role:xxx", "*"},
		"_wperm": types.S{"userid", "role:xxx"},
	}
	if utils.CompareArray(expect["_rperm"], result["_rperm"]) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	if utils.CompareArray(expect["_wperm"], result["_wperm"]) == false {
		t.Error("expect:", expect, "get result:", result)
	}
}

func Test_untransformObjectACL(t *testing.T) {
	var output types.M
	var result types.M
	var expect types.M
	/*************************************************/
	output = nil
	result = untransformObjectACL(output)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	output = types.M{}
	result = untransformObjectACL(output)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	output = types.M{"_rperm": "Incorrect type"}
	result = untransformObjectACL(output)
	expect = types.M{"ACL": types.M{}}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	output = types.M{"_rperm": types.S{"userid", "role:xxx"}}
	result = untransformObjectACL(output)
	expect = types.M{
		"ACL": types.M{
			"userid": types.M{
				"read": true,
			},
			"role:xxx": types.M{
				"read": true,
			},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	output = types.M{"_wperm": "Incorrect type"}
	result = untransformObjectACL(output)
	expect = types.M{"ACL": types.M{}}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	output = types.M{"_wperm": types.S{"userid", "role:xxx"}}
	result = untransformObjectACL(output)
	expect = types.M{
		"ACL": types.M{
			"userid": types.M{
				"write": true,
			},
			"role:xxx": types.M{
				"write": true,
			},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*************************************************/
	output = types.M{
		"_rperm": types.S{"userid", "role:xxx", "*"},
		"_wperm": types.S{"userid", "role:xxx"},
	}
	result = untransformObjectACL(output)
	expect = types.M{
		"ACL": types.M{
			"userid": types.M{
				"read":  true,
				"write": true,
			},
			"role:xxx": types.M{
				"read":  true,
				"write": true,
			},
			"*": types.M{
				"read": true,
			},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_transformAuthData(t *testing.T) {
	var className string
	var object types.M
	var expect types.M
	var schema types.M
	/*************************************************/
	className = "Other"
	object = types.M{"key": "value"}
	schema = nil
	transformAuthData(className, object, schema)
	expect = types.M{"key": "value"}
	if reflect.DeepEqual(expect, object) == false {
		t.Error("expect:", expect, "result:", object)
	}
	/*************************************************/
	className = "_User"
	object = nil
	schema = nil
	transformAuthData(className, object, schema)
	expect = nil
	if reflect.DeepEqual(expect, object) == false {
		t.Error("expect:", expect, "result:", object)
	}
	/*************************************************/
	className = "_User"
	object = types.M{}
	schema = nil
	transformAuthData(className, object, schema)
	expect = types.M{}
	if reflect.DeepEqual(expect, object) == false {
		t.Error("expect:", expect, "result:", object)
	}
	/*************************************************/
	className = "_User"
	object = types.M{"authData": nil}
	schema = nil
	transformAuthData(className, object, schema)
	expect = types.M{}
	if reflect.DeepEqual(expect, object) == false {
		t.Error("expect:", expect, "result:", object)
	}
	/*************************************************/
	className = "_User"
	object = types.M{"authData": 1024}
	schema = nil
	transformAuthData(className, object, schema)
	expect = types.M{}
	if reflect.DeepEqual(expect, object) == false {
		t.Error("expect:", expect, "result:", object)
	}
	/*************************************************/
	className = "_User"
	object = types.M{
		"authData": types.M{
			"facebook": nil,
		},
	}
	schema = nil
	transformAuthData(className, object, schema)
	expect = types.M{
		"_auth_data_facebook": types.M{
			"__op": "Delete",
		},
	}
	if reflect.DeepEqual(expect, object) == false {
		t.Error("expect:", expect, "result:", object)
	}
	/*************************************************/
	className = "_User"
	object = types.M{
		"authData": types.M{
			"facebook": 1024,
		},
	}
	schema = nil
	transformAuthData(className, object, schema)
	expect = types.M{
		"_auth_data_facebook": types.M{
			"__op": "Delete",
		},
	}
	if reflect.DeepEqual(expect, object) == false {
		t.Error("expect:", expect, "result:", object)
	}
	/*************************************************/
	className = "_User"
	object = types.M{
		"authData": types.M{
			"facebook": types.M{},
		},
	}
	schema = nil
	transformAuthData(className, object, schema)
	expect = types.M{
		"_auth_data_facebook": types.M{
			"__op": "Delete",
		},
	}
	if reflect.DeepEqual(expect, object) == false {
		t.Error("expect:", expect, "result:", object)
	}
	/*************************************************/
	className = "_User"
	object = types.M{
		"authData": types.M{
			"facebook": types.M{"id": "1024"},
			"twitter":  types.M{},
		},
	}
	schema = nil
	transformAuthData(className, object, schema)
	expect = types.M{
		"_auth_data_facebook": types.M{"id": "1024"},
		"_auth_data_twitter":  types.M{"__op": "Delete"},
	}
	if reflect.DeepEqual(expect, object) == false {
		t.Error("expect:", expect, "result:", object)
	}
	/*************************************************/
	className = "_User"
	object = types.M{
		"authData": types.M{
			"facebook": types.M{"id": "1024"},
			"twitter":  types.M{},
		},
	}
	schema = types.M{
		"fields": types.M{},
	}
	transformAuthData(className, object, schema)
	expect = types.M{
		"_auth_data_facebook": types.M{"id": "1024"},
		"_auth_data_twitter":  types.M{"__op": "Delete"},
	}
	if reflect.DeepEqual(expect, object) == false {
		t.Error("expect:", expect, "result:", object)
	}
	expect = types.M{
		"fields": types.M{
			"_auth_data_facebook": types.M{"type": "Object"},
		},
	}
	if reflect.DeepEqual(expect, schema) == false {
		t.Error("expect:", expect, "result:", schema)
	}
}

func Test_flattenUpdateOperatorsForCreate(t *testing.T) {
	var object types.M
	var err error
	var expect interface{}
	/**********************************************************/
	object = nil
	err = flattenUpdateOperatorsForCreate(object)
	expect = nil
	if err != nil || object != nil {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{"key": "value"}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{"key": "value"}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{"key": "value"},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{
		"key": types.M{"key": "value"},
	}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":   "Increment",
			"amount": 10.24,
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{
		"key": 10.24,
	}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":   "Increment",
			"amount": 1024,
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{
		"key": 1024,
	}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":   "Increment",
			"amount": "hello",
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = errs.E(errs.InvalidJSON, "objects to add must be an number")
	if err == nil || reflect.DeepEqual(err, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op": "Increment",
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = errs.E(errs.InvalidJSON, "objects to add must be an number")
	if err == nil || reflect.DeepEqual(err, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":    "Add",
			"objects": types.S{"abc", "def"},
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{
		"key": types.S{"abc", "def"},
	}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":    "Add",
			"objects": "hello",
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = errs.E(errs.InvalidJSON, "objects to add must be an array")
	if err == nil || reflect.DeepEqual(err, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":    "AddUnique",
			"objects": types.S{"abc", "def"},
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{
		"key": types.S{"abc", "def"},
	}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":    "AddUnique",
			"objects": "hello",
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = errs.E(errs.InvalidJSON, "objects to add must be an array")
	if err == nil || reflect.DeepEqual(err, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":    "Remove",
			"objects": types.S{"abc", "def"},
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{
		"key": types.S{},
	}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":    "Remove",
			"objects": "hello",
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = errs.E(errs.InvalidJSON, "objects to add must be an array")
	if err == nil || reflect.DeepEqual(err, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op": "Delete",
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op": "Other",
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = errs.E(errs.CommandUnavailable, "The Other operator is not supported yet.")
	if err == nil || reflect.DeepEqual(err, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
}

func initEnv() {
	Adapter = getAdapter()
	TomatoDBController = &DBController{}
}
