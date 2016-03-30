package orm

import "gopkg.in/mgo.v2/bson"

// Schema ...
type Schema struct {
	collection *MongoSchemaCollection
}

// AddClassIfNotExists 添加类定义
func (s *Schema) AddClassIfNotExists(className string, fields bson.M, classLevelPermissions bson.M) bson.M {
	// TODO
	return nil
}

func (s *Schema) reloadData() {
	// TODO
}

// MongoSchemaToSchemaAPIResponse ...
func MongoSchemaToSchemaAPIResponse(bson.M) bson.M {
	// TODO
	return nil
}

// Load 返回一个新的 Schema 结构体
func Load(collection *MongoSchemaCollection) *Schema {
	schema := &Schema{
		collection: collection,
	}
	schema.reloadData()
	return schema
}
