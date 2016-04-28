//Package orm 数据库操作模块，当前只对接了 MongoDB
package orm

import (
	"regexp"
	"strings"

	"github.com/lfq7413/tomato/errs"
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

// Destroy 从指定表中删除数据
func Destroy(className string, where types.M, options types.M) error {
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
		err := schema.validatePermission(className, aclGroup, "delete")
		return err
	}

	coll := AdaptiveCollection(className)
	mongoWhere := transformWhere(schema, className, where)
	// 组装 acl 查询条件，查找可被当前用户修改的对象
	if isMaster == false {
		writePerms := types.S{}
		// 可修改 不存在写权限字段的
		perm := types.M{
			"_wperm": types.M{"$exists": false},
		}
		writePerms = append(writePerms, perm)
		for _, acl := range aclGroup {
			// 可修改 写权限包含 当前用户角色与 id 的
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
	n, err := coll.deleteMany(mongoWhere)
	if err != nil {
		return err
	}
	// 排除 _Session，避免在修改密码时因为没有 Session 失败
	if n == 0 && className != "_Session" {
		return errs.E(errs.ObjectNotFound, "Object not found.")
	}

	return nil
}

// Update 更新对象
func Update(className string, where, data, options types.M) (types.M, error) {
	// 复制数据，不要修改原数据
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
		err := schema.validatePermission(className, aclGroup, "update")
		if err != nil {
			return nil, err
		}
	}
	// 处理 Relation
	handleRelationUpdates(className, utils.String(where["objectId"]), data)

	coll := AdaptiveCollection(className)
	mongoWhere := transformWhere(schema, className, where)
	// 组装 acl 查询条件，查找可被当前用户修改的对象
	if isMaster == false {
		writePerms := types.S{}
		// 可修改 不存在写权限字段的
		perm := types.M{
			"_wperm": types.M{"$exists": false},
		}
		writePerms = append(writePerms, perm)
		for _, acl := range aclGroup {
			// 可修改 写权限包含 当前用户角色与 id 的
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
	if result == nil || len(result) == 0 {
		return nil, errs.E(errs.ObjectNotFound, "Object not found.")
	}

	// 返回 数值增加的字段
	response := types.M{}
	if mongoUpdate["$inc"] != nil && utils.MapInterface(mongoUpdate["$inc"]) != nil {
		inc := utils.MapInterface(mongoUpdate["$inc"])
		for k := range inc {
			response[k] = result[k]
		}
	}

	return response, nil
}

// Create 创建对象
func Create(className string, data, options types.M) error {
	// 不要对原数据进行修改
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

	err := validateClassName(className)
	if err != nil {
		return err
	}

	schema := LoadSchema(nil)
	if isMaster == false {
		err := schema.validatePermission(className, aclGroup, "create")
		if err != nil {
			return err
		}
	}

	// 处理 Relation
	err = handleRelationUpdates(className, "", data)
	if err != nil {
		return err
	}

	coll := AdaptiveCollection(className)
	mongoObject := transformCreate(schema, className, data)
	return coll.insertOne(mongoObject)
}

// validateClassName 校验表名是否合法
func validateClassName(className string) error {
	if ClassNameIsValid(className) == false {
		return errs.E(errs.InvalidClassName, "invalid className: "+className)
	}
	return nil
}

// handleRelationUpdates 处理 Relation 相关操作
func handleRelationUpdates(className, objectID string, update types.M) error {
	objID := objectID
	if utils.String(update["objectId"]) != "" {
		objID = utils.String(update["objectId"])
	}

	// 定义处理函数
	// 传入参数 op 的格式如下
	// {
	//       "__op": "AddRelation",
	//       "objects": [
	//         {
	//           "__type": "Pointer",
	//           "className": "_User",
	//           "objectId": "8TOXdXf3tz"
	//         },
	//         {
	//           "__type": "Pointer",
	//           "className": "_User",
	//           "objectId": "g7y9tkhB7O"
	//         }
	//       ]
	// }
	var process func(op interface{}, key string) error
	process = func(op interface{}, key string) error {
		if op == nil || utils.MapInterface(op) == nil || utils.MapInterface(op)["__op"] == nil {
			return nil
		}
		opMap := utils.MapInterface(op)
		p := utils.String(opMap["__op"])
		if p == "AddRelation" {
			delete(update, key)
			// 添加 Relation 对象
			objects := utils.SliceInterface(opMap["objects"])
			for _, object := range objects {
				relationID := utils.String(utils.MapInterface(object)["objectId"])
				err := addRelation(key, className, objID, relationID)
				if err != nil {
					return err
				}
			}
		} else if p == "RemoveRelation" {
			delete(update, key)
			// 删除 Relation 对象
			objects := utils.SliceInterface(opMap["objects"])
			for _, object := range objects {
				relationID := utils.String(utils.MapInterface(object)["objectId"])
				err := removeRelation(key, className, objID, relationID)
				if err != nil {
					return err
				}
			}
		} else if p == "Batch" {
			// 批处理 Relation 对象
			ops := utils.SliceInterface(opMap["ops"])
			for _, x := range ops {
				err := process(x, key)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}

	for k, v := range update {
		err := process(v, k)
		if err != nil {
			return err
		}
	}
	return nil
}

// addRelation 把对象 id 加入 _Join 表，表名为 _Join:key:fromClassName
func addRelation(key, fromClassName, fromID, toID string) error {
	doc := types.M{
		"relatedId": toID,
		"owningId":  fromID,
	}
	className := "_Join:" + key + ":" + fromClassName
	coll := AdaptiveCollection(className)
	return coll.upsertOne(doc, doc)
}

// removeRelation 把对象 id 从 _Join 表中删除，表名为 _Join:key:fromClassName
func removeRelation(key, fromClassName, fromID, toID string) error {
	doc := types.M{
		"relatedId": toID,
		"owningId":  fromID,
	}
	className := "_Join:" + key + ":" + fromClassName
	coll := AdaptiveCollection(className)
	return coll.deleteOne(doc)
}

// ValidateObject 校验对象是否合法
func ValidateObject(className string, object, where, options types.M) error {
	schema := LoadSchema(nil)
	acl := []string{}
	if options["acl"] != nil {
		if v, ok := options["acl"].([]string); ok {
			acl = v
		}
	}

	err := canAddField(schema, className, object, acl)
	if err != nil {
		return err
	}

	err = schema.validateObject(className, object, where)
	if err != nil {
		return err
	}

	return nil
}

// LoadSchema 加载 Schema，仅加载一次，当 acceptor 返回 false 时，再从数据库读取一次
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
// 如果 key 字段的属性为 relation<classA> ，则返回 classA
func RedirectClassNameForKey(className, key string) string {
	schema := LoadSchema(nil)
	t := schema.getExpectedType(className, key)
	b, _ := regexp.MatchString(`^relation<(.*)>$`, t)
	if b {
		return className[len("relation<"):(len(className) - 1)]
	}
	return className
}

// canAddField 检测是否能添加字段到类上
func canAddField(schema *Schema, className string, object types.M, acl []string) error {
	if schema.data[className] == nil {
		return nil
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
		return schema.validatePermission(className, acl, "addField")
	}

	return nil
}

// keysForQuery 从查询条件中查找字段名
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

// reduceRelationKeys 处理查询条件中的 $relatedTo
// query 格式如下
// {
//     "$relatedTo":{
//         "object":{
//             "__type":"Pointer",
//             "className":"Post",
//             "objectId":"8TOXdXf3tz"
//         },
//         "key":"likes"
//     }
// }
// 表 Post 中的字段 likes 的类型为 relation<classA>
// 从 _Join:likes:Post 表中查询 Post id 对应的 classA id 列表，并添加到 query 中
// 替换后格式为
// {
//     "objectId":{
//         "$in":[
//             "id",
//             "id2"
//         ]
//     }
// }
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

// relatedIds 从 Join 表中查询 ids ，表名：_Join:key:className
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

// joinTableName 组装用于 relation 的 Join 表
func joinTableName(className, key string) string {
	return "_Join:" + key + ":" + className
}

// addInObjectIdsIds 添加 ids 到查询条件中
// 替换 $relatedTo 为：
// "objectId":{"$eq":"id"}
// 或者
// "objectId":{"$in":["id","id2"]}
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

// reduceInRelation 处理查询条件中，作用于 relation 类型字段上的 $in 或者等于某对象
// 例如 classA 中的 字段 key 为 relation<classB> 类型，查找 key 中包含指定 classB 对象的 classA
// query = {"key":{"$in":[]}}
func reduceInRelation(className string, query types.M, schema *Schema) types.M {
	// 处理 $or 数组中的数据，并替换回去
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

			// 从 Join 表中查找的 ids，替换查询条件
			ids := owningIds(className, key, relatedIds)
			delete(query, key)
			addInObjectIdsIds(ids, query)
		}
	}

	return query
}

// owningIds 从 Join 表中查询 relatedIds 对应的父对象
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

// untransformObject 从查询到的数据库对象转换出可返回给客户端的对象，并对 _User 表数据进行特殊处理
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
