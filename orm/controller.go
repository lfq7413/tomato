package orm

var adapter *MongoAdapter
var schemaPromise *Schema

func init() {
	adapter = &MongoAdapter{
		collectionList: []string{},
	}
}

// AdaptiveCollection ...
func AdaptiveCollection(className string) *MongoCollection {
	return adapter.adaptiveCollection(className)
}

// SchemaCollection 获取 Schema 表
func SchemaCollection() *MongoSchemaCollection {
	return adapter.schemaCollection()
}

// CollectionExists ...
func CollectionExists(className string) bool {
	return adapter.collectionExists(className)
}

// DropCollection ...
func DropCollection(className string) error {
	return adapter.dropCollection(className)
}

// Find ...
func Find(className string, where map[string]interface{}, options map[string]interface{}) []interface{} {
	// TODO
	return []interface{}{}
}

// Count ...
func Count(className string, where map[string]interface{}, options map[string]interface{}) int {
	// TODO
	return 0
}

// Destroy ...
func Destroy(className string, where map[string]interface{}, options map[string]interface{}) {
	// TODO
}

// Update ...
func Update(className string, where map[string]interface{}, data map[string]interface{}, options map[string]interface{}) error {
	// TODO
	return nil
}

// UpdateAll ...
func UpdateAll(className string, where map[string]interface{}, data map[string]interface{}, options map[string]interface{}) error {
	// TODO
	return nil
}

// Create ...
func Create(className string, data map[string]interface{}, options map[string]interface{}) error {
	// TODO
	return nil
}

// ValidateObject ...
func ValidateObject(className string, object, where, options map[string]interface{}) error {
	// TODO 处理错误
	schema := LoadSchema(nil)
	acl := []string{}
	if options["acl"] != nil {
		if v, ok := options["acl"].([]string); ok {
			acl = v
		}
	}

	canAddField(schema, className, object, acl)

	schema.validateObject(className, object, where)

	return nil
}

// LoadSchema 加载 Schema
func LoadSchema(acceptor func(*Schema) bool) *Schema {
	if schemaPromise == nil {
		collection := SchemaCollection()
		schemaPromise = Load(collection)
		return schemaPromise
	}

	if acceptor == nil {
		return schemaPromise
	}
	if acceptor(schemaPromise) {
		return schemaPromise
	}

	collection := SchemaCollection()
	schemaPromise = Load(collection)
	return schemaPromise
}

// canAddField ...
func canAddField(schema *Schema, className string, object map[string]interface{}, acl []string) {
	// TODO
}
