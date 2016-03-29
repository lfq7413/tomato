package controllers

// AnalyticsController ...
type AnalyticsController struct {
	ObjectsController
}

// HandleAppOpened ...
// @router /AppOpened [post]
func (a *AnalyticsController) HandleAppOpened() {
	// TODO
}

// HandleEvent ...
// @router /:eventName [post]
func (a *AnalyticsController) HandleEvent() {
	// TODO
}

// Get ...
// @router / [get]
func (a *AnalyticsController) Get() {
	a.ObjectsController.Get()
}

// Post ...
// @router / [post]
func (a *AnalyticsController) Post() {
	a.ObjectsController.Post()
}

// Delete ...
// @router / [delete]
func (a *AnalyticsController) Delete() {
	a.ObjectsController.Delete()
}

// Put ...
// @router / [put]
func (a *AnalyticsController) Put() {
	a.ObjectsController.Put()
}
