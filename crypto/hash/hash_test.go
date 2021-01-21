package hash

import (
	"crypto/md5"
	"crypto/sha256"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestHash(t *testing.T) {
	data := []byte("hello world")

	hashSha256, err := Hex(sha256.New, data)
	if err != nil {
		t.Fatal(err)
	}

	hashMd5, err := Hex(md5.New, data)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9", hashSha256)
	assert.Equal(t, "5eb63bbbe01eeed093cb22bb8f5acdc3", hashMd5)
}
