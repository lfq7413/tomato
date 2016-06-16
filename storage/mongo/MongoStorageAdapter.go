package mongo

import (
	"regexp"
	"strings"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/storage"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"

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
	rawCollection := m.collection(m.collectionPrefix + name)
	return newMongoCollection(rawCollection)
}

// schemaCollection 组装 _SCHEMA 表操作对象
func (m *MongoAdapter) schemaCollection() *MongoSchemaCollection {
	collection := m.adaptiveCollection(mongoSchemaCollectionName)
	return newMongoSchemaCollection(collection)
}

// ClassExists 检测数据库中是否存在指定类
func (m *MongoAdapter) ClassExists(name string) bool {
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

// SetClassLevelPermissions 设置类级别权限
func (m *MongoAdapter) SetClassLevelPermissions(className string, CLPs types.M) error {
	schemaCollection := m.schemaCollection()
	update := types.M{
		"$set": types.M{
			"_metadata": types.M{
				"class_permissions": CLPs,
			},
		},
	}
	return schemaCollection.updateSchema(className, update)
}

// CreateClass 创建类
func (m *MongoAdapter) CreateClass(className string, schema types.M) (types.M, error) {
	schemaCollection := m.schemaCollection()
	return schemaCollection.addSchema(className, utils.M(schema["fields"]), utils.M(schema["classLevelPermissions"]))
}

// AddFieldIfNotExists 添加字段定义
func (m *MongoAdapter) AddFieldIfNotExists(className, fieldName string, fieldType types.M) error {
	schemaCollection := m.schemaCollection()
	return schemaCollection.addFieldIfNotExists(className, fieldName, fieldType)
}

// DeleteClass 删除指定表
func (m *MongoAdapter) DeleteClass(className string) (types.M, error) {
	coll := m.adaptiveCollection(className)
	err := coll.drop()
	if err != nil {
		if err.Error() == "ns not found" {
			return nil, nil
		}
		return nil, err
	}
	schemaCollection := m.schemaCollection()
	return schemaCollection.findAndDeleteSchema(className)
}

// DeleteAllClasses 删除所有表，仅用于测试
func (m *MongoAdapter) DeleteAllClasses() error {
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
func (m *MongoAdapter) DeleteFields(className string, schema types.M, fieldNames []string) error {
	var fields types.M
	if schema != nil {
		fields = utils.M(schema["fields"])
	}
	mongoFormatNames := []string{}
	for _, fieldName := range fieldNames {
		if fields != nil {
			fieldType := utils.M(fields[fieldName])
			if fieldType != nil && utils.S(fieldType["type"]) == "Pointer" {
				mongoFormatNames = append(mongoFormatNames, "_p_"+fieldName)
			}
		}
		mongoFormatNames = append(mongoFormatNames, fieldName)
	}

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
	schemaCollection := m.schemaCollection()
	err = schemaCollection.updateSchema(className, schemaUpdate)
	if err != nil {
		return err
	}
	return nil
}

// CreateObject 创建对象
func (m *MongoAdapter) CreateObject(className string, schema, object types.M) error {
	mongoObject, err := m.transform.parseObjectToMongoObjectForCreate(className, object, schema)
	if err != nil {
		return err
	}
	coll := m.adaptiveCollection(className)
	return coll.insertOne(mongoObject)
}

// GetClass ...
func (m *MongoAdapter) GetClass(className string) (types.M, error) {
	return m.schemaCollection().findSchema(className)
}

// GetAllClasses ...
func (m *MongoAdapter) GetAllClasses() ([]types.M, error) {
	return m.schemaCollection().getAllSchemas()
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
func (m *MongoAdapter) DeleteObjectsByQuery(className string, schema, query types.M) error {
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
func (m *MongoAdapter) UpdateObjectsByQuery(className string, schema, query, update types.M) error {
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
func (m *MongoAdapter) FindOneAndUpdate(className string, schema, query, update types.M) (types.M, error) {
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
func (m *MongoAdapter) UpsertOneObject(className string, schema, query, update types.M) error {
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
func (m *MongoAdapter) Find(className string, schema, query, options types.M) ([]types.M, error) {
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
	results, err := coll.find(mongoWhere, options)
	if err != nil {
		return nil, err
	}
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
func (m *MongoAdapter) rawFind(className string, query types.M) ([]types.M, error) {
	coll := m.adaptiveCollection(className)
	return coll.find(query, types.M{})
}

// Count ...
func (m *MongoAdapter) Count(className string, schema, query types.M) (int, error) {
	coll := m.adaptiveCollection(className)
	mongoWhere, err := m.transform.transformWhere(className, query, schema)
	if err != nil {
		return 0, err
	}
	c := coll.count(mongoWhere, types.M{})
	return c, nil
}

// EnsureUniqueness 创建索引
func (m *MongoAdapter) EnsureUniqueness(className string, schema types.M, fieldNames []string) error {
	if fieldNames == nil {
		return nil
	}
	mongoFieldNames := []string{}
	for _, fieldName := range fieldNames {
		k := m.transform.transformKey(className, fieldName, schema)
		mongoFieldNames = append(mongoFieldNames, k)
	}
	coll := m.adaptiveCollection(className)
	err := coll.ensureSparseUniqueIndexInBackground(mongoFieldNames)
	return err
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
