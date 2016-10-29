package pubsub

import (
	"fmt"
	"testing"
	"time"
)

func Test_redis(t *testing.T) {
	sub := createRedisSubscriber("192.168.99.100:6379", "")
	pub := createRedisPublisher("192.168.99.100:6379", "")
	sub.Subscribe("afterSave")
	sub.Subscribe("afterDelete")
	sub.On("message", func(args ...string) {
		fmt.Println("redis msg", args)
	})
	pub.Publish("afterSave", "hello")
	pub.Publish("afterDelete", "hello")
	time.Sleep(500 * time.Millisecond)
	sub.Unsubscribe("afterSave")
	sub.Unsubscribe("afterDelete")
}
