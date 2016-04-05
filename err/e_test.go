package err

import "testing"

func TestE(t *testing.T) {
	err := E(20, "hello")
	if err.Error() != `{"code": 20,"error": "hello"}` {
		t.Error("err :", err.Error())
	}
}
