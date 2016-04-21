package orm

import "regexp"

import (
	"strings"

	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

var clpValidKeys = []string{"find", "get", "create", "update", "delete", "addField"}
var defaultClassLevelPermissions types.M
var defaultColumns map[string]types.M
var requiredColumns map[string][]string

func init() {
	defaultClassLevelPermissions = types.M{}
	for _, v := range clpValidKeys {
		defaultClassLevelPermissions[v] = types.M{
			"*": true,
		}
	}
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
	}
	requiredColumns = map[string][]string{
		"_Product": []string{"productIdentifier", "icon", "order", "title", "subtitle"},
		"_Role":    []string{"name", "ACL"},
	}
}

// Schema ...
type Schema struct {
	collection *MongoSchemaCollection
	data       types.M
	perms      types.M
}

// AddClassIfNotExists 添加类定义
func (s *Schema) AddClassIfNotExists(className string, fields types.M, classLevelPermissions types.M) (types.M, error) {
	if s.data[className] != nil {
		// TODO 类已存在
		return nil, nil
	}

	mongoObject := mongoSchemaFromFieldsAndClassNameAndCLP(fields, className, classLevelPermissions)
	if mongoObject["result"] == nil {
		// TODO 转换出现问题
		return nil, nil
	}
	err := s.collection.addSchema(className, utils.MapInterface(mongoObject["result"]))
	if err != nil {
		// TODO 出现错误
		return nil, nil
	}

	return utils.MapInterface(mongoObject["result"]), nil
}

// UpdateClass 更新类
func (s *Schema) UpdateClass(className string, submittedFields types.M, classLevelPermissions types.M) (types.M, error) {
	if s.data[className] == nil {
		// TODO 类不存在
		return nil, nil
	}
	existingFields := utils.CopyMap(utils.MapInterface(s.data[className]))
	existingFields["_id"] = className

	for name, v := range submittedFields {
		field := utils.MapInterface(v)
		op := utils.String(field["__op"])
		if existingFields[name] != nil && op != "Delete" {
			// TODO 字段已存在，不能更新
			return nil, nil
		}
		if existingFields[name] == nil && op == "Delete" {
			// TODO 字不存在，不能删除
			return nil, nil
		}
	}

	newSchema := buildMergedSchemaObject(existingFields, submittedFields)
	mongoObject := mongoSchemaFromFieldsAndClassNameAndCLP(newSchema, className, classLevelPermissions)
	if mongoObject["result"] == nil {
		// TODO 生成错误
		return nil, nil
	}

	insertedFields := []string{}
	for name, v := range submittedFields {
		field := utils.MapInterface(v)
		op := utils.String(field["__op"])
		if op == "Delete" {
			s.deleteField(name, className)
		} else {
			insertedFields = append(insertedFields, name)
		}
	}

	s.reloadData()

	mongoResult := utils.MapInterface(mongoObject["result"])
	for _, fieldName := range insertedFields {
		mongoType := utils.String(mongoResult[fieldName])
		s.validateField(className, fieldName, mongoType, false)
	}

	s.setPermissions(className, classLevelPermissions)

	return MongoSchemaToSchemaAPIResponse(mongoResult), nil
}

func (s *Schema) deleteField(fieldName string, className string) {
	if ClassNameIsValid(className) == false {
		// TODO 无效类名
		return
	}
	if fieldNameIsValid(fieldName) == false {
		// TODO 无效字段名
		return
	}
	if fieldNameIsValidForClass(fieldName, className) == false {
		// TODO 不能修改默认字段
		return
	}

	s.reloadData()

	hasClass := s.hasClass(className)
	if hasClass == false {
		// TODO 类不存在
		return
	}

	class := utils.MapInterface(s.data[className])
	if class[fieldName] == nil {
		// TODO 字段不存在
		return
	}

	name := utils.String(class[fieldName])
	if strings.HasPrefix(name, "relation<") {
		// 删除 _Join table
		DropCollection("_Join:" + fieldName + ":" + className)
	} else {
		collection := AdaptiveCollection(className)
		mongoFieldName := fieldName
		if strings.HasPrefix(name, "*") {
			mongoFieldName = "_p_" + fieldName
		}
		update := types.M{
			"$unset": types.M{mongoFieldName: nil},
		}
		collection.UpdateMany(types.M{}, update)
	}
	update := types.M{
		"$unset": types.M{fieldName: nil},
	}
	s.collection.updateSchema(className, update)
}

