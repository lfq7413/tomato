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
	tf := NewTransform()
	var restKey string
	var restValue interface{}
	var parseFormatSchema types.M
	var resultKey string
	var resultValue interface{}
	var err error
	var expectKey string
	var expectValue interface{}
	tmpTimeStr := utils.TimetoString(time.Now().UTC())
	tmpTime, _ := utils.StringtoTime(tmpTimeStr)
	/*************************************************/
	restKey = "objectId"
	restValue = "1024"
	parseFormatSchema = types.M{}
	resultKey, resultValue, err = tf.transformKeyValueForUpdate("", restKey, restValue, parseFormatSchema)
	expectKey = "_id"
	expectValue = "1024"
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectValue, "get result:", resultValue)
	}
	/*************************************************/
	restKey = "createdAt"
	restValue = tmpTimeStr
	parseFormatSchema = types.M{}
	resultKey, resultValue, err = tf.transformKeyValueForUpdate("", restKey, restValue, parseFormatSchema)
	expectKey = "_created_at"
	expectValue = tmpTime
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectValue, "get result:", resultValue)
	}
	/*************************************************/
	restKey = "updatedAt"
	restValue = tmpTimeStr
	parseFormatSchema = types.M{}
	resultKey, resultValue, err = tf.transformKeyValueForUpdate("", restKey, restValue, parseFormatSchema)
	expectKey = "_updated_at"
	expectValue = tmpTime
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectValue, "get result:", resultValue)
	}
	/*************************************************/
	restKey = "sessionToken"
	restValue = "abc"
	parseFormatSchema = types.M{}
	resultKey, resultValue, err = tf.transformKeyValueForUpdate("", restKey, restValue, parseFormatSchema)
	expectKey = "_session_token"
	expectValue = "abc"
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectValue, "get result:", resultValue)
	}
	/*************************************************/
	restKey = "expiresAt"
	restValue = tmpTimeStr
	parseFormatSchema = types.M{}
	resultKey, resultValue, err = tf.transformKeyValueForUpdate("", restKey, restValue, parseFormatSchema)
	expectKey = "_expiresAt"
	expectValue = tmpTime
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectValue, "get result:", resultValue)
	}
	/*************************************************/
	restKey = "_email_verify_token_expires_at"
	restValue = tmpTimeStr
	parseFormatSchema = types.M{}
	resultKey, resultValue, err = tf.transformKeyValueForUpdate("", restKey, restValue, parseFormatSchema)
	expectKey = "_email_verify_token_expires_at"
	expectValue = tmpTime
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectValue, "get result:", resultValue)
	}
	/*************************************************/
	restKey = "_rperm"
	restValue = "r"
	parseFormatSchema = types.M{}
	resultKey, resultValue, err = tf.transformKeyValueForUpdate("", restKey, restValue, parseFormatSchema)
	expectKey = "_rperm"
	expectValue = "r"
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectValue, "get result:", resultValue)
	}
	/*************************************************/
	restKey = "key"
	restValue = "value"
	parseFormatSchema = types.M{
		"fields": types.M{
			"key": types.M{
				"type": "Pointer",
			},
		},
	}
	resultKey, resultValue, err = tf.transformKeyValueForUpdate("", restKey, restValue, parseFormatSchema)
	expectKey = "_p_key"
	expectValue = "value"
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectValue, "get result:", resultValue)
	}
	/*************************************************/
	restKey = "key"
	restValue = types.M{
		"__type":    "Pointer",
		"className": "post",
		"objectId":  "1024",
	}
	parseFormatSchema = types.M{}
	resultKey, resultValue, err = tf.transformKeyValueForUpdate("", restKey, restValue, parseFormatSchema)
	expectKey = "_p_key"
	expectValue = "post$1024"
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectValue, "get result:", resultValue)
	}
	/*************************************************/
	restKey = "key"
	restValue = nil
	parseFormatSchema = types.M{}
	resultKey, resultValue, err = tf.transformKeyValueForUpdate("", restKey, restValue, parseFormatSchema)
	expectKey = "key"
	expectValue = nil
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectValue, "get result:", resultValue)
	}
	/*************************************************/
	restKey = "key"
	restValue = "value"
	parseFormatSchema = types.M{}
	resultKey, resultValue, err = tf.transformKeyValueForUpdate("", restKey, restValue, parseFormatSchema)
	expectKey = "key"
	expectValue = "value"
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectValue, "get result:", resultValue)
	}
	/*************************************************/
	restKey = "key"
	restValue = types.S{"value"}
	parseFormatSchema = types.M{}
	resultKey, resultValue, err = tf.transformKeyValueForUpdate("", restKey, restValue, parseFormatSchema)
	expectKey = "key"
	expectValue = types.S{"value"}
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectValue, "get result:", resultValue)
	}
	/*************************************************/
	restKey = "key"
	restValue = types.M{
		"__op":   "Increment",
		"amount": 10,
	}
	parseFormatSchema = types.M{}
	resultKey, resultValue, err = tf.transformKeyValueForUpdate("", restKey, restValue, parseFormatSchema)
	expectKey = "key"
	expectValue = types.M{
		"__op": "$inc",
		"arg":  10,
	}
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectValue, "get result:", resultValue)
	}
	/*************************************************/
	restKey = "key"
	restValue = types.M{
		"key": "value",
	}
	parseFormatSchema = types.M{}
	resultKey, resultValue, err = tf.transformKeyValueForUpdate("", restKey, restValue, parseFormatSchema)
	expectKey = "key"
	expectValue = types.M{
		"key": "value",
	}
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectValue, "get result:", resultValue)
	}
	/*************************************************/
	restKey = "key"
	restValue = []string{"hello"}
	parseFormatSchema = types.M{}
	resultKey, resultValue, err = tf.transformKeyValueForUpdate("", restKey, restValue, parseFormatSchema)
	expectKey = "key"
	expectValue = []string{"hello"}
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
}

