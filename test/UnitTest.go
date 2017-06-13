package test

import (
	"database/sql"
	"log"

	"gopkg.in/mgo.v2"
)

// MongoDBTestURL ...
const MongoDBTestURL = "192.168.99.100:27017/test"

// PostgreSQLTestURL ...
const PostgreSQLTestURL = "postgres://postgres:123456@192.168.99.100:5432/postgres?sslmode=disable"

// OpenMongoDBForTest ...
func OpenMongoDBForTest() *mgo.Database {
	session, err := mgo.Dial(MongoDBTestURL)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session.DB("")
}

var postgresDB *sql.DB

// OpenPostgreSQForTest ...
func OpenPostgreSQForTest() *sql.DB {
	if postgresDB != nil {
		return postgresDB
	}
	db, err := sql.Open("postgres", PostgreSQLTestURL)
	if err != nil {
		log.Fatal(err)
	}
	postgresDB = db
	return db
}
