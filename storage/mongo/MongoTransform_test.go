package mongo

import (
	"testing"
	"time"

	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

func Test_transformKey(t *testing.T) {
	tf := NewTransform()
	var fieldName string
	var schema types.M
	var result string
	/*********************case 01*********************/
	fieldName = "objectId"
	result = tf.transformKey("", fieldName, schema)
	if result != "_id" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
	/*********************case 02*********************/
	fieldName = "createdAt"
	result = tf.transformKey("", fieldName, schema)
	if result != "_created_at" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
	/*********************case 03*********************/
	fieldName = "updatedAt"
	result = tf.transformKey("", fieldName, schema)
	if result != "_updated_at" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
	/*********************case 04*********************/
	fieldName = "sessionToken"
	result = tf.transformKey("", fieldName, schema)
	if result != "_session_token" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
	/*********************case 05*********************/
	schema = nil
	fieldName = "user"
	result = tf.transformKey("", fieldName, schema)
	if result != "user" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
	/*********************case 06*********************/
	schema = types.M{
		"fields": "type is string",
	}
	fieldName = "user"
	result = tf.transformKey("", fieldName, schema)
	if result != "user" {
		t.Error("transform:", fieldName, "error!", "result:", result)
	}
	/*********************case 07*********************/
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
	/*********************case 08*********************/
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
	/*********************case 09*********************/
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
	// transformTopLevelAtom
	// transformInteriorValue
	// transformUpdateOperator
	// TODO
}

func Test_transformQueryKeyValue(t *testing.T) {
	// transformWhere
	// transformConstraint
	// transformTopLevelAtom
	// TODO
}

func Test_transformConstraint(t *testing.T) {
	// transformInteriorAtom
	// transformTopLevelAtom
	// TODO
}

func Test_transformTopLevelAtom(t *testing.T) {
	// dateCoder
	// bytesCoder
	// geoPointCoder
	// fileCoder
	// TODO
}

func Test_transformUpdateOperator(t *testing.T) {
	// transformInteriorAtom
	// TODO
}

func Test_parseObjectToMongoObjectForCreate(t *testing.T) {
	// transformAuthData
	// transformACL
	// parseObjectKeyValueToMongoObjectKeyValue
	// TODO
}

func Test_parseObjectKeyValueToMongoObjectKeyValue(t *testing.T) {
	// transformTopLevelAtom
	// transformInteriorValue
	// transformUpdateOperator
	// TODO
}

func Test_transformAuthData(t *testing.T) {
	// TODO
}

func Test_transformACL(t *testing.T) {
	// TODO
}

func Test_transformWhere(t *testing.T) {
	// transformQueryKeyValue
	// TODO
}

func Test_transformUpdate(t *testing.T) {
	// transformAuthData
	// transformACL
	// transformKeyValueForUpdate
	// TODO
}

func Test_nestedMongoObjectToNestedParseObject(t *testing.T) {
	// TODO
}

func Test_mongoObjectToParseObject(t *testing.T) {
	// nestedMongoObjectToNestedParseObject
	// dateCoder
	// bytesCoder
	// untransformACL
	// TODO
}

func Test_untransformACL(t *testing.T) {
	// TODO
}

func Test_transformInteriorAtom(t *testing.T) {
	// dateCoder
	// bytesCoder
	// TODO
}

func Test_transformInteriorValue(t *testing.T) {
	// transformInteriorAtom
	// TODO
}

func Test_dateCoder(t *testing.T) {
	// TODO
}

func Test_bytesCoder(t *testing.T) {
	// TODO
}

func Test_geoPointCoder(t *testing.T) {
	// TODO
}

func Test_fileCoder(t *testing.T) {
	// TODO
}

func Test_valueAsDate(t *testing.T) {
	var value interface{}
	var date time.Time
	var ok bool
	/*********************case 01*********************/
	value = 1024
	date, ok = valueAsDate(value)
	if ok {
		t.Error("value:", value, "ok:", ok, "date:", date)
	}
	/*********************case 02*********************/
	value = "Incorrect string time"
	date, ok = valueAsDate(value)
	if ok {
		t.Error("value:", value, "ok:", ok, "date:", date)
	}
	/*********************case 03*********************/
	value = "2006-01-02T15:04:05.000Z"
	date, ok = valueAsDate(value)
	if !ok || utils.TimetoString(date) != "2006-01-02T15:04:05.000Z" {
		t.Error("value:", value, "ok:", ok, "date:", date)
	}
	/*********************case 04*********************/
	value = time.Now().UTC()
	date, ok = valueAsDate(value)
	if !ok {
		t.Error("value:", value, "ok:", ok, "date:", date)
	}
}
