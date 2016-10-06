package pubsub

import (
	"fmt"
	"testing"
	"time"
)

func Test_redis(t *testing.T) {
	sub := createRedisSubscriber("192.168.99.100:6379")
	pub := createRedisPublisher("192.168.99.100:6379")
	sub.Subscribe("message")
	sub.On("message", func(args ...string) {
		fmt.Println(args)
	})
	pub.Publish("message", "hello")
	time.Sleep(500 * time.Millisecond)
	sub.Unsubscribe("message")
}
