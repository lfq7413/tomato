package controllers

// RolesController ...
type RolesController struct {
	ObjectsController
}

// HandleFind ...
// @router / [get]
func (r *RolesController) HandleFind() {
	r.ClassName = "_Role"
	r.ObjectsController.HandleFind()
}

// HandleGet ...
// @router /:objectId [get]
func (r *RolesController) HandleGet() {
	r.ClassName = "_Role"
	r.ObjectID = r.Ctx.Input.Param(":objectId")
	r.ObjectsController.HandleGet()
}

// HandleCreate ...
// @router / [post]
func (r *RolesController) HandleCreate() {
	r.ClassName = "_Role"
	r.ObjectsController.HandleCreate()
}

// HandleUpdate ...
// @router /:objectId [put]
func (r *RolesController) HandleUpdate() {
	r.ClassName = "_Role"
	r.ObjectID = r.Ctx.Input.Param(":objectId")
	r.ObjectsController.HandleUpdate()
}

// HandleDelete ...
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
