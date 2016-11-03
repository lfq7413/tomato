package rest

import (
	"reflect"
	"testing"

	"gopkg.in/mgo.v2"

	"github.com/lfq7413/tomato/cache"
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/storage"
	"github.com/lfq7413/tomato/storage/mongo"
	"github.com/lfq7413/tomato/types"
)

func Test_Execute(t *testing.T) {
	var schema types.M
	var object types.M
	var where types.M
	var options types.M
	var className string
	var q *Query
	var err error
	var result types.M
	var expect types.M
	/**********************************************************/
	initEnv()
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, object)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	className = "user"
	where = types.M{}
	options = types.M{}
	q, _ = NewQuery(Master(), className, where, options, nil)
	result, err = q.Execute()
	expect = types.M{
		"results": types.S{
			types.M{
				"objectId": "01",
				"key":      "hello",
			},
			types.M{
				"objectId": "02",
				"key":      "hello",
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "list"
	schema = types.M{
		"fields": types.M{
			"title": types.M{"type": "String"},
			"post": types.M{
				"type":        "Pointer",
				"targetClass": "post",
			},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"title":    "one",
		"post": types.M{
			"__type":    "Pointer",
			"className": "post",
			"objectId":  "2001",
		},
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "1002",
		"title":    "two",
		"post": types.M{
			"__type":    "Pointer",
			"className": "post",
			"objectId":  "2002",
		},
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "post"
	schema = types.M{
		"fields": types.M{
			"id": types.M{"type": "String"},
			"user": types.M{
				"type":        "Pointer",
				"targetClass": "user",
			},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"id":       "01",
		"user": types.M{
			"__type":    "Pointer",
			"className": "user",
			"objectId":  "3001",
		},
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2002",
		"id":       "02",
		"user": types.M{
			"__type":    "Pointer",
			"className": "user",
			"objectId":  "3002",
		},
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "user"
	schema = types.M{
		"fields": types.M{
			"name": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "3001",
		"name":     "joe",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "3002",
		"name":     "jack",
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "list"
	where = types.M{}
	options = types.M{"include": "post.user"}
	q, _ = NewQuery(Master(), className, where, options, nil)
	result, err = q.Execute()
	expect = types.M{
		"results": types.S{
			types.M{
				"objectId": "1001",
				"title":    "one",
				"post": types.M{
					"__type":    "Object",
					"className": "post",
					"objectId":  "2001",
					"id":        "01",
					"user": types.M{
						"__type":    "Object",
						"className": "user",
						"objectId":  "3001",
						"name":      "joe",
					},
				},
			},
			types.M{
				"objectId": "1002",
				"title":    "two",
				"post": types.M{
					"__type":    "Object",
					"className": "post",
					"objectId":  "2002",
					"id":        "02",
					"user": types.M{
						"__type":    "Object",
						"className": "user",
						"objectId":  "3002",
						"name":      "jack",
					},
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_BuildRestWhere(t *testing.T) {
	var className string
	var schema types.M
	var object types.M
	var q *Query
	var where types.M
	var err error
	var expect types.M
	/**********************************************************/
	initEnv()
	className = "Team"
	schema = types.M{
		"fields": types.M{
			"city":   types.M{"type": "String"},
			"winPct": types.M{"type": "Number"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"city":     "beijing",
		"winPct":   0.8,
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2002",
		"city":     "shanghai",
		"winPct":   0.7,
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2003",
		"city":     "guangzhou",
		"winPct":   0.4,
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "Post"
	schema = types.M{
		"fields": types.M{
			"title": types.M{"type": "String"},
			"image": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "3001",
		"title":    "one",
		"image":    "1.jpg",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "3002",
		"title":    "two",
		"image":    "2.jpg",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "3003",
		"title":    "three",
	}
	orm.Adapter.CreateObject(className, schema, object)
	where = types.M{
		"hometown": types.M{
			"$select": types.M{
				"query": types.M{
					"className": "Team",
					"where": types.M{
						"winPct": types.M{"$gt": 0.5},
					},
				},
				"key": "city",
			},
		},
		"post": types.M{
			"$inQuery": types.M{
				"where": types.M{
					"image": types.M{"$exists": true},
				},
				"className": "Post",
			},
		},
	}
	q, _ = NewQuery(Master(), "user", where, nil, nil)
	err = q.BuildRestWhere()
	expect = types.M{
		"hometown": types.M{
			"$in": types.S{"beijing", "shanghai"},
		},
		"post": types.M{
			"$in": types.S{
				types.M{
					"__type":    "Pointer",
					"className": "Post",
					"objectId":  "3001",
				},
				types.M{
					"__type":    "Pointer",
					"className": "Post",
					"objectId":  "3002",
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, q.Where) == false {
		t.Error("expect:", expect, "result:", q.Where, err)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_getUserAndRoleACL(t *testing.T) {
	var schema types.M
	var object types.M
	var className string
	var auth *Auth
	var q *Query
	var expect []string
	/**********************************************************/
	auth = Master()
	q, _ = NewQuery(auth, "user", nil, nil, nil)
	q.getUserAndRoleACL()
	if _, ok := q.findOptions["acl"]; ok {
		t.Error("findOptions[acl] exist")
	}
	/**********************************************************/
	auth = Nobody()
	q, _ = NewQuery(auth, "user", nil, nil, nil)
	q.getUserAndRoleACL()
	if q.findOptions["acl"] != nil {
		t.Error("findOptions[acl] is not nil")
	}
	/**********************************************************/
	cache.InitCache()
	initEnv()
	className = "_Role"
	schema = types.M{
		"fields": types.M{
			"name":  types.M{"type": "String"},
			"users": types.M{"type": "Relation", "targetClass": "_User"},
			"roles": types.M{"type": "Relation", "targetClass": "_Role"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"name":     "role1001",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "1002",
		"name":     "role1002",
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "_Join:roles:_Role"
	schema = types.M{
		"fields": types.M{
			"relatedId": types.M{"type": "String"},
			"owningId":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId":  "5001",
		"owningId":  "1002",
		"relatedId": "1001",
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "_Join:users:_Role"
	schema = types.M{
		"fields": types.M{
			"relatedId": types.M{"type": "String"},
			"owningId":  types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId":  "5002",
		"owningId":  "1001",
		"relatedId": "9001",
	}
	orm.Adapter.CreateObject(className, schema, object)
	auth = &Auth{
		IsMaster: false,
		User: types.M{
			"objectId": "9001",
		},
		FetchedRoles: false,
		RolePromise:  nil,
	}
	q, _ = NewQuery(auth, "user", nil, nil, nil)
	q.getUserAndRoleACL()
	expect = []string{"9001", "role:role1001", "role:role1002"}
	if reflect.DeepEqual(expect, q.findOptions["acl"]) == false {
		t.Error("expect:", expect, "result:", q.findOptions["acl"])
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_redirectClassNameForKey(t *testing.T) {
	var options types.M
	var q *Query
	var object types.M
	/**********************************************************/
	options = types.M{}
	q, _ = NewQuery(nil, "user", nil, options, nil)
	q.redirectClassNameForKey()
	if q.redirectKey != "" || q.redirectClassName != "" || q.className != "user" {
		t.Error("expect: empty result:", q.redirectKey, q.redirectClassName, q.className)
	}
	/**********************************************************/
	initEnv()
	options = types.M{"redirectClassNameForKey": "post"}
	q, _ = NewQuery(nil, "user", nil, options, nil)
	q.redirectClassNameForKey()
	if q.redirectKey != "post" || q.redirectClassName != "user" || q.className != "user" {
		t.Error("expect: empty result:", q.redirectKey, q.redirectClassName, q.className)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	object = types.M{
		"fields": types.M{
			"post": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("user", object)
	options = types.M{"redirectClassNameForKey": "post"}
	q, _ = NewQuery(nil, "user", nil, options, nil)
	q.redirectClassNameForKey()
	if q.redirectKey != "post" || q.redirectClassName != "user" || q.className != "user" {
		t.Error("expect: empty result:", q.redirectKey, q.redirectClassName, q.className)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	object = types.M{
		"fields": types.M{
			"post": types.M{
				"type":        "Relation",
				"targetClass": "postT",
			},
		},
	}
	orm.Adapter.CreateClass("user", object)
	options = types.M{"redirectClassNameForKey": "post"}
	q, _ = NewQuery(nil, "user", nil, options, nil)
	q.redirectClassNameForKey()
	if q.redirectKey != "post" || q.redirectClassName != "postT" || q.className != "postT" {
		t.Error("expect: empty result:", q.redirectKey, q.redirectClassName, q.className)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_validateClientClassCreation(t *testing.T) {
	var className string
	var q *Query
	var result error
	var expect error
	/**********************************************************/
	config.TConfig.AllowClientClassCreation = true
	className = "user"
	q, _ = NewQuery(nil, className, nil, nil, nil)
	result = q.validateClientClassCreation()
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	config.TConfig.AllowClientClassCreation = false
	className = "user"
	q, _ = NewQuery(Master(), className, nil, nil, nil)
	result = q.validateClientClassCreation()
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	config.TConfig.AllowClientClassCreation = false
	className = "_User"
	q, _ = NewQuery(nil, className, nil, nil, nil)
	result = q.validateClientClassCreation()
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	initEnv()
	object := types.M{
		"fields": types.M{
			"post": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("user", object)
	config.TConfig.AllowClientClassCreation = false
	className = "user"
	q, _ = NewQuery(nil, className, nil, nil, nil)
	result = q.validateClientClassCreation()
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	config.TConfig.AllowClientClassCreation = false
	className = "user"
	q, _ = NewQuery(nil, className, nil, nil, nil)
	result = q.validateClientClassCreation()
	expect = errs.E(errs.OperationForbidden, "This user is not allowed to access non-existent class: user")
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_replaceSelect(t *testing.T) {
	var className string
	var schema types.M
	var object types.M
	var q *Query
	var where types.M
	var err error
	var expectErr error
	var expect types.M
	/**********************************************************/
	where = types.M{"key": "hello"}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceSelect()
	expect = types.M{"key": "hello"}
	if err != nil || reflect.DeepEqual(expect, q.Where) == false {
		t.Error("expect:", expect, "result:", q.Where)
	}
	/**********************************************************/
	where = types.M{
		"hometown": types.M{
			"$select": "hello",
		},
	}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceSelect()
	expectErr = errs.E(errs.InvalidQuery, "improper usage of $select")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/**********************************************************/
	where = types.M{
		"hometown": types.M{
			"$select": types.M{},
		},
	}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceSelect()
	expectErr = errs.E(errs.InvalidQuery, "improper usage of $select")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/**********************************************************/
	where = types.M{
		"hometown": types.M{
			"$select": types.M{
				"query": "hello",
			},
		},
	}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceSelect()
	expectErr = errs.E(errs.InvalidQuery, "improper usage of $select")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/**********************************************************/
	where = types.M{
		"hometown": types.M{
			"$select": types.M{
				"query": "hello",
				"key":   "hello",
				"other": "world",
			},
		},
	}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceSelect()
	expectErr = errs.E(errs.InvalidQuery, "improper usage of $select")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/**********************************************************/
	where = types.M{
		"hometown": types.M{
			"$select": types.M{
				"query": "hello",
				"key":   "hello",
			},
		},
	}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceSelect()
	expectErr = errs.E(errs.InvalidQuery, "improper usage of $select")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/**********************************************************/
	where = types.M{
		"hometown": types.M{
			"$select": types.M{
				"query": types.M{},
				"key":   "hello",
			},
		},
	}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceSelect()
	expectErr = errs.E(errs.InvalidQuery, "improper usage of $select")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/**********************************************************/
	initEnv()
	where = types.M{
		"hometown": types.M{
			"$select": types.M{
				"query": types.M{
					"className": "Team",
				},
				"key": "city",
			},
		},
	}
	q, _ = NewQuery(Master(), "user", where, nil, nil)
	err = q.replaceSelect()
	expect = types.M{
		"hometown": types.M{
			"$in": types.S{},
		},
	}
	if err != nil || reflect.DeepEqual(expect, q.Where) == false {
		t.Error("expect:", expect, "result:", q.Where, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "Team"
	schema = types.M{
		"fields": types.M{
			"city": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"city":     "beijing",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2002",
		"city":     "shanghai",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2003",
		"city":     "guangzhou",
	}
	orm.Adapter.CreateObject(className, schema, object)
	where = types.M{
		"hometown": types.M{
			"$select": types.M{
				"query": types.M{
					"className": "Team",
				},
				"key": "city",
			},
		},
	}
	q, _ = NewQuery(Master(), "user", where, nil, nil)
	err = q.replaceSelect()
	expect = types.M{
		"hometown": types.M{
			"$in": types.S{"beijing", "shanghai", "guangzhou"},
		},
	}
	if err != nil || reflect.DeepEqual(expect, q.Where) == false {
		t.Error("expect:", expect, "result:", q.Where, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "Team"
	schema = types.M{
		"fields": types.M{
			"city":   types.M{"type": "String"},
			"winPct": types.M{"type": "Number"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"city":     "beijing",
		"winPct":   0.8,
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2002",
		"city":     "shanghai",
		"winPct":   0.7,
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2003",
		"city":     "guangzhou",
		"winPct":   0.4,
	}
	orm.Adapter.CreateObject(className, schema, object)
	where = types.M{
		"hometown": types.M{
			"$select": types.M{
				"query": types.M{
					"className": "Team",
					"where": types.M{
						"winPct": types.M{"$gt": 0.5},
					},
				},
				"key": "city",
			},
		},
	}
	q, _ = NewQuery(Master(), "user", where, nil, nil)
	err = q.replaceSelect()
	expect = types.M{
		"hometown": types.M{
			"$in": types.S{"beijing", "shanghai"},
		},
	}
	if err != nil || reflect.DeepEqual(expect, q.Where) == false {
		t.Error("expect:", expect, "result:", q.Where, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "Team"
	schema = types.M{
		"fields": types.M{
			"city":   types.M{"type": "String"},
			"winPct": types.M{"type": "Number"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"city":     "beijing",
		"winPct":   0.8,
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2002",
		"city":     "shanghai",
		"winPct":   0.7,
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2003",
		"city":     "guangzhou",
		"winPct":   0.4,
	}
	orm.Adapter.CreateObject(className, schema, object)
	where = types.M{
		"hometown": types.M{
			"$select": types.M{
				"query": types.M{
					"className": "Team",
					"where": types.M{
						"winPct": types.M{"$gt": 0.5},
					},
				},
				"key": "city",
			},
		},
		"hometown2": types.M{
			"$select": types.M{
				"query": types.M{
					"className": "Team",
					"where": types.M{
						"winPct": types.M{"$gt": 0.7},
					},
				},
				"key": "city",
			},
		},
	}
	q, _ = NewQuery(Master(), "user", where, nil, nil)
	err = q.replaceSelect()
	expect = types.M{
		"hometown": types.M{
			"$in": types.S{"beijing", "shanghai"},
		},
		"hometown2": types.M{
			"$in": types.S{"beijing"},
		},
	}
	if err != nil || reflect.DeepEqual(expect, q.Where) == false {
		t.Error("expect:", expect, "result:", q.Where, err)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_replaceDontSelect(t *testing.T) {
	var className string
	var schema types.M
	var object types.M
	var q *Query
	var where types.M
	var err error
	var expectErr error
	var expect types.M
	/**********************************************************/
	where = types.M{"key": "hello"}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceDontSelect()
	expect = types.M{"key": "hello"}
	if err != nil || reflect.DeepEqual(expect, q.Where) == false {
		t.Error("expect:", expect, "result:", q.Where)
	}
	/**********************************************************/
	where = types.M{
		"hometown": types.M{
			"$dontSelect": "hello",
		},
	}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceDontSelect()
	expectErr = errs.E(errs.InvalidQuery, "improper usage of $dontSelect")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/**********************************************************/
	where = types.M{
		"hometown": types.M{
			"$dontSelect": types.M{},
		},
	}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceDontSelect()
	expectErr = errs.E(errs.InvalidQuery, "improper usage of $dontSelect")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/**********************************************************/
	where = types.M{
		"hometown": types.M{
			"$dontSelect": types.M{
				"query": "hello",
			},
		},
	}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceDontSelect()
	expectErr = errs.E(errs.InvalidQuery, "improper usage of $dontSelect")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/**********************************************************/
	where = types.M{
		"hometown": types.M{
			"$dontSelect": types.M{
				"query": "hello",
				"key":   "hello",
				"other": "world",
			},
		},
	}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceDontSelect()
	expectErr = errs.E(errs.InvalidQuery, "improper usage of $dontSelect")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/**********************************************************/
	where = types.M{
		"hometown": types.M{
			"$dontSelect": types.M{
				"query": "hello",
				"key":   "hello",
			},
		},
	}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceDontSelect()
	expectErr = errs.E(errs.InvalidQuery, "improper usage of $dontSelect")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/**********************************************************/
	where = types.M{
		"hometown": types.M{
			"$dontSelect": types.M{
				"query": types.M{},
				"key":   "hello",
			},
		},
	}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceDontSelect()
	expectErr = errs.E(errs.InvalidQuery, "improper usage of $dontSelect")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/**********************************************************/
	initEnv()
	where = types.M{
		"hometown": types.M{
			"$dontSelect": types.M{
				"query": types.M{
					"className": "Team",
				},
				"key": "city",
			},
		},
	}
	q, _ = NewQuery(Master(), "user", where, nil, nil)
	err = q.replaceDontSelect()
	expect = types.M{
		"hometown": types.M{
			"$nin": types.S{},
		},
	}
	if err != nil || reflect.DeepEqual(expect, q.Where) == false {
		t.Error("expect:", expect, "result:", q.Where, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "Team"
	schema = types.M{
		"fields": types.M{
			"city": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"city":     "beijing",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2002",
		"city":     "shanghai",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2003",
		"city":     "guangzhou",
	}
	orm.Adapter.CreateObject(className, schema, object)
	where = types.M{
		"hometown": types.M{
			"$dontSelect": types.M{
				"query": types.M{
					"className": "Team",
				},
				"key": "city",
			},
		},
	}
	q, _ = NewQuery(Master(), "user", where, nil, nil)
	err = q.replaceDontSelect()
	expect = types.M{
		"hometown": types.M{
			"$nin": types.S{"beijing", "shanghai", "guangzhou"},
		},
	}
	if err != nil || reflect.DeepEqual(expect, q.Where) == false {
		t.Error("expect:", expect, "result:", q.Where, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "Team"
	schema = types.M{
		"fields": types.M{
			"city":   types.M{"type": "String"},
			"winPct": types.M{"type": "Number"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"city":     "beijing",
		"winPct":   0.8,
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2002",
		"city":     "shanghai",
		"winPct":   0.7,
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2003",
		"city":     "guangzhou",
		"winPct":   0.4,
	}
	orm.Adapter.CreateObject(className, schema, object)
	where = types.M{
		"hometown": types.M{
			"$dontSelect": types.M{
				"query": types.M{
					"className": "Team",
					"where": types.M{
						"winPct": types.M{"$gt": 0.5},
					},
				},
				"key": "city",
			},
		},
	}
	q, _ = NewQuery(Master(), "user", where, nil, nil)
	err = q.replaceDontSelect()
	expect = types.M{
		"hometown": types.M{
			"$nin": types.S{"beijing", "shanghai"},
		},
	}
	if err != nil || reflect.DeepEqual(expect, q.Where) == false {
		t.Error("expect:", expect, "result:", q.Where, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "Team"
	schema = types.M{
		"fields": types.M{
			"city":   types.M{"type": "String"},
			"winPct": types.M{"type": "Number"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"city":     "beijing",
		"winPct":   0.8,
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2002",
		"city":     "shanghai",
		"winPct":   0.7,
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2003",
		"city":     "guangzhou",
		"winPct":   0.4,
	}
	orm.Adapter.CreateObject(className, schema, object)
	where = types.M{
		"hometown": types.M{
			"$dontSelect": types.M{
				"query": types.M{
					"className": "Team",
					"where": types.M{
						"winPct": types.M{"$gt": 0.5},
					},
				},
				"key": "city",
			},
		},
		"hometown2": types.M{
			"$dontSelect": types.M{
				"query": types.M{
					"className": "Team",
					"where": types.M{
						"winPct": types.M{"$gt": 0.7},
					},
				},
				"key": "city",
			},
		},
	}
	q, _ = NewQuery(Master(), "user", where, nil, nil)
	err = q.replaceDontSelect()
	expect = types.M{
		"hometown": types.M{
			"$nin": types.S{"beijing", "shanghai"},
		},
		"hometown2": types.M{
			"$nin": types.S{"beijing"},
		},
	}
	if err != nil || reflect.DeepEqual(expect, q.Where) == false {
		t.Error("expect:", expect, "result:", q.Where, err)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_replaceInQuery(t *testing.T) {
	var className string
	var schema types.M
	var object types.M
	var q *Query
	var where types.M
	var err error
	var expectErr error
	var expect types.M
	/**********************************************************/
	where = types.M{"key": "hello"}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceInQuery()
	expect = types.M{"key": "hello"}
	if err != nil || reflect.DeepEqual(expect, q.Where) == false {
		t.Error("expect:", expect, "result:", q.Where)
	}
	/**********************************************************/
	where = types.M{
		"post": types.M{
			"$inQuery": "hello",
		},
	}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceInQuery()
	expectErr = errs.E(errs.InvalidQuery, "improper usage of $inQuery")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/**********************************************************/
	where = types.M{
		"post": types.M{
			"$inQuery": types.M{},
		},
	}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceInQuery()
	expectErr = errs.E(errs.InvalidQuery, "improper usage of $inQuery")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/**********************************************************/
	where = types.M{
		"post": types.M{
			"$inQuery": types.M{
				"where": "hello",
			},
		},
	}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceInQuery()
	expectErr = errs.E(errs.InvalidQuery, "improper usage of $inQuery")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/**********************************************************/
	where = types.M{
		"post": types.M{
			"$inQuery": types.M{
				"where":     types.M{},
				"className": 1024,
			},
		},
	}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceInQuery()
	expectErr = errs.E(errs.InvalidQuery, "improper usage of $inQuery")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/**********************************************************/
	initEnv()
	where = types.M{
		"post": types.M{
			"$inQuery": types.M{
				"where":     types.M{},
				"className": "Post",
			},
		},
	}
	q, _ = NewQuery(Master(), "user", where, nil, nil)
	err = q.replaceInQuery()
	expect = types.M{
		"post": types.M{
			"$in": types.S{},
		},
	}
	if err != nil || reflect.DeepEqual(expect, expect) == false {
		t.Error("expect:", expect, "result:", expect, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "Post"
	schema = types.M{
		"fields": types.M{
			"title": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"title":    "one",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2002",
		"title":    "two",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2003",
		"title":    "three",
	}
	orm.Adapter.CreateObject(className, schema, object)
	where = types.M{
		"post": types.M{
			"$inQuery": types.M{
				"where":     types.M{},
				"className": "Post",
			},
		},
	}
	q, _ = NewQuery(Master(), "user", where, nil, nil)
	err = q.replaceInQuery()
	expect = types.M{
		"post": types.M{
			"$in": types.S{
				types.M{
					"__type":    "Pointer",
					"className": "Post",
					"objectId":  "2001",
				},
				types.M{
					"__type":    "Pointer",
					"className": "Post",
					"objectId":  "2002",
				},
				types.M{
					"__type":    "Pointer",
					"className": "Post",
					"objectId":  "2003",
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, expect) == false {
		t.Error("expect:", expect, "result:", expect, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "Post"
	schema = types.M{
		"fields": types.M{
			"title": types.M{"type": "String"},
			"image": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"title":    "one",
		"image":    "1.jpg",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2002",
		"title":    "two",
		"image":    "2.jpg",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2003",
		"title":    "three",
	}
	orm.Adapter.CreateObject(className, schema, object)
	where = types.M{
		"post": types.M{
			"$inQuery": types.M{
				"where": types.M{
					"image": types.M{"$exists": true},
				},
				"className": "Post",
			},
		},
	}
	q, _ = NewQuery(Master(), "user", where, nil, nil)
	err = q.replaceInQuery()
	expect = types.M{
		"post": types.M{
			"$in": types.S{
				types.M{
					"__type":    "Pointer",
					"className": "Post",
					"objectId":  "2001",
				},
				types.M{
					"__type":    "Pointer",
					"className": "Post",
					"objectId":  "2002",
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, expect) == false {
		t.Error("expect:", expect, "result:", expect, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "Post"
	schema = types.M{
		"fields": types.M{
			"title":  types.M{"type": "String"},
			"image":  types.M{"type": "String"},
			"author": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"title":    "one",
		"image":    "1.jpg",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2002",
		"title":    "two",
		"image":    "2.jpg",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2003",
		"title":    "three",
		"author":   "joe",
	}
	orm.Adapter.CreateObject(className, schema, object)
	where = types.M{
		"post": types.M{
			"$inQuery": types.M{
				"where": types.M{
					"image": types.M{"$exists": true},
				},
				"className": "Post",
			},
		},
		"post2": types.M{
			"$inQuery": types.M{
				"where": types.M{
					"author": types.M{"$exists": true},
				},
				"className": "Post",
			},
		},
	}
	q, _ = NewQuery(Master(), "user", where, nil, nil)
	err = q.replaceInQuery()
	expect = types.M{
		"post": types.M{
			"$in": types.S{
				types.M{
					"__type":    "Pointer",
					"className": "Post",
					"objectId":  "2001",
				},
				types.M{
					"__type":    "Pointer",
					"className": "Post",
					"objectId":  "2002",
				},
			},
		},
		"post2": types.M{
			"$in": types.S{
				types.M{
					"__type":    "Pointer",
					"className": "Post",
					"objectId":  "2003",
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, expect) == false {
		t.Error("expect:", expect, "result:", expect, err)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_replaceNotInQuery(t *testing.T) {
	var className string
	var schema types.M
	var object types.M
	var q *Query
	var where types.M
	var err error
	var expectErr error
	var expect types.M
	/**********************************************************/
	where = types.M{"key": "hello"}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceNotInQuery()
	expect = types.M{"key": "hello"}
	if err != nil || reflect.DeepEqual(expect, q.Where) == false {
		t.Error("expect:", expect, "result:", q.Where)
	}
	/**********************************************************/
	where = types.M{
		"post": types.M{
			"$notInQuery": "hello",
		},
	}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceNotInQuery()
	expectErr = errs.E(errs.InvalidQuery, "improper usage of $notInQuery")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/**********************************************************/
	where = types.M{
		"post": types.M{
			"$notInQuery": types.M{},
		},
	}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceNotInQuery()
	expectErr = errs.E(errs.InvalidQuery, "improper usage of $notInQuery")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/**********************************************************/
	where = types.M{
		"post": types.M{
			"$notInQuery": types.M{
				"where": "hello",
			},
		},
	}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceNotInQuery()
	expectErr = errs.E(errs.InvalidQuery, "improper usage of $notInQuery")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/**********************************************************/
	where = types.M{
		"post": types.M{
			"$notInQuery": types.M{
				"where":     types.M{},
				"className": 1024,
			},
		},
	}
	q, _ = NewQuery(nil, "user", where, nil, nil)
	err = q.replaceNotInQuery()
	expectErr = errs.E(errs.InvalidQuery, "improper usage of $notInQuery")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	/**********************************************************/
	initEnv()
	where = types.M{
		"post": types.M{
			"$notInQuery": types.M{
				"where":     types.M{},
				"className": "Post",
			},
		},
	}
	q, _ = NewQuery(Master(), "user", where, nil, nil)
	err = q.replaceNotInQuery()
	expect = types.M{
		"post": types.M{
			"$nin": types.S{},
		},
	}
	if err != nil || reflect.DeepEqual(expect, expect) == false {
		t.Error("expect:", expect, "result:", expect, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "Post"
	schema = types.M{
		"fields": types.M{
			"title": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"title":    "one",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2002",
		"title":    "two",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2003",
		"title":    "three",
	}
	orm.Adapter.CreateObject(className, schema, object)
	where = types.M{
		"post": types.M{
			"$notInQuery": types.M{
				"where":     types.M{},
				"className": "Post",
			},
		},
	}
	q, _ = NewQuery(Master(), "user", where, nil, nil)
	err = q.replaceNotInQuery()
	expect = types.M{
		"post": types.M{
			"$nin": types.S{
				types.M{
					"__type":    "Pointer",
					"className": "Post",
					"objectId":  "2001",
				},
				types.M{
					"__type":    "Pointer",
					"className": "Post",
					"objectId":  "2002",
				},
				types.M{
					"__type":    "Pointer",
					"className": "Post",
					"objectId":  "2003",
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, expect) == false {
		t.Error("expect:", expect, "result:", expect, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "Post"
	schema = types.M{
		"fields": types.M{
			"title": types.M{"type": "String"},
			"image": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"title":    "one",
		"image":    "1.jpg",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2002",
		"title":    "two",
		"image":    "2.jpg",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2003",
		"title":    "three",
	}
	orm.Adapter.CreateObject(className, schema, object)
	where = types.M{
		"post": types.M{
			"$notInQuery": types.M{
				"where": types.M{
					"image": types.M{"$exists": true},
				},
				"className": "Post",
			},
		},
	}
	q, _ = NewQuery(Master(), "user", where, nil, nil)
	err = q.replaceNotInQuery()
	expect = types.M{
		"post": types.M{
			"$nin": types.S{
				types.M{
					"__type":    "Pointer",
					"className": "Post",
					"objectId":  "2001",
				},
				types.M{
					"__type":    "Pointer",
					"className": "Post",
					"objectId":  "2002",
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, expect) == false {
		t.Error("expect:", expect, "result:", expect, err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "Post"
	schema = types.M{
		"fields": types.M{
			"title":  types.M{"type": "String"},
			"image":  types.M{"type": "String"},
			"author": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"title":    "one",
		"image":    "1.jpg",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2002",
		"title":    "two",
		"image":    "2.jpg",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2003",
		"title":    "three",
		"author":   "joe",
	}
	orm.Adapter.CreateObject(className, schema, object)
	where = types.M{
		"post": types.M{
			"$notInQuery": types.M{
				"where": types.M{
					"image": types.M{"$exists": true},
				},
				"className": "Post",
			},
		},
		"post2": types.M{
			"$notInQuery": types.M{
				"where": types.M{
					"author": types.M{"$exists": true},
				},
				"className": "Post",
			},
		},
	}
	q, _ = NewQuery(Master(), "user", where, nil, nil)
	err = q.replaceNotInQuery()
	expect = types.M{
		"post": types.M{
			"$nin": types.S{
				types.M{
					"__type":    "Pointer",
					"className": "Post",
					"objectId":  "2001",
				},
				types.M{
					"__type":    "Pointer",
					"className": "Post",
					"objectId":  "2002",
				},
			},
		},
		"post2": types.M{
			"$nin": types.S{
				types.M{
					"__type":    "Pointer",
					"className": "Post",
					"objectId":  "2003",
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, expect) == false {
		t.Error("expect:", expect, "result:", expect, err)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_runFind(t *testing.T) {
	var object types.M
	var where types.M
	var options types.M
	var className string
	var q *Query
	var err error
	var expect types.S
	/**********************************************************/
	initEnv()
	where = types.M{}
	options = types.M{}
	className = "user"
	q, _ = NewQuery(Master(), className, where, options, nil)
	err = q.runFind()
	expect = types.S{}
	if err != nil || reflect.DeepEqual(expect, q.response["results"]) == false {
		t.Error("expect:", expect, "result:", q.response["results"], err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, object)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	initEnv()
	where = types.M{}
	options = types.M{}
	className = "user"
	q, _ = NewQuery(Master(), className, where, options, nil)
	err = q.runFind()
	expect = types.S{
		types.M{
			"objectId": "01",
			"key":      "hello",
		},
	}
	if err != nil || reflect.DeepEqual(expect, q.response["results"]) == false {
		t.Error("expect:", expect, "result:", q.response["results"], err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, object)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	initEnv()
	where = types.M{}
	options = types.M{"limit": 0}
	className = "user"
	q, _ = NewQuery(Master(), className, where, options, nil)
	err = q.runFind()
	expect = types.S{}
	if err != nil || reflect.DeepEqual(expect, q.response["results"]) == false {
		t.Error("expect:", expect, "result:", q.response["results"], err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, object)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	initEnv()
	where = types.M{}
	options = types.M{"limit": 1}
	className = "user"
	q, _ = NewQuery(Master(), className, where, options, nil)
	err = q.runFind()
	expect = types.S{
		types.M{
			"objectId": "01",
			"key":      "hello",
		},
	}
	if err != nil || reflect.DeepEqual(expect, q.response["results"]) == false {
		t.Error("expect:", expect, "result:", q.response["results"], err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, object)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	initEnv()
	where = types.M{}
	options = types.M{"skip": 1}
	className = "user"
	q, _ = NewQuery(Master(), className, where, options, nil)
	err = q.runFind()
	expect = types.S{
		types.M{
			"objectId": "02",
			"key":      "hello",
		},
	}
	if err != nil || reflect.DeepEqual(expect, q.response["results"]) == false {
		t.Error("expect:", expect, "result:", q.response["results"], err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
			"age": types.M{"type": "Number"},
		},
	}
	orm.Adapter.CreateClass(className, object)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
		"age":      2,
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
		"age":      3,
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "03",
		"key":      "hello",
		"age":      1,
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	initEnv()
	where = types.M{}
	options = types.M{"order": "age"}
	className = "user"
	q, _ = NewQuery(Master(), className, where, options, nil)
	err = q.runFind()
	expect = types.S{
		types.M{
			"objectId": "03",
			"key":      "hello",
			"age":      1,
		},
		types.M{
			"objectId": "01",
			"key":      "hello",
			"age":      2,
		},
		types.M{
			"objectId": "02",
			"key":      "hello",
			"age":      3,
		},
	}
	if err != nil || reflect.DeepEqual(expect, q.response["results"]) == false {
		t.Error("expect:", expect, "result:", q.response["results"], err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
			"age": types.M{"type": "Number"},
		},
	}
	orm.Adapter.CreateClass(className, object)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
		"age":      2,
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
		"age":      3,
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "03",
		"key":      "hello",
		"age":      1,
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	initEnv()
	where = types.M{}
	options = types.M{"keys": "age"}
	className = "user"
	q, _ = NewQuery(Master(), className, where, options, nil)
	err = q.runFind()
	expect = types.S{
		types.M{
			"objectId": "01",
			"age":      2,
		},
		types.M{
			"objectId": "02",
			"age":      3,
		},
		types.M{
			"objectId": "03",
			"age":      1,
		},
	}
	if err != nil || reflect.DeepEqual(expect, q.response["results"]) == false {
		t.Error("expect:", expect, "result:", q.response["results"], err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
			"age": types.M{"type": "Number"},
		},
	}
	orm.Adapter.CreateClass(className, object)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
		"age":      2,
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
		"age":      3,
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "03",
		"key":      "hello",
		"age":      1,
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	initEnv()
	where = types.M{}
	options = types.M{"keys": "age.id"}
	className = "user"
	q, _ = NewQuery(Master(), className, where, options, nil)
	err = q.runFind()
	expect = types.S{
		types.M{
			"objectId": "01",
			"age":      2,
		},
		types.M{
			"objectId": "02",
			"age":      3,
		},
		types.M{
			"objectId": "03",
			"age":      1,
		},
	}
	if err != nil || reflect.DeepEqual(expect, q.response["results"]) == false {
		t.Error("expect:", expect, "result:", q.response["results"], err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
			"age": types.M{"type": "Number"},
		},
	}
	orm.Adapter.CreateClass(className, object)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
		"age":      2,
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
		"age":      3,
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	initEnv()
	where = types.M{}
	options = types.M{}
	className = "user"
	q, _ = NewQuery(Master(), className, where, options, nil)
	q.redirectClassName = "post"
	err = q.runFind()
	expect = types.S{
		types.M{
			"className": "post",
			"objectId":  "01",
			"key":       "hello",
			"age":       2,
		},
		types.M{
			"className": "post",
			"objectId":  "02",
			"key":       "hello",
			"age":       3,
		},
	}
	if err != nil || reflect.DeepEqual(expect, q.response["results"]) == false {
		t.Error("expect:", expect, "result:", q.response["results"], err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "_User"
	object = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, object)
	object = types.M{
		"objectId": "01",
		"username": "joe",
		"password": "123456",
		"authData": types.M{
			"facebook": types.M{
				"id": "1001",
			},
		},
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"username": "jack",
		"password": "123456",
		"authData": types.M{
			"facebook": types.M{
				"id": "1002",
			},
		},
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	initEnv()
	where = types.M{}
	options = types.M{}
	className = "_User"
	q, _ = NewQuery(Master(), className, where, options, nil)
	err = q.runFind()
	expect = types.S{
		types.M{
			"objectId": "01",
			"username": "joe",
			"authData": types.M{
				"facebook": types.M{
					"id": "1001",
				},
			},
		},
		types.M{
			"objectId": "02",
			"username": "jack",
			"authData": types.M{
				"facebook": types.M{
					"id": "1002",
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, q.response["results"]) == false {
		t.Error("expect:", expect, "result:", q.response["results"], err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	config.TConfig.ServerURL = "http://127.0.0.1"
	config.TConfig.AppID = "1001"
	className = "_User"
	object = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
			"password": types.M{"type": "String"},
			"icon":     types.M{"type": "File"},
		},
	}
	orm.Adapter.CreateClass(className, object)
	object = types.M{
		"objectId": "01",
		"username": "joe",
		"password": "123456",
		"icon": types.M{
			"__type": "File",
			"name":   "icon1.jpg",
		},
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"username": "jack",
		"password": "123456",
		"icon": types.M{
			"__type": "File",
			"name":   "icon2.jpg",
		},
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	initEnv()
	where = types.M{}
	options = types.M{}
	className = "_User"
	q, _ = NewQuery(Master(), className, where, options, nil)
	err = q.runFind()
	expect = types.S{
		types.M{
			"objectId": "01",
			"username": "joe",
			"icon": types.M{
				"__type": "File",
				"name":   "icon1.jpg",
				"url":    "http://127.0.0.1/files/1001/icon1.jpg",
			},
		},
		types.M{
			"objectId": "02",
			"username": "jack",
			"icon": types.M{
				"__type": "File",
				"name":   "icon2.jpg",
				"url":    "http://127.0.0.1/files/1001/icon2.jpg",
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, q.response["results"]) == false {
		t.Error("expect:", expect, "result:", q.response["results"], err)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_runCount(t *testing.T) {
	var object types.M
	var options types.M
	var className string
	var q *Query
	var err error
	var expect int
	/**********************************************************/
	options = types.M{}
	className = "user"
	q, _ = NewQuery(Master(), className, types.M{}, options, nil)
	err = q.runCount()
	if err != nil || q.response["count"] != nil {
		t.Error("expect:", nil, "result:", q.response["count"], err)
	}
	/**********************************************************/
	initEnv()
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, object)
	options = types.M{"count": true}
	className = "user"
	q, _ = NewQuery(Master(), className, types.M{}, options, nil)
	err = q.runCount()
	expect = 0
	if err != nil || reflect.DeepEqual(q.response["count"], expect) == false {
		t.Error("expect:", nil, "result:", q.response["count"], err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, object)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	options = types.M{"count": true}
	className = "user"
	q, _ = NewQuery(Master(), className, types.M{}, options, nil)
	err = q.runCount()
	expect = 2
	if err != nil || reflect.DeepEqual(q.response["count"], expect) == false {
		t.Error("expect:", nil, "result:", q.response["count"], err)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, object)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	options = types.M{
		"count": true,
		"skip":  1,
		"limit": 1,
	}
	className = "user"
	q, _ = NewQuery(Master(), className, types.M{}, options, nil)
	err = q.runCount()
	expect = 2
	if err != nil || reflect.DeepEqual(q.response["count"], expect) == false {
		t.Error("expect:", nil, "result:", q.response["count"], err)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_handleInclude(t *testing.T) {
	var schema types.M
	var object types.M
	var className string
	var options types.M
	var q *Query
	var err error
	var expect types.M
	/**********************************************************/
	initEnv()
	className = "post"
	schema = types.M{
		"fields": types.M{
			"id": types.M{"type": "String"},
			"user": types.M{
				"type":        "Pointer",
				"targetClass": "user",
			},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "2001",
		"id":       "01",
		"user": types.M{
			"__type":    "Pointer",
			"className": "user",
			"objectId":  "3001",
		},
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "2002",
		"id":       "02",
		"user": types.M{
			"__type":    "Pointer",
			"className": "user",
			"objectId":  "3002",
		},
	}
	orm.Adapter.CreateObject(className, schema, object)
	className = "user"
	schema = types.M{
		"fields": types.M{
			"name": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "3001",
		"name":     "joe",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "3002",
		"name":     "jack",
	}
	orm.Adapter.CreateObject(className, schema, object)
	options = types.M{"include": "post.user"}
	q, _ = NewQuery(Master(), "list", types.M{}, options, nil)
	q.response = types.M{
		"results": types.S{
			types.M{
				"objectId": "1001",
				"post": types.M{
					"__type":    "Pointer",
					"className": "post",
					"objectId":  "2001",
				},
			},
			types.M{
				"objectId": "1002",
				"post": types.M{
					"__type":    "Pointer",
					"className": "post",
					"objectId":  "2002",
				},
			},
		},
	}
	err = q.handleInclude()
	expect = types.M{
		"results": types.S{
			types.M{
				"objectId": "1001",
				"post": types.M{
					"__type":    "Object",
					"className": "post",
					"objectId":  "2001",
					"id":        "01",
					"user": types.M{
						"__type":    "Object",
						"className": "user",
						"objectId":  "3001",
						"name":      "joe",
					},
				},
			},
			types.M{
				"objectId": "1002",
				"post": types.M{
					"__type":    "Object",
					"className": "post",
					"objectId":  "2002",
					"id":        "02",
					"user": types.M{
						"__type":    "Object",
						"className": "user",
						"objectId":  "3002",
						"name":      "jack",
					},
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, q.response) == false {
		t.Error("expect:", expect, "result:", q.response)
	}
	orm.TomatoDBController.DeleteEverything()
}

/////////////////////////////////////////////////////////////////

func Test_NewQuery(t *testing.T) {
	var auth *Auth
	var className string
	var where types.M
	var options types.M
	var clientSDK map[string]string
	var result *Query
	var err error
	var expectErr error
	var expect *Query
	/**********************************************************/
	auth = Nobody()
	className = "user"
	where = nil
	options = nil
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:              auth,
		className:         "user",
		Where:             types.M{},
		findOptions:       types.M{"acl": nil},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = nil
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:              auth,
		className:         "user",
		Where:             types.M{},
		findOptions:       types.M{},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = &Auth{
		IsMaster: false,
		User: types.M{
			"objectId": "1001",
		},
	}
	className = "user"
	where = nil
	options = nil
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:              auth,
		className:         "user",
		Where:             types.M{},
		findOptions:       types.M{"acl": []string{"1001"}},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Nobody()
	className = "_Session"
	where = nil
	options = nil
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expectErr = errs.E(errs.InvalidSessionToken, "This session token is invalid.")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", result, err)
	}
	/**********************************************************/
	auth = &Auth{
		IsMaster: false,
		User: types.M{
			"objectId": "1001",
		},
	}
	className = "_Session"
	where = nil
	options = nil
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:      auth,
		className: "_Session",
		Where: types.M{
			"$and": types.S{
				types.M{},
				types.M{
					"user": types.M{
						"__type":    "Pointer",
						"className": "_User",
						"objectId":  "1001",
					},
				},
			},
		},
		findOptions:       types.M{"acl": []string{"1001"}},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = types.M{"keys": 1024}
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:              auth,
		className:         "user",
		Where:             types.M{},
		restOptions:       types.M{"keys": 1024},
		findOptions:       types.M{},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = types.M{"keys": "post"}
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:              auth,
		className:         "user",
		Where:             types.M{},
		restOptions:       types.M{"keys": "post"},
		findOptions:       types.M{},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{"post", "objectId", "createdAt", "updatedAt"},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = types.M{"keys": "post,user"}
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:              auth,
		className:         "user",
		Where:             types.M{},
		restOptions:       types.M{"keys": "post,user"},
		findOptions:       types.M{},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{"post", "user", "objectId", "createdAt", "updatedAt"},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = types.M{"count": true}
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:              auth,
		className:         "user",
		Where:             types.M{},
		restOptions:       types.M{"count": true},
		findOptions:       types.M{},
		response:          types.M{},
		doCount:           true,
		include:           [][]string{},
		keys:              []string{},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = types.M{"skip": 10}
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:              auth,
		className:         "user",
		Where:             types.M{},
		restOptions:       types.M{"skip": 10},
		findOptions:       types.M{"skip": 10},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = types.M{"limit": 10}
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:              auth,
		className:         "user",
		Where:             types.M{},
		restOptions:       types.M{"limit": 10},
		findOptions:       types.M{"limit": 10},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = types.M{"order": "post,-user"}
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:              auth,
		className:         "user",
		Where:             types.M{},
		restOptions:       types.M{"order": "post,-user"},
		findOptions:       types.M{"sort": []string{"post", "-user"}},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = types.M{"include": "user.session,name.friend"}
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:        auth,
		className:   "user",
		Where:       types.M{},
		restOptions: types.M{"include": "user.session,name.friend"},
		findOptions: types.M{},
		response:    types.M{},
		doCount:     false,
		include: [][]string{
			[]string{"name"},
			[]string{"name", "friend"},
			[]string{"user"},
			[]string{"user", "session"},
		},
		keys:              []string{},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = types.M{"redirectClassNameForKey": "post"}
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:              auth,
		className:         "user",
		Where:             types.M{},
		restOptions:       types.M{"redirectClassNameForKey": "post"},
		findOptions:       types.M{},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{},
		redirectKey:       "post",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = types.M{"other": "hello"}
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expectErr = errs.E(errs.InvalidJSON, "bad option: other")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = types.M{
		"keys":                    "post,user",
		"count":                   true,
		"skip":                    10,
		"limit":                   10,
		"order":                   "post,-user",
		"include":                 "user.session,name.friend",
		"redirectClassNameForKey": "post",
	}
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:      auth,
		className: "user",
		Where:     types.M{},
		restOptions: types.M{
			"keys":                    "post,user",
			"count":                   true,
			"skip":                    10,
			"limit":                   10,
			"order":                   "post,-user",
			"include":                 "user.session,name.friend",
			"redirectClassNameForKey": "post",
		},
		findOptions: types.M{
			"skip":  10,
			"limit": 10,
			"sort":  []string{"post", "-user"},
		},
		response: types.M{},
		doCount:  true,
		include: [][]string{
			[]string{"name"},
			[]string{"name", "friend"},
			[]string{"user"},
			[]string{"user", "session"},
		},
		keys:              []string{"post", "user", "objectId", "createdAt", "updatedAt"},
		redirectKey:       "post",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
}

func Test_includePath(t *testing.T) {
	var className string
	var object types.M
	var auth *Auth
	var response types.M
	var path []string
	var restOptions types.M
	var err error
	var expect types.M
	/**********************************************************/
	initEnv()
	className = "post"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, object)
	object = types.M{
		"objectId": "2001",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "2002",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	auth = Master()
	response = types.M{
		"results": types.S{
			types.M{
				"objectId": "1001",
				"post": types.M{
					"__type":    "Pointer",
					"className": "post",
					"objectId":  "2001",
				},
			},
			types.M{
				"objectId": "1002",
				"post": types.M{
					"__type":    "Pointer",
					"className": "post",
					"objectId":  "2002",
				},
			},
		},
	}
	path = []string{"post"}
	err = includePath(auth, response, path, nil)
	expect = types.M{
		"results": types.S{
			types.M{
				"objectId": "1001",
				"post": types.M{
					"__type":    "Object",
					"className": "post",
					"objectId":  "2001",
					"key":       "hello",
				},
			},
			types.M{
				"objectId": "1002",
				"post": types.M{
					"__type":    "Object",
					"className": "post",
					"objectId":  "2002",
					"key":       "hello",
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, response) == false {
		t.Error("expect:", expect, "result:", response)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "post"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, object)
	object = types.M{
		"objectId": "2001",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	className = "postEx"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, object)
	object = types.M{
		"objectId": "3001",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	auth = Master()
	response = types.M{
		"results": types.S{
			types.M{
				"objectId": "1001",
				"post": types.M{
					"__type":    "Pointer",
					"className": "post",
					"objectId":  "2001",
				},
			},
			types.M{
				"objectId": "1002",
				"post": types.M{
					"__type":    "Pointer",
					"className": "postEx",
					"objectId":  "3001",
				},
			},
		},
	}
	path = []string{"post"}
	err = includePath(auth, response, path, nil)
	expect = types.M{
		"results": types.S{
			types.M{
				"objectId": "1001",
				"post": types.M{
					"__type":    "Object",
					"className": "post",
					"objectId":  "2001",
					"key":       "hello",
				},
			},
			types.M{
				"objectId": "1002",
				"post": types.M{
					"__type":    "Object",
					"className": "postEx",
					"objectId":  "3001",
					"key":       "hello",
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, response) == false {
		t.Error("expect:", expect, "result:", response)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "user"
	object = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, object)
	object = types.M{
		"objectId": "4001",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "4002",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	auth = Master()
	response = types.M{
		"results": types.S{
			types.M{
				"objectId": "1001",
				"post": types.M{
					"id": "1",
					"user": types.M{
						"__type":    "Pointer",
						"className": "user",
						"objectId":  "4001",
					},
				},
			},
			types.M{
				"objectId": "1002",
				"post": types.M{
					"id": "2",
					"user": types.M{
						"__type":    "Pointer",
						"className": "user",
						"objectId":  "4002",
					},
				},
			},
		},
	}
	path = []string{"post", "user"}
	err = includePath(auth, response, path, nil)
	expect = types.M{
		"results": types.S{
			types.M{
				"objectId": "1001",
				"post": types.M{
					"id": "1",
					"user": types.M{
						"__type":    "Object",
						"className": "user",
						"objectId":  "4001",
						"key":       "hello",
					},
				},
			},
			types.M{
				"objectId": "1002",
				"post": types.M{
					"id": "2",
					"user": types.M{
						"__type":    "Object",
						"className": "user",
						"objectId":  "4002",
						"key":       "hello",
					},
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, response) == false {
		t.Error("expect:", expect, "result:", response)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "_User"
	object = types.M{
		"fields": types.M{
			"username":     types.M{"type": "String"},
			"sessionToken": types.M{"type": "String"},
			"authData":     types.M{"type": "Object"},
		},
	}
	orm.Adapter.CreateClass(className, object)
	object = types.M{
		"objectId":     "2001",
		"username":     "joe",
		"sessionToken": "abc",
		"authData": types.M{
			"facebook": types.M{
				"id": "1024",
			},
		},
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId":     "2002",
		"username":     "jack",
		"sessionToken": "abc",
		"authData": types.M{
			"facebook": types.M{
				"id": "1024",
			},
		},
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	auth = Master()
	response = types.M{
		"results": types.S{
			types.M{
				"objectId": "1001",
				"user": types.M{
					"__type":    "Pointer",
					"className": "_User",
					"objectId":  "2001",
				},
			},
			types.M{
				"objectId": "1002",
				"user": types.M{
					"__type":    "Pointer",
					"className": "_User",
					"objectId":  "2002",
				},
			},
		},
	}
	path = []string{"user"}
	err = includePath(auth, response, path, nil)
	expect = types.M{
		"results": types.S{
			types.M{
				"objectId": "1001",
				"user": types.M{
					"__type":    "Object",
					"className": "_User",
					"objectId":  "2001",
					"username":  "joe",
					"authData": types.M{
						"facebook": types.M{
							"id": "1024",
						},
					},
				},
			},
			types.M{
				"objectId": "1002",
				"user": types.M{
					"__type":    "Object",
					"className": "_User",
					"objectId":  "2002",
					"username":  "jack",
					"authData": types.M{
						"facebook": types.M{
							"id": "1024",
						},
					},
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, response) == false {
		t.Error("expect:", expect, "result:", response)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "_User"
	object = types.M{
		"fields": types.M{
			"username":     types.M{"type": "String"},
			"sessionToken": types.M{"type": "String"},
			"authData":     types.M{"type": "Object"},
		},
	}
	orm.Adapter.CreateClass(className, object)
	object = types.M{
		"objectId":     "2001",
		"username":     "joe",
		"sessionToken": "abc",
		"authData": types.M{
			"facebook": types.M{
				"id": "1024",
			},
		},
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId":     "2002",
		"username":     "jack",
		"sessionToken": "abc",
		"authData": types.M{
			"facebook": types.M{
				"id": "1024",
			},
		},
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	auth = Nobody()
	response = types.M{
		"results": types.S{
			types.M{
				"objectId": "1001",
				"user": types.M{
					"__type":    "Pointer",
					"className": "_User",
					"objectId":  "2001",
				},
			},
			types.M{
				"objectId": "1002",
				"user": types.M{
					"__type":    "Pointer",
					"className": "_User",
					"objectId":  "2002",
				},
			},
		},
	}
	path = []string{"user"}
	err = includePath(auth, response, path, nil)
	expect = types.M{
		"results": types.S{
			types.M{
				"objectId": "1001",
				"user": types.M{
					"__type":    "Object",
					"className": "_User",
					"objectId":  "2001",
					"username":  "joe",
				},
			},
			types.M{
				"objectId": "1002",
				"user": types.M{
					"__type":    "Object",
					"className": "_User",
					"objectId":  "2002",
					"username":  "jack",
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, response) == false {
		t.Error("expect:", expect, "result:", response)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "user"
	object = types.M{
		"fields": types.M{
			"key":  types.M{"type": "String"},
			"name": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, object)
	object = types.M{
		"objectId": "4001",
		"key":      "hello",
		"name":     "joe",
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "4002",
		"key":      "hello",
		"name":     "jack",
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	auth = Master()
	response = types.M{
		"results": types.S{
			types.M{
				"objectId": "1001",
				"post": types.M{
					"id": "1",
					"user": types.M{
						"__type":    "Pointer",
						"className": "user",
						"objectId":  "4001",
					},
				},
			},
			types.M{
				"objectId": "1002",
				"post": types.M{
					"id": "2",
					"user": types.M{
						"__type":    "Pointer",
						"className": "user",
						"objectId":  "4002",
					},
				},
			},
		},
	}
	path = []string{"post", "user"}
	restOptions = types.M{"keys": "post.user.key"}
	err = includePath(auth, response, path, restOptions)
	expect = types.M{
		"results": types.S{
			types.M{
				"objectId": "1001",
				"post": types.M{
					"id": "1",
					"user": types.M{
						"__type":    "Object",
						"className": "user",
						"objectId":  "4001",
						"key":       "hello",
					},
				},
			},
			types.M{
				"objectId": "1002",
				"post": types.M{
					"id": "2",
					"user": types.M{
						"__type":    "Object",
						"className": "user",
						"objectId":  "4002",
						"key":       "hello",
					},
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, response) == false {
		t.Error("expect:", expect, "result:", response)
	}
	orm.TomatoDBController.DeleteEverything()
	/**********************************************************/
	initEnv()
	className = "user"
	object = types.M{
		"fields": types.M{
			"key":  types.M{"type": "String"},
			"name": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, object)
	object = types.M{
		"objectId": "4001",
		"key":      "hello",
		"name":     "joe",
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	object = types.M{
		"objectId": "4002",
		"key":      "hello",
		"name":     "jack",
	}
	orm.Adapter.CreateObject(className, types.M{}, object)
	auth = Master()
	response = types.M{
		"results": types.S{
			types.M{
				"objectId": "1001",
				"post": types.M{
					"id": "1",
					"user": types.M{
						"__type":    "Pointer",
						"className": "user",
						"objectId":  "4001",
					},
				},
			},
			types.M{
				"objectId": "1002",
				"post": types.M{
					"id": "2",
					"user": types.M{
						"__type":    "Pointer",
						"className": "user",
						"objectId":  "4002",
					},
				},
			},
		},
	}
	path = []string{"post", "user"}
	restOptions = types.M{"keys": "post.user.key,post.user.name"}
	err = includePath(auth, response, path, restOptions)
	expect = types.M{
		"results": types.S{
			types.M{
				"objectId": "1001",
				"post": types.M{
					"id": "1",
					"user": types.M{
						"__type":    "Object",
						"className": "user",
						"objectId":  "4001",
						"key":       "hello",
						"name":      "joe",
					},
				},
			},
			types.M{
				"objectId": "1002",
				"post": types.M{
					"id": "2",
					"user": types.M{
						"__type":    "Object",
						"className": "user",
						"objectId":  "4002",
						"key":       "hello",
						"name":      "jack",
					},
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, response) == false {
		t.Error("expect:", expect, "result:", response)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_findPointers(t *testing.T) {
	var object interface{}
	var path []string
	var result []types.M
	var expect []types.M
	/**********************************************************/
	object = nil
	path = nil
	result = findPointers(object, path)
	expect = []types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	object = "hello"
	path = nil
	result = findPointers(object, path)
	expect = []types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	object = types.M{
		"key": "hello",
	}
	path = nil
	result = findPointers(object, path)
	expect = []types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	object = types.M{
		"__type":   "Pointer",
		"objectId": "1001",
	}
	path = nil
	result = findPointers(object, path)
	expect = []types.M{
		types.M{
			"__type":   "Pointer",
			"objectId": "1001",
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	object = types.M{
		"key": "hello",
	}
	path = []string{"post"}
	result = findPointers(object, path)
	expect = []types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	object = types.M{
		"post": types.M{
			"__type":   "Pointer",
			"objectId": "1001",
		},
	}
	path = []string{"post"}
	result = findPointers(object, path)
	expect = []types.M{
		types.M{
			"__type":   "Pointer",
			"objectId": "1001",
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	object = types.S{
		types.M{
			"__type":   "Pointer",
			"objectId": "1001",
		},
		types.M{
			"__type":   "Pointer",
			"objectId": "1002",
		},
	}
	path = []string{}
	result = findPointers(object, path)
	expect = []types.M{
		types.M{
			"__type":   "Pointer",
			"objectId": "1001",
		},
		types.M{
			"__type":   "Pointer",
			"objectId": "1002",
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	object = types.M{
		"post": types.S{
			types.M{
				"__type":   "Pointer",
				"objectId": "1001",
			},
			types.M{
				"__type":   "Pointer",
				"objectId": "1002",
			},
		},
	}
	path = []string{"post"}
	result = findPointers(object, path)
	expect = []types.M{
		types.M{
			"__type":   "Pointer",
			"objectId": "1001",
		},
		types.M{
			"__type":   "Pointer",
			"objectId": "1002",
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	object = types.M{
		"post": types.M{
			"user": types.M{
				"__type":   "Pointer",
				"objectId": "1001",
			},
		},
	}
	path = []string{"post", "user"}
	result = findPointers(object, path)
	expect = []types.M{
		types.M{
			"__type":   "Pointer",
			"objectId": "1001",
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_replacePointers(t *testing.T) {
	var pointers []types.M
	var replace types.M
	var expect []types.M
	/**********************************************************/
	pointers = nil
	replace = nil
	replacePointers(pointers, replace)
	expect = nil
	if reflect.DeepEqual(expect, pointers) == false {
		t.Error("expect:", expect, "result:", pointers)
	}
	/**********************************************************/
	pointers = []types.M{nil}
	replace = types.M{
		"1001": types.M{
			"key": "hello",
		},
	}
	replacePointers(pointers, replace)
	expect = []types.M{nil}
	if reflect.DeepEqual(expect, pointers) == false {
		t.Error("expect:", expect, "result:", pointers)
	}
	/**********************************************************/
	pointers = []types.M{
		types.M{
			"post": "post",
		},
	}
	replace = types.M{
		"1001": types.M{
			"key": "hello",
		},
	}
	replacePointers(pointers, replace)
	expect = []types.M{
		types.M{
			"post": "post",
		},
	}
	if reflect.DeepEqual(expect, pointers) == false {
		t.Error("expect:", expect, "result:", pointers)
	}
	/**********************************************************/
	pointers = []types.M{
		types.M{
			"objectId": "1002",
			"post":     "post",
		},
	}
	replace = types.M{
		"1001": types.M{
			"key": "hello",
		},
	}
	replacePointers(pointers, replace)
	expect = []types.M{
		types.M{
			"objectId": "1002",
			"post":     "post",
		},
	}
	if reflect.DeepEqual(expect, pointers) == false {
		t.Error("expect:", expect, "result:", pointers)
	}
	/**********************************************************/
	pointers = []types.M{
		types.M{
			"objectId": "1001",
			"post":     "post",
		},
	}
	replace = types.M{
		"1001": types.M{
			"key": "hello",
		},
	}
	replacePointers(pointers, replace)
	expect = []types.M{
		types.M{
			"objectId": "1001",
			"post":     "post",
			"key":      "hello",
		},
	}
	if reflect.DeepEqual(expect, pointers) == false {
		t.Error("expect:", expect, "result:", pointers)
	}
	/**********************************************************/
	pointers = []types.M{
		types.M{
			"objectId": "1001",
			"post":     "post",
		},
		types.M{
			"objectId": "1002",
			"post":     "post",
		},
	}
	replace = types.M{
		"1001": types.M{
			"key": "hello",
		},
		"1002": types.M{
			"key": "hello",
		},
	}
	replacePointers(pointers, replace)
	expect = []types.M{
		types.M{
			"objectId": "1001",
			"post":     "post",
			"key":      "hello",
		},
		types.M{
			"objectId": "1002",
			"post":     "post",
			"key":      "hello",
		},
	}
	if reflect.DeepEqual(expect, pointers) == false {
		t.Error("expect:", expect, "result:", pointers)
	}
}

func Test_findObjectWithKey(t *testing.T) {
	var root interface{}
	var key string
	var result types.M
	var expect types.M
	/**********************************************************/
	root = nil
	key = "post"
	result = findObjectWithKey(root, key)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	root = types.M{
		"post": "hello",
	}
	key = "post"
	result = findObjectWithKey(root, key)
	expect = types.M{
		"post": "hello",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	root = types.M{
		"key": "hello",
	}
	key = "post"
	result = findObjectWithKey(root, key)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	root = types.M{
		"key": types.M{
			"post": "hello",
		},
	}
	key = "post"
	result = findObjectWithKey(root, key)
	expect = types.M{
		"post": "hello",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	root = types.M{
		"key": types.M{
			"key": "hello",
		},
	}
	key = "post"
	result = findObjectWithKey(root, key)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	root = types.S{
		types.M{
			"post": "hello",
		},
	}
	key = "post"
	result = findObjectWithKey(root, key)
	expect = types.M{
		"post": "hello",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	root = types.S{
		types.M{
			"key": "hello",
		},
	}
	key = "post"
	result = findObjectWithKey(root, key)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	root = types.S{
		types.M{
			"key": types.M{
				"post": "hello",
			},
		},
	}
	key = "post"
	result = findObjectWithKey(root, key)
	expect = types.M{
		"post": "hello",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_transformSelect(t *testing.T) {
	var selectObject types.M
	var key string
	var objects []types.M
	var expect types.M
	/**********************************************************/
	selectObject = nil
	key = "user"
	objects = nil
	transformSelect(selectObject, key, objects)
	expect = nil
	if reflect.DeepEqual(expect, selectObject) == false {
		t.Error("expect:", expect, "result:", selectObject)
	}
	/**********************************************************/
	selectObject = types.M{}
	key = "user"
	objects = nil
	transformSelect(selectObject, key, objects)
	expect = types.M{}
	if reflect.DeepEqual(expect, selectObject) == false {
		t.Error("expect:", expect, "result:", selectObject)
	}
	/**********************************************************/
	selectObject = types.M{
		"$select": "string",
	}
	key = "user"
	objects = nil
	transformSelect(selectObject, key, objects)
	expect = types.M{
		"$in": types.S{},
	}
	if reflect.DeepEqual(expect, selectObject) == false {
		t.Error("expect:", expect, "result:", selectObject)
	}
	/**********************************************************/
	selectObject = types.M{
		"$select": "string",
	}
	key = "user"
	objects = []types.M{}
	transformSelect(selectObject, key, objects)
	expect = types.M{
		"$in": types.S{},
	}
	if reflect.DeepEqual(expect, selectObject) == false {
		t.Error("expect:", expect, "result:", selectObject)
	}
	/**********************************************************/
	selectObject = types.M{
		"$select": "string",
	}
	key = "user"
	objects = []types.M{
		types.M{
			"user": "1001",
		},
		types.M{
			"key": "1002",
		},
	}
	transformSelect(selectObject, key, objects)
	expect = types.M{
		"$in": types.S{
			"1001",
		},
	}
	if reflect.DeepEqual(expect, selectObject) == false {
		t.Error("expect:", expect, "result:", selectObject)
	}
	/**********************************************************/
	selectObject = types.M{
		"$select": "string",
	}
	key = "user"
	objects = []types.M{
		types.M{
			"user": "1001",
		},
		types.M{
			"user": "1002",
		},
	}
	transformSelect(selectObject, key, objects)
	expect = types.M{
		"$in": types.S{
			"1001",
			"1002",
		},
	}
	if reflect.DeepEqual(expect, selectObject) == false {
		t.Error("expect:", expect, "result:", selectObject)
	}
	/**********************************************************/
	selectObject = types.M{
		"$select": "string",
		"$in": types.S{
			"1003",
		},
	}
	key = "user"
	objects = []types.M{
		types.M{
			"user": "1001",
		},
		types.M{
			"user": "1002",
		},
	}
	transformSelect(selectObject, key, objects)
	expect = types.M{
		"$in": types.S{
			"1003",
			"1001",
			"1002",
		},
	}
	if reflect.DeepEqual(expect, selectObject) == false {
		t.Error("expect:", expect, "result:", selectObject)
	}
}

func Test_transformDontSelect(t *testing.T) {
	var dontSelectObject types.M
	var key string
	var objects []types.M
	var expect types.M
	/**********************************************************/
	dontSelectObject = nil
	key = "user"
	objects = nil
	transformDontSelect(dontSelectObject, key, objects)
	expect = nil
	if reflect.DeepEqual(expect, dontSelectObject) == false {
		t.Error("expect:", expect, "result:", dontSelectObject)
	}
	/**********************************************************/
	dontSelectObject = types.M{}
	key = "user"
	objects = nil
	transformDontSelect(dontSelectObject, key, objects)
	expect = types.M{}
	if reflect.DeepEqual(expect, dontSelectObject) == false {
		t.Error("expect:", expect, "result:", dontSelectObject)
	}
	/**********************************************************/
	dontSelectObject = types.M{
		"$dontSelect": "string",
	}
	key = "user"
	objects = nil
	transformDontSelect(dontSelectObject, key, objects)
	expect = types.M{
		"$nin": types.S{},
	}
	if reflect.DeepEqual(expect, dontSelectObject) == false {
		t.Error("expect:", expect, "result:", dontSelectObject)
	}
	/**********************************************************/
	dontSelectObject = types.M{
		"$dontSelect": "string",
	}
	key = "user"
	objects = []types.M{}
	transformDontSelect(dontSelectObject, key, objects)
	expect = types.M{
		"$nin": types.S{},
	}
	if reflect.DeepEqual(expect, dontSelectObject) == false {
		t.Error("expect:", expect, "result:", dontSelectObject)
	}
	/**********************************************************/
	dontSelectObject = types.M{
		"$dontSelect": "string",
	}
	key = "user"
	objects = []types.M{
		types.M{
			"user": "1001",
		},
		types.M{
			"key": "1002",
		},
	}
	transformDontSelect(dontSelectObject, key, objects)
	expect = types.M{
		"$nin": types.S{
			"1001",
		},
	}
	if reflect.DeepEqual(expect, dontSelectObject) == false {
		t.Error("expect:", expect, "result:", dontSelectObject)
	}
	/**********************************************************/
	dontSelectObject = types.M{
		"$dontSelect": "string",
	}
	key = "user"
	objects = []types.M{
		types.M{
			"user": "1001",
		},
		types.M{
			"user": "1002",
		},
	}
	transformDontSelect(dontSelectObject, key, objects)
	expect = types.M{
		"$nin": types.S{
			"1001",
			"1002",
		},
	}
	if reflect.DeepEqual(expect, dontSelectObject) == false {
		t.Error("expect:", expect, "result:", dontSelectObject)
	}
	/**********************************************************/
	dontSelectObject = types.M{
		"$dontSelect": "string",
		"$nin": types.S{
			"1003",
		},
	}
	key = "user"
	objects = []types.M{
		types.M{
			"user": "1001",
		},
		types.M{
			"user": "1002",
		},
	}
	transformDontSelect(dontSelectObject, key, objects)
	expect = types.M{
		"$nin": types.S{
			"1003",
			"1001",
			"1002",
		},
	}
	if reflect.DeepEqual(expect, dontSelectObject) == false {
		t.Error("expect:", expect, "result:", dontSelectObject)
	}
}

func Test_transformInQuery(t *testing.T) {
	var inQueryObject types.M
	var className string
	var results []types.M
	var expect types.M
	/**********************************************************/
	inQueryObject = nil
	className = "user"
	results = nil
	transformInQuery(inQueryObject, className, results)
	expect = nil
	if reflect.DeepEqual(expect, inQueryObject) == false {
		t.Error("expect:", expect, "result:", inQueryObject)
	}
	/**********************************************************/
	inQueryObject = types.M{}
	className = "user"
	results = nil
	transformInQuery(inQueryObject, className, results)
	expect = types.M{}
	if reflect.DeepEqual(expect, inQueryObject) == false {
		t.Error("expect:", expect, "result:", inQueryObject)
	}
	/**********************************************************/
	inQueryObject = types.M{
		"$inQuery": "string",
	}
	className = "user"
	results = nil
	transformInQuery(inQueryObject, className, results)
	expect = types.M{
		"$in": types.S{},
	}
	if reflect.DeepEqual(expect, inQueryObject) == false {
		t.Error("expect:", expect, "result:", inQueryObject)
	}
	/**********************************************************/
	inQueryObject = types.M{
		"$inQuery": "string",
	}
	className = "user"
	results = []types.M{}
	transformInQuery(inQueryObject, className, results)
	expect = types.M{
		"$in": types.S{},
	}
	if reflect.DeepEqual(expect, inQueryObject) == false {
		t.Error("expect:", expect, "result:", inQueryObject)
	}
	/**********************************************************/
	inQueryObject = types.M{
		"$inQuery": "string",
	}
	className = "user"
	results = []types.M{
		types.M{
			"objectId": "1001",
		},
		types.M{
			"key": "1002",
		},
	}
	transformInQuery(inQueryObject, className, results)
	expect = types.M{
		"$in": types.S{
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1001",
			},
		},
	}
	if reflect.DeepEqual(expect, inQueryObject) == false {
		t.Error("expect:", expect, "result:", inQueryObject)
	}
	/**********************************************************/
	inQueryObject = types.M{
		"$inQuery": "string",
	}
	className = "user"
	results = []types.M{
		types.M{
			"objectId": "1001",
		},
		types.M{
			"objectId": "1002",
		},
	}
	transformInQuery(inQueryObject, className, results)
	expect = types.M{
		"$in": types.S{
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1001",
			},
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1002",
			},
		},
	}
	if reflect.DeepEqual(expect, inQueryObject) == false {
		t.Error("expect:", expect, "result:", inQueryObject)
	}
	/**********************************************************/
	inQueryObject = types.M{
		"$inQuery": "string",
		"$in": types.S{
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1003",
			},
		},
	}
	className = "user"
	results = []types.M{
		types.M{
			"objectId": "1001",
		},
		types.M{
			"objectId": "1002",
		},
	}
	transformInQuery(inQueryObject, className, results)
	expect = types.M{
		"$in": types.S{
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1003",
			},
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1001",
			},
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1002",
			},
		},
	}
	if reflect.DeepEqual(expect, inQueryObject) == false {
		t.Error("expect:", expect, "result:", inQueryObject)
	}
}

func Test_transformNotInQuery(t *testing.T) {
	var notInQueryObject types.M
	var className string
	var results []types.M
	var expect types.M
	/**********************************************************/
	notInQueryObject = nil
	className = "user"
	results = nil
	transformNotInQuery(notInQueryObject, className, results)
	expect = nil
	if reflect.DeepEqual(expect, notInQueryObject) == false {
		t.Error("expect:", expect, "result:", notInQueryObject)
	}
	/**********************************************************/
	notInQueryObject = types.M{}
	className = "user"
	results = nil
	transformNotInQuery(notInQueryObject, className, results)
	expect = types.M{}
	if reflect.DeepEqual(expect, notInQueryObject) == false {
		t.Error("expect:", expect, "result:", notInQueryObject)
	}
	/**********************************************************/
	notInQueryObject = types.M{
		"$notInQuery": "string",
	}
	className = "user"
	results = nil
	transformNotInQuery(notInQueryObject, className, results)
	expect = types.M{
		"$nin": types.S{},
	}
	if reflect.DeepEqual(expect, notInQueryObject) == false {
		t.Error("expect:", expect, "result:", notInQueryObject)
	}
	/**********************************************************/
	notInQueryObject = types.M{
		"$notInQuery": "string",
	}
	className = "user"
	results = []types.M{}
	transformNotInQuery(notInQueryObject, className, results)
	expect = types.M{
		"$nin": types.S{},
	}
	if reflect.DeepEqual(expect, notInQueryObject) == false {
		t.Error("expect:", expect, "result:", notInQueryObject)
	}
	/**********************************************************/
	notInQueryObject = types.M{
		"$notInQuery": "string",
	}
	className = "user"
	results = []types.M{
		types.M{
			"objectId": "1001",
		},
		types.M{
			"key": "1002",
		},
	}
	transformNotInQuery(notInQueryObject, className, results)
	expect = types.M{
		"$nin": types.S{
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1001",
			},
		},
	}
	if reflect.DeepEqual(expect, notInQueryObject) == false {
		t.Error("expect:", expect, "result:", notInQueryObject)
	}
	/**********************************************************/
	notInQueryObject = types.M{
		"$notInQuery": "string",
	}
	className = "user"
	results = []types.M{
		types.M{
			"objectId": "1001",
		},
		types.M{
			"objectId": "1002",
		},
	}
	transformNotInQuery(notInQueryObject, className, results)
	expect = types.M{
		"$nin": types.S{
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1001",
			},
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1002",
			},
		},
	}
	if reflect.DeepEqual(expect, notInQueryObject) == false {
		t.Error("expect:", expect, "result:", notInQueryObject)
	}
	/**********************************************************/
	notInQueryObject = types.M{
		"$notInQuery": "string",
		"$nin": types.S{
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1003",
			},
		},
	}
	className = "user"
	results = []types.M{
		types.M{
			"objectId": "1001",
		},
		types.M{
			"objectId": "1002",
		},
	}
	transformNotInQuery(notInQueryObject, className, results)
	expect = types.M{
		"$nin": types.S{
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1003",
			},
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1001",
			},
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1002",
			},
		},
	}
	if reflect.DeepEqual(expect, notInQueryObject) == false {
		t.Error("expect:", expect, "result:", notInQueryObject)
	}
}

func initEnv() {
	orm.InitOrm(getAdapter())
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
