package rest

import (
	"reflect"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// TypeBeforeSave 保存前回调
const TypeBeforeSave = "beforeSave"

// TypeAfterSave 保存后回调
const TypeAfterSave = "afterSave"

// TypeBeforeDelete 删除前回调
const TypeBeforeDelete = "beforeDelete"

// TypeAfterDelete 删除后回调
const TypeAfterDelete = "afterDelete"

// TriggerRequest ...
type TriggerRequest struct {
	TriggerName    string
	Object         types.M
	Original       types.M
	Master         bool
	User           types.M
	InstallationID string
}

// FunctionRequest ...
type FunctionRequest struct {
	Params         types.M
	Master         bool
	User           types.M
	InstallationID string
	Headers        map[string]string
}

// JobRequest ...
type JobRequest struct {
	Params types.M
	Master bool
	User   types.M
}

// Response ...
type Response interface {
	Success(response interface{})
	Error(code int, message string)
}

// TriggerHandler ...
type TriggerHandler func(TriggerRequest, Response)

// FunctionHandler ...
type FunctionHandler func(FunctionRequest, Response)

// ValidatorHandler ...
type ValidatorHandler func(FunctionRequest) bool

// JobHandler ...
type JobHandler func(JobRequest)

var triggers map[string]map[string]TriggerHandler
var functions map[string]FunctionHandler
var validators map[string]ValidatorHandler
var jobs map[string]JobHandler

func init() {
	triggers = map[string]map[string]TriggerHandler{
		TypeBeforeSave:   map[string]TriggerHandler{},
		TypeAfterSave:    map[string]TriggerHandler{},
		TypeBeforeDelete: map[string]TriggerHandler{},
		TypeAfterDelete:  map[string]TriggerHandler{},
	}
	functions = map[string]FunctionHandler{}
	validators = map[string]ValidatorHandler{}
	jobs = map[string]JobHandler{}
}

// AddFunction 添加函数到列表
func AddFunction(name string, handler FunctionHandler, validationHandler ValidatorHandler) {
	functions[name] = handler
	validators[name] = validationHandler
}

// AddJob 添加任务到列表
func AddJob(name string, handler JobHandler) {
	jobs[name] = handler
}

// AddTrigger 添加回调函数
func AddTrigger(triggerType string, className string, handler TriggerHandler) {
	triggers[triggerType][className] = handler
}

// RemoveFunction 从列表删除函数
func RemoveFunction(name string) {
	delete(functions, name)
	delete(validators, name)
}

// RemoveJob 从列表删除定时任务
func RemoveJob(name string) {
	delete(jobs, name)
}

// RemoveTrigger 从列表删除回调函数
func RemoveTrigger(triggerType string, className string) {
	delete(triggers[triggerType], className)
}

// Unregister 删除指定的云代码
func Unregister(category, name, triggerType string) {
	if category == "triggers" {
		if name != "" {
			delete(triggers[triggerType], name)
		} else {
			triggers[triggerType] = map[string]TriggerHandler{}
		}
	} else if category == "functions" {
		delete(functions, name)
	} else if category == "validators" {
		delete(validators, name)
	} else if category == "jobs" {
		delete(jobs, name)
	}
}

// UnregisterAll 删除所有注册的云代码
func UnregisterAll() {
	triggers = map[string]map[string]TriggerHandler{
		TypeBeforeSave:   map[string]TriggerHandler{},
		TypeAfterSave:    map[string]TriggerHandler{},
		TypeBeforeDelete: map[string]TriggerHandler{},
		TypeAfterDelete:  map[string]TriggerHandler{},
	}
	functions = map[string]FunctionHandler{}
	validators = map[string]ValidatorHandler{}
	jobs = map[string]JobHandler{}
}

// GetTrigger 获取回调函数
func GetTrigger(triggerType string, className string) TriggerHandler {
	if triggers == nil {
		return nil
	}
	if _, ok := triggers[triggerType]; ok == false {
		return nil
	}
	if v, ok := triggers[triggerType][className]; ok {
		return v
	}
	return nil
}

// TriggerExists 判断指定的回调函数是否存在
func TriggerExists(triggerType string, className string) bool {
	return GetTrigger(triggerType, className) != nil
}

// GetFunction 获取函数
func GetFunction(name string) FunctionHandler {
	if functions == nil {
		return nil
	}
	if v, ok := functions[name]; ok {
		return v
	}
	return nil
}

// GetValidator 获取校验函数
func GetValidator(name string) ValidatorHandler {
	if validators == nil {
		return nil
	}
	if v, ok := validators[name]; ok {
		return v
	}
	return nil
}

// GetJob 获取定时任务
func GetJob(name string) JobHandler {
	if jobs == nil {
		return nil
	}
	if v, ok := jobs[name]; ok {
		return v
	}
	return nil
}

type triggerResponse struct {
	request  TriggerRequest
	response types.M
	err      error
}

func (t *triggerResponse) Success(response interface{}) {
	t.response = utils.M(response)
	if t.response != nil &&
		reflect.DeepEqual(t.response, t.request.Object) == false &&
		t.request.TriggerName == TypeBeforeSave {
		return
	}
	t.response = types.M{}
	if t.request.TriggerName == TypeBeforeSave {
		t.response["object"] = t.request.Object
	}
}

func (t *triggerResponse) Error(code int, message string) {
	if code == 0 {
		code = errs.ScriptFailed
	}
	t.err = errs.E(code, message)
}

func getRequest(triggerType string, auth *Auth, parseObject, originalParseObject types.M) TriggerRequest {
	request := TriggerRequest{
		TriggerName: triggerType,
		Object:      parseObject,
		Master:      false,
	}

	if originalParseObject != nil {
		request.Original = originalParseObject
	}

	if auth == nil {
		return request
	}
	request.Master = auth.IsMaster
	request.User = auth.User
	request.InstallationID = auth.InstallationID

	return request
}

func getResponse(request TriggerRequest) *triggerResponse {
	response := &triggerResponse{
		request: request,
	}
	return response
}

func maybeRunTrigger(triggerType string, auth *Auth, parseObject, originalParseObject types.M) (types.M, error) {
	if parseObject == nil {
		return types.M{}, nil
	}

	trigger := GetTrigger(triggerType, utils.S(parseObject["className"]))
	if trigger == nil {
		return types.M{}, nil
	}
	request := getRequest(triggerType, auth, parseObject, originalParseObject)
	response := getResponse(request)
	trigger(request, response)
	return response.response, response.err
}

func inflate(data, restObject types.M) types.M {
	result := types.M{}
	if data != nil {
		for k, v := range data {
			result[k] = v
		}
	}
	if restObject != nil {
		for k, v := range restObject {
			result[k] = v
		}
	}

	return result
}
