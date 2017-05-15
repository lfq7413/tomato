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
	initEnv()
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
	className = "user"
	object = types.M{"key": "001"}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{"key": "002"}
	Adapter.CreateObject(className, types.M{}, object)
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

func Test_Find(t *testing.T) {
	initEnv()
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
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	className = "user"
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
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	className = "user"
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "03",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "03",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "Number"},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      3,
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      1,
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "03",
		"key":      2,
	}
	Adapter.CreateObject(className, types.M{}, object)
	className = "user"
	query = types.M{}
	options = types.M{"sort": []string{"key"}}
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{
		types.M{
			"objectId": "02",
			"key":      1,
		},
		types.M{
			"objectId": "03",
			"key":      2,
		},
		types.M{
			"objectId": "01",
			"key":      3,
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "Number"},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      3,
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      1,
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "03",
		"key":      2,
	}
	Adapter.CreateObject(className, types.M{}, object)
	className = "user"
	query = types.M{}
	options = types.M{"sort": []string{"-key"}}
	results, err = TomatoDBController.Find(className, query, options)
	expects = types.S{
		types.M{
			"objectId": "01",
			"key":      3,
		},
		types.M{
			"objectId": "03",
			"key":      2,
		},
		types.M{
			"objectId": "02",
			"key":      1,
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "Number"},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      3,
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      1,
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "03",
		"key":      2,
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "Number"},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      3,
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      1,
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "03",
		"key":      2,
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "Number"},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      3,
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      1,
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "03",
		"key":      2,
	}
	Adapter.CreateObject(className, types.M{}, object)
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
			"key":      3,
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	initEnv()
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
			"post": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "03",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
			"post": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "03",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
			"post": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "03",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
			"post": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "03",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
			"post": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "03",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "03",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "03",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	initEnv()
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
		"_rperm":   types.S{"role:1024"},
		"_wperm":   types.S{"role:1024"},
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, object)
	className = "_User"
	object = types.M{
		"objectId":         "01",
		"key":              "hello",
		"_hashed_password": "123456",
		"sessionToken":     "abcd",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"find": types.M{"role:1024": true},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	initEnv()
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"find": types.M{"role:2048": true},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	initEnv()
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"find":           types.M{"role:2048": true},
			"readUserFields": types.S{"key2"},
		},
	}
	Adapter.CreateClass(className, object)
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
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	className = "user"
	query = types.M{}
	options = types.M{"acl": []string{"role:1024", "123456789012345678901234"}}
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
	initEnv()
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"find": types.M{"role:1024": true},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
		"_rperm":   types.S{"role:1024"},
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"find": types.M{"role:1024": true},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
		"_rperm":   types.S{"role:1024"},
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
		"_rperm":   types.S{"role:2048"},
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"find": types.M{"role:1024": true},
		},
	}
	Adapter.CreateClass(className, object)
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
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
		"_rperm":   types.S{"role:2048"},
	}
	Adapter.CreateObject(className, types.M{}, object)
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
				"facebook": types.M{"id": "1024"},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results, err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "_User"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"find": types.M{"role:1024": true},
		},
	}
	Adapter.CreateClass(className, object)
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
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
		"_rperm":   types.S{"role:2048"},
	}
	Adapter.CreateObject(className, types.M{}, object)
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

