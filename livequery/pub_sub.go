package livequery

func createPublisher(pubType, pubURL string) publisher {
	return nil
}

func createSubscriber(pubType, pubURL string) subscriber {
	return nil
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
