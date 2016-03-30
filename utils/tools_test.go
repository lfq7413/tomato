package utils

import "testing"

func TestRegexp(t *testing.T) {
	s := "11@aa"
	s1 := "aa@cc.com"
	s2 := "ssss"
	s3 := "ssss@"
	s4 := "@ssss"
	if IsEmail(s) == false {
		t.Error(s, IsEmail(s))
	}
	if IsEmail(s1) == false {
		t.Error(s1, IsEmail(s1))
	}
	if IsEmail(s2) {
		t.Error(s2, IsEmail(s2))
	}
	if IsEmail(s3) {
		t.Error(s3, IsEmail(s3))
	}
	if IsEmail(s4) {
		t.Error(s4, IsEmail(s4))
	}
}

func TestDeepCopy(t *testing.T) {
	var src interface{}
	src = map[string]interface{}{
		"a": 1,
		"b": true,
		"c": "hello",
		"d": []string{"hello", "hi"},
		"e": map[string]interface{}{
			"f": "haha",
		},
	}
	var dst interface{}
	dst = DeepCopy(src)
	if v, ok := src.(map[string]interface{}); ok {
		e := v["e"]
		if v1, ok1 := e.(map[string]interface{}); ok1 {
			v1["g"] = "okok"
		}
	}
	if dst == nil {
		t.Error(src)
		t.Error(dst)
	}
}
