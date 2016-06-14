package mongo

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

// Transform ...
type Transform struct{}

// NewTransform ...
func NewTransform() *Transform {
	return &Transform{}
}

// transformKey ...
func (t *Transform) transformKey(className, fieldName string, schema types.M) string {
	// TODO className 没有用到
	switch fieldName {
	case "objectId":
		return "_id"

	case "createdAt":
		return "_created_at"

	case "updatedAt":
		return "_updated_at"

	case "sessionToken":
		return "_session_token"

	}

	if schema == nil {
		return fieldName
	}

	if fields := utils.M(schema["fields"]); fields != nil {
		if tp := utils.M(fields[fieldName]); tp != nil {
			if utils.S(tp["__type"]) == "Pointer" {
				fieldName = "_p_" + fieldName
			}
		}
	}

	return fieldName
}

// transformKeyValueForUpdate 把传入的键值对转换为数据库中保存的格式
// restKey API 格式的字段名
// restValue API 格式的值
// 返回转换成 数据库格式的字段名与值
func (t *Transform) transformKeyValueForUpdate(className, restKey string, restValue interface{}, parseFormatSchema types.M) (string, interface{}, error) {
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
	case "sessionToken", "_session_token":
		key = "_session_token"
	case "expiresAt", "_expiresAt":
		key = "_expiresAt"
		timeField = true
	case "_rperm", "_wperm":
		return key, restValue, nil

		// 其他字段名不做处理
	}

	// 处理特殊字段名
	var expected types.M
	if parseFormatSchema != nil {
		if fields := utils.M(parseFormatSchema["fields"]); fields != nil {
			expected = utils.M(fields[key])
		}
	}

	if expected != nil && utils.S(expected["type"]) == "Pointer" {
		// 期望类型为 *xxx
		// post ==> _p_post
		key = "_p_" + key
	} else if expected == nil {
		// 期望类型不存在，但是 restValue 中存在 "__type":"Pointer"
		if op := utils.M(restValue); op != nil {
			if utils.S(op["__type"]) == "Pointer" {
				key = "_p_" + key
			}
		}
	}

	if restValue == nil {
		return key, nil, nil
	}

	// 转换原子数据
	value, err := t.transformTopLevelAtom(restValue)
	if err != nil {
		return "", nil, err
	}
	if value != cannotTransform() {
		if timeField && utils.S(value) != "" {
			var err error
			value, err = utils.StringtoTime(utils.S(value))
			if err != nil {
				return "", nil, errs.E(errs.InvalidJSON, "Invalid Date value.")
			}
		}
		return key, value, nil
	}

	// 转换数组类型
	if valueArray := utils.A(restValue); valueArray != nil {
		outValue := types.S{}
		for _, restObj := range valueArray {
			v, err := t.transformInteriorValue(restObj)
			if err != nil {
				return "", nil, err
			}
			outValue = append(outValue, v)
		}
		return key, outValue, nil
	}

	// 处理更新操作中的 "_op"
	if value := utils.M(restValue); value != nil {
		if _, ok := value["__op"]; ok {
			v, err := t.transformUpdateOperator(restValue, false)
			if err != nil {
				return "", nil, err
			}
			return key, v, nil
		}
	}

	// 处理正常的对象
	if value := utils.M(restValue); value != nil {
		newValue := types.M{}
		for k, v := range value {
			r, err := t.transformInteriorValue(v)
			if err != nil {
				return "", nil, err
			}
			newValue[k] = r
		}
		return key, newValue, nil
	}

	return key, restValue, nil
}

// valueAsDate 校验并转换时间类型
func valueAsDate(value interface{}) (time.Time, bool) {
	if s, ok := value.(string); ok {
		t, err := utils.StringtoTime(s)
		if err != nil {
			return time.Time{}, false
		}
		return t, true
	}
	if t, ok := value.(time.Time); ok {
		return t, true
	}
	return time.Time{}, false
}

