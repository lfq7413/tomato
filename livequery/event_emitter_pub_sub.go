package livequery

type eventEmitterPublisher struct {
}

func (p *eventEmitterPublisher) publish(channel, message string) {

}

type eventEmitterSubscriber struct {
}

func (s *eventEmitterSubscriber) subscribe(channel string) {

}
func (s *eventEmitterSubscriber) unsubscribe(channel string) {

}

func createEventEmitterPublisher() *eventEmitterPublisher {
	return &eventEmitterPublisher{}
}

func createEventEmitterSubscriber() *eventEmitterSubscriber {
	return &eventEmitterSubscriber{}
}
