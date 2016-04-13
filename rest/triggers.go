package rest

import "github.com/lfq7413/tomato/types"

// TypeBeforeSave ...
const TypeBeforeSave = "beforeSave"

// TypeAfterSave ...
const TypeAfterSave = "afterSave"

// TypeBeforeDelete ...
const TypeBeforeDelete = "beforeDelete"

// TypeAfterDelete ...
const TypeAfterDelete = "afterDelete"

// RequestInfo ...
type RequestInfo struct {
	TriggerType    string
	Auth           *Auth
	NewObject      types.M
	OriginalObject types.M
}

// TypeFunction ...
type TypeFunction func(RequestInfo) types.M

type classeMap map[string]TypeFunction
type triggerMap map[string]classeMap
type functionMap map[string]TypeFunction

// Triggers 触发器列表
var Triggers triggerMap

// Functions 函数列表
var Functions functionMap

// Jobs 定时任务列表
var Jobs functionMap

func init() {
	Triggers = triggerMap{
		TypeBeforeSave:   classeMap{},
		TypeAfterSave:    classeMap{},
		TypeBeforeDelete: classeMap{},
		TypeAfterDelete:  classeMap{},
	}
	Functions = functionMap{}
	Jobs = functionMap{}
}

// AddFunction ...
func AddFunction(name string, function TypeFunction) {
	Functions[name] = function
}

// AddJob ...
func AddJob(name string, function TypeFunction) {
	Jobs[name] = function
}

// AddTrigger ...
func AddTrigger(triggerType string, className string, function TypeFunction) {
	Triggers[triggerType][className] = function
}

// RemoveFunction ...
func RemoveFunction(name string) {
	delete(Functions, name)
}

// RemoveJob ...
func RemoveJob(name string) {
	delete(Jobs, name)
}

// RemoveTrigger ...
func RemoveTrigger(triggerType string, className string) {
	delete(Triggers[triggerType], className)
}

// GetTrigger ...
func GetTrigger(triggerType string, className string) TypeFunction {
	if Triggers == nil {
		return nil
	}
	if _, ok := Triggers[triggerType]; ok == false {
		return nil
	}
	if v, ok := Triggers[triggerType][className]; ok {
		return v
	}
	return nil
}

// TriggerExists ...
func TriggerExists(triggerType string, className string) bool {
	return GetTrigger(triggerType, className) != nil
}

// GetFunction ...
func GetFunction(name string) TypeFunction {
	if Functions == nil {
		return nil
	}
	if v, ok := Functions[name]; ok {
		return v
	}
	return nil
}

// GetJob ...
func GetJob(name string) TypeFunction {
	if Jobs == nil {
		return nil
	}
	if v, ok := Jobs[name]; ok {
		return v
	}
	return nil
}

// RunTrigger ...
func RunTrigger(
	triggerType string,
	className string,
	auth *Auth,
	newObject types.M,
	originalObject types.M,
) types.M {
	trigger := GetTrigger(triggerType, className)
	if trigger == nil {
		return types.M{}
	}
	request := RequestInfo{
		TriggerType:    triggerType,
		Auth:           auth,
		NewObject:      newObject,
		OriginalObject: originalObject,
	}
	return trigger(request)
}
