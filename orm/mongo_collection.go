package orm

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoCollection ...
type MongoCollection struct {
	collection *mgo.Collection
}

func (m *MongoCollection) find(query interface{}, options map[string]interface{}) []bson.M {
	result, err := m.rawFind(query, options)
	if err != nil || result == nil {
		return []bson.M{}
	}
	return result
}

func (m *MongoCollection) rawFind(query interface{}, options map[string]interface{}) ([]bson.M, error) {
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
	var result []bson.M
	err := q.All(&result)
	return result, err
}

func (m *MongoCollection) count(query interface{}, options map[string]interface{}) int {
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

func (m *MongoCollection) findOneAndUpdate() {

}

func (m *MongoCollection) insertOne() {

}

func (m *MongoCollection) upsertOne() {

}

func (m *MongoCollection) updateOne() {

}

func (m *MongoCollection) updateMany() {

}

func (m *MongoCollection) deleteOne() {

}

func (m *MongoCollection) deleteMany() {

}

func (m *MongoCollection) drop() {

}
