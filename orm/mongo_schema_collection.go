package orm

import (
	"github.com/lfq7413/tomato/types"
)

// MongoSchemaCollection ...
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

// FindAndDeleteSchema ...
func (m *MongoSchemaCollection) FindAndDeleteSchema(name string) (types.M, error) {
	result, err := m.FindSchema(name)
	if err != nil {
		return nil, err
	}
	err = m.collection.deleteOne(mongoSchemaQueryFromNameQuery(name, nil))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MongoSchemaCollection) addSchema(name string, fields types.M) error {
	mongoObject := mongoSchemaObjectFromNameFields(name, fields)
	return m.collection.insertOne(mongoObject)
}

func (m *MongoSchemaCollection) updateSchema(name string, update types.M) error {
	return m.collection.updateOne(mongoSchemaQueryFromNameQuery(name, nil), update)
}

func (m *MongoSchemaCollection) upsertSchema(name string, query, update types.M) error {
	return m.collection.upsertOne(mongoSchemaQueryFromNameQuery(name, query), update)
}

func mongoSchemaQueryFromNameQuery(name string, query types.M) types.M {
	return mongoSchemaObjectFromNameFields(name, query)
}

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
