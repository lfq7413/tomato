package cloud

import (
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// RemoteDefine ...
func RemoteDefine(functionName string, functionHandlerURL, validatorHandlerURL string) {
	Define(functionName, getFunctionHandler(functionHandlerURL), getValidatorHandler(validatorHandlerURL))
}

// RemoteBeforeSave ...
func RemoteBeforeSave(className string, triggerHandlerURL string) error {
	return BeforeSave(className, getTriggerHandler(triggerHandlerURL))
}

// RemoteBeforeDelete ...
func RemoteBeforeDelete(className string, triggerHandlerURL string) error {
	return BeforeDelete(className, getTriggerHandler(triggerHandlerURL))
}

// RemoteAfterSave ...
func RemoteAfterSave(className string, triggerHandlerURL string) error {
	return AfterSave(className, getTriggerHandler(triggerHandlerURL))
}

// RemoteAfterDelete ...
func RemoteAfterDelete(className string, triggerHandlerURL string) error {
	return AfterDelete(className, getTriggerHandler(triggerHandlerURL))
}

func getFunctionHandler(url string) FunctionHandler {
	return func(request FunctionRequest, response Response) {
		params := types.M{
			"params":         request.Params,
			"master":         request.Master,
			"user":           request.User,
			"installationID": request.InstallationID,
			"headers":        request.Headers,
		}
		result := post(params, url)
		if result["code"] != nil && result["message"] != nil {
			var code int
			if v, ok := result["code"].(float64); ok {
				code = int(v)
			}
			var message string
			if v, ok := result["message"].(string); ok {
				message = v
			}
			response.Error(code, message)
			return
		}
		response.Success(result["result"])
	}
}

func getValidatorHandler(url string) ValidatorHandler {
	return func(request FunctionRequest) bool {
		params := types.M{
			"params":         request.Params,
			"master":         request.Master,
			"user":           request.User,
			"installationID": request.InstallationID,
			"headers":        request.Headers,
		}
		result := post(params, url)
		if v, ok := result["result"].(bool); ok {
			return v
		}
		return false
	}
}

func getTriggerHandler(url string) TriggerHandler {
	return func(request TriggerRequest, response Response) {
		params := types.M{
			"triggerName":    request.TriggerName,
			"object":         request.Object,
			"original":       request.Original,
			"master":         request.Master,
			"user":           request.User,
			"installationID": request.InstallationID,
		}
		result := post(params, url)
		if result["code"] != nil && result["message"] != nil {
			var code int
			if v, ok := result["code"].(float64); ok {
				code = int(v)
			}
			var message string
			if v, ok := result["message"].(string); ok {
				message = v
			}
			response.Error(code, message)
			return
		}
		if object := utils.M(result["object"]); object != nil {
			request.Object = object
		}
		response.Success(nil)
	}
}
