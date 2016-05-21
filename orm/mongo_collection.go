package orm

import (
	"strings"

	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
	"gopkg.in/mgo.v2"
)

// MongoCollection mongo 表操作对象
type MongoCollection struct {
	collection *mgo.Collection
}

// Find 执行查找操作，自动添加索引
func (m *MongoCollection) Find(query interface{}, options types.M) []types.M {
	result, err := m.rawFind(query, options)
	if err != nil || result == nil {
		return []types.M{}
	}
	// TODO 添加 geo 索引
	return result
}

// rawFind 执行原始查找操作，查找选项包括 sort、skip、limit
func (m *MongoCollection) rawFind(query interface{}, options types.M) ([]types.M, error) {
	q := m.collection.Find(query)
	if options["sort"] != nil {
		if sort, ok := options["sort"].([]string); ok {
			q = q.Sort(sort...)
		}
	}
	if options["skip"] != nil {
		if skip, ok := options["skip"].(float64); ok {
			q = q.Skip(int(skip))
		}
	}
	if options["limit"] != nil {
		if limit, ok := options["limit"].(float64); ok {
			q = q.Limit(int(limit))
		}
	}
	var result []types.M
	err := q.All(&result)
	return result, err
}

// Count 执行 count 操作，
func (m *MongoCollection) Count(query interface{}, options types.M) int {
	q := m.collection.Find(query)
	if options["sort"] != nil {
		if sort, ok := options["sort"].([]string); ok {
			q = q.Sort(sort...)
		}
	}
	if options["skip"] != nil {
		if skip, ok := options["skip"].(float64); ok {
			q = q.Skip(int(skip))
		}
	}
	if options["limit"] != nil {
		if limit, ok := options["limit"].(float64); ok {
			q = q.Limit(int(limit))
		}
	}
	n, err := q.Count()
	if err != nil {
		return 0
	}
	return n
}

// FindOneAndUpdate 查找并更新一个对象，返回更新后的对象
func (m *MongoCollection) FindOneAndUpdate(selector interface{}, update interface{}) types.M {

	var result types.M
	change := mgo.Change{
		Update:    update,
		ReturnNew: true,
	}
	info, err := m.collection.Find(selector).Apply(change, &result)
	if err != nil || info.Updated == 0 {
		return types.M{}
	}

	return result
}

// InsertOne 插入一个对象
func (m *MongoCollection) InsertOne(docs interface{}) error {
	return m.collection.Insert(docs)
}

// upsertOne 更新一个对象，如果要更新的对象不存在，则插入该对象
func (m *MongoCollection) upsertOne(selector interface{}, update interface{}) error {
	_, err := m.collection.Upsert(selector, update)
	return err
}

// UpdateOne 更新一个对象
func (m *MongoCollection) UpdateOne(selector interface{}, update interface{}) error {
	return m.collection.Update(selector, update)
}

// UpdateMany 更新多个对象
func (m *MongoCollection) UpdateMany(selector interface{}, update interface{}) error {
	_, err := m.collection.UpdateAll(selector, update)
	return err
}

// deleteOne 删除一个对象
func (m *MongoCollection) deleteOne(selector interface{}) error {
	return m.collection.Remove(selector)
}

// deleteMany 删除多个对象
func (m *MongoCollection) deleteMany(selector interface{}) (int, error) {
	info, err := m.collection.RemoveAll(selector)
	if err != nil {
		return 0, err
	}
	n := info.Removed
	return n, nil
}

// Drop 删除当前表
func (m *MongoCollection) Drop() error {
	return m.collection.DropCollection()
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
