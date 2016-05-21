package orm

import (
	"regexp"
	"strings"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// clpValidKeys 类级别的权限 列表
var clpValidKeys = []string{"find", "get", "create", "update", "delete", "addField"}

// defaultColumns 所有类的默认字段，以及系统类的默认字段
var defaultColumns map[string]types.M

// requiredColumns 类必须要有的字段
var requiredColumns map[string][]string

// SystemClasses 系统表
var SystemClasses = []string{"_User", "_Installation", "_Role", "_Session", "_Product"}

func init() {
	defaultColumns = map[string]types.M{
		"_Default": types.M{
			"objectId":  types.M{"type": "String"},
			"createdAt": types.M{"type": "Date"},
			"updatedAt": types.M{"type": "Date"},
			"ACL":       types.M{"type": "ACL"},
		},
		"_User": types.M{
			"username":      types.M{"type": "String"},
			"password":      types.M{"type": "String"},
			"authData":      types.M{"type": "Object"},
			"email":         types.M{"type": "String"},
			"emailVerified": types.M{"type": "Boolean"},
		},
		"_Installation": types.M{
			"installationId":   types.M{"type": "String"},
			"deviceToken":      types.M{"type": "String"},
			"channels":         types.M{"type": "Array"},
			"deviceType":       types.M{"type": "String"},
			"pushType":         types.M{"type": "String"},
			"GCMSenderId":      types.M{"type": "String"},
			"timeZone":         types.M{"type": "String"},
			"localeIdentifier": types.M{"type": "String"},
			"badge":            types.M{"type": "Number"},
		},
		"_Role": types.M{
			"name":  types.M{"type": "String"},
			"users": types.M{"type": "Relation", "targetClass": "_User"},
			"roles": types.M{"type": "Relation", "targetClass": "_Role"},
		},
		"_Session": types.M{
			"restricted":     types.M{"type": "Boolean"},
			"user":           types.M{"type": "Pointer", "targetClass": "_User"},
			"installationId": types.M{"type": "String"},
			"sessionToken":   types.M{"type": "String"},
			"expiresAt":      types.M{"type": "Date"},
			"createdWith":    types.M{"type": "Object"},
		},
		"_Product": types.M{
			"productIdentifier": types.M{"type": "String"},
			"download":          types.M{"type": "File"},
			"downloadName":      types.M{"type": "String"},
			"icon":              types.M{"type": "File"},
			"order":             types.M{"type": "Number"},
			"title":             types.M{"type": "String"},
			"subtitle":          types.M{"type": "String"},
		},
		"_PushStatus": types.M{
			"pushTime":      types.M{"type": "String"},
			"source":        types.M{"type": "String"}, // rest or webui
			"query":         types.M{"type": "String"}, // the stringified JSON query
			"payload":       types.M{"type": "Object"}, // the JSON payload,
			"title":         types.M{"type": "String"},
			"expiry":        types.M{"type": "Number"},
			"status":        types.M{"type": "String"},
			"numSent":       types.M{"type": "Number"},
			"numFailed":     types.M{"type": "Number"},
			"pushHash":      types.M{"type": "String"},
			"errorMessage":  types.M{"type": "Object"},
			"sentPerType":   types.M{"type": "Object"},
			"failedPerType": types.M{"type": "Object"},
		},
	}
	requiredColumns = map[string][]string{
		"_Product": []string{"productIdentifier", "icon", "order", "title", "subtitle"},
		"_Role":    []string{"name", "ACL"},
	}
}

// Schema schema 操作对象
type Schema struct {
	collection *MongoSchemaCollection
	data       types.M // data 保存类的字段信息，类型为数据库中保存的类型
	perms      types.M // perms 保存类的操作权限
}

// AddClassIfNotExists 添加类定义，包含默认的字段
func (s *Schema) AddClassIfNotExists(className string, fields types.M, classLevelPermissions types.M) (types.M, error) {
	if s.data[className] != nil {
		return nil, errs.E(errs.InvalidClassName, "Class "+className+" already exists.")
	}

	mongoObject, err := mongoSchemaFromFieldsAndClassNameAndCLP(fields, className, classLevelPermissions)
	if err != nil {
		return nil, err
	}
	err = s.collection.addSchema(className, utils.MapInterface(mongoObject["result"]))
	if err != nil {
		return nil, err
	}

	result := utils.MapInterface(mongoObject["result"])
	return MongoSchemaToParseSchema(result), nil
}

// UpdateClass 更新类
func (s *Schema) UpdateClass(className string, submittedFields types.M, classLevelPermissions types.M) (types.M, error) {
	if s.data[className] == nil {
		return nil, errs.E(errs.InvalidClassName, "Class "+className+" does not exist.")
	}
	// 组装已存在的字段
	existingFields := utils.CopyMap(utils.MapInterface(s.data[className]))
	existingFields["_id"] = className

	// 校验对字段的操作是否合法
	for name, v := range submittedFields {
		field := utils.MapInterface(v)
		op := utils.String(field["__op"])
		if existingFields[name] != nil && op != "Delete" {
			// 字段已存在，不能更新
			return nil, errs.E(errs.ClassNotEmpty, "Field "+name+" exists, cannot update.")
		}
		if existingFields[name] == nil && op == "Delete" {
			// 字段不存在，不能删除
			return nil, errs.E(errs.ClassNotEmpty, "Field "+name+" does not exist, cannot delete.")
		}
	}

	// 组装写入数据库的数据
	newSchema := buildMergedSchemaObject(existingFields, submittedFields)
	mongoObject, err := mongoSchemaFromFieldsAndClassNameAndCLP(newSchema, className, classLevelPermissions)
	if err != nil {
		return nil, err
	}

	// 删除指定字段，并统计需要插入的字段
	insertedFields := []string{}
	for name, v := range submittedFields {
		field := utils.MapInterface(v)
		op := utils.String(field["__op"])
		if op == "Delete" {
			err := s.deleteField(name, className)
			if err != nil {
				return nil, err
			}
		} else {
			insertedFields = append(insertedFields, name)
		}
	}

	// 重新加载修改过的数据
	s.reloadData()

	// 校验并插入字段
	mongoResult := utils.MapInterface(mongoObject["result"])
	for _, fieldName := range insertedFields {
		mongoType := utils.String(mongoResult[fieldName])
		err := s.validateField(className, fieldName, mongoType, false)
		if err != nil {
			return nil, err
		}
	}

	// 设置 CLP
	err = s.setPermissions(className, classLevelPermissions)
	if err != nil {
		return nil, err
	}

	// 把数据库格式的数据转换为 API 格式，并返回
	return MongoSchemaToParseSchema(mongoResult), nil
}

// deleteField 从类定义中删除指定的字段，并删除对象中的数据
func (s *Schema) deleteField(fieldName string, className string) error {
	if ClassNameIsValid(className) == false {
		return errs.E(errs.InvalidClassName, InvalidClassNameMessage(className))
	}
	if fieldNameIsValid(fieldName) == false {
		return errs.E(errs.InvalidKeyName, "invalid field name: "+fieldName)
	}
	if fieldNameIsValidForClass(fieldName, className) == false {
		return errs.E(errs.ChangedImmutableFieldError, "field "+fieldName+" cannot be changed")
	}

	s.reloadData()

	hasClass := s.hasClass(className)
	if hasClass == false {
		return errs.E(errs.InvalidClassName, "Class "+className+" does not exist.")
	}

	class := utils.MapInterface(s.data[className])
	if class[fieldName] == nil {
		return errs.E(errs.ClassNotEmpty, "Field "+fieldName+" does not exist, cannot delete.")
	}

	// 根据字段属性进行相应 对象数据 删除操作
	name := utils.String(class[fieldName])
	if strings.HasPrefix(name, "relation<") {
		// 删除 _Join table 数据
		err := DropCollection("_Join:" + fieldName + ":" + className)
		if err != nil {
			return err
		}
	} else {
		// 删除其他类型字段 对应的对象数据
		collection := AdaptiveCollection(className)
		mongoFieldName := fieldName
		if strings.HasPrefix(name, "*") {
			// Pointer 类型的字段名要添加前缀 _p_
			mongoFieldName = "_p_" + fieldName
		}
		update := types.M{
			"$unset": types.M{mongoFieldName: nil},
		}
		err := collection.UpdateMany(types.M{}, update)
		if err != nil {
			return err
		}
	}
	// 从 _SCHEMA 表中删除相应字段
	update := types.M{
		"$unset": types.M{fieldName: nil},
	}
	return s.collection.updateSchema(className, update)
}

// validateObject 校验对象是否合法
func (s *Schema) validateObject(className string, object, query types.M) error {
	geocount := 0
	err := s.validateClassName(className, false)
	if err != nil {
		return err
	}

	for k, v := range object {
		if v == nil {
			continue
		}
		expected, err := getType(v)
		if err != nil {
			return err
		}
		if expected == "geopoint" {
			geocount++
		}
		if geocount > 1 {
			// 只能有一个 geopoint
			return errs.E(errs.IncorrectType, "there can only be one geopoint field in a class")
		}
		if expected == "" {
			continue
		}
		// 校验字段与字段类型
		err = thenValidateField(s, className, k, expected)
		if err != nil {
			return err
		}
	}

	err = thenValidateRequiredColumns(s, className, object, query)
	if err != nil {
		return err
	}
	return nil
}

// validatePermission 校验对指定类的操作权限
func (s *Schema) validatePermission(className string, aclGroup []string, operation string) error {
	if s.perms[className] == nil && utils.MapInterface(s.perms[className])[operation] == nil {
		return nil
	}
	class := utils.MapInterface(s.perms[className])
	perms := utils.MapInterface(class[operation])
	// 当前操作的权限是公开的
	if _, ok := perms["*"]; ok {
		return nil
	}

	// 查找 acl 中的角色信息是否在权限列表中，找到一个即可
	found := false
	for _, v := range aclGroup {
		if _, ok := perms[v]; ok {
			found = true
			break
		}
	}
	if found == false {
		return errs.E(errs.ObjectNotFound, "Permission denied for this action.")
	}

	return nil
}

// validateClassName 校验类名 freeze 为 true 时，不进行更新
func (s *Schema) validateClassName(className string, freeze bool) error {
	if s.data[className] != nil {
		return nil
	}
	if freeze {
		return errs.E(errs.InvalidJSON, "schema is frozen, cannot add: "+className)
	}

	// 添加不存在的类定义
	err := s.collection.addSchema(className, types.M{})
	if err != nil {

	}
	s.reloadData()
	err = s.validateClassName(className, true)
	if err != nil {
		return errs.E(errs.InvalidJSON, "schema class name does not revalidate")
	}
	return nil
}

// validateRequiredColumns 校验必须的字段
func (s *Schema) validateRequiredColumns(className string, object, query types.M) error {
	columns := requiredColumns[className]
	if columns == nil || len(columns) == 0 {
		return nil
	}

	missingColumns := []string{}
	for _, column := range columns {
		if query != nil && query["objectId"] != nil {
			// 类必须的字段，不能进行删除操作
			if object[column] != nil && utils.MapInterface(object[column]) != nil {
				o := utils.MapInterface(object[column])
				if utils.String(o["__op"]) == "Delete" {
					missingColumns = append(missingColumns, column)
				}
			}
			continue
		}
		// 不能缺少必须的字段
		if object[column] == nil {
			missingColumns = append(missingColumns, column)
		}
	}

	if len(missingColumns) > 0 {
		return errs.E(errs.IncorrectType, missingColumns[0]+" is required.")
	}
	return nil
}

// validateField 校验并插入字段，freeze 为 true 时不进行修改
func (s *Schema) validateField(className, key, fieldtype string, freeze bool) error {
	// 检测 key 是否合法
	_, err := transformKey(s, className, key)
	if err != nil {
		return err
	}

	if strings.Index(key, ".") > 0 {
		key = strings.Split(key, ".")[0]
		fieldtype = "object"
	}

	expected := utils.String(utils.MapInterface(s.data[className])[key])
	if expected != "" {
		if expected == "map" {
			expected = "object"
		}
		if expected == key {
			return nil
		}
		// 类型不符
		return errs.E(errs.IncorrectType, "schema mismatch for "+className+"."+key+"; expected "+expected+" but got "+fieldtype)
	}

	if freeze {
		// 不能修改
		return errs.E(errs.InvalidJSON, "schema is frozen, cannot add "+key+" field")
	}

	// 没有当前要添加的字段，当字段类型为空时，不做更新
	if fieldtype == "" {
		return nil
	}

	if fieldtype == "geopoint" {
		// 只能有一个 geopoint
		fields := utils.MapInterface(s.data[className])
		for _, v := range fields {
			otherKey := utils.String(v)
			if otherKey == "geopoint" {
				return errs.E(errs.IncorrectType, "there can only be one geopoint field in a class")
			}
		}
	}

	// 当前没有该字段，更新 schema
	query := types.M{
		key: types.M{"$exists": false},
	}
	update := types.M{
		"$set": types.M{key: fieldtype},
	}
	err = s.collection.upsertSchema(className, query, update)
	if err != nil {
		// 失败时也需要重新加载数据，因为这时候可能有其他客户端更新了字段
		// s.reloadData()
		// return err
	}

	s.reloadData()
	// 再次尝试校验字段，本次不做更新
	err = s.validateField(className, key, fieldtype, true)
	if err != nil {
		// 字段依然无法校验通过
		return errs.E(errs.InvalidJSON, "schema key will not revalidate")
	}
	return nil
}

// setPermissions 给指定类设置权限
func (s *Schema) setPermissions(className string, perms types.M) error {
	if perms == nil {
		return nil
	}
	err := validateCLP(perms)
	if err != nil {
		return err
	}
	metadata := types.M{
		"_metadata": types.M{"class_permissions": perms},
	}
	update := types.M{
		"$set": metadata,
	}
	err = s.collection.updateSchema(className, update)
	if err != nil {
		return err
	}
	s.reloadData()
	return nil
}

// hasClass Schema 中是否存在类定义
func (s *Schema) hasClass(className string) bool {
	s.reloadData()
	return s.data[className] != nil
}

// hasKeys 指定类中是否存在指定字段
func (s *Schema) hasKeys(className string, keys []string) bool {
	for _, key := range keys {
		if s.data[className] == nil {
			return false
		}
		class := utils.MapInterface(s.data[className])
		if class[key] == nil {
			return false
		}
	}
	return true
}

// getExpectedType 获取期望的字段类型
func (s *Schema) getExpectedType(className, key string) string {
	if s.data != nil && s.data[className] != nil {
		cls := utils.MapInterface(s.data[className])
		return utils.String(cls[key])
	}
	return ""
}

// reloadData 从数据库加载表信息
func (s *Schema) reloadData() {
	s.data = types.M{}
	s.perms = types.M{}
	allSchemas, err := s.collection.GetAllSchemas()
	if err != nil {
		return
	}
	for _, schema := range allSchemas {
		// 组合默认字段
		parseFormatSchema := types.M{}
		for k, v := range defaultColumns["_Default"] {
			parseFormatSchema[k] = v
		}
		if defaultColumns[schema["className"].(string)] != nil {
			for k, v := range defaultColumns[schema["className"].(string)] {
				parseFormatSchema[k] = v
			}
		}
		if schema["fields"].(map[string]interface{}) != nil {
			for k, v := range schema["fields"].(map[string]interface{}) {
				parseFormatSchema[k] = v
			}
		}

		// 无需包含 ACL
		delete(parseFormatSchema, "ACL")
		// createdAt updatedAt 为 string 类型
		parseFormatSchema["createdAt"] = types.M{"type": "String"}
		parseFormatSchema["updatedAt"] = types.M{"type": "String"}

		// 转换为数据库存储格式
		mongoFormatSchema := types.M{}
		for k, v := range parseFormatSchema {
			mongoType, err := schemaAPITypeToMongoFieldType(v.(map[string]interface{}))
			if err != nil {
				continue
			}
			mongoFormatSchema[k] = mongoType
		}

		s.data[schema["className"].(string)] = mongoFormatSchema
		s.perms[schema["className"].(string)] = schema["classLevelPermissions"]
	}
}

// thenValidateField 校验字段，并且不对 schema 进行修改
func thenValidateField(schema *Schema, className, key, fieldtype string) error {
	return schema.validateField(className, key, fieldtype, true)
}

// thenValidateRequiredColumns 校验必须的字段
func thenValidateRequiredColumns(schema *Schema, className string, object, query types.M) error {
	return schema.validateRequiredColumns(className, object, query)
}

// getType 获取对象的格式
func getType(obj interface{}) (string, error) {
	switch obj.(type) {
	case bool:
		return "boolean", nil
	case string:
		return "string", nil
	case float64:
		return "number", nil
	case map[string]interface{}, []interface{}:
		return getObjectType(obj)
	default:
		return "", errs.E(errs.IncorrectType, "bad obj. can not get type")
	}
}

// getObjectType 获取对象格式 仅处理 slice 与 map
func getObjectType(obj interface{}) (string, error) {
	if utils.SliceInterface(obj) != nil {
		return "array", nil
	}
	if utils.MapInterface(obj) != nil {
		object := utils.MapInterface(obj)
		if object["__type"] != nil {
			t := utils.String(object["__type"])
			switch t {
			case "Pointer":
				if object["className"] != nil {
					return "*" + utils.String(object["className"]), nil
				}
			case "File":
				if object["name"] != nil {
					return "file", nil
				}
			case "Date":
				if object["iso"] != nil {
					return "date", nil
				}
			case "GeoPoint":
				if object["latitude"] != nil && object["longitude"] != nil {
					return "geopoint", nil
				}
			case "Bytes":
				if object["base64"] != nil {
					return "bytes", nil
				}
			default:
				// 无效的类型
				return "", errs.E(errs.IncorrectType, "This is not a valid "+t)
			}
		}
		if object["$ne"] != nil {
			return getObjectType(object["$ne"])
		}
		if object["__op"] != nil {
			op := utils.String(object["__op"])
			switch op {
			case "Increment":
				return "number", nil
			case "Delete":
				return "", nil
			case "Add", "AddUnique", "Remove":
				return "array", nil
			case "AddRelation", "RemoveRelation":
				objects := utils.SliceInterface(object["objects"])
				o := utils.MapInterface(objects[0])
				return "relation<" + utils.String(o["className"]) + ">", nil
			case "Batch":
				ops := utils.SliceInterface(object["ops"])
				return getObjectType(ops[0])
			default:
				// 无效操作
				return "", errs.E(errs.IncorrectType, "unexpected op: "+op)
			}
		}
	}

	return "object", nil
}

// mongoSchemaFromFieldsAndClassNameAndCLP 把字段属性转换为数据库中保存的类型
func mongoSchemaFromFieldsAndClassNameAndCLP(fields types.M, className string, classLevelPermissions types.M) (types.M, error) {
	// 校验类名与字段是否合法
	if ClassNameIsValid(className) == false {
		return nil, errs.E(errs.InvalidClassName, InvalidClassNameMessage(className))
	}
	for fieldName := range fields {
		if fieldNameIsValid(fieldName) == false {
			return nil, errs.E(errs.InvalidKeyName, "invalid field name: "+fieldName)
		}
		if fieldNameIsValidForClass(fieldName, className) == false {
			return nil, errs.E(errs.ChangedImmutableFieldError, "field "+fieldName+" cannot be added")
		}
	}

	mongoObject := types.M{
		"_id":       className,
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
	}

	// 添加默认字段
	if defaultColumns[className] != nil {
		for fieldName := range defaultColumns[className] {
			validatedField, err := schemaAPITypeToMongoFieldType(utils.MapInterface(defaultColumns[className][fieldName]))
			if err != nil {
				return nil, err
			}
			mongoObject[fieldName] = validatedField["result"]
		}
	}

	// 添加其他字段
	for fieldName := range fields {
		validatedField, err := schemaAPITypeToMongoFieldType(utils.MapInterface(defaultColumns[className][fieldName]))
		if err != nil {
			return nil, err
		}
		mongoObject[fieldName] = validatedField["result"]
	}

	// 处理 geopoint
	geoPoints := []string{}
	for k, v := range mongoObject {
		if utils.String(v) == "geopoint" {
			geoPoints = append(geoPoints, k)
		}
	}
	if len(geoPoints) > 1 {
		return nil, errs.E(errs.IncorrectType, "currently, only one GeoPoint field may exist in an object. Adding "+geoPoints[1]+" when "+geoPoints[0]+" already exists.")
	}

	// 校验类级别权限
	err := validateCLP(classLevelPermissions)
	if err != nil {
		return nil, err
	}
	// 添加 CLP
	var metadata types.M
	if mongoObject["_metadata"] == nil && utils.MapInterface(mongoObject["_metadata"]) == nil {
		metadata = types.M{}
	} else {
		metadata = utils.MapInterface(mongoObject["_metadata"])
	}
	if classLevelPermissions == nil {
		delete(metadata, "class_permissions")
	} else {
		metadata["class_permissions"] = classLevelPermissions
	}
	mongoObject["_metadata"] = metadata

	return types.M{
		"result": mongoObject,
	}, nil
}

// ClassNameIsValid 校验类名，可以是系统内置类、join 类
// 数字字母组合，以及下划线，但不能以下划线或字母开头
func ClassNameIsValid(className string) bool {
	for _, k := range SystemClasses {
		if className == k {
			return true
		}
	}
	return className == "_SCHEMA" ||
		joinClassIsValid(className) ||
		fieldNameIsValid(className) // 类名与字段名的规则相同
}

// InvalidClassNameMessage ...
func InvalidClassNameMessage(className string) string {
	return "Invalid classname: " + className + ", classnames can only have alphanumeric characters and _, and must start with an alpha character "
}

var joinClassRegex = `^_Join:[A-Za-z0-9_]+:[A-Za-z0-9_]+`

// joinClassIsValid 校验 join 表名， _Join:abc:abc
func joinClassIsValid(className string) bool {
	b, _ := regexp.MatchString(joinClassRegex, className)
	return b
}

var classAndFieldRegex = `^[A-Za-z][A-Za-z0-9_]*$`

// fieldNameIsValid 校验字段名或者类名，数字字母下划线，不以数字下划线开头
func fieldNameIsValid(fieldName string) bool {
	b, _ := regexp.MatchString(classAndFieldRegex, fieldName)
	return b
}

// fieldNameIsValidForClass 校验能否添加指定字段到类中
func fieldNameIsValidForClass(fieldName string, className string) bool {
	// 字段名不合法不能添加
	if fieldNameIsValid(fieldName) == false {
		return false
	}
	// 默认字段不能添加
	if defaultColumns["_Default"][fieldName] != nil {
		return false
	}
	// 当前类的默认字段不能添加
	if defaultColumns[className] != nil && defaultColumns[className][fieldName] != nil {
		return false
	}

	return true
}

// schemaAPITypeToMongoFieldType 把 API 格式的数据转换成 数据库存储的格式
func schemaAPITypeToMongoFieldType(t types.M) (types.M, error) {
	if utils.String(t["type"]) == "" {
		return nil, errs.E(errs.InvalidJSON, "invalid JSON")
	}
	apiType := utils.String(t["type"])

	// {"type":"Pointer", "targetClass":"abc"} => {"result":"*abc"}
	if apiType == "Pointer" {
		if t["targetClass"] == nil {
			return nil, errs.E(errs.MissingRequiredFieldError, "type Pointer needs a class name")
		}
		if utils.String(t["targetClass"]) == "" {
			return nil, errs.E(errs.InvalidJSON, "invalid targetClass")
		}
		targetClass := utils.String(t["targetClass"])
		if ClassNameIsValid(targetClass) == false {
			return nil, errs.E(errs.InvalidClassName, InvalidClassNameMessage(targetClass))
		}
		return types.M{"result": "*" + targetClass}, nil
	}

	// {"type":"Relation", "targetClass":"abc"} => {"result":"relation<abc>"}
	if apiType == "Relation" {
		if t["targetClass"] == nil {
			return nil, errs.E(errs.MissingRequiredFieldError, "type Relation needs a class name")
		}
		if utils.String(t["targetClass"]) == "" {
			return nil, errs.E(errs.InvalidJSON, "invalid targetClass")
		}
		targetClass := utils.String(t["targetClass"])
		if ClassNameIsValid(targetClass) == false {
			return nil, errs.E(errs.InvalidClassName, InvalidClassNameMessage(targetClass))
		}
		return types.M{"result": "relation<" + targetClass + ">"}, nil
	}

	switch apiType {
	case "Number":
		return types.M{"result": "number"}, nil
	case "String":
		return types.M{"result": "string"}, nil
	case "Boolean":
		return types.M{"result": "boolean"}, nil
	case "Date":
		return types.M{"result": "date"}, nil
	case "Object":
		return types.M{"result": "object"}, nil
	case "Array":
		return types.M{"result": "array"}, nil
	case "GeoPoint":
		return types.M{"result": "geopoint"}, nil
	case "File":
		return types.M{"result": "file"}, nil
	default:
		return nil, errs.E(errs.InvalidJSON, "invalid JSON")
	}
}

// validateCLP 校验类级别权限
// 正常的 perms 格式如下
// {
// 	"get":{
// 		"user24id":true,
// 		"role:xxx":true,
// 		"*":true,
// 	},
// 	"delete":{...},
// 	...
// }
func validateCLP(perms types.M) error {
	if perms == nil {
		return nil
	}

	for operation, perm := range perms {
		// 校验是否是系统规定的几种操作
		t := false
		for _, key := range clpValidKeys {
			if operation == key {
				t = true
				break
			}
		}
		if t == false {
			return errs.E(errs.InvalidJSON, operation+" is not a valid operation for class level permissions")
		}

		for key, p := range utils.MapInterface(perm) {
			err := verifyPermissionKey(key)
			if err != nil {
				return err
			}
			if v, ok := p.(bool); ok {
				if v == false {
					return errs.E(errs.InvalidJSON, "false is not a valid value for class level permissions "+operation+":"+key+":false")
				}
			} else {
				return errs.E(errs.InvalidJSON, "this perm is not a valid value for class level permissions "+operation+":"+key+":perm")
			}
		}
	}
	return nil
}

// 24 alpha numberic chars + uppercase
var userIDRegex = `^[a-zA-Z0-9]{24}$`

// Anything that start with role
var roleRegex = `^role:.*`

// * permission
var publicRegex = `^\*$`

var permissionKeyRegex = []string{userIDRegex, roleRegex, publicRegex}

// verifyPermissionKey 校验 CLP 中各种操作包含的角色名是否合法
// 可以是24位的用户 ID，可以是角色名 role:abc ,可以是公共权限 *
func verifyPermissionKey(key string) error {
	result := false
	for _, v := range permissionKeyRegex {
		b, _ := regexp.MatchString(v, key)
		result = result || b
	}
	if result == false {
		return errs.E(errs.InvalidJSON, key+" is not a valid key for class level permissions")
	}
	return nil
}

// buildMergedSchemaObject 组装数据库类型的 mongoObject 与 API 类型的 putRequest，
// 返回值中不包含默认字段，返回的是 API 类型的数据
func buildMergedSchemaObject(mongoObject types.M, putRequest types.M) types.M {
	newSchema := types.M{}

	sysSchemaField := []string{}
	id := utils.String(mongoObject["_id"])
	for k, v := range defaultColumns {
		// 如果是系统预定义的表，则取出默认字段
		if k == id {
			for key := range v {
				sysSchemaField = append(sysSchemaField, key)
			}
			break
		}
	}

	// 处理已经存在的字段
	for oldField, v := range mongoObject {
		// 仅处理以下五种字段以外的字段
		if oldField != "_id" &&
			oldField != "ACL" &&
			oldField != "updatedAt" &&
			oldField != "createdAt" &&
			oldField != "objectId" {
			// 不处理系统默认字段
			if len(sysSchemaField) > 0 {
				t := false
				for _, s := range sysSchemaField {
					if s == oldField {
						t = true
						break
					}
				}
				if t == true {
					continue
				}
			}
			// 处理要删除的字段，要删除的字段不加入返回数据中
			fieldIsDeleted := false
			if putRequest[oldField] != nil {
				op := utils.MapInterface(putRequest[oldField])
				if utils.String(op["__op"]) == "Delete" {
					fieldIsDeleted = true
				}
			}
			if fieldIsDeleted == false {
				newSchema[oldField] = mongoFieldToParseSchemaField(utils.String(v))
			}
		}
	}

	// 处理需要更新的字段
	for newField, v := range putRequest {
		op := utils.MapInterface(v)
		// 不处理 objectId，不处理要删除的字段，跳过系统默认字段，其余字段加入返回数据中
		if newField != "objectId" && utils.String(op["__op"]) != "Delete" {
			if len(sysSchemaField) > 0 {
				t := false
				for _, s := range sysSchemaField {
					if s == newField {
						t = true
						break
					}
				}
				if t == true {
					continue
				}
			}
			newSchema[newField] = v
		}
	}

	return newSchema
}

// Load 返回一个新的 Schema 结构体
func Load(collection *MongoSchemaCollection) *Schema {
	schema := &Schema{
		collection: collection,
	}
	schema.reloadData()
	return schema
}