func Test_Destroy(t *testing.T) {
	initEnv()
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
	object = types.M{
		"objectId": "1001",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "1002",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	className = "user"
	query = nil
	options = nil
	err = TomatoDBController.Destroy(className, query, options)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
	expects = []types.M{}
	if reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	object = types.M{
		"objectId": "1001",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "1002",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"delete": types.M{"role:2001": true},
		},
	}
	Adapter.CreateClass(className, object)
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
	initEnv()
	className = "user"
	object = types.M{
		"objectId": "1001",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "1002",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"delete": types.M{"role:1001": true},
		},
	}
	Adapter.CreateClass(className, object)
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
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
	expects = []types.M{}
	if reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	object = types.M{
		"objectId": "1001",
		"key":      "hello",
		"_wperm":   types.S{"role:1001"},
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "1002",
		"key":      "hello",
		"_wperm":   types.S{"role:2001"},
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"delete": types.M{"role:1001": true},
		},
	}
	Adapter.CreateClass(className, object)
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
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId": "1002",
			"key":      "hello",
			"_wperm":   []interface{}{"role:2001"},
		},
	}
	if reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	initEnv()
	className = "user"
	object = types.M{
		"objectId": "1001",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "1002",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"delete":          types.M{"role:2001": true},
			"writeUserFields": types.S{"key2"},
		},
	}
	Adapter.CreateClass(className, object)
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
	object = types.M{
		"objectId": "1001",
		"key":      "hello",
		"key2": types.M{
			"__type":    "Pointer",
			"className": "_User",
			"objectId":  "123456789012345678901234",
		},
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "1002",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"delete":          types.M{"role:2001": true},
			"writeUserFields": types.S{"key2"},
		},
	}
	Adapter.CreateClass(className, object)
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
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
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

func Test_Update(t *testing.T) {
	initEnv()
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
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
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
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId": "01",
			"key":      "haha",
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
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
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
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId": "01",
			"key":      "hello",
		},
		types.M{
			"objectId": "02",
			"key":      "hello",
		},
		types.M{
			"key": "haha",
		},
	}
	if err != nil || len(results) != 3 {
		t.Error("expect:", expects, "result:", results, err)
	} else {
		delete(results[2], "objectId")
		if reflect.DeepEqual(expects[2], results[2]) == false {
			t.Error("expect:", expects, "result:", results, err)
		}
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId": "01",
			"key":      "haha",
		},
		types.M{
			"objectId": "02",
			"key":      "hello",
		},
	}
	if err != nil || len(results) != 2 {
		t.Error("expect:", expects, "result:", results, err)
	} else {
		if reflect.DeepEqual(expects, results) == false {
			t.Error("expect:", expects, "result:", results, err)
		}
	}
	results, err = Adapter.Find("_Join:post:user", types.M{}, types.M{}, types.M{})
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
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId": "01",
			"key":      "haha",
			"_rperm":   []interface{}{"role:1024"},
			"_wperm":   []interface{}{"role:1024"},
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
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId": "01",
			"key":      "haha",
			"authData": types.M{
				"facebook": types.M{
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
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"objectId": "01",
		"key":      "hello",
		"key2":     10,
	}
	Adapter.CreateObject(className, types.M{}, object)
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
		"key2": 20,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId": "01",
			"key":      "haha",
			"key2":     20,
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
	object = types.M{
		"objectId": "01",
		"key":      "hello",
		"key2":     10,
	}
	Adapter.CreateObject(className, types.M{}, object)
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
		"key2":     20,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId": "01",
			"key":      "haha",
			"key2":     20,
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
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
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
	object = types.M{
		"fields": types.M{
			"key": "String",
		},
		"classLevelPermissions": types.M{
			"update": types.M{"role:1024": true},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
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
	initEnv()
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": "String",
		},
		"classLevelPermissions": types.M{
			"update": types.M{"role:2048": true},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"update":          types.M{"role:1024": true},
			"writeUserFields": types.S{"key2"},
		},
	}
	Adapter.CreateClass(className, object)
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
		"acl": []string{"role:2048", "123456789012345678901234"},
	}
	skipSanitization = false
	result, err = TomatoDBController.Update(className, query, update, options, skipSanitization)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"update":          types.M{"role:1024": true},
			"writeUserFields": types.S{"key2"},
		},
	}
	Adapter.CreateClass(className, object)
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
	initEnv()
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"update":          types.M{"role:1024": true},
			"writeUserFields": types.S{"key2"},
		},
	}
	Adapter.CreateClass(className, object)
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
	initEnv()
	className = "user"
	object = types.M{
		"objectId": "01",
		"key":      "hello",
		"_wperm":   types.S{"role:2048"},
	}
	Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
		"_wperm":   types.S{"role:1024"},
	}
	Adapter.CreateObject(className, types.M{}, object)
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
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId": "01",
			"key":      "haha",
			"_wperm":   []interface{}{"role:2048"},
		},
		types.M{
			"objectId": "02",
			"key":      "hello",
			"_wperm":   []interface{}{"role:1024"},
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

func Test_Create(t *testing.T) {
	initEnv()
	var className string
	var object types.M
	var options types.M
	var err error
	var expectErr error
	timeStr := utils.TimetoString(time.Now())
	var results []types.M
	var expects []types.M
	/*************************************************/
	className = "user"
	object = nil
	options = nil
	err = TomatoDBController.Create(className, object, options)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
	if len(results) != 1 {
		t.Error("expect:", 1, "result:", len(results))
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
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
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"create": types.M{"role:1001": true},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = nil
	options = nil
	err = TomatoDBController.Create(className, object, options)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
	if len(results) != 1 {
		t.Error("expect:", 1, "result:", len(results))
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	initEnv()
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"create": types.M{"role:1001": true},
		},
	}
	Adapter.CreateClass(className, object)
	className = "user"
	object = nil
	options = types.M{
		"acl": []string{"role:1001"},
	}
	err = TomatoDBController.Create(className, object, options)
	expectErr = nil
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
	if len(results) != 1 {
		t.Error("expect:", 1, "result:", len(results))
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
		"classLevelPermissions": types.M{
			"create": types.M{"role:1001": true},
		},
	}
	Adapter.CreateClass(className, object)
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
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId":  "1001",
			"createdAt": timeStr,
		},
	}
	if reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results)
	}
	results, err = Adapter.Find("_Join:key:user", types.M{}, types.M{}, types.M{})
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
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId":  "1001",
			"createdAt": timeStr,
			"authData": types.M{
				"facebook": types.M{
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
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId":  "1001",
			"createdAt": timeStr,
			"key":       10,
		},
	}
	if reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
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
	results, err = Adapter.Find(className, types.M{}, types.M{}, types.M{})
	expects = []types.M{
		types.M{
			"objectId":  "1001",
			"createdAt": timeStr,
			"_rperm":    []interface{}{"role:1001"},
			"_wperm":    []interface{}{"role:1001"},
		},
	}
	if reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results)
	}
	TomatoDBController.DeleteEverything()
}

