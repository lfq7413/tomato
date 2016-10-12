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
	// equalObject
	// compareNumber
	// inSlice
	// compareRegexp
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
	// TODO
}

func Test_compareRegexp(t *testing.T) {
	// TODO
}

func Test_inSlice(t *testing.T) {
	// TODO
	// equalObject
}

func Test_compareNumber(t *testing.T) {
	// TODO
}

func Test_equalObject(t *testing.T) {
	// TODO
}
