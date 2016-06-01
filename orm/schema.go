package orm

import (
	"regexp"
	"strings"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/storage"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// clpValidKeys 类级别的权限 列表
var clpValidKeys = []string{"find", "get", "create", "update", "delete", "addField", "readUserFields", "writeUserFields"}

// DefaultColumns 所有类的默认字段，以及系统类的默认字段
var DefaultColumns map[string]types.M

// requiredColumns 类必须要有的字段
var requiredColumns map[string][]string

// SystemClasses 系统表
var SystemClasses = []string{"_User", "_Installation", "_Role", "_Session", "_Product", "_PushStatus"}

func init() {
	DefaultColumns = map[string]types.M{
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
	collection storage.SchemaCollection
	dbAdapter  storage.Adapter
	data       types.M // data 保存类的字段信息，类型为 API 类型
	perms      types.M // perms 保存类的操作权限
}

// AddClassIfNotExists 添加类定义，包含默认的字段
func (s *Schema) AddClassIfNotExists(className string, fields types.M, classLevelPermissions types.M) (types.M, error) {
	err := s.validateNewClass(className, fields, classLevelPermissions)
	if err != nil {
		return nil, err
	}

	result, err := s.collection.AddSchema(className, fields, classLevelPermissions)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateClass 更新类
func (s *Schema) UpdateClass(className string, submittedFields types.M, classLevelPermissions types.M) (types.M, error) {
	if s.hasClass(className) == false {
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
	err := s.validateSchemaData(className, newSchema, classLevelPermissions)
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
	for _, fieldName := range insertedFields {
		fieldType := submittedFields[fieldName].(map[string]interface{})
		err := s.validateField(className, fieldName, fieldType, false)
		if err != nil {
			return nil, err
		}
	}

	// 设置 CLP
	err = s.setPermissions(className, classLevelPermissions, newSchema)
	if err != nil {
		return nil, err
	}

	return types.M{
		"className":             className,
		"fields":                s.data[className],
		"classLevelPermissions": s.perms[className],
	}, nil
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
	name := utils.MapInterface(class[fieldName])["type"].(string)
	if name == "Relation" {
		// 删除表数据与 schema 中的对应字段
		err := Adapter.DeleteFields(className, []string{fieldName}, []string{})
		if err != nil {
			return err
		}
		// 删除 _Join table 数据
		err = Adapter.DropCollection("_Join:" + fieldName + ":" + className)
		if err != nil {
			return err
		}
	}
	// 删除其他类型字段 对应的对象数据
	fieldNames := []string{fieldName}
	pointerFieldNames := []string{}

	if name == "Pointer" {
		pointerFieldNames = append(pointerFieldNames, fieldName)
	}
	return Adapter.DeleteFields(className, fieldNames, pointerFieldNames)
}

// validateObject 校验对象是否合法
func (s *Schema) validateObject(className string, object, query types.M) error {
	geocount := 0
	err := s.enforceClassExists(className, false)
	if err != nil {
		return err
	}

	for fieldName, v := range object {
		if v == nil {
			continue
		}
		expected, err := getType(v)
		if err != nil {
			return err
		}
		if expected == nil {
			continue
		}
		if expected["type"].(string) == "GeoPoint" {
			geocount++
		}
		if geocount > 1 {
			// 只能有一个 geopoint
			return errs.E(errs.IncorrectType, "there can only be one geopoint field in a class")
		}
		if fieldName == "ACL" {
			// 每个对象都隐含 ACL 字段
			continue
		}
		// 校验字段与字段类型
		err = thenValidateField(s, className, fieldName, expected)
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
	classPerms := utils.MapInterface(s.perms[className])
	perms := utils.MapInterface(classPerms[operation])
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

	if found {
		return nil
	}

	var permissionField string
	if operation == "get" || operation == "find" {
		permissionField = "readUserFields"
	} else {
		permissionField = "writeUserFields"
	}

	if permissionField == "writeUserFields" && operation == "create" {
		return errs.E(errs.OperationForbidden, "Permission denied for this action.")
	}

	if v, ok := classPerms[permissionField].([]interface{}); ok && len(v) > 0 {
		return nil
	}

	return errs.E(errs.OperationForbidden, "Permission denied for this action.")
}

// enforceClassExists 校验类名 freeze 为 true 时，不进行更新
func (s *Schema) enforceClassExists(className string, freeze bool) error {
	if s.data[className] != nil {
		return nil
	}
	if freeze {
		return errs.E(errs.InvalidJSON, "schema is frozen, cannot add: "+className)
	}

	// 添加不存在的类定义
	_, err := s.AddClassIfNotExists(className, types.M{}, types.M{})
	if err != nil {

	}
	s.reloadData()
	err = s.enforceClassExists(className, true)
	if err != nil {
		return errs.E(errs.InvalidJSON, "schema class name does not revalidate")
	}
	return nil
}

// validateNewClass 校验新建的类
func (s *Schema) validateNewClass(className string, fields types.M, classLevelPermissions types.M) error {
	if s.data[className] != nil {
		return errs.E(errs.InvalidClassName, "Class "+className+" already exists.")
	}

	if ClassNameIsValid(className) == false {
		return errs.E(errs.InvalidClassName, InvalidClassNameMessage(className))
	}

	return s.validateSchemaData(className, fields, classLevelPermissions)
}

// validateSchemaData 校验 Schema 数据
func (s *Schema) validateSchemaData(className string, fields types.M, classLevelPermissions types.M) error {
	for fieldName, v := range fields {
		if fieldNameIsValid(fieldName) == false {
			return errs.E(errs.InvalidKeyName, "invalid field name: "+fieldName)
		}
		if fieldNameIsValidForClass(fieldName, className) == false {
			return errs.E(errs.ChangedImmutableFieldError, "field "+fieldName+" cannot be added")
		}
		err := fieldTypeIsInvalid(v.(map[string]interface{}))
		if err != nil {
			return err
		}
	}

	if DefaultColumns[className] != nil {
		for fieldName, v := range DefaultColumns[className] {
			fields[fieldName] = v
		}
	}

	geoPoints := []string{}
	for key, v := range fields {
		if v != nil {
			fieldData := v.(map[string]interface{})
			if fieldData["type"].(string) == "GeoPoint" {
				geoPoints = append(geoPoints, key)
			}
		}
	}
	if len(geoPoints) > 1 {
		return errs.E(errs.IncorrectType, "currently, only one GeoPoint field may exist in an object. Adding "+geoPoints[1]+" when "+geoPoints[0]+" already exists.")
	}

	return validateCLP(classLevelPermissions, fields)
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
func (s *Schema) validateField(className, fieldName string, fieldtype types.M, freeze bool) error {
	s.reloadData()
	// 检测 fieldName 是否合法
	_, err := Transform.TransformKey(s, className, fieldName)
	if err != nil {
		return err
	}

	if strings.Index(fieldName, ".") > 0 {
		fieldName = strings.Split(fieldName, ".")[0]
		fieldtype = types.M{
			"type": "Object",
		}
	}

	expected := utils.MapInterface(utils.MapInterface(s.data[className])[fieldName])
	if expected != nil {
		if expected["type"].(string) == "map" {
			expected["type"] = "Object"
		}
		if expected["type"].(string) == fieldtype["type"].(string) && expected["targetClass"].(string) == fieldtype["targetClass"].(string) {
			return nil
		}
		// 类型不符
		return errs.E(errs.IncorrectType, "schema mismatch for "+className+"."+fieldName+"; expected "+expected["type"].(string)+" but got "+fieldtype["type"].(string))
	}

	if freeze {
		// 不能修改
		return errs.E(errs.InvalidJSON, "schema is frozen, cannot add "+fieldName+" field")
	}

	// 没有当前要添加的字段，当字段类型为空时，不做更新
	if fieldtype == nil || fieldtype["type"].(string) == "" {
		return nil
	}

	err = s.collection.AddFieldIfNotExists(className, fieldName, fieldtype)
	if err != nil {
		// 失败时也需要重新加载数据，因为这时候可能有其他客户端更新了字段
		// s.reloadData()
		// return err
	}

	s.reloadData()
	// 再次尝试校验字段，本次不做更新
	err = s.validateField(className, fieldName, fieldtype, true)
	if err != nil {
		// 字段依然无法校验通过
		return errs.E(errs.InvalidJSON, "schema key will not revalidate")
	}
	return nil
}

// setPermissions 给指定类设置权限
func (s *Schema) setPermissions(className string, perms types.M, newSchema types.M) error {
	if perms == nil {
		return nil
	}
	err := validateCLP(perms, newSchema)
	if err != nil {
		return err
	}
	metadata := types.M{
		"_metadata": types.M{"class_permissions": perms},
	}
	update := types.M{
		"$set": metadata,
	}
	err = s.collection.UpdateSchema(className, update)
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

// GetExpectedType 获取期望的字段类型
func (s *Schema) GetExpectedType(className, key string) types.M {
	if s.data != nil && s.data[className] != nil {
		cls := utils.MapInterface(s.data[className])
		return utils.MapInterface(cls[key])
	}
	return nil
}

// GetRelationFields 转换 relation 类型的字段为 API 格式
func (s *Schema) GetRelationFields(className string) types.M {
	relationFields := types.M{}
	if s.data != nil && s.data[className] != nil {
		classData := s.data[className].(map[string]interface{})
		for field, v := range classData {
			fieldType := v.(map[string]interface{})
			if fieldType["type"].(string) != "Relation" {
				continue
			}
			name := fieldType["targetClass"].(string)
			relationFields[field] = types.M{
				"__type":    "Relation",
				"className": name,
			}
		}
	}

	return relationFields
}

// reloadData 从数据库加载表信息
func (s *Schema) reloadData() {
	s.data = types.M{}
	s.perms = types.M{}
	allSchemas, err := s.GetAllSchemas()
	if err != nil {
		return
	}
	for _, schema := range allSchemas {
		s.data[schema["className"].(string)] = schema
		s.perms[schema["className"].(string)] = schema["classLevelPermissions"]
	}
}

// GetAllSchemas ...
func (s *Schema) GetAllSchemas() ([]types.M, error) {
	allSchemas, err := Adapter.GetAllSchemas()
	if err != nil {
		return nil, err
	}
	schems := []types.M{}
	for _, v := range allSchemas {
		schems = append(schems, injectDefaultSchema(v))
	}
	return schems, nil
}

// GetOneSchema ...
func (s *Schema) GetOneSchema(className string) (types.M, error) {
	schema, err := Adapter.GetOneSchema(className)
	if err != nil {
		return nil, err
	}
	return injectDefaultSchema(schema), nil
}

// thenValidateField 校验字段，并且不对 schema 进行修改
func thenValidateField(schema *Schema, className, key string, fieldtype types.M) error {
	return schema.validateField(className, key, fieldtype, true)
}

// thenValidateRequiredColumns 校验必须的字段
func thenValidateRequiredColumns(schema *Schema, className string, object, query types.M) error {
	return schema.validateRequiredColumns(className, object, query)
}

// getType 获取对象的格式
func getType(obj interface{}) (types.M, error) {
	switch obj.(type) {
	case bool:
		return types.M{"type": "Boolean"}, nil
	case string:
		return types.M{"type": "String"}, nil
	case float64:
		return types.M{"type": "Number"}, nil
	case map[string]interface{}, []interface{}:
		return getObjectType(obj)
	default:
		return nil, errs.E(errs.IncorrectType, "bad obj. can not get type")
	}
}

// getObjectType 获取对象格式 仅处理 slice 与 map
func getObjectType(obj interface{}) (types.M, error) {
	if utils.SliceInterface(obj) != nil {
		return types.M{"type": "Array"}, nil
	}
	if utils.MapInterface(obj) != nil {
		object := utils.MapInterface(obj)
		if object["__type"] != nil {
			t := utils.String(object["__type"])
			switch t {
			case "Pointer":
				if object["className"] != nil {
					return types.M{
						"type":        "Pointer",
						"targetClass": object["className"],
					}, nil
				}
			case "File":
				if object["name"] != nil {
					return types.M{"type": "File"}, nil
				}
			case "Date":
				if object["iso"] != nil {
					return types.M{"type": "Date"}, nil
				}
			case "GeoPoint":
				if object["latitude"] != nil && object["longitude"] != nil {
					return types.M{"type": "Geopoint"}, nil
				}
			case "Bytes":
				if object["base64"] != nil {
					return types.M{"type": "Bytes"}, nil
				}
			default:
				// 无效的类型
				return nil, errs.E(errs.IncorrectType, "This is not a valid "+t)
			}
		}
		if object["$ne"] != nil {
			return getObjectType(object["$ne"])
		}
		if object["__op"] != nil {
			op := utils.String(object["__op"])
			switch op {
			case "Increment":
				return types.M{"type": "Number"}, nil
			case "Delete":
				return nil, nil
			case "Add", "AddUnique", "Remove":
				return types.M{"type": "Array"}, nil
			case "AddRelation", "RemoveRelation":
				objects := utils.SliceInterface(object["objects"])
				o := utils.MapInterface(objects[0])
				return types.M{
					"type":        "Relation",
					"targetClass": utils.String(o["className"]),
				}, nil
			case "Batch":
				ops := utils.SliceInterface(object["ops"])
				return getObjectType(ops[0])
			default:
				// 无效操作
				return nil, errs.E(errs.IncorrectType, "unexpected op: "+op)
			}
		}
	}

	return types.M{"type": "object"}, nil
}

// ClassNameIsValid 校验类名，可以是系统内置类、join 类
// 数字字母组合，以及下划线，但不能以下划线或字母开头
func ClassNameIsValid(className string) bool {
	for _, k := range SystemClasses {
		if className == k {
			return true
		}
	}
	return joinClassIsValid(className) ||
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
	if DefaultColumns["_Default"][fieldName] != nil {
		return false
	}
	// 当前类的默认字段不能添加
	if DefaultColumns[className] != nil && DefaultColumns[className][fieldName] != nil {
		return false
	}

	return true
}

var validNonRelationOrPointerTypes = []string{
	"Number",
	"String",
	"Boolean",
	"Date",
	"Object",
	"Array",
	"GeoPoint",
	"File",
}

// fieldTypeIsInvalid 检测字段类型是否合法
func fieldTypeIsInvalid(t types.M) error {
	var invalidJSONError = errs.E(errs.InvalidJSON, "invalid JSON")
	fieldType := ""
	if v, ok := t["type"].(string); ok {
		fieldType = v
	} else {
		return invalidJSONError
	}
	targetClass := ""
	if fieldType == "Pointer" || fieldType == "Relation" {
		if _, ok := t["targetClass"]; ok == false {
			return errs.E(errs.MissingRequiredFieldError, "type "+fieldType+" needs a class name")
		}
		if v, ok := t["targetClass"].(string); ok {
			targetClass = v
		} else {
			return invalidJSONError
		}
		if ClassNameIsValid(targetClass) == false {
			return errs.E(errs.InvalidClassName, InvalidClassNameMessage(targetClass))
		}
		return nil
	}

	in := false
	for _, v := range validNonRelationOrPointerTypes {
		if fieldType == v {
			in = true
			break
		}
	}
	if in == false {
		return errs.E(errs.IncorrectType, "invalid field type: "+fieldType)
	}

	return nil
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
//  "readUserFields":{"aaa","bbb"}
// 	...
// }
func validateCLP(perms types.M, fields types.M) error {
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

		if operation == "readUserFields" || operation == "writeUserFields" {
			if p, ok := perm.([]interface{}); ok {
				for _, v := range p {
					key := v.(string)
					// 字段类型必须为指向 _User 的指针类型
					if fields[key] != nil {
						if t, ok := fields[key].(map[string]interface{}); ok {
							if t["type"].(string) == "Pointer" && t["targetClass"].(string) == "_User" {
								continue
							}
						}
					}
					return errs.E(errs.InvalidJSON, key+" is not a valid column for class level pointer permissions "+operation)
				}
				return nil
			}
			return errs.E(errs.InvalidJSON, "this perms[operation] is not a valid value for class level permissions "+operation)
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

// buildMergedSchemaObject 组装数据库类型的 existingFields 与 API 类型的 putRequest，
// 返回值中不包含默认字段，返回的是 API 类型的数据
func buildMergedSchemaObject(existingFields types.M, putRequest types.M) types.M {
	newSchema := types.M{}

	sysSchemaField := []string{}
	id := utils.String(existingFields["_id"])
	for k, v := range DefaultColumns {
		// 如果是系统预定义的表，则取出默认字段
		if k == id {
			for key := range v {
				sysSchemaField = append(sysSchemaField, key)
			}
			break
		}
	}

	// 处理已经存在的字段
	for oldField, v := range existingFields {
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
				newSchema[oldField] = v
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

// injectDefaultSchema 为 schema 添加默认字段
func injectDefaultSchema(schema types.M) types.M {
	newSchema := types.M{}
	fields := schema["fields"].(map[string]interface{})
	defaultFieldsSchema := DefaultColumns["_Default"]
	for k, v := range defaultFieldsSchema {
		fields[k] = v
	}
	defaultSchema := DefaultColumns[schema["className"].(string)]
	if defaultSchema != nil {
		for k, v := range defaultSchema {
			fields[k] = v
		}
	}
	newSchema["fields"] = fields
	newSchema["className"] = schema["className"]
	newSchema["classLevelPermissions"] = schema["classLevelPermissions"]

	return newSchema
}

// Load 返回一个新的 Schema 结构体
func Load(collection storage.SchemaCollection, adapter storage.Adapter) *Schema {
	schema := &Schema{
		collection: collection,
		dbAdapter:  adapter,
	}
	schema.reloadData()
	return schema
}
