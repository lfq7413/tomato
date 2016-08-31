package files

import (
	"reflect"
	"testing"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/types"
)

func Test_FileAdapter(t *testing.T) {
	adapter = newFileSystemAdapter("1001")
	hello := "hello world!"
	resp := CreateFile("hellol.txt", []byte(hello), "text/plain")
	if resp["url"] == "" || resp["name"] == "" {
		t.Error("expect:", "url+name", "result:", resp)
	}

	data, _ := GetFileData(resp["name"])
	if hello != string(data) {
		t.Error("expect:", hello, "result:", string(data))
	}

	err := DeleteFile(resp["name"])
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}

	adapter = newGridStoreAdapter()
	hello = "hello world!"
	resp = CreateFile("hellol.txt", []byte(hello), "text/plain")
	if resp["url"] == "" || resp["name"] == "" {
		t.Error("expect:", "url+name", "result:", resp)
	}

	data, _ = GetFileData(resp["name"])
	if hello != string(data) {
		t.Error("expect:", hello, "result:", string(data))
	}

	err = DeleteFile(resp["name"])
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
}

func Test_ExpandFilesInObject(t *testing.T) {
	var object, expect interface{}
	config.TConfig = &config.Config{
		ServerURL: "http://127.0.0.1",
		AppID:     "1001",
	}
	/*************************************************************/
	object = types.M{
		"file": types.M{
			"__type": "File",
			"name":   "hello.txt",
		},
	}
	ExpandFilesInObject(object)
	expect = types.M{
		"file": types.M{
			"__type": "File",
			"name":   "hello.txt",
			"url":    "http://127.0.0.1/files/1001/hello.txt",
		},
	}
	if reflect.DeepEqual(expect, object) == false {
		t.Error("expect:", expect, "result:", object)
	}
	/*************************************************************/
	object = types.M{
		"file": types.M{
			"__type": "File",
			"name":   "hello.txt",
			"url":    "http://127.0.0.1/files/1001/hello.txt",
		},
	}
	ExpandFilesInObject(object)
	expect = types.M{
		"file": types.M{
			"__type": "File",
			"name":   "hello.txt",
			"url":    "http://127.0.0.1/files/1001/hello.txt",
		},
	}
	if reflect.DeepEqual(expect, object) == false {
		t.Error("expect:", expect, "result:", object)
	}
	/*************************************************************/
	object = types.S{
		types.M{
			"file": types.M{
				"__type": "File",
				"name":   "hello.txt",
			},
		},
		types.M{
			"file": types.M{
				"__type": "File",
				"name":   "world.txt",
			},
		},
	}
	ExpandFilesInObject(object)
	expect = types.S{
		types.M{
			"file": types.M{
				"__type": "File",
				"name":   "hello.txt",
				"url":    "http://127.0.0.1/files/1001/hello.txt",
			},
		},
		types.M{
			"file": types.M{
				"__type": "File",
				"name":   "world.txt",
				"url":    "http://127.0.0.1/files/1001/world.txt",
			},
		},
	}
	if reflect.DeepEqual(expect, object) == false {
		t.Error("expect:", expect, "result:", object)
	}
}
