package orm

import "gopkg.in/mgo.v2/bson"

// MongoSchemaCollection ...
type MongoSchemaCollection struct {
	collection *MongoCollection
}

// GetAllSchemas 获取所有 Schema
func (m *MongoSchemaCollection) GetAllSchemas() ([]bson.M, error) {
	return m.collection.rawFind(map[string]interface{}{}, map[string]interface{}{})
}

func (m *MongoSchemaCollection) findSchema(name string) (bson.M, error) {
	options := bson.M{
		"limit": 1,
	}
	results, err := m.collection.rawFind(mongoSchemaQueryFromNameQuery(name, nil), options)
	if err != nil {
		return nil, err
	}
	if results == nil || len(results) == 0 {
		return bson.M{}, nil
	}
	return results[0], nil
}

func (m *MongoSchemaCollection) findAndDeleteSchema(name string) (bson.M, error) {
	result, err := m.findSchema(name)
	if err != nil {
		return nil, err
	}
	err = m.collection.deleteOne(mongoSchemaQueryFromNameQuery(name, nil))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MongoSchemaCollection) addSchema(name string, fields bson.M) error {
	mongoObject := mongoSchemaObjectFromNameFields(name, fields)
	return m.collection.insertOne(mongoObject)
}

func (m *MongoSchemaCollection) updateSchema(name string, update bson.M) error {
	return m.collection.updateOne(mongoSchemaQueryFromNameQuery(name, nil), update)
}

func (m *MongoSchemaCollection) upsertSchema(name string, update bson.M) error {
	return m.collection.upsertOne(mongoSchemaQueryFromNameQuery(name, nil), update)
}

func mongoSchemaQueryFromNameQuery(name string, query bson.M) bson.M {
	return mongoSchemaObjectFromNameFields(name, query)
}

func mongoSchemaObjectFromNameFields(name string, fields bson.M) bson.M {
	object := bson.M{
		"_id": name,
	}
	if fields != nil {
		for k, v := range fields {
			object[k] = v
		}
	}
	return object
}