// transformQueryKeyValue 转换查询请求中的键值对
func (t *Transform) transformQueryKeyValue(className, key string, value interface{}, schema types.M) (string, interface{}, error) {
	// TODO className 没有用到
	switch key {
	case "createdAt":
		if t, ok := valueAsDate(value); ok {
			return "_created_at", t, nil
		}
		key = "_created_at"

	case "updatedAt":
		if t, ok := valueAsDate(value); ok {
			return "_updated_at", t, nil
		}
		key = "_updated_at"

	case "expiresAt":
		if t, ok := valueAsDate(value); ok {
			return "_expiresAt", t, nil
		}
		key = "_expiresAt"

	case "objectId":
		return "_id", value, nil

	case "sessionToken":
		return "_session_token", value, nil

	case "_rperm", "_wperm", "_perishable_token", "_email_verify_token":
		return key, value, nil

	case "$or":
		if value == nil {
			return "$or", nil, nil
		}
		array := utils.A(value)
		if array == nil {
			return "$or", nil, nil
		}
		querys := types.S{}
		for _, subQuery := range array {
			sub := utils.M(subQuery)
			if sub == nil {
				continue
			}
			r, err := t.transformWhere(className, sub, schema)
			if err != nil {
				return "", nil, err
			}
			querys = append(querys, r)
		}
		return "$or", querys, nil

	case "$and":
		if value == nil {
			return "$and", nil, nil
		}
		array := utils.A(value)
		if array == nil {
			return "$and", nil, nil
		}
		querys := types.S{}
		for _, subQuery := range array {
			sub := utils.M(subQuery)
			if sub == nil {
				continue
			}
			r, err := t.transformWhere(className, sub, schema)
			if err != nil {
				return "", nil, err
			}
			querys = append(querys, r)
		}
		return "$and", querys, nil

	default:
		authDataMatch, _ := regexp.MatchString(`^authData\.([a-zA-Z0-9_]+)\.id$`, key)
		if authDataMatch {
			provider := key[len("authData."):(len(key) - len(".id"))]
			return "_auth_data_" + provider + ".id", value, nil
		}

	}

	var fieldType types.M
	if schema != nil {
		if fields := utils.M(schema["fields"]); fields != nil {
			fieldType = utils.M(fields[key])
		}
	}

	expectedTypeIsArray := schema != nil && fieldType != nil && utils.S(fieldType["type"]) == "Array"
	expectedTypeIsPointer := schema != nil && fieldType != nil && utils.S(fieldType["type"]) == "Pointer"

	if expectedTypeIsPointer {
		key = "_p_" + key
	} else if fieldType == nil {
		if v := utils.M(value); v != nil {
			if utils.S(v["__type"]) == "Pointer" {
				key = "_p_" + key
			}
		}
	}

	cValue, err := t.transformConstraint(value, expectedTypeIsArray)
	if err != nil {
		return "", nil, err
	}
	if cValue != cannotTransform() {
		return key, cValue, nil
	}

	if v := utils.A(value); v == nil && expectedTypeIsArray {
		return key, types.M{"$all": types.S{value}}, nil
	}

	aValue, err := t.transformTopLevelAtom(value)
	if err != nil {
		return "", nil, err
	}
	if aValue != cannotTransform() {
		return key, aValue, nil
	}
	return "", nil, errs.E(errs.InvalidJSON, "You cannot use this value as a query parameter.")
}

