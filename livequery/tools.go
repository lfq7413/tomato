package livequery

func queryHash(query M) string {
	return ""
}

func matchesQuery(object, query M) bool {

	if className, ok := query["className"]; ok {
		if object["className"].(string) != className {
			return false
		}
		return matchesQuery(object, query["where"].(map[string]interface{}))
	}

	for field, constraints := range query {
		if matchesKeyConstraints(object, field, constraints) == false {
			return false
		}
	}

	return true
}

func matchesKeyConstraints(object M, key string, constraints interface{}) bool {
	if key == "$or" {
		if querys, ok := constraints.([]interface{}); ok {
			for _, query := range querys {
				if q, ok := query.(map[string]interface{}); ok {
					if matchesQuery(object, q) {
						return true
					}
				}
			}
			return false
		}
		return false
	}

	if key == "$relatedTo" {
		return false
	}

	return true
}
