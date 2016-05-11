package livequery

func createPublisher(pubType, pubURL string) publisher {
	if useRedis(pubType) {
		return createEventEmitterPublisher()
	}
	return createEventEmitterPublisher()
}

func createSubscriber(subType, subURL string) subscriber {
	if useRedis(subType) {
		return createEventEmitterSubscriber()
	}
	return createEventEmitterSubscriber()
}

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
