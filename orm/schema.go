package orm

import (
	"regexp"

	"gopkg.in/mgo.v2/bson"
)
import "github.com/lfq7413/tomato/utils"
import "strings"

var clpValidKeys = []string{"find", "get", "create", "update", "delete", "addField"}
var defaultClassLevelPermissions bson.M
var defaultColumns map[string]bson.M

func init() {
	defaultClassLevelPermissions = bson.M{}
	for _, v := range clpValidKeys {
		defaultClassLevelPermissions[v] = bson.M{
			"*": true,
		}
	}
	defaultColumns = map[string]bson.M{
		"_Default": bson.M{
			"objectId":  bson.M{"type": "String"},
			"createdAt": bson.M{"type": "Date"},
			"updatedAt": bson.M{"type": "Date"},
			"ACL":       bson.M{"type": "ACL"},
		},
		"_User": bson.M{
			"username":      bson.M{"type": "String"},
			"password":      bson.M{"type": "String"},
			"authData":      bson.M{"type": "Object"},
			"email":         bson.M{"type": "String"},
			"emailVerified": bson.M{"type": "Boolean"},
		},
		"_Installation": bson.M{
			"installationId":   bson.M{"type": "String"},
			"deviceToken":      bson.M{"type": "String"},
			"channels":         bson.M{"type": "Array"},
			"deviceType":       bson.M{"type": "String"},
			"pushType":         bson.M{"type": "String"},
			"GCMSenderId":      bson.M{"type": "String"},
			"timeZone":         bson.M{"type": "String"},
			"localeIdentifier": bson.M{"type": "String"},
			"badge":            bson.M{"type": "Number"},
		},
		"_Role": bson.M{
			"name":  bson.M{"type": "String"},
			"users": bson.M{"type": "Relation", "targetClass": "_User"},
			"roles": bson.M{"type": "Relation", "targetClass": "_Role"},
		},
		"_Session": bson.M{
			"restricted":     bson.M{"type": "Boolean"},
			"user":           bson.M{"type": "Pointer", "targetClass": "_User"},
			"installationId": bson.M{"type": "String"},
			"sessionToken":   bson.M{"type": "String"},
			"expiresAt":      bson.M{"type": "Date"},
			"createdWith":    bson.M{"type": "Object"},
		},
		"_Product": bson.M{
			"productIdentifier": bson.M{"type": "String"},
			"download":          bson.M{"type": "File"},
			"downloadName":      bson.M{"type": "String"},
			"icon":              bson.M{"type": "File"},
			"order":             bson.M{"type": "Number"},
			"title":             bson.M{"type": "String"},
			"subtitle":          bson.M{"type": "String"},
		},
	}
}

// Schema ...
type Schema struct {
	collection *MongoSchemaCollection
	data       bson.M
	perms      bson.M
}

// AddClassIfNotExists 添加类定义
func (s *Schema) AddClassIfNotExists(className string, fields bson.M, classLevelPermissions bson.M) bson.M {
	// TODO
	if s.data[className] != nil {
		// TODO 类已存在
		return nil
	}

	mongoObject := mongoSchemaFromFieldsAndClassNameAndCLP(fields, className, classLevelPermissions)
	if mongoObject["result"] == nil {
		// TODO 转换出现问题
		return nil
	}
	err := s.collection.addSchema(className, utils.MapInterface(mongoObject["result"]))
	if err != nil {
		// TODO 出现错误
		return nil
	}

	return utils.MapInterface(mongoObject["result"])
}

func (s *Schema) reloadData() {
	// TODO
	s.data = bson.M{}
	s.perms = bson.M{}
	results, err := s.collection.GetAllSchemas()
	if err != nil {
		return
	}
	for _, obj := range results {
		className := ""
		classData := bson.M{}
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

// MongoSchemaToSchemaAPIResponse ...
func MongoSchemaToSchemaAPIResponse(schema bson.M) bson.M {
	result := bson.M{
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

func mongoSchemaAPIResponseFields(schema bson.M) bson.M {
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
	response := bson.M{}
	for _, v := range fieldNames {
		response[v] = mongoFieldTypeToSchemaAPIType(utils.String(schema[v]))
	}
	response["ACL"] = bson.M{
		"type": "ACL",
	}
	response["createdAt"] = bson.M{
		"type": "Date",
	}
	response["updatedAt"] = bson.M{
		"type": "Date",
	}
	response["objectId"] = bson.M{
		"type": "String",
	}
	return response
}

func mongoFieldTypeToSchemaAPIType(t string) bson.M {
	if t[0] == '*' {
		return bson.M{
			"type":        "Pointer",
			"targetClass": string(t[1:]),
		}
	}
	if strings.HasPrefix(t, "relation<") {
		return bson.M{
			"type":        "Relation",
			"targetClass": string(t[len("relation<") : len(t)-1]),
		}
	}
	switch t {
	case "number":
		return bson.M{
			"type": "Number",
		}
	case "string":
		return bson.M{
			"type": "String",
		}
	case "boolean":
		return bson.M{
			"type": "Boolean",
		}
	case "date":
		return bson.M{
			"type": "Date",
		}
	case "map":
		return bson.M{
			"type": "Object",
		}
	case "object":
		return bson.M{
			"type": "Object",
		}
	case "array":
		return bson.M{
			"type": "Array",
		}
	case "geopoint":
		return bson.M{
			"type": "GeoPoint",
		}
	case "file":
		return bson.M{
			"type": "File",
		}
	}

	return bson.M{}
}

func mongoSchemaFromFieldsAndClassNameAndCLP(fields bson.M, className string, classLevelPermissions bson.M) bson.M {
	if classNameIsValid(className) == false {
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

	mongoObject := bson.M{
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
	var metadata bson.M
	if mongoObject["_metadata"] == nil && utils.MapInterface(mongoObject["_metadata"]) == nil {
		metadata = bson.M{}
	} else {
		metadata = utils.MapInterface(mongoObject["_metadata"])
	}
	if classLevelPermissions == nil {
		delete(metadata, "class_permissions")
	} else {
		metadata["class_permissions"] = classLevelPermissions
	}
	mongoObject["_metadata"] = metadata

	return bson.M{
		"result": mongoObject,
	}
}

func classNameIsValid(className string) bool {
	return className == "_User" ||
		className == "_Installation" ||
		className == "_Session" ||
		className == "_Role" ||
		className == "_Product" ||
		joinClassIsValid(className) ||
		fieldNameIsValid(className)
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

func schemaAPITypeToMongoFieldType(t bson.M) bson.M {
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
		if classNameIsValid(targetClass) == false {
			// TODO 类名无效
			return nil
		}
		return bson.M{"result": "*" + targetClass}
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
		if classNameIsValid(targetClass) == false {
			// TODO 类名无效
			return nil
		}
		return bson.M{"result": "relation<" + targetClass + ">"}
	}
	switch apiType {
	case "Number":
		return bson.M{"result": "number"}
	case "String":
		return bson.M{"result": "string"}
	case "Boolean":
		return bson.M{"result": "boolean"}
	case "Date":
		return bson.M{"result": "date"}
	case "Object":
		return bson.M{"result": "object"}
	case "Array":
		return bson.M{"result": "array"}
	case "GeoPoint":
		return bson.M{"result": "geopoint"}
	case "File":
		return bson.M{"result": "file"}
	default:
		// TODO type 不正确
		return nil
	}
}

func validateCLP(classLevelPermissions bson.M) {

}

// Load 返回一个新的 Schema 结构体
func Load(collection *MongoSchemaCollection) *Schema {
	schema := &Schema{
		collection: collection,
	}
	schema.reloadData()
	return schema
}
