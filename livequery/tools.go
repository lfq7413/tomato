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

	// 不支持 relatedTo
	if key == "$relatedTo" {
		return false
	}

	// 只支持 key == "$or" 时，constraints 为数组的情况
	if _, ok := constraints.([]interface{}); ok {
		return false
	}

	var constraint M
	if v, ok := constraints.(map[string]interface{}); ok {
		constraint = v
	} else {
		if objects, ok := object[key].([]interface{}); ok {
			for _, o := range objects {
				if equalObject(o, constraints) {
					return true
				}
			}
			return false
		}
		return equalObject(object[key], constraints)
	}

	objectType := constraint["__type"].(string)
	if objectType != "" {

	}

	return true
}

// equalObject 仅比较基础类型：string float64 bool slice map
func equalObject(i1, i2 interface{}) bool {
	if v1, ok := i1.(string); ok {
		if v2, ok := i2.(string); ok {
			return v1 == v2
		}
		return false
	}

	if v1, ok := i1.(float64); ok {
		if v2, ok := i2.(float64); ok {
			return v1 == v2
		}
		return false
	}

	if v1, ok := i1.(bool); ok {
		if v2, ok := i2.(bool); ok {
			return v1 == v2
		}
		return false
	}

	if v1, ok := i1.([]interface{}); ok {
		if v2, ok := i2.([]interface{}); ok {
			if len(v1) != len(v2) {
				return false
			}
			for i := 0; i < len(v1); i++ {
				if equalObject(v1[i], v2[i]) == false {
					return false
				}
			}
			return true
		}
		return false
	}

	if v1, ok := i1.(map[string]interface{}); ok {
		if v2, ok := i2.(map[string]interface{}); ok {
			if len(v1) != len(v2) {
				return false
			}
			for k := range v1 {
				if equalObject(v1[k], v2[k]) == false {
					return false
				}
			}
			return true
		}
		return false
	}

	return false
}
