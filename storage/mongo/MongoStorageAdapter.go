package mongo

import (
	"regexp"
	"strings"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/storage"
	"github.com/lfq7413/tomato/types"

	"gopkg.in/mgo.v2"
)

const mongoSchemaCollectionName = "_SCHEMA"

// MongoAdapter mongo 数据库适配器
type MongoAdapter struct {
	collectionPrefix string
	collectionList   []string
	transform        *Transform
}

// NewMongoAdapter ...
func NewMongoAdapter(collectionPrefix string) *MongoAdapter {
	return &MongoAdapter{
		collectionPrefix: collectionPrefix,
		collectionList:   []string{},
		transform:        NewTransform(),
	}
}

// collection 获取指定表的操作对象
func (m *MongoAdapter) collection(name string) *mgo.Collection {
	return storage.TomatoDB.MongoDatabase.C(name)
}

// adaptiveCollection 组装 mongo 表操作对象
func (m *MongoAdapter) adaptiveCollection(name string) *MongoCollection {
	return &MongoCollection{
		collection: m.collection(m.collectionPrefix + name),
	}
}

// SchemaCollection 组装 _SCHEMA 表操作对象
func (m *MongoAdapter) SchemaCollection() storage.SchemaCollection {
	mongoCollection := &MongoCollection{
		collection: m.collection(m.collectionPrefix + mongoSchemaCollectionName),
	}
	return &MongoSchemaCollection{
		collection: mongoCollection,
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

// DeleteOneSchema 删除指定表
func (m *MongoAdapter) DeleteOneSchema(name string) error {
	// TODO 处理类不存在时的情况
	return m.collection(m.collectionPrefix + name).DropCollection()
}

// DeleteAllSchemas 删除所有表，仅用于测试
func (m *MongoAdapter) DeleteAllSchemas() error {
	collections := storageAdapterAllCollections(m)
	for _, collection := range collections {
		err := collection.drop()
		if err != nil {
			return err
		}
	}

	return nil
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
	collection := m.adaptiveCollection(className)
	err := collection.updateMany(types.M{}, collectionUpdate)
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

// CreateObject 创建对象
func (m *MongoAdapter) CreateObject(className string, object types.M, schema types.M) error {
	mongoObject, err := m.transform.parseObjectToMongoObjectForCreate(className, object, schema)
	if err != nil {
		return err
	}
	coll := m.adaptiveCollection(className)
	return coll.insertOne(mongoObject)
}

// GetOneSchema ...
func (m *MongoAdapter) GetOneSchema(className string) (types.M, error) {
	return m.SchemaCollection().FindSchema(className)
}

// GetAllSchemas ...
func (m *MongoAdapter) GetAllSchemas() ([]types.M, error) {
	return m.SchemaCollection().GetAllSchemas()
}

// getCollectionNames 获取数据库中当前已经存在的表名
func (m *MongoAdapter) getCollectionNames() []string {
	names, err := storage.TomatoDB.MongoDatabase.CollectionNames()
	if err == nil && names != nil {
		return names
	}
	return []string{}
}

// DeleteObjectsByQuery 删除符合条件的所有对象
func (m *MongoAdapter) DeleteObjectsByQuery(className string, query types.M, schema types.M) error {
	collection := m.adaptiveCollection(className)

	mongoWhere, err := m.transform.transformWhere(className, query, schema)
	if err != nil {
		return err
	}

	n, err := collection.deleteMany(mongoWhere)
	if err != nil {
		return errs.E(errs.InternalServerError, "Database adapter error")
	}
	if n == 0 {
		return errs.E(errs.ObjectNotFound, "Object not found.")
	}

	return nil
}

// UpdateObjectsByQuery ...
func (m *MongoAdapter) UpdateObjectsByQuery(className string, query, schema, update types.M) error {
	mongoUpdate, err := m.transform.transformUpdate(className, update, schema)
	if err != nil {
		return err
	}
	mongoWhere, err := m.transform.transformWhere(className, query, schema)
	if err != nil {
		return err
	}
	coll := m.adaptiveCollection(className)
	return coll.updateMany(mongoWhere, mongoUpdate)
}

// FindOneAndUpdate ...
func (m *MongoAdapter) FindOneAndUpdate(className string, query, schema, update types.M) (types.M, error) {
	mongoUpdate, err := m.transform.transformUpdate(className, update, schema)
	if err != nil {
		return nil, err
	}
	mongoWhere, err := m.transform.transformWhere(className, query, schema)
	if err != nil {
		return nil, err
	}
	coll := m.adaptiveCollection(className)
	return coll.findOneAndUpdate(mongoWhere, mongoUpdate), nil
}

// UpsertOneObject ...
func (m *MongoAdapter) UpsertOneObject(className string, query, schema, update types.M) error {
	mongoUpdate, err := m.transform.transformUpdate(className, update, schema)
	if err != nil {
		return err
	}
	mongoWhere, err := m.transform.transformWhere(className, query, schema)
	if err != nil {
		return err
	}
	coll := m.adaptiveCollection(className)
	return coll.upsertOne(mongoWhere, mongoUpdate)
}

// Find ...
func (m *MongoAdapter) Find(className string, query, schema, options types.M) ([]types.M, error) {
	mongoWhere, err := m.transform.transformWhere(className, query, schema)
	if err != nil {
		return nil, err
	}
	if options["sort"] != nil {
		keys := options["sort"].([]string)
		mongoSort := []string{}
		for _, key := range keys {
			var mongoKey string
			var prefix string

			if strings.HasPrefix(key, "-") {
				prefix = "-"
				key = key[1:]
			}

			mongoKey = prefix + m.transform.transformKey(className, key, schema)
			mongoSort = append(mongoSort, mongoKey)
		}
		options["sort"] = mongoSort
	}

	coll := m.adaptiveCollection(className)
	results := coll.find(mongoWhere, options)
	objects := []types.M{}
	for _, result := range results {
		r, err := m.transform.mongoObjectToParseObject(className, result, schema)
		if err != nil {
			return nil, err
		}
		objects = append(objects, r.(map[string]interface{}))
	}
	return objects, nil
}

// rawFind 仅用于测试
func (m *MongoAdapter) rawFind(className string, query types.M) []types.M {
	coll := m.adaptiveCollection(className)
	return coll.find(query, types.M{})
}

// Count ...
func (m *MongoAdapter) Count(className string, query, schema types.M) (int, error) {
	coll := m.adaptiveCollection(className)
	mongoWhere, err := m.transform.transformWhere(className, query, schema)
	if err != nil {
		return 0, err
	}
	c := coll.count(mongoWhere, types.M{})
	return c, nil
}

func storageAdapterAllCollections(m *MongoAdapter) []*MongoCollection {
	names := m.getCollectionNames()
	collections := []*MongoCollection{}

	for _, v := range names {
		if m, err := regexp.MatchString(`\.system\.`, v); err == nil && m {
			continue
		}
		if strings.HasPrefix(v, m.collectionPrefix) {
			collections = append(collections, m.adaptiveCollection(v[len(m.collectionPrefix):]))
		}
	}

	return collections
}
