package cache

import (
	"encoding/json"
	"sync"

	"github.com/garyburd/redigo/redis"
)

type redisCacheAdapter struct {
	conn redis.Conn
	mu   sync.Mutex
}

func newRedisMemoryCacheAdapter(address, password string) *redisCacheAdapter {
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
	m := &redisCacheAdapter{
		conn: c,
	}
	return m
}

func (m *redisCacheAdapter) get(key string) interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, _ := m.conn.Do("GET", key)
	if v == nil {
		return v
	}
	var value interface{}
	json.Unmarshal(v.([]byte), &value)
	return value
}

func (m *redisCacheAdapter) put(key string, value interface{}, ttl int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, _ := json.Marshal(value)
	m.conn.Do("SET", key, v)
}

func (m *redisCacheAdapter) del(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.conn.Do("DEL", key)
}

func (m *redisCacheAdapter) clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.conn.Do("FLUSHALL")
}

func (m *redisCacheAdapter) close() {
	m.conn.Close()
}
