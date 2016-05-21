package orm

import (
	"strings"

	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
	"gopkg.in/mgo.v2"
)

// MongoSchemaCollection _SCHEMA 表操作对象
type MongoSchemaCollection struct {
	collection *MongoCollection
}

// GetAllSchemas 获取所有 Schema，并转换为 API 格式的数据
func (m *MongoSchemaCollection) GetAllSchemas() ([]types.M, error) {
	results, err := m.collection.rawFind(types.M{}, types.M{})
	if err != nil {
		return nil, err
	}
	apiResults := []types.M{}
	for _, result := range results {
		apiResults = append(apiResults, MongoSchemaToParseSchema(result))
	}
	return m.collection.rawFind(types.M{}, types.M{})
}

// FindSchema 查找指定的 Schema
func (m *MongoSchemaCollection) FindSchema(name string) (types.M, error) {
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
	return results[0], nil
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

// addSchema 添加一个表定义
func (m *MongoSchemaCollection) addSchema(name string, fields types.M) error {
	mongoObject := mongoSchemaObjectFromNameFields(name, fields)
	return m.collection.InsertOne(mongoObject)
}

// updateSchema 更新一个表定义
func (m *MongoSchemaCollection) updateSchema(name string, update types.M) error {
	return m.collection.UpdateOne(mongoSchemaQueryFromNameQuery(name, nil), update)
}

// upsertSchema 更新或者插入一个表定义
func (m *MongoSchemaCollection) upsertSchema(name string, query, update types.M) error {
	return m.collection.upsertOne(mongoSchemaQueryFromNameQuery(name, query), update)
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
