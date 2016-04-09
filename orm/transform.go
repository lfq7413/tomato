package orm

import "gopkg.in/mgo.v2/bson"

import "github.com/lfq7413/tomato/utils"
import "regexp"
import "strings"
import "sort"
import "encoding/base64"

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

// transformConstraint 转换查询限制条件
func transformConstraint(constraint interface{}, inArray bool) interface{} {
	// TODO 需要根据 MongoDB 文档修正参数
	if constraint == nil && utils.MapInterface(constraint) == nil {
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
	answer := bson.M{}

	for _, key := range keys {
		switch key {
		case "$lt", "$lte", "$gt", "$gte", "$exists", "$ne", "$eq":
			answer[key] = transformAtom(object[key], true, bson.M{"inArray": inArray})

		case "$in", "$nin":
			arr := utils.SliceInterface(object[key])
			if arr == nil {
				// TODO 必须为数组
				return nil
			}
			answerArr := []interface{}{}
			for _, v := range arr {
				answerArr = append(answerArr, transformAtom(v, true, bson.M{}))
			}
			answer[key] = answerArr

		case "$all":
			arr := utils.SliceInterface(object[key])
			if arr == nil {
				// TODO 必须为数组
				return nil
			}
			answerArr := []interface{}{}
			for _, v := range arr {
				answerArr = append(answerArr, transformAtom(v, true, bson.M{"inArray": true}))
			}
			answer[key] = answerArr

		case "$regex":
			s := utils.String(object[key])
			if s == "" {
				// TODO 必须为字符串
				return nil
			}
			answer[key] = s

		case "$options":
			options := utils.String(object[key])
			if answer["$regex"] == nil || options == "" {
				// TODO 无效值
				return nil
			}
			b, _ := regexp.MatchString(`^[imxs]+$`, options)
			if b == false {
				// TODO 无效值
				return nil
			}
			answer[key] = options

		case "$nearSphere":
			point := utils.MapInterface(object[key])
			answer[key] = []interface{}{point["longitude"], point["latitude"]}

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
			// TODO 暂时不支持该参数
			return nil

		case "$within":
			within := utils.MapInterface(object[key])
			box := utils.SliceInterface(within["$box"])
			if box == nil || len(box) != 2 {
				// TODO 参数不正确
				return nil
			}
			box1 := utils.MapInterface(box[0])
			box2 := utils.MapInterface(box[1])
			answer[key] = bson.M{
				"$box": []interface{}{
					[]interface{}{box1["longitude"], box1["latitude"]},
					[]interface{}{box2["longitude"], box2["latitude"]},
				},
			}

		default:
			b, _ := regexp.MatchString(`^\$+`, key)
			if b {
				// TODO 无效参数
				return nil
			}
			return cannotTransform()
		}
	}

	return answer
}

func transformAtom(atom interface{}, force bool, options bson.M) interface{} {
	// TODO 处理错误
	if options == nil {
		options = bson.M{}
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
			return bson.M{
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

func transformUpdateOperator(operator interface{}, flatten bool) interface{} {
	// TODO 处理错误
	operatorMap := utils.MapInterface(operator)
	if operatorMap == nil || operatorMap["__op"] == nil {
		return cannotTransform()
	}

	op := utils.String(operatorMap["__op"])
	switch op {
	case "Delete":
		if flatten {
			return nil
		}
		return bson.M{
			"__op": "$unset",
			"arg":  "",
		}

	case "Increment":
		if _, ok := operatorMap["amount"].(float64); !ok {
			// TODO 必须为数字
			return nil
		}
		if flatten {
			return operatorMap["amount"]
		}
		return bson.M{
			"__op": "$inc",
			"arg":  operatorMap["amount"],
		}

	case "Add", "AddUnique":
		objects := utils.SliceInterface(operatorMap["objects"])
		if objects == nil {
			// TODO 必须为数组
			return nil
		}
		toAdd := []interface{}{}
		for _, obj := range objects {
			o := transformAtom(obj, true, bson.M{"inArray": true})
			toAdd = append(toAdd, o)
		}
		if flatten {
			return toAdd
		}
		mongoOp := bson.M{
			"Add":       "$push",
			"AddUnique": "$addToSet",
		}[op]
		return bson.M{
			"__op": mongoOp,
			"arg": bson.M{
				"$each": toAdd,
			},
		}

	case "Remove":
		objects := utils.SliceInterface(operatorMap["objects"])
		if objects == nil {
			// TODO 必须为数组
			return nil
		}
		toRemove := []interface{}{}
		for _, obj := range objects {
			o := transformAtom(obj, true, bson.M{"inArray": true})
			toRemove = append(toRemove, o)
		}
		if flatten {
			return []interface{}{}
		}
		return bson.M{
			"__op": "$pullAll",
			"arg":  toRemove,
		}

	default:
		// TODO 不支持的类型
		return nil
	}
}

// transformCreate ...
func transformCreate(schema *Schema, className string, create bson.M) bson.M {
	// TODO 处理错误
	if className == "_User" {
		create = transformAuthData(create)
	}
	mongoCreate := transformACL(create)
	for k, v := range create {
		k, v = transformKeyValue(schema, className, k, v, bson.M{})
		if v != nil {
			mongoCreate[k] = v
		}
	}
	return mongoCreate
}

func transformAuthData(restObject bson.M) bson.M {
	if restObject["authData"] != nil {
		authData := utils.MapInterface(restObject["authData"])
		for provider, v := range authData {
			restObject["_auth_data_"+provider] = v
		}
		delete(restObject, "authData")
	}
	return restObject
}

func transformACL(restObject bson.M) bson.M {
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

type dateCoder struct{}

func (d dateCoder) jsonToDatabase(json bson.M) interface{} {
	t, _ := utils.StringtoTime(utils.String(json["iso"]))
	return t
}

func (d dateCoder) isValidJSON(value bson.M) bool {
	return value != nil && utils.String(value["__type"]) == "Date"
}

type bytesCoder struct{}

func (b bytesCoder) databaseToJSON(object interface{}) interface{} {
	return nil
}

func (b bytesCoder) isValidDatabaseObject(object interface{}) bool {
	return false
}

func (b bytesCoder) jsonToDatabase(json bson.M) interface{} {
	by, _ := base64.StdEncoding.DecodeString(utils.String(json["base64"]))
	return by
}

func (b bytesCoder) isValidJSON(value bson.M) bool {
	return value != nil && utils.String(value["__type"]) == "Bytes"
}

type geoPointCoder struct{}

func (g geoPointCoder) databaseToJSON(object interface{}) interface{} {
	return nil
}

func (g geoPointCoder) isValidDatabaseObject(object interface{}) bool {
	return false
}

func (g geoPointCoder) jsonToDatabase(json bson.M) interface{} {
	return []interface{}{json["longitude"], json["latitude"]}
}

func (g geoPointCoder) isValidJSON(value bson.M) bool {
	return value != nil && utils.String(value["__type"]) == "GeoPoint"
}

type fileCoder struct{}

func (f fileCoder) databaseToJSON(object interface{}) interface{} {
	return nil
}

func (f fileCoder) isValidDatabaseObject(object interface{}) bool {
	return false
}

func (f fileCoder) jsonToDatabase(json bson.M) interface{} {
	return json["name"]
}

func (f fileCoder) isValidJSON(value bson.M) bool {
	return value != nil && utils.String(value["__type"]) == "File"
}
