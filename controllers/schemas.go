package controllers

import (
	"encoding/json"
	"strings"

	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// SchemasController ...
type SchemasController struct {
	ObjectsController
}

// HandleFind ...
// @router / [get]
func (s *SchemasController) HandleFind() {
	result, err := orm.SchemaCollection().GetAllSchemas()
	if err != nil && result == nil {
		s.Data["json"] = types.M{
			"results": types.S{},
		}
		s.ServeJSON()
		return
	}
	for i, v := range result {
		result[i] = orm.MongoSchemaToSchemaAPIResponse(v)
	}
	s.Data["json"] = types.M{
		"results": result,
	}
	s.ServeJSON()
}

// HandleGet ...
// @router /:className [get]
func (s *SchemasController) HandleGet() {
	className := s.Ctx.Input.Param(":className")
	result, err := orm.SchemaCollection().FindSchema(className)
	if err != nil && result == nil {
		// TODO 类不存在
		return
	}
	s.Data["json"] = result
	s.ServeJSON()
}

// HandleCreate ...
// @router /:className [post]
func (s *SchemasController) HandleCreate() {
	className := s.Ctx.Input.Param(":className")
	if s.Ctx.Input.RequestBody == nil {
		// TODO 数据为空
		return
	}
	var data types.M
	err := json.Unmarshal(s.Ctx.Input.RequestBody, &data)
	if err != nil {
		// TODO 解析错误
		return
	}
	bodyClassName := ""
	if data["className"] != nil && utils.String(data["className"]) != "" {
		bodyClassName = utils.String(data["className"])
	}
	if className != "" && bodyClassName != "" {
		if className != bodyClassName {
			// TODO 类名不一致
			return
		}
	}
	if className == "" {
		className = bodyClassName
	}
	if className == "" {
		// TODO 类名不能为空
		return
	}

	schema := orm.LoadSchema(nil)
	result := schema.AddClassIfNotExists(className, utils.MapInterface(data["fields"]), utils.MapInterface(data["classLevelPermissions"]))

	s.Data["json"] = orm.MongoSchemaToSchemaAPIResponse(result)
	s.ServeJSON()
}

// HandleUpdate ...
// @router /:className [put]
func (s *SchemasController) HandleUpdate() {
	className := s.Ctx.Input.Param(":className")
	if s.Ctx.Input.RequestBody == nil {
		// TODO 数据为空
		return
	}
	var data types.M
	err := json.Unmarshal(s.Ctx.Input.RequestBody, &data)
	if err != nil {
		// TODO 解析错误
		return
	}
	bodyClassName := ""
	if data["className"] != nil && utils.String(data["className"]) != "" {
		bodyClassName = utils.String(data["className"])
	}
	if className != bodyClassName {
		// TODO 类名不一致
		return
	}

	submittedFields := types.M{}
	if data["fields"] != nil && utils.MapInterface(data["fields"]) != nil {
		submittedFields = utils.MapInterface(data["fields"])
	}

	schema := orm.LoadSchema(nil)
	result := schema.UpdateClass(className, submittedFields, utils.MapInterface(data["classLevelPermissions"]))

	s.Data["json"] = result
	s.ServeJSON()
}

// HandleDelete ...
// @router /:className [delete]
func (s *SchemasController) HandleDelete() {
	className := s.Ctx.Input.Param(":className")
	if orm.ClassNameIsValid(className) == false {
		// TODO 类名无效
		return
	}

	exist := orm.CollectionExists(className)
	if exist == false {
		return
	}

	collection := orm.AdaptiveCollection(className)
	count := collection.Count(types.M{}, types.M{})
	if count > 0 {
		// TODO 类不为空
		return
	}
	collection.Drop()

	coll := orm.SchemaCollection()
	// TODO 处理错误
	document, _ := coll.FindAndDeleteSchema(className)
	if document != nil {
		removeJoinTables(document)
		// TODO 处理错误
	}
}

func removeJoinTables(mongoSchema types.M) error {
	for field, v := range mongoSchema {
		fieldType := utils.String(v)
		if strings.HasPrefix(fieldType, "relation<") {
			collectionName := "_Join:" + field + ":" + utils.String(mongoSchema["_id"])
			err := orm.DropCollection(collectionName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Delete ...
// @router / [delete]
func (s *SchemasController) Delete() {
	s.ObjectsController.Delete()
}

// Put ...
// @router / [put]
func (s *SchemasController) Put() {
	s.ObjectsController.Put()
}
