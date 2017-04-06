//Package orm 数据库操作模块，当前只对接了 MongoDB
package orm

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/lfq7413/tomato/cache"
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/storage"
	"github.com/lfq7413/tomato/storage/mongo"
	"github.com/lfq7413/tomato/storage/postgres"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// TomatoDBController ...
var TomatoDBController *DBController

// Adapter ...
var Adapter storage.Adapter

var schemaCache *cache.SchemaCache
var schemaPromise *Schema

// init 初始化 Mongo 适配器
func init() {
	if config.TConfig.DatabaseType == "MongoDB" {
		Adapter = mongo.NewMongoAdapter("tomato", storage.OpenMongoDB())
	} else if config.TConfig.DatabaseType == "PostgreSQL" {
		Adapter = postgres.NewPostgresAdapter("tomato", storage.OpenPostgreSQL())
	} else {
		// 默认连接 MongoDB
		Adapter = mongo.NewMongoAdapter("tomato", storage.OpenMongoDB())
	}
	schemaCache = cache.NewSchemaCache(config.TConfig.SchemaCacheTTL, config.TConfig.EnableSingleSchemaCache)
	TomatoDBController = &DBController{}
}

// DBController 数据库操作类
type DBController struct {
}

// CollectionExists 检测表是否存在
func (d *DBController) CollectionExists(className string) bool {
	return Adapter.ClassExists(className)
}

// PurgeCollection 清除类
func (d *DBController) PurgeCollection(className string) error {
	schema := d.LoadSchema(nil)
	sch, err := schema.GetOneSchema(className, false, nil)
	if err != nil {
		return err
	}
	return Adapter.DeleteObjectsByQuery(className, sch, types.M{})
}

