package utils

import "testing"
import "reflect"
import tp "github.com/lfq7413/tomato/livequery/t"

func Test_QueryHash(t *testing.T) {
	data := []struct {
		query  tp.M
		expect string
	}{
		{
			query: tp.M{
				"className": "Player",
				"where":     map[string]interface{}{},
			},
			expect: "Player:|[]",
		},
		{
			query: tp.M{
				"className": "Player",
				"where": map[string]interface{}{
					"name": "joe",
				},
			},
			expect: "Player:name|[joe]",
		},
		{
			query: tp.M{
				"className": "Player",
				"where": map[string]interface{}{
					"name": "joe",
					"age":  12,
				},
			},
			expect: "Player:age,name|[12 joe]",
		},
		{
			query: tp.M{
				"className": "Player",
				"where": map[string]interface{}{"$or": []interface{}{
					map[string]interface{}{
						"name": "joe",
					},
					map[string]interface{}{
						"age": "12",
					},
				}},
			},
			expect: "Player:age,name|[joe 12]",
		},
		{
			query: tp.M{
				"className": "Player",
				"where": map[string]interface{}{"$or": []interface{}{
					map[string]interface{}{
						"name": "joe",
					},
					map[string]interface{}{
						"name": "joe",
						"age":  "12",
					},
				}},
			},
			expect: "Player:age,name|[joe 12 joe]",
		},
	}

	for _, d := range data {
		result := QueryHash(d.query)
		if reflect.DeepEqual(d.expect, result) == false {
			t.Error("expect:", d.expect, "result:", result)
		}
	}
}

func Test_flattenOrQueries(t *testing.T) {
	data := []struct {
		where  tp.M
		expect []map[string]interface{}
	}{
		{
			where:  tp.M{},
			expect: nil,
		},
		{
			where:  tp.M{"$or": 1024},
			expect: nil,
		},
		{
			where: tp.M{"$or": []interface{}{
				map[string]interface{}{
					"name": "joe",
				},
				1024,
			}},
			expect: []map[string]interface{}{
				map[string]interface{}{
					"name": "joe",
				},
			},
		},
		{
			where: tp.M{"$or": []interface{}{
				map[string]interface{}{
					"name": "joe",
				},
				map[string]interface{}{
					"age": "20",
				},
			}},
			expect: []map[string]interface{}{
				map[string]interface{}{
					"name": "joe",
				},
				map[string]interface{}{
					"age": "20",
				},
			},
		},
	}

	for _, d := range data {
		result := flattenOrQueries(d.where)
		if reflect.DeepEqual(d.expect, result) == false {
			t.Error("expect:", d.expect, "result:", result)
		}
	}
}

