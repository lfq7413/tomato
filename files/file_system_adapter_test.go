package files

import (
	"testing"

	"github.com/lfq7413/tomato/config"
)
import "reflect"

func Test_fileSystemAdapter(t *testing.T) {
	f := newFileSystemAdapter("tomato")
	hello := "hello world!"
	err := f.createFile("hello.txt", []byte(hello), "text/plain")
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}

	data, _ := f.getFileData("hello.txt")
	if reflect.DeepEqual(hello, string(data)) == false {
		t.Error("expect:", hello, "result:", string(data))
	}

	config.TConfig = &config.Config{
		ServerURL: "http://127.0.0.1",
		AppID:     "1001",
	}
	loc := f.getFileLocation("hello.txt")
	if loc != "http://127.0.0.1/files/1001/hello.txt" {
		t.Error("expect:", "http://127.0.0.1/files/1001/hello.txt", "result:", loc)
	}

	err = f.deleteFile("hello.txt")
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
}
