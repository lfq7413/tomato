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
			query, err := transformWhere(schema, className, utils.MapInterface(v), nil)
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
			query, err := transformWhere(schema, className, utils.MapInterface(v), nil)
			if err != nil {
				return "", nil, err
			}
			mongoSubqueries = append(mongoSubqueries, query)
		}
		return "$and", mongoSubqueries, nil
	default:
		// 处理第三方 auth 数据，key 的格式为： authData.xxx.id
		// {
		// 	"authData.facebook.id":"abc123"
		// }
		// ==>
		// {
		// 	"_auth_data_.facebook.id":"abc123"
		// }
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
	var expected types.M
	if schema != nil {
		expected = schema.getExpectedType(className, key)
	}

	// 期望类型为 *xxx
	// post ==> _p_post
	if expected != nil && expected["type"].(string) == "Pointer" {
		key = "_p_" + key
	}
	// 期望类型不存在，但是 restValue 中存在 "__type":"Pointer"
	if expected == nil && restValue != nil {
		op := utils.MapInterface(restValue)
		if op != nil && op["__type"] != nil {
			if utils.String(op["__type"]) == "Pointer" {
				key = "_p_" + key
			}
		}
	}

	inArray := false
	// 期望类型为 array
	if expected != nil && expected["type"].(string) == "Array" {
		inArray = true
	}

	// 处理查询操作，转换限制条件
	if options["query"] != nil {
		value, err := transformConstraint(restValue, inArray)
		if err != nil {
			return "", nil, err
		}
		if value != cannotTransform() {
			return key, value, nil
		}
	}

	// 期望类型为 array，并且转换限制条件失败，并且 restValue 不为 array 类型时，转换为 $all
	if inArray && options["query"] != nil && utils.SliceInterface(restValue) == nil {
		return key, types.M{"$all": types.S{restValue}}, nil
	}

	// 转换原子数据
	value, err := transformAtom(restValue, false, options)
	if err != nil {
		return "", nil, err
	}
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
	value, err = transformUpdateOperator(restValue, flatten)
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
func transformConstraint(constraint interface{}, inArray bool) (interface{}, error) {
	// TODO 需要根据 MongoDB 文档修正参数
	if constraint == nil || utils.MapInterface(constraint) == nil {
		return cannotTransform(), nil
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
			var err error
			answer[key], err = transformAtom(object[key], true, types.M{"inArray": inArray})
			if err != nil {
				return nil, err
			}

		// 转换 包含、不包含 操作符
		case "$in", "$nin":
			arr := utils.SliceInterface(object[key])
			if arr == nil {
				// 必须为数组
				return nil, errs.E(errs.InvalidJSON, "bad "+key+" value")
			}
			answerArr := types.S{}
			for _, v := range arr {
				obj, err := transformAtom(v, true, types.M{"inArray": inArray})
				if err != nil {
					return nil, err
				}
				answerArr = append(answerArr, obj)
			}
			answer[key] = answerArr

		// 转换 包含所有 操作符，用于数组类型的字段
		case "$all":
			arr := utils.SliceInterface(object[key])
			if arr == nil {
				// 必须为数组
				return nil, errs.E(errs.InvalidJSON, "bad "+key+" value")
			}
			answerArr := types.S{}
			for _, v := range arr {
				obj, err := transformAtom(v, true, types.M{"inArray": true})
				if err != nil {
					return nil, err
				}
				answerArr = append(answerArr, obj)
			}
			answer[key] = answerArr

		// 转换 正则 操作符
		case "$regex":
			s := utils.String(object[key])
			if s == "" {
				// 必须为字符串
				return nil, errs.E(errs.InvalidJSON, "bad regex")
			}
			answer[key] = s

		// 转换 $options 操作符
		case "$options":
			options := utils.String(object[key])
			if answer["$regex"] == nil || options == "" {
				// 无效值
				return nil, errs.E(errs.InvalidQuery, "got a bad $options")
			}
			b, _ := regexp.MatchString(`^[imxs]+$`, options)
			if b == false {
				// 无效值
				return nil, errs.E(errs.InvalidQuery, "got a bad $options")
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
			return nil, errs.E(errs.CommandUnavailable, "the "+key+" constraint is not supported yet")

		case "$within":
			within := utils.MapInterface(object[key])
			box := utils.SliceInterface(within["$box"])
			if box == nil || len(box) != 2 {
				// 参数不正确
				return nil, errs.E(errs.InvalidJSON, "malformatted $within arg")
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
				return nil, errs.E(errs.InvalidJSON, "bad constraint: "+key)
			}
			return cannotTransform(), nil
		}
	}

	return answer, nil
}

