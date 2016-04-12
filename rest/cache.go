package rest

// users 用户信息缓存
var usersCache *CacheStore

func init() {
	usersCache = &CacheStore{
		dataStore: map[string]interface{}{},
	}
}

func clearCache() {
	usersCache.clear()
}

// CacheStore 缓存结构体
type CacheStore struct {
	dataStore map[string]interface{}
}

func (c *CacheStore) get(key string) interface{} {
	return c.dataStore[key]
}

func (c *CacheStore) set(key string, value interface{}) {
	c.dataStore[key] = value
}

func (c *CacheStore) remove(key string) {
	delete(c.dataStore, key)
}

func (c *CacheStore) clear() {
	c.dataStore = map[string]interface{}{}
}
