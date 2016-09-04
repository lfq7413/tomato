package controllers

// HooksController ...
type HooksController struct {
	ClassesController
}

// HandleGetAllFunctions ...
// @router /functions [get]
func (h *HooksController) HandleGetAllFunctions() {

}

// HandleGetFunction ...
// @router /functions/:functionName [get]
func (h *HooksController) HandleGetFunction() {

}

// HandleCreateFunction ...
// @router /functions [post]
func (h *HooksController) HandleCreateFunction() {

}

// HandleUpdateFunction ...
// @router /functions/:functionName [put]
func (h *HooksController) HandleUpdateFunction() {

}

// HandleGetAllTriggers ...
// @router /triggers [get]
func (h *HooksController) HandleGetAllTriggers() {

}

// HandleGetTrigger ...
// @router /triggers/:className/:triggerName [get]
func (h *HooksController) HandleGetTrigger() {

}

// HandleCreateTrigger ...
// @router /triggers [post]
func (h *HooksController) HandleCreateTrigger() {

}

// HandleUpdateTrigger ...
// @router /triggers/:className/:triggerName [put]
func (h *HooksController) HandleUpdateTrigger() {

}

// Get ...
// @router / [get]
func (h *HooksController) Get() {
	h.ClassesController.Get()
}

// Post ...
// @router / [post]
func (h *HooksController) Post() {
	h.ClassesController.Post()
}

// Delete ...
// @router / [delete]
func (h *HooksController) Delete() {
	h.ClassesController.Delete()
}

// Put ...
// @router / [put]
func (h *HooksController) Put() {
	h.ClassesController.Put()
}