func (s *Schema) validateObject(className string, object, query types.M) {
	// TODO 处理错误
	geocount := 0
	s.validateClassName(className, false)

	for k, v := range object {
		if v == nil {
			continue
		}
		expected := getType(v)
		if expected == "geopoint" {
			geocount++
		}
		if geocount > 1 {
			// TODO 只能有一个 geopoint
			return
		}
		if expected == "" {
			continue
		}
		thenValidateField(s, className, k, expected)
	}

	thenValidateRequiredColumns(s, className, object, query)
}

func (s *Schema) validatePermission(className string, aclGroup []string, operation string) {
	// TODO 处理错误
	if s.perms[className] == nil && utils.MapInterface(s.perms[className])[operation] == nil {
		return
	}
	class := utils.MapInterface(s.perms[className])
	perms := utils.MapInterface(class[operation])
	if _, ok := perms["*"]; ok {
		return
	}

	found := false
	for _, v := range aclGroup {
		if _, ok := perms[v]; ok {
			found = true
			break
		}
	}
	if found == false {
		// TODO 无权限
		return
	}
}

func (s *Schema) validateClassName(className string, freeze bool) {
	// TODO 处理错误
	if s.data[className] != nil {
		return
	}
	if freeze {
		// TODO 不能添加
		return
	}
	s.collection.addSchema(className, types.M{})
	s.reloadData()
	s.validateClassName(className, true)
	// TODO 处理上步错误
}

func (s *Schema) validateRequiredColumns(className string, object, query types.M) {
	// TODO 处理错误
	columns := requiredColumns[className]
	if columns == nil || len(columns) == 0 {
		return
	}

	missingColumns := []string{}
	for _, column := range columns {
		if query != nil && query["objectId"] != nil {
			if object[column] != nil && utils.MapInterface(object[column]) != nil {
				o := utils.MapInterface(object[column])
				if utils.String(o["__op"]) == "Delete" {
					missingColumns = append(missingColumns, column)
				}
			}
			continue
		}
		if object[column] == nil {
			missingColumns = append(missingColumns, column)
		}
	}

	if len(missingColumns) > 0 {
		// TODO 缺少字段
		return
	}
}

func (s *Schema) validateField(className, key, fieldtype string, freeze bool) {
	// TODO 检测 key 是否合法
	transformKey(s, className, key)

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
			return
		}
		// TODO 类型不符
		return
	}

	if freeze {
		// TODO 不能修改
		return
	}

	if fieldtype == "" {
		return
	}

	if fieldtype == "geopoint" {
		fields := utils.MapInterface(s.data[className])
		for _, v := range fields {
			otherKey := utils.String(v)
			if otherKey == "geopoint" {
				// TODO 只能有一个 geopoint
				return
			}
		}
	}

	query := types.M{
		key: types.M{"$exists": true},
	}
	update := types.M{
		"$set": types.M{key: fieldtype},
	}
	s.collection.upsertSchema(className, query, update)

	s.reloadData()
	s.validateField(className, key, fieldtype, true)

}

func (s *Schema) setPermissions(className string, perms types.M) {
	validateCLP(perms)
	metadata := types.M{
		"_metadata": types.M{"class_permissions": perms},
	}
	update := types.M{
		"$set": metadata,
	}
	s.collection.updateSchema(className, update)
	s.reloadData()
}

func (s *Schema) hasClass(className string) bool {
	s.reloadData()
	return s.data[className] != nil
}

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

func (s *Schema) getExpectedType(className, key string) string {
	if s.data != nil && s.data[className] != nil {
		cls := utils.MapInterface(s.data[className])
		return utils.String(cls[key])
	}
	return ""
}

