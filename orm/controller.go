package orm

import (
	"github.com/lfq7413/tomato/utils"
	"gopkg.in/mgo.v2/bson"
)

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
	// TODO 处理错误
	isMaster := false
	aclGroup := []string{}
	if options["acl"] == nil {
		isMaster = true
	} else {
		aclGroup = options["acl"].([]string)
	}

	schema := LoadSchema(nil)
	if isMaster == false {
		schema.validatePermission(className, aclGroup, "delete")
	}

	coll := AdaptiveCollection(className)
	mongoWhere := transformWhere(schema, className, where)
	// 组装查询条件，查找可被当前用户修改的对象
	if options["acl"] != nil {
		writePerms := []interface{}{}
		perm := bson.M{
			"_wperm": bson.M{"$exists": false},
		}
		writePerms = append(writePerms, perm)
		for _, acl := range aclGroup {
			perm = bson.M{
				"_wperm": bson.M{"$in": []string{acl}},
			}
			writePerms = append(writePerms, perm)
		}

		mongoWhere = bson.M{
			"$and": []interface{}{
				mongoWhere,
				bson.M{"$or": writePerms},
			},
		}
	}
	coll.deleteMany(mongoWhere)
	// TODO 处理返回错误
}

// Update ...
func Update(className string, where, data, options map[string]interface{}) (bson.M, error) {
	// TODO 处理错误
	data = utils.CopyMap(data)
	acceptor := func(schema *Schema) bool {
		keys := []string{}
		for k := range where {
			keys = append(keys, k)
		}
		return schema.hasKeys(className, keys)
	}
	isMaster := false
	aclGroup := []string{}
	if options["acl"] == nil {
		isMaster = true
	} else {
		aclGroup = options["acl"].([]string)
	}

	schema := LoadSchema(acceptor)
	if isMaster == false {
		schema.validatePermission(className, aclGroup, "update")
	}
	handleRelationUpdates(className, utils.String(where["objectId"]), data)

	coll := AdaptiveCollection(className)
	mongoWhere := transformWhere(schema, className, where)
	// 组装查询条件，查找可被当前用户修改的对象
	if options["acl"] != nil {
		writePerms := []interface{}{}
		perm := bson.M{
			"_wperm": bson.M{"$exists": false},
		}
		writePerms = append(writePerms, perm)
		for _, acl := range aclGroup {
			perm = bson.M{
				"_wperm": bson.M{"$in": []string{acl}},
			}
			writePerms = append(writePerms, perm)
		}

		mongoWhere = bson.M{
			"$and": []interface{}{
				mongoWhere,
				bson.M{"$or": writePerms},
			},
		}
	}
	mongoUpdate := transformUpdate(schema, className, data)

	result := coll.findOneAndUpdate(mongoWhere, mongoUpdate)
	// TODO 处理返回错误

	response := bson.M{}
	if mongoUpdate["$inc"] != nil && utils.MapInterface(mongoUpdate["$inc"]) != nil {
		inc := utils.MapInterface(mongoUpdate["$inc"])
		for k := range inc {
			response[k] = result[k]
		}
	}

	return response, nil
}

// Create ...
func Create(className string, data, options map[string]interface{}) error {
	// TODO 处理错误
	data = utils.CopyMap(data)
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
	// TODO 处理错误
	if ClassNameIsValid(className) == false {
		// TODO 无效类名
		return
	}
}

func handleRelationUpdates(className, objectID string, update map[string]interface{}) {
	// TODO 处理错误
	objID := objectID
	if utils.String(update["objectId"]) != "" {
		objID = utils.String(update["objectId"])
	}

	var process func(op interface{}, key string)
	process = func(op interface{}, key string) {
		if op == nil || utils.MapInterface(op) == nil || utils.MapInterface(op)["__op"] == nil {
			return
		}
		opMap := utils.MapInterface(op)
		p := utils.String(opMap["__op"])
		if p == "AddRelation" {
			delete(update, key)
			objects := utils.SliceInterface(opMap["objects"])
			for _, object := range objects {
				relationID := utils.String(utils.MapInterface(object)["objectId"])
				addRelation(key, className, objID, relationID)
			}
		} else if p == "RemoveRelation" {
			delete(update, key)
			objects := utils.SliceInterface(opMap["objects"])
			for _, object := range objects {
				relationID := utils.String(utils.MapInterface(object)["objectId"])
				removeRelation(key, className, objID, relationID)
			}
		} else if p == "Batch" {
			ops := utils.SliceInterface(opMap["ops"])
			for _, x := range ops {
				process(x, key)
			}
		}
	}

	for k, v := range update {
		process(v, k)
	}

}

func addRelation(key, fromClassName, fromID, toID string) {
	// TODO 处理错误
	doc := map[string]interface{}{
		"relatedId": toID,
		"owningId":  fromID,
	}
	className := "_Join:" + key + ":" + fromClassName
	coll := AdaptiveCollection(className)
	coll.upsertOne(doc, doc)
}

func removeRelation(key, fromClassName, fromID, toID string) {
	// TODO 处理错误
	doc := map[string]interface{}{
		"relatedId": toID,
		"owningId":  fromID,
	}
	className := "_Join:" + key + ":" + fromClassName
	coll := AdaptiveCollection(className)
	coll.deleteOne(doc)
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
