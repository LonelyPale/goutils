package goutils

import (
	"testing"

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
	assert.Equal(t, d, MergeSliceByte(a, b, c))
}

func TestReverseByte(t *testing.T) {
	a := []byte{1, 2, 3}
	ReverseByte(&a)
	assert.Equal(t, []byte{3, 2, 1}, a)
}
