package livequery

import (
	"encoding/json"

	"github.com/lfq7413/tomato/types"
)

type cloudCodePublisher struct {
	publisher publisher
}

func newCloudCodePublisher(pubType, pubURL string) *cloudCodePublisher {
	// TODO 后期添加真实的 Publisher
	return &cloudCodePublisher{
		publisher: createPublisher(pubType, pubURL),
	}
}

func (c *cloudCodePublisher) onCloudCodeAfterSave(request types.M) {
	c.onCloudCodeMessage("afterSave", request)
}

func (c *cloudCodePublisher) onCloudCodeAfterDelete(request types.M) {
	c.onCloudCodeMessage("afterDelete", request)
}

func (c *cloudCodePublisher) onCloudCodeMessage(messageType string, request types.M) {
	message := types.M{
		"currentParseObject": request["object"],
	}
	if request["original"] != nil {
		message["originalParseObject"] = request["original"]
	}
	res, err := json.Marshal(message)
	if err != nil {
		return
	}
	c.publisher.publish(messageType, string(res))
}
