package orm

import "gopkg.in/mgo.v2/bson"

func transformKey(schema *Schema, className, key string) string {
	// TODO
	return ""
}

// transformCreate ...
func transformCreate(schema *Schema, className string, create bson.M) bson.M {
	// TODO
	return nil
}

func transformWhere(schema *Schema, className string, where bson.M) bson.M {
	// TODO
	return nil
}

func transformUpdate(schema *Schema, className string, update bson.M) bson.M {
	// TODO
	return nil
}

func untransformObjectT(schema *Schema, className string, mongoObject interface{}, isNestedObject bool) interface{} {
	return nil
}
