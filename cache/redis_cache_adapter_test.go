package cache

import (
	"reflect"
	"testing"

	"github.com/lfq7413/tomato/types"
)

func Test_redis(t *testing.T) {
	var v interface{}
	cache := newRedisCacheAdapter("192.168.99.100:6379", "", 0)
	/*******************************************************************/
	cache.put("k1", "hello", 0)
	v = "hello"
	if reflect.DeepEqual(v, cache.get("k1")) == false {
		t.Error("get k1:", cache.get("k1"))
	}
	cache.del("k1")
	v = nil
	if reflect.DeepEqual(v, cache.get("k1")) == false {
		t.Error("get k1:", cache.get("k1"))
	}
	/*******************************************************************/
	cache.put("key1", 10, 0)
	v = 10.0
	if reflect.DeepEqual(v, cache.get("key1")) == false {
		t.Error("get key1:", cache.get("key1"))
	}
	/*******************************************************************/
	cache.put("key2", true, 0)
	v = true
	if reflect.DeepEqual(v, cache.get("key2")) == false {
		t.Error("get key2:", cache.get("key2"))
	}
	/*******************************************************************/
	cache.put("key3", []string{"a", "b"}, 0)
	v = []interface{}{"a", "b"}
	if reflect.DeepEqual(v, cache.get("key3")) == false {
		t.Error("get key3:", cache.get("key3"))
	}
	/*******************************************************************/
	cache.put("key4", map[string]bool{"key": true}, 0)
	v = map[string]interface{}{"key": true}
	if reflect.DeepEqual(v, cache.get("key4")) == false {
		t.Error("get key4:", cache.get("key4"))
	}
	/*******************************************************************/
	cache.put("key5", []types.M{types.M{"key": "types"}}, 0)
	v = []interface{}{map[string]interface{}{"key": "types"}}
	if reflect.DeepEqual(v, cache.get("key5")) == false {
		t.Error("get key5:", cache.get("key5"))
	}
	/*******************************************************************/
	cache.put("k2", map[string]interface{}{"key": "hello"}, 0)
	cache.put("k3", "hello world", 0)
	v = map[string]interface{}{"key": "hello"}
	if reflect.DeepEqual(v, cache.get("k2")) == false {
		t.Error("get k2:", cache.get("k2"))
	}
	v = "hello world"
	if reflect.DeepEqual(v, cache.get("k3")) == false {
		t.Error("get k3:", cache.get("k3"))
	}
	cache.clear()
	v = nil
	if reflect.DeepEqual(v, cache.get("k2")) == false {
		t.Error("get k2:", cache.get("k2"))
	}
	if reflect.DeepEqual(v, cache.get("k3")) == false {
		t.Error("get k3:", cache.get("k3"))
	}
}
