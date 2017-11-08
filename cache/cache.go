package cache

import (
	"strings"

	"github.com/lfq7413/tomato/config"
)

// Role ...
var Role *SubCache

// User ...
var User *SubCache

var adapter Adapter

func init() {
	a := config.TConfig.CacheAdapter
	if a == "InMemory" {
		adapter = newInMemoryCacheAdapter(5)
	} else if a == "Redis" {
		adapter = newRedisCacheAdapter(config.TConfig.RedisAddress, config.TConfig.RedisPassword, 0)
	} else if a == "Null" {
		adapter = newNullMemoryCacheAdapter()
	} else {
		adapter = newInMemoryCacheAdapter(5)
	}
	Role = &SubCache{
		prefix: "role",
	}
	User = &SubCache{
		prefix: "user",
	}
}

var keySeparatorChar = ":"

func joinKeys(keys ...string) string {
	return strings.Join(keys, keySeparatorChar)
}

func get(key string) interface{} {
	cacheKey := joinKeys(config.TConfig.AppID, key)
	return adapter.get(cacheKey)
}

func put(key string, value interface{}, ttl int64) {
	cacheKey := joinKeys(config.TConfig.AppID, key)
	adapter.put(cacheKey, value, ttl)
}

func del(key string) {
	cacheKey := joinKeys(config.TConfig.AppID, key)
	adapter.del(cacheKey)
}

func clear() {
	adapter.clear()
}

// SubCache ...
type SubCache struct {
	prefix string
}

// Get ...
func (c *SubCache) Get(key string) interface{} {
	cacheKey := joinKeys(c.prefix, key)
	return get(cacheKey)
}

// Put ...
func (c *SubCache) Put(key string, value interface{}, ttl int64) {
	cacheKey := joinKeys(c.prefix, key)
	put(cacheKey, value, ttl)
}

// Del ...
func (c *SubCache) Del(key string) {
	cacheKey := joinKeys(c.prefix, key)
	del(cacheKey)
}

// Clear ...
func (c *SubCache) Clear() {
	clear()
}

// Adapter ...
type Adapter interface {
	get(key string) interface{}
	put(key string, value interface{}, ttl int64)
	del(key string)
	clear()
}

// InitCache 仅用于测试
func InitCache() {
	adapter = newInMemoryCacheAdapter(5)
	Role = &SubCache{
		prefix: "role",
	}
	User = &SubCache{
		prefix: "user",
	}
}