// transformAtom 转换原子数据
// options.inArray 为 true，则不进行相应转换
// options.inObject 为 true，则不进行相应转换
// force 是否强制转换，true 时如果转换失败则返回错误
func transformAtom(atom interface{}, force bool, options types.M) (interface{}, error) {
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

	// 字符串、数字、布尔类型直接返回
	if _, ok := atom.(string); ok {
		return atom, nil
	}
	if _, ok := atom.(float64); ok {
		return atom, nil
	}
	if _, ok := atom.(bool); ok {
		return atom, nil
	}

	// 转换 "__type" 声明的类型
	if object, ok := atom.(map[string]interface{}); ok {
		if atom == nil || len(object) == 0 {
			return atom, nil
		}

		// Pointer 类型
		// {
		// 	"__type": "Pointer",
		// 	"className": "abc",
		// 	"objectId": "123"
		// }
		// ==> abc$123
		if utils.String(object["__type"]) == "Pointer" {
			if inArray == false && inObject == false {
				return utils.String(object["className"]) + "$" + utils.String(object["objectId"]), nil
			}
			return types.M{
				"__type":    "Pointer",
				"className": object["className"],
				"objectId":  object["objectId"],
			}, nil
		}

		// Date 类型
		// {
		// 	"__type": "Date",
		// 	"iso": "2015-03-01T15:59:11-07:00"
		// }
		// ==> 143123456789...
		d := dateCoder{}
		if d.isValidJSON(object) {
			return d.jsonToDatabase(object)
		}

		// Bytes 类型
		// {
		// 	"__type": "Bytes",
		// 	"base64": "aGVsbG8="
		// }
		// ==> hello
		b := bytesCoder{}
		if b.isValidJSON(object) {
			return b.jsonToDatabase(object)
		}

		// GeoPoint 类型
		// {
		// 	"__type": "GeoPoint",
		//  "longitude": -30.0,
		//	"latitude": 40.0
		// }
		// ==> [-30.0, 40.0]
		g := geoPointCoder{}
		if g.isValidJSON(object) {
			if inArray || inObject {
				return object, nil
			}
			return g.jsonToDatabase(object)
		}

		// File 类型
		// {
		// 	"__type": "File",
		// 	"name": "...hello.png"
		// }
		// ==> ...hello.png
		f := fileCoder{}
		if f.isValidJSON(object) {
			if inArray || inObject {
				return object, nil
			}
			return f.jsonToDatabase(object)
		}

		// 在数组或者对象中的元素无需转换
		if inArray || inObject {
			return atom, nil
		}

		if force {
			// 无效类型，"__type" 的值不支持
			return nil, errs.E(errs.InvalidJSON, "bad atom.")
		}
		return cannotTransform(), nil
	}

	// 其他类型无法转换
	return nil, errs.E(errs.InternalServerError, "really did not expect value: atom")
}

