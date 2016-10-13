package pubsub

// CreatePublisher 创建发布者，当前支持 EventEmitter 、 Redis
func CreatePublisher(pubType, pubURL string) Publisher {
	if useRedis(pubType) {
		return createRedisPublisher(pubURL)
	}
	return createEventEmitterPublisher()
}

// CreateSubscriber 创建订阅者，当前支持 EventEmitter 、 Redis
func CreateSubscriber(subType, subURL string) Subscriber {
	if useRedis(subType) {
		return createRedisSubscriber(subURL)
	}
	return createEventEmitterSubscriber()
}

// useRedis 判断类型是否为 redis
func useRedis(pubType string) bool {
	if pubType == "Redis" {
		return true
	}
	return false
}

// HandlerType ...
type HandlerType func(args ...string)

// Publisher ...
type Publisher interface {
	// Publish 向指定通道发送消息。当有对象保存或删除时，由 tomato 调用
	// channel 当前支持的通道包括：afterSave、afterDelete
	// message json 字符串，格式如下：
	// {
	// 	"currentParseObject": {...},
	// 	"originalParseObject": {...}
	// }
	Publish(channel, message string)
}

// Subscriber ...
type Subscriber interface {
	// Subscribe 订阅指定通道。由 LiveQueryServer 在初始化时调用
	// channel 当前支持的通道包括：afterSave、afterDelete
	Subscribe(channel string)

	// Unsubscribe 取消订阅指定通道
	Unsubscribe(channel string)

	// On 设置 从指定通道接收到消息时 的 回调函数
	// 统一从 message 通道获取数据
	// 传入 listener 的参数包含两个：
	// 第一个为已订阅的通道名称，当前为 afterSave 或者 afterDelete
	// 第二个为实际的对象数据，格式如下：
	// {
	// 	"currentParseObject": {...},
	// 	"originalParseObject": {...}
	// }
	On(channel string, listener HandlerType)
}
