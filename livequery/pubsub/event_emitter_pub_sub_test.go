package pubsub

import (
	"fmt"
	"testing"
	"time"
)

func Test_EventEmitter(t *testing.T) {
	sub := createEventEmitterSubscriber()
	pub := createEventEmitterPublisher()
	sub.Subscribe("message")
	sub.On("message", func(args ...string) {
		fmt.Println("event emitter msg", args)
	})
	pub.Publish("message", "hello")
	time.Sleep(500 * time.Millisecond)
	sub.Unsubscribe("message")
}
