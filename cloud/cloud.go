package cloud

import "github.com/lfq7413/tomato/rest"

// Define ...
func Define(functionName string, handler rest.FunctionHandler, validationHandler rest.ValidatorHandler) {
	rest.AddFunction(functionName, handler, validationHandler)
}

// BeforeSave ...
func BeforeSave(className string, handler rest.TriggerHandler) {
	rest.AddTrigger(rest.TypeBeforeSave, className, handler)
}

// BeforeDelete ...
func BeforeDelete(className string, handler rest.TriggerHandler) {
	rest.AddTrigger(rest.TypeBeforeDelete, className, handler)
}

// AfterSave ...
func AfterSave(className string, handler rest.TriggerHandler) {
	rest.AddTrigger(rest.TypeAfterSave, className, handler)
}

// AfterDelete ...
func AfterDelete(className string, handler rest.TriggerHandler) {
	rest.AddTrigger(rest.TypeAfterDelete, className, handler)
}

// RemoveHook ...
func RemoveHook(category, name, triggerType string) {
	rest.Unregister(category, name, triggerType)
}

// RemoveAllHooks ...
func RemoveAllHooks() {
	rest.UnregisterAll()
}
