package mongo

import (
	"reflect"
	"testing"
	"time"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

func Test_transformKey(t *testing.T) {
	tf := NewTransform()
	var fieldName string
	var schema types.M
	var result string
	/*************************************************/
	fieldName = "objectId"
	result = tf.transformKey("", fieldName, schema)
	if result != "_id" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
	/*************************************************/
	fieldName = "createdAt"
	result = tf.transformKey("", fieldName, schema)
	if result != "_created_at" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
	/*************************************************/
	fieldName = "updatedAt"
	result = tf.transformKey("", fieldName, schema)
	if result != "_updated_at" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
	/*************************************************/
	fieldName = "sessionToken"
	result = tf.transformKey("", fieldName, schema)
	if result != "_session_token" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
	/*************************************************/
	schema = nil
	fieldName = "user"
	result = tf.transformKey("", fieldName, schema)
	if result != "user" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
	/*************************************************/
	schema = types.M{
		"fields": "type is string",
	}
	fieldName = "user"
	result = tf.transformKey("", fieldName, schema)
	if result != "user" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
	/*************************************************/
	schema = types.M{
		"fields": types.M{
			"user": "type is string",
		},
	}
	fieldName = "user"
	result = tf.transformKey("", fieldName, schema)
	if result != "user" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
	/*************************************************/
	schema = types.M{
		"fields": types.M{
			"user": types.M{
				"__type": "String",
			},
		},
	}
	fieldName = "user"
	result = tf.transformKey("", fieldName, schema)
	if result != "user" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
	/*************************************************/
	schema = types.M{
		"fields": types.M{
			"user": types.M{
				"__type": "Pointer",
			},
		},
	}
	fieldName = "user"
	result = tf.transformKey("", fieldName, schema)
	if result != "_p_user" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
}

func Test_transformKeyValueForUpdate(t *testing.T) {
	// TODO
}

func Test_transformQueryKeyValue(t *testing.T) {
	// transformWhere
	// TODO
}

func Test_transformConstraint(t *testing.T) {
	tf := NewTransform()
	var constraint interface{}
	var inArray bool
	var result interface{}
	var err error
	var expect interface{}
	/*************************************************/
	constraint = nil
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = nil
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = 1024
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = nil
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{"$lt": 10}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{"$lt": 10}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{"$lt": 10}
	inArray = false
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{"$lt": 10}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{"$lt": types.M{"key": "value"}}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = errs.E(errs.InvalidJSON, "bad atom")
	if reflect.DeepEqual(err, expect) == false || result != nil {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	constraint = types.M{"$lt": types.M{"key": "value"}}
	inArray = false
	result, err = tf.transformConstraint(constraint, inArray)
	expect = errs.E(errs.InvalidJSON, "bad atom")
	if reflect.DeepEqual(err, expect) == false || result != nil {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	constraint = types.M{"$in": types.M{"key": "value"}}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = errs.E(errs.InvalidJSON, "bad "+"$in"+" value")
	if reflect.DeepEqual(err, expect) == false || result != nil {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	constraint = types.M{"$in": types.S{"hello", "world"}}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{"$in": types.S{"hello", "world"}}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{"$in": types.S{"hello", "world"}}
	inArray = false
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{"$in": types.S{"hello", "world"}}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{"$in": types.S{types.M{"key": "value"}}}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = errs.E(errs.InvalidJSON, "bad atom")
	if reflect.DeepEqual(err, expect) == false || result != nil {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	constraint = types.M{"$in": types.S{types.M{"key": "value"}}}
	inArray = false
	result, err = tf.transformConstraint(constraint, inArray)
	expect = errs.E(errs.InvalidJSON, "bad atom")
	if reflect.DeepEqual(err, expect) == false || result != nil {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	constraint = types.M{"$all": types.M{"key": "value"}}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = errs.E(errs.InvalidJSON, "bad "+"$all"+" value")
	if reflect.DeepEqual(err, expect) == false || result != nil {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	constraint = types.M{"$all": types.S{"hello", "world"}}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{"$all": types.S{"hello", "world"}}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{"$regex": 1024}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = errs.E(errs.InvalidJSON, "bad regex")
	if reflect.DeepEqual(err, expect) == false || result != nil {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	constraint = types.M{"$regex": "hello"}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{"$regex": "hello"}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{"$options": "imxs"}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = errs.E(errs.InvalidQuery, "got a bad $options")
	if reflect.DeepEqual(err, expect) == false || result != nil {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	constraint = types.M{"$options": 1024, "$regex": "hello"}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = errs.E(errs.InvalidQuery, "got a bad $options")
	if reflect.DeepEqual(err, expect) == false || result != nil {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	constraint = types.M{"$options": "hello", "$regex": "hello"}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = errs.E(errs.InvalidQuery, "got a bad $options")
	if reflect.DeepEqual(err, expect) == false || result != nil {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	constraint = types.M{"$options": "imxs", "$regex": "hello"}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{"$options": "imxs", "$regex": "hello"}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{"$nearSphere": "hello"}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{"$nearSphere": types.S{0, 0}}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{
		"$nearSphere": types.M{
			"longitude": 20,
			"latitude":  20,
		},
	}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{"$nearSphere": types.S{20, 20}}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{"$maxDistance": 0.26}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{"$maxDistance": 0.26}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{"$maxDistanceInRadians": 0.26}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{"$maxDistance": 0.26}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{"$maxDistanceInMiles": 16.0}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{"$maxDistance": 16.0 / 3959}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{"$maxDistanceInMiles": 16}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{"$maxDistance": 16.0 / 3959}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{"$maxDistanceInKilometers": 16.0}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{"$maxDistance": 16.0 / 6371}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{"$maxDistanceInKilometers": 16}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{"$maxDistance": 16.0 / 6371}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{"$select": "hello"}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = errs.E(errs.CommandUnavailable, "the "+"$select"+" constraint is not supported yet")
	if reflect.DeepEqual(err, expect) == false || result != nil {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	constraint = types.M{"$dontSelect": "hello"}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = errs.E(errs.CommandUnavailable, "the "+"$dontSelect"+" constraint is not supported yet")
	if reflect.DeepEqual(err, expect) == false || result != nil {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	constraint = types.M{"$within": "hello"}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = errs.E(errs.InvalidJSON, "malformatted $within arg")
	if reflect.DeepEqual(err, expect) == false || result != nil {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	constraint = types.M{"$within": types.M{}}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = errs.E(errs.InvalidJSON, "malformatted $within arg")
	if reflect.DeepEqual(err, expect) == false || result != nil {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	constraint = types.M{
		"$within": types.M{"$box": "hello"},
	}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = errs.E(errs.InvalidJSON, "malformatted $within arg")
	if reflect.DeepEqual(err, expect) == false || result != nil {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	constraint = types.M{
		"$within": types.M{
			"$box": types.S{"hello"},
		},
	}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = errs.E(errs.InvalidJSON, "malformatted $within arg")
	if reflect.DeepEqual(err, expect) == false || result != nil {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	constraint = types.M{
		"$within": types.M{
			"$box": types.S{"hello", "world"},
		},
	}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = errs.E(errs.InvalidJSON, "malformatted $within arg")
	if reflect.DeepEqual(err, expect) == false || result != nil {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	constraint = types.M{
		"$within": types.M{
			"$box": types.S{
				types.M{
					"longitude": 20,
					"latitude":  20,
				},
				types.M{
					"longitude": 30,
					"latitude":  30,
				},
			},
		},
	}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{
		"$geoWithin": types.M{
			"$box": types.S{
				types.S{20, 20},
				types.S{30, 30},
			},
		},
	}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{"$other": "hello"}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = errs.E(errs.InvalidJSON, "bad constraint: "+"$other")
	if reflect.DeepEqual(err, expect) == false || result != nil {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	constraint = types.M{"key": "hello"}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = nil
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
}

func Test_transformTopLevelAtom(t *testing.T) {
	tf := NewTransform()
	var atom interface{}
	var result interface{}
	var err error
	var expect interface{}
	/*************************************************/
	atom = nil
	result, err = tf.transformTopLevelAtom(atom)
	expect = nil
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	atom = "hello"
	result, err = tf.transformTopLevelAtom(atom)
	expect = "hello"
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	atom = 20.0
	result, err = tf.transformTopLevelAtom(atom)
	expect = 20.0
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	atom = 20
	result, err = tf.transformTopLevelAtom(atom)
	expect = 20
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	atom = true
	result, err = tf.transformTopLevelAtom(atom)
	expect = true
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	atom = time.Now()
	result, err = tf.transformTopLevelAtom(atom)
	expect = atom
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	atom = types.M{}
	result, err = tf.transformTopLevelAtom(atom)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	atom = types.M{
		"__type":    "Pointer",
		"className": "user",
		"objectId":  "1024",
	}
	result, err = tf.transformTopLevelAtom(atom)
	expect = "user$1024"
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	tmpTime := utils.TimetoString(time.Now().UTC())
	atom = types.M{
		"__type": "Date",
		"iso":    tmpTime,
	}
	result, err = tf.transformTopLevelAtom(atom)
	expect, _ = utils.StringtoTime(tmpTime)
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	atom = types.M{
		"__type": "Bytes",
		"base64": "aGVsbG8=",
	}
	result, err = tf.transformTopLevelAtom(atom)
	expect = []byte("hello")
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	atom = types.M{
		"__type":    "GeoPoint",
		"longitude": -30.0,
		"latitude":  40.0,
	}
	result, err = tf.transformTopLevelAtom(atom)
	expect = types.S{-30.0, 40.0}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	atom = types.M{
		"__type": "File",
		"name":   "...hello.png",
	}
	result, err = tf.transformTopLevelAtom(atom)
	expect = "...hello.png"
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	atom = types.M{
		"__type": "OtherType",
		"key":    "value",
	}
	result, err = tf.transformTopLevelAtom(atom)
	expect = nil
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	atom = types.S{"hello", "world"}
	result, err = tf.transformTopLevelAtom(atom)
	expect = nil
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	atom = []string{}
	result, err = tf.transformTopLevelAtom(atom)
	expect = errs.E(errs.InternalServerError, "really did not expect value: atom")
	if reflect.DeepEqual(err, expect) == false || result != nil {
		t.Error("expect:", expect, "get result:", err)
	}
}

func Test_transformUpdateOperator(t *testing.T) {
	tf := NewTransform()
	var operator interface{}
	var flatten bool
	var result interface{}
	var err error
	var expect interface{}
	/*************************************************/
	operator = 1024
	flatten = true
	result, err = tf.transformUpdateOperator(operator, flatten)
	expect = 1024
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	operator = types.M{"key": "value"}
	flatten = true
	result, err = tf.transformUpdateOperator(operator, flatten)
	expect = types.M{"key": "value"}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	operator = types.M{"__op": "Delete"}
	flatten = true
	result, err = tf.transformUpdateOperator(operator, flatten)
	expect = nil
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	operator = types.M{"__op": "Delete"}
	flatten = false
	result, err = tf.transformUpdateOperator(operator, flatten)
	expect = types.M{"__op": "$unset", "arg": ""}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	operator = types.M{"__op": "Increment", "amount": 10.0}
	flatten = true
	result, err = tf.transformUpdateOperator(operator, flatten)
	expect = 10.0
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	operator = types.M{"__op": "Increment", "amount": 10}
	flatten = true
	result, err = tf.transformUpdateOperator(operator, flatten)
	expect = 10
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	operator = types.M{"__op": "Increment", "amount": "10"}
	flatten = true
	result, err = tf.transformUpdateOperator(operator, flatten)
	expect = errs.E(errs.InvalidJSON, "incrementing must provide a number")
	if err == nil || reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	operator = types.M{"__op": "Increment", "amount": 10}
	flatten = false
	result, err = tf.transformUpdateOperator(operator, flatten)
	expect = types.M{"__op": "$inc", "arg": 10}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	operator = types.M{
		"__op":    "Add",
		"objects": "not an array",
	}
	flatten = true
	result, err = tf.transformUpdateOperator(operator, flatten)
	expect = errs.E(errs.InvalidJSON, "objects to add must be an array")
	if err == nil || reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	operator = types.M{
		"__op":    "Add",
		"objects": types.S{"hello", "world"},
	}
	flatten = true
	result, err = tf.transformUpdateOperator(operator, flatten)
	expect = types.S{"hello", "world"}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	operator = types.M{
		"__op":    "Add",
		"objects": types.S{"hello", "world"},
	}
	flatten = false
	result, err = tf.transformUpdateOperator(operator, flatten)
	expect = types.M{
		"__op": "$push",
		"arg": types.M{
			"$each": types.S{"hello", "world"},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	operator = types.M{
		"__op":    "AddUnique",
		"objects": "not an array",
	}
	flatten = true
	result, err = tf.transformUpdateOperator(operator, flatten)
	expect = errs.E(errs.InvalidJSON, "objects to add must be an array")
	if err == nil || reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	operator = types.M{
		"__op":    "AddUnique",
		"objects": types.S{"hello", "world"},
	}
	flatten = true
	result, err = tf.transformUpdateOperator(operator, flatten)
	expect = types.S{"hello", "world"}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	operator = types.M{
		"__op":    "AddUnique",
		"objects": types.S{"hello", "world"},
	}
	flatten = false
	result, err = tf.transformUpdateOperator(operator, flatten)
	expect = types.M{
		"__op": "$addToSet",
		"arg": types.M{
			"$each": types.S{"hello", "world"},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	operator = types.M{
		"__op":    "Remove",
		"objects": "not an array",
	}
	flatten = true
	result, err = tf.transformUpdateOperator(operator, flatten)
	expect = errs.E(errs.InvalidJSON, "objects to remove must be an array")
	if err == nil || reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	operator = types.M{
		"__op":    "Remove",
		"objects": types.S{"hello", "world"},
	}
	flatten = true
	result, err = tf.transformUpdateOperator(operator, flatten)
	expect = types.S{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	operator = types.M{
		"__op":    "Remove",
		"objects": types.S{"hello", "world"},
	}
	flatten = false
	result, err = tf.transformUpdateOperator(operator, flatten)
	expect = types.M{
		"__op": "$pullAll",
		"arg":  types.S{"hello", "world"},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	operator = types.M{
		"__op": "OtherOp",
	}
	flatten = true
	result, err = tf.transformUpdateOperator(operator, flatten)
	expect = errs.E(errs.CommandUnavailable, "the "+"OtherOp"+" operator is not supported yet")
	if err == nil || reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "get result:", err)
	}
}

func Test_parseObjectToMongoObjectForCreate(t *testing.T) {
	// TODO
}

func Test_parseObjectKeyValueToMongoObjectKeyValue(t *testing.T) {
	tf := NewTransform()
	var restKey string
	var restValue interface{}
	var schema types.M
	var resultKey string
	var resultValue interface{}
	var err error
	var expectKey string
	var expectValue interface{}
	tmpTimeStr := utils.TimetoString(time.Now().UTC())
	tmpTime, _ := utils.StringtoTime(tmpTimeStr)
	/*************************************************/
	restKey = "objectId"
	restValue = "123456"
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue("", restKey, restValue, schema)
	expectKey = "_id"
	expectValue = "123456"
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "createdAt"
	restValue = tmpTimeStr
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue("", restKey, restValue, schema)
	expectKey = "_created_at"
	expectValue = tmpTime
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "createdAt"
	restValue = tmpTime
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue("", restKey, restValue, schema)
	expectKey = "_created_at"
	expectValue = tmpTime
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "updatedAt"
	restValue = tmpTimeStr
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue("", restKey, restValue, schema)
	expectKey = "_updated_at"
	expectValue = tmpTime
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "updatedAt"
	restValue = tmpTime
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue("", restKey, restValue, schema)
	expectKey = "_updated_at"
	expectValue = tmpTime
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "expiresAt"
	restValue = tmpTimeStr
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue("", restKey, restValue, schema)
	expectKey = "_expiresAt"
	expectValue = tmpTime
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "expiresAt"
	restValue = tmpTime
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue("", restKey, restValue, schema)
	expectKey = "_expiresAt"
	expectValue = tmpTime
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "_email_verify_token"
	restValue = "abcd"
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue("", restKey, restValue, schema)
	expectKey = "_email_verify_token"
	expectValue = "abcd"
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "sessionToken"
	restValue = "abcd"
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue("", restKey, restValue, schema)
	expectKey = "_session_token"
	expectValue = "abcd"
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "authData.facebook.id"
	restValue = "abcd"
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue("", restKey, restValue, schema)
	expectValue = errs.E(errs.InvalidKeyName, "can only query on "+restKey)
	if reflect.DeepEqual(expectValue, err) == false {
		t.Error("expect:", expectValue, "get result:", err)
	}
	/*************************************************/
	restKey = "_auth_data_facebook"
	restValue = "abcd"
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue("", restKey, restValue, schema)
	expectKey = "_auth_data_facebook"
	expectValue = "abcd"
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "key"
	restValue = nil
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue("", restKey, restValue, schema)
	expectKey = "key"
	expectValue = nil
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "key"
	restValue = types.M{
		"__type":    "Pointer",
		"className": "user",
		"objectId":  "1024",
	}
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue("", restKey, restValue, schema)
	expectKey = "_p_key"
	expectValue = "user$1024"
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "key"
	restValue = types.M{
		"__type":    "OtherPointer",
		"className": "user",
		"objectId":  "1024",
	}
	schema = types.M{
		"fields": types.M{
			"key": types.M{
				"type": "Pointer",
			},
		},
	}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue("", restKey, restValue, schema)
	expectKey = "_p_key"
	expectValue = types.M{
		"__type":    "OtherPointer",
		"className": "user",
		"objectId":  "1024",
	}
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "ACL"
	restValue = "abcd"
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue("", restKey, restValue, schema)
	expectValue = errs.E(errs.InvalidKeyName, "There was a problem transforming an ACL.")
	if reflect.DeepEqual(expectValue, err) == false {
		t.Error("expect:", expectValue, "get result:", err)
	}
	/*************************************************/
	restKey = "key"
	restValue = "value"
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue("", restKey, restValue, schema)
	expectKey = "key"
	expectValue = "value"
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "key"
	restValue = types.S{"hello", "world"}
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue("", restKey, restValue, schema)
	expectKey = "key"
	expectValue = types.S{"hello", "world"}
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "key"
	restValue = types.M{
		"__op":   "Increment",
		"amount": 10,
	}
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue("", restKey, restValue, schema)
	expectKey = "key"
	expectValue = types.M{
		"__op": "$inc",
		"arg":  10,
	}
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "key"
	restValue = types.M{"key$": "value"}
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue("", restKey, restValue, schema)
	expectValue = errs.E(errs.InvalidNestedKey, "Nested keys should not contain the '$' or '.' characters")
	if reflect.DeepEqual(expectValue, err) == false {
		t.Error("expect:", expectValue, "get result:", err)
	}
	/*************************************************/
	restKey = "key"
	restValue = types.M{"key.": "value"}
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue("", restKey, restValue, schema)
	expectValue = errs.E(errs.InvalidNestedKey, "Nested keys should not contain the '$' or '.' characters")
	if reflect.DeepEqual(expectValue, err) == false {
		t.Error("expect:", expectValue, "get result:", err)
	}
	/*************************************************/
	restKey = "key"
	restValue = types.M{
		"subKey": "subValue",
	}
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue("", restKey, restValue, schema)
	expectKey = "key"
	expectValue = types.M{
		"subKey": "subValue",
	}
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
}

func Test_transformAuthData(t *testing.T) {
	tf := NewTransform()
	var restObject types.M
	var result types.M
	var expect types.M
	/*************************************************/
	restObject = nil
	result = tf.transformAuthData(restObject)
	expect = restObject
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	restObject = types.M{}
	result = tf.transformAuthData(restObject)
	expect = restObject
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	restObject = types.M{"authData": 1024}
	result = tf.transformAuthData(restObject)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	restObject = types.M{
		"authData": types.M{
			"facebook": nil,
		},
	}
	result = tf.transformAuthData(restObject)
	expect = types.M{
		"_auth_data_facebook": types.M{
			"__op": "Delete",
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	restObject = types.M{
		"authData": types.M{
			"facebook": 1024,
		},
	}
	result = tf.transformAuthData(restObject)
	expect = types.M{
		"_auth_data_facebook": types.M{
			"__op": "Delete",
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	restObject = types.M{
		"authData": types.M{
			"facebook": types.M{},
		},
	}
	result = tf.transformAuthData(restObject)
	expect = types.M{
		"_auth_data_facebook": types.M{
			"__op": "Delete",
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	restObject = types.M{
		"authData": types.M{
			"facebook": types.M{"id": "1024"},
			"twitter":  types.M{},
		},
	}
	result = tf.transformAuthData(restObject)
	expect = types.M{
		"_auth_data_facebook": types.M{"id": "1024"},
		"_auth_data_twitter":  types.M{"__op": "Delete"},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
}

func Test_transformACL(t *testing.T) {
	tf := NewTransform()
	var restObject types.M
	var result types.M
	var expect types.M
	/*************************************************/
	restObject = nil
	result = tf.transformACL(restObject)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	restObject = types.M{}
	result = tf.transformACL(restObject)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	restObject = types.M{
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
	result = tf.transformACL(restObject)
	expect = types.M{
		"_rperm": types.S{"userid", "role:xxx", "*"},
		"_wperm": types.S{"userid", "role:xxx"},
		"_acl": types.M{
			"userid": types.M{
				"r": true,
				"w": true,
			},
			"role:xxx": types.M{
				"r": true,
				"w": true,
			},
			"*": types.M{
				"r": true,
			},
		},
	}
	if utils.CompareArray(expect["_rperm"], result["_rperm"]) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	if utils.CompareArray(expect["_wperm"], result["_wperm"]) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	if reflect.DeepEqual(expect["_acl"], result["_acl"]) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	if reflect.DeepEqual(restObject, types.M{}) == false {
		t.Error("expect:", types.M{}, "get result:", restObject)
	}
}

func Test_transformWhere(t *testing.T) {
	// transformQueryKeyValue
	// TODO
}

func Test_transformUpdate(t *testing.T) {
	// transformKeyValueForUpdate
	// TODO
}

func Test_nestedMongoObjectToNestedParseObject(t *testing.T) {
	// TODO
}

func Test_mongoObjectToParseObject(t *testing.T) {
	// nestedMongoObjectToNestedParseObject
	// TODO
}

func Test_untransformACL(t *testing.T) {
	tf := NewTransform()
	var mongoObject types.M
	var result types.M
	var expect types.M
	/*************************************************/
	mongoObject = nil
	result = tf.untransformACL(mongoObject)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.M{}
	result = tf.untransformACL(mongoObject)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.M{}
	result = tf.untransformACL(mongoObject)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.M{"_rperm": "Incorrect type"}
	result = tf.untransformACL(mongoObject)
	expect = types.M{"ACL": types.M{}}
	if reflect.DeepEqual(expect, result) == false || reflect.DeepEqual(mongoObject, types.M{}) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.M{"_rperm": types.S{"userid", "role:xxx"}}
	result = tf.untransformACL(mongoObject)
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
	if reflect.DeepEqual(expect, result) == false || reflect.DeepEqual(mongoObject, types.M{}) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.M{"_wperm": "Incorrect type"}
	result = tf.untransformACL(mongoObject)
	expect = types.M{"ACL": types.M{}}
	if reflect.DeepEqual(expect, result) == false || reflect.DeepEqual(mongoObject, types.M{}) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.M{"_wperm": types.S{"userid", "role:xxx"}}
	result = tf.untransformACL(mongoObject)
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
	if reflect.DeepEqual(expect, result) == false || reflect.DeepEqual(mongoObject, types.M{}) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.M{
		"_rperm": types.S{"userid", "role:xxx", "*"},
		"_wperm": types.S{"userid", "role:xxx"},
	}
	result = tf.untransformACL(mongoObject)
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
	if reflect.DeepEqual(expect, result) == false || reflect.DeepEqual(mongoObject, types.M{}) == false {
		t.Error("expect:", expect, "get result:", result)
	}
}

func Test_transformInteriorAtom(t *testing.T) {
	tf := NewTransform()
	var atom interface{}
	var result interface{}
	var expect interface{}
	var err error
	/*************************************************/
	atom = nil
	result, err = tf.transformInteriorAtom(atom)
	if err != nil || result != nil {
		t.Error("expect:", nil, "get result:", result)
	}
	/*************************************************/
	atom = types.M{
		"__type":    "Pointer",
		"className": "user",
		"objectId":  "1024",
	}
	result, err = tf.transformInteriorAtom(atom)
	expect = types.M{
		"__type":    "Pointer",
		"className": "user",
		"objectId":  "1024",
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	tmpTime := utils.TimetoString(time.Now().UTC())
	atom = types.M{
		"__type": "Date",
		"iso":    tmpTime,
	}
	result, err = tf.transformInteriorAtom(atom)
	expect, _ = utils.StringtoTime(tmpTime)
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	atom = types.M{
		"__type": "Bytes",
		"base64": "aGVsbG8=",
	}
	result, err = tf.transformInteriorAtom(atom)
	expect = []byte("hello")
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	atom = 1024
	result, err = tf.transformInteriorAtom(atom)
	expect = 1024
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	atom = types.M{"key": "value"}
	result, err = tf.transformInteriorAtom(atom)
	expect = nil
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	atom = types.S{"hello", "world"}
	result, err = tf.transformInteriorAtom(atom)
	expect = nil
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
}

func Test_transformInteriorValue(t *testing.T) {
	tf := NewTransform()
	var restValue interface{}
	var result interface{}
	var err error
	var expect interface{}
	/*************************************************/
	restValue = nil
	result, err = tf.transformInteriorValue(restValue)
	expect = nil
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	restValue = types.M{"hello$world": "1024"}
	result, err = tf.transformInteriorValue(restValue)
	expect = errs.E(errs.InvalidNestedKey, "Nested keys should not contain the '$' or '.' characters")
	if reflect.DeepEqual(expect, err) == false || result != nil {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	restValue = types.M{"hello.world": "1024"}
	result, err = tf.transformInteriorValue(restValue)
	expect = errs.E(errs.InvalidNestedKey, "Nested keys should not contain the '$' or '.' characters")
	if reflect.DeepEqual(expect, err) == false || result != nil {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	restValue = "hello world"
	result, err = tf.transformInteriorValue(restValue)
	expect = "hello world"
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	restValue = types.M{
		"__type": "Bytes",
		"base64": "aGVsbG8=",
	}
	result, err = tf.transformInteriorValue(restValue)
	expect = []byte("hello")
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	restValue = types.S{
		"hello",
		types.M{
			"__type": "Bytes",
			"base64": "aGVsbG8=",
		},
	}
	result, err = tf.transformInteriorValue(restValue)
	expect = types.S{"hello", []byte("hello")}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	restValue = types.M{
		"__op":    "Add",
		"objects": types.S{"hello", "world"},
	}
	result, err = tf.transformInteriorValue(restValue)
	expect = types.S{"hello", "world"}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	restValue = types.M{
		"key1": "hello world",
		"key2": types.M{
			"__type": "Bytes",
			"base64": "aGVsbG8=",
		},
		"key3": types.S{
			"hello",
			types.M{
				"__type": "Bytes",
				"base64": "aGVsbG8=",
			},
		},
		"key4": types.M{
			"__op":    "Add",
			"objects": types.S{"hello", "world"},
		},
	}
	result, err = tf.transformInteriorValue(restValue)
	expect = types.M{
		"key1": "hello world",
		"key2": []byte("hello"),
		"key3": types.S{"hello", []byte("hello")},
		"key4": types.S{"hello", "world"},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
}

func Test_dateCoder(t *testing.T) {
	dc := dateCoder{}
	var databaseObject interface{}
	var jsonObject types.M
	var ok bool
	var expect interface{}
	var err error
	/*************************************************/
	databaseObject = "pic.jpg"
	jsonObject = dc.databaseToJSON(databaseObject)
	expect = types.M{
		"__type": "Date",
		"iso":    "",
	}
	if reflect.DeepEqual(jsonObject, expect) == false {
		t.Error("expect:", expect, "get jsonObject:", jsonObject)
	}
	/*************************************************/
	databaseObject = time.Now().UTC()
	jsonObject = dc.databaseToJSON(databaseObject)
	expect = types.M{
		"__type": "Date",
		"iso":    utils.TimetoString(databaseObject.(time.Time)),
	}
	if reflect.DeepEqual(jsonObject, expect) == false {
		t.Error("expect:", expect, "get jsonObject:", jsonObject)
	}
	/*************************************************/
	databaseObject = "pic.jpg"
	ok = dc.isValidDatabaseObject(databaseObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	databaseObject = time.Now().UTC()
	ok = dc.isValidDatabaseObject(databaseObject)
	if !ok {
		t.Error("expect:", "true", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{
		"__type": "Date",
		"iso":    "aabdcc",
	}
	databaseObject, err = dc.jsonToDatabase(jsonObject)
	expect = errs.E(errs.InvalidJSON, "invalid iso")
	if reflect.DeepEqual(err, expect) == false {
		t.Error("expect:", expect, "get:", err)
	}
	/*************************************************/
	tmpTimeStr := utils.TimetoString(time.Now().UTC())
	jsonObject = types.M{
		"__type": "Date",
		"iso":    tmpTimeStr,
	}
	databaseObject, err = dc.jsonToDatabase(jsonObject)
	expect, _ = utils.StringtoTime(tmpTimeStr)
	if err != nil || reflect.DeepEqual(expect, databaseObject) == false {
		t.Error("expect:", expect, "get:", databaseObject)
	}
	/*************************************************/
	jsonObject = nil
	ok = dc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{}
	ok = dc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{"__type": "Bytes"}
	ok = dc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{"__type": "Date"}
	ok = dc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{"__type": "Date", "iso": 1024}
	ok = dc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{"__type": "Date", "iso": "2006-01-02T15:04:05.000Z"}
	ok = dc.isValidJSON(jsonObject)
	if !ok {
		t.Error("expect:", "true", "get:", ok)
	}
}

func Test_bytesCoder(t *testing.T) {
	bc := bytesCoder{}
	var databaseObject interface{}
	var jsonObject types.M
	var ok bool
	var expect interface{}
	var err error
	/*************************************************/
	databaseObject = "pic.jpg"
	jsonObject = bc.databaseToJSON(databaseObject)
	expect = types.M{
		"__type": "Bytes",
		"base64": "",
	}
	if reflect.DeepEqual(jsonObject, expect) == false {
		t.Error("expect:", expect, "get jsonObject:", jsonObject)
	}
	/*************************************************/
	databaseObject = []byte("hello")
	jsonObject = bc.databaseToJSON(databaseObject)
	expect = types.M{
		"__type": "Bytes",
		"base64": "aGVsbG8=",
	}
	if reflect.DeepEqual(jsonObject, expect) == false {
		t.Error("expect:", expect, "get jsonObject:", jsonObject)
	}
	/*************************************************/
	databaseObject = "pic.jpg"
	ok = bc.isValidDatabaseObject(databaseObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	databaseObject = []byte("hello")
	ok = bc.isValidDatabaseObject(databaseObject)
	if !ok {
		t.Error("expect:", "true", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{
		"__type": "Bytes",
		"base64": "aabbcc",
	}
	databaseObject, err = bc.jsonToDatabase(jsonObject)
	expect = errs.E(errs.InvalidJSON, "invalid base64")
	if reflect.DeepEqual(err, expect) == false {
		t.Error("expect:", expect, "get:", err)
	}
	/*************************************************/
	jsonObject = types.M{
		"__type": "Bytes",
		"base64": "aGVsbG8=",
	}
	databaseObject, err = bc.jsonToDatabase(jsonObject)
	expect = []byte("hello")
	if err != nil || reflect.DeepEqual(expect, databaseObject) == false {
		t.Error("expect:", expect, "get:", databaseObject)
	}
	/*************************************************/
	jsonObject = nil
	ok = bc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{}
	ok = bc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{"__type": "Date"}
	ok = bc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{"__type": "Bytes"}
	ok = bc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{"__type": "Bytes", "base64": 1024}
	ok = bc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{"__type": "Bytes", "base64": "aGVsbG8="}
	ok = bc.isValidJSON(jsonObject)
	if !ok {
		t.Error("expect:", "true", "get:", ok)
	}
}

func Test_geoPointCoder(t *testing.T) {
	gpc := geoPointCoder{}
	var databaseObject interface{}
	var jsonObject types.M
	var ok bool
	var expect interface{}
	var err error
	/*************************************************/
	databaseObject = "Incorrect type"
	jsonObject = gpc.databaseToJSON(databaseObject)
	expect = types.M{
		"__type":    "GeoPoint",
		"longitude": 0,
		"latitude":  0,
	}
	if reflect.DeepEqual(jsonObject, expect) == false {
		t.Error("expect:", expect, "get jsonObject:", jsonObject)
	}
	/*************************************************/
	databaseObject = types.S{20, 20, 20}
	jsonObject = gpc.databaseToJSON(databaseObject)
	expect = types.M{
		"__type":    "GeoPoint",
		"longitude": 0,
		"latitude":  0,
	}
	if reflect.DeepEqual(jsonObject, expect) == false {
		t.Error("expect:", expect, "get jsonObject:", jsonObject)
	}
	/*************************************************/
	databaseObject = types.S{20, 20}
	jsonObject = gpc.databaseToJSON(databaseObject)
	expect = types.M{
		"__type":    "GeoPoint",
		"longitude": 20,
		"latitude":  20,
	}
	if reflect.DeepEqual(jsonObject, expect) == false {
		t.Error("expect:", expect, "get jsonObject:", jsonObject)
	}
	/*************************************************/
	databaseObject = "Incorrect type"
	ok = gpc.isValidDatabaseObject(databaseObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	databaseObject = types.S{20, 20, 20}
	ok = gpc.isValidDatabaseObject(databaseObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	databaseObject = types.S{"20", "20"}
	ok = gpc.isValidDatabaseObject(databaseObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	databaseObject = types.S{20, 20}
	ok = gpc.isValidDatabaseObject(databaseObject)
	if !ok {
		t.Error("expect:", "true", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{
		"__type":    "GeoPoint",
		"longitude": "20.0",
		"latitude":  20.0,
	}
	databaseObject, err = gpc.jsonToDatabase(jsonObject)
	expect = errs.E(errs.InvalidJSON, "invalid longitude")
	if reflect.DeepEqual(err, expect) == false {
		t.Error("expect:", expect, "get:", err)
	}
	/*************************************************/
	jsonObject = types.M{
		"__type":    "GeoPoint",
		"longitude": 20.0,
		"latitude":  "20.0",
	}
	databaseObject, err = gpc.jsonToDatabase(jsonObject)
	expect = errs.E(errs.InvalidJSON, "invalid latitude")
	if reflect.DeepEqual(err, expect) == false {
		t.Error("expect:", expect, "get:", err)
	}
	/*************************************************/
	jsonObject = types.M{
		"__type":    "GeoPoint",
		"longitude": 20.0,
		"latitude":  20.0,
	}
	databaseObject, err = gpc.jsonToDatabase(jsonObject)
	expect = types.S{20.0, 20.0}
	if err != nil || reflect.DeepEqual(databaseObject, expect) == false {
		t.Error("expect:", expect, "get:", databaseObject)
	}
	/*************************************************/
	jsonObject = nil
	ok = gpc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{}
	ok = gpc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{"__type": "Date"}
	ok = gpc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{"__type": "GeoPoint"}
	ok = gpc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{"__type": "GeoPoint", "longitude": 20}
	ok = gpc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{"__type": "GeoPoint", "longitude": 20, "latitude": 20}
	ok = gpc.isValidJSON(jsonObject)
	if !ok {
		t.Error("expect:", "true", "get:", ok)
	}
}

func Test_fileCoder(t *testing.T) {
	fc := fileCoder{}
	var databaseObject interface{}
	var jsonObject types.M
	var ok bool
	var expect interface{}
	/*************************************************/
	databaseObject = "pic.jpg"
	jsonObject = fc.databaseToJSON(databaseObject)
	expect = types.M{
		"__type": "File",
		"name":   "pic.jpg",
	}
	if reflect.DeepEqual(jsonObject, expect) == false {
		t.Error("expect:", expect, "get jsonObject:", jsonObject)
	}
	/*************************************************/
	databaseObject = 1024
	ok = fc.isValidDatabaseObject(databaseObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	databaseObject = "pic.jpg"
	ok = fc.isValidDatabaseObject(databaseObject)
	if !ok {
		t.Error("expect:", "true", "get:", ok)
	}
	/*************************************************/
	jsonObject = nil
	databaseObject, _ = fc.jsonToDatabase(jsonObject)
	if databaseObject != nil {
		t.Error("expect:", "nil", "get:", databaseObject)
	}
	/*************************************************/
	jsonObject = types.M{
		"__type": "File",
		"name":   "pic.jpg",
	}
	databaseObject, _ = fc.jsonToDatabase(jsonObject)
	if reflect.DeepEqual("pic.jpg", databaseObject) == false {
		t.Error("expect:", "pic.jpg", "get:", databaseObject)
	}
	/*************************************************/
	jsonObject = nil
	ok = fc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{}
	ok = fc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{"__type": "Date"}
	ok = fc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{"__type": "File"}
	ok = fc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{"__type": "File", "name": 1024}
	ok = fc.isValidJSON(jsonObject)
	if ok {
		t.Error("expect:", "false", "get:", ok)
	}
	/*************************************************/
	jsonObject = types.M{"__type": "File", "name": "pic.jpg"}
	ok = fc.isValidJSON(jsonObject)
	if !ok {
		t.Error("expect:", "true", "get:", ok)
	}
}

func Test_valueAsDate(t *testing.T) {
	var value interface{}
	var date time.Time
	var ok bool
	/*************************************************/
	value = 1024
	date, ok = valueAsDate(value)
	if ok {
		t.Error("value:", value, "date:", date, "expect: false", "get:", ok)
	}
	/*************************************************/
	value = "Incorrect string time"
	date, ok = valueAsDate(value)
	if ok {
		t.Error("value:", value, "date:", date, "expect: false", "get:", ok)
	}
	/*************************************************/
	value = "2006-01-02T15:04:05.000Z"
	date, ok = valueAsDate(value)
	if !ok || utils.TimetoString(date) != "2006-01-02T15:04:05.000Z" {
		t.Error("value:", value, "date:", date, "expect: true 2006-01-02T15:04:05.000Z", "get:", ok, utils.TimetoString(date))
	}
	/*************************************************/
	value = time.Now().UTC()
	date, ok = valueAsDate(value)
	if !ok {
		t.Error("value:", value, "date:", date, "expect: true", "get:", ok)
	}
}
