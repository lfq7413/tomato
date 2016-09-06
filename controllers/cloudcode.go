package controllers

import (
	"github.com/lfq7413/tomato/cloud"
	"github.com/lfq7413/tomato/types"
)

// CloudCodeController ...
type CloudCodeController struct {
	ClassesController
}

// HandleGet ...
// @router /jobs [get]
func (c *CloudCodeController) HandleGet() {
	jobs := cloud.GetJobs()
	jobNames := []string{}
	for n := range jobs {
		jobNames = append(jobNames, n)
	}
	c.Data["json"] = types.M{"jobName": jobNames}
	c.ServeJSON()
}

// Get ...
// @router / [get]
func (c *CloudCodeController) Get() {
	c.ClassesController.Get()
}

// Post ...
// @router / [post]
func (c *CloudCodeController) Post() {
	c.ClassesController.Post()
}

// Delete ...
// @router / [delete]
func (c *CloudCodeController) Delete() {
	c.ClassesController.Delete()
}

// Put ...
// @router / [put]
func (c *CloudCodeController) Put() {
	c.ClassesController.Put()
}
