package storage

import (
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/test"
	"gopkg.in/mgo.v2"
)

// OpenMongoDB 打开 MongoDB
func OpenMongoDB() *mgo.Database {
	// 此处仅用于测试
	if config.TConfig.DatabaseURI == "" {
		config.TConfig.DatabaseURI = test.MongoDBTestURL
	}

	session, err := mgo.Dial(config.TConfig.DatabaseURI)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session.DB("")
}
