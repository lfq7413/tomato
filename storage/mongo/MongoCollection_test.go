package mongo

import (
	"errors"
	"reflect"
	"testing"

	"github.com/lfq7413/tomato/types"

	"gopkg.in/mgo.v2"
)

func Test_find(t *testing.T) {
	db := openDB()
	defer db.Session.Close()
	mc := &MongoCollection{collection: db.C("obj")}
	var docs interface{}
	var query interface{}
	var options types.M
	var result []types.M
	var err error
	var expect interface{}
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "location": types.S{30, 30}}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "location": types.S{15, 15}}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "tom", "location": types.S{20, 20}}
	mc.insertOne(docs)
	docs = types.M{"_id": "004", "name": "ann", "location": types.S{10, 10}}
	mc.insertOne(docs)
	query = types.M{
		"location": types.M{
			"$nearSphere": types.M{
				"$geometry": types.M{
					"type":        "Point",
					"coordinates": types.S{10, 10},
				},
				"$maxDistance": 1000000,
			},
		},
	}
	options = types.M{}
	for i := 0; i < 100; i++ {
		result, err = mc.find(query, options)
		if err != nil {
			t.Error(err)
		}
	}
	result, err = mc.find(query, options)
	expect = []types.M{
		types.M{"_id": "004", "name": "ann", "location": []interface{}{10, 10}},
		types.M{"_id": "002", "name": "jack", "location": []interface{}{15, 15}},
	}
	if err != nil || result == nil || len(result) != 2 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result, err)
	}
	mc.drop()
	// TODO
}

