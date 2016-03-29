package orm

import (
	"fmt"

	"github.com/lfq7413/tomato/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	// TomatoDB ...
	TomatoDB *Database
)

func init() {
}

// Database ...
type Database struct {
	Session        *mgo.Session
	Database       *mgo.Database
	collectionList []string
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
	db := &Database{
		Session:        session,
		Database:       database,
		collectionList: []string{},
	}
	return db
}

// OpenDB ...
func OpenDB() {
	TomatoDB = NewDatabase(config.TConfig.DatabaseURI)
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
func (d *Database) FindOne(collection string, query interface{}) (bson.M, error) {
	var result bson.M
	err := d.Database.C(collection).Find(query).One(&result)
	return result, err
}

// Update ...
func (d *Database) Update(collection string, selector interface{}, update interface{}) error {
	err := d.Database.C(collection).Update(selector, update)
	return err
}

// Find ...
func (d *Database) Find(collection string, query interface{}) ([]bson.M, error) {
	var result []bson.M
	err := d.Database.C(collection).Find(query).All(&result)
	return result, err
}

// Remove ...
func (d *Database) Remove(collection string, selector interface{}) error {
	err := d.Database.C(collection).Remove(selector)
	return err
}

func (d *Database) collectionExists(className string) bool {
	// 先在内存中查询
	for _, v := range d.collectionList {
		if v == className {
			return true
		}
	}
	// 内存中不存在，则去数据库中查询一次，更新到内存中
	d.collectionList = d.getCollectionNames()
	if d.collectionList == nil {
		d.collectionList = []string{}
	}
	for _, v := range d.collectionList {
		if v == className {
			return true
		}
	}
	return false
}

func (d *Database) getCollectionNames() []string {
	names, err := d.Database.CollectionNames()
	if err == nil {
		return names
	}
	return []string{}
}
