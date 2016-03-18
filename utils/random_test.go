package utils

import "testing"

func TestCreateObjectID(t *testing.T) {
	id := CreateObjectID()
	l := len(id)
	if l != 32 {
		t.Error("CreateObjectID len is not 32!")
	}
}
