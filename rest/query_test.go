package rest

import (
	"reflect"
	"testing"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
)

func Test_Execute(t *testing.T) {
	// BuildRestWhere
	// runFind
	// runCount
	// handleInclude
	// TODO
}

func Test_BuildRestWhere(t *testing.T) {
	// getUserAndRoleACL
	// redirectClassNameForKey
	// validateClientClassCreation
	// replaceSelect
	// replaceDontSelect
	// replaceInQuery
	// replaceNotInQuery
	// TODO
}

func Test_getUserAndRoleACL(t *testing.T) {
	var auth *Auth
	var q *Query
	/**********************************************************/
	auth = Master()
	q, _ = NewQuery(auth, "user", nil, nil, nil)
	q.getUserAndRoleACL()
	if _, ok := q.findOptions["acl"]; ok {
		t.Error("findOptions[acl] exist")
	}
	/**********************************************************/
	auth = Nobody()
	q, _ = NewQuery(auth, "user", nil, nil, nil)
	q.getUserAndRoleACL()
	if q.findOptions["acl"] != nil {
		t.Error("findOptions[acl] is not nil")
	}
	// TODO
}

func Test_redirectClassNameForKey(t *testing.T) {
	// TODO
}

func Test_validateClientClassCreation(t *testing.T) {
	// TODO
}

func Test_replaceSelect(t *testing.T) {
	// Execute
	// TODO
}

func Test_replaceDontSelect(t *testing.T) {
	// Execute
	// TODO
}

func Test_replaceInQuery(t *testing.T) {
	// Execute
	// TODO
}

func Test_replaceNotInQuery(t *testing.T) {
	// Execute
	// TODO
}

func Test_runFind(t *testing.T) {
	// TODO
}

func Test_runCount(t *testing.T) {
	// TODO
}

func Test_handleInclude(t *testing.T) {
	// includePath
	// TODO
}

/////////////////////////////////////////////////////////////////

func Test_NewQuery(t *testing.T) {
	var auth *Auth
	var className string
	var where types.M
	var options types.M
	var clientSDK map[string]string
	var result *Query
	var err error
	var expectErr error
	var expect *Query
	/**********************************************************/
	auth = Nobody()
	className = "user"
	where = nil
	options = nil
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:              auth,
		className:         "user",
		Where:             types.M{},
		findOptions:       types.M{"acl": nil},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = nil
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:              auth,
		className:         "user",
		Where:             types.M{},
		findOptions:       types.M{},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = &Auth{
		IsMaster: false,
		User: types.M{
			"objectId": "1001",
		},
	}
	className = "user"
	where = nil
	options = nil
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:              auth,
		className:         "user",
		Where:             types.M{},
		findOptions:       types.M{"acl": []string{"1001"}},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Nobody()
	className = "_Session"
	where = nil
	options = nil
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expectErr = errs.E(errs.InvalidSessionToken, "This session token is invalid.")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", result, err)
	}
	/**********************************************************/
	auth = &Auth{
		IsMaster: false,
		User: types.M{
			"objectId": "1001",
		},
	}
	className = "_Session"
	where = nil
	options = nil
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:      auth,
		className: "_Session",
		Where: types.M{
			"$and": types.S{
				types.M{},
				types.M{
					"user": types.M{
						"__type":    "Pointer",
						"className": "_User",
						"objectId":  "1001",
					},
				},
			},
		},
		findOptions:       types.M{"acl": []string{"1001"}},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = types.M{"keys": 1024}
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:              auth,
		className:         "user",
		Where:             types.M{},
		findOptions:       types.M{},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = types.M{"keys": "post"}
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:              auth,
		className:         "user",
		Where:             types.M{},
		findOptions:       types.M{},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{"post", "objectId", "createdAt", "updatedAt"},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = types.M{"keys": "post,user"}
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:              auth,
		className:         "user",
		Where:             types.M{},
		findOptions:       types.M{},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{"post", "user", "objectId", "createdAt", "updatedAt"},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = types.M{"count": true}
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:              auth,
		className:         "user",
		Where:             types.M{},
		findOptions:       types.M{},
		response:          types.M{},
		doCount:           true,
		include:           [][]string{},
		keys:              []string{},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = types.M{"skip": 10}
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:              auth,
		className:         "user",
		Where:             types.M{},
		findOptions:       types.M{"skip": 10},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = types.M{"limit": 10}
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:              auth,
		className:         "user",
		Where:             types.M{},
		findOptions:       types.M{"limit": 10},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = types.M{"order": "post,-user"}
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:              auth,
		className:         "user",
		Where:             types.M{},
		findOptions:       types.M{"sort": []string{"post", "-user"}},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = types.M{"include": "user.session,name.friend"}
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:        auth,
		className:   "user",
		Where:       types.M{},
		findOptions: types.M{},
		response:    types.M{},
		doCount:     false,
		include: [][]string{
			[]string{"name"},
			[]string{"name", "friend"},
			[]string{"user"},
			[]string{"user", "session"},
		},
		keys:              []string{},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = types.M{"redirectClassNameForKey": "post"}
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:              auth,
		className:         "user",
		Where:             types.M{},
		findOptions:       types.M{},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{},
		redirectKey:       "post",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = types.M{"other": "hello"}
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expectErr = errs.E(errs.InvalidJSON, "bad option: other")
	if err == nil || reflect.DeepEqual(expectErr, err) == false {
		t.Error("expect:", expectErr, "result:", result, err)
	}
	/**********************************************************/
	auth = Master()
	className = "user"
	where = nil
	options = types.M{
		"keys":                    "post,user",
		"count":                   true,
		"skip":                    10,
		"limit":                   10,
		"order":                   "post,-user",
		"include":                 "user.session,name.friend",
		"redirectClassNameForKey": "post",
	}
	clientSDK = nil
	result, err = NewQuery(auth, className, where, options, clientSDK)
	expect = &Query{
		auth:      auth,
		className: "user",
		Where:     types.M{},
		findOptions: types.M{
			"skip":  10,
			"limit": 10,
			"sort":  []string{"post", "-user"},
		},
		response: types.M{},
		doCount:  true,
		include: [][]string{
			[]string{"name"},
			[]string{"name", "friend"},
			[]string{"user"},
			[]string{"user", "session"},
		},
		keys:              []string{"post", "user", "objectId", "createdAt", "updatedAt"},
		redirectKey:       "post",
		redirectClassName: "",
		clientSDK:         nil,
	}
	if err != nil || reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result, err)
	}
}

