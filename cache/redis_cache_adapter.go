package cache

import (
	"encoding/json"

	"github.com/garyburd/redigo/redis"
)

type redisCacheAdapter struct {
	conn redis.Conn
}

func newRedisMemoryCacheAdapter() *redisCacheAdapter {
	c, err := redis.Dial("tcp", "192.168.99.100:6379")
	if err != nil {
		panic(err)
	}
	m := &redisCacheAdapter{
		conn: c,
	}
	return m
}

func (m *redisCacheAdapter) get(key string) interface{} {
	v, _ := m.conn.Do("GET", key)
	if v == nil {
		return v
	}
	var value interface{}
	json.Unmarshal(v.([]byte), &value)
	return value
}

func (m *redisCacheAdapter) put(key string, value interface{}, ttl int64) {
	v, _ := json.Marshal(value)
	m.conn.Do("SET", key, v)
}

func (m *redisCacheAdapter) del(key string) {
	m.conn.Do("DEL", key)
}

func (m *redisCacheAdapter) clear() {
	m.conn.Do("FLUSHALL")
}

func (m *redisCacheAdapter) close() {
	m.conn.Close()
}