func (s *Schema) reloadData() {
	s.data = types.M{}
	s.perms = types.M{}
	results, err := s.collection.GetAllSchemas()
	if err != nil {
		return
	}
	for _, obj := range results {
		className := ""
		classData := types.M{}
		var permsData interface{}

		for k, v := range obj {
			switch k {
			case "_id":
				className = utils.String(v)
			case "_metadata":
				if v != nil && utils.MapInterface(v) != nil && utils.MapInterface(v)["class_permissions"] != nil {
					permsData = utils.MapInterface(v)["class_permissions"]
				}
			default:
				classData[k] = v
			}
		}

		if className != "" {
			s.data[className] = classData
			if permsData != nil {
				s.perms[className] = permsData
			}
		}
	}
}

func thenValidateField(schema *Schema, className, key, fieldtype string) {
	schema.validateField(className, key, fieldtype, false)
}

func thenValidateRequiredColumns(schema *Schema, className string, object, query types.M) {
	schema.validateRequiredColumns(className, object, query)
}

func getType(obj interface{}) string {
	switch obj.(type) {
	case bool:
		return "boolean"
	case string:
		return "string"
	case float64:
		return "number"
	case map[string]interface{}, []interface{}:
		return getObjectType(obj)
	default:
		// TODO 格式无效
		return ""
	}
}

func getObjectType(obj interface{}) string {
	if utils.SliceInterface(obj) != nil {
		return "array"
	}
	if utils.MapInterface(obj) != nil {
		object := utils.MapInterface(obj)
		if object["__type"] != nil {
			t := utils.String(object["__type"])
			switch t {
			case "Pointer":
				if object["className"] != nil {
					return "*" + utils.String(object["className"])
				}
			case "File":
				if object["name"] != nil {
					return "file"
				}
			case "Date":
				if object["iso"] != nil {
					return "date"
				}
			case "GeoPoint":
				if object["latitude"] != nil && object["longitude"] != nil {
					return "geopoint"
				}
			case "Bytes":
				if object["base64"] != nil {
					return "bytes"
				}
			default:
				// TODO 无效的类型
				return ""
			}
		}
		if object["$ne"] != nil {
			return getObjectType(object["$ne"])
		}
		if object["__op"] != nil {
			op := utils.String(object["__op"])
			switch op {
			case "Increment":
				return "number"
			case "Delete":
				return ""
			case "Add", "AddUnique", "Remove":
				return "array"
			case "AddRelation", "RemoveRelation":
				objects := utils.SliceInterface(object["objects"])
				o := utils.MapInterface(objects[0])
				return "relation<" + utils.String(o["className"]) + ">"
			case "Batch":
				ops := utils.SliceInterface(object["ops"])
				return getObjectType(ops[0])
			default:
				// TODO 无效操作
				return ""
			}
		}
	}

	return "object"
}

// MongoSchemaToSchemaAPIResponse ...
func MongoSchemaToSchemaAPIResponse(schema types.M) types.M {
	result := types.M{
		"className": schema["_id"],
		"fields":    mongoSchemaAPIResponseFields(schema),
	}

	classLevelPermissions := utils.CopyMap(defaultClassLevelPermissions)
	if schema["_metadata"] != nil && utils.MapInterface(schema["_metadata"]) != nil {
		metadata := utils.MapInterface(schema["_metadata"])
		if metadata["class_permissions"] != nil && utils.MapInterface(metadata["class_permissions"]) != nil {
			classPermissions := utils.MapInterface(metadata["class_permissions"])
			for k, v := range classPermissions {
				classLevelPermissions[k] = v
			}
		}
	}
	result["classLevelPermissions"] = classLevelPermissions

	return result
}

var nonFieldSchemaKeys = []string{"_id", "_metadata", "_client_permissions"}

