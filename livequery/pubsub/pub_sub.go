package pubsub

// CreatePublisher 创建发布者，当前仅支持 EventEmitter
func CreatePublisher(pubType, pubURL string) Publisher {
	if useRedis(pubType) {
		return createRedisPublisher()
	}
	return createEventEmitterPublisher()
}

// CreateSubscriber 创建订阅者，当前仅支持 EventEmitter
func CreateSubscriber(subType, subURL string) Subscriber {
	if useRedis(subType) {
		return createRedisSubscriber()
	}
	return createEventEmitterSubscriber()
}

// useRedis 判断类型是否为 redis
func useRedis(pubType string) bool {
	if pubType == "redis" {
		return true
	}
	return false
}

// HandlerType ...
type HandlerType func(args ...string)

// Publisher ...
type Publisher interface {
	Publish(channel, message string)
}

// Subscriber ...
type Subscriber interface {
	Subscribe(channel string)
	Unsubscribe(channel string)
	On(channel string, listener HandlerType)
}
