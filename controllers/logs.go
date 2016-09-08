package controllers

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/logger"
)

// LogsController ...
type LogsController struct {
	ClassesController
}

// HandleGet ...
// @router / [get]
func (l *LogsController) HandleGet() {
	if l.Auth.IsMaster == false {
		l.Data["json"] = errs.ErrorMessageToMap(errs.OperationForbidden, "Need master key!")
		l.ServeJSON()
		return
	}

	from := l.Ctx.Input.Param("from")
	until := l.Ctx.Input.Param("until")
	size := l.Ctx.Input.Param("size")
	if n := l.Ctx.Input.Param("size"); n != "" {
		size = n
	}
	order := l.Ctx.Input.Param("order")
	level := l.Ctx.Input.Param("level")

	options := map[string]string{
		"from":  from,
		"until": until,
		"size":  size,
		"order": order,
		"level": level,
	}
	result, err := logger.GetLogs(options)
	if err != nil {
		l.HandleError(err, 0)
		return
	}
	l.Data["json"] = result
	l.ServeJSON()
}

// Post ...
// @router / [post]
func (l *LogsController) Post() {
	l.ClassesController.Post()
}

// Delete ...
// @router / [delete]
func (l *LogsController) Delete() {
	l.ClassesController.Delete()
}

// Put ...
// @router / [put]
func (l *LogsController) Put() {
	l.ClassesController.Put()
}
