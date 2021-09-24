package crypto

import (
	"crypto/rand"
	"testing"
	"time"
)

func Test(t *testing.T) {
	iv := make([]byte, 12)
	for i := 0; i < 100; i++ {
		_, err := rand.Read(iv)
		if err != nil {
			t.Error(err)
		}
		t.Log(iv)
	}
}

func TestGenerateSecretKey(t *testing.T) {
	for i := 0; i < 10; i++ {
		startTime := time.Now()

		key, err := GenerateSecretKey(16)
		if err != nil {
			t.Fatal(err)
		}

		if len(key) != 16 {
			t.Fatal("gen length error")
		}

		t.Log(key.Hex(), "duration:", time.Since(startTime))
	}
}
