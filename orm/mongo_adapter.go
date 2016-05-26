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
	return TomatoDB.MongoDatabase.C(name)
}

// AdaptiveCollection 组装 mongo 表操作对象
func (m *MongoAdapter) AdaptiveCollection(name string) *MongoCollection {
	return &MongoCollection{
		collection: m.collection(m.collectionPrefix + name),
		transform:  m.transform,
	}
}

// SchemaCollection 组装 _SCHEMA 表操作对象
func (m *MongoAdapter) SchemaCollection() *MongoSchemaCollection {
	return &MongoSchemaCollection{
		collection: m.AdaptiveCollection(mongoSchemaCollectionName),
		transform:  m.transform,
	}
}

// CollectionExists 检测数据库中是否存在指定表
func (m *MongoAdapter) CollectionExists(name string) bool {
	name = m.collectionPrefix + name
	if m.collectionList == nil {
		m.collectionList = m.getCollectionNames()
	}
	// 先在内存中查询
	for _, v := range m.collectionList {
		if v == name {
			return true
		}
	}
	// 内存中不存在，则去数据库中查询一次，更新到内存中
	m.collectionList = m.getCollectionNames()
	for _, v := range m.collectionList {
		if v == name {
			return true
		}
	}
	return false
}

// DropCollection 删除指定表
func (m *MongoAdapter) DropCollection(name string) error {
	return m.collection(m.collectionPrefix + name).DropCollection()
}

// AllCollections 查找包含指定前缀的表集合，仅用于测试
func (m *MongoAdapter) AllCollections() []*MongoCollection {
	names := m.getCollectionNames()
	collections := []*MongoCollection{}

	for _, v := range names {
		if strings.HasPrefix(v, m.collectionPrefix) {
			collections = append(collections, m.AdaptiveCollection(v[len(m.collectionPrefix):]))
		}
	}

	return collections
}

// DeleteFields 删除字段
func (m *MongoAdapter) DeleteFields(className string, fieldNames, pointerFieldNames []string) error {
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
	collection := m.AdaptiveCollection(className)
	err := collection.UpdateMany(types.M{}, collectionUpdate)
	if err != nil {
		return err
	}
	// 更新 schema
	schemaCollection := m.SchemaCollection()
	err = schemaCollection.UpdateSchema(className, schemaUpdate)
	if err != nil {
		return err
	}
	return nil
}

// GetTransform ...
func (m *MongoAdapter) GetTransform() *MongoTransform {
	return m.transform
}

// getCollectionNames 获取数据库中当前已经存在的表名
func (m *MongoAdapter) getCollectionNames() []string {
	names, err := TomatoDB.MongoDatabase.CollectionNames()
	if err == nil && names != nil {
		return names
	}
	return []string{}
}
