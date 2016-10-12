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
	// TODO
	// matchesKeyConstraints
}

func Test_matchesKeyConstraints(t *testing.T) {
	// TODO
	// compareGeoPoint
	// compareBox
}

func Test_compareBox(t *testing.T) {
	// TODO
}

func Test_compareGeoPoint(t *testing.T) {
	// TODO
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
