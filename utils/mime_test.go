package utils

import "testing"

func TestExtName(t *testing.T) {
	s1 := "aa.jsp"
	s2 := "aa."
	s3 := ".aa"
	s4 := "aa"
	s5 := "aa.bb.cc"
	if ExtName(s1) != "jsp" {
		t.Error(s1, ExtName(s1))
	}
	if ExtName(s2) != "" {
		t.Error(s2, ExtName(s2))
	}
	if ExtName(s3) != "aa" {
		t.Error(s3, ExtName(s3))
	}
	if ExtName(s4) != "" {
		t.Error(s4, ExtName(s4))
	}
	if ExtName(s5) != "cc" {
		t.Error(s5, ExtName(s5))
	}
}

func TestIsFileName(t *testing.T) {
	s1 := "_aa.jsp"
	s2 := "aa."
	s3 := ".aa"
	s4 := "a a.~-_@"
	s5 := "aa.bb.cc!#$%^&*()"
	if IsFileName(s1) == false {
		t.Error(s1, IsFileName(s1))
	}
	if IsFileName(s2) == false {
		t.Error(s2, IsFileName(s2))
	}
	if IsFileName(s3) == true {
		t.Error(s3, IsFileName(s3))
	}
	if IsFileName(s4) == false {
		t.Error(s4, IsFileName(s4))
	}
	if IsFileName(s5) == true {
		t.Error(s5, IsFileName(s5))
	}
}

func TestLookupContentType(t *testing.T) {
	s1 := "aa.jpg"
	s2 := "aa.pkg"
	s3 := "aa.ipa"
	s4 := "aa"
	if LookupContentType(s1) != "image/jpeg" {
		t.Error(s1, LookupContentType(s1))
	}
	if LookupContentType(s2) != "application/octet-stream" {
		t.Error(s2, LookupContentType(s2))
	}
	if LookupContentType(s3) != "application/octet-stream" {
		t.Error(s3, LookupContentType(s3))
	}
	if LookupContentType(s4) != "application/octet-stream" {
		t.Error(s4, LookupContentType(s4))
	}
}

func TestLookupExtension(t *testing.T) {
	s1 := "application/vnd.android.package-archive"
	s2 := "application/octet-stream"
	s3 := "application/xxx"
	if LookupExtension(s1) != "apk" {
		t.Error(s1, LookupExtension(s1))
	}
	if LookupExtension(s2) != "bin" {
		t.Error(s2, LookupExtension(s2))
	}
	if LookupExtension(s3) != "" {
		t.Error(s3, LookupExtension(s3))
	}
}