func Test_rawFind(t *testing.T) {
	db := openDB()
	defer db.Session.Close()
	mc := &MongoCollection{collection: db.C("obj")}
	var docs interface{}
	var query interface{}
	var options types.M
	var result []types.M
	var err error
	var expect interface{}
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "tom", "age": 31}
	query = types.M{"name": "jone"}
	options = types.M{}
	result, err = mc.rawFind(query, options)
	expect = []types.M{}
	if err != nil || (result != nil && len(result) != 0) {
		t.Error("expect:", expect, "get result:", result)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "tom", "age": 31}
	mc.insertOne(docs)
	query = types.M{"name": "joe"}
	options = types.M{}
	result, err = mc.rawFind(query, options)
	expect = []types.M{
		types.M{"_id": "001", "name": "joe", "age": 25},
	}
	if err != nil || result == nil || len(result) != 1 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "tom", "age": 31}
	mc.insertOne(docs)
	query = types.M{}
	options = types.M{
		"sort": []string{"age"},
	}
	result, err = mc.rawFind(query, options)
	expect = []types.M{
		types.M{"_id": "001", "name": "joe", "age": 25},
		types.M{"_id": "002", "name": "jack", "age": 30},
		types.M{"_id": "003", "name": "tom", "age": 31},
	}
	if err != nil || result == nil || len(result) != 3 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "tom", "age": 31}
	mc.insertOne(docs)
	query = types.M{}
	options = types.M{
		"sort": []string{"-age"},
	}
	result, err = mc.rawFind(query, options)
	expect = []types.M{
		types.M{"_id": "003", "name": "tom", "age": 31},
		types.M{"_id": "002", "name": "jack", "age": 30},
		types.M{"_id": "001", "name": "joe", "age": 25},
	}
	if err != nil || result == nil || len(result) != 3 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "tom", "age": 31}
	mc.insertOne(docs)
	query = types.M{}
	options = types.M{
		"sort": []string{"name"},
	}
	result, err = mc.rawFind(query, options)
	expect = []types.M{
		types.M{"_id": "002", "name": "jack", "age": 30},
		types.M{"_id": "001", "name": "joe", "age": 25},
		types.M{"_id": "003", "name": "tom", "age": 31},
	}
	if err != nil || result == nil || len(result) != 3 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "tom", "age": 31}
	mc.insertOne(docs)
	query = types.M{}
	options = types.M{
		"sort": []string{"-name"},
	}
	result, err = mc.rawFind(query, options)
	expect = []types.M{
		types.M{"_id": "003", "name": "tom", "age": 31},
		types.M{"_id": "001", "name": "joe", "age": 25},
		types.M{"_id": "002", "name": "jack", "age": 30},
	}
	if err != nil || result == nil || len(result) != 3 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "tom", "age": 31}
	mc.insertOne(docs)
	query = types.M{}
	options = types.M{
		"skip": 1,
	}
	result, err = mc.rawFind(query, options)
	expect = []types.M{
		types.M{"_id": "002", "name": "jack", "age": 30},
		types.M{"_id": "003", "name": "tom", "age": 31},
	}
	if err != nil || result == nil || len(result) != 2 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "tom", "age": 31}
	mc.insertOne(docs)
	query = types.M{}
	options = types.M{
		"limit": 1,
	}
	result, err = mc.rawFind(query, options)
	expect = []types.M{
		types.M{"_id": "001", "name": "joe", "age": 25},
	}
	if err != nil || result == nil || len(result) != 1 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "tom", "age": 31}
	mc.insertOne(docs)
	docs = types.M{"_id": "004", "name": "ann", "age": 25}
	mc.insertOne(docs)
	query = types.M{
		"$or": types.S{types.M{"name": "joe"}, types.M{"age": 25}},
	}
	options = types.M{}
	result, err = mc.rawFind(query, options)
	expect = []types.M{
		types.M{"_id": "001", "name": "joe", "age": 25},
		types.M{"_id": "004", "name": "ann", "age": 25},
	}
	if err != nil || result == nil || len(result) != 2 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "tom", "age": 31}
	mc.insertOne(docs)
	docs = types.M{"_id": "004", "name": "ann", "age": 25}
	mc.insertOne(docs)
	query = types.M{
		"$and": types.S{types.M{"name": "joe"}, types.M{"age": 25}},
	}
	options = types.M{}
	result, err = mc.rawFind(query, options)
	expect = []types.M{
		types.M{"_id": "001", "name": "joe", "age": 25},
	}
	if err != nil || result == nil || len(result) != 1 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "skill": types.S{"one", "three", "five"}}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "skill": types.S{"two", "three", "six"}}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "tom", "skill": types.S{"one", "four", "six"}}
	mc.insertOne(docs)
	docs = types.M{"_id": "004", "name": "ann", "skill": types.S{"two", "four", "five"}}
	mc.insertOne(docs)
	query = types.M{
		"skill": types.M{"$all": types.S{"one", "four"}},
	}
	options = types.M{}
	result, err = mc.rawFind(query, options)
	expect = []types.M{
		types.M{"_id": "003", "name": "tom", "skill": []interface{}{"one", "four", "six"}},
	}
	if err != nil || result == nil || len(result) != 1 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "tom", "age": 31}
	mc.insertOne(docs)
	docs = types.M{"_id": "004", "name": "ann", "age": 25}
	mc.insertOne(docs)
	query = types.M{
		"age": types.M{"$lt": 30},
	}
	options = types.M{}
	result, err = mc.rawFind(query, options)
	expect = []types.M{
		types.M{"_id": "001", "name": "joe", "age": 25},
		types.M{"_id": "004", "name": "ann", "age": 25},
	}
	if err != nil || result == nil || len(result) != 2 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "tom", "age": 31, "skill": "one"}
	mc.insertOne(docs)
	docs = types.M{"_id": "004", "name": "ann", "age": 25}
	mc.insertOne(docs)
	query = types.M{
		"skill": types.M{"$exists": true},
	}
	options = types.M{}
	result, err = mc.rawFind(query, options)
	expect = []types.M{
		types.M{"_id": "003", "name": "tom", "age": 31, "skill": "one"},
	}
	if err != nil || result == nil || len(result) != 1 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "tom", "age": 31}
	mc.insertOne(docs)
	docs = types.M{"_id": "004", "name": "ann", "age": 25}
	mc.insertOne(docs)
	query = types.M{
		"age": types.M{"$in": types.S{30, 31}},
	}
	options = types.M{}
	result, err = mc.rawFind(query, options)
	expect = []types.M{
		types.M{"_id": "002", "name": "jack", "age": 30},
		types.M{"_id": "003", "name": "tom", "age": 31},
	}
	if err != nil || result == nil || len(result) != 2 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "tom", "age": 31}
	mc.insertOne(docs)
	docs = types.M{"_id": "004", "name": "ann", "age": 25}
	mc.insertOne(docs)
	query = types.M{
		"age": types.M{"$nin": types.S{30, 31}},
	}
	options = types.M{}
	result, err = mc.rawFind(query, options)
	expect = []types.M{
		types.M{"_id": "001", "name": "joe", "age": 25},
		types.M{"_id": "004", "name": "ann", "age": 25},
	}
	if err != nil || result == nil || len(result) != 2 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "tom", "age": 31}
	mc.insertOne(docs)
	docs = types.M{"_id": "004", "name": "ann", "age": 25}
	mc.insertOne(docs)
	query = types.M{
		"name": types.M{"$regex": `^j`},
	}
	options = types.M{}
	result, err = mc.rawFind(query, options)
	expect = []types.M{
		types.M{"_id": "001", "name": "joe", "age": 25},
		types.M{"_id": "002", "name": "jack", "age": 30},
	}
	if err != nil || result == nil || len(result) != 2 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "location": types.S{30, 30}}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "location": types.S{15, 15}}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "tom", "location": types.S{20, 20}}
	mc.insertOne(docs)
	docs = types.M{"_id": "004", "name": "ann", "location": types.S{10, 10}}
	mc.insertOne(docs)
	index := mgo.Index{
		Key:  []string{"$2dsphere:location"},
		Bits: 26,
	}
	mc.collection.EnsureIndex(index)
	query = types.M{
		"location": types.M{
			"$nearSphere": types.M{
				"$geometry": types.M{
					"type":        "Point",
					"coordinates": types.S{10, 10},
				},
				"$maxDistance": 1000000,
			},
		},
	}
	options = types.M{}
	result, err = mc.rawFind(query, options)
	expect = []types.M{
		types.M{"_id": "004", "name": "ann", "location": []interface{}{10, 10}},
		types.M{"_id": "002", "name": "jack", "location": []interface{}{15, 15}},
	}
	if err != nil || result == nil || len(result) != 2 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result, err)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "location": types.S{30, 30}}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "location": types.S{15, 15}}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "tom", "location": types.S{20, 20}}
	mc.insertOne(docs)
	docs = types.M{"_id": "004", "name": "ann", "location": types.S{10, 10}}
	mc.insertOne(docs)
	query = types.M{
		"location": types.M{
			"$geoWithin": types.M{
				"$box": types.S{
					types.S{5, 5},
					types.S{12, 12},
				},
			},
		},
	}
	options = types.M{}
	result, err = mc.rawFind(query, options)
	expect = []types.M{
		types.M{"_id": "004", "name": "ann", "location": []interface{}{10, 10}},
	}
	if err != nil || result == nil || len(result) != 1 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result, err)
	}
	mc.drop()
}

