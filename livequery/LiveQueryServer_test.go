package livequery

import (
	"reflect"
	"testing"

	tp "github.com/lfq7413/tomato/livequery/t"
)

func Test_getReadAccess(t *testing.T) {
	data := []struct {
		acl    tp.M
		id     string
		expect bool
	}{
		{
			acl:    nil,
			id:     "1024",
			expect: true,
		},
		{
			acl:    tp.M{},
			id:     "1024",
			expect: false,
		},
		{
			acl: tp.M{
				"*": map[string]interface{}{
					"read": true,
				},
			},
			id:     "1024",
			expect: false,
		},
		{
			acl: tp.M{
				"*": map[string]interface{}{
					"read": true,
				},
			},
			id:     "*",
			expect: true,
		},
		{
			acl: tp.M{
				"*": map[string]interface{}{
					"read": true,
				},
				"1024": map[string]interface{}{
					"read": true,
				},
			},
			id:     "1024",
			expect: true,
		},
		{
			acl: tp.M{
				"*": map[string]interface{}{
					"read": true,
				},
				"1024": map[string]interface{}{
					"read": true,
				},
			},
			id:     "2048",
			expect: false,
		},
	}

	for _, d := range data {
		result := getReadAccess(d.acl, d.id)
		if reflect.DeepEqual(d.expect, result) == false {
			t.Error("expect:", d.expect, "result:", result)
		}
	}
}
