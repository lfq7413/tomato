package pubsub

import (
	"fmt"
	"testing"
	"time"
)

func Test_EventEmitter(t *testing.T) {
	sub := createEventEmitterSubscriber()
	pub := createEventEmitterPublisher()
	sub.Subscribe("afterSave")
	sub.Subscribe("afterDelete")
	sub.On("message", func(args ...string) {
		fmt.Println("event emitter msg", args)
	})
	pub.Publish("afterSave", "hello")
	pub.Publish("afterDelete", "hello")
	time.Sleep(500 * time.Millisecond)
	sub.Unsubscribe("afterSave")
	sub.Unsubscribe("afterDelete")
}
