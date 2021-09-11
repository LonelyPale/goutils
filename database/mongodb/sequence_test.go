package mongodb

import (
	"sync"
	"testing"
)

func TestGetSequence(t *testing.T) {
	coll := client.DB().Collection("sequence")
	sequence := NewSequence(coll)

	n, err := sequence.Inc()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Sequence:", sequence)
	t.Log("Sequence: test ->", n)
}

func TestGetSequenceGo(t *testing.T) {
	coll := client.DB().Collection("sequence")
	sequence := NewSequence(coll)

	var group sync.WaitGroup
	for i := 0; i < 10; i++ {
		group.Add(1)
		go func(i int) {
			n, err := sequence.Inc()
			if err != nil {
				t.Fatal(err)
			}
			t.Log("Sequence ", i, ": test ->", n)
			group.Done()
		}(i)
	}
	group.Wait()
}
