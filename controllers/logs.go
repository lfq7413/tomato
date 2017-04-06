package controllers

import "github.com/lfq7413/tomato/logger"

// LogsController ...
type LogsController struct {
	ClassesController
}

// HandleGet ...
// @router / [get]
func (l *LogsController) HandleGet() {
	if l.EnforceMasterKeyAccess() == false {
		return
	}

	from := l.Query["from"]
	until := l.Query["until"]
	size := l.Query["size"]
	if n := l.Query["n"]; n != "" {
		size = n
	}
	order := l.Query["order"]
	level := l.Query["level"]

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
