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

func Test_enforceRoleSecurity(t *testing.T) {
	var method, className string
	var auth *Auth
	var err, expect error
	/********************************************************/
	method = "delete"
	className = "_Installation"
	auth = Nobody()
	err = enforceRoleSecurity(method, className, auth)
	expect = errs.E(errs.OperationForbidden, "Clients aren't allowed to perform the delete operation on the installation collection.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/********************************************************/
	method = "find"
	className = "_Installation"
	auth = Nobody()
	err = enforceRoleSecurity(method, className, auth)
	expect = errs.E(errs.OperationForbidden, "Clients aren't allowed to perform the find operation on the installation collection.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	/********************************************************/
	method = "find"
	className = "_Installation"
	auth = Master()
	err = enforceRoleSecurity(method, className, auth)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
}

func Test_Find(t *testing.T) {
	var object, schema types.M
	var className string
	var result, expect types.M
	var err error
	/********************************************************/
	initEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
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
	initEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
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
	initEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
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

func Test_Get(t *testing.T) {
	var object, schema types.M
	var className string
	var result, expect types.M
	var err error
	/********************************************************/
	initEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
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
	initEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
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

func Test_Delete(t *testing.T) {
	var object, schema types.M
	var auth *Auth
	var className, objectID string
	var err, expect error
	/********************************************************/
	initEnv()
	className = "_User"
	auth = Nobody()
	objectID = "01"
	err = Delete(auth, className, objectID, nil)
	expect = errs.E(errs.SessionMissing, "insufficient auth to delete user")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/********************************************************/
	initEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
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
	err = Delete(auth, className, objectID, nil)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/********************************************************/
	initEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
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
	err = Delete(auth, className, objectID, nil)
	expect = errs.E(errs.ObjectNotFound, "Object not found.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/********************************************************/
	initEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
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
	err = Delete(auth, className, objectID, nil)
	expect = nil
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
	/********************************************************/
	initEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"key": types.M{"type": "String"},
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
	err = Delete(auth, className, objectID, nil)
	expect = errs.E(errs.ObjectNotFound, "Object not found for delete.")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "result:", err)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_Create(t *testing.T) {
	var auth *Auth
	var className string
	var object types.M
	var result types.M
	var err error
	/********************************************************/
	initEnv()
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

func Test_Update(t *testing.T) {
	var auth *Auth
	var className, objectID string
	var object, schema types.M
	var result types.M
	var err, expectErr error
	var results, expects types.S
	/********************************************************/
	initEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"name": types.M{"type": "String"},
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
		t.Error("expect:", nil, "result:", result)
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
	initEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"name": types.M{"type": "String"},
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
	initEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"name": types.M{"type": "String"},
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
		t.Error("expect:", nil, "result:", result)
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
	initEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"name": types.M{"type": "String"},
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
