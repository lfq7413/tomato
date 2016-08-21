package rest

import (
	"github.com/lfq7413/tomato/cloud"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

func getRequest(triggerType string, auth *Auth, parseObject, originalParseObject types.M) cloud.TriggerRequest {
	request := cloud.TriggerRequest{
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

func getResponse(request cloud.TriggerRequest) *cloud.TriggerResponse {
	response := &cloud.TriggerResponse{
		Request: request,
	}
	return response
}

func maybeRunTrigger(triggerType string, auth *Auth, parseObject, originalParseObject types.M) (types.M, error) {
	if parseObject == nil {
		return types.M{}, nil
	}

	trigger := cloud.GetTrigger(triggerType, utils.S(parseObject["className"]))
	if trigger == nil {
		return types.M{}, nil
	}
	request := getRequest(triggerType, auth, parseObject, originalParseObject)
	response := getResponse(request)
	trigger(request, response)
	return response.Response, response.Err
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
