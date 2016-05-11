package livequery

import "github.com/lfq7413/tomato/types"

type cloudCodePublisher struct {
}

func newCloudCodePublisher() *cloudCodePublisher {
	return &cloudCodePublisher{}
}

func (c *cloudCodePublisher) onCloudCodeAfterSave(request types.M) {

}

func (c *cloudCodePublisher) onCloudCodeAfterDelete(request types.M) {

}

func (c *cloudCodePublisher) onCloudCodeMessage(messageType string, request types.M) {

}
