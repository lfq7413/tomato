package livequery

import "github.com/lfq7413/tomato/dependencies/lru"

type sessionTokenCache struct {
	cache *lru.Cache
}

func newSessionTokenCache() *sessionTokenCache {
	return &sessionTokenCache{
		cache: lru.New(10000),
	}
}

func (s *sessionTokenCache) getUserID(sessionToken string) string {
	// TODO
	if v, ok := s.cache.Get(sessionToken); ok {
		TLog.verbose("Fetch userId", v, "of sessionToken", sessionToken, "from Cache")
		return v.(string)
	}
	return ""
}
