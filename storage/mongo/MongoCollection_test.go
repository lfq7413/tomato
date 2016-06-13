package mongo

import (
	"reflect"
	"testing"

	"github.com/lfq7413/tomato/types"

	"gopkg.in/mgo.v2"
)

func Test_find(t *testing.T) {
	// TODO
}

func Test_rawFind(t *testing.T) {
	// TODO
}

func Test_count(t *testing.T) {
	// TODO
}

func Test_findOneAndUpdate(t *testing.T) {
	// TODO
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
	if err != nil || result == nil || len(result) != 1 {
		t.Error("expect:", expect, "get result:", result, err)
	}
	if reflect.DeepEqual(result[0], expect) == false {
		t.Error("expect:", expect, "get result:", result[0], err)
	}

	mc.drop()
}

func Test_upsertOne(t *testing.T) {
	// TODO
}

func Test_updateOne(t *testing.T) {
	// TODO
}

func Test_updateMany(t *testing.T) {
	// TODO
}

func Test_deleteOne(t *testing.T) {
	// TODO
}

func Test_deleteMany(t *testing.T) {
	// TODO
}

func Test_drop(t *testing.T) {
	// TODO
}

func openDB() *mgo.Database {
	session, err := mgo.Dial("192.168.99.100:27017/test")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session.DB("")
}
