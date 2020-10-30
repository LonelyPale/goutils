package goutils

import "testing"

func TestUUID(t *testing.T) {
	id, err := UUID()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)
}
