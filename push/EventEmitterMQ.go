package push

import "github.com/lfq7413/tomato/livequery/pubsub"

var emitter = pubsub.NewEventEmitter()
var subscriptions = map[string]pubsub.HandlerType{}

func unsubscribe(channel string) {
	if subscriptions[channel] == nil {
		return
	}
	emitter.RemoveListener(channel, subscriptions[channel])
	delete(subscriptions, channel)
}

// Publisher ...
type Publisher struct {
	emitter *pubsub.EventEmitter
}

// Publish ...
func (p *Publisher) Publish(channel string, message string) {
	p.emitter.Emit(channel, message)
}

// Consumer ...
type Consumer struct {
	pubsub.EventEmitter
	emitter *pubsub.EventEmitter
}

// Subscribe ...
func (c *Consumer) Subscribe(channel string) {
	unsubscribe(channel)
	var handler = func(args ...string) {
		allArgs := []string{channel}
		allArgs = append(allArgs, args...)
		c.Emit("message", allArgs...)
	}
	subscriptions[channel] = handler
	c.emitter.On(channel, handler)
}

// Unsubscribe ...
func (c *Consumer) Unsubscribe(channel string) {
	unsubscribe(channel)
}

// On ...
func (c *Consumer) On(channel string, listener pubsub.HandlerType) {
	c.EventEmitter.On(channel, listener)
}

// CreatePublisher ...
func CreatePublisher() pubsub.Publisher {
	return &Publisher{
		emitter: emitter,
	}
}

// CreateSubscriber ...
func CreateSubscriber() pubsub.Subscriber {
	c := &Consumer{
		emitter: emitter,
	}
	c.Init()
	return c
}
