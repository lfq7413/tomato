package controllers

// InstallationsController 处理 /installations 接口的请求
type InstallationsController struct {
	ClassesController
}

// HandleFind 处理查找已安装设备请求
// @router / [get]
func (i *InstallationsController) HandleFind() {
	i.ClassName = "_Installation"
	i.ClassesController.HandleFind()
}

// HandleGet 处理获取指定设备信息请求
// @router /:objectId [get]
func (i *InstallationsController) HandleGet() {
	i.ClassName = "_Installation"
	i.ObjectID = i.Ctx.Input.Param(":objectId")
	i.ClassesController.HandleGet()
}

// HandleCreate 处理添加设备请求
// @router / [post]
func (i *InstallationsController) HandleCreate() {
	i.ClassName = "_Installation"
	i.ClassesController.HandleCreate()
}

// HandleUpdate 处理更新设备信息请求
// @router /:objectId [put]
func (i *InstallationsController) HandleUpdate() {
	i.ClassName = "_Installation"
	i.ObjectID = i.Ctx.Input.Param(":objectId")
	i.ClassesController.HandleUpdate()
}

// HandleDelete 处理删除指定设备请求
// @router /:objectId [delete]
func (i *InstallationsController) HandleDelete() {
	i.ClassName = "_Installation"
	i.ObjectID = i.Ctx.Input.Param(":objectId")
	i.ClassesController.HandleDelete()
}

// Delete ...
// @router / [delete]
func (i *InstallationsController) Delete() {
	i.ClassesController.Delete()
}

// Put ...
// @router / [put]
func (i *InstallationsController) Put() {
	i.ClassesController.Put()
}