func mongoSchemaAPIResponseFields(schema types.M) types.M {
	fieldNames := []string{}
	for k := range schema {
		t := false
		for _, v := range nonFieldSchemaKeys {
			if k == v {
				t = true
				break
			}
		}
		if t == false {
			fieldNames = append(fieldNames, k)
		}
	}
	response := types.M{}
	for _, v := range fieldNames {
		response[v] = mongoFieldTypeToSchemaAPIType(utils.String(schema[v]))
	}
	response["ACL"] = types.M{
		"type": "ACL",
	}
	response["createdAt"] = types.M{
		"type": "Date",
	}
	response["updatedAt"] = types.M{
		"type": "Date",
	}
	response["objectId"] = types.M{
		"type": "String",
	}
	return response
}

func mongoFieldTypeToSchemaAPIType(t string) types.M {
	if t[0] == '*' {
		return types.M{
			"type":        "Pointer",
			"targetClass": string(t[1:]),
		}
	}
	if strings.HasPrefix(t, "relation<") {
		return types.M{
			"type":        "Relation",
			"targetClass": string(t[len("relation<") : len(t)-1]),
		}
	}
	switch t {
	case "number":
		return types.M{
			"type": "Number",
		}
	case "string":
		return types.M{
			"type": "String",
		}
	case "boolean":
		return types.M{
			"type": "Boolean",
		}
	case "date":
		return types.M{
			"type": "Date",
		}
	case "map":
		return types.M{
			"type": "Object",
		}
	case "object":
		return types.M{
			"type": "Object",
		}
	case "array":
		return types.M{
			"type": "Array",
		}
	case "geopoint":
		return types.M{
			"type": "GeoPoint",
		}
	case "file":
		return types.M{
			"type": "File",
		}
	}

	return types.M{}
}

func mongoSchemaFromFieldsAndClassNameAndCLP(fields types.M, className string, classLevelPermissions types.M) types.M {
	if ClassNameIsValid(className) == false {
		// TODO 无效类名
		return nil
	}
	for fieldName := range fields {
		if fieldNameIsValid(fieldName) == false {
			// TODO 无效字段名
			return nil
		}
		if fieldNameIsValidForClass(fieldName, className) == false {
			// TODO 无法添加字段
			return nil
		}
	}

	mongoObject := types.M{
		"_id":       className,
		"objectId":  "string",
		"updatedAt": "string",
		"createdAt": "string",
	}

	if defaultColumns[className] != nil {
		for fieldName := range defaultColumns[className] {
			validatedField := schemaAPITypeToMongoFieldType(utils.MapInterface(defaultColumns[className][fieldName]))
			if validatedField["result"] == nil {
				// TODO 转换错误
				return nil
			}
			mongoObject[fieldName] = validatedField["result"]
		}
	}

	for fieldName := range fields {
		validatedField := schemaAPITypeToMongoFieldType(utils.MapInterface(defaultColumns[className][fieldName]))
		if validatedField["result"] == nil {
			// TODO 转换错误
			return nil
		}
		mongoObject[fieldName] = validatedField["result"]
	}

	geoPoints := []string{}
	for k, v := range mongoObject {
		if utils.String(v) == "geopoint" {
			geoPoints = append(geoPoints, k)
		}
	}
	if len(geoPoints) > 1 {
		// TODO 只能有一个 geoPoint
		return nil
	}

	validateCLP(classLevelPermissions)
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
	}
}

// ClassNameIsValid ...
func ClassNameIsValid(className string) bool {
	return className == "_User" ||
		className == "_Installation" ||
		className == "_Session" ||
		className == "_Role" ||
		className == "_Product" ||
		joinClassIsValid(className) ||
		fieldNameIsValid(className)
}

// InvalidClassNameMessage ...
func InvalidClassNameMessage(className string) string {
	return "Invalid classname: " + className + ", classnames can only have alphanumeric characters and _, and must start with an alpha character "
}

var joinClassRegex = `^_Join:[A-Za-z0-9_]+:[A-Za-z0-9_]+`

func joinClassIsValid(className string) bool {
	b, _ := regexp.MatchString(joinClassRegex, className)
	return b
}

var classAndFieldRegex = `^[A-Za-z][A-Za-z0-9_]*$`

func fieldNameIsValid(fieldName string) bool {
	b, _ := regexp.MatchString(classAndFieldRegex, fieldName)
	return b
}

