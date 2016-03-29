package orm

import "gopkg.in/mgo.v2"

// MongoCollection ...
type MongoCollection struct {
	collection *mgo.Collection
}