// transformUpdateOperator 转换更新请求中的操作
// flatten 为 true 时，不再组装，直接返回实际数据
func transformUpdateOperator(operator interface{}, flatten bool) (interface{}, error) {
	// 具体操作放在 "__op" 中
	operatorMap := utils.MapInterface(operator)
	if operatorMap == nil || operatorMap["__op"] == nil {
		return cannotTransform(), nil
	}

	op := utils.String(operatorMap["__op"])
	switch op {
	// 删除字段操作
	// {
	// 	"__op":"Delete"
	// }
	// ==>
	// {
	// 	"__op": "$unset",
	// 	"arg":  ""
	// }
	case "Delete":
		if flatten {
			return nil, nil
		}
		return types.M{
			"__op": "$unset",
			"arg":  "",
		}, nil

	// 数值增加操作
	// {
	// 	"__op":"Increment",
	// 	"amount":10
	// }
	// ==>
	// {
	// 	"__op": "$inc",
	// 	"arg":10
	// }
	case "Increment":
		if _, ok := operatorMap["amount"].(float64); !ok {
			// 必须为数字
			return nil, errs.E(errs.InvalidJSON, "incrementing must provide a number")
		}
		if flatten {
			return operatorMap["amount"], nil
		}
		return types.M{
			"__op": "$inc",
			"arg":  operatorMap["amount"],
		}, nil

	// 增加对象操作
	// {
	// 	"__op":"Add"
	// 	"objects":[{...},{...}]
	// }
	// ==>
	// {
	// 	"__op":"$push",
	// 	"arg":{
	// 		"$each":[{...},{...}]
	// 	}
	// }
	case "Add", "AddUnique":
		objects := utils.SliceInterface(operatorMap["objects"])
		if objects == nil {
			// 必须为数组
			return nil, errs.E(errs.InvalidJSON, "objects to add must be an array")
		}
		toAdd := types.S{}
		for _, obj := range objects {
			o, err := transformAtom(obj, true, types.M{"inArray": true})
			if err != nil {
				return nil, err
			}
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

	// 删除对象操作
	// {
	// 	"__op":"Remove",
	// 	"objects":[{...},{...}]
	// }
	// ==>
	// {
	// 	"__op": "$pullAll",
	// 	"arg":[{...},{...}]
	// }
	case "Remove":
		objects := utils.SliceInterface(operatorMap["objects"])
		if objects == nil {
			// 必须为数组
			return nil, errs.E(errs.InvalidJSON, "objects to remove must be an array")
		}
		toRemove := types.S{}
		for _, obj := range objects {
			o, err := transformAtom(obj, true, types.M{"inArray": true})
			if err != nil {
				return nil, err
			}
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
		// 不支持的类型
		return nil, errs.E(errs.CommandUnavailable, "the "+op+" op is not supported yet")
	}
}

// transformCreate 转换 create 数据
func transformCreate(schema *Schema, className string, create types.M) (types.M, error) {
	// 转换第三方登录数据
	if className == "_User" {
		create = transformAuthData(create)
	}
	// 转换权限数据，转换完成之后仅包含权限信息
	mongoCreate := transformACL(create)

	// 转换其他字段并添加
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
func transformAuthData(restObject types.M) types.M {
	if restObject["authData"] != nil {
		authData := utils.MapInterface(restObject["authData"])
		for provider, v := range authData {
			if v == nil || utils.MapInterface(v) == nil || len(utils.MapInterface(v)) == 0 {
				restObject["_auth_data_"+provider] = types.M{
					"__op": "Delete",
				}
			} else {
				restObject["_auth_data_"+provider] = v
			}
		}
		delete(restObject, "authData")
	}
	return restObject
}

// transformACL 转换生成权限信息，并删除源数据中的 ACL 字段
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
// 	"_wperm":["userid","role:xxx"]
// }
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
func transformWhere(schema *Schema, className string, where types.M, options types.M) (types.M, error) {
	if options == nil || len(options) == 0 {
		options = types.M{"validate": true}
	}
	mongoWhere := types.M{}
	if where["ACL"] != nil {
		// 不能查询 ACL
		return nil, errs.E(errs.InvalidQuery, "Cannot query on ACL.")
	}

	transformKeyOptions := types.M{
		"query":    true,
		"validate": options["validate"],
	}
	for k, v := range where {
		key, value, err := transformKeyValue(schema, className, k, v, transformKeyOptions)
		if err != nil {
			return nil, err
		}
		mongoWhere[key] = value
	}

	return mongoWhere, nil
}

// transformUpdate 转换 update 数据
func transformUpdate(schema *Schema, className string, update types.M) (types.M, error) {
	if update == nil {
		// 更新数据不能为空
		return nil, errs.E(errs.InvalidJSON, "got empty restUpdate")
	}
	// 处理第三方登录数据
	if className == "_User" {
		update = transformAuthData(update)
	}

	mongoUpdate := types.M{}
	// 转换并设置权限信息
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

	// 转换 update 中的其他数据
	for k, v := range update {
		key, value, err := transformKeyValue(schema, className, k, v, types.M{"update": true})
		if err != nil {
			return nil, err
		}

		op := utils.MapInterface(value)
		if op != nil && op["__op"] != nil {
			// 处理带 "__op" 的数据，如下：
			// {
			// 	"size":{
			// 		"__op":"$inc",
			// 		"arg":3
			// 	}
			// }
			// ==>
			// {
			// 	"$inc":{
			// 		"size",3
			// 	}
			// }
			opKey := utils.String(op["__op"])
			opValue := types.M{}
			if mongoUpdate[opKey] != nil {
				opValue = utils.MapInterface(mongoUpdate[opKey])
			}
			opValue[key] = op["arg"]
			mongoUpdate[opKey] = opValue
		} else {
			// 转换其他数据
			// {
			// 	"name":"joe"
			// }
			// ==>
			// {
			// 	"$set":{
			// 		"name":"joe"
			// 	}
			// }
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

var specialKeysForUntransform = []string{
	"_id",
	"_hashed_password",
	"_acl",
	"_email_verify_token",
	"_perishable_token",
	"_tombstone",
	"_session_token",
	"updatedAt",
	"_updated_at",
	"createdAt",
	"_created_at",
	"expiresAt",
	"_expiresAt",
}

// untransformObjectT  把数据库类型数据转换为 API 格式
func untransformObjectT(schema *Schema, className string, mongoObject interface{}, isNestedObject bool) (interface{}, error) {
	if mongoObject == nil {
		return mongoObject, nil
	}

	// 转换基本类型
	switch mongoObject.(type) {
	case string, float64, bool:
		return mongoObject, nil

	case []interface{}:
		results := types.S{}
		objs := mongoObject.([]interface{})
		for _, o := range objs {
			res, err := untransformObjectT(schema, className, o, true)
			if err != nil {
				return nil, err
			}
			results = append(results, res)
		}
		return results, nil
	}

	// 日期格式
	// {
	// 	"__type": "Date",
	// 	"iso": "2015-03-01T15:59:11-07:00"
	// }
	d := dateCoder{}
	if d.isValidDatabaseObject(mongoObject) {
		return d.databaseToJSON(mongoObject), nil
	}

	// byte 数组
	// {
	// 	"__type": "Bytes",
	// 	"base64": "aGVsbG8="
	// }
	b := bytesCoder{}
	if b.isValidDatabaseObject(mongoObject) {
		return b.databaseToJSON(mongoObject), nil
	}

	// 转换对象类型
	if object, ok := mongoObject.(map[string]interface{}); ok {
		// 转换权限信息
		restObject := untransformACL(object)
		for key, value := range object {
			if isNestedObject {
				in := false
				for _, v := range specialKeysForUntransform {
					if key == v {
						in = true
						break
					}
				}
				if in {
					r, err := untransformObjectT(schema, className, value, true)
					if err != nil {
						return nil, err
					}
					restObject[key] = r
					continue
				}
			}
			switch key {
			case "_id":
				restObject["objectId"] = value

			case "_hashed_password":
				restObject["password"] = value

			case "_acl", "_email_verify_token", "_perishable_token", "_tombstone":
				// 以上字段不添加到结果中

			case "_session_token":
				restObject["sessionToken"] = value

			// 时间类型转换为 ISO8601 标准的字符串
			case "updatedAt", "_updated_at":
				restObject["updatedAt"] = utils.TimetoString(value.(time.Time))

			case "createdAt", "_created_at":
				restObject["createdAt"] = utils.TimetoString(value.(time.Time))

			case "expiresAt", "_expiresAt":
				restObject["expiresAt"] = utils.TimetoString(value.(time.Time))

			default:
				// 处理第三方登录数据
				// {
				// 	"_auth_data_facebook":{...}
				// }
				// ==>
				// {
				// 	"authData":{
				// 		"facebook":{...}
				// 	}
				// }
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

				// 处理指针类型的字段
				// {
				// 	"_p_post":"abc$123"
				// }
				// ==>
				// {
				// 	"__type":    "Pointer",
				// 	"className": "post",
				// 	"objectId":  "123"
				// }
				if strings.HasPrefix(key, "_p_") {
					newKey := key[3:]
					expected := schema.getExpectedType(className, newKey)
					if expected == nil {
						// 不在 schema 中的指针类型，丢弃
						break
					}
					if expected != nil && expected["type"].(string) != "Pointer" {
						// schema 中对应的位置不是指针类型，丢弃
						break
					}
					if value == nil {
						break
					}
					objData := strings.Split(value.(string), "$")
					newClass := ""
					if expected != nil {
						newClass = expected["targetClass"].(string)
					} else {
						newClass = objData[0]
					}
					if newClass != objData[0] {
						// 指向了错误的类
						return nil, errs.E(errs.InternalServerError, "pointer to incorrect className")
					}
					restObject[newKey] = types.M{
						"__type":    "Pointer",
						"className": objData[0],
						"objectId":  objData[1],
					}
					break
				} else if isNestedObject == false && strings.HasPrefix(key, "_") && key != "__type" {
					// 转换错误
					return nil, errs.E(errs.InternalServerError, "bad key in untransform: "+key)
				} else {
					// TODO 此处可能会有问题，isNestedObject == true 时，即子对象也会进来
					// 但是拿子对象的 key 无法从 className 中查询有效的类型
					// 所以当子对象的某个 key 与 className 中的某个 key 相同时，可能出问题
					expectedType := schema.getExpectedType(className, key)
					// file 类型
					// {
					// 	"__type": "File",
					// 	"name":   "hello.jpg"
					// }
					f := fileCoder{}
					if expectedType != nil && expectedType["type"].(string) == "File" && f.isValidDatabaseObject(value) {
						restObject[key] = f.databaseToJSON(value)
						break
					}
					// geopoint 类型
					// {
					// 	"__type":    "GeoPoint",
					// 	"longitude": 30,
					// 	"latitude":  40
					// }
					g := geoPointCoder{}
					if expectedType != nil && expectedType["type"].(string) == "geopoint" && g.isValidDatabaseObject(value) {
						restObject[key] = g.databaseToJSON(value)
						break
					}
				}
				// 转换子对象
				res, err := untransformObjectT(schema, className, value, true)
				if err != nil {
					return nil, err
				}
				restObject[key] = res
			}
		}

		if isNestedObject == false {
			relationFields := schema.getRelationFields(className)
			for k, v := range relationFields {
				restObject[k] = v
			}
		}

		return restObject, nil
	}

	// 无法转换
	return nil, errs.E(errs.InternalServerError, "unknown object type")
}

// untransformACL 把数据库格式的权限信息转换为 API 格式
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

// transformSelect 转换对象中的 $select
func transformSelect(selectObject types.M, key string, objects []types.M) {
	values := []interface{}{}
	for _, result := range objects {
		values = append(values, result[key])
	}

	delete(selectObject, "$select")
	var in []interface{}
	if v, ok := selectObject["$in"].([]interface{}); ok {
		in = v
		in = append(in, values...)
	} else {
		in = values
	}
	selectObject["$in"] = in
}

// transformDontSelect 转换对象中的 $dontSelect
func transformDontSelect(dontSelectObject types.M, key string, objects []types.M) {
	values := []interface{}{}
	for _, result := range objects {
		values = append(values, result[key])
	}

	delete(dontSelectObject, "$dontSelect")
	var nin []interface{}
	if v, ok := dontSelectObject["$nin"].([]interface{}); ok {
		nin = v
		nin = append(nin, values...)
	} else {
		nin = values
	}
	dontSelectObject["$nin"] = nin
}

// transformInQuery 转换对象中的 $inQuery
func transformInQuery(inQueryObject types.M, className string, results []types.M) {
	values := []interface{}{}
	for _, result := range results {
		o := types.M{
			"__type":    "Pointer",
			"className": className,
			"objectId":  result["objectId"],
		}
		values = append(values, o)
	}

	delete(inQueryObject, "$inQuery")
	var in []interface{}
	if v, ok := inQueryObject["$in"].([]interface{}); ok {
		in = v
		in = append(in, values...)
	} else {
		in = values
	}
	inQueryObject["$in"] = in
}

func cannotTransform() interface{} {
	return nil
}

// dateCoder Date 类型数据处理
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

func (d dateCoder) jsonToDatabase(json types.M) (interface{}, error) {
	t, err := utils.StringtoTime(utils.String(json["iso"]))
	if err != nil {
		return nil, errs.E(errs.InvalidJSON, "invalid iso")
	}
	return t, nil
}

func (d dateCoder) isValidJSON(value types.M) bool {
	return value != nil && utils.String(value["__type"]) == "Date" && utils.String(value["iso"]) != ""
}

// bytesCoder Bytes 类型处理
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

func (b bytesCoder) jsonToDatabase(json types.M) (interface{}, error) {
	by, err := base64.StdEncoding.DecodeString(utils.String(json["base64"]))
	if err != nil {
		return nil, errs.E(errs.InvalidJSON, "invalid base64")
	}
	return by, nil
}

func (b bytesCoder) isValidJSON(value types.M) bool {
	return value != nil && utils.String(value["__type"]) == "Bytes" && utils.String(value["base64"]) != ""
}

// geoPointCoder GeoPoint 类型处理
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

func (g geoPointCoder) jsonToDatabase(json types.M) (interface{}, error) {
	if _, ok := json["longitude"].(float64); ok == false {
		return nil, errs.E(errs.InvalidJSON, "invalid longitude")
	}
	if _, ok := json["latitude"].(float64); ok == false {
		return nil, errs.E(errs.InvalidJSON, "invalid latitude")
	}
	return types.S{json["longitude"], json["latitude"]}, nil
}

func (g geoPointCoder) isValidJSON(value types.M) bool {
	return value != nil && utils.String(value["__type"]) == "GeoPoint" && value["longitude"] != nil && value["latitude"] != nil
}

// fileCoder File 类型处理
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

func (f fileCoder) jsonToDatabase(json types.M) (interface{}, error) {
	return json["name"], nil
}

func (f fileCoder) isValidJSON(value types.M) bool {
	return value != nil && utils.String(value["__type"]) == "File" && utils.String(value["name"]) != ""
}
