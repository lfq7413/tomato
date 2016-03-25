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
	a.Controller.Get()
}

// Post ...
// @router / [post]
func (a *AnalyticsController) Post() {
	a.Controller.Post()
}

// Delete ...
// @router / [delete]
func (a *AnalyticsController) Delete() {
	a.Controller.Delete()
}

// Put ...
// @router / [put]
func (a *AnalyticsController) Put() {
	a.Controller.Put()
}
