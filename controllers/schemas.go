package controllers

import (
	"strings"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// SchemasController 处理 /schemas 接口的请求
type SchemasController struct {
	ObjectsController
}

// Prepare 访问 /schemas 接口需要 master key
func (s *SchemasController) Prepare() {
	s.ObjectsController.Prepare()
	if s.Auth.IsMaster == false {
		s.Data["json"] = errs.ErrorMessageToMap(errs.OperationForbidden, "Need master key!")
		s.ServeJSON()
	}
}

// HandleFind 处理 schema 查找请求
// @router / [get]
func (s *SchemasController) HandleFind() {
	schema := orm.TomatoDBController.LoadSchema()
	schemas, err := schema.GetAllSchemas()
	if err != nil {
		s.Data["json"] = types.M{
			"results": types.S{},
		}
		s.ServeJSON()
		return
	}
	s.Data["json"] = types.M{
		"results": schemas,
	}
	s.ServeJSON()
}

// HandleGet 处理查找指定的类请求
// @router /:className [get]
func (s *SchemasController) HandleGet() {
	className := s.Ctx.Input.Param(":className")
	schema := orm.TomatoDBController.LoadSchema()
	sch, err := schema.GetOneSchema(className, false)
	if err != nil {
		s.Data["json"] = errs.ErrorMessageToMap(errs.InvalidClassName, "Class "+className+" does not exist.")
		s.ServeJSON()
		return
	}
	s.Data["json"] = sch
	s.ServeJSON()
}

// HandleCreate 处理创建类请求，同时可匹配 / 的 POST 请求
// @router /:className [post]
func (s *SchemasController) HandleCreate() {
	className := s.Ctx.Input.Param(":className")
	var data = s.JSONBody
	if data == nil {
		s.Data["json"] = errs.ErrorMessageToMap(errs.InvalidJSON, "request body is empty")
		s.ServeJSON()
		return
	}

	bodyClassName := ""
	if data["className"] != nil && utils.S(data["className"]) != "" {
		bodyClassName = utils.S(data["className"])
	}
	if className != "" && bodyClassName != "" {
		if className != bodyClassName {
			s.Data["json"] = errs.ErrorMessageToMap(errs.InvalidClassName, "Class name mismatch between "+bodyClassName+" and "+className+".")
			s.ServeJSON()
			return
		}
	}
	if className == "" {
		className = bodyClassName
	}
	if className == "" {
		s.Data["json"] = errs.ErrorMessageToMap(errs.MissingRequiredFieldError, "POST schemas needs a class name.")
		s.ServeJSON()
		return
	}

	schema := orm.TomatoDBController.LoadSchema()
	result, err := schema.AddClassIfNotExists(className, utils.M(data["fields"]), utils.M(data["classLevelPermissions"]))
	if err != nil {
		s.Data["json"] = errs.ErrorToMap(err)
		s.ServeJSON()
		return
	}

	s.Data["json"] = result
	s.ServeJSON()
}

// HandleUpdate 处理更新类请求
// @router /:className [put]
func (s *SchemasController) HandleUpdate() {
	className := s.Ctx.Input.Param(":className")
	var data = s.JSONBody
	if data == nil {
		s.Data["json"] = errs.ErrorMessageToMap(errs.InvalidJSON, "request body is empty")
		s.ServeJSON()
		return
	}

	bodyClassName := ""
	if data["className"] != nil && utils.S(data["className"]) != "" {
		bodyClassName = utils.S(data["className"])
	}
	if className != bodyClassName {
		s.Data["json"] = errs.ErrorMessageToMap(errs.InvalidClassName, "Class name mismatch between "+bodyClassName+" and "+className+".")
		s.ServeJSON()
		return
	}

	submittedFields := types.M{}
	if data["fields"] != nil && utils.M(data["fields"]) != nil {
		submittedFields = utils.M(data["fields"])
	}

	schema := orm.TomatoDBController.LoadSchema()
	result, err := schema.UpdateClass(className, submittedFields, utils.M(data["classLevelPermissions"]))
	if err != nil {
		s.Data["json"] = errs.ErrorToMap(err)
		s.ServeJSON()
		return
	}

	s.Data["json"] = result
	s.ServeJSON()
}

// HandleDelete 处理删除指定类请求
// @router /:className [delete]
func (s *SchemasController) HandleDelete() {
	className := s.Ctx.Input.Param(":className")
	if orm.ClassNameIsValid(className) == false {
		s.Data["json"] = errs.ErrorMessageToMap(errs.InvalidClassName, orm.InvalidClassNameMessage(className))
		s.ServeJSON()
		return
	}

	err := orm.TomatoDBController.DeleteSchema(className)
	if err != nil {
		s.Data["json"] = errs.ErrorToMap(err)
		s.ServeJSON()
		return
	}

	// 从 _SCHEMA 表中删除类信息，清除相关的 _Join 表
	coll := orm.TomatoDBController.SchemaCollection()
	document, err := coll.FindAndDeleteSchema(className)
	if err != nil {
		s.Data["json"] = errs.ErrorToMap(err)
		s.ServeJSON()
		return
	}
	if document != nil {
		err = removeJoinTables(document)
		if err != nil {
			s.Data["json"] = errs.ErrorToMap(err)
			s.ServeJSON()
			return
		}
	}
	s.Data["json"] = types.M{}
	s.ServeJSON()
	return
}

// removeJoinTables 清除类中的所有关联表
// 需要查找的类型： "field":"relation<otherClass>"
// 需要删除的表明： "_Join:field:className"
func removeJoinTables(mongoSchema types.M) error {
	for field, v := range mongoSchema {
		fieldType := utils.S(v)
		if field != "_metadata" && strings.HasPrefix(fieldType, "relation<") {
			collectionName := "_Join:" + field + ":" + utils.S(mongoSchema["_id"])
			err := orm.Adapter.DeleteOneSchema(collectionName)
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

// injectDefaultSchema 为 schema 添加默认字段
func injectDefaultSchema(schema types.M) types.M {
	defaultSchema := orm.DefaultColumns[schema["className"].(string)]
	if defaultSchema != nil {
		fields := schema["fields"].(map[string]interface{})
		for k, v := range defaultSchema {
			fields[k] = v
		}
		schema["fields"] = fields
	}
	return schema
}
