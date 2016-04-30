package orm

import (
	"github.com/lfq7413/tomato/types"
	"gopkg.in/mgo.v2"
)

// MongoSchemaCollection _SCHEMA 表操作对象
type MongoSchemaCollection struct {
	collection *MongoCollection
}

// GetAllSchemas 获取所有 Schema
func (m *MongoSchemaCollection) GetAllSchemas() ([]types.M, error) {
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
	return m.collection.insertOne(mongoObject)
}

// updateSchema 更新一个表定义
func (m *MongoSchemaCollection) updateSchema(name string, update types.M) error {
	return m.collection.updateOne(mongoSchemaQueryFromNameQuery(name, nil), update)
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