func Test_includePath(t *testing.T) {
	// Execute
	// TODO
}

func Test_findPointers(t *testing.T) {
	var object interface{}
	var path []string
	var result []types.M
	var expect []types.M
	/**********************************************************/
	object = nil
	path = nil
	result = findPointers(object, path)
	expect = []types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	object = "hello"
	path = nil
	result = findPointers(object, path)
	expect = []types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	object = types.M{
		"key": "hello",
	}
	path = nil
	result = findPointers(object, path)
	expect = []types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	object = types.M{
		"__type":   "Pointer",
		"objectId": "1001",
	}
	path = nil
	result = findPointers(object, path)
	expect = []types.M{
		types.M{
			"__type":   "Pointer",
			"objectId": "1001",
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	object = types.M{
		"key": "hello",
	}
	path = []string{"post"}
	result = findPointers(object, path)
	expect = []types.M{}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	object = types.M{
		"post": types.M{
			"__type":   "Pointer",
			"objectId": "1001",
		},
	}
	path = []string{"post"}
	result = findPointers(object, path)
	expect = []types.M{
		types.M{
			"__type":   "Pointer",
			"objectId": "1001",
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	object = types.S{
		types.M{
			"__type":   "Pointer",
			"objectId": "1001",
		},
		types.M{
			"__type":   "Pointer",
			"objectId": "1002",
		},
	}
	path = []string{}
	result = findPointers(object, path)
	expect = []types.M{
		types.M{
			"__type":   "Pointer",
			"objectId": "1001",
		},
		types.M{
			"__type":   "Pointer",
			"objectId": "1002",
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	object = types.M{
		"post": types.S{
			types.M{
				"__type":   "Pointer",
				"objectId": "1001",
			},
			types.M{
				"__type":   "Pointer",
				"objectId": "1002",
			},
		},
	}
	path = []string{"post"}
	result = findPointers(object, path)
	expect = []types.M{
		types.M{
			"__type":   "Pointer",
			"objectId": "1001",
		},
		types.M{
			"__type":   "Pointer",
			"objectId": "1002",
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	object = types.M{
		"post": types.M{
			"user": types.M{
				"__type":   "Pointer",
				"objectId": "1001",
			},
		},
	}
	path = []string{"post", "user"}
	result = findPointers(object, path)
	expect = []types.M{
		types.M{
			"__type":   "Pointer",
			"objectId": "1001",
		},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_replacePointers(t *testing.T) {
	var pointers []types.M
	var replace types.M
	var expect []types.M
	/**********************************************************/
	pointers = nil
	replace = nil
	replacePointers(pointers, replace)
	expect = nil
	if reflect.DeepEqual(expect, pointers) == false {
		t.Error("expect:", expect, "result:", pointers)
	}
	/**********************************************************/
	pointers = []types.M{nil}
	replace = types.M{
		"1001": types.M{
			"key": "hello",
		},
	}
	replacePointers(pointers, replace)
	expect = []types.M{nil}
	if reflect.DeepEqual(expect, pointers) == false {
		t.Error("expect:", expect, "result:", pointers)
	}
	/**********************************************************/
	pointers = []types.M{
		types.M{
			"post": "post",
		},
	}
	replace = types.M{
		"1001": types.M{
			"key": "hello",
		},
	}
	replacePointers(pointers, replace)
	expect = []types.M{
		types.M{
			"post": "post",
		},
	}
	if reflect.DeepEqual(expect, pointers) == false {
		t.Error("expect:", expect, "result:", pointers)
	}
	/**********************************************************/
	pointers = []types.M{
		types.M{
			"objectId": "1002",
			"post":     "post",
		},
	}
	replace = types.M{
		"1001": types.M{
			"key": "hello",
		},
	}
	replacePointers(pointers, replace)
	expect = []types.M{
		types.M{
			"objectId": "1002",
			"post":     "post",
		},
	}
	if reflect.DeepEqual(expect, pointers) == false {
		t.Error("expect:", expect, "result:", pointers)
	}
	/**********************************************************/
	pointers = []types.M{
		types.M{
			"objectId": "1001",
			"post":     "post",
		},
	}
	replace = types.M{
		"1001": types.M{
			"key": "hello",
		},
	}
	replacePointers(pointers, replace)
	expect = []types.M{
		types.M{
			"objectId": "1001",
			"post":     "post",
			"key":      "hello",
		},
	}
	if reflect.DeepEqual(expect, pointers) == false {
		t.Error("expect:", expect, "result:", pointers)
	}
	/**********************************************************/
	pointers = []types.M{
		types.M{
			"objectId": "1001",
			"post":     "post",
		},
		types.M{
			"objectId": "1002",
			"post":     "post",
		},
	}
	replace = types.M{
		"1001": types.M{
			"key": "hello",
		},
		"1002": types.M{
			"key": "hello",
		},
	}
	replacePointers(pointers, replace)
	expect = []types.M{
		types.M{
			"objectId": "1001",
			"post":     "post",
			"key":      "hello",
		},
		types.M{
			"objectId": "1002",
			"post":     "post",
			"key":      "hello",
		},
	}
	if reflect.DeepEqual(expect, pointers) == false {
		t.Error("expect:", expect, "result:", pointers)
	}
}

func Test_findObjectWithKey(t *testing.T) {
	var root interface{}
	var key string
	var result types.M
	var expect types.M
	/**********************************************************/
	root = nil
	key = "post"
	result = findObjectWithKey(root, key)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	root = types.M{
		"post": "hello",
	}
	key = "post"
	result = findObjectWithKey(root, key)
	expect = types.M{
		"post": "hello",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	root = types.M{
		"key": "hello",
	}
	key = "post"
	result = findObjectWithKey(root, key)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	root = types.M{
		"key": types.M{
			"post": "hello",
		},
	}
	key = "post"
	result = findObjectWithKey(root, key)
	expect = types.M{
		"post": "hello",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	root = types.M{
		"key": types.M{
			"key": "hello",
		},
	}
	key = "post"
	result = findObjectWithKey(root, key)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	root = types.S{
		types.M{
			"post": "hello",
		},
	}
	key = "post"
	result = findObjectWithKey(root, key)
	expect = types.M{
		"post": "hello",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	root = types.S{
		types.M{
			"key": "hello",
		},
	}
	key = "post"
	result = findObjectWithKey(root, key)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/**********************************************************/
	root = types.S{
		types.M{
			"key": types.M{
				"post": "hello",
			},
		},
	}
	key = "post"
	result = findObjectWithKey(root, key)
	expect = types.M{
		"post": "hello",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_transformSelect(t *testing.T) {
	var selectObject types.M
	var key string
	var objects []types.M
	var expect types.M
	/**********************************************************/
	selectObject = nil
	key = "user"
	objects = nil
	transformSelect(selectObject, key, objects)
	expect = nil
	if reflect.DeepEqual(expect, selectObject) == false {
		t.Error("expect:", expect, "result:", selectObject)
	}
	/**********************************************************/
	selectObject = types.M{}
	key = "user"
	objects = nil
	transformSelect(selectObject, key, objects)
	expect = types.M{}
	if reflect.DeepEqual(expect, selectObject) == false {
		t.Error("expect:", expect, "result:", selectObject)
	}
	/**********************************************************/
	selectObject = types.M{
		"$select": "string",
	}
	key = "user"
	objects = nil
	transformSelect(selectObject, key, objects)
	expect = types.M{
		"$in": types.S{},
	}
	if reflect.DeepEqual(expect, selectObject) == false {
		t.Error("expect:", expect, "result:", selectObject)
	}
	/**********************************************************/
	selectObject = types.M{
		"$select": "string",
	}
	key = "user"
	objects = []types.M{}
	transformSelect(selectObject, key, objects)
	expect = types.M{
		"$in": types.S{},
	}
	if reflect.DeepEqual(expect, selectObject) == false {
		t.Error("expect:", expect, "result:", selectObject)
	}
	/**********************************************************/
	selectObject = types.M{
		"$select": "string",
	}
	key = "user"
	objects = []types.M{
		types.M{
			"user": "1001",
		},
		types.M{
			"key": "1002",
		},
	}
	transformSelect(selectObject, key, objects)
	expect = types.M{
		"$in": types.S{
			"1001",
		},
	}
	if reflect.DeepEqual(expect, selectObject) == false {
		t.Error("expect:", expect, "result:", selectObject)
	}
	/**********************************************************/
	selectObject = types.M{
		"$select": "string",
	}
	key = "user"
	objects = []types.M{
		types.M{
			"user": "1001",
		},
		types.M{
			"user": "1002",
		},
	}
	transformSelect(selectObject, key, objects)
	expect = types.M{
		"$in": types.S{
			"1001",
			"1002",
		},
	}
	if reflect.DeepEqual(expect, selectObject) == false {
		t.Error("expect:", expect, "result:", selectObject)
	}
	/**********************************************************/
	selectObject = types.M{
		"$select": "string",
		"$in": types.S{
			"1003",
		},
	}
	key = "user"
	objects = []types.M{
		types.M{
			"user": "1001",
		},
		types.M{
			"user": "1002",
		},
	}
	transformSelect(selectObject, key, objects)
	expect = types.M{
		"$in": types.S{
			"1003",
			"1001",
			"1002",
		},
	}
	if reflect.DeepEqual(expect, selectObject) == false {
		t.Error("expect:", expect, "result:", selectObject)
	}
}

func Test_transformDontSelect(t *testing.T) {
	var dontSelectObject types.M
	var key string
	var objects []types.M
	var expect types.M
	/**********************************************************/
	dontSelectObject = nil
	key = "user"
	objects = nil
	transformDontSelect(dontSelectObject, key, objects)
	expect = nil
	if reflect.DeepEqual(expect, dontSelectObject) == false {
		t.Error("expect:", expect, "result:", dontSelectObject)
	}
	/**********************************************************/
	dontSelectObject = types.M{}
	key = "user"
	objects = nil
	transformDontSelect(dontSelectObject, key, objects)
	expect = types.M{}
	if reflect.DeepEqual(expect, dontSelectObject) == false {
		t.Error("expect:", expect, "result:", dontSelectObject)
	}
	/**********************************************************/
	dontSelectObject = types.M{
		"$dontSelect": "string",
	}
	key = "user"
	objects = nil
	transformDontSelect(dontSelectObject, key, objects)
	expect = types.M{
		"$nin": types.S{},
	}
	if reflect.DeepEqual(expect, dontSelectObject) == false {
		t.Error("expect:", expect, "result:", dontSelectObject)
	}
	/**********************************************************/
	dontSelectObject = types.M{
		"$dontSelect": "string",
	}
	key = "user"
	objects = []types.M{}
	transformDontSelect(dontSelectObject, key, objects)
	expect = types.M{
		"$nin": types.S{},
	}
	if reflect.DeepEqual(expect, dontSelectObject) == false {
		t.Error("expect:", expect, "result:", dontSelectObject)
	}
	/**********************************************************/
	dontSelectObject = types.M{
		"$dontSelect": "string",
	}
	key = "user"
	objects = []types.M{
		types.M{
			"user": "1001",
		},
		types.M{
			"key": "1002",
		},
	}
	transformDontSelect(dontSelectObject, key, objects)
	expect = types.M{
		"$nin": types.S{
			"1001",
		},
	}
	if reflect.DeepEqual(expect, dontSelectObject) == false {
		t.Error("expect:", expect, "result:", dontSelectObject)
	}
	/**********************************************************/
	dontSelectObject = types.M{
		"$dontSelect": "string",
	}
	key = "user"
	objects = []types.M{
		types.M{
			"user": "1001",
		},
		types.M{
			"user": "1002",
		},
	}
	transformDontSelect(dontSelectObject, key, objects)
	expect = types.M{
		"$nin": types.S{
			"1001",
			"1002",
		},
	}
	if reflect.DeepEqual(expect, dontSelectObject) == false {
		t.Error("expect:", expect, "result:", dontSelectObject)
	}
	/**********************************************************/
	dontSelectObject = types.M{
		"$dontSelect": "string",
		"$nin": types.S{
			"1003",
		},
	}
	key = "user"
	objects = []types.M{
		types.M{
			"user": "1001",
		},
		types.M{
			"user": "1002",
		},
	}
	transformDontSelect(dontSelectObject, key, objects)
	expect = types.M{
		"$nin": types.S{
			"1003",
			"1001",
			"1002",
		},
	}
	if reflect.DeepEqual(expect, dontSelectObject) == false {
		t.Error("expect:", expect, "result:", dontSelectObject)
	}
}

func Test_transformInQuery(t *testing.T) {
	var inQueryObject types.M
	var className string
	var results []types.M
	var expect types.M
	/**********************************************************/
	inQueryObject = nil
	className = "user"
	results = nil
	transformInQuery(inQueryObject, className, results)
	expect = nil
	if reflect.DeepEqual(expect, inQueryObject) == false {
		t.Error("expect:", expect, "result:", inQueryObject)
	}
	/**********************************************************/
	inQueryObject = types.M{}
	className = "user"
	results = nil
	transformInQuery(inQueryObject, className, results)
	expect = types.M{}
	if reflect.DeepEqual(expect, inQueryObject) == false {
		t.Error("expect:", expect, "result:", inQueryObject)
	}
	/**********************************************************/
	inQueryObject = types.M{
		"$inQuery": "string",
	}
	className = "user"
	results = nil
	transformInQuery(inQueryObject, className, results)
	expect = types.M{
		"$in": types.S{},
	}
	if reflect.DeepEqual(expect, inQueryObject) == false {
		t.Error("expect:", expect, "result:", inQueryObject)
	}
	/**********************************************************/
	inQueryObject = types.M{
		"$inQuery": "string",
	}
	className = "user"
	results = []types.M{}
	transformInQuery(inQueryObject, className, results)
	expect = types.M{
		"$in": types.S{},
	}
	if reflect.DeepEqual(expect, inQueryObject) == false {
		t.Error("expect:", expect, "result:", inQueryObject)
	}
	/**********************************************************/
	inQueryObject = types.M{
		"$inQuery": "string",
	}
	className = "user"
	results = []types.M{
		types.M{
			"objectId": "1001",
		},
		types.M{
			"key": "1002",
		},
	}
	transformInQuery(inQueryObject, className, results)
	expect = types.M{
		"$in": types.S{
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1001",
			},
		},
	}
	if reflect.DeepEqual(expect, inQueryObject) == false {
		t.Error("expect:", expect, "result:", inQueryObject)
	}
	/**********************************************************/
	inQueryObject = types.M{
		"$inQuery": "string",
	}
	className = "user"
	results = []types.M{
		types.M{
			"objectId": "1001",
		},
		types.M{
			"objectId": "1002",
		},
	}
	transformInQuery(inQueryObject, className, results)
	expect = types.M{
		"$in": types.S{
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1001",
			},
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1002",
			},
		},
	}
	if reflect.DeepEqual(expect, inQueryObject) == false {
		t.Error("expect:", expect, "result:", inQueryObject)
	}
	/**********************************************************/
	inQueryObject = types.M{
		"$inQuery": "string",
		"$in": types.S{
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1003",
			},
		},
	}
	className = "user"
	results = []types.M{
		types.M{
			"objectId": "1001",
		},
		types.M{
			"objectId": "1002",
		},
	}
	transformInQuery(inQueryObject, className, results)
	expect = types.M{
		"$in": types.S{
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1003",
			},
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1001",
			},
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1002",
			},
		},
	}
	if reflect.DeepEqual(expect, inQueryObject) == false {
		t.Error("expect:", expect, "result:", inQueryObject)
	}
}

func Test_transformNotInQuery(t *testing.T) {
	var notInQueryObject types.M
	var className string
	var results []types.M
	var expect types.M
	/**********************************************************/
	notInQueryObject = nil
	className = "user"
	results = nil
	transformNotInQuery(notInQueryObject, className, results)
	expect = nil
	if reflect.DeepEqual(expect, notInQueryObject) == false {
		t.Error("expect:", expect, "result:", notInQueryObject)
	}
	/**********************************************************/
	notInQueryObject = types.M{}
	className = "user"
	results = nil
	transformNotInQuery(notInQueryObject, className, results)
	expect = types.M{}
	if reflect.DeepEqual(expect, notInQueryObject) == false {
		t.Error("expect:", expect, "result:", notInQueryObject)
	}
	/**********************************************************/
	notInQueryObject = types.M{
		"$notInQuery": "string",
	}
	className = "user"
	results = nil
	transformNotInQuery(notInQueryObject, className, results)
	expect = types.M{
		"$nin": types.S{},
	}
	if reflect.DeepEqual(expect, notInQueryObject) == false {
		t.Error("expect:", expect, "result:", notInQueryObject)
	}
	/**********************************************************/
	notInQueryObject = types.M{
		"$notInQuery": "string",
	}
	className = "user"
	results = []types.M{}
	transformNotInQuery(notInQueryObject, className, results)
	expect = types.M{
		"$nin": types.S{},
	}
	if reflect.DeepEqual(expect, notInQueryObject) == false {
		t.Error("expect:", expect, "result:", notInQueryObject)
	}
	/**********************************************************/
	notInQueryObject = types.M{
		"$notInQuery": "string",
	}
	className = "user"
	results = []types.M{
		types.M{
			"objectId": "1001",
		},
		types.M{
			"key": "1002",
		},
	}
	transformNotInQuery(notInQueryObject, className, results)
	expect = types.M{
		"$nin": types.S{
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1001",
			},
		},
	}
	if reflect.DeepEqual(expect, notInQueryObject) == false {
		t.Error("expect:", expect, "result:", notInQueryObject)
	}
	/**********************************************************/
	notInQueryObject = types.M{
		"$notInQuery": "string",
	}
	className = "user"
	results = []types.M{
		types.M{
			"objectId": "1001",
		},
		types.M{
			"objectId": "1002",
		},
	}
	transformNotInQuery(notInQueryObject, className, results)
	expect = types.M{
		"$nin": types.S{
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1001",
			},
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1002",
			},
		},
	}
	if reflect.DeepEqual(expect, notInQueryObject) == false {
		t.Error("expect:", expect, "result:", notInQueryObject)
	}
	/**********************************************************/
	notInQueryObject = types.M{
		"$notInQuery": "string",
		"$nin": types.S{
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1003",
			},
		},
	}
	className = "user"
	results = []types.M{
		types.M{
			"objectId": "1001",
		},
		types.M{
			"objectId": "1002",
		},
	}
	transformNotInQuery(notInQueryObject, className, results)
	expect = types.M{
		"$nin": types.S{
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1003",
			},
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1001",
			},
			types.M{
				"__type":    "Pointer",
				"className": "user",
				"objectId":  "1002",
			},
		},
	}
	if reflect.DeepEqual(expect, notInQueryObject) == false {
		t.Error("expect:", expect, "result:", notInQueryObject)
	}
}
