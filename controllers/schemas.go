package controllers

import (
	"encoding/json"

	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/utils"
)

import "gopkg.in/mgo.v2/bson"

// SchemasController ...
type SchemasController struct {
	ObjectsController
}

// HandleFind ...
// @router / [get]
func (s *SchemasController) HandleFind() {
	result, err := orm.SchemaCollection().GetAllSchemas()
	if err != nil && result == nil {
		s.Data["json"] = bson.M{
			"results": []interface{}{},
		}
		s.ServeJSON()
		return
	}
	for i, v := range result {
		result[i] = orm.MongoSchemaToSchemaAPIResponse(v)
	}
	s.Data["json"] = bson.M{
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
	var data bson.M
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
	if className == "" && bodyClassName == "" {
		// TODO 类名不能为空
		return
	}
	if className == "" {
		className = bodyClassName
	}

	schema := orm.LoadSchema()
	result := schema.AddClassIfNotExists(className, utils.MapInterface(data["fields"]), utils.MapInterface(data["classLevelPermissions"]))

	s.Data["json"] = orm.MongoSchemaToSchemaAPIResponse(result)
	s.ServeJSON()
}

// HandleUpdate ...
// @router /:className [put]
func (s *SchemasController) HandleUpdate() {
	s.ObjectsController.Put()
}

// HandleDelete ...
// @router /:className [delete]
func (s *SchemasController) HandleDelete() {
	s.ObjectsController.Delete()
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
