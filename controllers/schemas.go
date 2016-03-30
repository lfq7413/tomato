package controllers

import "github.com/lfq7413/tomato/orm"
import "github.com/lfq7413/tomato/schema"
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
		result[i] = schema.MongoSchemaToSchemaAPIResponse(v)
	}
	s.Data["json"] = bson.M{
		"results": result,
	}
	s.ServeJSON()
}

// HandleGet ...
// @router /:className [get]
func (s *SchemasController) HandleGet() {
	s.ObjectsController.Get()
}

// HandleCreate ...
// @router / [post]
func (s *SchemasController) HandleCreate() {
	s.ObjectsController.Post()
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
