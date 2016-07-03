package orm

import (
	"reflect"
	"testing"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
)

func Test_CollectionExists(t *testing.T) {
	// TODO
}

func Test_PurgeCollection(t *testing.T) {
	// LoadSchema
	// TODO
}

func Test_Find(t *testing.T) {
	// LoadSchema
	// reduceRelationKeys
	// reduceInRelation
	// addPointerPermissions
	// addReadACL
	// validateQuery
	// untransformObjectACL
	// filterSensitiveData
	// TODO
}

func Test_Destroy(t *testing.T) {
	// LoadSchema
	// addPointerPermissions
	// addWriteACL
	// validateQuery
	// TODO
}

func Test_Update(t *testing.T) {
	// LoadSchema
	// handleRelationUpdates
	// addPointerPermissions
	// addWriteACL
	// validateQuery
	// transformObjectACL
	// transformAuthData
	// sanitizeDatabaseResult
	// TODO
}

func Test_Create(t *testing.T) {
	// transformObjectACL
	// validateClassName
	// LoadSchema
	// handleRelationUpdates
	// transformAuthData
	// TODO
}

func Test_validateClassName(t *testing.T) {
	// TODO
}

func Test_handleRelationUpdates(t *testing.T) {
	// addRelation
	// removeRelation
	// TODO
}

func Test_addRelation(t *testing.T) {
	// TODO
}

func Test_removeRelation(t *testing.T) {
	// TODO
}

func Test_ValidateObject(t *testing.T) {
	// LoadSchema
	// canAddField
	// TODO
}

func Test_LoadSchema(t *testing.T) {
	// TODO
}

func Test_DeleteEverything(t *testing.T) {
	// TODO
}

func Test_RedirectClassNameForKey(t *testing.T) {
	// LoadSchema
	// TODO
}

func Test_canAddField(t *testing.T) {
	// TODO
}

func Test_reduceRelationKeys(t *testing.T) {
	// relatedIds
	// addInObjectIdsIds
	// TODO
}

func Test_relatedIds(t *testing.T) {
	// TODO
}

func Test_addInObjectIdsIds(t *testing.T) {
	// TODO
}

func Test_addNotInObjectIdsIds(t *testing.T) {
	// TODO
}

func Test_reduceInRelation(t *testing.T) {
	// owningIds
	// addNotInObjectIdsIds
	// addInObjectIdsIds
	// TODO
}

func Test_owningIds(t *testing.T) {
	// TODO
}

func Test_DeleteSchema(t *testing.T) {
	// LoadSchema
	// CollectionExists
	// TODO
}

func Test_addPointerPermissions(t *testing.T) {
	// TODO
}

//////////////////////////////////////////////////////

func Test_sanitizeDatabaseResult(t *testing.T) {
	// TODO
}

func Test_keysForQuery(t *testing.T) {
	// TODO
}

func Test_joinTableName(t *testing.T) {
	// TODO
}

func Test_filterSensitiveData(t *testing.T) {
	// TODO
}

func Test_addWriteACL(t *testing.T) {
	// TODO
}

func Test_addReadACL(t *testing.T) {
	// TODO
}

func Test_validateQuery(t *testing.T) {
	// TODO
}

func Test_transformObjectACL(t *testing.T) {
	// TODO
}

func Test_untransformObjectACL(t *testing.T) {
	// TODO
}

func Test_transformAuthData(t *testing.T) {
	// TODO
}

func Test_flattenUpdateOperatorsForCreate(t *testing.T) {
	var object types.M
	var err error
	var expect interface{}
	/**********************************************************/
	object = nil
	err = flattenUpdateOperatorsForCreate(object)
	expect = nil
	if err != nil || object != nil {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{"key": "value"}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{"key": "value"}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{"key": "value"},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{
		"key": types.M{"key": "value"},
	}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":   "Increment",
			"amount": 10.24,
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{
		"key": 10.24,
	}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":   "Increment",
			"amount": 1024,
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{
		"key": 1024,
	}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":   "Increment",
			"amount": "hello",
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = errs.E(errs.InvalidJSON, "objects to add must be an number")
	if err == nil || reflect.DeepEqual(err, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op": "Increment",
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = errs.E(errs.InvalidJSON, "objects to add must be an number")
	if err == nil || reflect.DeepEqual(err, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":    "Add",
			"objects": types.S{"abc", "def"},
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{
		"key": types.S{"abc", "def"},
	}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":    "Add",
			"objects": "hello",
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = errs.E(errs.InvalidJSON, "objects to add must be an array")
	if err == nil || reflect.DeepEqual(err, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":    "AddUnique",
			"objects": types.S{"abc", "def"},
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{
		"key": types.S{"abc", "def"},
	}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":    "AddUnique",
			"objects": "hello",
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = errs.E(errs.InvalidJSON, "objects to add must be an array")
	if err == nil || reflect.DeepEqual(err, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":    "Remove",
			"objects": types.S{"abc", "def"},
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{
		"key": types.S{},
	}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op":    "Remove",
			"objects": "hello",
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = errs.E(errs.InvalidJSON, "objects to add must be an array")
	if err == nil || reflect.DeepEqual(err, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op": "Delete",
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = types.M{}
	if err != nil || reflect.DeepEqual(object, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
	/**********************************************************/
	object = types.M{
		"key": types.M{
			"__op": "Other",
		},
	}
	err = flattenUpdateOperatorsForCreate(object)
	expect = errs.E(errs.CommandUnavailable, "The Other operator is not supported yet.")
	if err == nil || reflect.DeepEqual(err, expect) == false {
		t.Error("expect:", expect, "result:", object, err)
	}
}
