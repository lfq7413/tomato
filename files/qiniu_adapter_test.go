package files

import "testing"

func Test_qiniuAdapter(t *testing.T) {
	f := newQiniuAdapter()
	hello := "hello world!"
	err := f.createFile("hello-test2.txt", []byte(hello), "text/plain")
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}

	err = f.deleteFile("hello.txt")
	if err != nil {
		t.Error("expect:", nil, "result:", err)
	}
}
