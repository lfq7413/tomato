package storage

import "github.com/lfq7413/tomato/types"

// Adapter 数据库操作适配器接口
type Adapter interface {
	ClassExists(name string) bool
	SetClassLevelPermissions(className string, CLPs types.M) error
	CreateClass(className string, schema types.M) (types.M, error)
	AddFieldIfNotExists(className, fieldName string, fieldType types.M) error
	DeleteClass(className string) (types.M, error)
	DeleteAllClasses() error
	DeleteFields(className string, schema types.M, fieldNames []string) error
	CreateObject(className string, schema, object types.M) error
	GetAllClasses() ([]types.M, error)
	GetClass(className string) (types.M, error)
	DeleteObjectsByQuery(className string, schema, query types.M) error
	Find(className string, schema, query, options types.M) ([]types.M, error)
	Count(className string, schema, query types.M) (int, error)
	UpdateObjectsByQuery(className string, schema, query, update types.M) error
	FindOneAndUpdate(className string, schema, query, update types.M) (types.M, error)
	UpsertOneObject(className string, schema, query, update types.M) error
	EnsureUniqueness(className string, schema types.M, fieldNames []string) error
}