func Test_MatchesQuery(t *testing.T) {
	data := []struct {
		object tp.M
		query  tp.M
		expect bool
	}{
		{
			object: tp.M{
				"name":  "joe",
				"names": []interface{}{"joe", "tom"},
				"user": map[string]interface{}{
					"__type":    "Pointer",
					"className": "user",
					"objectId":  "1024",
				},
				"time": map[string]interface{}{
					"__type": "Date",
					"iso":    "2016-09-28T08:33:34.551Z",
				},
				"times": []interface{}{
					map[string]interface{}{
						"__type": "Date",
						"iso":    "2016-09-28T08:33:34.551Z",
					},
					map[string]interface{}{
						"__type": "Date",
						"iso":    "2017-09-28T08:33:34.551Z",
					},
				},
				"age1":  20,
				"age2":  20,
				"age3":  20,
				"age4":  20,
				"age5":  []interface{}{15, 20, 25},
				"age6":  20,
				"name2": "jack",
				"location": map[string]interface{}{
					"longitude": 0.0,
					"latitude":  0.0,
				},
				"location2": map[string]interface{}{
					"longitude": 10.0,
					"latitude":  10.0,
				},
			},
			query: tp.M{
				"name":  "joe",
				"names": "joe",
				"user": map[string]interface{}{
					"__type":    "Pointer",
					"className": "user",
					"objectId":  "1024",
				},
				"time": map[string]interface{}{
					"__type": "Date",
					"iso":    "2016-09-28T08:33:34.551Z",
				},
				"times": map[string]interface{}{
					"__type": "Date",
					"iso":    "2016-09-28T08:33:34.551Z",
				},
				"age1": map[string]interface{}{
					"$lt": 25,
				},
				"age2": map[string]interface{}{
					"$ne": 25,
				},
				"age3": map[string]interface{}{
					"$in": []interface{}{15, 20, 25},
				},
				"age4": map[string]interface{}{
					"$nin": []interface{}{15, 25},
				},
				"age5": map[string]interface{}{
					"$all": []interface{}{15, 25},
				},
				"age6": map[string]interface{}{
					"$exists": true,
				},
				"name2": map[string]interface{}{
					"$regex": "j*",
				},
				"location": map[string]interface{}{
					"$nearSphere": map[string]interface{}{
						"longitude": 90.0,
						"latitude":  0.0,
					},
					"$maxDistance": 2.0,
				},
				"location2": map[string]interface{}{
					"$within": map[string]interface{}{
						"$box": []interface{}{
							map[string]interface{}{
								"longitude": 0.0,
								"latitude":  0.0,
							},
							map[string]interface{}{
								"longitude": 20.0,
								"latitude":  20.0,
							},
						},
					},
				},
			},
			expect: true,
		},
		{
			object: tp.M{
				"age": 15,
			},
			query: tp.M{
				"$or": []interface{}{
					map[string]interface{}{
						"age": 20,
					},
					map[string]interface{}{
						"age": 15,
					},
				},
			},
			expect: true,
		},
	}

	for _, d := range data {
		result := MatchesQuery(d.object, d.query)
		if reflect.DeepEqual(d.expect, result) == false {
			t.Error("expect:", d.expect, "result:", result)
		}
	}
}

