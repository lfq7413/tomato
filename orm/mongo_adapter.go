package orm

import "gopkg.in/mgo.v2"
import "strings"

const mongoSchemaCollectionName = "_SCHEMA"

// MongoAdapter ...
type MongoAdapter struct {
	collectionList []string
}

func (m *MongoAdapter) collection(name string) *mgo.Collection {
	return TomatoDB.Database.C(name)
}

func (m *MongoAdapter) adaptiveCollection(name string) *MongoCollection {
	return &MongoCollection{
		collection: m.collection(name),
	}
}

func (m *MongoAdapter) schemaCollection() *MongoSchemaCollection {
	return &MongoSchemaCollection{
		collection: m.adaptiveCollection(mongoSchemaCollectionName),
	}
}

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

func (m *MongoAdapter) dropCollection(name string) error {
	return m.collection(name).DropCollection()
}

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
