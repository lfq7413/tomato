package orm

import (
	"reflect"
	"testing"
	"time"

	"github.com/lfq7413/tomato/cache"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

func TestPostgres_CollectionExists(t *testing.T) {
	initPostgresEnv()
	var object types.M
	var schema types.M
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
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	Adapter.CreateObject(className, schema, object)
	className = "user"
	result = TomatoDBController.CollectionExists(className)
	expect = true
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	Adapter.DeleteAllClasses()
}

func TestPostgres_PurgeCollection(t *testing.T) {
	initPostgresEnv()
	var schema types.M
	var object types.M
	var className string
	var err error
	var expect error
	var resluts []types.M
	var expects []types.M
	/*************************************************/
	className = "user"
	err = TomatoDBController.PurgeCollection(className)
	expect = errs.E(errs.ObjectNotFound, "Object not found.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	Adapter.CreateClass("user", schema)
	className = "user"
	object = types.M{"key": "001"}
	Adapter.CreateObject(className, schema, object)
	object = types.M{"key": "002"}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	err = TomatoDBController.PurgeCollection(className)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	resluts, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
	expects = []types.M{}
	if reflect.DeepEqual(expects, resluts) == false {
		t.Error("expect:", expects, "result:", resluts)
	}
	TomatoDBController.DeleteEverything()
}

func TestPostgres_Find(t *testing.T) {
	initPostgresEnv()
	var schema types.M
	var object types.M
	var className string
	var query types.M
	var options types.M
	var results types.S
	var err error
	var expects types.S
	var expectErr error
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "post"
	query = nil
	options = nil
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "post"
	query = nil
	options = types.M{"count": true}
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{0}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = nil
	options = nil
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{
		types.M{
			"objectId": "01",
			"key":      "hello",
		},
		types.M{
			"objectId": "02",
			"key":      "hello",
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "03",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = nil
	options = types.M{"skip": 1}
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{
		types.M{
			"objectId": "02",
			"key":      "hello",
		},
		types.M{
			"objectId": "03",
			"key":      "hello",
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "03",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = nil
	options = types.M{"limit": 2}
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{
		types.M{
			"objectId": "01",
			"key":      "hello",
		},
		types.M{
			"objectId": "02",
			"key":      "hello",
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{"objectId": "01"}
	options = nil
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{
		types.M{
			"objectId": "01",
			"key":      "hello",
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{"objectId": "03"}
	options = nil
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "Number"},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      3,
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      1,
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "03",
		"key":      2,
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{}
	options = types.M{"sort": []string{"key"}}
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{
		types.M{
			"objectId": "02",
			"key":      1.0,
		},
		types.M{
			"objectId": "03",
			"key":      2.0,
		},
		types.M{
			"objectId": "01",
			"key":      3.0,
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "Number"},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      3,
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      1,
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "03",
		"key":      2,
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{}
	options = types.M{"sort": []string{"-key"}}
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{
		types.M{
			"objectId": "01",
			"key":      3.0,
		},
		types.M{
			"objectId": "03",
			"key":      2.0,
		},
		types.M{
			"objectId": "02",
			"key":      1.0,
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "Number"},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      3,
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      1,
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "03",
		"key":      2,
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{}
	options = types.M{"sort": []string{"@key"}}
	results, err = TomatoDBController.Find(className, query, options)
	expectErr = errs.E(errs.InvalidKeyName, "Invalid field name: @key")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "Number"},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      3,
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      1,
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "03",
		"key":      2,
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{}
	options = types.M{"sort": []string{"authData.facebook.id"}}
	results, err = TomatoDBController.Find(className, query, options)
	expectErr = errs.E(errs.InvalidKeyName, "Cannot sort by authData.facebook.id")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "Number"},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      3,
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      1,
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "03",
		"key":      2,
	}
	Adapter.CreateObject(className, schema, object)
	className = "post"
	schema = types.M{
		"fields": types.M{
			"user": types.M{"type": "Relation", "targetClass": "user"},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "_Join:user:post"
	object = types.M{
		"relatedId": "01",
		"owningId":  "2001",
	}
	Adapter.CreateObject(className, relationSchema, object)
	className = "user"
	query = types.M{
		"$relatedTo": types.M{
			"key": "user",
			"object": types.M{
				"className": "post",
				"objectId":  "2001",
			},
		},
	}
	options = nil
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{
		types.M{
			"objectId": "01",
			"key":      3.0,
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"post": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "03",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "_Join:post:user"
	object = types.M{
		"relatedId": "2001",
		"owningId":  "01",
	}
	Adapter.CreateObject(className, relationSchema, object)
	object = types.M{
		"relatedId": "2002",
		"owningId":  "02",
	}
	Adapter.CreateObject(className, relationSchema, object)
	object = types.M{
		"relatedId": "2003",
		"owningId":  "03",
	}
	Adapter.CreateObject(className, relationSchema, object)
	className = "user"
	query = types.M{
		"post": types.M{
			"__type":    "Pointer",
			"className": "post",
			"objectId":  "2001",
		},
	}
	options = nil
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{
		types.M{
			"objectId": "01",
			"key":      "hello",
			"post": types.M{
				"__type":    "Relation",
				"className": "post",
			},
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"post": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "03",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "_Join:post:user"
	object = types.M{
		"relatedId": "2001",
		"owningId":  "01",
	}
	Adapter.CreateObject(className, relationSchema, object)
	object = types.M{
		"relatedId": "2002",
		"owningId":  "02",
	}
	Adapter.CreateObject(className, relationSchema, object)
	object = types.M{
		"relatedId": "2003",
		"owningId":  "03",
	}
	Adapter.CreateObject(className, relationSchema, object)
	className = "user"
	query = types.M{
		"post": types.M{
			"$in": types.S{
				types.M{
					"__type":    "Pointer",
					"className": "post",
					"objectId":  "2001",
				},
			},
		},
	}
	options = nil
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{
		types.M{
			"objectId": "01",
			"key":      "hello",
			"post": types.M{
				"__type":    "Relation",
				"className": "post",
			},
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"post": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "03",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "_Join:post:user"
	object = types.M{
		"relatedId": "2001",
		"owningId":  "01",
	}
	Adapter.CreateObject(className, relationSchema, object)
	object = types.M{
		"relatedId": "2002",
		"owningId":  "02",
	}
	Adapter.CreateObject(className, relationSchema, object)
	object = types.M{
		"relatedId": "2003",
		"owningId":  "03",
	}
	Adapter.CreateObject(className, relationSchema, object)
	className = "user"
	query = types.M{
		"post": types.M{
			"$nin": types.S{
				types.M{
					"__type":    "Pointer",
					"className": "post",
					"objectId":  "2001",
				},
			},
		},
	}
	options = nil
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{
		types.M{
			"objectId": "02",
			"key":      "hello",
			"post": types.M{
				"__type":    "Relation",
				"className": "post",
			},
		},
		types.M{
			"objectId": "03",
			"key":      "hello",
			"post": types.M{
				"__type":    "Relation",
				"className": "post",
			},
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"post": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "03",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "_Join:post:user"
	object = types.M{
		"relatedId": "2001",
		"owningId":  "01",
	}
	Adapter.CreateObject(className, relationSchema, object)
	object = types.M{
		"relatedId": "2002",
		"owningId":  "02",
	}
	Adapter.CreateObject(className, relationSchema, object)
	object = types.M{
		"relatedId": "2003",
		"owningId":  "03",
	}
	Adapter.CreateObject(className, relationSchema, object)
	className = "user"
	query = types.M{
		"post": types.M{
			"$ne": types.M{
				"__type":    "Pointer",
				"className": "post",
				"objectId":  "2001",
			},
		},
	}
	options = nil
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{
		types.M{
			"objectId": "02",
			"key":      "hello",
			"post": types.M{
				"__type":    "Relation",
				"className": "post",
			},
		},
		types.M{
			"objectId": "03",
			"key":      "hello",
			"post": types.M{
				"__type":    "Relation",
				"className": "post",
			},
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"post": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "03",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "_Join:post:user"
	object = types.M{
		"relatedId": "2001",
		"owningId":  "01",
	}
	Adapter.CreateObject(className, relationSchema, object)
	object = types.M{
		"relatedId": "2002",
		"owningId":  "02",
	}
	Adapter.CreateObject(className, relationSchema, object)
	object = types.M{
		"relatedId": "2003",
		"owningId":  "03",
	}
	Adapter.CreateObject(className, relationSchema, object)
	className = "user"
	query = types.M{
		"post": types.M{
			"$eq": types.M{
				"__type":    "Pointer",
				"className": "post",
				"objectId":  "2001",
			},
		},
	}
	options = nil
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{
		types.M{
			"objectId": "01",
			"key":      "hello",
			"post": types.M{
				"__type":    "Relation",
				"className": "post",
			},
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "03",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{
		"@key": "hello",
	}
	options = nil
	results, err = TomatoDBController.Find(className, query, options)
	expectErr = errs.E(errs.InvalidKeyName, "Invalid key name: @key")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "03",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{}
	options = types.M{"count": true}
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{3}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"_wperm":   types.M{"type": "Array"},
			"_rperm":   types.M{"type": "Array"},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
		"_rperm":   types.S{"role:1024"},
		"_wperm":   types.S{"role:1024"},
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{}
	options = nil
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{
		types.M{
			"objectId": "01",
			"key":      "hello",
			"ACL": types.M{
				"role:1024": types.M{
					"read":  true,
					"write": true,
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "_User"
	schema = types.M{
		"fields": types.M{
			"objectId":         types.M{"type": "String"},
			"key":              types.M{"type": "String"},
			"_hashed_password": types.M{"type": "String"},
			"sessionToken":     types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "_User"
	object = types.M{
		"objectId":         "01",
		"key":              "hello",
		"_hashed_password": "123456",
		"sessionToken":     "abcd",
	}
	Adapter.CreateObject(className, schema, object)
	className = "_User"
	query = types.M{}
	options = nil
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{
		types.M{
			"objectId": "01",
			"key":      "hello",
			"password": "123456",
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"_wperm":   types.M{"type": "Array"},
			"_rperm":   types.M{"type": "Array"},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{}
	options = types.M{"acl": []string{"role:1024", "123456789012345678901234"}}
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{
		types.M{
			"objectId": "01",
			"key":      "hello",
		},
		types.M{
			"objectId": "02",
			"key":      "hello",
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"_wperm":   types.M{"type": "Array"},
			"_rperm":   types.M{"type": "Array"},
		},
		"classLevelPermissions": types.M{
			"find": types.M{"role:1024": true},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{}
	options = types.M{"acl": []string{"role:1024", "123456789012345678901234"}}
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{
		types.M{
			"objectId": "01",
			"key":      "hello",
		},
		types.M{
			"objectId": "02",
			"key":      "hello",
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"_wperm":   types.M{"type": "Array"},
			"_rperm":   types.M{"type": "Array"},
		},
		"classLevelPermissions": types.M{
			"find": types.M{"role:2048": true},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{}
	options = types.M{"acl": []string{"role:1024", "123456789012345678901234"}}
	results, err = TomatoDBController.Find(className, query, options)
	expectErr = errs.E(errs.OperationForbidden, "Permission denied for action find on class user.")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"key2":     types.M{"type": "Pointer", "targetClass": "_User"},
			"_wperm":   types.M{"type": "Array"},
			"_rperm":   types.M{"type": "Array"},
		},
		"classLevelPermissions": types.M{
			"find":           types.M{"role:2048": true},
			"readUserFields": types.S{"key2"},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
		"key2": types.M{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  "123456789012345678901234",
		},
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{}
	options = types.M{"acl": []string{"role:1024", "123456789012345678901234"}}
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{
		types.M{
			"objectId": "01",
			"key":      "hello",
			"key2": types.M{
				"__type":    "Pointer",
				"className": "_User",
				"objectId":  "123456789012345678901234",
			},
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"_wperm":   types.M{"type": "Array"},
			"_rperm":   types.M{"type": "Array"},
		},
		"classLevelPermissions": types.M{
			"find": types.M{"role:1024": true},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
		"_rperm":   types.S{"role:1024"},
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{}
	options = types.M{"acl": []string{"role:1024", "123456789012345678901234"}}
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{
		types.M{
			"objectId": "01",
			"key":      "hello",
			"ACL": types.M{
				"role:1024": types.M{
					"read": true,
				},
			},
		},
		types.M{
			"objectId": "02",
			"key":      "hello",
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"_wperm":   types.M{"type": "Array"},
			"_rperm":   types.M{"type": "Array"},
		},
		"classLevelPermissions": types.M{
			"find": types.M{"role:1024": true},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
		"_rperm":   types.S{"role:1024"},
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
		"_rperm":   types.S{"role:2048"},
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{}
	options = types.M{"acl": []string{"role:1024", "123456789012345678901234"}}
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{
		types.M{
			"objectId": "01",
			"key":      "hello",
			"ACL": types.M{
				"role:1024": types.M{
					"read": true,
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "_User"
	schema = types.M{
		"fields": types.M{
			"objectId":         types.M{"type": "String"},
			"key":              types.M{"type": "String"},
			"_hashed_password": types.M{"type": "String"},
			"sessionToken":     types.M{"type": "String"},
			"authData":         types.M{"type": "Object"},
			"_wperm":           types.M{"type": "Array"},
			"_rperm":           types.M{"type": "Array"},
		},
		"classLevelPermissions": types.M{
			"find": types.M{"role:1024": true},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "_User"
	object = types.M{
		"objectId":         "123456789012345678901234",
		"key":              "hello",
		"_rperm":           types.S{"role:1024"},
		"_hashed_password": "123456",
		"sessionToken":     "abcd",
		"authData": types.M{
			"facebook": types.M{"id": "1024"},
		},
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
		"_rperm":   types.S{"role:2048"},
	}
	Adapter.CreateObject(className, schema, object)
	className = "_User"
	query = types.M{}
	options = types.M{"acl": []string{"role:1024", "123456789012345678901234"}}
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{
		types.M{
			"objectId": "123456789012345678901234",
			"key":      "hello",
			"ACL": types.M{
				"role:1024": types.M{
					"read": true,
				},
			},
			"password": "123456",
			"authData": types.M{
				"facebook": map[string]interface{}{"id": "1024"},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "_User"
	schema = types.M{
		"fields": types.M{
			"objectId":         types.M{"type": "String"},
			"key":              types.M{"type": "String"},
			"_hashed_password": types.M{"type": "String"},
			"sessionToken":     types.M{"type": "String"},
			"authData":         types.M{"type": "Object"},
			"_wperm":           types.M{"type": "Array"},
			"_rperm":           types.M{"type": "Array"},
		},
		"classLevelPermissions": types.M{
			"find": types.M{"role:1024": true},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "_User"
	object = types.M{
		"objectId":         "123456789012345678904321",
		"key":              "hello",
		"_rperm":           types.S{"role:1024"},
		"_hashed_password": "123456",
		"sessionToken":     "abcd",
		"authData": types.M{
			"facebook": types.M{"id": "1024"},
		},
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
		"_rperm":   types.S{"role:2048"},
	}
	Adapter.CreateObject(className, schema, object)
	className = "_User"
	query = types.M{}
	options = types.M{"acl": []string{"role:1024", "123456789012345678901234"}}
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{
		types.M{
			"objectId": "123456789012345678904321",
			"key":      "hello",
			"ACL": types.M{
				"role:1024": types.M{
					"read": true,
				},
			},
			"password": "123456",
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
}

func TestPostgres_Destroy(t *testing.T) {
	initPostgresEnv()
	var schema types.M
	var object types.M
	var className string
	var query types.M
	var options types.M
	var err error
	var expectErr error
	var results []types.M
	var expects []types.M
	/*************************************************/
	className = "user"
	query = nil
	options = nil
	err = TomatoDBController.Destroy(className, query, options)
	expectErr = errs.E(errs.ObjectNotFound, "Object not found.")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "1002",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = nil
	options = nil
	err = TomatoDBController.Destroy(className, query, options)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	expects = []types.M{}
	if reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"_wperm":   types.M{"type": "Array"},
			"_rperm":   types.M{"type": "Array"},
		},
		"classLevelPermissions": types.M{
			"delete": types.M{"role:2001": true},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "1002",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = nil
	options = types.M{
		"acl": []string{"123456789012345678901234", "role:1001"},
	}
	err = TomatoDBController.Destroy(className, query, options)
	expectErr = errs.E(errs.OperationForbidden, "Permission denied for action delete on class user.")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"_wperm":   types.M{"type": "Array"},
			"_rperm":   types.M{"type": "Array"},
		},
		"classLevelPermissions": types.M{
			"delete": types.M{"role:1001": true},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "1002",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = nil
	options = types.M{
		"acl": []string{"123456789012345678901234", "role:1001"},
	}
	err = TomatoDBController.Destroy(className, query, options)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	expects = []types.M{}
	if reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"_wperm":   types.M{"type": "Array"},
			"_rperm":   types.M{"type": "Array"},
		},
		"classLevelPermissions": types.M{
			"delete": types.M{"role:1001": true},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"key":      "hello",
		"_wperm":   types.S{"role:1001"},
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "1002",
		"key":      "hello",
		"_wperm":   types.S{"role:2001"},
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = nil
	options = types.M{
		"acl": []string{"123456789012345678901234", "role:1001"},
	}
	err = TomatoDBController.Destroy(className, query, options)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId": "1002",
			"key":      "hello",
			"_wperm":   types.S{"role:2001"},
		},
	}
	if reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"_wperm":   types.M{"type": "Array"},
			"_rperm":   types.M{"type": "Array"},
		},
		"classLevelPermissions": types.M{
			"delete":          types.M{"role:2001": true},
			"writeUserFields": types.S{"key2"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "1002",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = nil
	options = types.M{
		"acl": []string{"role:1001"},
	}
	err = TomatoDBController.Destroy(className, query, options)
	expectErr = errs.E(errs.ObjectNotFound, "Object not found.")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"key2":     types.M{"type": "Pointer", "targetClass": "_User"},
			"_wperm":   types.M{"type": "Array"},
			"_rperm":   types.M{"type": "Array"},
		},
		"classLevelPermissions": types.M{
			"delete":          types.M{"role:2001": true},
			"writeUserFields": types.S{"key2"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"key":      "hello",
		"key2": types.M{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  "123456789012345678901234",
		},
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "1002",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = nil
	options = types.M{
		"acl": []string{"123456789012345678901234", "role:1001"},
	}
	err = TomatoDBController.Destroy(className, query, options)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId": "1002",
			"key":      "hello",
		},
	}
	if reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "_Session"
	query = nil
	options = nil
	err = TomatoDBController.Destroy(className, query, options)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	TomatoDBController.DeleteEverything()
}

func TestPostgres_Update(t *testing.T) {
	initPostgresEnv()
	var schema types.M
	var object types.M
	var className string
	var query types.M
	var update types.M
	var options types.M
	var skipSanitization bool
	var result types.M
	var err error
	var expectErr error
	var expect types.M
	var results []types.M
	var expects []types.M
	/*************************************************/
	className = "user"
	query = nil
	update = nil
	options = nil
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	query = types.M{}
	update = nil
	options = nil
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{"objectId": "01"}
	update = types.M{"key": "haha"}
	options = nil
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId": "01",
			"key":      "haha",
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{"key": "hello"}
	update = types.M{"key": "haha"}
	options = nil
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId": "01",
			"key":      "haha",
		},
		types.M{
			"objectId": "02",
			"key":      "haha",
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{"key": "hello"}
	update = types.M{"key": "haha"}
	options = types.M{
		"many": true,
	}
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId": "01",
			"key":      "haha",
		},
		types.M{
			"objectId": "02",
			"key":      "haha",
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"key": "hello",
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"key": "helloworld",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{"key": "haha"}
	update = types.M{"key": "haha"}
	options = types.M{
		"upsert": true,
	}
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"key": "hello",
		},
		types.M{
			"key": "helloworld",
		},
		types.M{
			"key": "haha",
		},
	}
	if err != nil || len(results) != 3 {
		t.Error("expect:", expects, "result:", results, err)
	} else {
		if reflect.DeepEqual(expects[2], results[2]) == false {
			t.Error("expect:", expects, "result:", results, err)
		}
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"post":     types.M{"type": "Relation", "targetClass": "post"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{"objectId": "01"}
	update = types.M{
		"key": "haha",
		"post": types.M{
			"__op": "AddRelation",
			"objects": types.S{
				types.M{
					"__type":    "Pointer",
					"className": "post",
					"objectId":  "2001",
				},
			},
		},
	}
	options = nil
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId": "01",
			"key":      "haha",
			"post":     types.M{"__type": "Relation", "className": "post"},
		},
	}
	if err != nil || len(results) != 1 {
		t.Error("expect:", expects, "result:", results, err)
	} else {
		if reflect.DeepEqual(expects, results) == false {
			t.Error("expect:", expects, "result:", results, err)
		}
	}
	results, err = Adapter.Find("_Join:post:user", relationSchema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"relatedId": "2001",
			"owningId":  "01",
		},
	}
	if err != nil || len(results) != 1 {
		t.Error("expect:", expects, "result:", results, err)
	} else {
		delete(results[0], "objectId")
		if reflect.DeepEqual(expects, results) == false {
			t.Error("expect:", expects, "result:", results, err)
		}
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{"@key": "01"}
	update = types.M{
		"key": "haha",
	}
	options = nil
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expectErr = errs.E(errs.InvalidKeyName, "Invalid key name: @key")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{"objectId": "01"}
	update = types.M{
		"authData.facebook.id": "haha",
	}
	options = nil
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expectErr = errs.E(errs.InvalidKeyName, "Invalid field name for update: authData.facebook.id")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{"objectId": "01"}
	update = types.M{
		"_abc": "haha",
	}
	options = nil
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expectErr = errs.E(errs.InvalidKeyName, "Invalid field name for update: _abc")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{"objectId": "01"}
	update = types.M{
		"key": types.M{
			"a$b": "hello",
		},
	}
	options = nil
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expectErr = errs.E(errs.InvalidNestedKey, "Nested keys should not contain the '$' or '.' characters")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{"objectId": "01"}
	update = types.M{
		"key": types.M{
			"a.b": "hello",
		},
	}
	options = nil
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expectErr = errs.E(errs.InvalidNestedKey, "Nested keys should not contain the '$' or '.' characters")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"_wperm":   types.M{"type": "Array"},
			"_rperm":   types.M{"type": "Array"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{"objectId": "01"}
	update = types.M{
		"key": "haha",
		"ACL": types.M{
			"role:1024": types.M{
				"read":  true,
				"write": true,
			},
		},
	}
	options = nil
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId": "01",
			"key":      "haha",
			"_rperm":   types.S{"role:1024"},
			"_wperm":   types.S{"role:1024"},
		},
	}
	if err != nil || len(results) != 1 {
		t.Error("expect:", expects, "result:", results, err)
	} else {
		if reflect.DeepEqual(expects, results) == false {
			t.Error("expect:", expects, "result:", results, err)
		}
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"authData": types.M{"type": "Object"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{"objectId": "01"}
	update = types.M{
		"key": "haha",
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	options = nil
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId": "01",
			"key":      "haha",
			"authData": types.M{
				"facebook": map[string]interface{}{
					"id": "1001",
				},
			},
		},
	}
	if err != nil || len(results) != 1 {
		t.Error("expect:", expects, "result:", results, err)
	} else {
		if reflect.DeepEqual(expects, results) == false {
			t.Error("expect:", expects, "result:", results, err)
		}
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{"objectId": "02"}
	update = types.M{
		"key": "haha",
	}
	options = nil
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expectErr = errs.E(errs.ObjectNotFound, "Object not found.")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"key2":     types.M{"type": "Number"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
		"key2":     10,
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{"objectId": "01"}
	update = types.M{
		"key": "haha",
		"key2": types.M{
			"__op":   "Increment",
			"amount": 10,
		},
	}
	options = nil
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expect = types.M{
		"key2": 20.0,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId": "01",
			"key":      "haha",
			"key2":     20.0,
		},
	}
	if err != nil || len(results) != 1 {
		t.Error("expect:", expects, "result:", results, err)
	} else {
		if reflect.DeepEqual(expects, results) == false {
			t.Error("expect:", expects, "result:", results, err)
		}
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"key2":     types.M{"type": "Number"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
		"key2":     10,
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{"objectId": "01"}
	update = types.M{
		"key": "haha",
		"key2": types.M{
			"__op":   "Increment",
			"amount": 10,
		},
	}
	options = nil
	skipSanitization = true
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expect = types.M{
		"objectId": "01",
		"key":      "haha",
		"key2":     20.0,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId": "01",
			"key":      "haha",
			"key2":     20.0,
		},
	}
	if err != nil || len(results) != 1 {
		t.Error("expect:", expects, "result:", results, err)
	} else {
		if reflect.DeepEqual(expects, results) == false {
			t.Error("expect:", expects, "result:", results, err)
		}
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"_wperm":   types.M{"type": "Array"},
			"_rperm":   types.M{"type": "Array"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{"objectId": "01"}
	update = types.M{
		"key": "haha",
	}
	options = types.M{
		"acl": []string{"role:1024"},
	}
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId": "01",
			"key":      "haha",
		},
	}
	if err != nil || len(results) != 1 {
		t.Error("expect:", expects, "result:", results, err)
	} else {
		if reflect.DeepEqual(expects, results) == false {
			t.Error("expect:", expects, "result:", results, err)
		}
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"_wperm":   types.M{"type": "Array"},
			"_rperm":   types.M{"type": "Array"},
		},
		"classLevelPermissions": types.M{
			"update": types.M{"role:1024": true},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{"objectId": "01"}
	update = types.M{
		"key": "haha",
	}
	options = types.M{
		"acl": []string{"role:1024"},
	}
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId": "01",
			"key":      "haha",
		},
	}
	if err != nil || len(results) != 1 {
		t.Error("expect:", expects, "result:", results, err)
	} else {
		if reflect.DeepEqual(expects, results) == false {
			t.Error("expect:", expects, "result:", results, err)
		}
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"_wperm":   types.M{"type": "Array"},
			"_rperm":   types.M{"type": "Array"},
		},
		"classLevelPermissions": types.M{
			"update": types.M{"role:2048": true},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{"objectId": "01"}
	update = types.M{
		"key": "haha",
	}
	options = types.M{
		"acl": []string{"role:1024"},
	}
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expectErr = errs.E(errs.OperationForbidden, "Permission denied for action update on class user.")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"key2":     types.M{"type": "Pointer", "targetClass": "_User"},
			"_wperm":   types.M{"type": "Array"},
			"_rperm":   types.M{"type": "Array"},
		},
		"classLevelPermissions": types.M{
			"update":          types.M{"role:1024": true},
			"writeUserFields": types.S{"key2"},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
		"key2": types.M{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  "123456789012345678901234",
		},
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{"objectId": "01"}
	update = types.M{
		"key": "haha",
	}
	options = types.M{
		"acl": []string{"role:2048", "123456789012345678901234"},
	}
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId": "01",
			"key":      "haha",
			"key2": types.M{
				"__type":    "Pointer",
				"className": "_User",
				"objectId":  "123456789012345678901234",
			},
		},
	}
	if err != nil || len(results) != 1 {
		t.Error("expect:", expects, "result:", results, err)
	} else {
		if reflect.DeepEqual(expects, results) == false {
			t.Error("expect:", expects, "result:", results, err)
		}
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"key2":     types.M{"type": "Pointer", "targetClass": "_User"},
			"_wperm":   types.M{"type": "Array"},
			"_rperm":   types.M{"type": "Array"},
		},
		"classLevelPermissions": types.M{
			"update":          types.M{"role:1024": true},
			"writeUserFields": types.S{"key2"},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
		"key2": types.M{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  "123456789012345678901234",
		},
	}
	Adapter.CreateObject(className, types.M{}, object)
	className = "user"
	query = types.M{"objectId": "01"}
	update = types.M{
		"key": "haha",
	}
	options = types.M{
		"acl": []string{"role:2048"},
	}
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"key2":     types.M{"type": "Pointer", "targetClass": "_User"},
			"_wperm":   types.M{"type": "Array"},
			"_rperm":   types.M{"type": "Array"},
		},
		"classLevelPermissions": types.M{
			"update":          types.M{"role:1024": true},
			"writeUserFields": types.S{"key2"},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
		"key2": types.M{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  "123456789012345678901234",
		},
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{"objectId": "01"}
	update = types.M{
		"key": "haha",
	}
	options = types.M{
		"acl": []string{"role:2048", "123456789012345678900000"},
	}
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expectErr = errs.E(errs.ObjectNotFound, "Object not found.")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
			"_wperm":   types.M{"type": "Array"},
			"_rperm":   types.M{"type": "Array"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
		"_wperm":   types.S{"role:2048"},
	}
	Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
		"_wperm":   types.S{"role:1024"},
	}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	query = types.M{"objectId": "01"}
	update = types.M{
		"key": "haha",
	}
	options = types.M{
		"acl": []string{"role:2048", "123456789012345678901234"},
	}
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId": "02",
			"key":      "hello",
			"_wperm":   types.S{"role:1024"},
		},
		types.M{
			"objectId": "01",
			"key":      "haha",
			"_wperm":   types.S{"role:2048"},
		},
	}
	if err != nil || len(results) != 2 {
		t.Error("expect:", expects, "result:", results, err)
	} else {
		if reflect.DeepEqual(expects, results) == false {
			t.Error("expect:", expects, "result:", results, err)
		}
	}
	TomatoDBController.DeleteEverything()
}

func TestPostgres_Create(t *testing.T) {
	initPostgresEnv()
	var className string
	var schema types.M
	var object types.M
	var options types.M
	var err error
	var expectErr error
	timeStr := utils.TimetoString(time.Now())
	var results []types.M
	var expects []types.M
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = nil
	options = nil
	err = TomatoDBController.Create(className, object, options)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	if len(results) != 0 {
		t.Error("expect:", 0, "result:", len(results))
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId":  types.M{"type": "String"},
			"createdAt": types.M{"type": "Date"},
			"updatedAt": types.M{"type": "Date"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId":  "01",
		"createdAt": timeStr,
		"updatedAt": timeStr,
	}
	options = nil
	err = TomatoDBController.Create(className, object, options)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId":  "01",
			"createdAt": timeStr,
			"updatedAt": timeStr,
		},
	}
	if reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "@user"
	object = nil
	options = nil
	err = TomatoDBController.Create(className, object, options)
	expectErr = errs.E(errs.InvalidClassName, "invalid className: @user")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"create": types.M{"role:1001": true},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{"key": "hello"}
	options = nil
	err = TomatoDBController.Create(className, object, options)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	if len(results) != 1 {
		t.Error("expect:", 1, "result:", len(results))
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"create": types.M{"role:1001": true},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = types.M{"key": "hello"}
	options = types.M{
		"acl": []string{"role:1001"},
	}
	err = TomatoDBController.Create(className, object, options)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	if len(results) != 1 {
		t.Error("expect:", 1, "result:", len(results))
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"create": types.M{"role:1001": true},
		},
	}
	Adapter.CreateClass(className, schema)
	className = "user"
	object = nil
	options = types.M{
		"acl": []string{"role:2001"},
	}
	err = TomatoDBController.Create(className, object, options)
	expectErr = errs.E(errs.OperationForbidden, "Permission denied for action create on class user.")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId":  types.M{"type": "String"},
			"createdAt": types.M{"type": "Date"},
			"key":       types.M{"type": "Relation", "targetClass": "post"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId":  "1001",
		"createdAt": timeStr,
		"key": types.M{
			"__op": "AddRelation",
			"objects": types.S{
				types.M{
					"__type":    "Pointer",
					"className": "post",
					"objectId":  "2001",
				},
			},
		},
	}
	options = nil
	err = TomatoDBController.Create(className, object, options)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId":  "1001",
			"createdAt": timeStr,
			"key":       types.M{"__type": "Relation", "className": "post"},
		},
	}
	if reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results)
	}
	results, err = Adapter.Find("_Join:key:user", relationSchema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"relatedId": "2001",
			"owningId":  "1001",
		},
	}
	delete(results[0], "objectId")
	if reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "_User"
	schema = types.M{
		"fields": types.M{
			"objectId":  types.M{"type": "String"},
			"createdAt": types.M{"type": "Date"},
			"authData":  types.M{"type": "Object"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId":  "1001",
		"createdAt": timeStr,
		"authData": types.M{
			"facebook": types.M{
				"id": "1024",
			},
		},
	}
	options = nil
	err = TomatoDBController.Create(className, object, options)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId":  "1001",
			"createdAt": timeStr,
			"authData": types.M{
				"facebook": map[string]interface{}{
					"id": "1024",
				},
			},
		},
	}
	if reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId":  types.M{"type": "String"},
			"createdAt": types.M{"type": "Date"},
			"key":       types.M{"type": "Number"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId":  "1001",
		"createdAt": timeStr,
		"key": types.M{
			"__op":   "Increment",
			"amount": 10,
		},
	}
	options = nil
	err = TomatoDBController.Create(className, object, options)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId":  "1001",
			"createdAt": timeStr,
			"key":       10.0,
		},
	}
	if reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId":  types.M{"type": "String"},
			"createdAt": types.M{"type": "Date"},
			"_rperm":    types.M{"type": "Array"},
			"_wperm":    types.M{"type": "Array"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId":  "1001",
		"createdAt": timeStr,
		"ACL": types.M{
			"role:1001": types.M{
				"read":  true,
				"write": true,
			},
		},
	}
	options = nil
	err = TomatoDBController.Create(className, object, options)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	results, err = Adapter.Find(className, schema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId":  "1001",
			"createdAt": timeStr,
			"_rperm":    types.S{"role:1001"},
			"_wperm":    types.S{"role:1001"},
		},
	}
	if reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results)
	}
	TomatoDBController.DeleteEverything()
}

func TestPostgres_validateClassName(t *testing.T) {
	initPostgresEnv()
	var className string
	var err error
	var expect error
	/*************************************************/
	className = "user"
	err = TomatoDBController.validateClassName(className)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/*************************************************/
	className = "@user"
	err = TomatoDBController.validateClassName(className)
	expect = errs.E(errs.InvalidClassName, "invalid className: @user")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
}

func TestPostgres_handleRelationUpdates(t *testing.T) {
	initPostgresEnv()
	var className string
	var objectID string
	var update types.M
	var err error
	var expectErr error
	var expect types.M
	var object types.M
	var results []types.M
	var expects []types.M
	/*************************************************/
	className = "user"
	objectID = "1001"
	update = nil
	err = TomatoDBController.handleRelationUpdates(className, objectID, update)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	className = "user"
	objectID = "1001"
	update = types.M{"key": "hello"}
	err = TomatoDBController.handleRelationUpdates(className, objectID, update)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	expect = types.M{"key": "hello"}
	if reflect.DeepEqual(expect, update) == false {
		t.Error("expect:", expect, "result:", update)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"fields": types.M{
			"key":  types.M{"type": "String"},
			"post": types.M{"type": "Relation"},
		},
	}
	Adapter.CreateClass("user", object)
	className = "user"
	objectID = "1001"
	update = types.M{
		"key": "hello",
		"post": types.M{
			"__op": "AddRelation",
			"objects": types.S{
				types.M{
					"__type":    "Pointer",
					"className": "post",
					"objectId":  "2001",
				},
				types.M{
					"__type":    "Pointer",
					"className": "post",
					"objectId":  "2002",
				},
			},
		},
	}
	err = TomatoDBController.handleRelationUpdates(className, objectID, update)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	expect = types.M{"key": "hello"}
	if reflect.DeepEqual(expect, update) == false {
		t.Error("expect:", expect, "result:", update)
	}
	results, err = Adapter.Find("_Join:post:user", relationSchema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"relatedId": "2001",
			"owningId":  "1001",
		},
		types.M{
			"relatedId": "2002",
			"owningId":  "1001",
		},
	}
	if err != nil || len(results) != 2 {
		t.Error("expect:", expects, "result:", results, err)
	} else {
		delete(results[0], "objectId")
		delete(results[1], "objectId")
		if reflect.DeepEqual(expects, results) == false {
			t.Error("expect:", expects, "result:", results, err)
		}
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"fields": types.M{
			"key":  types.M{"type": "String"},
			"post": types.M{"type": "Relation"},
		},
	}
	Adapter.CreateClass("user", object)
	object = types.M{
		"relatedId": "2001",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:post:user", relationSchema, object)
	object = types.M{
		"relatedId": "2002",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:post:user", relationSchema, object)
	className = "user"
	objectID = "1001"
	update = types.M{
		"key": "hello",
		"post": types.M{
			"__op": "RemoveRelation",
			"objects": types.S{
				types.M{
					"__type":    "Pointer",
					"className": "post",
					"objectId":  "2001",
				},
			},
		},
	}
	err = TomatoDBController.handleRelationUpdates(className, objectID, update)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	expect = types.M{"key": "hello"}
	if reflect.DeepEqual(expect, update) == false {
		t.Error("expect:", expect, "result:", update)
	}
	results, err = Adapter.Find("_Join:post:user", relationSchema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"relatedId": "2002",
			"owningId":  "1001",
		},
	}
	if err != nil || len(results) != 1 {
		t.Error("expect:", expects, "result:", results, err)
	} else {
		if reflect.DeepEqual(expects, results) == false {
			t.Error("expect:", expects, "result:", results, err)
		}
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"fields": types.M{
			"key":  types.M{"type": "String"},
			"post": types.M{"type": "Relation"},
		},
	}
	Adapter.CreateClass("user", object)
	object = types.M{
		"relatedId": "2001",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:post:user", relationSchema, object)
	object = types.M{
		"relatedId": "2002",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:post:user", relationSchema, object)
	className = "user"
	objectID = "1001"
	update = types.M{
		"key": "hello",
		"post": types.M{
			"__op": "Batch",
			"ops": types.S{
				types.M{
					"__op": "RemoveRelation",
					"objects": types.S{
						types.M{
							"__type":    "Pointer",
							"className": "post",
							"objectId":  "2001",
						},
					},
				},
				types.M{
					"__op": "AddRelation",
					"objects": types.S{
						types.M{
							"__type":    "Pointer",
							"className": "post",
							"objectId":  "2003",
						},
					},
				},
			},
		},
	}
	err = TomatoDBController.handleRelationUpdates(className, objectID, update)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	expect = types.M{"key": "hello"}
	if reflect.DeepEqual(expect, update) == false {
		t.Error("expect:", expect, "result:", update)
	}
	results, err = Adapter.Find("_Join:post:user", relationSchema, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"relatedId": "2002",
			"owningId":  "1001",
		},
		types.M{
			"relatedId": "2003",
			"owningId":  "1001",
		},
	}
	if err != nil || len(results) != 2 {
		t.Error("expect:", expects, "result:", results, err)
	} else {
		delete(results[0], "objectId")
		delete(results[1], "objectId")
		if reflect.DeepEqual(expects, results) == false {
			t.Error("expect:", expects, "result:", results, err)
		}
	}
	Adapter.DeleteAllClasses()
}

func TestPostgres_addRelation(t *testing.T) {
	initPostgresEnv()
	var object types.M
	var key, fromClassName, fromID, toID string
	var err error
	var expect error
	var results []types.M
	var expectRes types.M
	/*************************************************/
	object = types.M{
		"fields": types.M{
			"key":  types.M{"type": "String"},
			"post": types.M{"type": "Relation"},
		},
	}
	Adapter.CreateClass("user", object)
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
		"fields": types.M{
			"key":  types.M{"type": "String"},
			"post": types.M{"type": "Relation"},
		},
	}
	Adapter.CreateClass("user", object)
	object = types.M{
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
		"relatedId": "2001",
		"owningId":  "1001",
	}
	if len(results) != 1 || reflect.DeepEqual(expectRes, results[0]) == false {
		t.Error("expect:", expectRes, "result:", results[0])
	}
	Adapter.DeleteAllClasses()
}

func TestPostgres_removeRelation(t *testing.T) {
	initPostgresEnv()
	var object types.M
	var key, fromClassName, fromID, toID string
	var err error
	var expect error
	var results []types.M
	/*************************************************/
	object = types.M{
		"fields": types.M{
			"key":  types.M{"type": "String"},
			"post": types.M{"type": "Relation"},
		},
	}
	Adapter.CreateClass("user", object)
	key = "post"
	fromClassName = "user"
	fromID = "1001"
	toID = "2001"
	err = TomatoDBController.removeRelation(key, fromClassName, fromID, toID)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	results, err = Adapter.Find("_Join:post:user", relationSchema, types.M{}, types.M{})
	if err != nil || len(results) != 0 {
		t.Error("expect:", nil, "result:", results, err)
	}
	Adapter.DeleteAllClasses()
	/*************************************************/
	object = types.M{
		"fields": types.M{
			"key":  types.M{"type": "String"},
			"post": types.M{"type": "Relation"},
		},
	}
	Adapter.CreateClass("user", object)
	object = types.M{
		"relatedId": "2001",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:post:user", relationSchema, object)
	key = "post"
	fromClassName = "user"
	fromID = "1001"
	toID = "2001"
	err = TomatoDBController.removeRelation(key, fromClassName, fromID, toID)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	results, err = Adapter.Find("_Join:post:user", relationSchema, types.M{}, types.M{})
	if err != nil || len(results) != 0 {
		t.Error("expect:", nil, "result:", results, err)
	}
	Adapter.DeleteAllClasses()
}

func TestPostgres_ValidateObject(t *testing.T) {
	initPostgresEnv()
	var className string
	var object types.M
	var query types.M
	var options types.M
	var err error
	var expect error
	/*************************************************/
	className = "user"
	object = types.M{}
	query = types.M{}
	options = types.M{}
	err = TomatoDBController.ValidateObject(className, object, query, options)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"addField": types.M{"role:1001": true},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{"key1": "hello"}
	query = types.M{}
	options = types.M{}
	err = TomatoDBController.ValidateObject(className, object, query, options)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"addField": types.M{"role:1001": true},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{"key1": "hello"}
	query = types.M{}
	options = types.M{
		"acl": []string{"2001"},
	}
	err = TomatoDBController.ValidateObject(className, object, query, options)
	expect = errs.E(errs.OperationForbidden, "Permission denied for action addField on class user.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"addField": types.M{"role:1001": true},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{"key1": "hello"}
	query = types.M{}
	options = types.M{
		"acl": []string{"role:1001"},
	}
	err = TomatoDBController.ValidateObject(className, object, query, options)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	TomatoDBController.DeleteEverything()
}

func TestPostgres_LoadSchema(t *testing.T) {
	initPostgresEnv()
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
	result = TomatoDBController.LoadSchema(nil)
	expect = types.M{
		"key":       map[string]interface{}{"type": "String"},
		"objectId":  types.M{"type": "String"},
		"createdAt": types.M{"type": "Date"},
		"updatedAt": types.M{"type": "Date"},
		"ACL":       types.M{"type": "ACL"},
	}
	if reflect.DeepEqual(expect, result.data["user"]) == false {
		t.Error("expect:", expect, "result:", result.data["user"])
	}
	expect = types.M{
		"get":      map[string]interface{}{"*": true},
		"find":     types.M{"*": true},
		"create":   types.M{"*": true},
		"update":   types.M{"*": true},
		"delete":   types.M{"*": true},
		"addField": types.M{"*": true},
	}
	if reflect.DeepEqual(expect, result.perms["user"]) == false {
		t.Error("expect:", expect, "result:", result.perms["user"])
	}
	Adapter.DeleteAllClasses()
}

func TestPostgres_DeleteEverything(t *testing.T) {
	// 测试用例与 Adapter.DeleteAllClasses 类似
}

func TestPostgres_RedirectClassNameForKey(t *testing.T) {
	initPostgresEnv()
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
	initPostgresEnv()
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

func TestPostgres_canAddField(t *testing.T) {
	initPostgresEnv()
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
	schema = getPostgresSchema()
	schema.reloadData(nil)
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
	schema = getPostgresSchema()
	schema.reloadData(nil)
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
	schema = getPostgresSchema()
	schema.reloadData(nil)
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
			"get":      types.M{"*": true},
			"addField": types.M{},
		},
	}
	Adapter.CreateClass("user", object)
	schema = getPostgresSchema()
	schema.reloadData(nil)
	className = "user"
	object = types.M{
		"key":  "hello",
		"key1": "hello",
	}
	acl = nil
	err = TomatoDBController.canAddField(schema, className, object, acl)
	expect = errs.E(errs.OperationForbidden, "Permission denied for action addField on class user.")
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
	schema = getPostgresSchema()
	schema.reloadData(nil)
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
	schema = getPostgresSchema()
	schema.reloadData(nil)
	className = "user"
	object = types.M{
		"key":  "hello",
		"key1": "hello",
	}
	acl = []string{"role:1024"}
	err = TomatoDBController.canAddField(schema, className, object, acl)
	expect = errs.E(errs.OperationForbidden, "Permission denied for action addField on class user.")
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
	schema = getPostgresSchema()
	schema.reloadData(nil)
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

func TestPostgres_reduceRelationKeys(t *testing.T) {
	initPostgresEnv()
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
		"fields": types.M{
			"key": types.M{"type": "Relation"},
		},
	}
	Adapter.CreateClass("Post", object)
	object = types.M{
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
		"fields": types.M{
			"key": types.M{"type": "Relation"},
		},
	}
	Adapter.CreateClass("Post", object)
	object = types.M{
		"relatedId": "2001",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:key:Post", relationSchema, object)
	object = types.M{
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

func TestPostgres_relatedIds(t *testing.T) {
	initPostgresEnv()
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
		"fields": types.M{
			"key":  types.M{"type": "String"},
			"name": types.M{"type": "Relation"},
		},
	}
	Adapter.CreateClass("user", object)
	object = types.M{
		"relatedId": "01",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:name:user", relationSchema, object)
	object = types.M{
		"relatedId": "02",
		"owningId":  "1002",
	}
	Adapter.CreateObject("_Join:name:user", relationSchema, object)
	object = types.M{
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

func TestPostgres_addInObjectIdsIds(t *testing.T) {
	initPostgresEnv()
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

func TestPostgres_addNotInObjectIdsIds(t *testing.T) {
	initPostgresEnv()
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
	/*************************************************/
	ids = types.S{"1024"}
	query = types.M{"objectId": "2048"}
	result = TomatoDBController.addNotInObjectIdsIds(ids, query)
	expect = types.M{
		"objectId": types.M{
			"$eq":  "2048",
			"$nin": types.S{"1024"},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func TestPostgres_reduceInRelation(t *testing.T) {
	initPostgresEnv()
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
	schema = getPostgresSchema()
	schema.reloadData(nil)
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
	schema = getPostgresSchema()
	schema.reloadData(nil)
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
	schema = getPostgresSchema()
	schema.reloadData(nil)
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
	schema = getPostgresSchema()
	schema.reloadData(nil)
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
	schema = getPostgresSchema()
	schema.reloadData(nil)
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
	schema = getPostgresSchema()
	schema.reloadData(nil)
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
	schema = getPostgresSchema()
	schema.reloadData(nil)
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
	schema = getPostgresSchema()
	schema.reloadData(nil)
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
		"relatedId": "2001",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:key:user", relationSchema, object)
	object = types.M{
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
	schema = getPostgresSchema()
	schema.reloadData(nil)
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
		"relatedId": "2001",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:key:user", relationSchema, object)
	object = types.M{
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
	schema = getPostgresSchema()
	schema.reloadData(nil)
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

func TestPostgres_owningIds(t *testing.T) {
	initPostgresEnv()
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
		"fields": types.M{
			"key":  types.M{"type": "String"},
			"name": types.M{"type": "Relation"},
		},
	}
	Adapter.CreateClass("user", object)
	object = types.M{
		"relatedId": "01",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:name:user", relationSchema, object)
	object = types.M{
		"relatedId": "02",
		"owningId":  "1002",
	}
	Adapter.CreateObject("_Join:name:user", relationSchema, object)
	object = types.M{
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

func TestPostgres_DeleteSchema(t *testing.T) {
	initPostgresEnv()
	var schema types.M
	var object types.M
	var className string
	var err error
	var expectErr error
	/*************************************************/
	className = "user"
	err = TomatoDBController.DeleteSchema(className)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{"key": "hello"}
	Adapter.CreateObject(className, schema, object)
	className = "user"
	err = TomatoDBController.DeleteSchema(className)
	expectErr = errs.E(errs.ClassNotEmpty, "Class user is not empty, contains 1 objects, cannot drop schema.")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{"key": "hello"}
	Adapter.CreateObject(className, schema, object)
	Adapter.DeleteObjectsByQuery(className, schema, types.M{"key": "hello"})
	className = "user"
	err = TomatoDBController.DeleteSchema(className)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	if Adapter.ClassExists(className) == true {
		t.Error("expect:", false, "result:", true)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{"key": "hello"}
	Adapter.CreateObject(className, schema, object)
	Adapter.DeleteObjectsByQuery(className, schema, types.M{"key": "hello"})
	className = "user"
	err = TomatoDBController.DeleteSchema(className)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	if Adapter.ClassExists(className) == true {
		t.Error("expect:", false, "result:", true)
	}
	object, err = Adapter.GetClass(className)
	if err != nil || (object != nil && len(object) != 0) {
		t.Error("expect:", nil, "result:", object, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
			"key1": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
		},
	}
	Adapter.CreateClass(className, schema)
	object = types.M{"key": "hello"}
	Adapter.CreateObject(className, schema, object)
	Adapter.DeleteObjectsByQuery(className, schema, types.M{"key": "hello"})
	object = types.M{
		"relatedId": "2001",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:key1:user", types.M{}, object)
	className = "user"
	err = TomatoDBController.DeleteSchema(className)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	if Adapter.ClassExists(className) == true {
		t.Error("expect:", false, "result:", true)
	}
	object, err = Adapter.GetClass(className)
	if err != nil || (object != nil && len(object) != 0) {
		t.Error("expect:", nil, "result:", object, err)
	}
	if Adapter.ClassExists("_Join:key1:user") == true {
		t.Error("expect:", false, "result:", true)
	}
	TomatoDBController.DeleteEverything()
}

func TestPostgres_addPointerPermissions(t *testing.T) {
	initPostgresEnv()
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
	schema = getPostgresSchema()
	schema.reloadData(nil)
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
	schema = getPostgresSchema()
	schema.reloadData(nil)
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
	schema = getPostgresSchema()
	schema.reloadData(nil)
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
	schema = getPostgresSchema()
	schema.reloadData(nil)
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
	schema = getPostgresSchema()
	schema.reloadData(nil)
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
	schema = getPostgresSchema()
	schema.reloadData(nil)
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
	schema = getPostgresSchema()
	schema.reloadData(nil)
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

func initPostgresEnv() {
	Adapter = getPostgresAdapter()
	schemaCache = cache.NewSchemaCache(5, false)
	TomatoDBController = &DBController{}
}