func Test_count(t *testing.T) {
	db := openDB()
	defer db.Session.Close()
	mc := &MongoCollection{collection: db.C("obj")}
	var docs interface{}
	var selector interface{}
	var count int
	var expect int
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "joe", "age": 31}
	mc.insertOne(docs)
	selector = types.M{"name": "joe"}
	count = mc.count(selector, nil)
	expect = 2
	if count != expect {
		t.Error("expect:", expect, "get result:", count)
	}
	/********************************************************/
	selector = types.M{"name": "jack"}
	count = mc.count(selector, nil)
	expect = 1
	if count != expect {
		t.Error("expect:", expect, "get result:", count)
	}
	/********************************************************/
	selector = types.M{"name": "tom"}
	count = mc.count(selector, nil)
	expect = 0
	if count != expect {
		t.Error("expect:", expect, "get result:", count)
	}
	mc.drop()
}

func Test_findOneAndUpdate(t *testing.T) {
	db := openDB()
	defer db.Session.Close()
	mc := &MongoCollection{collection: db.C("obj")}
	var docs interface{}
	var selector interface{}
	var update interface{}
	var obj types.M
	var result []types.M
	var err error
	var expect interface{}
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "joe", "age": 31}
	mc.insertOne(docs)
	selector = types.M{"name": "joe"}
	update = types.M{"$set": types.M{"age": 35}}
	obj = mc.findOneAndUpdate(selector, update)
	expect = types.M{"_id": "001", "name": "joe", "age": 35}
	if reflect.DeepEqual(obj, expect) == false {
		t.Error("expect:", expect, "get result:", obj)
	}
	result, err = mc.rawFind(selector, nil)
	expect = []types.M{
		types.M{"_id": "001", "name": "joe", "age": 35},
		types.M{"_id": "003", "name": "joe", "age": 31},
	}
	if err != nil || result == nil || len(result) != 2 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result, err)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "joe", "age": 31}
	mc.insertOne(docs)
	selector = types.M{"name": "tom"}
	update = types.M{"$set": types.M{"age": 35}}
	obj = mc.findOneAndUpdate(selector, update)
	expect = types.M{}
	if reflect.DeepEqual(obj, expect) == false {
		t.Error("expect:", expect, "get result:", obj)
	}
	result, err = mc.rawFind(nil, nil)
	expect = []types.M{
		types.M{"_id": "001", "name": "joe", "age": 25},
		types.M{"_id": "002", "name": "jack", "age": 30},
		types.M{"_id": "003", "name": "joe", "age": 31},
	}
	if err != nil || result == nil || len(result) != 3 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result, err)
	}
	mc.drop()
}

