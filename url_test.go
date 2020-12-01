package goutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLJoin(t *testing.T) {
	assert.Equal(t, "http://127.0.0.1:8080/test/", URLJoin("http://", "/127.0.0.1:8080/", "/test/"))
	assert.Equal(t, "http://127.0.0.1:8080/test/", URLJoin("http://", "127.0.0.1:8080/", "test/"))
	assert.Equal(t, "http://127.0.0.1:8080/test/", URLJoin("http://", "127.0.0.1:8080", "test/"))
}
