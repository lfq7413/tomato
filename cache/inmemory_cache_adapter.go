package cache

import "time"

// TODO 增加定时清理过期缓存的操作

type inMemoryCacheAdapter struct {
	ttl   int64
	cache map[string]*recordCache
}

const defaultCacheTTL = 5 * 1000

func newInMemoryCacheAdapter(ttl int64) *inMemoryCacheAdapter {
	if ttl == 0 {
		ttl = defaultCacheTTL
	}
	m := &inMemoryCacheAdapter{
		ttl:   ttl,
		cache: map[string]*recordCache{},
	}
	return m
}

func (m *inMemoryCacheAdapter) get(key string) interface{} {
	if record, ok := m.cache[key]; ok {
		if record.expire >= time.Now().UnixNano() {
			return record.value
		}
		delete(m.cache, key)
		return nil
	}
	return nil
}

// put ttl 的单位为毫秒，为 0 时表示使用默认的时长
func (m *inMemoryCacheAdapter) put(key string, value interface{}, ttl int64) {
	if ttl == 0 {
		ttl = m.ttl
	}

	record := &recordCache{
		value:  value,
		expire: ttl*10e6 + time.Now().UnixNano(),
	}

	m.cache[key] = record
}

func (m *inMemoryCacheAdapter) del(key string) {
	delete(m.cache, key)
}

func (m *inMemoryCacheAdapter) clear() {
	m.cache = map[string]*recordCache{}
}

type recordCache struct {
	expire int64
	value  interface{}
}
