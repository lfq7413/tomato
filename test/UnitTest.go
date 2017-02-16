package test

import (
	"gopkg.in/mgo.v2"
)

// MongoDBTestURL ...
const MongoDBTestURL = "192.168.99.100:27017/test"

// OpenMongoDBForTest ...
func OpenMongoDBForTest() *mgo.Database {
	session, err := mgo.Dial(MongoDBTestURL)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session.DB("")
}
