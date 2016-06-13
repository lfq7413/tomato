package mongo

import (
	"github.com/lfq7413/tomato/types"
	"gopkg.in/mgo.v2"
)

// MongoCollection mongo 表操作对象
type MongoCollection struct {
	collection *mgo.Collection
}

// find 执行查找操作，自动添加索引
func (m *MongoCollection) find(query interface{}, options types.M) []types.M {
	result, err := m.rawFind(query, options)
	if err != nil || result == nil {
		return []types.M{}
	}
	// TODO 添加 geo 索引
	return result
}

// rawFind 执行原始查找操作，查找选项包括 sort、skip、limit
func (m *MongoCollection) rawFind(query interface{}, options types.M) ([]types.M, error) {
	if options == nil {
		options = types.M{}
	}
	q := m.collection.Find(query)
	if options["sort"] != nil {
		if sort, ok := options["sort"].([]string); ok {
			q = q.Sort(sort...)
		}
	}
	if options["skip"] != nil {
		if skip, ok := options["skip"].(float64); ok {
			q = q.Skip(int(skip))
		} else if skip, ok := options["skip"].(int); ok {
			q = q.Skip(skip)
		}
	}
	if options["limit"] != nil {
		if limit, ok := options["limit"].(float64); ok {
			q = q.Limit(int(limit))
		} else if limit, ok := options["limit"].(int); ok {
			q = q.Limit(limit)
		}
	}
	var result []types.M
	err := q.All(&result)
	return result, err
}

// count 执行 count 操作
func (m *MongoCollection) count(query interface{}, options types.M) int {
	if options == nil {
		options = types.M{}
	}
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

// findOneAndUpdate 查找并更新一个对象，返回更新后的对象
func (m *MongoCollection) findOneAndUpdate(selector interface{}, update interface{}) types.M {

	var result types.M
	change := mgo.Change{
		Update:    update,
		ReturnNew: true,
	}
	info, err := m.collection.Find(selector).Apply(change, &result)
	if err != nil || info.Updated == 0 {
		return types.M{}
	}

	return result
}

// insertOne 插入一个对象
func (m *MongoCollection) insertOne(docs interface{}) error {
	return m.collection.Insert(docs)
}

// upsertOne 更新一个对象，如果要更新的对象不存在，则插入该对象
func (m *MongoCollection) upsertOne(selector interface{}, update interface{}) error {
	_, err := m.collection.Upsert(selector, update)
	return err
}

// updateOne 更新一个对象
func (m *MongoCollection) updateOne(selector interface{}, update interface{}) error {
	return m.collection.Update(selector, update)
}

// updateMany 更新多个对象
func (m *MongoCollection) updateMany(selector interface{}, update interface{}) error {
	_, err := m.collection.UpdateAll(selector, update)
	return err
}

// deleteOne 删除一个对象
func (m *MongoCollection) deleteOne(selector interface{}) error {
	return m.collection.Remove(selector)
}

// deleteMany 删除多个对象
func (m *MongoCollection) deleteMany(selector interface{}) (int, error) {
	info, err := m.collection.RemoveAll(selector)
	if err != nil {
		return 0, err
	}
	n := info.Removed
	return n, nil
}

// drop 删除当前表
func (m *MongoCollection) drop() error {
	return m.collection.DropCollection()
}