// transformConstraint 转换查询限制条件，处理的操作符类似 "$lt", "$gt" 等
// inArray 表示该字段是否为数组类型
func (t *Transform) transformConstraint(constraint interface{}, inArray bool) (interface{}, error) {
	// TODO 需要根据 MongoDB 文档修正参数
	if constraint == nil || utils.M(constraint) == nil {
		return cannotTransform(), nil
	}

	// keys is the constraints in reverse alphabetical order.
	// This is a hack so that:
	//   $regex is handled before $options
	//   $nearSphere is handled before $maxDistance
	object := utils.M(constraint)
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
			if inArray {
				answer[key], err = t.transformInteriorAtom(object[key])
				if err != nil {
					return nil, err
				}
			} else {
				answer[key], err = t.transformTopLevelAtom(object[key])
				if err != nil {
					return nil, err
				}
			}
			if answer[key] == cannotTransform() {
				return nil, errs.E(errs.InvalidJSON, "bad atom")
			}

		// 转换 包含、不包含 操作符
		case "$in", "$nin":
			arr := utils.A(object[key])
			if arr == nil {
				// 必须为数组
				return nil, errs.E(errs.InvalidJSON, "bad "+key+" value")
			}
			answerArr := types.S{}
			for _, value := range arr {
				var result interface{}
				var err error
				if inArray {
					result, err = t.transformInteriorAtom(value)
					if err != nil {
						return nil, err
					}
				} else {
					result, err = t.transformTopLevelAtom(value)
					if err != nil {
						return nil, err
					}
				}
				if result == cannotTransform() {
					return nil, errs.E(errs.InvalidJSON, "bad atom")
				}
				answerArr = append(answerArr, result)
			}
			answer[key] = answerArr

		// 转换 包含所有 操作符，用于数组类型的字段
		case "$all":
			arr := utils.A(object[key])
			if arr == nil {
				// 必须为数组
				return nil, errs.E(errs.InvalidJSON, "bad "+key+" value")
			}
			answerArr := types.S{}
			for _, v := range arr {
				obj, err := t.transformInteriorAtom(v)
				if err != nil {
					return nil, err
				}
				answerArr = append(answerArr, obj)
			}
			answer[key] = answerArr

		// 转换 正则 操作符
		case "$regex":
			s := utils.S(object[key])
			if s == "" {
				// 必须为字符串
				return nil, errs.E(errs.InvalidJSON, "bad regex")
			}
			answer[key] = s

		// 转换 $options 操作符
		case "$options":
			options := utils.S(object[key])
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
			point := utils.M(object[key])
			if point == nil {
				answer[key] = types.M{
					"$geometry": types.M{
						"type":        "Point",
						"coordinates": types.S{0, 0},
					},
					"$maxDistance": 0,
				}
			} else {
				answer[key] = types.M{
					"$geometry": types.M{
						"type":        "Point",
						"coordinates": types.S{point["longitude"], point["latitude"]},
					},
				}
			}

		// 转换 最大距离 操作符，单位是弧度，需要转换为 米
		// 在 MongoDB 2.6 之后，操作符 $nearSphere 最新的使用格式中，$maxDistance 的单位是米
		// $maxDistance 仅能用于 2dsphere 格式的索引
		case "$maxDistance", "$maxDistanceInRadians":
			var nearSphere types.M
			if nearSphere = utils.M(answer["$nearSphere"]); nearSphere == nil {
				break
			}
			var distance float64
			if v, ok := object[key].(float64); ok {
				distance = v * 6371 * 1000
			} else if v, ok := object[key].(int); ok {
				distance = float64(v) * 6371 * 1000
			}
			nearSphere["$maxDistance"] = distance
			answer["$nearSphere"] = nearSphere

		// 以下三项在 SDK 中未使用，但是在 REST API 中使用了
		case "$maxDistanceInMiles":
			var nearSphere types.M
			if nearSphere = utils.M(answer["$nearSphere"]); nearSphere == nil {
				break
			}
			var distance float64
			if v, ok := object[key].(float64); ok {
				distance = v * 1.609344 * 1000
			} else if v, ok := object[key].(int); ok {
				distance = float64(v) * 1.609344 * 1000
			}
			nearSphere["$maxDistance"] = distance
			answer["$nearSphere"] = nearSphere

		case "$maxDistanceInKilometers":
			var nearSphere types.M
			if nearSphere = utils.M(answer["$nearSphere"]); nearSphere == nil {
				break
			}
			var distance float64
			if v, ok := object[key].(float64); ok {
				distance = v * 1000
			} else if v, ok := object[key].(int); ok {
				distance = float64(v) * 1000
			}
			nearSphere["$maxDistance"] = distance
			answer["$nearSphere"] = nearSphere

		case "$select", "$dontSelect":
			// 暂时不支持该参数
			return nil, errs.E(errs.CommandUnavailable, "the "+key+" constraint is not supported yet")

		case "$within":
			within := utils.M(object[key])
			if within == nil {
				return nil, errs.E(errs.InvalidJSON, "malformatted $within arg")
			}
			box := utils.A(within["$box"])
			if box == nil || len(box) != 2 {
				// 参数不正确
				return nil, errs.E(errs.InvalidJSON, "malformatted $within arg")
			}
			box1 := utils.M(box[0])
			box2 := utils.M(box[1])
			if box1 == nil || box2 == nil {
				return nil, errs.E(errs.InvalidJSON, "malformatted $within arg")
			}
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

// transformTopLevelAtom 转换顶层的原子数据
func (t *Transform) transformTopLevelAtom(atom interface{}) (interface{}, error) {
	if atom == nil {
		return atom, nil
	}
	// 字符串、数字、布尔类型直接返回
	if _, ok := atom.(string); ok {
		return atom, nil
	}
	if _, ok := atom.(float64); ok {
		return atom, nil
	}
	if _, ok := atom.(int); ok {
		return atom, nil
	}
	if _, ok := atom.(bool); ok {
		return atom, nil
	}
	if _, ok := atom.(time.Time); ok {
		return atom, nil
	}

	// 转换 "__type" 声明的类型
	if object := utils.M(atom); object != nil {
		if len(object) == 0 {
			return atom, nil
		}

		// Pointer 类型
		// {
		// 	"__type": "Pointer",
		// 	"className": "abc",
		// 	"objectId": "123"
		// }
		// ==> abc$123
		if utils.S(object["__type"]) == "Pointer" {
			return utils.S(object["className"]) + "$" + utils.S(object["objectId"]), nil
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
			return f.jsonToDatabase(object)
		}

		return cannotTransform(), nil
	}

	// 数组类型不转换
	if object := utils.A(atom); object != nil {
		return cannotTransform(), nil
	}

	// 其他类型无法转换
	return cannotTransform(), nil
}

// transformUpdateOperator 转换更新请求中的操作
// flatten 为 true 时，不再组装，直接返回实际数据
func (t *Transform) transformUpdateOperator(operator interface{}, flatten bool) (interface{}, error) {
	// 具体操作放在 "__op" 中
	operatorMap := utils.M(operator)
	if operatorMap == nil || operatorMap["__op"] == nil {
		return operator, nil
	}

	op := utils.S(operatorMap["__op"])
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
		number := false
		if _, ok := operatorMap["amount"].(float64); ok {
			number = true
		}
		if _, ok := operatorMap["amount"].(int); ok {
			number = true
		}
		if number == false {
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
		objects := utils.A(operatorMap["objects"])
		if objects == nil {
			// 必须为数组
			return nil, errs.E(errs.InvalidJSON, "objects to add must be an array")
		}
		toAdd := types.S{}
		for _, obj := range objects {
			o, err := t.transformInteriorAtom(obj)
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
		objects := utils.A(operatorMap["objects"])
		if objects == nil {
			// 必须为数组
			return nil, errs.E(errs.InvalidJSON, "objects to remove must be an array")
		}
		toRemove := types.S{}
		for _, obj := range objects {
			o, err := t.transformInteriorAtom(obj)
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
		return nil, errs.E(errs.CommandUnavailable, "the "+op+" operator is not supported yet")
	}
}

// parseObjectToMongoObjectForCreate 转换 create 数据
func (t *Transform) parseObjectToMongoObjectForCreate(className string, create types.M, schema types.M) (types.M, error) {
	if create == nil {
		return nil, nil
	}
	// 转换第三方登录数据
	if className == "_User" {
		create = t.transformAuthData(create)
	}
	// 转换权限数据，转换完成之后仅包含权限信息
	mongoCreate := t.transformACL(create)

	// 转换其他字段并添加
	for k, v := range create {
		key, value, err := t.parseObjectKeyValueToMongoObjectKeyValue(className, k, v, schema)
		if err != nil {
			return nil, err
		}
		if value != nil {
			mongoCreate[key] = value
		}
	}
	return mongoCreate, nil
}

func (t *Transform) parseObjectKeyValueToMongoObjectKeyValue(className string, restKey string, restValue interface{}, schema types.M) (string, interface{}, error) {
	// TODO className 没有用到
	var transformedValue interface{}
	var coercedToDate interface{}
	var err error
	switch restKey {
	case "objectId":
		return "_id", restValue, nil

	case "createdAt":
		transformedValue, err = t.transformTopLevelAtom(restValue)
		if err != nil {
			return "", nil, err
		}
		if v, ok := transformedValue.(string); ok {
			coercedToDate, err = utils.StringtoTime(v)
			if err != nil {
				return "", nil, err
			}
		} else {
			coercedToDate = transformedValue
		}
		return "_created_at", coercedToDate, nil

	case "updatedAt":
		transformedValue, err = t.transformTopLevelAtom(restValue)
		if err != nil {
			return "", nil, err
		}
		if v, ok := transformedValue.(string); ok {
			coercedToDate, err = utils.StringtoTime(v)
			if err != nil {
				return "", nil, err
			}
		} else {
			coercedToDate = transformedValue
		}
		return "_updated_at", coercedToDate, nil

	case "expiresAt":
		transformedValue, err = t.transformTopLevelAtom(restValue)
		if err != nil {
			return "", nil, err
		}
		if v, ok := transformedValue.(string); ok {
			coercedToDate, err = utils.StringtoTime(v)
			if err != nil {
				return "", nil, err
			}
		} else {
			coercedToDate = transformedValue
		}
		return "_expiresAt", coercedToDate, nil

	case "_rperm", "_wperm", "_email_verify_token", "_hashed_password", "_perishable_token":
		return restKey, restValue, nil

	case "sessionToken":
		return "_session_token", restValue, nil

	default:
		m, _ := regexp.MatchString(`^authData\.([a-zA-Z0-9_]+)\.id$`, restKey)
		if m {
			return "", nil, errs.E(errs.InvalidKeyName, "can only query on "+restKey)
		}
		m, _ = regexp.MatchString(`^_auth_data_[a-zA-Z0-9_]+$`, restKey)
		if m {
			return restKey, restValue, nil
		}
	}

	if restValue == nil {
		return restKey, restValue, nil
	}

	if v := utils.M(restValue); v != nil {
		if ty, ok := v["__type"]; ok {
			if ty.(string) != "Bytes" {
				if ty.(string) == "Pointer" {
					restKey = "_p_" + restKey
				} else if schema != nil {
					if fields := utils.M(schema["fields"]); fields != nil {
						if t := utils.M(fields[restKey]); t != nil {
							if utils.S(t["type"]) == "Pointer" {
								restKey = "_p_" + restKey
							}
						}
					}
				}
			}
		}
	}

	if restKey == "ACL" {
		return "", nil, errs.E(errs.InvalidKeyName, "There was a problem transforming an ACL.")
	}

	value, err := t.transformTopLevelAtom(restValue)
	if err != nil {
		return "", nil, err
	}
	if value != cannotTransform() {
		return restKey, value, nil
	}

	if s := utils.A(restValue); s != nil {
		value := types.S{}
		for _, restObj := range s {
			v, err := t.transformInteriorValue(restObj)
			if err != nil {
				return "", nil, err
			}
			value = append(value, v)
		}
		return restKey, value, nil
	}

	// 处理更新操作中的 "_op"
	if value := utils.M(restValue); value != nil {
		if _, ok := value["__op"]; ok {
			v, err := t.transformUpdateOperator(restValue, false)
			if err != nil {
				return "", nil, err
			}
			return restKey, v, nil
		}
	}

	// 处理正常的对象
	if value := utils.M(restValue); value != nil {
		for k := range value {
			if strings.Index(k, "$") > -1 || strings.Index(k, ".") > -1 {
				return "", nil, errs.E(errs.InvalidNestedKey, "Nested keys should not contain the '$' or '.' characters")
			}
		}
	}

	if value := utils.M(restValue); value != nil {
		newValue := types.M{}
		for k, v := range value {
			r, err := t.transformInteriorValue(v)
			if err != nil {
				return "", nil, err
			}
			newValue[k] = r
		}
		return restKey, newValue, nil
	}

	return restKey, restValue, nil
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
func (t *Transform) transformAuthData(restObject types.M) types.M {
	if restObject != nil && restObject["authData"] != nil {
		authData := utils.M(restObject["authData"])
		if authData != nil {
			for provider, v := range authData {
				if v == nil || utils.M(v) == nil || len(utils.M(v)) == 0 {
					restObject["_auth_data_"+provider] = types.M{
						"__op": "Delete",
					}
				} else {
					restObject["_auth_data_"+provider] = v
				}
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
// 	"_wperm":["userid","role:xxx"],
// 	"_acl":{
// 		"userid":{
// 			"r": true,
// 			"w": true,
// 		},
// 		"role:xxx":{
// 			"r": true,
// 			"w": true,
// 		},
// 		"*":{
// 			"r": true,
// 		},
// 	},
// }
func (t *Transform) transformACL(restObject types.M) types.M {
	output := types.M{}
	if restObject == nil || restObject["ACL"] == nil {
		return output
	}

	acl := utils.M(restObject["ACL"])
	rperm := types.S{}
	wperm := types.S{}
	_acl := types.M{}
	for entry, v := range acl {
		perm := utils.M(v)
		if perm != nil {
			a := types.M{}
			if perm["read"] != nil {
				rperm = append(rperm, entry)
				a["r"] = true
			}
			if perm["write"] != nil {
				wperm = append(wperm, entry)
				a["w"] = true
			}
			_acl[entry] = a
		}
	}
	output["_rperm"] = rperm
	output["_wperm"] = wperm
	output["_acl"] = _acl

	delete(restObject, "ACL")
	return output
}

// transformWhere 转换 where 查询数据，返回数据库格式的数据
func (t *Transform) transformWhere(className string, where, schema types.M) (types.M, error) {
	if where == nil {
		return nil, nil
	}
	// TODO className 没有用到
	mongoWhere := types.M{}
	for k, v := range where {
		key, value, err := t.transformQueryKeyValue(className, k, v, schema)
		if err != nil {
			return nil, err
		}
		mongoWhere[key] = value
	}

	return mongoWhere, nil
}

// transformUpdate 转换 update 数据
func (t *Transform) transformUpdate(className string, update types.M, parseFormatSchema types.M) (types.M, error) {
	if update == nil {
		return nil, nil
	}
	// 处理第三方登录数据
	if className == "_User" {
		update = t.transformAuthData(update)
	}

	mongoUpdate := types.M{}
	// 转换并设置权限信息
	acl := t.transformACL(update)
	if acl["_rperm"] != nil || acl["_wperm"] != nil || acl["_acl"] != nil {
		set := types.M{}
		if acl["_rperm"] != nil {
			set["_rperm"] = acl["_rperm"]
		}
		if acl["_wperm"] != nil {
			set["_wperm"] = acl["_wperm"]
		}
		if acl["_acl"] != nil {
			set["_acl"] = acl["_acl"]
		}
		mongoUpdate["$set"] = set
	}

	// 转换 update 中的其他数据
	for k, v := range update {
		key, value, err := t.transformKeyValueForUpdate(className, k, v, parseFormatSchema)
		if err != nil {
			return nil, err
		}

		op := utils.M(value)
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
			opKey := utils.S(op["__op"])
			opValue := types.M{}
			if p := utils.M(mongoUpdate[opKey]); p != nil {
				opValue = p
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
			if s := utils.M(mongoUpdate["$set"]); s != nil {
				set = s
			}
			set[key] = value
			mongoUpdate["$set"] = set
		}
	}

	return mongoUpdate, nil
}

func (t *Transform) nestedMongoObjectToNestedParseObject(mongoObject interface{}) (interface{}, error) {
	if mongoObject == nil {
		return mongoObject, nil
	}

	// 转换基本类型
	switch mongoObject.(type) {
	case string, float64, int, bool:
		return mongoObject, nil

	}

	if objs := utils.A(mongoObject); objs != nil {
		results := types.S{}
		for _, o := range objs {
			res, err := t.nestedMongoObjectToNestedParseObject(o)
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
	if object := utils.M(mongoObject); object != nil {
		newObject := types.M{}
		for k, v := range object {
			r, err := t.nestedMongoObjectToNestedParseObject(v)
			if err != nil {
				return nil, err
			}
			newObject[k] = r
		}
		return newObject, nil
	}

	// 无法转换
	return nil, errs.E(errs.InternalServerError, "unknown object type")
}

func (t *Transform) mongoObjectToParseObject(className string, mongoObject interface{}, schema types.M) (interface{}, error) {
	// TODO className 没有用到
	if mongoObject == nil {
		return mongoObject, nil
	}

	// 转换基本类型
	switch mongoObject.(type) {
	case string, float64, int, bool:
		return mongoObject, nil

	}

	if objs := utils.A(mongoObject); objs != nil {
		results := types.S{}
		for _, o := range objs {
			res, err := t.nestedMongoObjectToNestedParseObject(o)
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
	if object := utils.M(mongoObject); object != nil {
		// 转换权限信息
		restObject := t.untransformACL(object)
		for key, value := range object {
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
				if t, ok := value.(time.Time); ok {
					restObject["updatedAt"] = utils.TimetoString(t)
				} else {
					restObject["updatedAt"] = value
				}

			case "createdAt", "_created_at":
				if t, ok := value.(time.Time); ok {
					restObject["createdAt"] = utils.TimetoString(t)
				} else {
					restObject["createdAt"] = value
				}

			case "expiresAt", "_expiresAt":
				if t, ok := value.(time.Time); ok {
					restObject["expiresAt"] = utils.TimetoString(t)
				} else {
					restObject["expiresAt"] = value
				}

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
						authData = utils.M(restObject["authData"])
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
				// 	"className": "abc",
				// 	"objectId":  "123"
				// }
				if strings.HasPrefix(key, "_p_") {
					newKey := key[3:]
					var expected types.M
					if schema != nil {
						if fields := utils.M(schema["fields"]); fields != nil {
							expected = utils.M(fields[newKey])
						}
					}
					if expected == nil {
						// 不在 schema 中的指针类型，丢弃
						break
					}
					if utils.S(expected["type"]) != "Pointer" {
						// schema 中对应的位置不是指针类型，丢弃
						break
					}
					if value == nil {
						break
					}
					objData := strings.Split(utils.S(value), "$")
					if len(objData) != 2 {
						break
					}
					if utils.S(expected["targetClass"]) != objData[0] {
						// 指向了错误的类
						return nil, errs.E(errs.InternalServerError, "pointer to incorrect className")
					}
					restObject[newKey] = types.M{
						"__type":    "Pointer",
						"className": objData[0],
						"objectId":  objData[1],
					}
					break
				} else if strings.HasPrefix(key, "_") && key != "__type" {
					// 转换错误
					return nil, errs.E(errs.InternalServerError, "bad key in untransform: "+key)
				} else {
					var expectedType types.M
					if schema != nil {
						if fields := utils.M(schema["fields"]); fields != nil {
							expectedType = utils.M(fields[key])
						}
					}
					// file 类型
					// {
					// 	"__type": "File",
					// 	"name":   "hello.jpg"
					// }
					f := fileCoder{}
					if expectedType != nil && utils.S(expectedType["type"]) == "File" && f.isValidDatabaseObject(value) {
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
					if expectedType != nil && utils.S(expectedType["type"]) == "GeoPoint" && g.isValidDatabaseObject(value) {
						restObject[key] = g.databaseToJSON(value)
						break
					}
				}
				// 转换子对象
				res, err := t.nestedMongoObjectToNestedParseObject(value)
				if err != nil {
					return nil, err
				}
				restObject[key] = res
			}
		}

		relationFields := types.M{}
		if schema != nil {
			if fields := utils.M(schema["fields"]); fields != nil {
				for fieldName, v := range fields {
					if fieldType := utils.M(v); fieldType != nil {
						if utils.S(fieldType["type"]) == "Relation" {
							relationFields[fieldName] = types.M{
								"__type":    "Relation",
								"className": fieldType["targetClass"],
							}
						}
					}
				}
			}
		}

		for k, v := range relationFields {
			restObject[k] = v
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
func (t *Transform) untransformACL(mongoObject types.M) types.M {
	output := types.M{}
	if mongoObject == nil {
		return output
	}
	if mongoObject["_rperm"] == nil && mongoObject["_wperm"] == nil {
		return output
	}

	acl := types.M{}
	rperm := types.S{}
	wperm := types.S{}
	if mongoObject["_rperm"] != nil {
		rperm = utils.A(mongoObject["_rperm"])
	}
	if mongoObject["_wperm"] != nil {
		wperm = utils.A(mongoObject["_wperm"])
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
	delete(mongoObject, "_rperm")
	delete(mongoObject, "_wperm")

	return output
}

func cannotTransform() interface{} {
	return nil
}

// dateCoder Date 类型数据处理
type dateCoder struct{}

func (d dateCoder) databaseToJSON(object interface{}) types.M {
	if data, ok := object.(time.Time); ok {
		json := utils.TimetoString(data)

		return types.M{
			"__type": "Date",
			"iso":    json,
		}
	}
	return types.M{
		"__type": "Date",
		"iso":    "",
	}
}

func (d dateCoder) isValidDatabaseObject(object interface{}) bool {
	if _, ok := object.(time.Time); ok {
		return true
	}
	return false
}

func (d dateCoder) jsonToDatabase(json types.M) (interface{}, error) {
	t, err := utils.StringtoTime(utils.S(json["iso"]))
	if err != nil {
		return nil, errs.E(errs.InvalidJSON, "invalid iso")
	}
	return t, nil
}

func (d dateCoder) isValidJSON(value types.M) bool {
	return value != nil && utils.S(value["__type"]) == "Date" && utils.S(value["iso"]) != ""
}

// bytesCoder Bytes 类型处理
type bytesCoder struct{}

func (b bytesCoder) databaseToJSON(object interface{}) types.M {
	if data, ok := object.([]byte); ok {
		json := base64.StdEncoding.EncodeToString(data)

		return types.M{
			"__type": "Bytes",
			"base64": json,
		}
	}
	return types.M{
		"__type": "Bytes",
		"base64": "",
	}
}

func (b bytesCoder) isValidDatabaseObject(object interface{}) bool {
	if _, ok := object.([]byte); ok {
		return true
	}
	return false
}

func (b bytesCoder) jsonToDatabase(json types.M) (interface{}, error) {
	by, err := base64.StdEncoding.DecodeString(utils.S(json["base64"]))
	if err != nil {
		return nil, errs.E(errs.InvalidJSON, "invalid base64")
	}
	return by, nil
}

func (b bytesCoder) isValidJSON(value types.M) bool {
	return value != nil && utils.S(value["__type"]) == "Bytes" && utils.S(value["base64"]) != ""
}

// geoPointCoder GeoPoint 类型处理
type geoPointCoder struct{}

func (g geoPointCoder) databaseToJSON(object interface{}) types.M {
	// 只校验空、数组长度，不再校验数据类型，在 isValidDatabaseObject 中校验
	var points types.S
	if points = utils.A(object); points == nil || len(points) != 2 {
		return types.M{
			"__type":    "GeoPoint",
			"longitude": 0,
			"latitude":  0,
		}
	}
	return types.M{
		"__type":    "GeoPoint",
		"longitude": points[0],
		"latitude":  points[1],
	}
}

func (g geoPointCoder) isValidDatabaseObject(object interface{}) bool {
	var points types.S
	if points = utils.A(object); points == nil || len(points) != 2 {
		return false
	}
	p1 := false
	if _, ok := points[0].(float64); ok {
		p1 = true
	}
	if _, ok := points[0].(int); ok {
		p1 = true
	}
	p2 := false
	if _, ok := points[1].(float64); ok {
		p2 = true
	}
	if _, ok := points[1].(int); ok {
		p2 = true
	}
	if p1 && p2 {
		return true
	}
	return false
}

func (g geoPointCoder) jsonToDatabase(json types.M) (interface{}, error) {
	// 默认已经是合法的 geoPoint json
	p1 := false
	if _, ok := json["longitude"].(float64); ok {
		p1 = true
	}
	if _, ok := json["longitude"].(int); ok {
		p1 = true
	}
	if p1 == false {
		return nil, errs.E(errs.InvalidJSON, "invalid longitude")
	}
	p2 := false
	if _, ok := json["latitude"].(float64); ok {
		p2 = true
	}
	if _, ok := json["latitude"].(int); ok {
		p2 = true
	}
	if p2 == false {
		return nil, errs.E(errs.InvalidJSON, "invalid latitude")
	}
	return types.S{json["longitude"], json["latitude"]}, nil
}

func (g geoPointCoder) isValidJSON(value types.M) bool {
	return value != nil && utils.S(value["__type"]) == "GeoPoint" && value["longitude"] != nil && value["latitude"] != nil
}

// fileCoder File 类型处理
type fileCoder struct{}

func (f fileCoder) databaseToJSON(object interface{}) types.M {
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
	if json == nil {
		return nil, nil
	}
	return json["name"], nil
}

func (f fileCoder) isValidJSON(value types.M) bool {
	return value != nil && utils.S(value["__type"]) == "File" && utils.S(value["name"]) != ""
}

// transformInteriorAtom 转换基本类型，以及 Pointer Date Bytes 三种类型，数组与 map 类型无法转换
func (t *Transform) transformInteriorAtom(atom interface{}) (interface{}, error) {
	if atom == nil {
		return atom, nil
	}

	if a := utils.M(atom); a != nil {
		if utils.S(a["__type"]) == "Pointer" {
			return types.M{
				"__type":    "Pointer",
				"className": a["className"],
				"objectId":  a["objectId"],
			}, nil
		}
	}

	// Date 类型
	// {
	// 	"__type": "Date",
	// 	"iso": "2015-03-01T15:59:11-07:00"
	// }
	// ==> 143123456789...
	d := dateCoder{}
	if a := utils.M(atom); a != nil {
		if d.isValidJSON(a) {
			return d.jsonToDatabase(a)
		}
	}

	// Bytes 类型
	// {
	// 	"__type": "Bytes",
	// 	"base64": "aGVsbG8="
	// }
	// ==> hello
	b := bytesCoder{}
	if a := utils.M(atom); a != nil {
		if b.isValidJSON(a) {
			return b.jsonToDatabase(a)
		}
	}

	if a := utils.M(atom); a != nil {
		return cannotTransform(), nil
	}
	if a := utils.A(atom); a != nil {
		return cannotTransform(), nil
	}

	return atom, nil
}

// transformInteriorValue 转换 map 或者 数组 内的值
func (t *Transform) transformInteriorValue(restValue interface{}) (interface{}, error) {
	if restValue == nil {
		return restValue, nil
	}

	if value := utils.M(restValue); value != nil {
		for k := range value {
			if strings.Index(k, "$") > -1 || strings.Index(k, ".") > -1 {
				return nil, errs.E(errs.InvalidNestedKey, "Nested keys should not contain the '$' or '.' characters")
			}
		}
	}

	value, err := t.transformInteriorAtom(restValue)
	if err != nil {
		return nil, err
	}
	if value != cannotTransform() {
		return value, err
	}

	if value := utils.A(restValue); value != nil {
		newValue := types.S{}
		for _, v := range value {
			r, err := t.transformInteriorValue(v)
			if err != nil {
				return nil, err
			}
			newValue = append(newValue, r)
		}
		return newValue, nil
	}

	if value := utils.M(restValue); value != nil {
		if _, ok := value["__op"]; ok {
			return t.transformUpdateOperator(restValue, true)
		}
	}

	if value := utils.M(restValue); value != nil {
		newValue := types.M{}
		for k, v := range value {
			r, err := t.transformInteriorValue(v)
			if err != nil {
				return nil, err
			}
			newValue[k] = r
		}
		return newValue, nil
	}

	return restValue, nil
}