// Find 从指定表中查询数据，查询到的数据放入 list 中
// 如果查询的是 count ，结果也会放入 list，并且只有这一个元素
// options 中的选项包括：skip、limit、sort、keys、count、acl
func (d *DBController) Find(className string, query, options types.M) (types.S, error) {
	if options == nil {
		options = types.M{}
	}
	if query == nil {
		query = types.M{}
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
	if v, ok := options["op"].(string); ok && v != "" {
		op = v
	} else {
		if _, ok := query["objectId"].(string); ok {
			if len(query) == 1 {
				op = "get"
			} else {
				op = "find"
			}
		} else {
			op = "find"
		}
	}

	classExists := true

	schema := d.LoadSchema(nil)
	parseFormatSchema, err := schema.GetOneSchema(className, isMaster, nil)
	if err != nil {
		return nil, err
	}
	if len(parseFormatSchema) == 0 {
		classExists = false
		parseFormatSchema["fields"] = types.M{}
	}

	if keys, ok := options["sort"].([]string); ok {
		for i, key := range keys {
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

			if match, _ := regexp.MatchString(`^authData\.([a-zA-Z0-9_]+)\.id$`, key); match {
				return nil, errs.E(errs.InvalidKeyName, "Cannot sort by "+key)
			}

			if fieldNameIsValid(key) == false {
				return nil, errs.E(errs.InvalidKeyName, "Invalid field name: "+key)
			}

			keys[i] = prefix + key
		}
		options["sort"] = keys
	}

	// 校验当前用户是否能对表进行 find 或者 get 操作
	if isMaster == false {
		err := schema.validatePermission(className, aclGroup, op)
		if err != nil {
			return nil, err
		}
	}

	// 处理 $relatedTo
	query = d.reduceRelationKeys(className, query)
	// 处理 relation 字段上的 $in
	query = d.reduceInRelation(className, query, schema)

	if isMaster == false {
		query = d.addPointerPermissions(schema, className, op, query, aclGroup)
	}
	if query == nil {
		if op == "get" {
			return nil, errs.E(errs.ObjectNotFound, "Object not found.")
		}
		// 如果需要计算 count ，则默认返回  0
		if options["count"] != nil {
			return types.S{0}, nil
		}
		return types.S{}, nil
	}

	// 组装 acl 查询条件，查找可被当前用户访问的对象
	if isMaster == false {
		query = addReadACL(query, aclGroup)
	}

	err = validateQuery(query)
	if err != nil {
		return nil, err
	}

	// 获取 count
	if options["count"] != nil {
		if classExists == false {
			return types.S{0}, nil
		}
		count, err := Adapter.Count(className, parseFormatSchema, query)
		if err != nil {
			return nil, err
		}
		return types.S{count}, nil
	}

	if classExists == false {
		return types.S{}, nil
	}

	// 执行查询操作
	objects, err := Adapter.Find(className, parseFormatSchema, query, options)
	if err != nil {
		return nil, err
	}
	results := types.S{}
	for _, object := range objects {
		object = untransformObjectACL(object)
		result := filterSensitiveData(isMaster, aclGroup, className, object)
		results = append(results, result)
	}
	return results, nil
}

// Destroy 从指定表中删除数据
func (d *DBController) Destroy(className string, query types.M, options types.M) error {
	if query == nil {
		query = types.M{}
	}
	if options == nil {
		options = types.M{}
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

	schema := d.LoadSchema(nil)
	if isMaster == false {
		err := schema.validatePermission(className, aclGroup, "delete")
		if err != nil {
			return err
		}
	}

	if isMaster == false {
		query = d.addPointerPermissions(schema, className, "delete", query, aclGroup)
		if query == nil {
			return errs.E(errs.ObjectNotFound, "Object not found.")
		}
	}

	if isMaster == false {
		query = addWriteACL(query, aclGroup)
	}

	err := validateQuery(query)
	if err != nil {
		return err
	}

	parseFormatSchema, err := schema.GetOneSchema(className, false, nil)
	if err != nil {
		return err
	}
	if len(parseFormatSchema) == 0 {
		parseFormatSchema["fields"] = types.M{}
	}

	err = Adapter.DeleteObjectsByQuery(className, parseFormatSchema, query)
	if err != nil {
		// 排除 _Session，避免在修改密码时因为没有 Session 失败
		if className == "_Session" && errs.GetErrorCode(err) == errs.ObjectNotFound {
			return nil
		}
		return err
	}

	return nil
}

var specialKeysForUpdate = map[string]bool{
	"_hashed_password":               true,
	"_perishable_token":              true,
	"_email_verify_token":            true,
	"_email_verify_token_expires_at": true,
	"_account_lockout_expires_at":    true,
	"_failed_login_count":            true,
	"_perishable_token_expires_at":   true,
	"_password_changed_at":           true,
	"_password_history":              true,
}

// Update 更新对象
// options 中的参数包括：acl、many、upsert
// skipSanitization 默认为 false
func (d *DBController) Update(className string, query, update, options types.M, skipSanitization bool) (types.M, error) {
	if len(query) == 0 {
		return types.M{}, nil
	}
	if len(update) == 0 {
		return types.M{}, nil
	}
	if options == nil {
		options = types.M{}
	}
	originalUpdate := update
	// 复制数据，不要修改原数据
	update = utils.CopyMap(update)

	var many bool
	if v, ok := options["many"].(bool); ok {
		many = v
	}
	var upsert bool
	if v, ok := options["upsert"].(bool); ok {
		upsert = v
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

	schema := d.LoadSchema(nil)
	if isMaster == false {
		err := schema.validatePermission(className, aclGroup, "update")
		if err != nil {
			return nil, err
		}
	}
	// 处理 Relation
	d.handleRelationUpdates(className, utils.S(query["objectId"]), update)

	// 添加用户权限
	if isMaster == false {
		query = d.addPointerPermissions(schema, className, "update", query, aclGroup)
	}
	if query == nil {
		return types.M{}, nil
	}

	// 组装 acl 查询条件，查找可被当前用户修改的对象
	if isMaster == false {
		query = addWriteACL(query, aclGroup)
	}

	err := validateQuery(query)
	if err != nil {
		return nil, err
	}

	sch, err := schema.GetOneSchema(className, false, nil)
	if err != nil {
		return nil, err
	}
	if len(sch) == 0 {
		sch["fields"] = types.M{}
	}

	for fieldName, v := range update {
		if match, _ := regexp.MatchString(`^authData\.([a-zA-Z0-9_]+)\.id$`, fieldName); match {
			return nil, errs.E(errs.InvalidKeyName, "Invalid field name for update: "+fieldName)
		}
		fieldName = strings.Split(fieldName, ".")[0]
		if fieldNameIsValid(fieldName) == false && specialKeysForUpdate[fieldName] == false {
			return nil, errs.E(errs.InvalidKeyName, "Invalid field name for update: "+fieldName)
		}

		if updateOperation := utils.M(v); updateOperation != nil {
			for innerKey := range updateOperation {
				if strings.Index(innerKey, "$") > -1 || strings.Index(innerKey, ".") > -1 {
					return nil, errs.E(errs.InvalidNestedKey, "Nested keys should not contain the '$' or '.' characters")
				}
			}
		}
	}

	update = transformObjectACL(update)
	transformAuthData(className, update, sch)
	var result types.M
	if many {
		err := Adapter.UpdateObjectsByQuery(className, sch, query, update)
		if err != nil {
			return nil, err
		}
		result = types.M{}
	} else if upsert {
		err := Adapter.UpsertOneObject(className, sch, query, update)
		if err != nil {
			return nil, err
		}
		result = types.M{}
	} else {
		var err error
		result, err = Adapter.FindOneAndUpdate(className, sch, query, update)
		if err != nil {
			return nil, err
		}
	}

	// 不处理 many 、 upsert 时的操作结果，仅处理 FindOneAndUpdate 的结果
	if many == false && upsert == false && len(result) == 0 {
		return nil, errs.E(errs.ObjectNotFound, "Object not found.")
	}

	if skipSanitization {
		return result, nil
	}

	// 返回经过修改的字段
	response := sanitizeDatabaseResult(originalUpdate, result)

	return response, nil
}

// sanitizeDatabaseResult 处理数据库返回结果
func sanitizeDatabaseResult(originalObject, result types.M) types.M {
	response := types.M{}
	if originalObject == nil || result == nil {
		return response
	}

	// 检测是否是对字段的操作
	for key, value := range originalObject {
		if keyUpdate := utils.M(value); keyUpdate != nil {
			if op := utils.S(keyUpdate["__op"]); op != "" {
				if op == "Add" || op == "AddUnique" || op == "Remove" || op == "Increment" {
					// 只把操作的字段放入返回结果中
					if v, ok := result[key]; ok {
						response[key] = v
					}
				}
			}
		}
	}

	return response
}

// Create 创建对象
func (d *DBController) Create(className string, object, options types.M) error {
	if options == nil {
		options = types.M{}
	}
	if object == nil {
		object = types.M{}
	}
	// 复制数据，不要修改原数据
	object = utils.CopyMapM(object)

	object = transformObjectACL(object)

	if v, ok := object["createdAt"]; ok {
		object["createdAt"] = types.M{
			"__type": "Date",
			"iso":    v,
		}
	}
	if v, ok := object["updatedAt"]; ok {
		object["updatedAt"] = types.M{
			"__type": "Date",
			"iso":    v,
		}
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

	err := d.validateClassName(className)
	if err != nil {
		return err
	}

	schema := d.LoadSchema(nil)
	if isMaster == false {
		err := schema.validatePermission(className, aclGroup, "create")
		if err != nil {
			return err
		}
	}

	// 处理 Relation
	err = d.handleRelationUpdates(className, "", object)
	if err != nil {
		return err
	}

	err = schema.EnforceClassExists(className)
	if err != nil {
		return err
	}

	schema.reloadData(nil)

	sch, err := schema.GetOneSchema(className, true, nil)
	if err != nil {
		return err
	}

	transformAuthData(className, object, sch)
	flattenUpdateOperatorsForCreate(object)

	// 无需调用 sanitizeDatabaseResult
	return Adapter.CreateObject(className, convertSchemaToAdapterSchema(sch), object)
}

// validateClassName 校验表名是否合法
func (d *DBController) validateClassName(className string) error {
	if ClassNameIsValid(className) == false {
		return errs.E(errs.InvalidClassName, "invalid className: "+className)
	}
	return nil
}

// handleRelationUpdates 处理 Relation 相关操作
// TODO 修改 handleRelationUpdates 调用时机：在 update 或 create 成功之后再调用
func (d *DBController) handleRelationUpdates(className, objectID string, update types.M) error {
	if update == nil {
		return nil
	}
	objID := objectID
	if utils.S(update["objectId"]) != "" {
		objID = utils.S(update["objectId"])
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
		if op == nil || utils.M(op) == nil || utils.M(op)["__op"] == nil {
			return nil
		}
		opMap := utils.M(op)
		p := utils.S(opMap["__op"])
		if p == "AddRelation" {
			delete(update, key)
			// 添加 Relation 对象
			if objects := utils.A(opMap["objects"]); objects != nil {
				for _, object := range objects {
					if obj := utils.M(object); obj != nil {
						if relationID := utils.S(obj["objectId"]); relationID != "" {
							err := d.addRelation(key, className, objID, relationID)
							if err != nil {
								return err
							}
						}
					}
				}
			}
		} else if p == "RemoveRelation" {
			delete(update, key)
			// 删除 Relation 对象
			if objects := utils.A(opMap["objects"]); objects != nil {
				for _, object := range objects {
					if obj := utils.M(object); obj != nil {
						if relationID := utils.S(obj["objectId"]); relationID != "" {
							err := d.removeRelation(key, className, objID, relationID)
							if err != nil {
								return err
							}
						}
					}
				}
			}
		} else if p == "Batch" {
			delete(update, key)
			// 批处理 Relation 对象
			if ops := utils.A(opMap["ops"]); ops != nil {
				for _, x := range ops {
					err := process(x, key)
					if err != nil {
						return err
					}
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

var relationSchema = types.M{
	"fields": types.M{
		"relatedId": types.M{"type": "String"},
		"owningId":  types.M{"type": "String"},
	},
}

// addRelation 把对象 id 加入 _Join 表，表名为 _Join:key:fromClassName
func (d *DBController) addRelation(key, fromClassName, fromID, toID string) error {
	doc := types.M{
		"relatedId": toID,
		"owningId":  fromID,
	}
	className := "_Join:" + key + ":" + fromClassName
	return Adapter.UpsertOneObject(className, relationSchema, doc, doc)
}

// removeRelation 把对象 id 从 _Join 表中删除，表名为 _Join:key:fromClassName
func (d *DBController) removeRelation(key, fromClassName, fromID, toID string) error {
	doc := types.M{
		"relatedId": toID,
		"owningId":  fromID,
	}
	className := "_Join:" + key + ":" + fromClassName
	err := Adapter.DeleteObjectsByQuery(className, relationSchema, doc)
	if err != nil {
		if errs.GetErrorCode(err) == errs.ObjectNotFound {
			return nil
		}
		return err
	}
	return nil
}

// ValidateObject 校验对象是否合法
func (d *DBController) ValidateObject(className string, object, query, options types.M) error {
	schema := d.LoadSchema(nil)

	if options == nil {
		options = types.M{}
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

	if !isMaster {
		err := d.canAddField(schema, className, object, aclGroup)
		if err != nil {
			return err
		}
	}

	err := schema.validateObject(className, object, query)
	if err != nil {
		return err
	}

	return nil
}

// LoadSchema 加载 Schema，仅加载一次
func (d *DBController) LoadSchema(options types.M) *Schema {
	if options == nil {
		options = types.M{"clearCache": false}
	}
	if c, ok := options["clearCache"].(bool); ok && c {
		schemaPromise = Load(Adapter, schemaCache, options)
		return schemaPromise
	}
	if schemaPromise == nil {
		schemaPromise = Load(Adapter, schemaCache, options)
	}
	return schemaPromise
}

// DeleteEverything 删除所有表数据，仅用于测试
func (d *DBController) DeleteEverything() {
	schemaCache.Clear()
	schemaPromise = nil
	Adapter.DeleteAllClasses()
}

// RedirectClassNameForKey 返回指定类的字段所对应的类型
// 如果 key 字段的属性为 relation<classA> ，则返回 classA
func (d *DBController) RedirectClassNameForKey(className, key string) string {
	schema := d.LoadSchema(nil)
	t := schema.getExpectedType(className, key)
	if t != nil && utils.S(t["type"]) == "Relation" {
		return utils.S(t["targetClass"])
	}
	return className
}

// canAddField 检测是否能添加字段到类上
func (d *DBController) canAddField(schema *Schema, className string, object types.M, acl []string) error {
	if schema == nil || schema.data == nil || schema.data[className] == nil {
		return nil
	}

	classSchema := utils.M(schema.data[className])

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
// 已知父对象，查找子对象
func (d *DBController) reduceRelationKeys(className string, query types.M) types.M {
	if query == nil {
		return query
	}
	// 处理 $or 数组中的数据，并替换回去
	if query["$or"] != nil {
		var subQuerys types.S
		subQuerys = utils.A(query["$or"])
		for i, v := range subQuerys {
			aQuery := utils.M(v)
			subQuerys[i] = d.reduceRelationKeys(className, aQuery)
		}
		query["$or"] = subQuerys
		return query
	}

	if r, ok := query["$relatedTo"]; ok {
		delete(query, "$relatedTo")
		relatedTo := utils.M(r)
		if relatedTo == nil {
			return query
		}
		key := utils.S(relatedTo["key"])
		object := utils.M(relatedTo["object"])
		if key == "" || object == nil {
			return query
		}
		objClassName := utils.S(object["className"])
		objID := utils.S(object["objectId"])
		ids := d.relatedIds(objClassName, key, objID)
		query = d.addInObjectIdsIds(ids, query)
		query = d.reduceRelationKeys(className, query)
	}

	return query
}

// relatedIds 从 Join 表中查询 ids ，表名：_Join:key:className
func (d *DBController) relatedIds(className, key, owningID string) types.S {
	ids := types.S{}
	results, err := Adapter.Find(joinTableName(className, key), relationSchema, types.M{"owningId": owningID}, types.M{})
	if err != nil {
		return ids
	}

	for _, result := range results {
		id := result["relatedId"]
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
func (d *DBController) addInObjectIdsIds(ids types.S, query types.M) types.M {
	// 当两个参数均为空时，不进行处理
	if ids == nil && query == nil {
		return query
	}
	// 参数有一个不为空时，即进行处理
	if ids == nil {
		ids = types.S{}
	}
	if query == nil {
		query = types.M{}
	}
	// coll 保存四种类型的 ID 列表：idsFromString、idsFromEq、idsFromIn、ids
	coll := map[string]types.S{}
	idsFromString := types.S{}
	if id := utils.S(query["objectId"]); id != "" {
		idsFromString = append(idsFromString, id)
	}
	coll["idsFromString"] = idsFromString

	idsFromEq := types.S{}
	if eqid := utils.M(query["objectId"]); eqid != nil {
		if id, ok := eqid["$eq"]; ok {
			// TODO 结果中包含了 $in ，是否需要删除 $eq ？
			if v := utils.S(id); v != "" {
				idsFromEq = append(idsFromEq, v)
			}
		}
	}
	coll["idsFromEq"] = idsFromEq

	idsFromIn := types.S{}
	if inid := utils.M(query["objectId"]); inid != nil {
		if id, ok := inid["$in"]; ok {
			if v := utils.A(id); v != nil {
				idsFromIn = append(idsFromIn, v...)
			}
		}
	}
	coll["idsFromIn"] = idsFromIn

	coll["ids"] = ids

	// 统计 idsFromString idsFromEq idsFromIn ids 中的共同元素加入到 $in 中
	// max 以上4个集合中不为空的个数，也就是说 某个 objectId 出现的次数应该等于 max 才能加入到 $in 中查询
	max := 0
	for k, v := range coll {
		// 删除空集合
		if len(v) > 0 {
			max++
		} else {
			delete(coll, k)
		}
	}

	// idsColl 统计每个 objectId 出现的次数
	idsColl := map[string]int{}
	for _, c := range coll {
		// 从每个集合中取出 objectId
		// idColl 统计当前集合中出现过的 ID ，用于去重
		idColl := map[string]int{}
		for _, v := range c {
			id := utils.S(v)
			// 并去除重复，只处理未出现过的 ID
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

	// queryIn 统计出现次数为 max 的 objectId
	queryIn := types.S{}
	for k, v := range idsColl {
		if v == max {
			queryIn = append(queryIn, k)
		}
	}

	// 替换 objectId
	if objectID := utils.M(query["objectId"]); objectID != nil {
		objectID["$in"] = queryIn
	} else {
		query["objectId"] = types.M{"$in": queryIn}
	}
	return query
}

// addNotInObjectIdsIds 添加 ids 到查询条件中，应该取 $ne $nin ids 的并集
// 替换 objectId 为：
// "objectId":{"$nin":["id","id2"]}
func (d *DBController) addNotInObjectIdsIds(ids types.S, query types.M) types.M {
	// 当两个参数均为空时，不进行处理
	if ids == nil && query == nil {
		return query
	}
	// 参数有一个不为空时，即进行处理
	if ids == nil {
		ids = types.S{}
	}
	if query == nil {
		query = types.M{}
	}
	// coll 保存两种类型的 ID 列表：idsFromNin、ids
	coll := map[string]types.S{}
	idsFromNin := types.S{}
	if ninid := utils.M(query["objectId"]); ninid != nil {
		if id, ok := ninid["$nin"]; ok {
			if v := utils.A(id); v != nil {
				idsFromNin = append(idsFromNin, v...)
			}
		}
	}
	coll["idsFromNin"] = idsFromNin

	coll["ids"] = ids

	idsColl := map[string]int{}
	for _, c := range coll {
		// 从每个集合中取出 objectId
		for _, v := range c {
			id := utils.S(v)
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

	// 替换 objectId
	if v, ok := query["objectId"]; ok == false || v == nil {
		query["objectId"] = types.M{"$nin": queryNin}
	} else if v, ok := query["objectId"].(string); ok {
		query["objectId"] = types.M{
			"$eq":  v,
			"$nin": queryNin,
		}
	} else if objectID := utils.M(query["objectId"]); objectID != nil {
		objectID["$nin"] = queryNin
	} else {
		query["objectId"] = types.M{"$nin": queryNin}
	}
	return query
}

// reduceInRelation 处理查询条件中，作用于 relation 类型字段上的 $in $ne $nin $eq 或者等于某对象
// 例如 classA 中的 字段 key 为 relation<classB> 类型，查找 key 中包含指定 classB 对象的 classA
// query = {"key":{"$in":[]}}
// 已知子对象，查找父对象
func (d *DBController) reduceInRelation(className string, query types.M, schema *Schema) types.M {
	if query == nil {
		return query
	}
	// 处理 $or 数组中的数据，并替换回去
	if query["$or"] != nil {
		var ors types.S
		ors = utils.A(query["$or"])
		for i, v := range ors {
			aQuery := utils.M(v)
			ors[i] = d.reduceInRelation(className, aQuery, schema)
		}
		query["$or"] = ors
		return query
	}

	for key, v := range query {
		op := utils.M(v)
		if op != nil && (op["$in"] != nil || op["$ne"] != nil || op["$nin"] != nil || op["$eq"] != nil || utils.S(op["__type"]) == "Pointer") {
			if schema == nil {
				continue
			}
			// 只处理 relation 类型
			t := schema.getExpectedType(className, key)
			if t == nil || utils.S(t["type"]) != "Relation" {
				continue
			}

			// 取出所有限制条件
			relatedIds := []types.S{}
			isNegation := []bool{}
			for constraintKey, value := range op {
				if constraintKey == "objectId" {
					if utils.S(value) != "" {
						ids := types.S{value}
						relatedIds = append(relatedIds, ids)
						isNegation = append(isNegation, false)
					}
				} else if constraintKey == "$in" {
					in := utils.A(value)
					ids := types.S{}
					for _, v := range in {
						if r := utils.M(v); r != nil {
							if utils.S(r["objectId"]) != "" {
								ids = append(ids, r["objectId"])
							}
						}
					}
					// 只计算有效的 objectId
					if len(ids) > 0 {
						relatedIds = append(relatedIds, ids)
						isNegation = append(isNegation, false)
					}
				} else if constraintKey == "$nin" {
					nin := utils.A(value)
					ids := types.S{}
					for _, v := range nin {
						if r := utils.M(v); r != nil {
							if utils.S(r["objectId"]) != "" {
								ids = append(ids, r["objectId"])
							}
						}
					}
					// 只计算有效的 objectId
					if len(ids) > 0 {
						relatedIds = append(relatedIds, ids)
						isNegation = append(isNegation, true)
					}
				} else if constraintKey == "$ne" {
					if ne := utils.M(value); ne != nil {
						if utils.S(ne["objectId"]) != "" {
							ids := types.S{ne["objectId"]}
							relatedIds = append(relatedIds, ids)
							isNegation = append(isNegation, true)
						}
					}
				} else if constraintKey == "$eq" {
					if eq := utils.M(value); eq != nil {
						if utils.S(eq["objectId"]) != "" {
							ids := types.S{eq["objectId"]}
							relatedIds = append(relatedIds, ids)
							isNegation = append(isNegation, false)
						}
					}
				}
			}

			delete(query, key)

			// 应用所有限制条件
			for i, relatedID := range relatedIds {
				// 此处 relatedID 含有至少一个元素
				// 从 Join 表中查找的 ids，替换查询条件
				ids := d.owningIds(className, key, relatedID)
				if isNegation[i] {
					query = d.addNotInObjectIdsIds(ids, query)
				} else {
					query = d.addInObjectIdsIds(ids, query)
				}
			}
		}
	}

	return query
}

// owningIds 从 Join 表中查询 relatedIds 对应的父对象
func (d *DBController) owningIds(className, key string, relatedIds types.S) types.S {
	ids := types.S{}
	query := types.M{
		"relatedId": types.M{
			"$in": relatedIds,
		},
	}
	results, err := Adapter.Find(joinTableName(className, key), relationSchema, query, types.M{})
	if err != nil {
		return ids
	}

	for _, result := range results {
		ids = append(ids, result["owningId"])
	}
	return ids
}

// filterSensitiveData 对 _User 表数据进行特殊处理
func filterSensitiveData(isMaster bool, aclGroup []string, className string, object types.M) types.M {
	if className != "_User" {
		return object
	}
	if object == nil {
		return object
	}
	// 以下单独处理 _User 类
	if _, ok := object["_hashed_password"]; ok {
		object["password"] = object["_hashed_password"]
		delete(object, "_hashed_password")
	}

	delete(object, "sessionToken")

	if isMaster {
		return object
	}

	delete(object, "_email_verify_token")
	delete(object, "_perishable_token")
	delete(object, "_perishable_token_expires_at")
	delete(object, "_tombstone")
	delete(object, "_email_verify_token_expires_at")
	delete(object, "_failed_login_count")
	delete(object, "_account_lockout_expires_at")
	delete(object, "_password_changed_at")

	// 当前用户返回所有信息
	if aclGroup == nil {
		aclGroup = []string{}
	}
	id := utils.S(object["objectId"])
	for _, v := range aclGroup {
		if v == id {
			return object
		}
	}
	delete(object, "authData")
	return object
}

// DeleteSchema 删除类
func (d *DBController) DeleteSchema(className string) error {
	schemaController := d.LoadSchema(types.M{"clearCache": true})
	schema, err := schemaController.GetOneSchema(className, false, types.M{"clearCache": true})
	if err != nil {
		return err
	}
	if schema == nil || len(schema) == 0 {
		schema = types.M{"fields": types.M{}}
	}

	exist := d.CollectionExists(className)
	if exist {
		count, err := Adapter.Count(className, types.M{"fields": types.M{}}, types.M{})
		if err != nil {
			return err
		}
		if count > 0 {
			return errs.E(errs.ClassNotEmpty, "Class "+className+" is not empty, contains "+strconv.Itoa(count)+" objects, cannot drop schema.")
		}
	}

	result, err := Adapter.DeleteClass(className)
	if err != nil {
		return err
	}
	if result != nil {
		if fields := utils.M(schema["fields"]); fields != nil {
			for fieldName, v := range fields {
				if fieldType := utils.M(v); fieldType != nil {
					if utils.S(fieldType["type"]) == "Relation" {
						_, err = Adapter.DeleteClass(joinTableName(className, fieldName))
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}
	d.LoadSchema(types.M{"clearCache": true})
	return nil
}

// addPointerPermissions 添加查询用户权限，perms[className][readUserFields] 中保存的是字段名，该字段中的内容是：有权限进行读操作的用户
func (d *DBController) addPointerPermissions(schema *Schema, className string, operation string, query types.M, aclGroup []string) types.M {
	if schema == nil {
		return query
	}

	if schema.testBaseCLP(className, aclGroup, operation) {
		return query
	}

	perms := schema.perms[className]
	// 根据当前操作确定是读还是写
	var field string
	if operation == "get" || operation == "find" {
		field = "readUserFields"
	} else {
		field = "writeUserFields"
	}

	// 取出用户 ID
	if aclGroup == nil {
		aclGroup = []string{}
	}
	userACL := []string{}
	for _, acl := range aclGroup {
		if strings.HasPrefix(acl, "role:") == false && acl != "*" {
			userACL = append(userACL, acl)
		}
	}

	// 查找保存 有读或写权限 用户 的字段，然后使用 and 与原查询请求进行拼装
	if perms != nil {
		if permsMap := utils.M(perms); permsMap != nil {
			if permFields := utils.A(permsMap[field]); permFields != nil {
				if permFields != nil && len(permFields) > 0 {
					// 用户 ID 有多个时，表示没有正确处理 ACL
					if len(userACL) != 1 {
						return nil
					}
					userID := userACL[0]
					userPointer := types.M{
						"__type":    "Pointer",
						"className": "_User",
						"objectId":  userID,
					}

					// 使用 and 拼装请求
					ors := []types.M{}
					for _, key := range permFields {
						q := types.M{
							utils.S(key): userPointer,
						}
						and := types.M{
							"$and": types.S{q, query},
						}
						ors = append(ors, and)
					}
					// 有多个权限字段时，使用 or 再次拼装
					if len(ors) > 1 {
						return types.M{"$or": ors}
					}
					return ors[0]
				}
			}
		}

	}
	return query
}

// PerformInitialization 初始化数据库索引
func (d *DBController) PerformInitialization() {
	requiredUserFields := types.M{}
	defaultUserColumns := types.M{}
	for k, v := range DefaultColumns["_Default"] {
		defaultUserColumns[k] = v
	}
	for k, v := range DefaultColumns["_User"] {
		defaultUserColumns[k] = v
	}
	requiredUserFields["fields"] = defaultUserColumns
	d.LoadSchema(nil).EnforceClassExists("_User")
	Adapter.EnsureUniqueness("_User", requiredUserFields, []string{"username"})
	Adapter.EnsureUniqueness("_User", requiredUserFields, []string{"email"})
	Adapter.PerformInitialization(types.M{"VolatileClassesSchemas": volatileClassesSchemas()})
}

func addWriteACL(query types.M, acl []string) types.M {
	if query == nil {
		query = types.M{}
	}
	if acl == nil {
		acl = []string{}
	}
	newQuery := utils.CopyMap(query)
	writePerms := types.S{nil}
	for _, a := range acl {
		writePerms = append(writePerms, a)
	}
	newQuery["_wperm"] = types.M{"$in": writePerms}
	return newQuery
}

func addReadACL(query types.M, acl []string) types.M {
	if query == nil {
		query = types.M{}
	}
	if acl == nil {
		acl = []string{}
	}
	newQuery := utils.CopyMap(query)
	orParts := types.S{nil, "*"}
	for _, a := range acl {
		orParts = append(orParts, a)
	}
	newQuery["_rperm"] = types.M{"$in": orParts}
	return newQuery
}

var specialQuerykeys = map[string]bool{
	"$and":                           true,
	"$or":                            true,
	"_rperm":                         true,
	"_wperm":                         true,
	"_perishable_token":              true,
	"_perishable_token_expires_at":   true,
	"_email_verify_token":            true,
	"_email_verify_token_expires_at": true,
	"_account_lockout_expires_at":    true,
	"_failed_login_count":            true,
	"_password_changed_at":           true,
}

func validateQuery(query types.M) error {
	if query == nil {
		return nil
	}

	if _, ok := query["ACL"]; ok {
		return errs.E(errs.InvalidQuery, "Cannot query on ACL.")
	}

	if or, ok := query["$or"]; ok {
		if arr := utils.A(or); arr != nil {
			for _, a := range arr {
				subQuery := utils.M(a)
				if subQuery == nil {
					return errs.E(errs.InvalidQuery, "Bad $or format - invalid sub query.")
				}
				err := validateQuery(subQuery)
				if err != nil {
					return err
				}
			}
		} else {
			return errs.E(errs.InvalidQuery, "Bad $or format - use an array value.")
		}
	}

	if and, ok := query["$and"]; ok {
		if arr := utils.A(and); arr != nil {
			for _, a := range arr {
				subQuery := utils.M(a)
				if subQuery == nil {
					return errs.E(errs.InvalidQuery, "Bad $and format - invalid sub query.")
				}
				err := validateQuery(subQuery)
				if err != nil {
					return err
				}
			}
		} else {
			return errs.E(errs.InvalidQuery, "Bad $and format - use an array value.")
		}
	}

	for key := range query {
		// 检测 $options 是否为 imxs
		if condition := utils.M(query[key]); condition != nil {
			if condition["$regex"] != nil {
				if op, ok := condition["$options"].(string); ok {
					b, _ := regexp.MatchString(`^[imxs]+$`, op)
					if b == false {
						// 无效值
						return errs.E(errs.InvalidQuery, "Bad $options value for query: "+op)
					}
				}
			}
		}

		if specialQuerykeys[key] == true {
			continue
		}
		match, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_\.]*$`, key)
		if match == false {
			return errs.E(errs.InvalidKeyName, "Invalid key name: "+key)
		}
	}

	return nil
}

// transformObjectACL 转换对象中的 ACL 字段
// {
// 	"ACL":{
// 		"userid":{
// 			"read":true,
// 			"write":true
// 		},
// 		"role:xxx":{
// 			"read":true,
// 			"write":true
// 		}
// 		"*":{
// 			"read":true
// 		}
// 	}
// }
// ==>
// {
// 	"_rperm":["userid","role:xxx","*"],
// 	"_wperm":["userid","role:xxx"],
// }
func transformObjectACL(result types.M) types.M {
	if result == nil {
		return result
	}

	if _, ok := result["ACL"]; ok == false {
		return result
	}

	acl := utils.M(result["ACL"])
	if acl == nil {
		delete(result, "ACL")
		return result
	}
	rperm := types.S{}
	wperm := types.S{}
	for entry, v := range acl {
		perm := utils.M(v)
		if perm != nil {
			if perm["read"] != nil {
				rperm = append(rperm, entry)
			}
			if perm["write"] != nil {
				wperm = append(wperm, entry)
			}
		}
	}
	result["_rperm"] = rperm
	result["_wperm"] = wperm
	delete(result, "ACL")

	return result
}

// untransformObjectACL 把数据库格式的 ACL 转换为 API 格式
// {
// 	"_rperm":["userid","role:xxx","*"],
// 	"_wperm":["userid","role:xxx"]
// }
// ==>
// {
// 	"ACL":{
// 		"userid":{
// 			"read":true,
// 			"write":true
// 		},
// 		"role:xxx":{
// 			"read":true,
// 			"write":true
// 		}
// 		"*":{
// 			"read":true
// 		}
// 	}
// }
func untransformObjectACL(output types.M) types.M {
	if output == nil {
		return output
	}

	if output["_rperm"] == nil && output["_wperm"] == nil {
		return output
	}

	acl := types.M{}
	rperm := types.S{}
	wperm := types.S{}
	if output["_rperm"] != nil {
		rperm = utils.A(output["_rperm"])
	}
	if output["_wperm"] != nil {
		wperm = utils.A(output["_wperm"])
	}
	if rperm != nil {
		for _, v := range rperm {
			entry := v.(string)
			if acl[entry] == nil {
				acl[entry] = types.M{"read": true}
			} else {
				var per types.M
				per = utils.M(acl[entry])
				per["read"] = true
				acl[entry] = per
			}
		}
	}
	if wperm != nil {
		for _, v := range wperm {
			entry := v.(string)
			if acl[entry] == nil {
				acl[entry] = types.M{"write": true}
			} else {
				var per types.M
				per = utils.M(acl[entry])
				per["write"] = true
				acl[entry] = per
			}
		}
	}
	output["ACL"] = acl
	delete(output, "_rperm")
	delete(output, "_wperm")

	return output
}

// transformAuthData 转换第三方登录数据
// {
// 	"authData": {
// 		"facebook": {...}
// 	}
// }
// ==>
// {
// 	"_auth_data_facebook": {...}
// }
func transformAuthData(className string, object, schema types.M) {
	if className == "_User" && object != nil {
		if _, ok := object["authData"]; ok == false {
			return
		}
		if authData := utils.M(object["authData"]); authData != nil {
			for provider, providerData := range authData {
				fieldName := "_auth_data_" + provider

				if providerData == nil || utils.M(providerData) == nil || len(utils.M(providerData)) == 0 {
					object[fieldName] = types.M{
						"__op": "Delete",
					}
				} else {
					object[fieldName] = providerData
					if schema != nil {
						if fields := utils.M(schema["fields"]); fields != nil {
							fields[fieldName] = types.M{"type": "Object"}
						}
					}
				}
			}
		}
		delete(object, "authData")
	}
}

// flattenUpdateOperatorsForCreate 处理 Create 数据中的 __op 操作符
func flattenUpdateOperatorsForCreate(object types.M) error {
	if object == nil {
		return nil
	}
	for key, v := range object {
		if value := utils.M(v); value != nil && utils.S(value["__op"]) != "" {
			switch utils.S(value["__op"]) {
			case "Increment":
				if a, ok := value["amount"].(float64); ok {
					object[key] = a
				} else if a, ok := value["amount"].(int); ok {
					object[key] = a
				} else {
					return errs.E(errs.InvalidJSON, "objects to add must be an number")
				}

			case "Add":
				if objects := utils.A(value["objects"]); objects != nil {
					object[key] = value["objects"]
				} else {
					return errs.E(errs.InvalidJSON, "objects to add must be an array")
				}

			case "AddUnique":
				if objects := utils.A(value["objects"]); objects != nil {
					object[key] = value["objects"]
				} else {
					return errs.E(errs.InvalidJSON, "objects to add must be an array")
				}

			case "Remove":
				if objects := utils.A(value["objects"]); objects != nil {
					object[key] = types.S{}
				} else {
					return errs.E(errs.InvalidJSON, "objects to add must be an array")
				}

			case "Delete":
				delete(object, key)

			default:
				return errs.E(errs.CommandUnavailable, "The "+utils.S(value["__op"])+" operator is not supported yet.")
			}
		}
	}
	return nil
}

// InitOrm 初始化 orm ，仅用于测试
func InitOrm(a storage.Adapter) {
	Adapter = a
	schemaCache = cache.NewSchemaCache(5, false)
	TomatoDBController = &DBController{}
}
