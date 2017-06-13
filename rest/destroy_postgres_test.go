package rest

import (
	"reflect"
	"testing"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/types"
)

func TestPostgres_Destroy(t *testing.T) {
	var schema, object types.M
	var auth *Auth
	var className string
	var query types.M
	var originalData types.M
	var d *Destroy
	var err, expectErr error
	var results, expect types.S
	/*****************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"username": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"username": "joe",
	}
	orm.Adapter.CreateObject(className, schema, object)
	auth = Master()
	className = "user"
	query = types.M{"objectId": "1001"}
	originalData = types.M{"username": "joe"}
	d = NewDestroy(auth, className, query, originalData)
	err = d.Execute()
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	results, _ = orm.TomatoDBController.Find(className, types.M{}, types.M{})
	expect = types.S{}
	if reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results)
	}
	orm.TomatoDBController.DeleteEverything()
	/*****************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"username": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"username": "joe",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "1002",
		"username": "joe",
	}
	orm.Adapter.CreateObject(className, schema, object)
	auth = Master()
	className = "user"
	query = types.M{"objectId": "1001"}
	originalData = types.M{"username": "joe"}
	d = NewDestroy(auth, className, query, originalData)
	err = d.Execute()
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
	results, _ = orm.TomatoDBController.Find(className, types.M{}, types.M{})
	expect = types.S{
		types.M{
			"objectId": "1002",
			"username": "joe",
		},
	}
	if reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results)
	}
	orm.TomatoDBController.DeleteEverything()
	/*****************************************************/
	initPostgresEnv()
	className = "user"
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"username": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass(className, schema)
	object = types.M{
		"objectId": "1001",
		"username": "joe",
	}
	orm.Adapter.CreateObject(className, schema, object)
	object = types.M{
		"objectId": "1002",
		"username": "joe",
	}
	orm.Adapter.CreateObject(className, schema, object)
	auth = Master()
	className = "user"
	query = types.M{"objectId": "1003"}
	originalData = types.M{"username": "joe"}
	d = NewDestroy(auth, className, query, originalData)
	err = d.Execute()
	expectErr = errs.E(errs.ObjectNotFound, "Object not found.")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", err)
	}
	results, _ = orm.TomatoDBController.Find(className, types.M{}, types.M{})
	expect = types.S{
		types.M{
			"objectId": "1001",
			"username": "joe",
		},
		types.M{
			"objectId": "1002",
			"username": "joe",
		},
	}
	if reflect.DeepEqual(expect, results) == false {
		t.Error("expect:", expect, "result:", results)
	}
	orm.TomatoDBController.DeleteEverything()
}
