package orm

import (
	"strings"

	"github.com/lfq7413/tomato/types"

	"gopkg.in/mgo.v2"
)

const mongoSchemaCollectionName = "_SCHEMA"

// MongoAdapter mongo 数据库适配器
type MongoAdapter struct {
	collectionPrefix string
	collectionList   []string
	transform        *MongoTransform
}

// NewMongoAdapter ...
func NewMongoAdapter(collectionPrefix string) *MongoAdapter {
	return &MongoAdapter{
		collectionPrefix: collectionPrefix,
		collectionList:   []string{},
		transform:        NewMongoTransform(),
	}
}

// collection 获取指定表的操作对象
func (m *MongoAdapter) collection(name string) *mgo.Collection {
	return TomatoDB.Database.C(name)
}

// adaptiveCollection 组装 mongo 表操作对象
func (m *MongoAdapter) adaptiveCollection(name string) *MongoCollection {
	return &MongoCollection{
		collection: m.collection(m.collectionPrefix + name),
		transform:  m.transform,
	}
}

// schemaCollection 组装 _SCHEMA 表操作对象
func (m *MongoAdapter) schemaCollection() *MongoSchemaCollection {
	return &MongoSchemaCollection{
		collection: m.adaptiveCollection(mongoSchemaCollectionName),
		transform:  m.transform,
	}
}

// collectionExists 检测数据库中是否存在指定表
func (m *MongoAdapter) collectionExists(name string) bool {
	name = m.collectionPrefix + name
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
	return m.collection(m.collectionPrefix + name).DropCollection()
}

// allCollections 查找包含指定前缀的表集合，仅用于测试
func (m *MongoAdapter) allCollections() []*mgo.Collection {
	names := TomatoDB.getCollectionNames()
	collections := []*mgo.Collection{}

	for _, v := range names {
		if strings.HasPrefix(v, m.collectionPrefix) {
			collections = append(collections, m.collection(v))
		}
	}

	return collections
}

// deleteFields 删除字段
func (m *MongoAdapter) deleteFields(className string, fieldNames, pointerFieldNames []string) error {
	// 查找非指针字段名
	nonPointerFieldNames := []string{}
	for _, fieldName := range fieldNames {
		in := false
		for _, pointerFieldName := range pointerFieldNames {
			if fieldName == pointerFieldName {
				in = true
				break
			}
		}
		if in == false {
			nonPointerFieldNames = append(nonPointerFieldNames, fieldName)
		}
	}
	// 转换为 mongo 格式
	var mongoFormatNames []string
	for _, pointerFieldName := range pointerFieldNames {
		nonPointerFieldNames = append(nonPointerFieldNames, "_p_"+pointerFieldName)
	}
	mongoFormatNames = nonPointerFieldNames

	// 组装表数据更新语句
	unset := types.M{}
	for _, name := range mongoFormatNames {
		unset[name] = nil
	}
	collectionUpdate := types.M{"$unset": unset}

	// 组装 schema 更新语句
	unset2 := types.M{}
	for _, name := range fieldNames {
		unset[name] = nil
	}
	schemaUpdate := types.M{"$unset": unset2}

	// 更新表数据
	collection := m.adaptiveCollection(className)
	err := collection.UpdateMany(types.M{}, collectionUpdate)
	if err != nil {
		return err
	}
	// 更新 schema
	schemaCollection := m.schemaCollection()
	err = schemaCollection.updateSchema(className, schemaUpdate)
	if err != nil {
		return err
	}
	return nil
}
