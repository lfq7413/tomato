package utils

import "testing"

func TestStringtoTime(t *testing.T) {
	s := "2016-02-28T13:25:05.123Z"
	time, _ := StringtoTime(s)
	s1 := TimetoString(time)
	if s != s1 {
		t.Error("StringtoTime error")
	}
}

func TestUnixmillitoTime(t *testing.T) {
	var m int64
	m = 1456637363615
	time := UnixmillitoTime(m)
	m1 := TimetoUnixmilli(time)
	if m != m1 {
		t.Error("UnixmillitoTime error")
	}
}

func TestStringtoUnixmilli(t *testing.T) {
	s := "2016-02-28T13:25:05.123Z"
	m, _ := StringtoUnixmilli(s)
	s1 := UnixmillitoString(m)
	if s != s1 {
		t.Error("StringtoUnixmilli error")
	}
}

func TestUnixmillitoString(t *testing.T) {
	var m int64
	m = 0
	s := UnixmillitoString(m)
	m1, _ := StringtoUnixmilli(s)
	if m != m1 {
		t.Error("UnixmillitoString error")
	}
}
