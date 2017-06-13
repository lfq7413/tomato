package rest

import (
	"reflect"
	"testing"

	"github.com/lfq7413/tomato/cloud"
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

func TestPostgres_Find(t *testing.T) {
	var object, schema types.M
	var className string
	var result, expect types.M
	var err error
	/********************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, schema, object)
	result, err = Find(Master(), className, types.M{}, types.M{}, nil)
	expect = types.M{
		"results": types.S{
			types.M{
				"objectId": "01",
				"key":      "hello",
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/********************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "02",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, schema, object)
	result, err = Find(Master(), className, types.M{}, types.M{"count": true}, nil)
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
		"count": 2,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/********************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, schema, object)
	result, err = Find(Master(), className, types.M{"objectId": "02"}, types.M{}, nil)
	expect = types.M{
		"results": types.S{},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
}

func TestPostgres_Get(t *testing.T) {
	var object, schema types.M
	var className string
	var result, expect types.M
	var err error
	/********************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, schema, object)
	result, err = Get(Master(), className, "01", types.M{}, nil)
	expect = types.M{
		"results": types.S{
			types.M{
				"objectId": "01",
				"key":      "hello",
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/********************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, schema, object)
	result, err = Get(Master(), className, "02", types.M{}, nil)
	expect = types.M{
		"results": types.S{},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
}

func TestPostgres_Delete(t *testing.T) {
	var object, schema types.M
	var auth *Auth
	var className, objectID string
	var err, expect error
	/********************************************************/
	initPostgresEnv()
	className = "_User"
	auth = Nobody()
	objectID = "01"
	err = Delete(auth, className, objectID)
	expect = errs.E(errs.SessionMissing, "insufficient auth to delete user")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/********************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, schema, object)
	auth = Master()
	objectID = "01"
	err = Delete(auth, className, objectID)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/********************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, schema, object)
	auth = Master()
	objectID = "02"
	err = Delete(auth, className, objectID)
	expect = errs.E(errs.ObjectNotFound, "Object not found.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/********************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, schema, object)
	cloud.AfterDelete(className, func(cloud.TriggerRequest, cloud.Response) {})
	auth = Master()
	objectID = "01"
	err = Delete(auth, className, objectID)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/********************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"key":      types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"key":      "hello",
	}
	orm.Adapter.CreateObject(className, schema, object)
	cloud.AfterDelete(className, func(cloud.TriggerRequest, cloud.Response) {})
	auth = Master()
	objectID = "02"
	err = Delete(auth, className, objectID)
	expect = errs.E(errs.ObjectNotFound, "Object not found for delete.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
}

func TestPostgres_Create(t *testing.T) {
	var auth *Auth
	var className string
	var object types.M
	var result types.M
	var err error
	/********************************************************/
	initPostgresEnv()
	auth = Master()
	className = "user"
	object = types.M{
		"name": "joe",
		"age":  "12",
	}
	config.TConfig.ServerURL = "http://127.0.0.1/v1"
	result, err = Create(auth, className, object, nil)
	if err != nil || result == nil {
		t.Error("expect:", nil, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
}

func TestPostgres_Update(t *testing.T) {
	var auth *Auth
	var className, objectID string
	var object, schema types.M
	var result types.M
	var err, expectErr error
	var results, expects types.S
	/********************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId":  types.M{"type": "String"},
			"name":      types.M{"type": "String"},
			"updatedAt": types.M{"type": "Date"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"name":     "joe",
	}
	orm.Adapter.CreateObject(className, schema, object)
	auth = Master()
	objectID = "01"
	object = types.M{
		"name": "jack",
	}
	result, err = Update(auth, className, objectID, object, nil)
	if err != nil || result == nil {
		t.Error("expect:", nil, "result:", result, err)
	}
	results, _ = orm.TomatoDBController.Find(className, types.M{}, types.M{})
	expects = types.S{
		types.M{
			"objectId": "01",
			"name":     "jack",
		},
	}
	delete(utils.M(results[0]), "updatedAt")
	if reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results)
	}
	orm.TomatoDBController.DeleteEverything()
	/********************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId":  types.M{"type": "String"},
			"name":      types.M{"type": "String"},
			"updatedAt": types.M{"type": "Date"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"name":     "joe",
	}
	orm.Adapter.CreateObject(className, schema, object)
	auth = Master()
	objectID = "02"
	object = types.M{
		"name": "jack",
	}
	result, err = Update(auth, className, objectID, object, nil)
	expectErr = errs.E(errs.ObjectNotFound, "Object not found.")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/********************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId":  types.M{"type": "String"},
			"name":      types.M{"type": "String"},
			"updatedAt": types.M{"type": "Date"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"name":     "joe",
	}
	orm.Adapter.CreateObject(className, schema, object)
	cloud.AfterSave(className, func(cloud.TriggerRequest, cloud.Response) {})
	auth = Master()
	objectID = "01"
	object = types.M{
		"name": "jack",
	}
	result, err = Update(auth, className, objectID, object, nil)
	if err != nil || result == nil {
		t.Error("expect:", nil, "result:", result, err)
	}
	results, _ = orm.TomatoDBController.Find(className, types.M{}, types.M{})
	expects = types.S{
		types.M{
			"objectId": "01",
			"name":     "jack",
		},
	}
	delete(utils.M(results[0]), "updatedAt")
	if reflect.DeepEqual(expects, results) == false {
		t.Error("expect:", expects, "result:", results)
	}
	orm.TomatoDBController.DeleteEverything()
	/********************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"name":     types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "01",
		"name":     "joe",
	}
	orm.Adapter.CreateObject(className, schema, object)
	cloud.AfterSave(className, func(cloud.TriggerRequest, cloud.Response) {})
	auth = Master()
	objectID = "02"
	object = types.M{
		"name": "jack",
	}
	result, err = Update(auth, className, objectID, object, nil)
	expectErr = errs.E(errs.ObjectNotFound, "Object not found for update.")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
}
