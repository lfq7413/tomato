package server

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_getUser(t *testing.T) {
	TomatoInfo = map[string]string{
		"appId":     "test",
		"clientKey": "test",
		"serverURL": "http://127.0.0.1:8080/v1",
	}
	sessionToken := "59C3697129E46DE6F1CED31B8FB2B862"
	user, err := getUser(sessionToken)
	if err != nil {
		t.Error(err)
	}
	if user == nil {
		t.Error("user is null")
	}
	if reflect.DeepEqual(user["objectId"], "57d7c2013cdd0164775cea4f") == false {
		fmt.Println(user["objectId"], err)
		t.Error("expect:", "57d7c2013cdd0164775cea4f", "result:", user["objectId"])
	}
}
