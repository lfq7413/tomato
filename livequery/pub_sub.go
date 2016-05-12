package livequery

// createPublisher 创建发布者，当前仅支持 EventEmitter
func createPublisher(pubType, pubURL string) publisher {
	if useRedis(pubType) {
		// TODO 后期添加 Redis 支持
		return createEventEmitterPublisher()
	}
	return createEventEmitterPublisher()
}

// createSubscriber 创建订阅者，当前仅支持 EventEmitter
func createSubscriber(subType, subURL string) subscriber {
	if useRedis(subType) {
		// TODO 后期添加 Redis 支持
		return createEventEmitterSubscriber()
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

type publisher interface {
	publish(channel, message string)
}

type subscriber interface {
	subscribe(channel string)
	unsubscribe(channel string)
}
