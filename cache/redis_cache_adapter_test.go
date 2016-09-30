package cache

import (
	"fmt"
	"testing"
)

func Test_redis(t *testing.T) {
	cache := newRedisMemoryCacheAdapter()
	cache.put("k1", "hello", 0)
	fmt.Println("get k1:", cache.get("k1"))
	cache.del("k1")
	fmt.Println("get k1:", cache.get("k1"))

	cache.put("k2", map[string]interface{}{"key": "hello"}, 0)
	cache.put("k3", "hello world", 0)
	fmt.Println("get k2:", cache.get("k2"))
	fmt.Println("get k3:", cache.get("k3"))
	cache.clear()
	fmt.Println("get k2:", cache.get("k2"))
	fmt.Println("get k3:", cache.get("k3"))
	cache.close()
}
