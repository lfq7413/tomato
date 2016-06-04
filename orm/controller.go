//Package orm 数据库操作模块，当前只对接了 MongoDB
package orm

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/storage"
	"github.com/lfq7413/tomato/storage/mongo"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// TomatoDBController ...
var TomatoDBController *DBController

// Adapter ...
var Adapter storage.Adapter

// Transform ...
var Transform storage.Transform
var schemaPromise *Schema

// init 初始化 Mongo 适配器
func init() {
	Adapter = mongo.NewMongoAdapter("tomato")
	Transform = Adapter.GetTransform()
	TomatoDBController = &DBController{
		skipValidation: false,
	}
}

// DBController 数据库操作类
type DBController struct {
	skipValidation bool
}

// WithoutValidation 返回不进行字段校验的数据库操作对象
func (d DBController) WithoutValidation() *DBController {
	return &DBController{
		skipValidation: true,
	}
}

// SchemaCollection 获取 Schema 表
func (d DBController) SchemaCollection() storage.SchemaCollection {
	return Adapter.SchemaCollection()
}

// CollectionExists 检测表是否存在
func (d DBController) CollectionExists(className string) bool {
	return Adapter.CollectionExists(className)
}

// Find 从指定表中查询数据，查询到的数据放入 list 中
// 如果查询的是 count ，结果也会放入 list，并且只有这一个元素
// options 中的选项包括：skip、limit、sort、count、acl
func (d DBController) Find(className string, where, options types.M) (types.S, error) {
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

	isMaster := false
	aclGroup := []string{}
	if acl, ok := options["acl"]; ok {
		if v, ok := acl.([]string); ok {
			aclGroup = v
		}
	} else {
		isMaster = true
	}

	var op string
	if _, ok := where["objectId"].(string); ok {
		if len(where) == 1 {
			op = "get"
		} else {
			op = "find"
		}
	} else {
		op = "find"
	}

	schema := d.LoadSchema()
	parseFormatSchema, err := schema.GetOneSchema(className, false)
	if err != nil {
		return nil, err
	}
	if len(parseFormatSchema) == 0 {
		parseFormatSchema["fields"] = types.M{}
	}

	if options["sort"] != nil {
		sortKeys := []string{}
		keys := options["sort"].([]string)
		for _, key := range keys {
			mongoKey := ""
			// sort 中的 key ，如果是要按倒序排列，则会加前缀 "-" ，所以要对其进行处理
			var prefix string
			if strings.HasPrefix(key, "-") {
				prefix = "-"
				key = key[1:]
			}

			if key == "_created_at" {
				key = "createdAt"
			} else if key == "_updated_at" {
				key = "updatedAt"
			}

			if fieldNameIsValid(key) == false {
				return nil, errs.E(errs.InvalidKeyName, "Invalid field name: "+key)
			}

			if match, _ := regexp.MatchString(`^authData\.([a-zA-Z0-9_]+)\.id$`, key); match {
				return nil, errs.E(errs.InvalidKeyName, "Cannot sort by "+key)
			}

			k := Transform.TransformKey(className, key, parseFormatSchema)
			mongoKey = prefix + k

			sortKeys = append(sortKeys, mongoKey)
		}
		mongoOptions["sort"] = sortKeys
	}

	// 校验当前用户是否能对表进行 find 或者 get 操作
	if isMaster == false {
		err := schema.validatePermission(className, aclGroup, op)
		if err != nil {
			return nil, err
		}
	}

	// 处理 $relatedTo
	d.reduceRelationKeys(className, where)
	// 处理 relation 字段上的 $in
	d.reduceInRelation(className, where, schema)

	coll := Adapter.AdaptiveCollection(className)

	if isMaster == false {
		where = d.addPointerPermissions(schema, className, op, where, aclGroup)
	}
	if where == nil {
		if op == "get" {
			return nil, errs.E(errs.ObjectNotFound, "Object not found.")
		}
		return types.S{}, nil
	}

	// 组装 acl 查询条件，查找可被当前用户访问的对象
	if isMaster == false {
		where = addReadACL(where, aclGroup)
	}

	mongoWhere, err := Transform.TransformWhere(className, where, nil, parseFormatSchema)
	if err != nil {
		return nil, err
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
		result, err := d.untransformObject(schema, isMaster, aclGroup, className, r)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil

}

// Destroy 从指定表中删除数据
func (d DBController) Destroy(className string, where types.M, options types.M) error {
	isMaster := false
	aclGroup := []string{}
	if acl, ok := options["acl"]; ok {
		if v, ok := acl.([]string); ok {
			aclGroup = v
		}
	} else {
		isMaster = true
	}

	schema := d.LoadSchema()
	if isMaster == false {
		err := schema.validatePermission(className, aclGroup, "delete")
		return err
	}

	if isMaster == false {
		where = d.addPointerPermissions(schema, className, "delete", where, aclGroup)
		if where == nil {
			return errs.E(errs.ObjectNotFound, "Object not found.")
		}
	}

	if isMaster == false {
		where = addWriteACL(where, aclGroup)
	}

	parseFormatSchema, err := schema.GetOneSchema(className, false)
	if err != nil {
		return err
	}
	if len(parseFormatSchema) == 0 {
		parseFormatSchema["fields"] = types.M{}
	}

	err = Adapter.DeleteObjectsByQuery(className, where, !d.skipValidation, parseFormatSchema)
	if err != nil {
		msg := err.Error()
		msg = strings.Replace(msg, " ", "", -1)
		// 排除 _Session，避免在修改密码时因为没有 Session 失败
		if className == "_Session" && strings.Index(msg, `"code":`+strconv.Itoa(errs.ObjectNotFound)) > -1 {
			return nil
		}
		return err
	}

	return nil
}

// Update 更新对象
// options 中的参数包括：acl、many、upsert
func (d DBController) Update(className string, where, data, options types.M) (types.M, error) {
	if options == nil {
		options = types.M{}
	}
	originalUpdate := data
	// 复制数据，不要修改原数据
	data = utils.CopyMap(data)

	many := options["many"].(bool)
	upsert := options["upsert"].(bool)

	isMaster := false
	aclGroup := []string{}
	if acl, ok := options["acl"]; ok {
		if v, ok := acl.([]string); ok {
			aclGroup = v
		}
	} else {
		isMaster = true
	}

	schema := d.LoadSchema()
	if isMaster == false {
		err := schema.validatePermission(className, aclGroup, "update")
		if err != nil {
			return nil, err
		}
	}
	// 处理 Relation
	d.handleRelationUpdates(className, utils.String(where["objectId"]), data)

	coll := Adapter.AdaptiveCollection(className)

	// 添加用户权限
	if isMaster == false {
		where = d.addPointerPermissions(schema, className, "update", where, aclGroup)
	}
	if where == nil {
		return types.M{}, nil
	}

	// 组装 acl 查询条件，查找可被当前用户修改的对象
	if isMaster == false {
		where = addWriteACL(where, aclGroup)
	}

	parseFormatSchema, err := schema.GetOneSchema(className, false)
	if err != nil {
		return nil, err
	}
	if len(parseFormatSchema) == 0 {
		parseFormatSchema["fields"] = types.M{}
	}

	mongoWhere, err := Transform.TransformWhere(className, where, types.M{"validate": !d.skipValidation}, parseFormatSchema)
	if err != nil {
		return nil, err
	}
	mongoUpdate, err := Transform.TransformUpdate(schema, className, data, types.M{"validate": !d.skipValidation})
	if err != nil {
		return nil, err
	}
	var result types.M
	if many {
		err := coll.UpdateMany(mongoWhere, mongoUpdate)
		if err != nil {
			return nil, err
		}
		result = types.M{}
	} else if upsert {
		err := coll.UpsertOne(mongoWhere, mongoUpdate)
		if err != nil {
			return nil, err
		}
		result = types.M{}
	} else {
		result = coll.FindOneAndUpdate(mongoWhere, mongoUpdate)
	}

	if result == nil {
		return nil, errs.E(errs.ObjectNotFound, "Object not found.")
	}

	if d.skipValidation {
		return result, nil
	}

	// 返回经过修改的字段
	response := sanitizeDatabaseResult(originalUpdate, result)

	return response, nil
}

// sanitizeDatabaseResult 处理数据库返回结果
func sanitizeDatabaseResult(originalObject, result types.M) types.M {
	response := types.M{}
	if result == nil {
		return response
	}

	// 检测是否是对字段的操作
	for key, value := range originalObject {
		if value != nil && utils.MapInterface(value) != nil {
			keyUpdate := utils.MapInterface(value)
			if keyUpdate["__op"] != nil {
				op := utils.String(keyUpdate["__op"])
				if op == "Add" || op == "AddUnique" || op == "Remove" || op == "Increment" {
					// 只把操作的字段放入返回结果中
					response[key] = result[key]
				}
			}
		}
	}

	return response
}

// Create 创建对象
func (d DBController) Create(className string, data, options types.M) error {
	if options == nil {
		options = types.M{}
	}
	// 不要对原数据进行修改
	data = utils.CopyMap(data)

	isMaster := false
	aclGroup := []string{}
	if acl, ok := options["acl"]; ok {
		if v, ok := acl.([]string); ok {
			aclGroup = v
		}
	} else {
		isMaster = true
	}

	err := d.validateClassName(className)
	if err != nil {
		return err
	}

	schema := d.LoadSchema()
	if isMaster == false {
		err := schema.validatePermission(className, aclGroup, "create")
		if err != nil {
			return err
		}
	}

	// 处理 Relation
	err = d.handleRelationUpdates(className, "", data)
	if err != nil {
		return err
	}

	err = schema.enforceClassExists(className, false)
	if err != nil {
		return err
	}

	sch, err := schema.GetOneSchema(className, true)
	if err != nil {
		return err
	}

	// 无需调用 sanitizeDatabaseResult
	return Adapter.CreateObject(className, data, schema, sch)
}

// validateClassName 校验表名是否合法
func (d DBController) validateClassName(className string) error {
	if d.skipValidation {
		return nil
	}
	if ClassNameIsValid(className) == false {
		return errs.E(errs.InvalidClassName, "invalid className: "+className)
	}
	return nil
}

// handleRelationUpdates 处理 Relation 相关操作
func (d DBController) handleRelationUpdates(className, objectID string, update types.M) error {
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
				err := d.addRelation(key, className, objID, relationID)
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
				err := d.removeRelation(key, className, objID, relationID)
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
func (d DBController) addRelation(key, fromClassName, fromID, toID string) error {
	doc := types.M{
		"relatedId": toID,
		"owningId":  fromID,
	}
	className := "_Join:" + key + ":" + fromClassName
	coll := Adapter.AdaptiveCollection(className)
	return coll.UpsertOne(doc, doc)
}

// removeRelation 把对象 id 从 _Join 表中删除，表名为 _Join:key:fromClassName
func (d DBController) removeRelation(key, fromClassName, fromID, toID string) error {
	doc := types.M{
		"relatedId": toID,
		"owningId":  fromID,
	}
	className := "_Join:" + key + ":" + fromClassName
	coll := Adapter.AdaptiveCollection(className)
	return coll.DeleteOne(doc)
}

// ValidateObject 校验对象是否合法
func (d DBController) ValidateObject(className string, object, where, options types.M) error {
	schema := d.LoadSchema()

	isMaster := false
	aclGroup := []string{}
	if acl, ok := options["acl"]; ok {
		if v, ok := acl.([]string); ok {
			aclGroup = v
		}
	} else {
		isMaster = true
	}

	if isMaster {
		return nil
	}

	err := d.canAddField(schema, className, object, aclGroup)
	if err != nil {
		return err
	}

	err = schema.validateObject(className, object, where)
	if err != nil {
		return err
	}

	return nil
}

// LoadSchema 加载 Schema，仅加载一次
func (d DBController) LoadSchema() *Schema {
	if schemaPromise == nil {
		collection := d.SchemaCollection()
		schemaPromise = Load(collection, Adapter)
	}
	return schemaPromise
}

// MongoFind 直接执行数据库查询，仅用于测试
func (d *DBController) MongoFind(className string, query, options types.M) []types.M {
	coll := Adapter.AdaptiveCollection(className)
	return coll.Find(query, options)
}

// DeleteEverything 删除所有表数据，仅用于测试
func (d DBController) DeleteEverything() {
	schemaPromise = nil
	Adapter.DeleteAllSchemas()
}

// RedirectClassNameForKey 返回指定类的字段所对应的类型
// 如果 key 字段的属性为 relation<classA> ，则返回 classA
func (d DBController) RedirectClassNameForKey(className, key string) string {
	schema := d.LoadSchema()
	t := schema.GetExpectedType(className, key)
	if t != nil && t["type"].(string) == "Relation" {
		return t["targetClass"].(string)
	}
	return className
}

// canAddField 检测是否能添加字段到类上
func (d DBController) canAddField(schema *Schema, className string, object types.M, acl []string) error {
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
func (d DBController) reduceRelationKeys(className string, query types.M) {
	if query["$or"] != nil {
		subQuerys := utils.SliceInterface(query["$or"])
		for _, v := range subQuerys {
			aQuery := utils.MapInterface(v)
			d.reduceRelationKeys(className, aQuery)
		}
		return
	}

	if query["$relatedTo"] != nil {
		relatedTo := utils.MapInterface(query["$relatedTo"])
		key := utils.String(relatedTo["key"])
		object := utils.MapInterface(relatedTo["object"])
		objClassName := utils.String(object["className"])
		objID := utils.String(object["objectId"])
		ids := d.relatedIds(objClassName, key, objID)
		delete(query, "$relatedTo")
		d.addInObjectIdsIds(ids, query)
		d.reduceRelationKeys(className, query)
	}

}

// relatedIds 从 Join 表中查询 ids ，表名：_Join:key:className
func (d DBController) relatedIds(className, key, owningID string) types.S {
	coll := Adapter.AdaptiveCollection(joinTableName(className, key))
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

// addInObjectIdsIds 添加 ids 到查询条件中, 应该取 objectId $eq $in ids 的交集
// 替换 objectId 为：
// "objectId":{"$in":["id","id2"]}
func (d DBController) addInObjectIdsIds(ids types.S, query types.M) {
	coll := map[string]types.S{}
	idsFromString := types.S{}
	if id, ok := query["objectId"].(string); ok {
		idsFromString = append(idsFromString, id)
	}
	coll["idsFromString"] = idsFromString

	idsFromEq := types.S{}
	if eqid, ok := query["objectId"].(map[string]interface{}); ok {
		if id, ok := eqid["$eq"]; ok {
			idsFromEq = append(idsFromEq, id.(string))
		}
	}
	coll["idsFromEq"] = idsFromEq

	idsFromIn := types.S{}
	if inid, ok := query["objectId"].(map[string]interface{}); ok {
		if id, ok := inid["$in"]; ok {
			idsFromIn = append(idsFromIn, id.([]interface{})...)
		}
	}
	coll["idsFromIn"] = idsFromIn

	if ids != nil {
		coll["ids"] = ids
	}

	// 统计 idsFromString idsFromEq idsFromIn ids 中的共同元素加入到 $in 中
	max := 0 // 以上4个集合中不为空的个数，也就是说 某个 objectId 出现的次数应该等于 max 才能加入到 $in 中查询
	for k, v := range coll {
		// 删除空集合
		if len(v) > 0 {
			max++
		} else {
			delete(coll, k)
		}
	}
	idsColl := map[string]int{} // 统计每个 objectId 出现的次数
	for _, c := range coll {
		// 从每个集合中取出 objectId
		idColl := map[string]int{}
		for _, v := range c {
			id := v.(string)
			// 并去除重复
			if _, ok := idColl[id]; ok == false {
				idColl[id] = 0

				// 加入到 idsColl 中，并增加出现次数
				if i, ok := idsColl[id]; ok {
					idsColl[id] = i + 1
				} else {
					idsColl[id] = 1
				}
			}
		}
	}
	queryIn := types.S{} // 统计出现次数为 max 的 objectId
	for k, v := range idsColl {
		if v == max {
			queryIn = append(queryIn, k)
		}
	}

	if v, ok := query["objectId"]; ok {
		if _, ok := v.(string); ok {
			query["objectId"] = types.M{}
		}
	} else {
		query["objectId"] = types.M{}
	}
	id := query["objectId"].(map[string]interface{})
	id["$in"] = queryIn

	query["objectId"] = id
}

// addNotInObjectIdsIds 添加 ids 到查询条件中，应该取 $ne $nin ids 的并集
// 替换 objectId 为：
// "objectId":{"$nin":["id","id2"]}
func (d DBController) addNotInObjectIdsIds(ids types.S, query types.M) {
	coll := map[string]types.S{}
	idsFromNin := types.S{}
	if ninid, ok := query["objectId"].(map[string]interface{}); ok {
		if id, ok := ninid["$nin"]; ok {
			idsFromNin = append(idsFromNin, id.([]interface{})...)
		}
	}
	coll["idsFromNin"] = idsFromNin

	if ids != nil {
		coll["ids"] = ids
	}

	idsColl := map[string]int{}
	for _, c := range coll {
		// 从每个集合中取出 objectId
		for _, v := range c {
			id := v.(string)
			// 并去除重复
			if _, ok := idsColl[id]; ok == false {
				idsColl[id] = 0
			}
		}
	}

	queryNin := types.S{}
	for k := range idsColl {
		queryNin = append(queryNin, k)
	}

	if v, ok := query["objectId"]; ok {
		if _, ok := v.(string); ok {
			query["objectId"] = types.M{}
		}
	} else {
		query["objectId"] = types.M{}
	}
	id := query["objectId"].(map[string]interface{})
	id["$nin"] = queryNin

	query["objectId"] = id
}

// reduceInRelation 处理查询条件中，作用于 relation 类型字段上的 $in $ne $nin $eq 或者等于某对象
// 例如 classA 中的 字段 key 为 relation<classB> 类型，查找 key 中包含指定 classB 对象的 classA
// query = {"key":{"$in":[]}}
func (d DBController) reduceInRelation(className string, query types.M, schema *Schema) types.M {
	// 处理 $or 数组中的数据，并替换回去
	if query["$or"] != nil {
		ors := utils.SliceInterface(query["$or"])
		for i, v := range ors {
			aQuery := utils.MapInterface(v)
			aQuery = d.reduceInRelation(className, aQuery, schema)
			ors[i] = aQuery
		}
		query["$or"] = ors
		return query
	}

	for key, v := range query {
		op := utils.MapInterface(v)
		if op != nil && (op["$in"] != nil || op["$ne"] != nil || op["$nin"] != nil || op["$eq"] != nil || utils.String(op["__type"]) == "Pointer") {
			// 只处理 relation 类型
			t := schema.GetExpectedType(className, key)
			if t == nil || t["type"].(string) != "Relation" {
				return query
			}

			// 取出所有限制条件
			relatedIds := []types.S{}
			isNegation := []bool{}
			for constraintKey, value := range op {
				if constraintKey == "objectId" {
					ids := types.S{value}
					relatedIds = append(relatedIds, ids)
					isNegation = append(isNegation, false)
				} else if constraintKey == "$in" {
					in := utils.SliceInterface(value)
					ids := types.S{}
					for _, v := range in {
						r := utils.MapInterface(v)
						ids = append(ids, r["objectId"])
					}
					relatedIds = append(relatedIds, ids)
					isNegation = append(isNegation, false)
				} else if constraintKey == "$nin" {
					nin := utils.SliceInterface(value)
					ids := types.S{}
					for _, v := range nin {
						r := utils.MapInterface(v)
						ids = append(ids, r["objectId"])
					}
					relatedIds = append(relatedIds, ids)
					isNegation = append(isNegation, true)
				} else if constraintKey == "$ne" {
					ne := utils.MapInterface(value)
					ids := types.S{ne["objectId"]}
					relatedIds = append(relatedIds, ids)
					isNegation = append(isNegation, true)
				} else if constraintKey == "$eq" {
					eq := utils.MapInterface(value)
					ids := types.S{eq["objectId"]}
					relatedIds = append(relatedIds, ids)
					isNegation = append(isNegation, false)
				}
			}

			delete(query, key)

			// 应用所有限制条件
			for i, relatedID := range relatedIds {
				// 从 Join 表中查找的 ids，替换查询条件
				ids := d.owningIds(className, key, relatedID)
				if isNegation[i] {
					d.addNotInObjectIdsIds(ids, query)
				} else {
					d.addInObjectIdsIds(ids, query)
				}
			}
		}
	}

	return query
}

// owningIds 从 Join 表中查询 relatedIds 对应的父对象
func (d DBController) owningIds(className, key string, relatedIds types.S) types.S {
	coll := Adapter.AdaptiveCollection(joinTableName(className, key))
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
func (d DBController) untransformObject(schema *Schema, isMaster bool, aclGroup []string, className string, mongoObject types.M) (types.M, error) {
	res, err := Transform.UntransformObject(schema, className, mongoObject, false)
	if err != nil {
		return nil, err
	}
	object := utils.MapInterface(res)
	if className != "_User" {
		return object, nil
	}
	// 以下单独处理 _User 类
	delete(object, "sessionToken")
	if isMaster {
		return object, nil
	}
	// 当前用户返回所有信息
	id := utils.String(object["objectId"])
	for _, v := range aclGroup {
		if v == id {
			return object, nil
		}
	}
	delete(object, "authData")
	return object, nil
}

// DeleteSchema 删除类
func (d *DBController) DeleteSchema(className string) error {
	exist := d.CollectionExists(className)
	if exist == false {
		return nil
	}
	coll := Adapter.AdaptiveCollection(className)
	count := coll.Count(types.M{}, types.M{})
	if count > 0 {
		return errs.E(errs.ClassNotEmpty, "Class "+className+" is not empty, contains "+strconv.Itoa(count)+" objects, cannot drop schema.")
	}
	return Adapter.DeleteOneSchema(className)
}

// addPointerPermissions 添加查询用户权限，perms[className][readUserFields] 中保存的是字段名，该字段中的内容是：有权限进行读操作的用户
func (d *DBController) addPointerPermissions(schema *Schema, className string, operation string, query types.M, aclGroup []string) types.M {
	perms := schema.perms[className]
	var field string
	if operation == "get" || operation == "find" {
		field = "readUserFields"
	} else {
		field = "writeUserFields"
	}
	userACL := []string{}
	for _, acl := range aclGroup {
		if strings.HasPrefix(acl, "role:") == false && acl != "*" {
			userACL = append(userACL, acl)
		}
	}

	if perms != nil && utils.MapInterface(perms)[field] != nil {
		permFields := utils.MapInterface(perms)[field].([]interface{})
		if permFields != nil && len(permFields) > 0 {
			if len(userACL) != 1 {
				return nil
			}
			userID := userACL[0]
			userPointer := types.M{
				"__type":    "Pointer",
				"className": "_User",
				"objectId":  userID,
			}

			ors := []types.M{}
			for _, key := range permFields {
				q := types.M{
					key.(string): userPointer,
				}
				and := types.M{
					"$and": types.S{q, query},
				}
				ors = append(ors, and)
			}
			if len(ors) > 1 {
				return types.M{"$or": ors}
			}
			return ors[0]
		}
	}
	return query
}

func addWriteACL(query types.M, acl []string) types.M {
	newQuery := utils.CopyMap(query)
	writePerms := types.S{nil}
	for _, a := range acl {
		writePerms = append(writePerms, a)
	}
	newQuery["_wperm"] = types.M{"$in": writePerms}
	return newQuery
}

func addReadACL(query types.M, acl []string) types.M {
	newQuery := utils.CopyMap(query)
	orParts := types.S{nil, "*"}
	for _, a := range acl {
		orParts = append(orParts, a)
	}
	newQuery["_rperm"] = types.M{"$in": orParts}
	return newQuery
}

var specialQuerykeys = []string{"$and", "$or", "_rperm", "_wperm", "_perishable_token", "_email_verify_token"}

func validateQuery(query types.M) error {
	if query == nil {
		return nil
	}

	if _, ok := query["ACL"]; ok {
		return errs.E(errs.InvalidQuery, "Cannot query on ACL.")
	}

	if or, ok := query["$or"]; ok {
		if arr, ok := or.([]interface{}); ok {
			for _, a := range arr {
				err := validateQuery(a.(map[string]interface{}))
				if err != nil {
					return err
				}
			}
		} else {
			return errs.E(errs.InvalidQuery, "Bad $or format - use an array value.")
		}
	}

	if and, ok := query["$and"]; ok {
		if arr, ok := and.([]interface{}); ok {
			for _, a := range arr {
				err := validateQuery(a.(map[string]interface{}))
				if err != nil {
					return err
				}
			}
		} else {
			return errs.E(errs.InvalidQuery, "Bad $and format - use an array value.")
		}
	}

	for key := range query {
		include := false
		for _, v := range specialQuerykeys {
			if key == v {
				include = true
				break
			}
		}
		if include == false {
			match, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_\.]*$`, key)
			if match == false {
				return errs.E(errs.InvalidKeyName, "Invalid key name: "+key)
			}
		}
	}

	return nil
}
