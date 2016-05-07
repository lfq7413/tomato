package orm

import (
	"encoding/base64"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// transformKey 把 key 转换为数据库中保存的格式
func transformKey(schema *Schema, className, key string) (string, error) {
	k, _, err := transformKeyValue(schema, className, key, nil, nil)
	if err != nil {
		return "", err
	}
	return k, nil
}

// transformKeyValue 把传入的键值对转换为数据库中保存的格式
// restKey API 格式的字段名
// restValue API 格式的值
// options 设置项
// options 有以下几种选择：
// query: true 表示 restValue 中包含类似 $lt 的查询限制条件
// update: true 表示 restValue 中包含 __op 操作，类似 Add、Delete，需要进行转换
// validate: true 表示需要校验字段名
// 返回转换成 数据库格式的字段名与值
func transformKeyValue(schema *Schema, className, restKey string, restValue interface{}, options types.M) (string, interface{}, error) {
	if options == nil {
		options = types.M{}
	}

	// 检测 key 是否为 内置字段
	key := restKey
	timeField := false
	switch key {
	case "objectId", "_id":
		key = "_id"
	case "createdAt", "_created_at":
		key = "_created_at"
		timeField = true
	case "updatedAt", "_updated_at":
		key = "_updated_at"
		timeField = true
	case "_email_verify_token":
		key = "_email_verify_token"
	case "_perishable_token":
		key = "_perishable_token"
	case "sessionToken", "_session_token":
		key = "_session_token"
	case "expiresAt", "_expiresAt":
		key = "_expiresAt"
		timeField = true
	case "_rperm", "_wperm":
		return key, restValue, nil
	case "$or":
		if options["query"] == nil {
			// 只有查询时才能使用 $or
			return "", nil, errs.E(errs.InvalidKeyName, "you can only use $or in queries")
		}
		querys := utils.SliceInterface(restValue)
		if querys == nil {
			// 待转换值必须为数组类型
			return "", nil, errs.E(errs.InvalidQuery, "bad $or format - use an array value")
		}
		mongoSubqueries := types.S{}
		// 转换 where 查询条件
		for _, v := range querys {
			query, err := transformWhere(schema, className, utils.MapInterface(v))
			if err != nil {
				return "", nil, err
			}
			mongoSubqueries = append(mongoSubqueries, query)
		}
		return "$or", mongoSubqueries, nil
	case "$and":
		if options["query"] == nil {
			// 只有查询时才能使用 and
			return "", "", errs.E(errs.InvalidKeyName, "you can only use $and in queries")
		}
		querys := utils.SliceInterface(restValue)
		if querys == nil {
			// 待转换值必须为数组类型
			return "", nil, errs.E(errs.InvalidQuery, "bad $and format - use an array value")
		}
		mongoSubqueries := types.S{}
		// 转换 where 查询条件
		for _, v := range querys {
			query, err := transformWhere(schema, className, utils.MapInterface(v))
			if err != nil {
				return "", nil, err
			}
			mongoSubqueries = append(mongoSubqueries, query)
		}
		return "$and", mongoSubqueries, nil
	default:
		// 处理第三方 auth 数据，key 的格式为： authData.xxx.id
		authDataMatch, _ := regexp.MatchString(`^authData\.([a-zA-Z0-9_]+)\.id$`, key)
		if authDataMatch {
			if options["query"] != nil {
				// 取出 authData.xxx.id 中的 xxx ，转换为 _auth_data_.xxx.id
				provider := key[len("authData."):(len(key) - len(".id"))]
				return "_auth_data_" + provider + ".id", restValue, nil
			}
			// 只能将其应用查询操作
			return "", nil, errs.E(errs.InvalidKeyName, "can only query on "+key)
		}

		// 校验字段名
		if options["validate"] != nil {
			keyMatch, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_\.]*$`, key)
			if keyMatch == false {
				// 无效的键名
				return "", nil, errs.E(errs.InvalidKeyName, "invalid key name: "+key)
			}
		}
		// 其他字段名不做处理
	}

	// 处理特殊字段名
	expected := ""
	if schema != nil {
		expected = schema.getExpectedType(className, key)
	}

	// 期望类型为 *xxx
	if expected != "" && strings.HasPrefix(expected, "*") {
		key = "_p_" + key
	}
	// 期望类型不存在，但是 restValue 中存在 "__type":"Pointer"
	if expected == "" && restValue != nil {
		op := utils.MapInterface(restValue)
		if op != nil && op["__type"] != nil {
			if utils.String(op["__type"]) == "Pointer" {
				key = "_p_" + key
			}
		}
	}

	inArray := false
	// 期望类型为 array
	if expected == "array" {
		inArray = true
	}

	// 处理查询操作，转换限制条件
	if options["query"] != nil {
		value := transformConstraint(restValue, inArray)
		if value != cannotTransform() {
			return key, value, nil
		}
	}

	// 期望类型为 array，并且转换限制条件失败，并且 restValue 不为 array 类型时，转换为 $all
	if inArray && options["query"] != nil && utils.SliceInterface(restValue) == nil {
		return key, types.M{"$all": types.S{restValue}}, nil
	}

	// 转换原子数据
	value := transformAtom(restValue, false, options)
	if value != cannotTransform() {
		if timeField && utils.String(value) != "" {
			var err error
			value, err = utils.StringtoTime(utils.String(value))
			if err != nil {
				return "", nil, errs.E(errs.InvalidJSON, "Invalid Date value.")
			}
		}
		return key, value, nil
	}

	// ACL 应该在此之前处理，如果依然出现，则返回错误
	if key == "ACL" {
		// 不能在此转换 ACL
		return "", "", errs.E(errs.InvalidKeyName, "There was a problem transforming an ACL.")
	}

	// 转换数组类型
	if valueArray, ok := restValue.([]interface{}); ok {
		if options["query"] != nil {
			// 查询时不能为数组
			return "", "", errs.E(errs.InvalidJSON, "cannot use array as query param")
		}
		outValue := types.S{}
		for _, restObj := range valueArray {
			_, v, err := transformKeyValue(schema, className, restKey, restObj, types.M{"inArray": true})
			if err != nil {
				return "", nil, err
			}
			outValue = append(outValue, v)
		}
		return key, outValue, nil
	}

	// 处理更新操作
	var flatten bool
	if options["update"] == nil {
		flatten = true
	} else {
		flatten = false
	}
	// 处理更新操作中的 "_op"
	value, err := transformUpdateOperator(restValue, flatten)
	if err != nil {
		return "", nil, err
	}
	if value != cannotTransform() {
		return key, value, nil
	}

	// 处理正常的对象
	normalValue := types.M{}
	for subRestKey, subRestValue := range utils.MapInterface(restValue) {
		k, v, err := transformKeyValue(schema, className, subRestKey, subRestValue, types.M{"inObject": true})
		if err != nil {
			return "", nil, err
		}
		normalValue[k] = v
	}
	return key, normalValue, nil
}

// transformConstraint 转换查询限制条件，处理的操作符类似 "$lt", "$gt" 等
// inArray 表示该字段是否为数组类型
func transformConstraint(constraint interface{}, inArray bool) interface{} {
	// TODO 需要根据 MongoDB 文档修正参数
	if constraint == nil || utils.MapInterface(constraint) == nil {
		return cannotTransform()
	}

	// keys is the constraints in reverse alphabetical order.
	// This is a hack so that:
	//   $regex is handled before $options
	//   $nearSphere is handled before $maxDistance
	object := utils.MapInterface(constraint)
	keys := []string{}
	for k := range object {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	answer := types.M{}

	for _, key := range keys {
		switch key {
		// 转换 小于、大于、存在、等于、不等于 操作符
		case "$lt", "$lte", "$gt", "$gte", "$exists", "$ne", "$eq":
			answer[key] = transformAtom(object[key], true, types.M{"inArray": inArray})

		// 转换 包含、不包含 操作符
		case "$in", "$nin":
			arr := utils.SliceInterface(object[key])
			if arr == nil {
				// 必须为数组
				return errs.E(errs.InvalidJSON, "bad "+key+" value")
			}
			answerArr := types.S{}
			for _, v := range arr {
				answerArr = append(answerArr, transformAtom(v, true, types.M{}))
			}
			answer[key] = answerArr

		// 转换 包含所有 操作符，用于数组类型的字段
		case "$all":
			arr := utils.SliceInterface(object[key])
			if arr == nil {
				// 必须为数组
				return errs.E(errs.InvalidJSON, "bad "+key+" value")
			}
			answerArr := types.S{}
			for _, v := range arr {
				answerArr = append(answerArr, transformAtom(v, true, types.M{"inArray": true}))
			}
			answer[key] = answerArr

		// 转换 正则 操作符
		case "$regex":
			s := utils.String(object[key])
			if s == "" {
				// 必须为字符串
				return errs.E(errs.InvalidJSON, "bad regex")
			}
			answer[key] = s

		// 转换 $options 操作符
		case "$options":
			options := utils.String(object[key])
			if answer["$regex"] == nil || options == "" {
				// 无效值
				return errs.E(errs.InvalidQuery, "got a bad $options")
			}
			b, _ := regexp.MatchString(`^[imxs]+$`, options)
			if b == false {
				// 无效值
				return errs.E(errs.InvalidQuery, "got a bad $options")
			}
			answer[key] = options

		// 转换 附近 操作符
		case "$nearSphere":
			point := utils.MapInterface(object[key])
			answer[key] = types.S{point["longitude"], point["latitude"]}

		// 转换 最大距离 操作符，单位是弧度
		case "$maxDistance":
			answer[key] = object[key]

		// 以下三项在 SDK 中未使用，但是在 REST API 中使用了
		case "$maxDistanceInRadians":
			answer["$maxDistance"] = object[key]
		case "$maxDistanceInMiles":
			var distance float64
			if v, ok := object[key].(float64); ok {
				distance = v / 3959
			}
			answer["$maxDistance"] = distance
		case "$maxDistanceInKilometers":
			var distance float64
			if v, ok := object[key].(float64); ok {
				distance = v / 6371
			}
			answer["$maxDistance"] = distance

		case "$select", "$dontSelect":
			// 暂时不支持该参数
			return errs.E(errs.CommandUnavailable, "the "+key+" constraint is not supported yet")

		case "$within":
			within := utils.MapInterface(object[key])
			box := utils.SliceInterface(within["$box"])
			if box == nil || len(box) != 2 {
				// 参数不正确
				return errs.E(errs.InvalidJSON, "malformatted $within arg")
			}
			box1 := utils.MapInterface(box[0])
			box2 := utils.MapInterface(box[1])
			// MongoDB 2.4 中 $within 替换为了 $geoWithin
			answer["$geoWithin"] = types.M{
				"$box": types.S{
					types.S{box1["longitude"], box1["latitude"]},
					types.S{box2["longitude"], box2["latitude"]},
				},
			}

		default:
			b, _ := regexp.MatchString(`^\$+`, key)
			if b {
				// 其他以 $ 开头的操作符为无效参数
				return errs.E(errs.InvalidJSON, "bad constraint: "+key)
			}
			return cannotTransform()
		}
	}

	return answer
}

func transformAtom(atom interface{}, force bool, options types.M) interface{} {
	// TODO 处理错误
	if options == nil {
		options = types.M{}
	}
	inArray := false
	inObject := false
	if v, ok := options["inArray"].(bool); ok {
		inArray = v
	}
	if v, ok := options["inObject"].(bool); ok {
		inObject = v
	}

	if _, ok := atom.(string); ok {
		return atom
	}
	if _, ok := atom.(float64); ok {
		return atom
	}
	if _, ok := atom.(bool); ok {
		return atom
	}

	if object, ok := atom.(map[string]interface{}); ok {
		if atom == nil || len(object) == 0 {
			return atom
		}

		if utils.String(object["__type"]) == "Pointer" {
			if inArray == false && inObject == false {
				return utils.String(object["className"]) + "$" + utils.String(object["objectId"])
			}
			return types.M{
				"__type":    "Pointer",
				"className": object["className"],
				"objectId":  object["objectId"],
			}
		}

		d := dateCoder{}
		if d.isValidJSON(object) {
			return d.jsonToDatabase(object)
		}
		b := bytesCoder{}
		if b.isValidJSON(object) {
			return b.jsonToDatabase(object)
		}
		g := geoPointCoder{}
		if g.isValidJSON(object) {
			if inArray || inObject {
				return object
			}
			return g.jsonToDatabase(object)
		}
		f := fileCoder{}
		if f.isValidJSON(object) {
			if inArray || inObject {
				return object
			}
			return f.jsonToDatabase(object)
		}

		if force {
			// TODO 无效类型
			return nil
		}
		return cannotTransform()
	}

	return cannotTransform()
}

func transformUpdateOperator(operator interface{}, flatten bool) (interface{}, error) {
	// TODO 处理错误
	operatorMap := utils.MapInterface(operator)
	if operatorMap == nil || operatorMap["__op"] == nil {
		return cannotTransform(), nil
	}

	op := utils.String(operatorMap["__op"])
	switch op {
	case "Delete":
		if flatten {
			return nil, nil
		}
		return types.M{
			"__op": "$unset",
			"arg":  "",
		}, nil

	case "Increment":
		if _, ok := operatorMap["amount"].(float64); !ok {
			// TODO 必须为数字
			return nil, nil
		}
		if flatten {
			return operatorMap["amount"], nil
		}
		return types.M{
			"__op": "$inc",
			"arg":  operatorMap["amount"],
		}, nil

	case "Add", "AddUnique":
		objects := utils.SliceInterface(operatorMap["objects"])
		if objects == nil {
			// TODO 必须为数组
			return nil, nil
		}
		toAdd := types.S{}
		for _, obj := range objects {
			o := transformAtom(obj, true, types.M{"inArray": true})
			toAdd = append(toAdd, o)
		}
		if flatten {
			return toAdd, nil
		}
		mongoOp := types.M{
			"Add":       "$push",
			"AddUnique": "$addToSet",
		}[op]
		return types.M{
			"__op": mongoOp,
			"arg": types.M{
				"$each": toAdd,
			},
		}, nil

	case "Remove":
		objects := utils.SliceInterface(operatorMap["objects"])
		if objects == nil {
			// TODO 必须为数组
			return nil, nil
		}
		toRemove := types.S{}
		for _, obj := range objects {
			o := transformAtom(obj, true, types.M{"inArray": true})
			toRemove = append(toRemove, o)
		}
		if flatten {
			return types.S{}, nil
		}
		return types.M{
			"__op": "$pullAll",
			"arg":  toRemove,
		}, nil

	default:
		// TODO 不支持的类型
		return nil, nil
	}
}

// transformCreate ...
func transformCreate(schema *Schema, className string, create types.M) (types.M, error) {
	// TODO 处理错误
	if className == "_User" {
		create = transformAuthData(create)
	}
	mongoCreate := transformACL(create)
	for k, v := range create {
		key, value, err := transformKeyValue(schema, className, k, v, types.M{})
		if err != nil {
			return nil, err
		}
		if value != nil {
			mongoCreate[key] = value
		}
	}
	return mongoCreate, nil
}

func transformAuthData(restObject types.M) types.M {
	if restObject["authData"] != nil {
		authData := utils.MapInterface(restObject["authData"])
		for provider, v := range authData {
			restObject["_auth_data_"+provider] = v
		}
		delete(restObject, "authData")
	}
	return restObject
}

func transformACL(restObject types.M) types.M {
	output := types.M{}
	if restObject["ACL"] == nil {
		return output
	}

	acl := utils.MapInterface(restObject["ACL"])
	rperm := types.S{}
	wperm := types.S{}
	for entry, v := range acl {
		perm := utils.MapInterface(v)
		if perm["read"] != nil {
			rperm = append(rperm, entry)
		}
		if perm["write"] != nil {
			wperm = append(wperm, entry)
		}
	}
	output["_rperm"] = rperm
	output["_wperm"] = wperm

	delete(restObject, "ACL")
	return output
}

// transformWhere 转换 where 查询数据，返回数据库格式的数据
func transformWhere(schema *Schema, className string, where types.M) (types.M, error) {
	mongoWhere := types.M{}
	if where["ACL"] != nil {
		// 不能查询 ACL
		return nil, errs.E(errs.InvalidQuery, "Cannot query on ACL.")
	}

	for k, v := range where {
		options := types.M{
			"query":    true,
			"validate": true,
		}
		key, value, err := transformKeyValue(schema, className, k, v, options)
		if err != nil {
			return nil, err
		}
		mongoWhere[key] = value
	}

	return mongoWhere, nil
}

func transformUpdate(schema *Schema, className string, update types.M) (types.M, error) {
	// TODO 处理错误
	if update == nil {
		// TODO 更新数据不能为空
		return nil, nil
	}
	if className == "_User" {
		update = transformAuthData(update)
	}

	mongoUpdate := types.M{}
	acl := transformACL(update)
	if acl["_rperm"] != nil || acl["_wperm"] != nil {
		set := types.M{}
		if acl["_rperm"] != nil {
			set["_rperm"] = acl["_rperm"]
		}
		if acl["_wperm"] != nil {
			set["_wperm"] = acl["_wperm"]
		}
		mongoUpdate["$set"] = set
	}

	for k, v := range update {
		key, value, err := transformKeyValue(schema, className, k, v, types.M{"update": true})
		if err != nil {
			return nil, err
		}

		op := utils.MapInterface(value)
		if op != nil && op["__op"] != nil {
			opKey := utils.String(op["__op"])
			opValue := types.M{}
			if mongoUpdate[opKey] != nil {
				opValue = utils.MapInterface(mongoUpdate[opKey])
			}
			opValue[key] = op["arg"]
			mongoUpdate[opKey] = opValue
		} else {
			set := types.M{}
			if mongoUpdate["$set"] != nil {
				set = utils.MapInterface(mongoUpdate["$set"])
			}
			set[key] = value
			mongoUpdate["$set"] = set
		}
	}

	return mongoUpdate, nil
}

func untransformObjectT(schema *Schema, className string, mongoObject interface{}, isNestedObject bool) interface{} {
	// TODO 处理错误
	if mongoObject == nil {
		return mongoObject
	}

	switch mongoObject.(type) {
	case string, float64, bool:
		return mongoObject

	case []interface{}:
		results := types.S{}
		objs := mongoObject.([]interface{})
		for _, o := range objs {
			results = append(results, untransformObjectT(schema, className, o, false))
		}
		return results
	}

	d := dateCoder{}
	if d.isValidDatabaseObject(mongoObject) {
		return d.databaseToJSON(mongoObject)
	}

	b := bytesCoder{}
	if b.isValidDatabaseObject(mongoObject) {
		return b.databaseToJSON(mongoObject)
	}

	if object, ok := mongoObject.(map[string]interface{}); ok {
		restObject := untransformACL(object)
		for key, value := range object {
			switch key {
			case "_id":
				restObject["objectId"] = value

			case "_hashed_password":
				restObject["password"] = value

			case "_acl", "_email_verify_token", "_perishable_token", "_tombstone":

			case "_session_token":
				restObject["sessionToken"] = value

			case "updatedAt", "_updated_at":
				restObject["updatedAt"] = utils.TimetoString(value.(time.Time))

			case "createdAt", "_created_at":
				restObject["createdAt"] = utils.TimetoString(value.(time.Time))

			case "expiresAt", "_expiresAt":
				restObject["expiresAt"] = utils.TimetoString(value.(time.Time))

			default:
				// 处理第三方登录数据
				authDataMatch, _ := regexp.MatchString(`^_auth_data_([a-zA-Z0-9_]+)$`, key)
				if authDataMatch {
					provider := key[len("_auth_data_"):]
					authData := types.M{}
					if restObject["authData"] != nil {
						authData = utils.MapInterface(restObject["authData"])
					}
					authData[provider] = value
					restObject["authData"] = authData
					break
				}

				if strings.HasPrefix(key, "_p_") {
					newKey := key[3:]
					expected := schema.getExpectedType(className, newKey)
					if expected == "" {
						// 不在 schema 中的指针类型，丢弃
						break
					}
					if expected != "" && strings.HasPrefix(expected, "*") == false {
						// schema 中对应的位置不是置身类型，丢弃
						break
					}
					if value == nil {
						break
					}
					objData := strings.Split(value.(string), "$")
					newClass := ""
					if expected != "" {
						newClass = expected[1:]
					} else {
						newClass = objData[0]
					}
					if newClass != objData[0] {
						// TODO 指向了错误的类
						return nil
					}
					restObject[newKey] = types.M{
						"__type":    "Pointer",
						"className": objData[0],
						"objectId":  objData[1],
					}
					break
				} else if isNestedObject == false && strings.HasPrefix(key, "_") && key != "__type" {
					// TODO 转换错误
					return nil
				} else {
					expectedType := schema.getExpectedType(className, key)
					f := fileCoder{}
					if expectedType == "file" && f.isValidDatabaseObject(value) {
						restObject[key] = f.databaseToJSON(value)
						break
					}
					g := geoPointCoder{}
					if expectedType == "geopoint" && g.isValidDatabaseObject(value) {
						restObject[key] = g.databaseToJSON(value)
						break
					}
				}
				restObject[key] = untransformObjectT(schema, className, value, true)
			}
		}
		return restObject
	}

	// TODO 无法转换
	return nil
}

func untransformACL(mongoObject types.M) types.M {
	output := types.M{}
	if mongoObject["_rperm"] == nil && mongoObject["_wperm"] == nil {
		return output
	}

	acl := types.M{}
	rperm := types.S{}
	wperm := types.S{}
	if mongoObject["_rperm"] != nil {
		rperm = utils.SliceInterface(mongoObject["_rperm"])
	}
	if mongoObject["_wperm"] != nil {
		wperm = utils.SliceInterface(mongoObject["_wperm"])
	}
	for _, v := range rperm {
		entry := v.(string)
		if acl[entry] == nil {
			acl[entry] = types.M{"read": true}
		} else {
			per := utils.MapInterface(acl[entry])
			per["read"] = true
			acl[entry] = per
		}
	}
	for _, v := range wperm {
		entry := v.(string)
		if acl[entry] == nil {
			acl[entry] = types.M{"write": true}
		} else {
			per := utils.MapInterface(acl[entry])
			per["write"] = true
			acl[entry] = per
		}
	}
	output["ACL"] = acl
	delete(mongoObject, "_rperm")
	delete(mongoObject, "_wperm")

	return output
}

func cannotTransform() interface{} {
	return nil
}

type dateCoder struct{}

func (d dateCoder) databaseToJSON(object interface{}) interface{} {
	data := object.(time.Time)
	json := utils.TimetoString(data)

	return types.M{
		"__type": "Date",
		"iso":    json,
	}
}

func (d dateCoder) isValidDatabaseObject(object interface{}) bool {
	if _, ok := object.(time.Time); ok {
		return true
	}
	return false
}

func (d dateCoder) jsonToDatabase(json types.M) interface{} {
	// TODO 处理错误
	t, _ := utils.StringtoTime(utils.String(json["iso"]))
	return t
}

func (d dateCoder) isValidJSON(value types.M) bool {
	return value != nil && utils.String(value["__type"]) == "Date"
}

type bytesCoder struct{}

func (b bytesCoder) databaseToJSON(object interface{}) interface{} {
	data := object.([]byte)
	json := base64.StdEncoding.EncodeToString(data)

	return types.M{
		"__type": "Bytes",
		"base64": json,
	}
}

func (b bytesCoder) isValidDatabaseObject(object interface{}) bool {
	if _, ok := object.([]byte); ok {
		return true
	}
	return false
}

func (b bytesCoder) jsonToDatabase(json types.M) interface{} {
	// TODO 处理错误
	by, _ := base64.StdEncoding.DecodeString(utils.String(json["base64"]))
	return by
}

func (b bytesCoder) isValidJSON(value types.M) bool {
	return value != nil && utils.String(value["__type"]) == "Bytes"
}

type geoPointCoder struct{}

func (g geoPointCoder) databaseToJSON(object interface{}) interface{} {
	v := object.([]interface{})
	return types.M{
		"__type":    "GeoPoint",
		"longitude": v[0],
		"latitude":  v[1],
	}
}

func (g geoPointCoder) isValidDatabaseObject(object interface{}) bool {
	if v, ok := object.([]interface{}); ok {
		if len(v) == 2 {
			return true
		}
	}
	return false
}

func (g geoPointCoder) jsonToDatabase(json types.M) interface{} {
	return types.S{json["longitude"], json["latitude"]}
}

func (g geoPointCoder) isValidJSON(value types.M) bool {
	return value != nil && utils.String(value["__type"]) == "GeoPoint"
}

type fileCoder struct{}

func (f fileCoder) databaseToJSON(object interface{}) interface{} {
	return types.M{
		"__type": "File",
		"name":   object,
	}
}

func (f fileCoder) isValidDatabaseObject(object interface{}) bool {
	if _, ok := object.(string); ok {
		return true
	}
	return false
}

func (f fileCoder) jsonToDatabase(json types.M) interface{} {
	return json["name"]
}

func (f fileCoder) isValidJSON(value types.M) bool {
	return value != nil && utils.String(value["__type"]) == "File"
}
