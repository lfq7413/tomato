package mongo

import (
	"strings"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
	"gopkg.in/mgo.v2"
)

// MongoSchemaCollection _SCHEMA 表操作对象
type MongoSchemaCollection struct {
	collection *MongoCollection
	transform  *MongoTransform
}

// GetAllSchemas 获取所有 Schema，并转换为 API 格式的数据
func (m *MongoSchemaCollection) GetAllSchemas() ([]types.M, error) {
	results, err := m.collection.RawFind(types.M{}, types.M{})
	if err != nil {
		return nil, err
	}
	apiResults := []types.M{}
	for _, result := range results {
		apiResults = append(apiResults, MongoSchemaToParseSchema(result))
	}
	return apiResults, nil
}

// FindSchema 查找指定的 Schema
func (m *MongoSchemaCollection) FindSchema(name string) (types.M, error) {
	options := types.M{
		"limit": 1,
	}
	results, err := m.collection.RawFind(mongoSchemaQueryFromNameQuery(name, nil), options)
	if err != nil {
		return nil, err
	}
	if results == nil || len(results) == 0 {
		return types.M{}, nil
	}
	return MongoSchemaToParseSchema(results[0]), nil
}

// FindAndDeleteSchema 查找并删除指定的表定义
func (m *MongoSchemaCollection) FindAndDeleteSchema(name string) (types.M, error) {

	var result types.M
	change := mgo.Change{
		Remove:    true,
		ReturnNew: false,
	}
	selector := mongoSchemaQueryFromNameQuery(name, nil)
	info, err := m.collection.collection.Find(selector).Apply(change, &result)
	if err != nil {
		return nil, err
	}
	if info.Removed == 0 {
		return types.M{}, nil
	}

	return result, nil
}

// AddSchema 添加一个表定义
func (m *MongoSchemaCollection) AddSchema(name string, fields types.M, classLevelPermissions types.M) (types.M, error) {
	mongoSchema, err := mongoSchemaFromFieldsAndClassNameAndCLP(fields, name, classLevelPermissions)
	if err != nil {
		return nil, err
	}
	mongoObject := mongoSchemaObjectFromNameFields(name, mongoSchema)
	return MongoSchemaToParseSchema(mongoObject), m.collection.InsertOne(mongoObject)
}

// UpdateSchema 更新一个表定义
func (m *MongoSchemaCollection) UpdateSchema(name string, update types.M) error {
	return m.collection.UpdateOne(mongoSchemaQueryFromNameQuery(name, nil), update)
}

// upsertSchema 更新或者插入一个表定义
func (m *MongoSchemaCollection) upsertSchema(name string, query, update types.M) error {
	return m.collection.UpsertOne(mongoSchemaQueryFromNameQuery(name, query), update)
}

// AddFieldIfNotExists 更新字段
func (m *MongoSchemaCollection) AddFieldIfNotExists(className string, fieldName string, fieldType types.M) error {
	schema, err := m.FindSchema(className)
	if err != nil {
		return err
	}
	if schema == nil || len(schema) == 0 {
		return nil
	}
	if fieldType["type"].(string) == "GeoPoint" {
		fields := schema["fields"].(map[string]interface{})
		for _, v := range fields {
			existingField := v.(map[string]interface{})
			if existingField["type"].(string) == "GeoPoint" {
				return errs.E(errs.IncorrectType, "MongoDB only supports one GeoPoint field in a class.")
			}
		}
	}

	query := types.M{
		fieldName: types.M{"$exists": false},
	}
	date := types.M{
		fieldName: parseFieldTypeToMongoFieldType(fieldType),
	}
	update := types.M{
		"$set": date,
	}
	return m.upsertSchema(className, query, update)
}

// mongoSchemaQueryFromNameQuery 从表名及查询条件组装 mongo 查询对象
func mongoSchemaQueryFromNameQuery(name string, query types.M) types.M {
	return mongoSchemaObjectFromNameFields(name, query)
}

// mongoSchemaObjectFromNameFields 从表名及字段列表组装 mongo 查询对象
func mongoSchemaObjectFromNameFields(name string, fields types.M) types.M {
	object := types.M{
		"_id": name,
	}
	if fields != nil {
		for k, v := range fields {
			object[k] = v
		}
	}
	return object
}

