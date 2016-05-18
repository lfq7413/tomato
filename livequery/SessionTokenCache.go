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
	if v, ok := s.cache.Get(sessionToken); ok {
		TLog.verbose("Fetch userId", v, "of sessionToken", sessionToken, "from Cache")
		return v.(string)
	}

	user, err := getUser(sessionToken)
	if err != nil {
		TLog.error("Can not fetch userId for sessionToken", sessionToken, ", error", err.Error())
		return ""
	}

	TLog.verbose("Fetch userId", user["objectId"], "of sessionToken", sessionToken, "from Parse")
	userID := user["objectId"].(string)
	s.cache.Add(sessionToken, userID)
	return userID
}
