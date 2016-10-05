package pubsub

import "github.com/garyburd/redigo/redis"

type redisPublisher struct {
	c redis.Conn
}

func (r *redisPublisher) Publish(channel, message string) {
	r.c.Do("PUBLISH", channel, message)
}

type redisSubscriber struct {
	psc       redis.PubSubConn
	listeners map[string]HandlerType // TODO 增加并发锁
}

func (r *redisSubscriber) Subscribe(channel string) {
	r.psc.Subscribe(channel)
}

func (r *redisSubscriber) Unsubscribe(channel string) {
	r.psc.Unsubscribe(channel)
}

func (r *redisSubscriber) On(channel string, listener HandlerType) {
	r.listeners[channel] = listener
}

func (r *redisSubscriber) receive() {
	go func() {
		for {
			switch n := r.psc.Receive().(type) {
			case redis.Message:
				listener := r.listeners[n.Channel]
				if listener != nil {
					listener(n.Channel, string(n.Data))
				}
			}
		}
	}()
}

func createRedisPublisher(address string) *redisPublisher {
	c, err := redis.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	return &redisPublisher{
		c: c,
	}
}

func createRedisSubscriber(address string) *redisSubscriber {
	c, err := redis.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	r := &redisSubscriber{
		psc:       redis.PubSubConn{Conn: c},
		listeners: map[string]HandlerType{},
	}
	r.receive()
	return r
}
