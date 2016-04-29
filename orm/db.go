package orm

import (
	"fmt"

	"github.com/lfq7413/tomato/config"
	"gopkg.in/mgo.v2"
)

var (
	// TomatoDB 全局可访问的数据库操作结构体
	TomatoDB *Database
)

func init() {
}

// Database 封装 mongo 数据库对象
type Database struct {
	Session        *mgo.Session
	Database       *mgo.Database
	collectionList []string
}

// NewDatabase 创建数据库结构体
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

// OpenDB 打开数据库
func OpenDB() {
	TomatoDB = NewDatabase(config.TConfig.DatabaseURI)
}

// CloseDB 关闭数据库
func CloseDB() {
	TomatoDB.Session.Close()
}

// getCollectionNames 获取数据库中当前已经存在的表名
func (d *Database) getCollectionNames() []string {
	names, err := d.Database.CollectionNames()
	if err == nil && names != nil {
		return names
	}
	return []string{}
}
