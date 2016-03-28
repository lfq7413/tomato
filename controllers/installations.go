package controllers

// InstallationsController ...
type InstallationsController struct {
	ObjectsController
}

// HandleFind ...
// @router / [get]
func (i *InstallationsController) HandleFind() {
	i.ClassName = "_Installation"
	i.ObjectsController.HandleFind()
}

// HandleGet ...
// @router /:objectId [get]
func (i *InstallationsController) HandleGet() {
	i.ClassName = "_Installation"
	i.ObjectID = i.Ctx.Input.Param(":objectId")
	i.ObjectsController.HandleGet()
}

// HandleCreate ...
// @router / [post]
func (i *InstallationsController) HandleCreate() {
	i.ClassName = "_Installation"
	i.ObjectsController.HandleCreate()
}

// HandleUpdate ...
// @router /:objectId [put]
func (i *InstallationsController) HandleUpdate() {
	i.ClassName = "_Installation"
	i.ObjectID = i.Ctx.Input.Param(":objectId")
	i.ObjectsController.HandleUpdate()
}

// HandleDelete ...
// @router /:objectId [delete]
func (i *InstallationsController) HandleDelete() {
	i.ClassName = "_Installation"
	i.ObjectID = i.Ctx.Input.Param(":objectId")
	i.ObjectsController.HandleDelete()
}

// Delete ...
// @router / [delete]
func (i *InstallationsController) Delete() {
	i.Controller.Delete()
}

// Put ...
// @router / [put]
func (i *InstallationsController) Put() {
	i.Controller.Put()
}