func Test_matchesKeyConstraints(t *testing.T) {
	data := []struct {
		object      tp.M
		key         string
		constraints interface{}
		expect      bool
	}{
		{
			object:      nil,
			key:         "name",
			constraints: nil,
			expect:      false,
		},
		{
			object:      tp.M{},
			key:         "name",
			constraints: nil,
			expect:      false,
		},
		{
			object: tp.M{
				"post": map[string]interface{}{
					"user": "joe",
				},
			},
			key:         "post.user",
			constraints: "joe",
			expect:      true,
		},
		{
			object: tp.M{
				"post": map[string]interface{}{
					"user": "joe",
				},
			},
			key:         "post.user",
			constraints: "jack",
			expect:      false,
		},
		{
			object: tp.M{
				"post": map[string]interface{}{
					"user": map[string]interface{}{
						"id": "1024",
					},
				},
			},
			key:         "post.user.id",
			constraints: "1024",
			expect:      true,
		},
		{
			object: tp.M{
				"age": 15,
			},
			key: "$or",
			constraints: []interface{}{
				map[string]interface{}{
					"age": 20,
				},
				map[string]interface{}{
					"age": 15,
				},
			},
			expect: true,
		},
		{
			object: tp.M{
				"age": 25,
			},
			key: "$or",
			constraints: []interface{}{
				map[string]interface{}{
					"age": 20,
				},
				map[string]interface{}{
					"age": 15,
				},
			},
			expect: false,
		},
		{
			object:      tp.M{},
			key:         "$relatedTo",
			constraints: 1024,
			expect:      false,
		},
		{
			object:      tp.M{},
			key:         "name",
			constraints: []interface{}{},
			expect:      false,
		},
		{
			object: tp.M{
				"name": "joe",
			},
			key:         "name",
			constraints: "joe",
			expect:      true,
		},
		{
			object: tp.M{
				"name": []interface{}{"joe", "tom"},
			},
			key:         "name",
			constraints: "joe",
			expect:      true,
		},
		{
			object: tp.M{
				"user": "joe",
			},
			key: "user",
			constraints: map[string]interface{}{
				"__type": "Pointer",
			},
			expect: false,
		},
		{
			object: tp.M{
				"user": map[string]interface{}{
					"__type":    "Pointer",
					"className": "user",
					"objectId":  "1024",
				},
			},
			key: "user",
			constraints: map[string]interface{}{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1024",
			},
			expect: true,
		},
		{
			object: tp.M{
				"user": map[string]interface{}{
					"__type":    "Pointer",
					"className": "user",
					"objectId":  "1024",
				},
			},
			key: "user",
			constraints: map[string]interface{}{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "2048",
			},
			expect: false,
		},
		{
			object: tp.M{
				"user": []interface{}{
					map[string]interface{}{
						"__type":    "Pointer",
						"className": "user",
						"objectId":  "1024",
					},
					map[string]interface{}{
						"__type":    "Pointer",
						"className": "user",
						"objectId":  "2048",
					},
				},
			},
			key: "user",
			constraints: map[string]interface{}{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1024",
			},
			expect: true,
		},
		{
			object: tp.M{
				"time": map[string]interface{}{
					"__type": "Date",
					"iso":    "2016-09-28T08:33:34.551Z",
				},
			},
			key: "time",
			constraints: map[string]interface{}{
				"__type": "Date",
				"iso":    "2016-09-28T08:33:34.551Z",
			},
			expect: true,
		},
		{
			object: tp.M{
				"time": []interface{}{
					map[string]interface{}{
						"__type": "Date",
						"iso":    "2016-09-28T08:33:34.551Z",
					},
					map[string]interface{}{
						"__type": "Date",
						"iso":    "2017-09-28T08:33:34.551Z",
					},
				},
			},
			key: "time",
			constraints: map[string]interface{}{
				"__type": "Date",
				"iso":    "2016-09-28T08:33:34.551Z",
			},
			expect: true,
		},
		{
			object: tp.M{
				"age": 20,
			},
			key: "age",
			constraints: map[string]interface{}{
				"$lt": 25,
			},
			expect: true,
		},
		{
			object: tp.M{
				"age": 20,
			},
			key: "age",
			constraints: map[string]interface{}{
				"$lte": 20,
			},
			expect: true,
		},
		{
			object: tp.M{
				"age": 20,
			},
			key: "age",
			constraints: map[string]interface{}{
				"$gt": 15,
			},
			expect: true,
		},
		{
			object: tp.M{
				"age": 20,
			},
			key: "age",
			constraints: map[string]interface{}{
				"$gte": 20,
			},
			expect: true,
		},
		{
			object: tp.M{
				"age": 20,
			},
			key: "age",
			constraints: map[string]interface{}{
				"$ne": 15,
			},
			expect: true,
		},
		{
			object: tp.M{
				"age": 20,
			},
			key: "age",
			constraints: map[string]interface{}{
				"$in": []interface{}{15, 20, 25},
			},
			expect: true,
		},
		{
			object: tp.M{
				"age": 20,
			},
			key: "age",
			constraints: map[string]interface{}{
				"$nin": []interface{}{15, 25},
			},
			expect: true,
		},
		{
			object: tp.M{
				"age": []interface{}{15, 20, 25},
			},
			key: "age",
			constraints: map[string]interface{}{
				"$all": []interface{}{15, 25},
			},
			expect: true,
		},
		{
			object: tp.M{
				"age": 20,
			},
			key: "age",
			constraints: map[string]interface{}{
				"$exists": true,
			},
			expect: true,
		},
		{
			object: tp.M{
				"age": 20,
			},
			key: "name",
			constraints: map[string]interface{}{
				"$exists": false,
			},
			expect: true,
		},
		{
			object: tp.M{
				"name": "joe",
			},
			key: "name",
			constraints: map[string]interface{}{
				"$regex": "j*",
			},
			expect: true,
		},
		{
			object: tp.M{
				"location": map[string]interface{}{
					"longitude": 0.0,
					"latitude":  0.0,
				},
			},
			key: "location",
			constraints: map[string]interface{}{
				"$nearSphere": map[string]interface{}{
					"longitude": 90.0,
					"latitude":  0.0,
				},
				"$maxDistance": 2.0,
			},
			expect: true,
		},
		{
			object: tp.M{
				"location": map[string]interface{}{
					"longitude": 10.0,
					"latitude":  10.0,
				},
			},
			key: "location",
			constraints: map[string]interface{}{
				"$within": map[string]interface{}{
					"$box": []interface{}{
						map[string]interface{}{
							"longitude": 0.0,
							"latitude":  0.0,
						},
						map[string]interface{}{
							"longitude": 20.0,
							"latitude":  20.0,
						},
					},
				},
			},
			expect: true,
		},
		{
			object: tp.M{
				"name": "joe",
			},
			key: "name",
			constraints: map[string]interface{}{
				"$select": 1024,
			},
			expect: false,
		},
		{
			object: tp.M{
				"name": "joe",
			},
			key: "name",
			constraints: map[string]interface{}{
				"$dontSelect": 1024,
			},
			expect: false,
		},
	}

	for _, d := range data {
		result := matchesKeyConstraints(d.object, d.key, d.constraints)
		if reflect.DeepEqual(d.expect, result) == false {
			t.Error("expect:", d.expect, "result:", result)
		}
	}
}

