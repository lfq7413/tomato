package cloud

import "github.com/lfq7413/tomato/types"

// RemoteDefine ...
func RemoteDefine(functionName string, functionHandlerURL, validatorHandlerURL string) {
	Define(functionName, GetFunctionHandler(functionHandlerURL), GetValidatorHandler(validatorHandlerURL))
}

// RemoteBeforeSave ...
func RemoteBeforeSave(className string, triggerHandlerURL string) error {
	return BeforeSave(className, GetTriggerHandler(triggerHandlerURL))
}

// RemoteBeforeDelete ...
func RemoteBeforeDelete(className string, triggerHandlerURL string) error {
	return BeforeDelete(className, GetTriggerHandler(triggerHandlerURL))
}

// RemoteAfterSave ...
func RemoteAfterSave(className string, triggerHandlerURL string) error {
	return AfterSave(className, GetTriggerHandler(triggerHandlerURL))
}

// RemoteAfterDelete ...
func RemoteAfterDelete(className string, triggerHandlerURL string) error {
	return AfterDelete(className, GetTriggerHandler(triggerHandlerURL))
}

// GetFunctionHandler ...
func GetFunctionHandler(url string) FunctionHandler {
	return func(request FunctionRequest, response Response) {
		params := types.M{
			"params":         request.Params,
			"master":         request.Master,
			"user":           request.User,
			"installationID": request.InstallationID,
			"headers":        request.Headers,
		}
		result, err := post(params, url)
		if err != nil {
			response.Error(err["code"].(int), err["message"].(string))
			return
		}
		response.Success(result["result"])
	}
}

// GetValidatorHandler ...
func GetValidatorHandler(url string) ValidatorHandler {
	return func(request FunctionRequest) bool {
		params := types.M{
			"params":         request.Params,
			"master":         request.Master,
			"user":           request.User,
			"installationID": request.InstallationID,
			"headers":        request.Headers,
		}
		result, _ := post(params, url)
		if v, ok := result["result"].(bool); ok {
			return v
		}
		return false
	}
}

// GetTriggerHandler ...
func GetTriggerHandler(url string) TriggerHandler {
	return func(request TriggerRequest, response Response) {
		params := types.M{
			"triggerName":    request.TriggerName,
			"object":         request.Object,
			"original":       request.Original,
			"master":         request.Master,
			"user":           request.User,
			"installationID": request.InstallationID,
		}
		result, err := post(params, url)
		if err != nil {
			response.Error(err["code"].(int), err["message"].(string))
			return
		}
		if request.TriggerName == TypeBeforeSave {
			delete(result, "createdAt")
			delete(result, "updatedAt")
			request.Object = result
		}
		response.Success(nil)
	}
}
