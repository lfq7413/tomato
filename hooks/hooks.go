package hooks

import (
	"github.com/lfq7413/tomato/cloud"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

const defaultHooksCollectionName = "_Hooks"

// Load ...
func Load() {
	hooks, _ := getHooks(types.M{}, types.M{})
	for _, v := range hooks {
		if hook := utils.M(v); hook != nil {
			addHookToTriggers(hook)
		}
	}
}

// GetFunction ...
func GetFunction(functionName string) (types.M, error) {
	results, err := getHooks(types.M{"functionName": functionName}, types.M{"limit": 1})
	if err != nil {
		return nil, err
	}
	if results == nil || len(results) != 1 {
		return nil, nil
	}
	return utils.M(results[0]), nil
}

// GetFunctions ...
func GetFunctions() (types.S, error) {
	results, err := getHooks(types.M{"functionName": types.M{"$exists": true}}, types.M{})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// GetTrigger ...
func GetTrigger(className, triggerName string) (types.M, error) {
	results, err := getHooks(types.M{"className": className, "triggerName": triggerName}, types.M{"limit": 1})
	if err != nil {
		return nil, err
	}
	if results == nil || len(results) != 1 {
		return nil, nil
	}
	return utils.M(results[0]), nil
}

// GetTriggers ...
func GetTriggers() (types.S, error) {
	results, err := getHooks(types.M{"className": types.M{"$exists": true}, "triggerName": types.M{"$exists": true}}, types.M{})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// DeleteFunction ...
func DeleteFunction(functionName string) error {
	cloud.RemoveFunction(functionName)
	return removeHooks(types.M{"functionName": functionName})
}

// DeleteTrigger ...
func DeleteTrigger(className, triggerName string) error {
	cloud.RemoveTrigger(triggerName, className)
	return removeHooks(types.M{"className": className, "triggerName": triggerName})
}

func getHooks(query, options types.M) (types.S, error) {
	results, err := orm.TomatoDBController.Find(defaultHooksCollectionName, query, options)
	if err != nil {
		return nil, err
	}
	for _, v := range results {
		if result := utils.M(v); result != nil {
			delete(result, "objectId")
		}
	}
	return results, nil
}

func removeHooks(query types.M) error {
	return orm.TomatoDBController.Destroy(defaultHooksCollectionName, query, types.M{})
}

func saveHook(hook types.M) (types.M, error) {
	var query types.M
	if hook["functionName"] != nil && hook["url"] != nil {
		query = types.M{
			"functionName": hook["functionName"],
		}
	} else if hook["triggerName"] != nil && hook["className"] != nil && hook["url"] != nil {
		query = types.M{
			"triggerName": hook["triggerName"],
			"className":   hook["className"],
		}
	} else {
		return nil, errs.E(errs.WebhookError, "invalid hook declaration")
	}

	return orm.TomatoDBController.Update(defaultHooksCollectionName, query, hook, types.M{"upsert": true}, false)
}

func addHookToTriggers(hook types.M) {
	if hook["className"] != nil {
		cloud.AddTrigger(utils.S(hook["triggerName"]), utils.S(hook["className"]), cloud.GetTriggerHandler(utils.S(hook["url"])))
	}
	cloud.AddFunction(utils.S(hook["functionName"]), cloud.GetFunctionHandler(utils.S(hook["url"])), nil)
}

func addHook(hook types.M) (types.M, error) {
	addHookToTriggers(hook)
	return saveHook(hook)
}

func createOrUpdateHook(aHook types.M) (types.M, error) {
	var hook types.M
	if aHook != nil && aHook["functionName"] != nil && aHook["url"] != nil {
		hook = types.M{
			"functionName": aHook["functionName"],
			"url":          aHook["url"],
		}
	} else if aHook != nil && aHook["className"] != nil && aHook["url"] != nil && aHook["triggerName"] != nil {
		hook = types.M{
			"className":   aHook["className"],
			"triggerName": aHook["triggerName"],
			"url":         aHook["url"],
		}
	} else {
		return nil, errs.E(errs.WebhookError, "invalid hook declaration")
	}

	return addHook(hook)
}

// CreateHook ...
func CreateHook(aHook types.M) (types.M, error) {
	if aHook["functionName"] != nil {
		result, _ := GetFunction(utils.S(aHook["functionName"]))
		if result != nil {
			return nil, errs.E(errs.WebhookError, "function name: "+utils.S(aHook["functionName"])+" already exits")
		}
		return createOrUpdateHook(aHook)
	} else if aHook["className"] != nil && aHook["triggerName"] != nil {
		result, _ := GetTrigger(utils.S(aHook["className"]), utils.S(aHook["triggerName"]))
		if result != nil {
			return nil, errs.E(errs.WebhookError, "class "+utils.S(aHook["className"])+" already has trigger "+utils.S(aHook["triggerName"]))
		}
		return createOrUpdateHook(aHook)
	}
	return nil, errs.E(errs.WebhookError, "invalid hook declaration")
}

// UpdateHook ...
func UpdateHook(aHook types.M) (types.M, error) {
	if aHook["functionName"] != nil {
		result, _ := GetFunction(utils.S(aHook["functionName"]))
		if result == nil {
			return nil, errs.E(errs.WebhookError, "no function named: "+utils.S(aHook["functionName"])+" is defined")
		}
		return createOrUpdateHook(aHook)
	} else if aHook["className"] != nil && aHook["triggerName"] != nil {
		result, _ := GetTrigger(utils.S(aHook["className"]), utils.S(aHook["triggerName"]))
		if result == nil {
			return nil, errs.E(errs.WebhookError, "class "+utils.S(aHook["className"])+" does not exist")
		}
		return createOrUpdateHook(aHook)
	}
	return nil, errs.E(errs.WebhookError, "invalid hook declaration")
}
