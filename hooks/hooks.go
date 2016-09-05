package hooks

import "github.com/lfq7413/tomato/types"

// Load ...
func Load() {

}

// GetFunction ...xxx
func GetFunction(functionName string) (types.M, error) {
	return nil, nil
}

// GetFunctions ...xxx
func GetFunctions() (types.S, error) {
	return nil, nil
}

// GetTrigger ...xxx
func GetTrigger(className, triggerName string) (types.M, error) {
	return nil, nil
}

// GetTriggers ...xxx
func GetTriggers() (types.S, error) {
	return nil, nil
}

// DeleteFunction ...xxx
func DeleteFunction(functionName string) error {
	return nil
}

// DeleteTrigger ...xxx
func DeleteTrigger(className, triggerName string) error {
	return nil
}

func getHooks(query, options types.M) (types.S, error) {
	return nil, nil
}

func removeHooks(query types.M) error {
	return nil
}

func saveHook(hook types.M) (types.M, error) {
	return nil, nil
}

func addHookToTriggers(hook types.M) error {
	return nil
}

func addHook(hook types.M) (types.M, error) {
	return nil, nil
}

func createOrUpdateHook(aHook types.M) (types.M, error) {
	return nil, nil
}

// CreateHook ...xxx
func CreateHook(aHook types.M) (types.M, error) {
	return nil, nil
}

// UpdateHook ...xxx
func UpdateHook(aHook types.M) (types.M, error) {
	return nil, nil
}

func wrapToHTTPRequest(hook types.M, key string) {

}
