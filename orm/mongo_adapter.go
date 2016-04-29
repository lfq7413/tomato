package orm

import (
	"strings"

	"gopkg.in/mgo.v2"
)

const mongoSchemaCollectionName = "_SCHEMA"

// MongoAdapter mongo 数据库适配器
type MongoAdapter struct {
	collectionList []string
}

// collection 获取指定表的操作对象
func (m *MongoAdapter) collection(name string) *mgo.Collection {
	return TomatoDB.Database.C(name)
}

// adaptiveCollection 组装 mongo 表操作对象
func (m *MongoAdapter) adaptiveCollection(name string) *MongoCollection {
	return &MongoCollection{
		collection: m.collection(name),
	}
}

// schemaCollection 组装 _SCHEMA 表操作对象
func (m *MongoAdapter) schemaCollection() *MongoSchemaCollection {
	return &MongoSchemaCollection{
		collection: m.adaptiveCollection(mongoSchemaCollectionName),
	}
}

// collectionExists 检测数据库中是否存在指定表
func (m *MongoAdapter) collectionExists(name string) bool {
	if m.collectionList == nil {
		m.collectionList = TomatoDB.getCollectionNames()
	}
	// 先在内存中查询
	for _, v := range m.collectionList {
		if v == name {
			return true
		}
	}
	// 内存中不存在，则去数据库中查询一次，更新到内存中
	m.collectionList = TomatoDB.getCollectionNames()
	for _, v := range m.collectionList {
		if v == name {
			return true
		}
	}
	return false
}

// dropCollection 删除指定表
func (m *MongoAdapter) dropCollection(name string) error {
	return m.collection(name).DropCollection()
}

// collectionsContaining 查找包含指定前缀的表集合
func (m *MongoAdapter) collectionsContaining(match string) []*mgo.Collection {
	names := TomatoDB.getCollectionNames()
	collections := []*mgo.Collection{}

	for _, v := range names {
		if strings.HasPrefix(v, match) {
			collections = append(collections, m.collection(v))
		}
	}

	return collections
}
