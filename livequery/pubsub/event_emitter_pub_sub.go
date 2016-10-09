package pubsub

var emitter = NewEventEmitter()

// eventEmitterPublisher 使用 EventEmitter 实现的发布者
type eventEmitterPublisher struct {
	emitter *EventEmitter
}

func (p *eventEmitterPublisher) Publish(channel, message string) {
	p.emitter.Emit(channel, message)
}

// eventEmitterPublisher 使用 EventEmitter 实现的订阅者
type eventEmitterSubscriber struct {
	EventEmitter
	emitter       *EventEmitter
	subscriptions map[string]HandlerType
}

func (s *eventEmitterSubscriber) Subscribe(channel string) {
	var handler HandlerType
	handler = func(args ...string) {
		a := append([]string{channel}, args...)
		s.Emit("message", a...)
	}
	s.subscriptions[channel] = handler
	s.emitter.On(channel, handler)
}

func (s *eventEmitterSubscriber) Unsubscribe(channel string) {
	if handler, ok := s.subscriptions[channel]; ok {
		s.emitter.RemoveListener(channel, handler)
		delete(s.subscriptions, channel)
	}
}

func (s *eventEmitterSubscriber) On(channel string, listener HandlerType) {
	s.EventEmitter.On(channel, listener)
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
	s.events = map[string]map[int]HandlerType{}
	return s
}
