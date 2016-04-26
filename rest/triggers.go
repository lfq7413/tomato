package rest

import "github.com/lfq7413/tomato/types"

// TypeBeforeSave 保存前回调
const TypeBeforeSave = "beforeSave"

// TypeAfterSave 保存后回调
const TypeAfterSave = "afterSave"

// TypeBeforeDelete 删除前回调
const TypeBeforeDelete = "beforeDelete"

// TypeAfterDelete 删除后回调
const TypeAfterDelete = "afterDelete"

// RequestInfo 请求参数信息
// TriggerType 回调类型
// Auth 当前请求的权限信息
// NewObject 要保存或修改的数据，也可以时函数的传入参数
// OriginalObject 原始数据
type RequestInfo struct {
	TriggerType    string
	Auth           *Auth
	NewObject      types.M
	OriginalObject types.M
}

// TypeFunction 函数类型，返回数据封装到 types.M 中，
// 返回数据格式如下：
// {
// 	"object":{...}
// }
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

// AddFunction 添加函数到列表
func AddFunction(name string, function TypeFunction) {
	Functions[name] = function
}

// AddJob 添加任务到列表
func AddJob(name string, function TypeFunction) {
	Jobs[name] = function
}

// AddTrigger 添加回调函数
func AddTrigger(triggerType string, className string, function TypeFunction) {
	Triggers[triggerType][className] = function
}

// RemoveFunction 从列表删除函数
func RemoveFunction(name string) {
	delete(Functions, name)
}

// RemoveJob 从列表删除定时任务
func RemoveJob(name string) {
	delete(Jobs, name)
}

// RemoveTrigger 从列表删除回调函数
func RemoveTrigger(triggerType string, className string) {
	delete(Triggers[triggerType], className)
}

// GetTrigger 获取回调函数
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

// TriggerExists 判断指定的回调函数是否存在
func TriggerExists(triggerType string, className string) bool {
	return GetTrigger(triggerType, className) != nil
}

// GetFunction 获取函数
func GetFunction(name string) TypeFunction {
	if Functions == nil {
		return nil
	}
	if v, ok := Functions[name]; ok {
		return v
	}
	return nil
}

// GetJob 获取定时任务
func GetJob(name string) TypeFunction {
	if Jobs == nil {
		return nil
	}
	if v, ok := Jobs[name]; ok {
		return v
	}
	return nil
}

// RunTrigger 运行回调
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
