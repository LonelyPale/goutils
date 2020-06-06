package uuid

import "testing"

func TestNew(t *testing.T) {
	uuid := New().String()
	t.Log(uuid)
}