func Test_compareBox(t *testing.T) {
	data := []struct {
		compareTo interface{}
		point     interface{}
		expect    bool
	}{
		{
			compareTo: "hello",
			point: map[string]interface{}{
				"longitude": 10.0,
				"latitude":  10.0,
			},
			expect: false,
		},
		{
			compareTo: map[string]interface{}{
				"$box": "hello",
			},
			point: map[string]interface{}{
				"longitude": 10.0,
				"latitude":  10.0,
			},
			expect: false,
		},
		{
			compareTo: map[string]interface{}{
				"$box": []interface{}{1},
			},
			point: map[string]interface{}{
				"longitude": 10.0,
				"latitude":  10.0,
			},
			expect: false,
		},
		{
			compareTo: map[string]interface{}{
				"$box": []interface{}{1, 2},
			},
			point: map[string]interface{}{
				"longitude": 10.0,
				"latitude":  10.0,
			},
			expect: false,
		},
		{
			compareTo: map[string]interface{}{
				"$box": []interface{}{
					map[string]interface{}{
						"longitude": 10.0,
						"latitude":  0.0,
					},
					map[string]interface{}{
						"longitude": 0.0,
						"latitude":  0.0,
					},
				},
			},
			point: map[string]interface{}{
				"longitude": 10.0,
				"latitude":  10.0,
			},
			expect: false,
		},
		{
			compareTo: map[string]interface{}{
				"$box": []interface{}{
					map[string]interface{}{
						"longitude": 0.0,
						"latitude":  10.0,
					},
					map[string]interface{}{
						"longitude": 10.0,
						"latitude":  0.0,
					},
				},
			},
			point: map[string]interface{}{
				"longitude": 10.0,
				"latitude":  10.0,
			},
			expect: false,
		},
		{
			compareTo: map[string]interface{}{
				"$box": []interface{}{
					map[string]interface{}{
						"longitude": 0.0,
						"latitude":  0.0,
					},
					map[string]interface{}{
						"longitude": 10.0,
						"latitude":  10.0,
					},
				},
			},
			point:  "hello",
			expect: false,
		},
		{
			compareTo: map[string]interface{}{
				"$box": []interface{}{
					map[string]interface{}{
						"longitude": 0.0,
						"latitude":  0.0,
					},
					map[string]interface{}{
						"longitude": 10.0,
						"latitude":  10.0,
					},
				},
			},
			point: map[string]interface{}{
				"longitude": 20.0,
				"latitude":  20.0,
			},
			expect: false,
		},
		{
			compareTo: map[string]interface{}{
				"$box": []interface{}{
					map[string]interface{}{
						"longitude": 0.0,
						"latitude":  0.0,
					},
					map[string]interface{}{
						"longitude": 10.0,
						"latitude":  10.0,
					},
				},
			},
			point: map[string]interface{}{
				"longitude": 5.0,
				"latitude":  5.0,
			},
			expect: true,
		},
	}

	for _, d := range data {
		result := compareBox(d.compareTo, d.point)
		if reflect.DeepEqual(d.expect, result) == false {
			t.Error("expect:", d.expect, "result:", result)
		}
	}
}

