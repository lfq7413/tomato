package orm

import "gopkg.in/mgo.v2/bson"
import "github.com/lfq7413/tomato/utils"
import "regexp"
import "strings"

// transformKey 把 key 转换为数据库中保存的格式
func transformKey(schema *Schema, className, key string) string {
	k, _ := transformKeyValue(schema, className, key, nil, nil)
	return k
}

// transformKeyValue 把传入的键值对转换为数据库中保存的格式
func transformKeyValue(schema *Schema, className, restKey string, restValue interface{}, options bson.M) (string, interface{}) {
	if options == nil {
		options = bson.M{}
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
		return key, restValue
	case "$or":
		if options["query"] == nil {
			// TODO 只有查询时才能使用 or
			return "", nil
		}
		querys := utils.SliceInterface(restValue)
		if querys == nil {
			// TODO 待转换值必须为数组类型
			return "", nil
		}
		mongoSubqueries := []interface{}{}
		for _, v := range querys {
			query := transformWhere(schema, className, utils.MapInterface(v))
			mongoSubqueries = append(mongoSubqueries, query)
		}
		return "$or", mongoSubqueries
	case "$and":
		if options["query"] == nil {
			// TODO 只有查询时才能使用 and
			return "", nil
		}
		querys := utils.SliceInterface(restValue)
		if querys == nil {
			// TODO 待转换值必须为数组类型
			return "", nil
		}
		mongoSubqueries := []interface{}{}
		for _, v := range querys {
			query := transformWhere(schema, className, utils.MapInterface(v))
			mongoSubqueries = append(mongoSubqueries, query)
		}
		return "$and", mongoSubqueries
	default:
		// 处理第三方 auth 数据
		authDataMatch, _ := regexp.MatchString(`^authData\.([a-zA-Z0-9_]+)\.id$`, key)
		if authDataMatch {
			if options["query"] != nil {
				provider := key[len("authData."):(len(key) - len(".id"))]
				return "_auth_data_" + provider + ".id", restKey
			}
			// TODO 只能将其应用查询操作
			return "", nil
		}
		// 默认处理
		if options["validate"] != nil {
			keyMatch, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_\.]*$`, key)
			if keyMatch == false {
				// TODO 无效的键名
				return "", nil
			}
		}
	}

	// 处理特殊键值
	expected := ""
	if schema != nil {
		expected = schema.getExpectedType(className, key)
	}
	// 处理指向其他对象的字段
	if expected != "" && strings.HasPrefix(expected, "*") {
		key = "_p_" + key
	}
	if expected == "" && restValue != nil {
		op := utils.MapInterface(restValue)
		if op != nil && op["__type"] != nil {
			if utils.String(op["__type"]) == "Pointer" {
				key = "_p_" + key
			}
		}
	}

	inArray := false
	if expected == "array" {
		inArray = true
	}
	// 处理查询操作
	if options["query"] != nil {
		value := transformConstraint(restValue, inArray)
		if value != cannotTransform() {
			return key, value
		}
	}
	if inArray && options["query"] != nil && utils.SliceInterface(restValue) == nil {
		return key, bson.M{"$all": []interface{}{restValue}}
	}

	// 处理原子数据
	value := transformAtom(restValue, false, options)
	if value != cannotTransform() {
		if timeField && utils.String(value) != "" {
			value, _ = utils.StringtoTime(utils.String(value))
		}
		return key, value
	}

	// ACL 在此之前处理，如果依然出现，则返回错误
	if key == "ACL" {
		// TODO 不能在此转换 ACL
		return "", nil
	}

	// 处理数组类型
	if valueArray, ok := restValue.([]interface{}); ok {
		if options["query"] != nil {
			// TODO 查询时不能为数组
			return "", nil
		}
		outValue := []interface{}{}
		for _, restObj := range valueArray {
			_, v := transformKeyValue(schema, className, restKey, restObj, bson.M{"inArray": true})
			outValue = append(outValue, v)
		}
		return key, outValue
	}

	// 处理更新操作
	var flatten bool
	if options["update"] == nil {
		flatten = true
	} else {
		flatten = false
	}
	value = transformUpdateOperator(restValue, flatten)
	if value != cannotTransform() {
		return key, value
	}

	// 处理正常的对象
	normalValue := bson.M{}
	for subRestKey, subRestValue := range utils.MapInterface(restValue) {
		k, v := transformKeyValue(schema, className, subRestKey, subRestValue, bson.M{"inObject": true})
		normalValue[k] = v
	}
	return key, normalValue
}

func transformConstraint(constraint interface{}, inArray bool) interface{} {
	// TODO
	return nil
}

func transformAtom(atom interface{}, force bool, options bson.M) interface{} {
	// TODO
	return nil
}

func transformUpdateOperator(operator interface{}, flatten bool) interface{} {
	// TODO
	return nil
}

// transformCreate ...
func transformCreate(schema *Schema, className string, create bson.M) bson.M {
	// TODO
	return nil
}

func transformWhere(schema *Schema, className string, where bson.M) bson.M {
	// TODO
	return nil
}

func transformUpdate(schema *Schema, className string, update bson.M) bson.M {
	// TODO
	return nil
}

func untransformObjectT(schema *Schema, className string, mongoObject interface{}, isNestedObject bool) interface{} {
	// TODO
	return nil
}

func cannotTransform() interface{} {
	return nil
}
