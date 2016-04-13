package orm

import (
	"github.com/lfq7413/tomato/types"
	"gopkg.in/mgo.v2"
)

// MongoCollection ...
type MongoCollection struct {
	collection *mgo.Collection
}

func (m *MongoCollection) find(query interface{}, options map[string]interface{}) []types.M {
	result, err := m.rawFind(query, options)
	if err != nil || result == nil {
		return []types.M{}
	}
	return result
}

func (m *MongoCollection) rawFind(query interface{}, options map[string]interface{}) ([]types.M, error) {
	q := m.collection.Find(query)
	if options["sort"] != nil {
		if sort, ok := options["sort"].([]string); ok {
			q = q.Sort(sort...)
		}
	}
	if options["skip"] != nil {
		if skip, ok := options["skip"].(float64); ok {
			q = q.Skip(int(skip))
		}
	}
	if options["limit"] != nil {
		if limit, ok := options["limit"].(float64); ok {
			q = q.Limit(int(limit))
		}
	}
	var result []types.M
	err := q.All(&result)
	return result, err
}

// Count ...
func (m *MongoCollection) Count(query interface{}, options map[string]interface{}) int {
	q := m.collection.Find(query)
	if options["sort"] != nil {
		if sort, ok := options["sort"].([]string); ok {
			q = q.Sort(sort...)
		}
	}
	if options["skip"] != nil {
		if skip, ok := options["skip"].(float64); ok {
			q = q.Skip(int(skip))
		}
	}
	if options["limit"] != nil {
		if limit, ok := options["limit"].(float64); ok {
			q = q.Limit(int(limit))
		}
	}
	n, err := q.Count()
	if err != nil {
		return 0
	}
	return n
}

// 当前框架 Update 时的 selector 中，只包含 objectid ，所以更新之后再去查找，找到的为同一对象
func (m *MongoCollection) findOneAndUpdate(selector interface{}, update interface{}) types.M {
	// TODO 使用 Apply 实现
	err := m.collection.Update(selector, update)
	if err != nil {
		return types.M{}
	}
	var result types.M
	err = m.collection.Find(selector).One(&result)
	if err != nil || result == nil {
		return types.M{}
	}
	return result
}

func (m *MongoCollection) insertOne(docs interface{}) error {
	return m.collection.Insert(docs)
}

func (m *MongoCollection) upsertOne(selector interface{}, update interface{}) error {
	_, err := m.collection.Upsert(selector, update)
	return err
}

func (m *MongoCollection) updateOne(selector interface{}, update interface{}) error {
	return m.collection.Update(selector, update)
}

// UpdateMany ...
func (m *MongoCollection) UpdateMany(selector interface{}, update interface{}) error {
	_, err := m.collection.UpdateAll(selector, update)
	return err
}

func (m *MongoCollection) deleteOne(selector interface{}) error {
	return m.collection.Remove(selector)
}

func (m *MongoCollection) deleteMany(selector interface{}) error {
	_, err := m.collection.RemoveAll(selector)
	return err
}

// Drop ...
func (m *MongoCollection) Drop() error {
	return m.collection.DropCollection()
}
