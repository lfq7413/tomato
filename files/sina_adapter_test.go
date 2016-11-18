package files

import "testing"

func Test_sina(t *testing.T) {
	f := newSinaAdapter()
	hello := "hello world!"
	err := f.createFile("hello-test.txt", []byte(hello), "text/plain")
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}

	err = f.deleteFile("hello.txt")
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
}
