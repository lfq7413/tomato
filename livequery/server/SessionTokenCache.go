package server

import (
	"github.com/lfq7413/tomato/dependencies/lru"
	"github.com/lfq7413/tomato/livequery/utils"
)

// SessionTokenCache ...
type SessionTokenCache struct {
	cache *lru.Cache
}

// NewSessionTokenCache ...
func NewSessionTokenCache() *SessionTokenCache {
	return &SessionTokenCache{
		cache: lru.New(10000),
	}
}

// GetUserID ...
func (s *SessionTokenCache) GetUserID(sessionToken string) string {
	if v, ok := s.cache.Get(sessionToken); ok {
		utils.TLog.Verbose("Fetch userId", v, "of sessionToken", sessionToken, "from Cache")
		return v.(string)
	}

	user, err := getUser(sessionToken)
	if err != nil {
		utils.TLog.Error("Can not fetch userId for sessionToken", sessionToken, ", error", err.Error())
		return ""
	}

	utils.TLog.Verbose("Fetch userId", user["objectId"], "of sessionToken", sessionToken, "from Parse")
	userID := user["objectId"].(string)
	s.cache.Add(sessionToken, userID)
	return userID
}
