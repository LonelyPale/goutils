package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectIDFrom(t *testing.T) {
	id := NewObjectID()
	idstr := id.Hex()
	t.Log(idstr)

	objid, err := ObjectIDFrom(id)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, id, objid)

	objid, err = ObjectIDFrom(idstr)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, id, objid)

	objid, err = ObjectIDFrom(&id)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, id, objid)

	objid, err = ObjectIDFrom(&idstr)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, id, objid)
}
