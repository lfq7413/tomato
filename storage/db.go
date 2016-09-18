package storage

import (
	"github.com/lfq7413/tomato/config"
	"gopkg.in/mgo.v2"
)

// TomatoDB 全局可访问的数据库操作结构体
var TomatoDB *Database

// Database 封装 mongo 数据库对象
type Database struct {
	MongoSession  *mgo.Session
	MongoDatabase *mgo.Database
}

func init() {
	OpenDB()
}

// newMongoDB 创建 MongoDB 数据库连接
func newMongoDB(url string) *Database {
	// 此处仅用于测试
	if url == "" {
		url = "192.168.99.100:27017/test"
	}

	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	database := session.DB("")
	db := &Database{
		MongoSession:  session,
		MongoDatabase: database,
	}
	return db
}

// OpenDB 打开数据库
func OpenDB() {
	TomatoDB = newMongoDB(config.TConfig.DatabaseURI)
}

// CloseDB 关闭数据库
func CloseDB() {
	TomatoDB.MongoSession.Close()
}