func Test_validateClassName(t *testing.T) {
	initEnv()
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

func Test_handleRelationUpdates(t *testing.T) {
	initEnv()
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
		"_id":       "1",
		"relatedId": "2001",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:post:user", relationSchema, object)
	object = types.M{
		"_id":       "2",
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
			"objectId":  "2",
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
		"_id":       "1",
		"relatedId": "2001",
		"owningId":  "1001",
	}
	Adapter.CreateObject("_Join:post:user", relationSchema, object)
	object = types.M{
		"_id":       "2",
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
	initEnv()
	var object types.M
	var key, fromClassName, fromID, toID string
	var err error
	var expect error
	var results []types.M
	/*************************************************/
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
		"objectId":  "01",
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

func Test_ValidateObject(t *testing.T) {
	initEnv()
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
	result = TomatoDBController.LoadSchema(nil)
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

func Test_DeleteEverything(t *testing.T) {
	// 测试用例与 Adapter.DeleteAllClasses 类似
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
	initEnv()
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
	schema = getSchema()
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
	schema = getSchema()
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
	schema = getSchema()
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
	schema = getSchema()
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
	schema = getSchema()
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
	schema = getSchema()
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
	initEnv()
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
	object = types.M{"key": "hello"}
	Adapter.CreateObject(className, types.M{}, object)
	className = "user"
	err = TomatoDBController.DeleteSchema(className)
	expectErr = errs.E(errs.ClassNotEmpty, "Class user is not empty, contains 1 objects, cannot drop schema.")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	TomatoDBController.DeleteEverything()
	/*************************************************/
	className = "user"
	object = types.M{"key": "hello"}
	Adapter.CreateObject(className, types.M{}, object)
	Adapter.DeleteObjectsByQuery(className, types.M{}, types.M{"key": "hello"})
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	Adapter.CreateClass(className, object)
	object = types.M{"key": "hello"}
	Adapter.CreateObject(className, types.M{}, object)
	Adapter.DeleteObjectsByQuery(className, types.M{}, types.M{"key": "hello"})
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
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
			"key1": types.M{
				"type":        "Relation",
				"targetClass": "post",
			},
		},
	}
	Adapter.CreateClass(className, object)
	object = types.M{"key": "hello"}
	Adapter.CreateObject(className, types.M{}, object)
	Adapter.DeleteObjectsByQuery(className, types.M{}, types.M{"key": "hello"})
	object = types.M{
		"_id":       "01",
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
	schema = getSchema()
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
	schema = getSchema()
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
	schema = getSchema()
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
	schema = getSchema()
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
	schema = getSchema()
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
	schema = getSchema()
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
	/*************************************************/
	originalObject = types.M{
		"key": types.M{
			"__op":   "Increment",
			"amount": 10,
		},
		"key2": "hello",
	}
	object = types.M{
		"key2": "hello",
	}
	result = sanitizeDatabaseResult(originalObject, object)
	expect = types.M{}
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
		"objectId":                       "1024",
		"_hashed_password":               "1024",
		"_email_verify_token":            "abc",
		"_perishable_token":              "abc",
		"_perishable_token_expires_at":   "abc",
		"_password_changed_at":           "abc",
		"_tombstone":                     "abc",
		"_email_verify_token_expires_at": "abc",
		"_failed_login_count":            "abc",
		"_account_lockout_expires_at":    "abc",
		"sessionToken":                   "abc",
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
	/*************************************************/
	query = types.M{
		"$or": types.S{
			types.M{"key": "hello"},
			types.M{"key1": "hello"},
		},
	}
	acl = []string{"role:1024"}
	result = addReadACL(query, acl)
	expect = types.M{
		"$or": types.S{
			types.M{"key": "hello"},
			types.M{"key1": "hello"},
		},
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
	var expectQuery types.M
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
		"_rperm":                         "hello",
		"_wperm":                         "hello",
		"_perishable_token":              "hello",
		"_email_verify_token":            "hello",
		"_email_verify_token_expires_at": "hello",
		"_account_lockout_expires_at":    "hello",
		"_failed_login_count":            "hello",
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
		"$or": types.S{
			types.M{"key": "value"},
			types.M{"key": "value"},
		},
		"_rperm": types.M{
			"$in": types.S{nil, "*", "role:1024"},
		},
	}
	err = validateQuery(query)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	expectQuery = types.M{
		"$or": types.S{
			types.M{
				"key": "value",
				"_rperm": types.M{
					"$in": types.S{nil, "*", "role:1024"},
				},
			},
			types.M{
				"key": "value",
				"_rperm": types.M{
					"$in": types.S{nil, "*", "role:1024"},
				},
			},
		},
	}
	if reflect.DeepEqual(query, expectQuery) == false {
		t.Error("expect:", expectQuery, "result:", query)
	}
	/*************************************************/
	query = types.M{
		"$or": types.S{
			types.M{"key": "value"},
			types.M{"key": "value"},
		},
		"loc": types.M{
			"$nearSphere": types.M{
				"__type":    "GeoPoint",
				"latitude":  40,
				"longitude": -30,
			},
		},
	}
	err = validateQuery(query)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	expectQuery = types.M{
		"$or": types.S{
			types.M{"key": "value"},
			types.M{"key": "value"},
		},
		"loc": types.M{
			"$nearSphere": types.M{
				"__type":    "GeoPoint",
				"latitude":  40,
				"longitude": -30,
			},
		},
	}
	if reflect.DeepEqual(query, expectQuery) == false {
		t.Error("expect:", expectQuery, "result:", query)
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
	schemaCache = cache.NewSchemaCache(5, false)
	TomatoDBController = &DBController{}
}
