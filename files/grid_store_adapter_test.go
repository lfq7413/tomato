package files

import (
	"testing"

	"gopkg.in/mgo.v2"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/storage"
)
import "reflect"

func Test_gridStoreAdapter(t *testing.T) {
	storage.TomatoDB = newMongoDB("192.168.99.100:27017/test")
	f := newGridStoreAdapter()
	hello := "hello world!"
	err := f.createFile("hello.txt", []byte(hello), "text/plain")
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}

	data := f.getFileData("hello.txt")
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

func newMongoDB(url string) *storage.Database {
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	database := session.DB("")
	db := &storage.Database{
		MongoSession:  session,
		MongoDatabase: database,
	}
	return db
}