func Test_insertOne(t *testing.T) {
	db := openDB()
	defer db.Session.Close()
	mc := &MongoCollection{collection: db.C("obj")}
	var docs interface{}
	var result []types.M
	var err error
	var expect interface{}
	/********************************************************/
	docs = types.M{
		"_id":      "1024",
		"name":     "joe",
		"age":      12,
		"male":     true,
		"number":   12.5,
		"location": types.S{10, 20.5},
		"user":     types.M{"name": "jack"},
	}
	mc.insertOne(docs)
	result, err = mc.rawFind(nil, nil)
	expect = types.M{
		"_id":      "1024",
		"name":     "joe",
		"age":      12,
		"male":     true,
		"number":   12.5,
		"location": []interface{}{10, 20.5},
		"user":     types.M{"name": "jack"},
	}
	if err != nil || result == nil || len(result) != 1 || reflect.DeepEqual(result[0], expect) == false {
		t.Error("expect:", expect, "get result:", result, err)
	}
	mc.drop()
}

func Test_upsertOne(t *testing.T) {
	db := openDB()
	defer db.Session.Close()
	mc := &MongoCollection{collection: db.C("obj")}
	var docs interface{}
	var selector interface{}
	var update interface{}
	var result []types.M
	var err error
	var expect interface{}
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	selector = types.M{"name": "joe"}
	update = types.M{"$set": types.M{"age": 35}}
	err = mc.upsertOne(selector, update)
	if err != nil {
		t.Error("expect:", nil, "get result:", err)
	}
	result, err = mc.rawFind(selector, nil)
	expect = types.M{"_id": "001", "name": "joe", "age": 35}
	if err != nil || result == nil || len(result) != 1 || reflect.DeepEqual(result[0], expect) == false {
		t.Error("expect:", expect, "get result:", result, err)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	selector = types.M{"name": "tom"}
	update = types.M{"$set": types.M{"_id": "003", "age": 35}}
	err = mc.upsertOne(selector, update)
	if err != nil {
		t.Error("expect:", nil, "get result:", err)
	}
	result, err = mc.rawFind(selector, nil)
	expect = types.M{"_id": "003", "name": "tom", "age": 35}
	if err != nil || result == nil || len(result) != 1 || reflect.DeepEqual(result[0], expect) == false {
		t.Error("expect:", expect, "get result:", result, err)
	}
	mc.drop()
}

func Test_updateOne(t *testing.T) {
	db := openDB()
	defer db.Session.Close()
	mc := &MongoCollection{collection: db.C("obj")}
	var docs interface{}
	var selector interface{}
	var update interface{}
	var result []types.M
	var err error
	var expect interface{}
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	selector = types.M{"name": "joe"}
	update = types.M{"$set": types.M{"age": 35}}
	err = mc.updateOne(selector, update)
	if err != nil {
		t.Error("expect:", nil, "get result:", err)
	}
	result, err = mc.rawFind(selector, nil)
	expect = types.M{"_id": "001", "name": "joe", "age": 35}
	if err != nil || result == nil || len(result) != 1 || reflect.DeepEqual(result[0], expect) == false {
		t.Error("expect:", expect, "get result:", result, err)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	selector = types.M{"name": "joe"}
	update = types.M{"$unset": types.M{"age": ""}}
	err = mc.updateOne(selector, update)
	if err != nil {
		t.Error("expect:", nil, "get result:", err)
	}
	result, err = mc.rawFind(selector, nil)
	expect = types.M{"_id": "001", "name": "joe"}
	if err != nil || result == nil || len(result) != 1 || reflect.DeepEqual(result[0], expect) == false {
		t.Error("expect:", expect, "get result:", result, err)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	selector = types.M{"name": "joe"}
	update = types.M{"$inc": types.M{"age": 10}}
	err = mc.updateOne(selector, update)
	if err != nil {
		t.Error("expect:", nil, "get result:", err)
	}
	result, err = mc.rawFind(selector, nil)
	expect = types.M{"_id": "001", "name": "joe", "age": 35}
	if err != nil || result == nil || len(result) != 1 || reflect.DeepEqual(result[0], expect) == false {
		t.Error("expect:", expect, "get result:", result, err)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "skill": types.S{"one"}}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "skill": types.S{"one"}}
	mc.insertOne(docs)
	selector = types.M{"name": "joe"}
	update = types.M{
		"$push": types.M{
			"skill": types.M{
				"$each": types.S{"two", "three"},
			},
		},
	}
	err = mc.updateOne(selector, update)
	if err != nil {
		t.Error("expect:", nil, "get result:", err)
	}
	result, err = mc.rawFind(selector, nil)
	expect = types.M{"_id": "001", "name": "joe", "skill": []interface{}{"one", "two", "three"}}
	if err != nil || result == nil || len(result) != 1 || reflect.DeepEqual(result[0], expect) == false {
		t.Error("expect:", expect, "get result:", result, err)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "skill": types.S{"one"}}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "skill": types.S{"one"}}
	mc.insertOne(docs)
	selector = types.M{"name": "joe"}
	update = types.M{
		"$addToSet": types.M{
			"skill": types.M{
				"$each": types.S{"two", "three"},
			},
		},
	}
	err = mc.updateOne(selector, update)
	if err != nil {
		t.Error("expect:", nil, "get result:", err)
	}
	result, err = mc.rawFind(selector, nil)
	expect = types.M{"_id": "001", "name": "joe", "skill": []interface{}{"one", "two", "three"}}
	if err != nil || result == nil || len(result) != 1 || reflect.DeepEqual(result[0], expect) == false {
		t.Error("expect:", expect, "get result:", result, err)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "skill": types.S{"one"}}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "skill": types.S{"one"}}
	mc.insertOne(docs)
	selector = types.M{"name": "joe"}
	update = types.M{
		"$addToSet": types.M{
			"skill": types.M{
				"$each": types.S{"one", "two"},
			},
		},
	}
	err = mc.updateOne(selector, update)
	if err != nil {
		t.Error("expect:", nil, "get result:", err)
	}
	result, err = mc.rawFind(selector, nil)
	expect = types.M{"_id": "001", "name": "joe", "skill": []interface{}{"one", "two"}}
	if err != nil || result == nil || len(result) != 1 || reflect.DeepEqual(result[0], expect) == false {
		t.Error("expect:", expect, "get result:", result, err)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "skill": types.S{"one", "two", "three"}}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "skill": types.S{"one"}}
	mc.insertOne(docs)
	selector = types.M{"name": "joe"}
	update = types.M{
		"$pullAll": types.M{
			"skill": types.S{"two", "three"},
		},
	}
	err = mc.updateOne(selector, update)
	if err != nil {
		t.Error("expect:", nil, "get result:", err)
	}
	result, err = mc.rawFind(selector, nil)
	expect = types.M{"_id": "001", "name": "joe", "skill": []interface{}{"one"}}
	if err != nil || result == nil || len(result) != 1 || reflect.DeepEqual(result[0], expect) == false {
		t.Error("expect:", expect, "get result:", result, err)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	selector = types.M{"name": "tom"}
	update = types.M{"$set": types.M{"age": 35}}
	err = mc.updateOne(selector, update)
	if err.Error() != "not found" {
		t.Error("expect:", nil, "get result:", err)
	}
	mc.drop()
}

