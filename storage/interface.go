package storage

import "github.com/lfq7413/tomato/types"

// Adapter 数据库操作适配器接口
type Adapter interface {
	AdaptiveCollection(name string) Collection
	SchemaCollection() SchemaCollection
	CollectionExists(name string) bool
	DeleteOneSchema(name string) error
	DeleteAllSchemas() error
	DeleteFields(className string, fieldNames, pointerFieldNames []string) error
	CreateObject(className string, object types.M, schema types.M) error
	GetAllSchemas() ([]types.M, error)
	GetOneSchema(className string) (types.M, error)
	DeleteObjectsByQuery(className string, query types.M, schema types.M) error
	Find(className string, query, schema, options types.M) ([]types.M, error)
	Count(className string, query, schema types.M) (int, error)
	UpdateObjectsByQuery(className string, query, schema, update types.M) error
	FindOneAndUpdate(className string, query, schema, update types.M) (types.M, error)
	UpsertOneObject(className string, query, schema, update types.M) error
}

// Collection 集合操作接口
type Collection interface {
	Find(query interface{}, options types.M) []types.M
	Count(query interface{}, options types.M) int
	FindOneAndUpdate(selector interface{}, update interface{}) types.M
	InsertOne(docs interface{}) error
	UpsertOne(selector interface{}, update interface{}) error
	UpdateMany(selector interface{}, update interface{}) error
	DeleteOne(selector interface{}) error
	DeleteMany(selector interface{}) (int, error)
	Drop() error
}

// SchemaCollection Schema 集合操作接口
type SchemaCollection interface {
	GetAllSchemas() ([]types.M, error)
	FindSchema(name string) (types.M, error)
	FindAndDeleteSchema(name string) (types.M, error)
	AddSchema(name string, fields types.M, classLevelPermissions types.M) (types.M, error)
	UpdateSchema(name string, update types.M) error
	AddFieldIfNotExists(className string, fieldName string, fieldType types.M) error
}
