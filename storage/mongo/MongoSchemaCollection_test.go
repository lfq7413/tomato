package mongo

import (
	"testing"

	"github.com/lfq7413/tomato/types"

	"gopkg.in/mgo.v2"
)

func Test_getAllSchemas(t *testing.T) {
	// mongoSchemaToParseSchema
	// TODO
}

func Test_findSchema(t *testing.T) {
	// mongoSchemaQueryFromNameQuery
	// mongoSchemaToParseSchema
	// TODO
}

func Test_findAndDeleteSchema(t *testing.T) {
	// mongoSchemaQueryFromNameQuery
	// TODO
}

func Test_addSchema(t *testing.T) {
	// mongoSchemaFromFieldsAndClassNameAndCLP
	// mongoSchemaObjectFromNameFields
	// mongoSchemaToParseSchema
	// TODO
}

func Test_updateSchema(t *testing.T) {
	// TODO
}

func Test_upsertSchema(t *testing.T) {
	// TODO
}

func Test_addFieldIfNotExists(t *testing.T) {
	// findSchema
	// upsertSchema
	// TODO
}

func Test_mongoSchemaQueryFromNameQuery(t *testing.T) {
	// mongoSchemaObjectFromNameFields
	// TODO
}

func Test_mongoSchemaObjectFromNameFields(t *testing.T) {
	// TODO
}

func Test_mongoFieldToParseSchemaField(t *testing.T) {
	// TODO
}

func Test_mongoSchemaFieldsToParseSchemaFields(t *testing.T) {
	// mongoFieldToParseSchemaField
	// TODO
}

func Test_mongoSchemaToParseSchema(t *testing.T) {
	// mongoSchemaFieldsToParseSchemaFields
	// TODO
}

func Test_parseFieldTypeToMongoFieldType(t *testing.T) {
	var fieldType types.M
	var result string
	var expect string
	/*****************************************************/
	fieldType = nil
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = ""
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = ""
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type": "Pointer",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "*"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type": "Relation",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "relation<>"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type":        "Pointer",
		"targetClass": "user",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "*user"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type":        "Relation",
		"targetClass": "user",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "relation<user>"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type": "Number",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "number"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type": "String",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "string"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type": "Boolean",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "boolean"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type": "Date",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "date"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type": "Object",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "object"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type": "Array",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "array"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type": "GeoPoint",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "geopoint"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type": "File",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = "file"
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
	/*****************************************************/
	fieldType = types.M{
		"type": "Other",
	}
	result = parseFieldTypeToMongoFieldType(fieldType)
	expect = ""
	if result != expect {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_mongoSchemaFromFieldsAndClassNameAndCLP(t *testing.T) {
	// TODO
}

func getSchemaCollection(db *mgo.Database) *MongoSchemaCollection {
	mc := newMongoCollection(db.C("SCHEMA"))
	msc := newMongoSchemaCollection(mc)
	return msc
}