func Test_updateMany(t *testing.T) {
	db := openDB()
	defer db.Session.Close()
	mc := &MongoCollection{collection: db.C("obj")}
	var docs interface{}
	var selector interface{}
	var update interface{}
	var result []types.M
	var err error
	var expect interface{}
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "joe", "age": 31}
	mc.insertOne(docs)
	selector = types.M{"name": "joe"}
	update = types.M{"$set": types.M{"age": 35}}
	err = mc.updateMany(selector, update)
	if err != nil {
		t.Error("expect:", nil, "get result:", err)
	}
	result, err = mc.rawFind(selector, nil)
	expect = []types.M{
		types.M{"_id": "001", "name": "joe", "age": 35},
		types.M{"_id": "003", "name": "joe", "age": 35},
	}
	if err != nil || result == nil || len(result) != 2 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result, err)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "joe", "age": 31}
	mc.insertOne(docs)
	selector = types.M{"name": "tom"}
	update = types.M{"$set": types.M{"age": 35}}
	err = mc.updateMany(selector, update)
	if err != nil {
		t.Error("expect:", nil, "get result:", err)
	}
	result, err = mc.rawFind(nil, nil)
	expect = []types.M{
		types.M{"_id": "001", "name": "joe", "age": 25},
		types.M{"_id": "002", "name": "jack", "age": 30},
		types.M{"_id": "003", "name": "joe", "age": 31},
	}
	if err != nil || result == nil || len(result) != 3 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result, err)
	}
	mc.drop()
}

