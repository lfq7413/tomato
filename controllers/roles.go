package controllers

// RolesController 处理 /roles 接口的请求
type RolesController struct {
	ObjectsController
}

// HandleFind 处理查找 role 请求
// @router / [get]
func (r *RolesController) HandleFind() {
	r.ClassName = "_Role"
	r.ObjectsController.HandleFind()
}

// HandleGet 处理获取指定 role 请求
// @router /:objectId [get]
func (r *RolesController) HandleGet() {
	r.ClassName = "_Role"
	r.ObjectID = r.Ctx.Input.Param(":objectId")
	r.ObjectsController.HandleGet()
}

// HandleCreate 处理创建 role 请求
// @router / [post]
func (r *RolesController) HandleCreate() {
	r.ClassName = "_Role"
	r.ObjectsController.HandleCreate()
}

// HandleUpdate 处理更新指定 role 请求
// @router /:objectId [put]
func (r *RolesController) HandleUpdate() {
	r.ClassName = "_Role"
	r.ObjectID = r.Ctx.Input.Param(":objectId")
	r.ObjectsController.HandleUpdate()
}

// HandleDelete 处理删除指定 role 请求
// @router /:objectId [delete]
func (r *RolesController) HandleDelete() {
	r.ClassName = "_Role"
	r.ObjectID = r.Ctx.Input.Param(":objectId")
	r.ObjectsController.HandleDelete()
}

// Put ...
// @router / [put]
func (r *RolesController) Put() {
	r.ObjectsController.Put()
}

// Delete ...
// @router / [delete]
func (r *RolesController) Delete() {
	r.ObjectsController.Delete()
}