// mongoFieldToParseSchemaField 把数据库格式的字段类型转换为 API 格式
func mongoFieldToParseSchemaField(t string) types.M {
	// *abc ==> {"type":"Pointer", "targetClass":"abc"}
	if t[0] == '*' {
		return types.M{
			"type":        "Pointer",
			"targetClass": string(t[1:]),
		}
	}
	// relation<abc> ==> {"type":"Relation", "targetClass":"abc"}
	if strings.HasPrefix(t, "relation<") {
		return types.M{
			"type":        "Relation",
			"targetClass": string(t[len("relation<") : len(t)-1]),
		}
	}
	switch t {
	case "number":
		return types.M{
			"type": "Number",
		}
	case "string":
		return types.M{
			"type": "String",
		}
	case "boolean":
		return types.M{
			"type": "Boolean",
		}
	case "date":
		return types.M{
			"type": "Date",
		}
	case "map":
		return types.M{
			"type": "Object",
		}
	case "object":
		return types.M{
			"type": "Object",
		}
	case "array":
		return types.M{
			"type": "Array",
		}
	case "geopoint":
		return types.M{
			"type": "GeoPoint",
		}
	case "file":
		return types.M{
			"type": "File",
		}
	case "bytes":
		return types.M{
			"type": "Bytes",
		}
	}

	return types.M{}
}

var nonFieldSchemaKeys = []string{"_id", "_metadata", "_client_permissions"}

// mongoSchemaFieldsToParseSchemaFields 转换数据库格式的字段到 API类型，排除掉 nonFieldSchemaKeys 中的字段
func mongoSchemaFieldsToParseSchemaFields(schema types.M) types.M {
	fieldNames := []string{}
	for k := range schema {
		t := false
		// 排除 nonFieldSchemaKeys
		for _, v := range nonFieldSchemaKeys {
			if k == v {
				t = true
				break
			}
		}
		if t == false {
			fieldNames = append(fieldNames, k)
		}
	}
	response := types.M{}
	// 转换普通字段
	for _, v := range fieldNames {
		response[v] = mongoFieldToParseSchemaField(utils.String(schema[v]))
	}
	// 转换默认字段
	response["ACL"] = types.M{
		"type": "ACL",
	}
	response["createdAt"] = types.M{
		"type": "Date",
	}
	response["updatedAt"] = types.M{
		"type": "Date",
	}
	response["objectId"] = types.M{
		"type": "String",
	}
	return response
}

// defaultCLPS 默认的类级别权限
var defaultCLPS = types.M{
	"find":     types.M{"*": true},
	"get":      types.M{"*": true},
	"create":   types.M{"*": true},
	"update":   types.M{"*": true},
	"delete":   types.M{"*": true},
	"addField": types.M{"*": true},
}

// MongoSchemaToParseSchema 把数据库格式的数据转换为 API 格式
func MongoSchemaToParseSchema(schema types.M) types.M {
	result := types.M{
		"className": schema["_id"],
		"fields":    mongoSchemaFieldsToParseSchemaFields(schema),
	}

	// 复制 schema["_metadata"]["class_permissions"] 到 classLevelPermissions 中
	classLevelPermissions := utils.CopyMap(defaultCLPS)
	if schema["_metadata"] != nil && utils.MapInterface(schema["_metadata"]) != nil {
		metadata := utils.MapInterface(schema["_metadata"])
		if metadata["class_permissions"] != nil && utils.MapInterface(metadata["class_permissions"]) != nil {
			classPermissions := utils.MapInterface(metadata["class_permissions"])
			for k, v := range classPermissions {
				classLevelPermissions[k] = v
			}
		}
	}
	result["classLevelPermissions"] = classLevelPermissions

	return result
}

// parseFieldTypeToMongoFieldType 返回数据库中存储的字段类型
func parseFieldTypeToMongoFieldType(t types.M) string {
	fieldType := t["type"].(string)
	targetClass := ""
	if fieldType == "Pointer" || fieldType == "Relation" {
		targetClass = t["targetClass"].(string)
	}
	switch fieldType {
	case "Pointer":
		return "*" + targetClass
	case "Relation":
		return "relation<" + targetClass + ">"
	case "Number":
		return "number"
	case "String":
		return "string"
	case "Boolean":
		return "boolean"
	case "Date":
		return "date"
	case "Object":
		return "object"
	case "Array":
		return "array"
	case "GeoPoint":
		return "geopoint"
	case "File":
		return "file"
	default:
		return ""
	}
}

// mongoSchemaFromFieldsAndClassNameAndCLP 把字段属性转换为数据库中保存的类型
func mongoSchemaFromFieldsAndClassNameAndCLP(fields types.M, className string, classLevelPermissions types.M) (types.M, error) {
	mongoObject := types.M{
		"_id":       className,
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
	}

	// 添加其他字段
	for fieldName, v := range fields {
		mongoObject[fieldName] = parseFieldTypeToMongoFieldType(v.(map[string]interface{}))
	}

	// 添加 CLP
	var metadata types.M
	if mongoObject["_metadata"] == nil && utils.MapInterface(mongoObject["_metadata"]) == nil {
		metadata = types.M{}
	} else {
		metadata = utils.MapInterface(mongoObject["_metadata"])
	}
	if classLevelPermissions == nil {
		delete(metadata, "class_permissions")
	} else {
		metadata["class_permissions"] = classLevelPermissions
	}
	mongoObject["_metadata"] = metadata

	return mongoObject, nil
}
