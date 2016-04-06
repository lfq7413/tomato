package orm

import "github.com/lfq7413/tomato/utils"

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
func Create(className string, data, options map[string]interface{}) error {
	// TODO 处理错误
	isMaster := false
	aclGroup := []string{}
	if options["acl"] == nil {
		isMaster = true
	} else {
		aclGroup = options["acl"].([]string)
	}

	validateClassName(className)

	schema := LoadSchema(nil)
	if isMaster == false {
		schema.validatePermission(className, aclGroup, "create")
	}

	handleRelationUpdates(className, "", data)

	coll := AdaptiveCollection(className)
	mongoObject := transformCreate(schema, className, data)
	coll.insertOne(mongoObject)

	return nil
}

func validateClassName(className string) {
	// TODO
}

func handleRelationUpdates(className, objectID string, updatemap map[string]interface{}) {
	// TODO
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
	// TODO 处理错误
	if schema.data[className] == nil {
		return
	}
	classSchema := utils.MapInterface(schema.data[className])

	schemaFields := []string{}
	for k := range classSchema {
		schemaFields = append(schemaFields, k)
	}
	// 收集新增的字段
	newKeys := []string{}
	for k := range object {
		t := true
		for _, v := range schemaFields {
			if k == v {
				t = false
				break
			}
		}
		if t {
			newKeys = append(newKeys, k)
		}
	}

	if len(newKeys) > 0 {
		schema.validatePermission(className, acl, "addField")
	}
}
