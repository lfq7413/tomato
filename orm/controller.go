//Package orm 数据库操作模块，当前只对接了 MongoDB
package orm

import (
	"regexp"
	"strings"

	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

var adapter *MongoAdapter
var schemaPromise *Schema

// init 初始化 Mongo 适配器
func init() {
	adapter = &MongoAdapter{
		collectionList: []string{},
	}
}

// AdaptiveCollection 获取要操作的表，以便后续操作
func AdaptiveCollection(className string) *MongoCollection {
	return adapter.adaptiveCollection(className)
}

// SchemaCollection 获取 Schema 表
func SchemaCollection() *MongoSchemaCollection {
	return adapter.schemaCollection()
}

// CollectionExists 检测表是否存在
func CollectionExists(className string) bool {
	return adapter.collectionExists(className)
}

// DropCollection 删除指定表
func DropCollection(className string) error {
	return adapter.dropCollection(className)
}

// Find 从指定表中查询数据，查询到的数据放入 list 中
// 如果查询的是 count ，结果也会放入 list，并且只有这一个元素
// options 中的选项包括：skip、limit、sort、count、acl
func Find(className string, where, options types.M) (types.S, error) {
	// TODO 处理错误
	if options == nil {
		options = types.M{}
	}
	if where == nil {
		where = types.M{}
	}

	// 组装数据库查询设置项
	mongoOptions := types.M{}
	if options["skip"] != nil {
		mongoOptions["skip"] = options["skip"]
	}
	if options["limit"] != nil {
		mongoOptions["limit"] = options["limit"]
	}

	var isMaster bool
	if _, ok := options["acl"]; ok {
		isMaster = false
	} else {
		// 不存在键值 acl 时，即为 Master
		isMaster = true
	}
	var aclGroup []string
	if options["acl"] == nil {
		aclGroup = []string{}
	} else {
		aclGroup = options["acl"].([]string)
	}

	// 检测查询条件中的 key 在表中是否存在
	acceptor := func(schema *Schema) bool {
		return schema.hasKeys(className, keysForQuery(where))
	}
	schema := LoadSchema(acceptor)

	if options["sort"] != nil {
		sortKeys := []string{}
		keys := options["sort"].([]string)
		for _, key := range keys {
			mongoKey := ""
			// sort 中的 key ，如果是要按倒序排列，则会加前缀 "-" ，所以要对其进行处理
			if strings.HasPrefix(key, "-") {
				mongoKey = "-" + transformKey(schema, className, key[1:])
			} else {
				mongoKey = transformKey(schema, className, key)
			}
			sortKeys = append(sortKeys, mongoKey)
		}
		mongoOptions["sort"] = sortKeys
	}

	// 校验当前用户是否能对表进行 find 或者 get 操作
	if isMaster == false {
		op := "find"
		if len(where) == 1 && where["objectId"] != nil && utils.String(where["objectId"]) != "" {
			op = "get"
		}
		err := schema.validatePermission(className, aclGroup, op)
		if err != nil {
			return nil, err
		}
	}

	// 处理 $relatedTo
	reduceRelationKeys(className, where)
	// 处理 relation 字段上的 $in
	reduceInRelation(className, where, schema)

	coll := AdaptiveCollection(className)
	mongoWhere := transformWhere(schema, className, where)
	// 组装 acl 查询条件，查找可被当前用户访问的对象
	if isMaster == false {
		queryPerms := types.S{}
		// 可查询 不存在读权限字段的
		perm := types.M{
			"_rperm": types.M{"$exists": false},
		}
		queryPerms = append(queryPerms, perm)
		// 可查询 读权限包含 * 的
		perm = types.M{
			"_rperm": types.M{"$in": []string{"*"}},
		}
		queryPerms = append(queryPerms, perm)
		for _, acl := range aclGroup {
			// 可查询 读权限包含 当前用户角色与 id 的
			perm = types.M{
				"_rperm": types.M{"$in": []string{acl}},
			}
			queryPerms = append(queryPerms, perm)
		}

		mongoWhere = types.M{
			"$and": types.S{
				mongoWhere,
				types.M{"$or": queryPerms},
			},
		}
	}

	// 获取 count
	if options["count"] != nil {
		delete(mongoOptions, "limit")
		count := coll.Count(mongoWhere, mongoOptions)
		return types.S{count}, nil
	}

	// 执行查询操作
	mongoResults := coll.Find(mongoWhere, mongoOptions)
	results := types.S{}
	for _, r := range mongoResults {
		result := untransformObject(schema, isMaster, aclGroup, className, r)
		results = append(results, result)
	}
	return results, nil

}

// Destroy ...
func Destroy(className string, where types.M, options types.M) {
	// TODO 处理错误
	var isMaster bool
	if _, ok := options["acl"]; ok {
		isMaster = false
	} else {
		isMaster = true
	}
	var aclGroup []string
	if options["acl"] == nil {
		aclGroup = []string{}
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
		writePerms := types.S{}
		perm := types.M{
			"_wperm": types.M{"$exists": false},
		}
		writePerms = append(writePerms, perm)
		for _, acl := range aclGroup {
			perm = types.M{
				"_wperm": types.M{"$in": []string{acl}},
			}
			writePerms = append(writePerms, perm)
		}

		mongoWhere = types.M{
			"$and": types.S{
				mongoWhere,
				types.M{"$or": writePerms},
			},
		}
	}
	coll.deleteMany(mongoWhere)
	// TODO 处理返回错误
}

// Update ...
func Update(className string, where, data, options types.M) (types.M, error) {
	// TODO 处理错误
	data = utils.CopyMap(data)
	acceptor := func(schema *Schema) bool {
		keys := []string{}
		for k := range where {
			keys = append(keys, k)
		}
		return schema.hasKeys(className, keys)
	}
	var isMaster bool
	if _, ok := options["acl"]; ok {
		isMaster = false
	} else {
		isMaster = true
	}
	var aclGroup []string
	if options["acl"] == nil {
		aclGroup = []string{}
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
		writePerms := types.S{}
		perm := types.M{
			"_wperm": types.M{"$exists": false},
		}
		writePerms = append(writePerms, perm)
		for _, acl := range aclGroup {
			perm = types.M{
				"_wperm": types.M{"$in": []string{acl}},
			}
			writePerms = append(writePerms, perm)
		}

		mongoWhere = types.M{
			"$and": types.S{
				mongoWhere,
				types.M{"$or": writePerms},
			},
		}
	}
	mongoUpdate := transformUpdate(schema, className, data)

	result := coll.FindOneAndUpdate(mongoWhere, mongoUpdate)
	// TODO 处理返回错误

	response := types.M{}
	if mongoUpdate["$inc"] != nil && utils.MapInterface(mongoUpdate["$inc"]) != nil {
		inc := utils.MapInterface(mongoUpdate["$inc"])
		for k := range inc {
			response[k] = result[k]
		}
	}

	return response, nil
}

// Create ...
func Create(className string, data, options types.M) error {
	// TODO 处理错误
	data = utils.CopyMap(data)
	var isMaster bool
	if _, ok := options["acl"]; ok {
		isMaster = false
	} else {
		isMaster = true
	}
	var aclGroup []string
	if options["acl"] == nil {
		aclGroup = []string{}
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

func handleRelationUpdates(className, objectID string, update types.M) {
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
	doc := types.M{
		"relatedId": toID,
		"owningId":  fromID,
	}
	className := "_Join:" + key + ":" + fromClassName
	coll := AdaptiveCollection(className)
	coll.upsertOne(doc, doc)
}

func removeRelation(key, fromClassName, fromID, toID string) {
	// TODO 处理错误
	doc := types.M{
		"relatedId": toID,
		"owningId":  fromID,
	}
	className := "_Join:" + key + ":" + fromClassName
	coll := AdaptiveCollection(className)
	coll.deleteOne(doc)
}

// ValidateObject ...
func ValidateObject(className string, object, where, options types.M) error {
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

// RedirectClassNameForKey 返回指定类的字段所对应的类型
func RedirectClassNameForKey(className, key string) string {
	schema := LoadSchema(nil)
	t := schema.getExpectedType(className, key)
	b, _ := regexp.MatchString(`^relation<(.*)>$`, t)
	if b {
		return className[len("relation<"):(len(className) - 1)]
	}
	return className
}

// canAddField ...
func canAddField(schema *Schema, className string, object types.M, acl []string) {
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

func keysForQuery(query types.M) []string {
	answer := []string{}

	var s interface{}
	if query["$and"] != nil {
		s = query["$and"]
	} else {
		s = query["$or"]
	}

	if s != nil {
		sublist := utils.SliceInterface(s)
		for _, v := range sublist {
			subquery := utils.MapInterface(v)
			answer = append(answer, keysForQuery(subquery)...)
		}
		return answer
	}

	for k := range query {
		answer = append(answer, k)
	}

	return answer
}

func reduceRelationKeys(className string, query types.M) {
	if query["$or"] != nil {
		subQuerys := utils.SliceInterface(query["$or"])
		for _, v := range subQuerys {
			aQuery := utils.MapInterface(v)
			reduceRelationKeys(className, aQuery)
		}
		return
	}

	if query["$relatedTo"] != nil {
		relatedTo := utils.MapInterface(query["$relatedTo"])
		key := utils.String(relatedTo["key"])
		object := utils.MapInterface(relatedTo["object"])
		objClassName := utils.String(object["className"])
		objID := utils.String(object["objectId"])
		ids := relatedIds(objClassName, key, objID)
		delete(query, "$relatedTo")
		addInObjectIdsIds(ids, query)
		reduceRelationKeys(className, query)
	}

}

func relatedIds(className, key, owningID string) types.S {
	coll := AdaptiveCollection(joinTableName(className, key))
	results := coll.Find(types.M{"owningId": owningID}, types.M{})
	ids := types.S{}
	for _, r := range results {
		id := r["relatedId"]
		ids = append(ids, id)
	}
	return ids
}

func joinTableName(className, key string) string {
	return "_Join:" + key + ":" + className
}

func addInObjectIdsIds(ids types.S, query types.M) {
	if id, ok := query["objectId"].(string); ok {
		query["objectId"] = types.M{"$eq": id}
	}

	objectID := utils.MapInterface(query["objectId"])
	if objectID == nil {
		objectID = types.M{}
	}

	queryIn := types.S{}
	if objectID["$in"] != nil && utils.SliceInterface(objectID["$in"]) != nil {
		in := utils.SliceInterface(objectID["$in"])
		queryIn = append(queryIn, in...)
	}
	if ids != nil {
		queryIn = append(queryIn, ids...)
	}
	objectID["$in"] = queryIn
	query["objectId"] = objectID
}

func reduceInRelation(className string, query types.M, schema *Schema) types.M {
	if query["$or"] != nil {
		ors := utils.SliceInterface(query["$or"])
		for i, v := range ors {
			aQuery := utils.MapInterface(v)
			aQuery = reduceInRelation(className, aQuery, schema)
			ors[i] = aQuery
		}
		query["$or"] = ors
		return query
	}

	for key, v := range query {
		op := utils.MapInterface(v)
		if v != nil && (op["$in"] != nil || utils.String(op["__type"]) == "Pointer") {
			// 只处理 relation 类型
			t := schema.getExpectedType(className, key)
			match := false
			if t != "" {
				b, _ := regexp.MatchString("^relation<(.*)>$", t)
				match = b
			}
			if match == false {
				return query
			}

			relatedIds := types.S{}
			if op["$in"] != nil {
				ors := utils.SliceInterface(op["$in"])
				for _, v := range ors {
					r := utils.MapInterface(v)
					relatedIds = append(relatedIds, r["objectId"])
				}
			} else {
				relatedIds = append(relatedIds, op["objectId"])
			}

			ids := owningIds(className, key, relatedIds)
			delete(query, key)
			addInObjectIdsIds(ids, query)
		}
	}

	return query
}

func owningIds(className, key string, relatedIds types.S) types.S {
	coll := AdaptiveCollection(joinTableName(className, key))
	query := types.M{
		"relatedId": types.M{
			"$in": relatedIds,
		},
	}
	results := coll.Find(query, types.M{})
	ids := types.S{}
	for _, r := range results {
		ids = append(ids, r["owningId"])
	}
	return ids
}

func untransformObject(schema *Schema, isMaster bool, aclGroup []string, className string, mongoObject types.M) types.M {
	res := untransformObjectT(schema, className, mongoObject, false)
	object := utils.MapInterface(res)
	if className != "_User" {
		return object
	}
	// 以下单独处理 _User 类
	if isMaster {
		return object
	}
	// 当前用户返回所有信息
	id := utils.String(object["objectId"])
	for _, v := range aclGroup {
		if v == id {
			return object
		}
	}
	// 其他用户删除相关信息
	delete(object, "authData")
	delete(object, "sessionToken")
	return object
}