func Test_deleteOne(t *testing.T) {
	db := openDB()
	defer db.Session.Close()
	mc := &MongoCollection{collection: db.C("obj")}
	var docs interface{}
	var selector interface{}
	var result []types.M
	var err error
	var expect interface{}
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	selector = types.M{"name": "joe"}
	err = mc.deleteOne(selector)
	if err != nil {
		t.Error("expect:", nil, "get result:", err)
	}
	result, err = mc.rawFind(nil, nil)
	expect = types.M{"_id": "002", "name": "jack", "age": 30}
	if err != nil || result == nil || len(result) != 1 || reflect.DeepEqual(result[0], expect) == false {
		t.Error("expect:", expect, "get result:", result, err)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	selector = types.M{"name": "tom"}
	err = mc.deleteOne(selector)
	if err.Error() != "not found" {
		t.Error("expect:", nil, "get result:", err)
	}
	mc.drop()
}

func Test_deleteMany(t *testing.T) {
	db := openDB()
	defer db.Session.Close()
	mc := &MongoCollection{collection: db.C("obj")}
	var docs interface{}
	var selector interface{}
	var result []types.M
	var count int
	var err error
	var expect interface{}
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "joe", "age": 31}
	mc.insertOne(docs)
	selector = types.M{"name": "joe"}
	count, err = mc.deleteMany(selector)
	if err != nil || count != 2 {
		t.Error("expect:", nil, "get result:", err, count)
	}
	result, err = mc.rawFind(nil, nil)
	expect = []types.M{
		types.M{"_id": "002", "name": "jack", "age": 30},
	}
	if err != nil || result == nil || len(result) != 1 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result, err)
	}
	mc.drop()
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "joe", "age": 31}
	mc.insertOne(docs)
	selector = types.M{"name": "tom"}
	count, err = mc.deleteMany(selector)
	if err != nil || count != 0 {
		t.Error("expect:", nil, "get result:", err, count)
	}
	result, err = mc.rawFind(nil, nil)
	expect = []types.M{
		types.M{"_id": "001", "name": "joe", "age": 25},
		types.M{"_id": "002", "name": "jack", "age": 30},
		types.M{"_id": "003", "name": "joe", "age": 31},
	}
	if err != nil || result == nil || len(result) != 3 || reflect.DeepEqual(result, expect) == false {
		t.Error("expect:", expect, "get result:", result, err)
	}
	mc.drop()
}

func Test_drop(t *testing.T) {
	db := openDB()
	defer db.Session.Close()
	mc := &MongoCollection{collection: db.C("obj")}
	var docs interface{}
	var result []types.M
	var err error
	var expect interface{}
	/********************************************************/
	docs = types.M{"_id": "001", "name": "joe", "age": 25}
	mc.insertOne(docs)
	docs = types.M{"_id": "002", "name": "jack", "age": 30}
	mc.insertOne(docs)
	docs = types.M{"_id": "003", "name": "joe", "age": 31}
	mc.insertOne(docs)
	err = mc.drop()
	if err != nil {
		t.Error("expect:", nil, "get result:", err)
	}
	result, err = mc.rawFind(nil, nil)
	if err != nil || (result != nil && len(result) != 0) {
		t.Error("expect:", expect, "get result:", result, err)
	}
	/********************************************************/
	err = mc.drop()
	expectErr := errors.New("ns not found")
	if err == nil || err.Error() != expectErr.Error() {
		t.Error("expect:", expect, "get result:", err)
	}
}

func openDB() *mgo.Database {
	session, err := mgo.Dial("192.168.99.100:27017/test")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session.DB("")
}
