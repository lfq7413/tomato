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
	// transformInteriorValue
	// TODO
}

func Test_transformQueryKeyValue(t *testing.T) {
	// transformWhere
	// transformConstraint
	// TODO
}

func Test_transformConstraint(t *testing.T) {
	// TODO
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
	atom = time.Now()
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
	// parseObjectKeyValueToMongoObjectKeyValue
	// TODO
}

func Test_parseObjectKeyValueToMongoObjectKeyValue(t *testing.T) {
	// transformInteriorValue
	// TODO
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

	// TODO
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
