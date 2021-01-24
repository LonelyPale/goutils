package crypto

import (
	"encoding/hex"
	"testing"
	"time"
)

func TestBcrypt(t *testing.T) {
	t1 := time.Now()
	data := []byte("hello world")
	salt := []byte("123")
	hash, err := Bcrypt(data, salt)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("elapsed: ", time.Since(t1), len(hash), hash)
	t.Log(hex.EncodeToString(hash))
}

func TestBcryptSimple(t *testing.T) {
	t1 := time.Now()
	data := []byte("hello world")
	salt := []byte("123")
	hash, err := BcryptSimple(data, salt)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("elapsed: ", time.Since(t1), len(hash), hash)
	t.Log(hex.EncodeToString(hash))
}
