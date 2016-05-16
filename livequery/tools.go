package livequery

func queryHash(query M) string {
	return ""
}

func matchesQuery(object, query M) bool {
	if object["className"].(string) != query["className"].(string) {
		return false
	}

	where := query["where"].(map[string]interface{})

	for field, constraints := range where {
		if matchesKeyConstraints(object, field, constraints) == false {
			return false
		}
	}

	return true
}

func matchesKeyConstraints(object M, key string, constraints interface{}) bool {
	return false
}