func Test_compareGeoPoint(t *testing.T) {
	data := []struct {
		p1          interface{}
		p2          interface{}
		maxDistance interface{}
		expect      bool
	}{
		{
			p1: 1024,
			p2: map[string]interface{}{
				"longitude": 10.0,
				"latitude":  10.0,
			},
			maxDistance: nil,
			expect:      false,
		},
		{
			p1: map[string]interface{}{
				"longitude": 10.0,
				"latitude":  10.0,
			},
			p2:          1024,
			maxDistance: nil,
			expect:      false,
		},
		{
			p1: map[string]interface{}{
				"longitude": 10.0,
				"latitude":  10.0,
			},
			p2: map[string]interface{}{
				"longitude": 10.0,
				"latitude":  20.0,
			},
			maxDistance: nil,
			expect:      true,
		},
		{
			p1: map[string]interface{}{
				"longitude": 10.0,
				"latitude":  10.0,
			},
			p2: map[string]interface{}{
				"longitude": 10.0,
				"latitude":  20.0,
			},
			maxDistance: "hello",
			expect:      true,
		},
		{
			p1: map[string]interface{}{
				"longitude": 0.0,
				"latitude":  0.0,
			},
			p2: map[string]interface{}{
				"longitude": 90.0,
				"latitude":  0.0,
			},
			maxDistance: 2.0,
			expect:      true,
		},
		{
			p1: map[string]interface{}{
				"longitude": 0.0,
				"latitude":  0.0,
			},
			p2: map[string]interface{}{
				"longitude": 90.0,
				"latitude":  0.0,
			},
			maxDistance: 1.0,
			expect:      false,
		},
	}

	for _, d := range data {
		result := compareGeoPoint(d.p1, d.p2, d.maxDistance)
		if reflect.DeepEqual(d.expect, result) == false {
			t.Error("expect:", d.expect, "result:", result)
		}
	}
}

func Test_distance(t *testing.T) {
	data := []struct {
		x1     float64
		y1     float64
		x2     float64
		y2     float64
		expect float64
	}{
		{x1: 0, y1: 0, x2: 0, y2: 0, expect: 0},
		{x1: 0, y1: 0, x2: 180, y2: 0, expect: 3.141592653589793},
		{x1: 0, y1: 0, x2: 0, y2: 90, expect: 1.5707963267948966},
	}

	for _, d := range data {
		result := distance(d.x1, d.y1, d.x2, d.y2)
		if reflect.DeepEqual(d.expect, result) == false {
			t.Error("expect:", d.expect, "result:", result)
		}
	}
}

func Test_compareRegexp(t *testing.T) {
	data := []struct {
		exp    interface{}
		object interface{}
		expect bool
	}{
		{exp: "hello", object: 1024, expect: false},
		{exp: 1024, object: "hello", expect: false},
		{exp: "hello", object: "hello", expect: true},
		{exp: "hell*", object: "hello", expect: true},
		{exp: "hell*", object: "hi", expect: false},
	}

	for _, d := range data {
		result := compareRegexp(d.exp, d.object)
		if reflect.DeepEqual(d.expect, result) == false {
			t.Error("expect:", d.expect, "result:", result)
		}
	}
}

