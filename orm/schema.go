package orm

import "gopkg.in/mgo.v2/bson"
import "github.com/lfq7413/tomato/utils"

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
func MongoSchemaToSchemaAPIResponse(bson.M) bson.M {
	// TODO
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
