package livequery

var emitter = NewEventEmitter()

type eventEmitterPublisher struct {
	emitter *EventEmitter
}

func (p *eventEmitterPublisher) publish(channel, message string) {
	p.emitter.Emit(channel, message)
}

type eventEmitterSubscriber struct {
	EventEmitter
	emitter       *EventEmitter
	subscriptions map[string]HandlerType
}

func (s *eventEmitterSubscriber) subscribe(channel string) {
	var handler HandlerType
	handler = func(args ...string) {
		a := append([]string{channel}, args...)
		s.Emit("message", a...)
	}
	s.subscriptions[channel] = handler
	s.emitter.On(channel, handler)
}

func (s *eventEmitterSubscriber) unsubscribe(channel string) {
	if handler, ok := s.subscriptions[channel]; ok {
		s.emitter.RemoveListener(channel, handler)
		delete(s.subscriptions, channel)
	}
}

func createEventEmitterPublisher() *eventEmitterPublisher {
	return &eventEmitterPublisher{
		emitter: emitter,
	}
}

func createEventEmitterSubscriber() *eventEmitterSubscriber {
	s := &eventEmitterSubscriber{
		emitter:       emitter,
		subscriptions: map[string]HandlerType{},
	}
	s.events = map[string][]HandlerType{}
	return s
}
