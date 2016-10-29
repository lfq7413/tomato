package pubsub

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

type redisPublisher struct {
	address  string
	password string
	p        *redis.Pool
}

func (r *redisPublisher) Publish(channel, message string) {
	r.do("PUBLISH", channel, message)
}

func (r *redisPublisher) connectInit() {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", r.address)
		if err != nil {
			return nil, err
		}

		if r.password != "" {
			if _, err := c.Do("AUTH", r.password); err != nil {
				c.Close()
				return nil, err
			}
		}

		return
	}
	// initialize a new pool
	r.p = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}
}

func (r *redisPublisher) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c := r.p.Get()
	defer c.Close()

	return c.Do(commandName, args...)
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
	m := &redisPublisher{
		address:  address,
		password: password,
	}
	m.connectInit()
	c := m.p.Get()
	defer c.Close()
	if c.Err() != nil {
		panic(c.Err())
	}

	return m
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
