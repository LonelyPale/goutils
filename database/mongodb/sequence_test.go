package mongodb

import (
	"sync"
	"testing"
)

func TestGetSequence(t *testing.T) {
	name := "test"
	n, err := GetSequence(name)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Sequence: test ->", n)
}

func TestGetSequenceGo(t *testing.T) {
	var group sync.WaitGroup
	for i := 0; i < 10; i++ {
		group.Add(1)
		go func(i int) {
			name := "test"
			n, err := GetSequence(name)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("Sequence ", i, ": test ->", n)
			group.Done()
		}(i)
	}
	group.Wait()
}
