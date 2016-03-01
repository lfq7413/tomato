package utils

import "testing"

func TestCreateObjectId(t *testing.T) {
	id := CreateObjectId()
	l := len(id)
	if l != 32 {
		t.Error("CreateObjectId len is not 32!")
	}
}
