package cloud

import "errors"

var restrictedClassNames = []string{"_Session"}

func validateClassNameForTriggers(className string) error {
	for _, v := range restrictedClassNames {
		if v == className {
			return errors.New("Triggers are not supported for " + className + " class.")
		}
	}
	return nil
}

// Define ...
func Define(functionName string, handler FunctionHandler, validationHandler ValidatorHandler) {
	AddFunction(functionName, handler, validationHandler)
}

// BeforeSave ...
func BeforeSave(className string, handler TriggerHandler) error {
	err := validateClassNameForTriggers(className)
	if err != nil {
		return err
	}
	AddTrigger(TypeBeforeSave, className, handler)
	return nil
}

// BeforeDelete ...
func BeforeDelete(className string, handler TriggerHandler) error {
	err := validateClassNameForTriggers(className)
	if err != nil {
		return err
	}
	AddTrigger(TypeBeforeDelete, className, handler)
	return nil
}

// AfterSave ...
func AfterSave(className string, handler TriggerHandler) error {
	err := validateClassNameForTriggers(className)
	if err != nil {
		return err
	}
	AddTrigger(TypeAfterSave, className, handler)
	return nil
}

// AfterDelete ...
func AfterDelete(className string, handler TriggerHandler) error {
	err := validateClassNameForTriggers(className)
	if err != nil {
		return err
	}
	AddTrigger(TypeAfterDelete, className, handler)
	return nil
}

// RemoveHook ...
func RemoveHook(category, name, triggerType string) {
	Unregister(category, name, triggerType)
}

// RemoveAllHooks ...
func RemoveAllHooks() {
	UnregisterAll()
}
