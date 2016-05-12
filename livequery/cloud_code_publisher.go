package livequery

import (
	"encoding/json"

	"github.com/lfq7413/tomato/types"
)

// cloudCodePublisher 云代码发布者，当前支持发布 afterSave 与 afterDelete 通知
type cloudCodePublisher struct {
	publisher publisher
}

// newCloudCodePublisher 创建云代码发布者，其中的 Publisher 当前仅支持实验性质的 EventEmitter
func newCloudCodePublisher(pubType, pubURL string) *cloudCodePublisher {
	// TODO 后期添加更多 Publisher
	return &cloudCodePublisher{
		publisher: createPublisher(pubType, pubURL),
	}
}

// onCloudCodeAfterSave 对象保存时调用，request 中包含修改前与修改后的数据
func (c *cloudCodePublisher) onCloudCodeAfterSave(request types.M) {
	c.onCloudCodeMessage("afterSave", request)
}

// onCloudCodeAfterDelete 对象删除时调用，request 中包含要删除的数据
func (c *cloudCodePublisher) onCloudCodeAfterDelete(request types.M) {
	c.onCloudCodeMessage("afterDelete", request)
}

// onCloudCodeMessage 向发送者发送通知消息
// 组装之后的 message 为 JSON 格式：
// {
// 	"currentParseObject": {...},
// 	"originalParseObject": {...}
// }
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
