package client

import "testing"

func Test_satisfies(t *testing.T) {
	var clientVersion, compatiblityVersion string
	var result bool
	var expect bool
	/******************************************************/
	clientVersion = "3.5"
	compatiblityVersion = ">=1.2.3"
	result = satisfies(clientVersion, compatiblityVersion)
	expect = true
	if expect != result {
		t.Error(clientVersion, compatiblityVersion, "expect:", expect, "result:", result)
	}
	/******************************************************/
	clientVersion = "1.2"
	compatiblityVersion = ">=3"
	result = satisfies(clientVersion, compatiblityVersion)
	expect = false
	if expect != result {
		t.Error(clientVersion, compatiblityVersion, "expect:", expect, "result:", result)
	}
	/******************************************************/
	clientVersion = "1.2"
	compatiblityVersion = ">=1.2"
	result = satisfies(clientVersion, compatiblityVersion)
	expect = true
	if expect != result {
		t.Error(clientVersion, compatiblityVersion, "expect:", expect, "result:", result)
	}
	/******************************************************/
	clientVersion = "3.5"
	compatiblityVersion = "<=1.2.3"
	result = satisfies(clientVersion, compatiblityVersion)
	expect = false
	if expect != result {
		t.Error(clientVersion, compatiblityVersion, "expect:", expect, "result:", result)
	}
	/******************************************************/
	clientVersion = "1.2"
	compatiblityVersion = "<=3"
	result = satisfies(clientVersion, compatiblityVersion)
	expect = true
	if expect != result {
		t.Error(clientVersion, compatiblityVersion, "expect:", expect, "result:", result)
	}
	/******************************************************/
	clientVersion = "1.2"
	compatiblityVersion = "<=1.2"
	result = satisfies(clientVersion, compatiblityVersion)
	expect = true
	if expect != result {
		t.Error(clientVersion, compatiblityVersion, "expect:", expect, "result:", result)
	}
	/******************************************************/
	clientVersion = "3.5"
	compatiblityVersion = ">1.2.3"
	result = satisfies(clientVersion, compatiblityVersion)
	expect = true
	if expect != result {
		t.Error(clientVersion, compatiblityVersion, "expect:", expect, "result:", result)
	}
	/******************************************************/
	clientVersion = "1.2"
	compatiblityVersion = ">3"
	result = satisfies(clientVersion, compatiblityVersion)
	expect = false
	if expect != result {
		t.Error(clientVersion, compatiblityVersion, "expect:", expect, "result:", result)
	}
	/******************************************************/
	clientVersion = "1.2"
	compatiblityVersion = ">1.2"
	result = satisfies(clientVersion, compatiblityVersion)
	expect = false
	if expect != result {
		t.Error(clientVersion, compatiblityVersion, "expect:", expect, "result:", result)
	}
	/******************************************************/
	clientVersion = "3.5"
	compatiblityVersion = "<1.2.3"
	result = satisfies(clientVersion, compatiblityVersion)
	expect = false
	if expect != result {
		t.Error(clientVersion, compatiblityVersion, "expect:", expect, "result:", result)
	}
	/******************************************************/
	clientVersion = "1.2"
	compatiblityVersion = "<3"
	result = satisfies(clientVersion, compatiblityVersion)
	expect = true
	if expect != result {
		t.Error(clientVersion, compatiblityVersion, "expect:", expect, "result:", result)
	}
	/******************************************************/
	clientVersion = "1.2"
	compatiblityVersion = "<1.2"
	result = satisfies(clientVersion, compatiblityVersion)
	expect = false
	if expect != result {
		t.Error(clientVersion, compatiblityVersion, "expect:", expect, "result:", result)
	}
	/******************************************************/
	clientVersion = "3.5"
	compatiblityVersion = "=1.2.3"
	result = satisfies(clientVersion, compatiblityVersion)
	expect = false
	if expect != result {
		t.Error(clientVersion, compatiblityVersion, "expect:", expect, "result:", result)
	}
	/******************************************************/
	clientVersion = "1.2"
	compatiblityVersion = "=3"
	result = satisfies(clientVersion, compatiblityVersion)
	expect = false
	if expect != result {
		t.Error(clientVersion, compatiblityVersion, "expect:", expect, "result:", result)
	}
	/******************************************************/
	clientVersion = "1.2"
	compatiblityVersion = "=1.2"
	result = satisfies(clientVersion, compatiblityVersion)
	expect = true
	if expect != result {
		t.Error(clientVersion, compatiblityVersion, "expect:", expect, "result:", result)
	}
	/******************************************************/
	clientVersion = "3.5"
	compatiblityVersion = "1.2.3"
	result = satisfies(clientVersion, compatiblityVersion)
	expect = false
	if expect != result {
		t.Error(clientVersion, compatiblityVersion, "expect:", expect, "result:", result)
	}
	/******************************************************/
	clientVersion = "1.2"
	compatiblityVersion = "3"
	result = satisfies(clientVersion, compatiblityVersion)
	expect = false
	if expect != result {
		t.Error(clientVersion, compatiblityVersion, "expect:", expect, "result:", result)
	}
	/******************************************************/
	clientVersion = "1.2"
	compatiblityVersion = "1.2"
	result = satisfies(clientVersion, compatiblityVersion)
	expect = true
	if expect != result {
		t.Error(clientVersion, compatiblityVersion, "expect:", expect, "result:", result)
	}
}
