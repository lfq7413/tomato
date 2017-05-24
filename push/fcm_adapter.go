package push

import (
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"

	"github.com/NaySoftware/go-fcm"
)

type fcmPushAdapter struct {
	validPushTypes []string
	serverKey      string
}

func newFCMPush() *fcmPushAdapter {
	f := &fcmPushAdapter{
		validPushTypes: []string{"ios", "osx", "tvos", "android", "fcm"},
		serverKey:      config.TConfig.FCMServerKey,
	}
	return f
}

func (f *fcmPushAdapter) send(body types.M, installations types.S, pushStatus string) []types.M {
	devices := []types.M{}
	deviceTokens := []string{}
	results := []types.M{}

	for _, installation := range installations {
		if dev := utils.M(installation); dev != nil {
			devices = append(devices, dev)
			deviceTokens = append(deviceTokens, utils.S(dev["deviceToken"]))
		}
	}

	status, err := f.sendToRegistrationTokens(deviceTokens, body)

	if err != nil {
		for _, device := range devices {
			result := types.M{
				"device":      device,
				"transmitted": false,
				"response":    map[string]string{"error": err.Error()},
			}
			results = append(results, result)
		}
		return results
	}

	multicastID := status.MulticastId
	pushResults := status.Results

	for index := range deviceTokens {
		var pushResult map[string]string
		if pushResults != nil && index < len(pushResults) {
			pushResult = pushResults[index]
		} else {
			pushResult = nil
		}
		device := devices[index]

		resolution := types.M{
			"device":       device,
			"multicast_id": multicastID,
			"response":     pushResult,
		}

		if pushResult == nil || pushResult["error"] != "" {
			resolution["transmitted"] = false
		} else {
			resolution["transmitted"] = true
		}

		results = append(results, resolution)
	}

	return results
}

func (f *fcmPushAdapter) getValidPushTypes() []string {
	return f.validPushTypes
}

func (f *fcmPushAdapter) sendToRegistrationTokens(tokens []string, body types.M) (*fcm.FcmResponseStatus, error) {
	// TODO 转换 body 中的数据到 FCM 支持的格式
	c := fcm.NewFcmClient(f.serverKey)
	c.NewFcmRegIdsMsg(tokens, body)

	return c.Send()
}

// classifyInstallations 对设备按照推送类型进行分类
func classifyInstallations(installations types.S, validPushTypes []string) map[string][]types.M {
	deviceMap := map[string][]types.M{}
	for _, validPushType := range validPushTypes {
		deviceMap[validPushType] = []types.M{}
	}

	for _, installation := range installations {
		if dev := utils.M(installation); dev != nil {
			deviceToken := utils.S(dev["deviceToken"])
			if deviceToken == "" {
				continue
			}

			pushType := utils.S(dev["pushType"])
			deviceType := utils.S(dev["deviceType"])
			var devices []types.M
			var tp string

			if d := deviceMap[pushType]; d != nil {
				devices = d
				tp = pushType
			} else if d := deviceMap[deviceType]; d != nil {
				devices = d
				tp = deviceType
			} else {
				devices = nil
			}

			if devices != nil {
				device := types.M{
					"deviceToken":   deviceToken,
					"deviceType":    deviceType,
					"appIdentifier": dev["appIdentifier"],
				}
				devices = append(devices, device)
				deviceMap[tp] = devices
			}
		}
	}

	return deviceMap
}
