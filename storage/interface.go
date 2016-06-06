package storage

import "github.com/lfq7413/tomato/types"

// Schema 上层需要实现的 Schema 接口，用于 Transform 中
type Schema interface {
	GetExpectedType(className, key string) types.M
}

// Transform API 格式与数据库格式之间转换的接口
type Transform interface {
	TransformKey(className, fieldName string, schema types.M) string
	TransformWhere(className string, where, schema types.M) (types.M, error)
	TransformUpdate(schema Schema, className string, update, options types.M) (types.M, error)
}

// Adapter 数据库操作适配器接口
type Adapter interface {
	AdaptiveCollection(name string) Collection
	SchemaCollection() SchemaCollection
	CollectionExists(name string) bool
	DeleteOneSchema(name string) error
	DeleteAllSchemas() error
	DeleteFields(className string, fieldNames, pointerFieldNames []string) error
	CreateObject(className string, object types.M, parseFormatSchema types.M) error
	GetTransform() Transform
	GetAllSchemas() ([]types.M, error)
	GetOneSchema(className string) (types.M, error)
	DeleteObjectsByQuery(className string, query types.M, schema types.M) error
	Find(className string, query, schema, options types.M) ([]types.M, error)
	Count(className string, query, schema types.M) (int, error)
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
