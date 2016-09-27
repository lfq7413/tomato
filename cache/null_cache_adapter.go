package cache

type nullCacheAdapter struct {
}

func newNullMemoryCacheAdapter() *nullCacheAdapter {
	m := &nullCacheAdapter{}
	return m
}

func (m *nullCacheAdapter) get(key string) interface{} {
	return nil
}

func (m *nullCacheAdapter) put(key string, value interface{}, ttl int64) {
}

func (m *nullCacheAdapter) del(key string) {
}

func (m *nullCacheAdapter) clear() {
}
