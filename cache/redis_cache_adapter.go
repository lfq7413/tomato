package cache

import (
	"encoding/json"
	"time"

	"github.com/garyburd/redigo/redis"
)

type redisCacheAdapter struct {
	address  string
	password string
	p        *redis.Pool
}

func newRedisMemoryCacheAdapter(address, password string) *redisCacheAdapter {
	m := &redisCacheAdapter{
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

func (m *redisCacheAdapter) connectInit() {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", m.address)
		if err != nil {
			return nil, err
		}

		if m.password != "" {
			if _, err := c.Do("AUTH", m.password); err != nil {
				c.Close()
				return nil, err
			}
		}

		return
	}
	// initialize a new pool
	m.p = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}
}

func (m *redisCacheAdapter) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c := m.p.Get()
	defer c.Close()

	return c.Do(commandName, args...)
}

func (m *redisCacheAdapter) get(key string) interface{} {
	v, _ := m.do("GET", key)
	if v == nil {
		return v
	}
	var value interface{}
	json.Unmarshal(v.([]byte), &value)
	return value
}

func (m *redisCacheAdapter) put(key string, value interface{}, ttl int64) {
	v, _ := json.Marshal(value)
	m.do("SET", key, v)
}

func (m *redisCacheAdapter) del(key string) {
	m.do("DEL", key)
}

func (m *redisCacheAdapter) clear() {
	m.do("FLUSHALL")
}
