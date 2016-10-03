package utils

import (
	"math"
	"regexp"

	"github.com/lfq7413/tomato/livequery/t"
)

// QueryHash ...
func QueryHash(query t.M) string {
	return ""
}

// MatchesQuery ...
func MatchesQuery(object, query t.M) bool {

	if className, ok := query["className"]; ok {
		if object["className"].(string) != className {
			return false
		}
		return MatchesQuery(object, query["where"].(map[string]interface{}))
	}

	for field, constraints := range query {
		if matchesKeyConstraints(object, field, constraints) == false {
			return false
		}
	}

	return true
}

func matchesKeyConstraints(object t.M, key string, constraints interface{}) bool {
	if constraints == nil {
		return false
	}
	if key == "$or" {
		if querys, ok := constraints.([]interface{}); ok {
			for _, query := range querys {
				if q, ok := query.(map[string]interface{}); ok {
					if MatchesQuery(object, q) {
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

	var constraint t.M
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
		if objectType == "Pointer" {
			o := object[key].(map[string]interface{})
			if o == nil {
				return false
			}
			class1 := constraint["className"]
			class2 := o["className"]
			id1 := constraint["objectId"]
			id2 := o["objectId"]
			return class1 == class2 && id1 == id2
		}

		if objects, ok := object[key].([]interface{}); ok {
			for _, object := range objects {
				if equalObject(object, constraint) {
					return true
				}
			}
			return false
		}
		return equalObject(object[key], constraint)
	}

	for condition, compareTo := range constraint {
		switch condition {
		case "$lt", "$lte", "$gt", "$gte":
			if compareNumber(object[key], compareTo, condition) == false {
				return false
			}
		case "$ne":
			if equalObject(object[key], compareTo) {
				return false
			}
		case "$in":
			if inSlice(compareTo, object[key]) == false {
				return false
			}
		case "$nin":
			if inSlice(compareTo, object[key]) {
				return false
			}
		case "$all":
			if compareToObjects, ok := compareTo.([]interface{}); ok {
				for _, compareToObject := range compareToObjects {
					if inSlice(object[key], compareToObject) == false {
						return false
					}
				}
			} else {
				return false
			}
		case "$exists":
			var propertyExists bool
			if _, ok := object[key]; ok {
				propertyExists = true
			} else {
				propertyExists = false
			}
			var existenceIsRequired bool
			if v, ok := constraint["$exists"].(bool); ok {
				existenceIsRequired = v
			} else {
				break
			}
			if (!propertyExists && existenceIsRequired) || (propertyExists && !existenceIsRequired) {
				return false
			}
		case "$regex":
			if compareRegexp(compareTo, object[key]) == false {
				return false
			}
		case "$nearSphere":
			if compareGeoPoint(compareTo, object[key], constraint["$maxDistance"]) == false {
				return false
			}
		case "$within":
			if compareBox(compareTo, object[key]) == false {
				return false
			}
		case "$options":
		case "$maxDistance":
		case "$select":
			return false
		case "$dontSelect":
			return false
		default:
			return false
		}
	}

	return true
}

func compareBox(compareTo, point interface{}) bool {
	southWest := map[string]float64{}
	northEast := map[string]float64{}
	geoPoint := map[string]float64{}

	box := []interface{}{}
	if within, ok := compareTo.(map[string]interface{}); ok {
		if b, ok := within["$box"].([]interface{}); ok {
			box = b
		} else {
			return false
		}
	} else {
		return false
	}

	if len(box) != 2 {
		return false
	}

	if p, ok := box[0].(map[string]interface{}); ok {
		southWest["longitude"] = p["longitude"].(float64)
		southWest["latitude"] = p["latitude"].(float64)
	} else {
		return false
	}

	if p, ok := box[1].(map[string]interface{}); ok {
		northEast["longitude"] = p["longitude"].(float64)
		northEast["latitude"] = p["latitude"].(float64)
	} else {
		return false
	}

	if southWest["latitude"] > northEast["latitude"] ||
		southWest["longitude"] > northEast["longitude"] {
		return false
	}

	if p, ok := point.(map[string]interface{}); ok {
		geoPoint["longitude"] = p["longitude"].(float64)
		geoPoint["latitude"] = p["latitude"].(float64)
	} else {
		return false
	}

	return geoPoint["latitude"] > southWest["latitude"] &&
		geoPoint["latitude"] < northEast["latitude"] &&
		geoPoint["longitude"] > southWest["longitude"] &&
		geoPoint["longitude"] < northEast["longitude"]
}

func compareGeoPoint(p1, p2, maxDistance interface{}) bool {
	if v1, ok := p1.(map[string]interface{}); ok {
		if v2, ok := p2.(map[string]interface{}); ok {
			if x1, ok := v1["longitude"].(float64); ok {
				if y1, ok := v1["latitude"].(float64); ok {
					if x2, ok := v2["longitude"].(float64); ok {
						if y2, ok := v2["latitude"].(float64); ok {
							d := distance(x1, y1, x2, y2)
							if maxDistance == nil {
								return true
							}
							if max, ok := maxDistance.(float64); ok {
								return d <= max
							}
						}
					}
				}
			}
		}
	}
	return false
}

// distance 计算两点之间的弧度值，两点经纬度分别为 (x1,y1), (x2,y2)
func distance(x1, y1, x2, y2 float64) float64 {
	x1 = x1 * math.Pi / 180
	y1 = y1 * math.Pi / 180
	x2 = x2 * math.Pi / 180
	y2 = y2 * math.Pi / 180
	// d=R*arcos[cos(Y1)*cos(Y2)*cos(X1-X2)+sin(Y1)*sin(Y2)]
	return math.Acos(math.Cos(y1)*math.Cos(y2)*math.Cos(x1-x2) + math.Sin(y1)*math.Sin(y2))
}

// compareRegexp 比较 object 是否符合正则表达式 exp
func compareRegexp(exp, object interface{}) bool {
	if c, ok := exp.(string); ok {
		if o, ok := object.(string); ok {
			matched, err := regexp.MatchString(c, o)
			if err != nil {
				return false
			}
			return matched
		}
		return false
	}
	return false
}

// inSlice 判断 slice 中是否存在指定的对象
// s 必须为 []interface{} 类型
// o 要查找的对象
func inSlice(s, o interface{}) bool {
	if objects, ok := s.([]interface{}); ok {
		for _, object := range objects {
			if equalObject(object, o) {
				return true
			}
		}
		return false
	}
	return false
}

// compareNumber i1, i2 均为数字类型时才比较大小
func compareNumber(i1, i2 interface{}, op string) bool {
	if v1, ok := i1.(float64); ok {
		if v2, ok := i2.(float64); ok {
			switch op {
			case "$lt":
				return v1 < v2
			case "$lte":
				return v1 <= v2
			case "$gt":
				return v1 > v2
			case "$gte":
				return v1 >= v2
			default:
				return false
			}
		}
		return false
	}
	return false
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
