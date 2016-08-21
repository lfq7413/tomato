package cloud

// Define ...
func Define(functionName string, handler FunctionHandler, validationHandler ValidatorHandler) {
	AddFunction(functionName, handler, validationHandler)
}

// BeforeSave ...
func BeforeSave(className string, handler TriggerHandler) {
	AddTrigger(TypeBeforeSave, className, handler)
}

// BeforeDelete ...
func BeforeDelete(className string, handler TriggerHandler) {
	AddTrigger(TypeBeforeDelete, className, handler)
}

// AfterSave ...
func AfterSave(className string, handler TriggerHandler) {
	AddTrigger(TypeAfterSave, className, handler)
}

// AfterDelete ...
func AfterDelete(className string, handler TriggerHandler) {
	AddTrigger(TypeAfterDelete, className, handler)
}

// RemoveHook ...
func RemoveHook(category, name, triggerType string) {
	Unregister(category, name, triggerType)
}

// RemoveAllHooks ...
func RemoveAllHooks() {
	UnregisterAll()
}
