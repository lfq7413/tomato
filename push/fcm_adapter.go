package push

import (
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"

	"time"

	"fmt"

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
	deviceMap := classifyInstallations(installations, f.validPushTypes)
	results := []types.M{}

	loop := func(pushType string) {
		devices := deviceMap[pushType]
		if len(devices) == 0 {
			return
		}
		deviceTokens := []string{}

		for _, device := range devices {
			deviceTokens = append(deviceTokens, utils.S(device["deviceToken"]))
		}

		var status *fcm.FcmResponseStatus
		var err error

		switch pushType {
		case "ios", "tvos", "osx":
			status, err = f.sendToiOSDevices(deviceTokens, body)
		case "android", "fcm":
			status, err = f.sendToAndroidDevices(deviceTokens, body)
		}

		if err != nil {
			for _, device := range devices {
				result := types.M{
					"device":      device,
					"transmitted": false,
					"response":    map[string]string{"error": err.Error()},
				}
				results = append(results, result)
			}
			return
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

	}

	for pushType := range deviceMap {
		loop(pushType)
	}

	return results
}

func (f *fcmPushAdapter) getValidPushTypes() []string {
	return f.validPushTypes
}

func (f *fcmPushAdapter) sendToiOSDevices(tokens []string, body types.M) (*fcm.FcmResponseStatus, error) {
	c := fcm.NewFcmClient(f.serverKey)
	c.SetPriority(fcm.Priority_HIGH)

	if t, ok := body["expiration_time"].(int64); ok {
		timeToLive := (t - time.Now().Unix()) / 1000
		if timeToLive < 0 {
			timeToLive = 0
		}
		c.SetTimeToLive(int(timeToLive))
	}

	pushData := utils.M(body["data"])
	if pushData == nil {
		pushData = types.M{}
	}
	payload := fcm.NotificationPayload{}
	data := types.M{}
	for key, v := range pushData {
		switch key {
		case "alert":
			payload.Body = utils.S(v)
		case "badge":
			payload.Badge = fmt.Sprintf("%v", v)
		case "sound":
			payload.Sound = utils.S(v)
		case "content-available":
			c.SetContentAvailable(true)
		case "category":
			payload.ClickAction = utils.S(v)
		case "uri":
		case "title":
			payload.Title = utils.S(v)
		default:
			data[key] = v
		}
	}
	c.SetNotificationPayload(&payload)

	c.NewFcmRegIdsMsg(tokens, data)
	return c.Send()
}

func (f *fcmPushAdapter) sendToAndroidDevices(tokens []string, body types.M) (*fcm.FcmResponseStatus, error) {
	c := fcm.NewFcmClient(f.serverKey)
	c.SetPriority(fcm.Priority_HIGH)

	if t, ok := body["expiration_time"].(int64); ok {
		timeToLive := (t - time.Now().Unix()) / 1000
		if timeToLive < 0 {
			timeToLive = 0
		}
		c.SetTimeToLive(int(timeToLive))
	}

	pushData := utils.M(body["data"])
	if pushData == nil {
		pushData = types.M{}
	}
	payload := fcm.NotificationPayload{}
	for key, v := range pushData {
		switch key {
		case "alert":
			payload.Body = utils.S(v)
		case "badge":
		case "sound":
			payload.Sound = utils.S(v)
		case "content-available":
		case "category":
		case "uri":
			payload.ClickAction = utils.S(v)
		case "title":
			payload.Title = utils.S(v)
		default:
		}
	}
	c.SetNotificationPayload(&payload)

	data := types.M{
		"data":    body["data"],
		"push_id": utils.CreateString(10),
		"time":    utils.TimetoString(time.Now().UTC()),
	}

	c.NewFcmRegIdsMsg(tokens, data)
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
