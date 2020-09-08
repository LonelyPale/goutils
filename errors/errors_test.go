package errors

import "testing"

func TestWrapf(t *testing.T) {
	e1 := New("Error-1")
	e2 := Wrapf(e1, "code: %d, %s.", 404, "not find")
	t.Log(e2)
}
