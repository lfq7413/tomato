package orm

import (
	"encoding/base64"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// transformKey 把 key 转换为数据库中保存的格式
func transformKey(schema *Schema, className, key string) string {
	k, _ := transformKeyValue(schema, className, key, nil, nil)
	return k
}

// transformKeyValue 把传入的键值对转换为数据库中保存的格式
func transformKeyValue(schema *Schema, className, restKey string, restValue interface{}, options types.M) (string, interface{}) {
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
		mongoSubqueries := types.S{}
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
		mongoSubqueries := types.S{}
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
		return key, types.M{"$all": types.S{restValue}}
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
		outValue := types.S{}
		for _, restObj := range valueArray {
			_, v := transformKeyValue(schema, className, restKey, restObj, types.M{"inArray": true})
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
	normalValue := types.M{}
	for subRestKey, subRestValue := range utils.MapInterface(restValue) {
		k, v := transformKeyValue(schema, className, subRestKey, subRestValue, types.M{"inObject": true})
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
	answer := types.M{}

	for _, key := range keys {
		switch key {
		case "$lt", "$lte", "$gt", "$gte", "$exists", "$ne", "$eq":
			answer[key] = transformAtom(object[key], true, types.M{"inArray": inArray})

		case "$in", "$nin":
			arr := utils.SliceInterface(object[key])
			if arr == nil {
				// TODO 必须为数组
				return nil
			}
			answerArr := types.S{}
			for _, v := range arr {
				answerArr = append(answerArr, transformAtom(v, true, types.M{}))
			}
			answer[key] = answerArr

		case "$all":
			arr := utils.SliceInterface(object[key])
			if arr == nil {
				// TODO 必须为数组
				return nil
			}
			answerArr := types.S{}
			for _, v := range arr {
				answerArr = append(answerArr, transformAtom(v, true, types.M{"inArray": true}))
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
			answer[key] = types.S{point["longitude"], point["latitude"]}

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
			answer[key] = types.M{
				"$box": types.S{
					types.S{box1["longitude"], box1["latitude"]},
					types.S{box2["longitude"], box2["latitude"]},
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
		return types.M{
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
		return types.M{
			"__op": "$inc",
			"arg":  operatorMap["amount"],
		}

	case "Add", "AddUnique":
		objects := utils.SliceInterface(operatorMap["objects"])
		if objects == nil {
			// TODO 必须为数组
			return nil
		}
		toAdd := types.S{}
		for _, obj := range objects {
			o := transformAtom(obj, true, types.M{"inArray": true})
			toAdd = append(toAdd, o)
		}
		if flatten {
			return toAdd
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
		}

	case "Remove":
		objects := utils.SliceInterface(operatorMap["objects"])
		if objects == nil {
			// TODO 必须为数组
			return nil
		}
		toRemove := types.S{}
		for _, obj := range objects {
			o := transformAtom(obj, true, types.M{"inArray": true})
			toRemove = append(toRemove, o)
		}
		if flatten {
			return types.S{}
		}
		return types.M{
			"__op": "$pullAll",
			"arg":  toRemove,
		}

	default:
		// TODO 不支持的类型
		return nil
	}
}

// transformCreate ...
func transformCreate(schema *Schema, className string, create types.M) types.M {
	// TODO 处理错误
	if className == "_User" {
		create = transformAuthData(create)
	}
	mongoCreate := transformACL(create)
	for k, v := range create {
		key, value := transformKeyValue(schema, className, k, v, types.M{})
		if value != nil {
			mongoCreate[key] = value
		}
	}
	return mongoCreate
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

func transformWhere(schema *Schema, className string, where types.M) types.M {
	// TODO 处理错误
	mongoWhere := types.M{}
	if where["ACL"] != nil {
		// TODO 不能查询 ACL
		return nil
	}
	for k, v := range where {
		options := types.M{
			"query":    true,
			"validate": true,
		}
		key, value := transformKeyValue(schema, className, k, v, options)
		mongoWhere[key] = value
	}

	return mongoWhere
}

func transformUpdate(schema *Schema, className string, update types.M) types.M {
	// TODO 处理错误
	if update == nil {
		// TODO 更新数据不能为空
		return nil
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
		key, value := transformKeyValue(schema, className, k, v, types.M{"update": true})

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

	return mongoUpdate
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
