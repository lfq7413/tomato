package livequery

type sessionTokenCache struct {
	cache map[string]string
}

func (s *sessionTokenCache) getUserID(sessionToken string) string {
	// TODO
	return ""
}
