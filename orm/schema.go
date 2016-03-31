package orm

import "gopkg.in/mgo.v2/bson"
import "github.com/lfq7413/tomato/utils"

var clpValidKeys = []string{"find", "get", "create", "update", "delete", "addField"}
var defaultClassLevelPermissions bson.M

func init() {
	defaultClassLevelPermissions = bson.M{}
	for _, v := range clpValidKeys {
		defaultClassLevelPermissions[v] = bson.M{
			"*": true,
		}
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
	// TODO
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
	return nil
}

func mongoSchemaFromFieldsAndClassNameAndCLP(fields bson.M, className string, classLevelPermissions bson.M) bson.M {
	return nil
}

// Load 返回一个新的 Schema 结构体
func Load(collection *MongoSchemaCollection) *Schema {
	schema := &Schema{
		collection: collection,
	}
	schema.reloadData()
	return schema
}
