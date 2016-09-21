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
}

func newMongoSchemaCollection(collection *MongoCollection) *MongoSchemaCollection {
	return &MongoSchemaCollection{
		collection: collection,
	}
}

// getAllSchemas 获取所有 Schema，并转换为 API 格式的数据
func (m *MongoSchemaCollection) getAllSchemas() ([]types.M, error) {
	results, err := m.collection.rawFind(types.M{}, types.M{})
	if err != nil {
		return nil, err
	}
	apiResults := []types.M{}
	if results == nil {
		return apiResults, nil
	}
	for _, result := range results {
		apiResults = append(apiResults, mongoSchemaToParseSchema(result))
	}
	return apiResults, nil
}

// findSchema 查找指定的 Schema
func (m *MongoSchemaCollection) findSchema(name string) (types.M, error) {
	options := types.M{
		"limit": 1,
	}
	results, err := m.collection.rawFind(mongoSchemaQueryFromNameQuery(name, nil), options)
	if err != nil {
		return nil, err
	}
	if results == nil || len(results) == 0 {
		return types.M{}, nil
	}
	return mongoSchemaToParseSchema(results[0]), nil
}

// findAndDeleteSchema 查找并删除指定的表定义
func (m *MongoSchemaCollection) findAndDeleteSchema(name string) (types.M, error) {

	var result types.M
	change := mgo.Change{
		Remove:    true,
		ReturnNew: false,
	}
	selector := mongoSchemaQueryFromNameQuery(name, nil)
	info, err := m.collection.collection.Find(selector).Apply(change, &result)
	if err != nil {
		if err.Error() == "not found" {
			return types.M{}, nil
		}
		return nil, err
	}
	if info.Removed == 0 {
		return types.M{}, nil
	}

	return result, nil
}

// addSchema 添加一个表定义
// 仅用于测试，逻辑转移到 MongoStorageAdapter/CreateClass
func (m *MongoSchemaCollection) addSchema(name string, fields types.M, classLevelPermissions types.M) (types.M, error) {
	mongoSchema := mongoSchemaFromFieldsAndClassNameAndCLP(fields, name, classLevelPermissions)
	mongoObject := mongoSchemaQueryFromNameQuery(name, mongoSchema)
	// 处理 insertOne 失败的情况，数据库插入失败，检测是否是因为键值重复造成的错误
	err := m.collection.insertOne(mongoObject)
	if err != nil {
		if errs.GetErrorCode(err) == errs.DuplicateValue {
			return nil, errs.E(errs.DuplicateValue, "Class already exists.")
		}
		return nil, err
	}
	return mongoSchemaToParseSchema(mongoObject), err
}

// updateSchema 更新一个表定义
func (m *MongoSchemaCollection) updateSchema(name string, update types.M) error {
	return m.collection.updateOne(mongoSchemaQueryFromNameQuery(name, nil), update)
}

// upsertSchema 更新或者插入一个表定义
func (m *MongoSchemaCollection) upsertSchema(name string, query, update types.M) error {
	return m.collection.upsertOne(mongoSchemaQueryFromNameQuery(name, query), update)
}

// addFieldIfNotExists 更新字段
func (m *MongoSchemaCollection) addFieldIfNotExists(className string, fieldName string, fieldType types.M) error {
	schema, err := m.findSchema(className)
	if err != nil {
		return err
	}
	if schema == nil || len(schema) == 0 {
		return nil
	}

	if fieldType != nil && utils.S(fieldType["type"]) == "GeoPoint" {
		if fields := utils.M(schema["fields"]); fields != nil {
			for _, v := range fields {
				if existingField := utils.M(v); existingField != nil {
					if utils.S(existingField["type"]) == "GeoPoint" {
						return errs.E(errs.IncorrectType, "MongoDB only supports one GeoPoint field in a class.")
					}
				}
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
	object := types.M{
		"_id": name,
	}
	if query != nil {
		for k, v := range query {
			object[k] = v
		}
	}
	return object
}

// mongoFieldToParseSchemaField 把数据库格式的字段类型转换为 API 格式
func mongoFieldToParseSchemaField(t string) types.M {
	if t == "" {
		return types.M{}
	}
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
	response := types.M{}
	if schema == nil {
		return response
	}
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
	// 转换普通字段
	for _, v := range fieldNames {
		response[v] = mongoFieldToParseSchemaField(utils.S(schema[v]))
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

var emptyCLPS = types.M{
	"find":     types.M{},
	"get":      types.M{},
	"create":   types.M{},
	"update":   types.M{},
	"delete":   types.M{},
	"addField": types.M{},
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

// mongoSchemaToParseSchema 把数据库格式的数据转换为 API 格式
func mongoSchemaToParseSchema(schema types.M) types.M {
	if schema == nil {
		return types.M{}
	}
	// 复制 schema["_metadata"]["class_permissions"] 到 classLevelPermissions 中
	var clps types.M
	clps = utils.CopyMap(defaultCLPS)
	if metadata := utils.M(schema["_metadata"]); metadata != nil {
		if classPermissions := utils.M(metadata["class_permissions"]); classPermissions != nil {
			// clps = utils.CopyMap(emptyCLPS)
			// 不存在的 action 默认为公共权限
			for k, v := range classPermissions {
				clps[k] = v
			}
		}
	}

	return types.M{
		"className":             schema["_id"],
		"fields":                mongoSchemaFieldsToParseSchemaFields(schema),
		"classLevelPermissions": clps,
	}
}

// parseFieldTypeToMongoFieldType 返回数据库中存储的字段类型
func parseFieldTypeToMongoFieldType(t types.M) string {
	if t == nil {
		return ""
	}
	fieldType := utils.S(t["type"])
	targetClass := ""
	if fieldType == "Pointer" || fieldType == "Relation" {
		targetClass = utils.S(t["targetClass"])
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
