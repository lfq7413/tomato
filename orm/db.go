package orm

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

var (
	// TomatoDB ...
	TomatoDB *Database
)

func init() {
}

// Database is sds.
type Database struct {
	Session  *mgo.Session
	Database *mgo.Database
}

// NewDatabase ...
func NewDatabase(url string) *Database {
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	fmt.Println("connect success")
	session.SetMode(mgo.Monotonic, true)
	database := session.DB("")
	db := &Database{Session: session, Database: database}
	return db
}

// OpenDB ...
func OpenDB(url string) {
	TomatoDB = NewDatabase(url)
}

// CloseDB ...
func CloseDB() {
	TomatoDB.Session.Close()
}

// Insert ...
func (d *Database) Insert(collection string, docs interface{}) error {
	err := d.Database.C(collection).Insert(docs)
	return err
}

// FindOne ...
func (d *Database) FindOne(collection string, query interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := d.Database.C(collection).Find(query).One(&result)
	return result, err
}
