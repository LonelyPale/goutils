package goutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMergeSliceByte(t *testing.T) {
	a := []byte{1, 2, 3}
	b := []byte{4, 5, 6}
	c := []byte{7, 8, 9}
	d := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	var bs []byte

	assert.Equal(t, bs, MergeSliceByte())
	assert.Equal(t, a, MergeSliceByte(a))
	assert.Equal(t, a, MergeSliceByte(a, nil))
	assert.Equal(t, d, MergeSliceByte(a, b, c))

	var count1 time.Duration
	for i := 0; i < 1000000; i++ {
		t1 := time.Now()
		mergeSliceCopy(a, b, c)
		count1 += time.Since(t1)
	}
	t.Log("copy count1: ", count1/1000000)

	var count2 time.Duration
	for i := 0; i < 1000000; i++ {
		t1 := time.Now()
		mergeSliceAppend(a, b, c)
		count2 += time.Since(t1)
	}
	t.Log("append count2: ", count2/1000000)
}

func TestReverseByte(t *testing.T) {
	a := []byte{1, 2, 3}
	ReverseByte(a)
	assert.Equal(t, []byte{3, 2, 1}, a)
}
