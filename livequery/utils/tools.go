package utils

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"

	"github.com/lfq7413/tomato/livequery/t"
)

// QueryHash 计算 query 的 hash
func QueryHash(query t.M) string {
	where := query["where"].(map[string]interface{})
	if where == nil {
		where = map[string]interface{}{}
	}

	list := flattenOrQueries(where)
	columns := []string{}
	values := []interface{}{}
	if len(list) > 0 {
		uniqueColumns := map[string]bool{}
		for _, subQuery := range list {
			keys := []string{}
			for k := range subQuery {
				keys = append(keys, k)
			}
			sort.Sort(sort.StringSlice(keys))
			for _, k := range keys {
				values = append(values, subQuery[k])
				uniqueColumns[k] = true
			}
		}
		for k := range uniqueColumns {
			columns = append(columns, k)
		}
		sort.Sort(sort.StringSlice(columns))
	} else {
		for k := range where {
			columns = append(columns, k)
		}
		sort.Sort(sort.StringSlice(columns))
		for _, k := range columns {
			values = append(values, where[k])
		}
	}

	sections := []string{strings.Join(columns, ","), fmt.Sprint(values)}

	return query["className"].(string) + ":" + strings.Join(sections, "|")
}

// flattenOrQueries 处理 where 中的 $or 字段
func flattenOrQueries(where t.M) []map[string]interface{} {
	ors := where["$or"]
	if ors == nil {
		return nil
	}
	if v, ok := ors.([]interface{}); ok {
		accum := []map[string]interface{}{}
		for _, sub := range v {
			if query, ok := sub.(map[string]interface{}); ok {
				accum = append(accum, query)
			}
		}
		return accum
	}
	return nil
}

// MatchesQuery 检测对象是否符合订阅条件
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

// matchesKeyConstraints 检测对象中的字段是否符合指定的条件
func matchesKeyConstraints(object t.M, key string, constraints interface{}) bool {
	if constraints == nil {
		return false
	}
	// 处理 $or ，有一处符合即可
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

	// 处理 constraints 不为 array map 时的情况
	var constraint t.M
	if v, ok := constraints.(map[string]interface{}); ok {
		constraint = v
	} else {
		// 当 object[key] 为数组时，只要有一个符合即可
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

	// 处理 constraint 为 __type 类型时的情况
	var objectType string
	if v, ok := constraint["__type"].(string); ok {
		objectType = v
	}
	if objectType != "" {
		if objectType == "Pointer" {
			var o map[string]interface{}
			if v, ok := object[key].(map[string]interface{}); ok {
				o = v
			}
			if o == nil {
				return false
			}
			class1 := constraint["className"]
			class2 := o["className"]
			id1 := constraint["objectId"]
			id2 := o["objectId"]
			return equalObject(class1, class2) && equalObject(id1, id2)
		}

		// 如果 object[key] 为数组，只要一个符合条件即可
		if objs, ok := object[key].([]interface{}); ok {
			for _, obj := range objs {
				if equalObject(obj, constraints) {
					return true
				}
			}
			return false
		}
		return equalObject(object[key], constraints)
	}

	// 处理 constraint 包含限制条件时的情况
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

// compareBox 校验一点是否在指定区域内
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
		if v, ok := p["longitude"].(float64); ok {
			southWest["longitude"] = v
		}
		if v, ok := p["latitude"].(float64); ok {
			southWest["latitude"] = v
		}
	} else {
		return false
	}

	if p, ok := box[1].(map[string]interface{}); ok {
		if v, ok := p["longitude"].(float64); ok {
			northEast["longitude"] = v
		}
		if v, ok := p["latitude"].(float64); ok {
			northEast["latitude"] = v
		}
	} else {
		return false
	}

	if southWest["latitude"] > northEast["latitude"] ||
		southWest["longitude"] > northEast["longitude"] {
		return false
	}

	if p, ok := point.(map[string]interface{}); ok {
		if v, ok := p["longitude"].(float64); ok {
			geoPoint["longitude"] = v
		}
		if v, ok := p["latitude"].(float64); ok {
			geoPoint["latitude"] = v
		}
	} else {
		return false
	}

	return geoPoint["latitude"] > southWest["latitude"] &&
		geoPoint["latitude"] < northEast["latitude"] &&
		geoPoint["longitude"] > southWest["longitude"] &&
		geoPoint["longitude"] < northEast["longitude"]
}

// compareGeoPoint 比较两点是否相邻
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
							return true
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
	var v1 float64
	var v2 float64

	switch i1.(type) {
	case float64:
		v1 = i1.(float64)
	case int:
		v1 = float64(i1.(int))
	default:
		return false
	}

	switch i2.(type) {
	case float64:
		v2 = i2.(float64)
	case int:
		v2 = float64(i2.(int))
	default:
		return false
	}

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

	if v1, ok := i1.(int); ok {
		if v2, ok := i2.(int); ok {
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
