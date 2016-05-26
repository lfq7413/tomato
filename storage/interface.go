package storage

// Schema 上层需要实现的 Schema 接口，用于 Transform 中
type Schema interface {
	GetExpectedType(className, key string) map[string]interface{}
	GetRelationFields(className string) map[string]interface{}
}

// Transform API 格式与数据库格式之间转换的接口
type Transform interface {
	TransformKey(schema Schema, className, key string) (string, error)
	TransformWhere(schema Schema, className string, where, options map[string]interface{}) (map[string]interface{}, error)
	TransformUpdate(schema Schema, className string, update, options map[string]interface{}) (map[string]interface{}, error)
	TransformCreate(schema Schema, className string, create map[string]interface{}) (map[string]interface{}, error)
	AddReadACL(mongoWhere interface{}, acl []string) map[string]interface{}
	AddWriteACL(mongoWhere interface{}, acl []string) map[string]interface{}
	UntransformObject(schema Schema, className string, mongoObject interface{}, isNestedObject bool) (interface{}, error)
}

// Adapter 数据库操作适配器接口
type Adapter interface {
	AdaptiveCollection(name string) Collection
	SchemaCollection() SchemaCollection
	CollectionExists(name string) bool
	DropCollection(name string) error
	AllCollections() []Collection
	DeleteFields(className string, fieldNames, pointerFieldNames []string) error
	GetTransform() Transform
}

// Collection 集合操作接口
type Collection interface {
	Find(query interface{}, options map[string]interface{}) []map[string]interface{}
	Count(query interface{}, options map[string]interface{}) int
	FindOneAndUpdate(selector interface{}, update interface{}) map[string]interface{}
	InsertOne(docs interface{})
	UpsertOne(selector interface{}, update interface{}) error
	UpdateMany(selector interface{}, update interface{}) error
	DeleteOne(selector interface{}) error
	DeleteMany(selector interface{}) (int, error)
	Drop() error
}

// SchemaCollection Schema 集合操作接口
type SchemaCollection interface {
	GetAllSchemas() ([]map[string]interface{}, error)
	FindSchema(name string) (map[string]interface{}, error)
	FindAndDeleteSchema(name string) (map[string]interface{}, error)
	AddSchema(name string, fields map[string]interface{}, classLevelPermissions map[string]interface{}) (map[string]interface{}, error)
	UpdateSchema(name string, update map[string]interface{}) error
	UpdateField(className string, fieldName string, fieldType map[string]interface{}) error
}