func Test_transformQueryKeyValue(t *testing.T) {
	tf := NewTransform()
	var key string
	var value interface{}
	var schema types.M
	var resultKey string
	var resultValue interface{}
	var err error
	var expectKey string
	var expectValue interface{}
	tmpTimeStr := utils.TimetoString(time.Now().UTC())
	tmpTime, _ := utils.StringtoTime(tmpTimeStr)
	/*************************************************/
	key = "createdAt"
	value = tmpTimeStr
	schema = types.M{}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectKey = "_created_at"
	expectValue = tmpTime
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
	/*************************************************/
	key = "updatedAt"
	value = tmpTimeStr
	schema = types.M{}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectKey = "_updated_at"
	expectValue = tmpTime
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
	/*************************************************/
	key = "expiresAt"
	value = tmpTimeStr
	schema = types.M{}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectKey = "_expiresAt"
	expectValue = tmpTime
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
	/*************************************************/
	key = "_email_verify_token_expires_at"
	value = tmpTimeStr
	schema = types.M{}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectKey = "_email_verify_token_expires_at"
	expectValue = tmpTime
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
	/*************************************************/
	key = "objectId"
	value = "1024"
	schema = types.M{}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectKey = "_id"
	expectValue = "1024"
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
	/*************************************************/
	key = "sessionToken"
	value = "abc"
	schema = types.M{}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectKey = "_session_token"
	expectValue = "abc"
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
	/*************************************************/
	key = "_rperm"
	value = "abc"
	schema = types.M{}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectKey = "_rperm"
	expectValue = "abc"
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
	/*************************************************/
	key = "$or"
	value = nil
	schema = types.M{}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectKey = "$or"
	expectValue = nil
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
	/*************************************************/
	key = "$or"
	value = "hello"
	schema = types.M{}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectKey = "$or"
	expectValue = nil
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
	/*************************************************/
	key = "$or"
	value = types.S{types.M{"name": "joe"}}
	schema = types.M{}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectKey = "$or"
	expectValue = types.S{types.M{"name": "joe"}}
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
	/*************************************************/
	key = "$or"
	value = types.S{types.M{"name": "joe"}, "hello"}
	schema = types.M{}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectKey = "$or"
	expectValue = types.S{types.M{"name": "joe"}}
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
	/*************************************************/
	key = "$or"
	value = types.S{types.M{"name": "joe"}, types.M{"age": 25}}
	schema = types.M{}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectKey = "$or"
	expectValue = types.S{types.M{"name": "joe"}, types.M{"age": 25}}
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
	/*************************************************/
	key = "$and"
	value = nil
	schema = types.M{}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectKey = "$and"
	expectValue = nil
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
	/*************************************************/
	key = "$and"
	value = "hello"
	schema = types.M{}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectKey = "$and"
	expectValue = nil
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
	/*************************************************/
	key = "$and"
	value = types.S{types.M{"name": "joe"}}
	schema = types.M{}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectKey = "$and"
	expectValue = types.S{types.M{"name": "joe"}}
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
	/*************************************************/
	key = "$and"
	value = types.S{types.M{"name": "joe"}, "hello"}
	schema = types.M{}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectKey = "$and"
	expectValue = types.S{types.M{"name": "joe"}}
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
	/*************************************************/
	key = "$and"
	value = types.S{types.M{"name": "joe"}, types.M{"age": 25}}
	schema = types.M{}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectKey = "$and"
	expectValue = types.S{types.M{"name": "joe"}, types.M{"age": 25}}
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
	/*************************************************/
	key = "authData.facebook.id"
	value = "1024"
	schema = types.M{}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectKey = "_auth_data_facebook.id"
	expectValue = "1024"
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
	/*************************************************/
	key = "user"
	value = "1024"
	schema = types.M{
		"fields": types.M{
			"user": types.M{
				"type": "Pointer",
			},
		},
	}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectKey = "_p_user"
	expectValue = "1024"
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
	/*************************************************/
	key = "user"
	value = types.M{
		"__type":    "Pointer",
		"className": "_User",
		"objectId":  "1024",
	}
	schema = types.M{}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectKey = "_p_user"
	expectValue = "_User$1024"
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
	/*************************************************/
	key = "age"
	value = types.M{
		"$lt": 25,
		"$gt": 20,
	}
	schema = types.M{}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectKey = "age"
	expectValue = types.M{
		"$lt": 25,
		"$gt": 20,
	}
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
	/*************************************************/
	key = "skill"
	value = "one"
	schema = types.M{
		"fields": types.M{
			"skill": types.M{
				"type": "Array",
			},
		},
	}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectKey = "skill"
	expectValue = types.M{
		"$all": types.S{"one"},
	}
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
	/*************************************************/
	key = "key"
	value = "value"
	schema = types.M{}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectKey = "key"
	expectValue = "value"
	if err != nil || resultKey != expectKey || reflect.DeepEqual(resultValue, expectValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue, err)
	}
	/*************************************************/
	key = "key"
	value = types.M{"key": "value"}
	schema = types.M{}
	resultKey, resultValue, err = tf.transformQueryKeyValue("", key, value, schema)
	expectValue = errs.E(errs.InvalidJSON, "You cannot use this value as a query parameter.")
	if reflect.DeepEqual(err, expectValue) == false {
		t.Error("expect:", expectValue, "get result:", err)
	}
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
	expect = types.M{"$options": "imxs"}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
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
	expect = types.M{
		"$nearSphere": types.M{
			"$geometry": types.M{
				"type":        "Point",
				"coordinates": types.S{0, 0},
			},
			"$maxDistance": 0,
		},
	}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{
		"$nearSphere": types.M{
			"longitude": 30,
			"latitude":  20,
		},
	}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{
		"$nearSphere": types.M{
			"$geometry": types.M{
				"type":        "Point",
				"coordinates": types.S{30, 20},
			},
		},
	}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{"$maxDistance": 0.26}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{
		"$nearSphere": types.M{
			"longitude": 30,
			"latitude":  20,
		},
		"$maxDistance": 0.26,
	}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{
		"$nearSphere": types.M{
			"$geometry": types.M{
				"type":        "Point",
				"coordinates": types.S{30, 20},
			},
			"$maxDistance": 0.26 * 6371 * 1000,
		},
	}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{
		"$nearSphere": types.M{
			"longitude": 30,
			"latitude":  20,
		},
		"$maxDistanceInRadians": 0.26,
	}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{
		"$nearSphere": types.M{
			"$geometry": types.M{
				"type":        "Point",
				"coordinates": types.S{30, 20},
			},
			"$maxDistance": 0.26 * 6371 * 1000,
		},
	}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{
		"$nearSphere": types.M{
			"longitude": 30,
			"latitude":  20,
		},
		"$maxDistanceInMiles": 16.0,
	}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{
		"$nearSphere": types.M{
			"$geometry": types.M{
				"type":        "Point",
				"coordinates": types.S{30, 20},
			},
			"$maxDistance": 16.0 * 1.609344 * 1000,
		},
	}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{
		"$nearSphere": types.M{
			"longitude": 30,
			"latitude":  20,
		},
		"$maxDistanceInMiles": 16,
	}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{
		"$nearSphere": types.M{
			"$geometry": types.M{
				"type":        "Point",
				"coordinates": types.S{30, 20},
			},
			"$maxDistance": 16 * 1.609344 * 1000,
		},
	}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{
		"$nearSphere": types.M{
			"longitude": 30,
			"latitude":  20,
		},
		"$maxDistanceInKilometers": 16.0,
	}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{
		"$nearSphere": types.M{
			"$geometry": types.M{
				"type":        "Point",
				"coordinates": types.S{30, 20},
			},
			"$maxDistance": 16.0 * 1000,
		},
	}
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	constraint = types.M{
		"$nearSphere": types.M{
			"longitude": 30,
			"latitude":  20,
		},
		"$maxDistanceInKilometers": 16,
	}
	inArray = true
	result, err = tf.transformConstraint(constraint, inArray)
	expect = types.M{
		"$nearSphere": types.M{
			"$geometry": types.M{
				"type":        "Point",
				"coordinates": types.S{30, 20},
			},
			"$maxDistance": 16.0 * 1000,
		},
	}
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
	expect = nil
	if err != nil || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
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
	tf := NewTransform()
	var className string
	var create types.M
	var schema types.M
	var result types.M
	var err error
	var expect types.M
	/*************************************************/
	className = "_User"
	create = nil
	schema = types.M{}
	result, err = tf.parseObjectToMongoObjectForCreate(className, create, schema)
	expect = nil
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	className = "_User"
	create = types.M{
		"_auth_data_facebook": types.M{
			"id": "1024",
		},
		"name": "joe",
	}
	schema = types.M{}
	result, err = tf.parseObjectToMongoObjectForCreate(className, create, schema)
	expect = types.M{
		"_auth_data_facebook": types.M{
			"id": "1024",
		},
		"name": "joe",
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	className = "post"
	create = types.M{
		"name":   "joe",
		"_rperm": types.S{"userid"},
		"_wperm": types.S{"userid"},
	}
	schema = types.M{}
	result, err = tf.parseObjectToMongoObjectForCreate(className, create, schema)
	expect = types.M{
		"_rperm": types.S{"userid"},
		"_wperm": types.S{"userid"},
		"_acl": types.M{
			"userid": types.M{
				"r": true,
				"w": true,
			},
		},
		"name": "joe",
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	tmpTimeStr := utils.TimetoString(time.Now().UTC())
	tmpTime, _ := utils.StringtoTime(tmpTimeStr)
	className = "_User"
	create = types.M{
		"objectId": "1024",
		"name":     "joe",
		"_rperm":   types.S{"userid"},
		"_wperm":   types.S{"userid"},
		"_auth_data_facebook": types.M{
			"id": "1024",
		},
		"createdAt":        tmpTimeStr,
		"updatedAt":        tmpTimeStr,
		"_hashed_password": "password",
	}
	schema = types.M{}
	result, err = tf.parseObjectToMongoObjectForCreate(className, create, schema)
	expect = types.M{
		"_id":    "1024",
		"name":   "joe",
		"_rperm": types.S{"userid"},
		"_wperm": types.S{"userid"},
		"_acl": types.M{
			"userid": types.M{
				"r": true,
				"w": true,
			},
		},
		"_auth_data_facebook": types.M{
			"id": "1024",
		},
		"_created_at":      tmpTime,
		"_updated_at":      tmpTime,
		"_hashed_password": "password",
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
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
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
	expectKey = "_id"
	expectValue = "123456"
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "createdAt"
	restValue = tmpTimeStr
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
	expectKey = "createdAt"
	expectValue = tmpTimeStr
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "createdAt"
	restValue = tmpTime
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
	expectKey = "createdAt"
	expectValue = tmpTime
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "createdAt"
	restValue = types.M{
		"__type": "Date",
		"iso":    tmpTimeStr,
	}
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
	expectKey = "createdAt"
	expectValue = tmpTime
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "updatedAt"
	restValue = tmpTimeStr
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
	expectKey = "updatedAt"
	expectValue = tmpTimeStr
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "updatedAt"
	restValue = tmpTime
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
	expectKey = "updatedAt"
	expectValue = tmpTime
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "updatedAt"
	restValue = types.M{
		"__type": "Date",
		"iso":    tmpTimeStr,
	}
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
	expectKey = "updatedAt"
	expectValue = tmpTime
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "expiresAt"
	restValue = tmpTimeStr
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
	expectKey = "_expiresAt"
	expectValue = tmpTime
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "expiresAt"
	restValue = tmpTime
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
	expectKey = "_expiresAt"
	expectValue = tmpTime
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "_email_verify_token_expires_at"
	restValue = tmpTimeStr
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
	expectKey = "_email_verify_token_expires_at"
	expectValue = tmpTime
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "_email_verify_token_expires_at"
	restValue = tmpTime
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
	expectKey = "_email_verify_token_expires_at"
	expectValue = tmpTime
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "_email_verify_token"
	restValue = "abcd"
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
	expectKey = "_email_verify_token"
	expectValue = "abcd"
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "sessionToken"
	restValue = "abcd"
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
	expectKey = "_session_token"
	expectValue = "abcd"
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "authData.facebook.id"
	restValue = "abcd"
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
	expectValue = errs.E(errs.InvalidKeyName, "can only query on "+restKey)
	if reflect.DeepEqual(expectValue, err) == false {
		t.Error("expect:", expectValue, "get result:", err)
	}
	/*************************************************/
	restKey = "_auth_data_facebook"
	restValue = "abcd"
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
	expectKey = "_auth_data_facebook"
	expectValue = "abcd"
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "key"
	restValue = nil
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
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
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
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
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
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
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
	expectValue = errs.E(errs.InvalidKeyName, "There was a problem transforming an ACL.")
	if reflect.DeepEqual(expectValue, err) == false {
		t.Error("expect:", expectValue, "get result:", err)
	}
	/*************************************************/
	restKey = "key"
	restValue = "value"
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
	expectKey = "key"
	expectValue = "value"
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "key"
	restValue = types.S{"hello", "world"}
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
	expectKey = "key"
	expectValue = types.S{"hello", "world"}
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
	/*************************************************/
	restKey = "key"
	restValue = types.M{"key$": "value"}
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
	expectValue = errs.E(errs.InvalidNestedKey, "Nested keys should not contain the '$' or '.' characters")
	if reflect.DeepEqual(expectValue, err) == false {
		t.Error("expect:", expectValue, "get result:", err)
	}
	/*************************************************/
	restKey = "key"
	restValue = types.M{"key.": "value"}
	schema = types.M{}
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
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
	resultKey, resultValue, err = tf.parseObjectKeyValueToMongoObjectKeyValue(restKey, restValue, schema)
	expectKey = "key"
	expectValue = types.M{
		"subKey": "subValue",
	}
	if err != nil || expectKey != resultKey || reflect.DeepEqual(expectValue, resultValue) == false {
		t.Error("expect:", expectKey, expectValue, "get result:", resultKey, resultValue)
	}
}

func Test_addLegacyACL(t *testing.T) {
	tf := NewTransform()
	var restObject types.M
	var result types.M
	var expect types.M
	/*************************************************/
	restObject = nil
	result = tf.addLegacyACL(restObject)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	restObject = types.M{}
	result = tf.addLegacyACL(restObject)
	expect = types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	restObject = types.M{
		"key": "value",
	}
	result = tf.addLegacyACL(restObject)
	expect = types.M{
		"key": "value",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	restObject = types.M{
		"key":    "value",
		"_wperm": types.S{"id", "role"},
	}
	result = tf.addLegacyACL(restObject)
	expect = types.M{
		"key":    "value",
		"_wperm": types.S{"id", "role"},
		"_acl": types.M{
			"id": types.M{
				"w": true,
			},
			"role": types.M{
				"w": true,
			},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	restObject = types.M{
		"key":    "value",
		"_rperm": types.S{"id", "role"},
	}
	result = tf.addLegacyACL(restObject)
	expect = types.M{
		"key":    "value",
		"_rperm": types.S{"id", "role"},
		"_acl": types.M{
			"id": types.M{
				"r": true,
			},
			"role": types.M{
				"r": true,
			},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	restObject = types.M{
		"key":    "value",
		"_rperm": types.S{"id", "role"},
		"_wperm": types.S{"id"},
	}
	result = tf.addLegacyACL(restObject)
	expect = types.M{
		"key":    "value",
		"_rperm": types.S{"id", "role"},
		"_wperm": types.S{"id"},
		"_acl": types.M{
			"id": types.M{
				"r": true,
				"w": true,
			},
			"role": types.M{
				"r": true,
			},
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
}

func Test_transformWhere(t *testing.T) {
	tf := NewTransform()
	var where types.M
	var schema types.M
	var result types.M
	var err error
	var expect types.M
	tmpTimeStr := utils.TimetoString(time.Now().UTC())
	tmpTime, _ := utils.StringtoTime(tmpTimeStr)
	/*************************************************/
	where = nil
	schema = types.M{}
	result, err = tf.transformWhere("", where, schema)
	expect = nil
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	where = types.M{"key": types.M{"key": "value"}}
	schema = types.M{}
	result, err = tf.transformWhere("", where, schema)
	expectErr := errs.E(errs.InvalidJSON, "You cannot use this value as a query parameter.")
	if reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "get result:", err)
	}
	/*************************************************/
	where = types.M{
		"objectId":             "1024",
		"createdAt":            tmpTimeStr,
		"authData.facebook.id": "1024",
		"user":                 "jack",
		"number": types.M{
			"$lt": 25,
			"$gt": 20,
		},
		"skill": "one",
		"$and": types.S{
			types.M{"name": "joe"},
			types.M{"age": 25},
		},
		"key": "value",
	}
	schema = types.M{
		"fields": types.M{
			"user": types.M{
				"type": "Pointer",
			},
			"skill": types.M{
				"type": "Array",
			},
		},
	}
	result, err = tf.transformWhere("", where, schema)
	expect = types.M{
		"_id":                    "1024",
		"_created_at":            tmpTime,
		"_auth_data_facebook.id": "1024",
		"_p_user":                "jack",
		"number": types.M{
			"$lt": 25,
			"$gt": 20,
		},
		"skill": types.M{
			"$all": types.S{"one"},
		},
		"$and": types.S{
			types.M{"name": "joe"},
			types.M{"age": 25},
		},
		"key": "value",
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
}

func Test_transformUpdate(t *testing.T) {
	tf := NewTransform()
	var className string
	var update types.M
	var parseFormatSchema types.M
	var result types.M
	var err error
	var expect types.M
	/*************************************************/
	className = "post"
	update = nil
	parseFormatSchema = types.M{}
	result, err = tf.transformUpdate(className, update, parseFormatSchema)
	expect = nil
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	className = "_User"
	update = types.M{
		"_auth_data_facebook": types.M{
			"id": "1024",
		},
	}
	parseFormatSchema = types.M{}
	result, err = tf.transformUpdate(className, update, parseFormatSchema)
	expect = types.M{
		"$set": types.M{
			"_auth_data_facebook": types.M{
				"id": "1024",
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	className = "post"
	update = types.M{
		"_rperm": types.S{"userid"},
		"_wperm": types.S{"userid"},
	}
	parseFormatSchema = types.M{}
	result, err = tf.transformUpdate(className, update, parseFormatSchema)
	expect = types.M{
		"$set": types.M{
			"_rperm": types.S{"userid"},
			"_wperm": types.S{"userid"},
			"_acl": types.M{
				"userid": types.M{
					"r": true,
					"w": true,
				},
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	className = "post"
	update = types.M{
		"number": types.M{
			"__op":   "Increment",
			"amount": 10,
		},
		"user": types.M{
			"__op":   "Increment",
			"amount": 1,
		},
	}
	parseFormatSchema = types.M{}
	result, err = tf.transformUpdate(className, update, parseFormatSchema)
	expect = types.M{
		"$inc": types.M{
			"number": 10,
			"user":   1,
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	className = "_User"
	update = types.M{
		"_auth_data_facebook": types.M{
			"id": "1024",
		},
		"_rperm": types.S{"userid"},
		"_wperm": types.S{"userid"},
		"number": types.M{
			"__op":   "Increment",
			"amount": 10,
		},
		"user": types.M{
			"__op":   "Increment",
			"amount": 1,
		},
		"name": types.M{
			"__op": "Delete",
		},
		"follower": types.M{
			"__op":    "Add",
			"objects": types.S{"one", "two"},
		},
		"like": types.M{
			"__op":    "Remove",
			"objects": types.S{"one", "two"},
		},
		"age": 10,
	}
	parseFormatSchema = types.M{}
	result, err = tf.transformUpdate(className, update, parseFormatSchema)
	expect = types.M{
		"$set": types.M{
			"_auth_data_facebook": types.M{
				"id": "1024",
			},
			"_rperm": types.S{"userid"},
			"_wperm": types.S{"userid"},
			"_acl": types.M{
				"userid": types.M{
					"r": true,
					"w": true,
				},
			},
			"age": 10,
		},
		"$inc": types.M{
			"number": 10,
			"user":   1,
		},
		"$unset": types.M{
			"name": "",
		},
		"$push": types.M{
			"follower": types.M{
				"$each": types.S{"one", "two"},
			},
		},
		"$pullAll": types.M{
			"like": types.S{"one", "two"},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
}

func Test_nestedMongoObjectToNestedParseObject(t *testing.T) {
	tf := NewTransform()
	var mongoObject interface{}
	var result interface{}
	var err error
	var expect interface{}
	/*************************************************/
	mongoObject = nil
	result, err = tf.nestedMongoObjectToNestedParseObject(mongoObject)
	expect = nil
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = "hello"
	result, err = tf.nestedMongoObjectToNestedParseObject(mongoObject)
	expect = "hello"
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = 10.0
	result, err = tf.nestedMongoObjectToNestedParseObject(mongoObject)
	expect = 10.0
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = 10
	result, err = tf.nestedMongoObjectToNestedParseObject(mongoObject)
	expect = 10
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = true
	result, err = tf.nestedMongoObjectToNestedParseObject(mongoObject)
	expect = true
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.S{"hello", "world"}
	result, err = tf.nestedMongoObjectToNestedParseObject(mongoObject)
	expect = types.S{"hello", "world"}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	tmpTimeStr := utils.TimetoString(time.Now().UTC())
	tmpTime, _ := utils.StringtoTime(tmpTimeStr)
	mongoObject = tmpTime
	result, err = tf.nestedMongoObjectToNestedParseObject(mongoObject)
	expect = types.M{
		"__type": "Date",
		"iso":    tmpTimeStr,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = []byte("hello")
	result, err = tf.nestedMongoObjectToNestedParseObject(mongoObject)
	expect = types.M{
		"__type": "Bytes",
		"base64": "aGVsbG8=",
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.M{
		"name":  "joe",
		"age":   25,
		"m":     false,
		"skill": types.S{"skill1", "skill2"},
		"date":  tmpTime,
		"file":  []byte("hello"),
	}
	result, err = tf.nestedMongoObjectToNestedParseObject(mongoObject)
	expect = types.M{
		"name":  "joe",
		"age":   25,
		"m":     false,
		"skill": types.S{"skill1", "skill2"},
		"date": types.M{
			"__type": "Date",
			"iso":    tmpTimeStr,
		},
		"file": types.M{
			"__type": "Bytes",
			"base64": "aGVsbG8=",
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = []string{"hello"}
	result, err = tf.nestedMongoObjectToNestedParseObject(mongoObject)
	expect = errs.E(errs.InternalServerError, "unknown object type")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "get result:", err)
	}
}

func Test_mongoObjectToParseObject(t *testing.T) {
	tf := NewTransform()
	var mongoObject interface{}
	var schema types.M
	var result interface{}
	var err error
	var expect interface{}
	tmpTimeStr := utils.TimetoString(time.Now().UTC())
	tmpTime, _ := utils.StringtoTime(tmpTimeStr)
	/*************************************************/
	mongoObject = nil
	schema = types.M{}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = nil
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = "hello"
	schema = types.M{}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = "hello"
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = 10.0
	schema = types.M{}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = 10.0
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = 10
	schema = types.M{}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = 10
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = true
	schema = types.M{}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = true
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.S{"hello", "world"}
	schema = types.M{}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = types.S{"hello", "world"}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = tmpTime
	schema = types.M{}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = types.M{
		"__type": "Date",
		"iso":    tmpTimeStr,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = []byte("hello")
	schema = types.M{}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = types.M{
		"__type": "Bytes",
		"base64": "aGVsbG8=",
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.M{
		"_rperm": types.S{"userid"},
		"_wperm": types.S{"userid"},
	}
	schema = types.M{}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = types.M{
		"_rperm": types.S{"userid"},
		"_wperm": types.S{"userid"},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.M{
		"_id":              "1024",
		"_hashed_password": "password",
		"_acl":             "acl",
		"_session_token":   "abc",
		"_updated_at":      tmpTime,
		"_created_at":      tmpTime,
	}
	schema = types.M{}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = types.M{
		"objectId":         "1024",
		"_hashed_password": "password",
		"sessionToken":     "abc",
		"updatedAt":        tmpTimeStr,
		"createdAt":        tmpTimeStr,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.M{
		"_auth_data_facebook": types.M{
			"id": "1024",
		},
	}
	schema = types.M{}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = types.M{
		"authData": types.M{
			"facebook": types.M{
				"id": "1024",
			},
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.M{
		"_p_post": "abc$123",
	}
	schema = nil
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.M{
		"_p_post": "abc$123",
	}
	schema = types.M{}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.M{
		"_p_post": "abc$123",
	}
	schema = types.M{
		"fields": types.M{},
	}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.M{
		"_p_post": "abc$123",
	}
	schema = types.M{
		"fields": types.M{
			"post": types.M{
				"type": "Date",
			},
		},
	}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.M{
		"_p_post": nil,
	}
	schema = types.M{
		"fields": types.M{
			"post": types.M{
				"type": "Pointer",
			},
		},
	}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.M{
		"_p_post": "abc",
	}
	schema = types.M{
		"fields": types.M{
			"post": types.M{
				"type": "Pointer",
			},
		},
	}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.M{
		"_p_post": "abc$123",
	}
	schema = types.M{
		"fields": types.M{
			"post": types.M{
				"type":        "Pointer",
				"targetClass": "def",
			},
		},
	}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = errs.E(errs.InternalServerError, "pointer to incorrect className")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	mongoObject = types.M{
		"_p_post": "abc$123",
	}
	schema = types.M{
		"fields": types.M{
			"post": types.M{
				"type":        "Pointer",
				"targetClass": "abc",
			},
		},
	}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = types.M{
		"post": types.M{
			"__type":    "Pointer",
			"className": "abc",
			"objectId":  "123",
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.M{
		"_other": "hello",
	}
	schema = types.M{}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = errs.E(errs.InternalServerError, "bad key in untransform: "+"_other")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "get result:", err)
	}
	/*************************************************/
	mongoObject = types.M{
		"icon": "hello.jpg",
	}
	schema = types.M{
		"fields": types.M{
			"icon": types.M{
				"type": "File",
			},
		},
	}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = types.M{
		"icon": types.M{
			"__type": "File",
			"name":   "hello.jpg",
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.M{
		"location": types.S{30.0, 40.0},
	}
	schema = types.M{
		"fields": types.M{
			"location": types.M{
				"type": "GeoPoint",
			},
		},
	}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = types.M{
		"location": types.M{
			"__type":    "GeoPoint",
			"longitude": 30.0,
			"latitude":  40.0,
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.M{
		"key": "value",
	}
	schema = types.M{}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = types.M{
		"key": "value",
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = types.M{
		"key": "value",
	}
	schema = types.M{
		"fields": types.M{
			"user": types.M{
				"type":        "Relation",
				"targetClass": "_User",
			},
		},
	}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = types.M{
		"key": "value",
		"user": types.M{
			"__type":    "Relation",
			"className": "_User",
		},
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	/*************************************************/
	mongoObject = []string{"hello"}
	schema = types.M{}
	result, err = tf.mongoObjectToParseObject("", mongoObject, schema)
	expect = errs.E(errs.InternalServerError, "unknown object type")
	if reflect.DeepEqual(expect, err) == false {
		t.Error("expect:", expect, "get result:", err)
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
