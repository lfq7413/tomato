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
	listeners []HandlerType
}

func (r *redisSubscriber) Subscribe(channel string) {
	r.psc.Subscribe(channel)
}

func (r *redisSubscriber) Unsubscribe(channel string) {
	r.psc.Unsubscribe(channel)
}

func (r *redisSubscriber) On(channel string, listener HandlerType) {
	r.listeners = append(r.listeners, listener)
}

func (r *redisSubscriber) receive() {
	go func() {
		for {
			switch n := r.psc.Receive().(type) {
			case redis.Message:
				for _, listener := range r.listeners {
					go listener(n.Channel, string(n.Data))
				}
			}
		}
	}()
}

func createRedisPublisher(address, password string) *redisPublisher {
	c, err := redis.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	if password != "" {
		if _, err := c.Do("AUTH", password); err != nil {
			c.Close()
			panic(err)
		}
	}
	return &redisPublisher{
		c: c,
	}
}

func createRedisSubscriber(address, password string) *redisSubscriber {
	c, err := redis.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	if password != "" {
		if _, err := c.Do("AUTH", password); err != nil {
			c.Close()
			panic(err)
		}
	}
	r := &redisSubscriber{
		psc:       redis.PubSubConn{Conn: c},
		listeners: []HandlerType{},
	}
	r.receive()
	return r
}
