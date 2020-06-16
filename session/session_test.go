package session

import (
	"testing"
)

func Test(t *testing.T) {
	store, err := NewMemoryStore()
	if err != nil {
		t.Fatal(err)
	}

	s := store.New()
	t.Log("id1:", s.ID())

	s.Set("name", "tom")

	name := s.Get("name")
	t.Log("name1:", name)

	if err := s.Save(); err != nil {
		t.Fatal(err)
	}

	ss, err := store.Get(s.ID())
	if err != nil {
		t.Fatal(err)
	}
	t.Log("id2:", ss.ID())

	name2 := s.Get("name")
	t.Log("name2:", name2)
}
