package orm

import "gopkg.in/mgo.v2/bson"

// Schema ...
type Schema struct {
}

// AddClassIfNotExists 添加类定义
func (s *Schema) AddClassIfNotExists(className string, fields bson.M, classLevelPermissions bson.M) bson.M {
	return nil
}

// MongoSchemaToSchemaAPIResponse ...
func MongoSchemaToSchemaAPIResponse(bson.M) bson.M {
	// TODO
	return nil
}