func Test_inSlice(t *testing.T) {
	data := []struct {
		s      interface{}
		o      interface{}
		expect bool
	}{
		{s: "hello", o: 1, expect: false},
		{s: []interface{}{1, 2, 3}, o: 4, expect: false},
		{s: []interface{}{1, 2, 3}, o: 3, expect: true},
	}

	for _, d := range data {
		result := inSlice(d.s, d.o)
		if reflect.DeepEqual(d.expect, result) == false {
			t.Error("expect:", d.expect, "result:", result)
		}
	}
}

func Test_compareNumber(t *testing.T) {
	data := []struct {
		i1     interface{}
		i2     interface{}
		op     string
		expect bool
	}{
		{i1: 10.0, i2: 20.0, op: "$lt", expect: true},
		{i1: 10, i2: 20.0, op: "$lt", expect: true},
		{i1: 10.0, i2: 20, op: "$lt", expect: true},
		{i1: "hi", i2: 20, op: "$lt", expect: false},
		{i1: 10, i2: "hi", op: "$lt", expect: false},
		{i1: 10, i2: 20, op: "$lt", expect: true},
		{i1: 20, i2: 20, op: "$lte", expect: true},
		{i1: 20, i2: 10, op: "$gt", expect: true},
		{i1: 20, i2: 20, op: "$gte", expect: true},
		{i1: 30, i2: 20, op: "$lt", expect: false},
		{i1: 30, i2: 20, op: "$lte", expect: false},
		{i1: 20, i2: 30, op: "$gt", expect: false},
		{i1: 20, i2: 30, op: "$gte", expect: false},
		{i1: 30, i2: 20, op: "$other", expect: false},
	}

	for _, d := range data {
		result := compareNumber(d.i1, d.i2, d.op)
		if reflect.DeepEqual(d.expect, result) == false {
			t.Error("expect:", d.expect, "result:", result)
		}
	}
}

func Test_equalObject(t *testing.T) {
	data := []struct {
		i1     interface{}
		i2     interface{}
		expect bool
	}{
		{i1: "hello", i2: "hello", expect: true},
		{i1: "hello", i2: 11.0, expect: false},
		{i1: "hello", i2: 1024, expect: false},
		{i1: 10.0, i2: 10.0, expect: true},
		{i1: 10.0, i2: 11.0, expect: false},
		{i1: 10.0, i2: "hi", expect: false},
		{i1: 10, i2: 10, expect: true},
		{i1: 10, i2: 11, expect: false},
		{i1: 10, i2: "hi", expect: false},
		{i1: true, i2: true, expect: true},
		{i1: true, i2: false, expect: false},
		{i1: true, i2: "hi", expect: false},
		{
			i1:     []interface{}{1, 2, 3},
			i2:     []interface{}{1, 2},
			expect: false,
		},
		{
			i1:     []interface{}{1, 2, 3},
			i2:     []interface{}{1, 2, 4},
			expect: false,
		},
		{
			i1:     []interface{}{1, 2, 3},
			i2:     []interface{}{1, 2, 3},
			expect: true,
		},
		{
			i1:     []interface{}{1, 2, 3},
			i2:     "hi",
			expect: false,
		},
		{
			i1: map[string]interface{}{
				"name": "joe",
				"age":  12,
			},
			i2: map[string]interface{}{
				"name": "joe",
			},
			expect: false,
		},
		{
			i1: map[string]interface{}{
				"name": "joe",
				"age":  12,
			},
			i2: map[string]interface{}{
				"name": "joe",
				"age":  20,
			},
			expect: false,
		},
		{
			i1: map[string]interface{}{
				"name": "joe",
				"age":  12,
			},
			i2: map[string]interface{}{
				"name": "joe",
				"age":  12,
			},
			expect: true,
		},
		{
			i1: map[string]interface{}{
				"name": "joe",
				"age":  12,
			},
			i2:     "hi",
			expect: false,
		},
	}
	for _, d := range data {
		result := equalObject(d.i1, d.i2)
		if reflect.DeepEqual(d.expect, result) == false {
			t.Error("expect:", d.expect, "result:", result)
		}
	}
}
