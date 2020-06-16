package session

import "testing"

func Test(t *testing.T) {
	s, err := NewSession()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("id:", s.ID())

	s.Set("name", "tom")

	name := s.Get("name")
	t.Log(name)

	if err := s.Save(); err != nil {
		t.Fatal(err)
	}

	ss, err := NewSession(s.ID())
	if err != nil {
		t.Fatal(err)
	}
	t.Log("id2:", ss.ID())

	name2 := s.Get("name")
	t.Log(name2)

	if err := s.Save(); err != nil {
		t.Fatal(err)
	}
}
