package utils

import "testing"

func TestPassword(t *testing.T) {
	s := Hash("pass")
	if s != "d74ff0ee8da3b9806b18c877dbf29bbde50b5bd8e4dad7a3a725000feb82e8f1" {
		t.Error("Hash error", s)
	}
}

func TestCompare(t *testing.T) {
	b := Compare("pass", "d74ff0ee8da3b9806b18c877dbf29bbde50b5bd8e4dad7a3a725000feb82e8f1")
	if b == false {
		t.Error("Compare error", b)
	}
}