func fieldNameIsValidForClass(fieldName string, className string) bool {
	if fieldNameIsValid(fieldName) == false {
		return false
	}
	if defaultColumns["_Default"][fieldName] != nil {
		return false
	}
	if defaultColumns[className] != nil && defaultColumns[className][fieldName] != nil {
		return false
	}

	return false
}

func schemaAPITypeToMongoFieldType(t types.M) types.M {
	if utils.String(t["type"]) == "" {
		// TODO type 无效
		return nil
	}
	apiType := utils.String(t["type"])

	if apiType == "Pointer" {
		if t["targetClass"] == nil {
			// TODO 需要 targetClass
			return nil
		}
		if utils.String(t["targetClass"]) == "" {
			// TODO targetClass 无效
			return nil
		}
		targetClass := utils.String(t["targetClass"])
		if ClassNameIsValid(targetClass) == false {
			// TODO 类名无效
			return nil
		}
		return types.M{"result": "*" + targetClass}
	}
	if apiType == "Relation" {
		if t["targetClass"] == nil {
			// TODO 需要 targetClass
			return nil
		}
		if utils.String(t["targetClass"]) == "" {
			// TODO targetClass 无效
			return nil
		}
		targetClass := utils.String(t["targetClass"])
		if ClassNameIsValid(targetClass) == false {
			// TODO 类名无效
			return nil
		}
		return types.M{"result": "relation<" + targetClass + ">"}
	}
	switch apiType {
	case "Number":
		return types.M{"result": "number"}
	case "String":
		return types.M{"result": "string"}
	case "Boolean":
		return types.M{"result": "boolean"}
	case "Date":
		return types.M{"result": "date"}
	case "Object":
		return types.M{"result": "object"}
	case "Array":
		return types.M{"result": "array"}
	case "GeoPoint":
		return types.M{"result": "geopoint"}
	case "File":
		return types.M{"result": "file"}
	default:
		// TODO type 不正确
		return nil
	}
}

func validateCLP(perms types.M) {
	if perms == nil {
		return
	}

	for operation, perm := range perms {
		t := false
		for _, key := range clpValidKeys {
			if operation == key {
				t = true
				break
			}
		}
		if t == false {
			// TODO 不是有效操作
			return
		}

		for key, p := range utils.MapInterface(perm) {
			verifyPermissionKey(key)
			if v, ok := p.(bool); ok {
				if v == false {
					// TODO 值无效
					return
				}
			} else {
				// TODO 值无效
				return
			}
		}
	}
}

// 24 alpha numberic chars + uppercase
var userIDRegex = `^[a-zA-Z0-9]{24}$`

// Anything that start with role
var roleRegex = `^role:.*`

// * permission
var publicRegex = `^\*$`

var permissionKeyRegex = []string{userIDRegex, roleRegex, publicRegex}

func verifyPermissionKey(key string) {
	result := false
	for _, v := range permissionKeyRegex {
		b, _ := regexp.MatchString(v, key)
		result = result || b
	}
	if result == false {
		// TODO 无效的权限名称
		return
	}
}

func buildMergedSchemaObject(mongoObject types.M, putRequest types.M) types.M {
	newSchema := types.M{}

	sysSchemaField := []string{}
	id := utils.String(mongoObject["_id"])
	for k, v := range defaultColumns {
		if k == id {
			for key := range v {
				sysSchemaField = append(sysSchemaField, key)
			}
			break
		}
	}

	for oldField, v := range mongoObject {
		if oldField != "_id" &&
			oldField != "ACL" &&
			oldField != "updatedAt" &&
			oldField != "createdAt" &&
			oldField != "objectId" {
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
			fieldIsDeleted := false
			if putRequest[oldField] != nil {
				op := utils.MapInterface(putRequest[oldField])
				if utils.String(op["__op"]) == "Delete" {
					fieldIsDeleted = true
				}
			}
			if fieldIsDeleted == false {
				newSchema[oldField] = mongoFieldTypeToSchemaAPIType(utils.String(v))
			}
		}
	}

	for newField, v := range putRequest {
		op := utils.MapInterface(v)
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
