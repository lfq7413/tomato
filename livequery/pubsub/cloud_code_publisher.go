package pubsub

import (
	"encoding/json"

	"github.com/lfq7413/tomato/livequery/server"
	"github.com/lfq7413/tomato/livequery/t"
)

// CloudCodePublisher 云代码发布者，当前支持发布 afterSave 与 afterDelete 通知
type CloudCodePublisher struct {
	publisher Publisher
}

// NewCloudCodePublisher 创建云代码发布者，其中的 Publisher 当前仅支持实验性质的 EventEmitter
func NewCloudCodePublisher(pubType, pubURL, pubConfig string) *CloudCodePublisher {
	return &CloudCodePublisher{
		publisher: CreatePublisher(pubType, pubURL, pubConfig),
	}
}

// OnCloudCodeAfterSave 对象保存时调用，request 中包含修改前与修改后的数据
func (c *CloudCodePublisher) OnCloudCodeAfterSave(request t.M) {
	c.onCloudCodeMessage(server.TomatoInfo["appId"]+"afterSave", request)
}

// OnCloudCodeAfterDelete 对象删除时调用，request 中包含要删除的数据
func (c *CloudCodePublisher) OnCloudCodeAfterDelete(request t.M) {
	c.onCloudCodeMessage(server.TomatoInfo["appId"]+"afterDelete", request)
}

// onCloudCodeMessage 向发送者发送通知消息
// 组装之后的 message 为 JSON 格式：
// {
// 	"currentParseObject": {...},
// 	"originalParseObject": {...}
// }
func (c *CloudCodePublisher) onCloudCodeMessage(messageType string, request t.M) {
	message := t.M{
		"currentParseObject": request["object"],
	}
	if request["original"] != nil {
		message["originalParseObject"] = request["original"]
	}
	res, err := json.Marshal(message)
	if err != nil {
		return
	}
	c.publisher.Publish(messageType, string(res))
}
