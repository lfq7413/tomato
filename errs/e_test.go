package errs

import "testing"
import "errors"

func TestE(t *testing.T) {
	err := E(20, "hello")
	if err.Error() != `{"code": 20,"error": "hello"}` {
		t.Error("err :", err.Error())
	}
}

func Test_GetErrorCode(t *testing.T) {
	var err error
	var result int
	var expect int
	err = E(20, "hello")
	result = GetErrorCode(err)
	expect = 20
	if result != expect {
		t.Error("expect:", expect, "reslut:", result)
	}

	err = errors.New("hello")
	result = GetErrorCode(err)
	expect = 0
	if result != expect {
		t.Error("expect:", expect, "reslut:", result)
	}
}

func Test_GetErrorMessage(t *testing.T) {
	var err error
	var result string
	var expect string
	err = E(20, "hello")
	result = GetErrorMessage(err)
	expect = "hello"
	if result != expect {
		t.Error("expect:", expect, "reslut:", result)
	}

	err = errors.New("hello")
	result = GetErrorMessage(err)
	expect = "hello"
	if result != expect {
		t.Error("expect:", expect, "reslut:", result)
	}
}
